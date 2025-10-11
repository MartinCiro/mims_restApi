package main

import (
	"fmt"
	"log"
	config "scrapper_go_email/config"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func main() {
	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	fmt.Printf("Server running on port %d (Debug: %t)\n", cfg.Port, cfg.Debug)

	// Conectar al servidor IMAP de Gmail
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Logout()

	// Iniciar sesión usando la configuración
	email := cfg.GmailEmail
	password := cfg.GmailAppPassword

	if email == "" || password == "" {
		log.Fatal("Gmail credentials not found in configuration")
	}

	if err := c.Login(email, password); err != nil {
		log.Fatal(err)
	}

	// Seleccionar bandeja de entrada
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	// Buscar mensajes no leídos
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{"\\Seen"}
	uids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	if len(uids) == 0 {
		fmt.Println("No hay mensajes no leídos")
		return
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(uids...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody}, messages)
	}()

	fmt.Printf("Mensajes no leídos: %d\n", len(uids))
	for msg := range messages {
		fmt.Printf("Asunto: %s\n", msg.Envelope.Subject)
		fmt.Printf("De: %s\n", msg.Envelope.From[0].Address())
		fmt.Printf("Fecha: %s\n", msg.Envelope.Date)
		fmt.Println("---")
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
}
