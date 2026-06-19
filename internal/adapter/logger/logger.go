package logger

import (
	"log"
	"os"
)

// Logger — простая реализация логгера
type Logger struct {
	level string
}

// NewLogger создаёт новый логгер
func NewLogger(level string) *Logger {
	return &Logger{level: level}
}

func (l *Logger) Info(msg string, keysAndValues ...interface{}) {
	log.Printf("[INFO] %s %v\n", msg, keysAndValues)
}

func (l *Logger) Error(msg string, keysAndValues ...interface{}) {
	log.Printf("[ERROR] %s %v\n", msg, keysAndValues)
}

func (l *Logger) Debug(msg string, keysAndValues ...interface{}) {
	if l.level == "debug" {
		log.Printf("[DEBUG] %s %v\n", msg, keysAndValues)
	}
}

func (l *Logger) Warn(msg string, keysAndValues ...interface{}) {
	log.Printf("[WARN] %s %v\n", msg, keysAndValues)
}

func (l *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	log.Printf("[FATAL] %s %v\n", msg, keysAndValues)
	os.Exit(1)
}
