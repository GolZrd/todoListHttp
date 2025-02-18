package kafka

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

// Создадим структуру для консьюмера
type Consumer struct {
	consumer *kafka.Consumer
	logger   *logrus.Logger
	stop     bool // Флаг для остановки консьюмера
}

// Создадим конструктор для консьюмера, также как и в продюсере, сюда передаем адреса наших брокеров
func NewConcumer(adress []string, topic string, consumerGroup string, logfile string) (*Consumer, error) {
	// Создаем конфиг для консьюмера
	cfg := &kafka.ConfigMap{
		"bootstrap.servers":        adress,        // Список брокеров
		"group.id":                 consumerGroup, // Группа, которая будет обрабатывать сообщения
		"session.timeout.ms":       6000,          // Таймаут сессии, в миллисекундах
		"enable.auto.offset.store": false,         // Отключаем автоматическое сохранение оффсета
		"enable.auto.commit":       true,          // Включаем автоматическое обновление оффсета
		"auto.commit.interval.ms":  5000,          // Интервал обновления оффсета, в миллисекундах
	}

	// Создаем сам консьюмер
	c, err := kafka.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Подписываемся на топик
	err = c.Subscribe(topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	// Создаем логгер, который будет сохранять все сообщения в лог файле
	logger := logrus.New()
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to create request log file: %v", err)
	}
	logger.Out = file
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Level = logrus.InfoLevel

	return &Consumer{consumer: c, logger: logger}, nil
}

// Метод для получения сообщений консьюмером
func (c *Consumer) Srart() {
	for {
		if c.stop { // Проверяем, что консьюмер не остановлен
			break
		}
		// Читаем сообщение
		kafkaMsg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Error("Failed to read message: ", err)
		}
		// Проверяем, что сообщение не пустое
		if kafkaMsg == nil {
			continue
		}
		var event map[string]interface{}
		err = json.Unmarshal(kafkaMsg.Value, &event)
		if err != nil {
			logrus.Error("Failed to unmarshal event: ", err)
		}
		c.logger.WithFields(logrus.Fields{"message": event}).Info("Kafka event received")

		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil { // Сохраняем оффсет сообщения
			logrus.Error("Failed to store message: ", err)
			continue
		}
	}
}

// Функуия для остановки консьюмера
func (c *Consumer) Stop() error {
	c.stop = true // Переводим флаг в true

	if _, err := c.consumer.Commit(); err != nil { // Когда останавливаем консьюмер, необходимо сохранить оффсет обработанных сообщений в кафке
		return fmt.Errorf("failed to commit offsets: %w", err)
	}
	logrus.Info("Commited offset")
	return c.consumer.Close() // Закрываем консьюмер
}
