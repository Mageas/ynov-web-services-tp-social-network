package logger

import (
	"log"
	"os"
)

// Logger provides structured logging capabilities
type Logger struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

// New creates a new Logger instance
func New() *Logger {
	return &Logger{
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Info logs an informational message
func (l *Logger) Info(message string, args ...interface{}) {
	l.infoLog.Printf(message, args...)
}

// Error logs an error message
func (l *Logger) Error(message string, args ...interface{}) {
	l.errorLog.Printf(message, args...)
}

// Fatal logs a fatal error and exits
func (l *Logger) Fatal(message string, args ...interface{}) {
	l.errorLog.Fatalf(message, args...)
}
