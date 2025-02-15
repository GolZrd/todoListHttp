package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Создание логгера для запросов
func NewRequestLogger(filepath string) *logrus.Logger {
	// создаем логгер
	logger := logrus.New()

	// Проверяем существование папки для логов, если ее нет создаем
	if err := os.MkdirAll("logs", 0755); err != nil {
		logrus.Fatalf("Failed to create logs directory: %v", err)
	}

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to create request log file: %v", err)
	}

	logger.Out = file
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.InfoLevel

	return logger
}

// Создание логгера для ответов
func NewResponseLogger(filepath string) *logrus.Logger {
	logger := logrus.New()
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to create response log file: %v", err)
	}

	logger.Out = file
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.InfoLevel

	return logger
}

// Создание логгера для ошибок
func NewErrorLogger(filepath string) *logrus.Logger {
	logger := logrus.New()
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to create error log file: %v", err)
	}

	logger.Out = file
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.ErrorLevel

	return logger
}
