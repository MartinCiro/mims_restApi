package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "modernc.org/sqlite"

	wsp "api_go/handlers"
)

var (
	client *whatsmeow.Client
	wg     sync.WaitGroup
)

// Estructuras para las respuestas JSON
type QRResponse struct {
	QRCode    string `json:"qr_code,omitempty"`
	Connected bool   `json:"connected"`
	Success   bool   `json:"success,omitempty"`
	Message   string `json:"message,omitempty"`
}

// Variables globales para el QR
var (
	currentQRCode string
	isConnected   bool
	qrMutex       sync.RWMutex
)

func main() {
	// Iniciar servidor web en una goroutine
	go startWebServer()

	// Crear contexto
	ctx := context.Background()

	// Configuración del logger
	dbLog := waLog.Stdout("Database", "DEBUG", true)

	// Configuración de la base de datos (SQLite)
	container, err := sqlstore.New(ctx, "sqlite", "file:whatsapp_session.db?_pragma=foreign_keys(1)&_pragma=journal_mode=WAL&_pragma=busy_timeout=10000&_pragma=synchronous=NORMAL&_pragma=cache_size=10000", dbLog)
	if err != nil {
		log.Fatal("Error creando almacenamiento:", err)
	}

	// Si no hay dispositivos guardados, registrar uno nuevo
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		log.Fatal("Error obteniendo dispositivo:", err)
	}

	// Configuración del cliente
	clientLog := waLog.Stdout("Client", "INFO", true)
	client = whatsmeow.NewClient(deviceStore, clientLog)

	// Registrar handlers de eventos
	wsp.RegisterWhatsAppHandlers(client)

	// Conectar cliente
	if client.Store.ID == nil {
		// No hay sesión guardada, necesitamos iniciar sesión
		fmt.Println("🚀 Iniciando nueva sesión de WhatsApp...")

		// Obtener canal QR
		qrChan, err := client.GetQRChannel(ctx)
		if err != nil {
			log.Fatal("Error obteniendo canal QR:", err)
		}

		// Conectar el cliente
		err = client.Connect()
		if err != nil {
			log.Fatal("Error conectando:", err)
		}

		// Manejar el QR tanto en terminal como en web
		go handleQRChannel(qrChan)

	} else {
		// Sesión ya guardada, conectar directamente
		fmt.Println("🔗 Reconectando a WhatsApp...")
		err = client.Connect()
		if err != nil {
			log.Fatal("Error conectando:", err)
		}
		fmt.Println("✅ Reconectado a WhatsApp exitosamente")
		setConnected(true)
	}

	// Verificar conexión
	if client.IsConnected() {
		fmt.Println("✅ WhatsApp conectado y listo")
		setConnected(true)
	} else {
		fmt.Println("⚠️ Cliente no conectado")
		os.Exit(1)
	}

	// Esperar señal para terminar
	fmt.Println("🟢 Bot ejecutándose. Presiona Ctrl+C para salir")
	fmt.Println("🌐 Servidor web corriendo en http://localhost:8080")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n🔴 Desconectando...")
	client.Disconnect()
	wg.Wait()
	fmt.Println("👋 Sesión finalizada")
}

func handleQRChannel(qrChan <-chan whatsmeow.QRChannelItem) {
	for evt := range qrChan {
		if evt.Event == "code" {
			// Mostrar QR en terminal
			fmt.Println("📱 Escanea el código QR que aparece a continuación con WhatsApp:")
			fmt.Println("1. Abre WhatsApp en tu teléfono")
			fmt.Println("2. Ve a Ajustes → WhatsApp Web → Escanear código")
			fmt.Println("3. Escanea el código QR de abajo\n")
			fmt.Println(evt.Code)
			fmt.Println("\n⏰ Esperando escaneo...")

			// Generar QR como imagen para la web
			setQRCode(evt.Code)

		} else if evt.Event == "success" {
			fmt.Println("✅ ¡Sesión iniciada correctamente!")
			setConnected(true)
			clearQRCode()
			break
		} else if evt.Event == "timeout" {
			fmt.Println("❌ Tiempo agotado. Reinicia la aplicación.")
			setQRCode("") // Limpiar QR
			os.Exit(1)
		}
	}
}

// Funciones para manejar el estado del QR
func setQRCode(code string) {
	qrMutex.Lock()
	defer qrMutex.Unlock()
	currentQRCode = code
	isConnected = false
}

func clearQRCode() {
	qrMutex.Lock()
	defer qrMutex.Unlock()
	currentQRCode = ""
}

func setConnected(connected bool) {
	qrMutex.Lock()
	defer qrMutex.Unlock()
	isConnected = connected
	if connected {
		currentQRCode = ""
	}
}

func getQRStatus() (string, bool) {
	qrMutex.RLock()
	defer qrMutex.RUnlock()
	return currentQRCode, isConnected
}

// Servidor Web
func startWebServer() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/qr-data", qrDataHandler)
	http.HandleFunc("/refresh-qr", refreshQRHandler)

	fmt.Println("🌐 Iniciando servidor web en puerto 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
		<html>
		<head>
			<title>WhatsApp QR Code</title>
			<meta charset="utf-8">
		</head>
		<body>
			<h1>WhatsApp QR Code</h1>
			<p>1. Abre WhatsApp en tu teléfono</p>
			<p>2. Ve a Ajustes → WhatsApp Web → Escanear código</p>
			<p>3. Escanea el código QR de abajo</p>
			
			<div id="qrContainer">
				<p id="qrStatus">Cargando código QR...</p>
			</div>
			
			<button onclick="refreshQR()">Actualizar QR</button>
			<p id="message"></p>

			<script>
				function loadQR() {
					fetch('/qr-data')
						.then(response => response.json())
						.then(data => {
							const container = document.getElementById('qrContainer');
							const status = document.getElementById('qrStatus');
							const message = document.getElementById('message');
							
							if (data.qr_code) {
								container.innerHTML = '<img src="' + data.qr_code + '" alt="QR Code">';
								status.textContent = 'Código QR listo para escanear';
								message.textContent = '';
							} else if (data.connected) {
								container.innerHTML = '<div style="color: green; font-size: 24px;">✅ WhatsApp Conectado</div>';
								status.textContent = '';
								message.textContent = '✅ Sesión iniciada correctamente. Puedes cerrar esta página.';
							} else {
								status.textContent = 'Esperando código QR...';
								message.textContent = '';
								setTimeout(loadQR, 2000);
							}
						})
						.catch(error => {
							console.error('Error:', error);
							setTimeout(loadQR, 2000);
						});
				}

				function refreshQR() {
					fetch('/refresh-qr', { method: 'POST' })
						.then(response => response.json())
						.then(data => {
							if (data.success) {
								document.getElementById('qrStatus').textContent = 'Solicitando nuevo QR...';
								setTimeout(loadQR, 1000);
							}
						});
				}

				// Cargar QR inicialmente
				loadQR();
				
				// Verificar estado cada 3 segundos
				setInterval(loadQR, 3000);
			</script>
		</body>
		</html>`

	tmpl = strings.ReplaceAll(tmpl, "&#34;", "\"")
	tmpl = strings.ReplaceAll(tmpl, "&#39;", "'")

	t, err := template.New("index").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, nil)
}

func qrDataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	qrCode, connected := getQRStatus()

	if connected {
		json.NewEncoder(w).Encode(QRResponse{
			Connected: true,
			Message:   "WhatsApp conectado exitosamente",
		})
		return
	}

	if qrCode == "" {
		json.NewEncoder(w).Encode(QRResponse{
			Connected: false,
			Message:   "Esperando código QR...",
		})
		return
	}

	// Generar código QR como Data URL
	qrPNG, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		json.NewEncoder(w).Encode(QRResponse{
			Connected: false,
			Message:   "Error generando QR: " + err.Error(),
		})
		return
	}

	// Convertir a Data URL
	qrBase64 := base64.StdEncoding.EncodeToString(qrPNG)
	qrDataURL := "data:image/png;base64," + qrBase64

	json.NewEncoder(w).Encode(QRResponse{
		QRCode:    qrDataURL,
		Connected: false,
		Message:   "Código QR listo para escanear",
	})
}

func refreshQRHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Limpiar el QR actual para forzar regeneración
	clearQRCode()

	json.NewEncoder(w).Encode(QRResponse{
		Success: true,
		Message: "Solicitando nuevo código QR...",
	})
}
