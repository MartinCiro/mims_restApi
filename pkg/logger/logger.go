package logger

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Logger es un wrapper alrededor de slog.Logger
type Logger struct {
	*slog.Logger
}

// Configuración global
var globalLogger *Logger

// Inicialización por defecto
func init() {
	// Por defecto usamos texto legible para desarrollo
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Simplificar el timestamp para desarrollo
			if a.Key == slog.TimeKey {
				if t, ok := a.Value.Any().(time.Time); ok {
					a.Value = slog.StringValue(t.Format("15:04:05"))
				}
			}
			return a
		},
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	globalLogger = &Logger{slog.New(handler)}
}

// SetupDevelopment configura el logger para desarrollo
func SetupDevelopment() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	globalLogger = &Logger{slog.New(handler)}
}

// SetupProduction configura el logger para producción (JSON)
func SetupProduction() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	globalLogger = &Logger{slog.New(handler)}
}

// SetLevel cambia el nivel de log globalmente
func SetLevel(level slog.Level) {
	var handler slog.Handler

	if globalLogger.Handler().Enabled(context.Background(), slog.LevelDebug) {
		// Si estaba en modo texto, mantener texto
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		// Si estaba en modo JSON, mantener JSON
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	globalLogger = &Logger{slog.New(handler)}
}

// Helper para agregar información del caller
func withCaller(logger *Logger) *Logger {
	_, file, line, _ := runtime.Caller(2)
	return &Logger{logger.With(
		slog.String("file", file),
		slog.Int("line", line),
	)}
}

// Métodos de nivel con caller information

func Debug(msg string, args ...interface{}) {
	withCaller(globalLogger).Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	withCaller(globalLogger).Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	withCaller(globalLogger).Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	withCaller(globalLogger).Error(msg, args...)
}

func Fatal(msg string, args ...interface{}) {
	withCaller(globalLogger).Error(msg, args...)
	os.Exit(1)
}

// Métodos para crear loggers con contexto

func With(args ...interface{}) *Logger {
	return &Logger{globalLogger.Logger.With(args...)}
}

func WithGroup(name string) *Logger {
	return &Logger{globalLogger.Logger.WithGroup(name)}
}
