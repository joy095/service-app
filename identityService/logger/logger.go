package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggerConfig holds the configuration for logging
type LoggerConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Level      logrus.Level
}

// NewLogger initializes a new logger
func NewLogger(config LoggerConfig) *logrus.Logger {
	logger := logrus.New()

	// Configure log rotation
	logger.SetOutput(&lumberjack.Logger{
		Filename:   config.Filename,
		MaxSize:    config.MaxSize,    // MB
		MaxBackups: config.MaxBackups, // Number of old logs to keep
		MaxAge:     config.MaxAge,     // Days
		Compress:   true,              // Compress old logs
	})

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	logger.SetLevel(config.Level)
	logger.SetOutput(os.Stdout)

	return logger
}

// Initialize Loggers
var (
	InfoLogger  *logrus.Logger
	ErrorLogger *logrus.Logger
)

// InitLoggers initializes both info and error loggers
func InitLoggers() {
	InfoLogger = NewLogger(LoggerConfig{
		Filename:   "logs/info.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Level:      logrus.InfoLevel,
	})

	ErrorLogger = NewLogger(LoggerConfig{
		Filename:   "logs/error.log",
		MaxSize:    10,
		MaxBackups: 5,
		MaxAge:     30,
		Level:      logrus.ErrorLevel,
	})
}
