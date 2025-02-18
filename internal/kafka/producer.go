package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Создадим структуру для продюссера
type Producer struct {
	producer *kafka.Producer
}

// Создаем новый продюсер, сюда передаем адреса наших брокеров, то есть в какие будем отправлять сообщения
func NewProducer(adress string) (*Producer, error) {
	// Создаем конфиг для продюссера
	conf := &kafka.ConfigMap{
		"bootstrap.servers": adress, // Список брокеров
	}
	// Создаем продюссера
	p, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	return &Producer{producer: p}, nil
}

// Метод для отправки сообщений продюссером
func (p *Producer) Produce(action, topic string) error {
	// Формируем сообщение
	event := map[string]interface{}{
		"timestamp": time.Now(),
		"action":    action,
	}

	msg, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	kafkaMsg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny}, // Указываем топик и партицию
		Value:          msg,                                                                // Наше сообщение
		Key:            nil,
	}

	// Создаем канал для получения ответа от продюссера, в этот канал приходит сообщение об успешности отправки
	kafkaChan := make(chan kafka.Event)

	// Отправляем сообщение
	err = p.producer.Produce(kafkaMsg, kafkaChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	// Теперь необходимо обработать ответ от продюссера
	e := <-kafkaChan

	// Проверяем тип у переданного события
	switch ev := e.(type) {
	case *kafka.Message: // Если это kafka.Message, значит отправка прошла успешно, и возвращаем nil
		return nil
	case kafka.Error:
		return ev // Если это ошибка(kafka.Error), то возвращаем эту ошибку
	default:
		return fmt.Errorf("unexpected event type: %T", ev) // В противном случае возвращаем ошибку
	}

}

// Метод для закрытия продюссера
// Если мы закрываем продюссера то мы хотим дождаться, что отправленные сообщения будут обработаны, исопльзуя метод Flush
func (p *Producer) Close() {
	p.producer.Flush(5000)
	p.producer.Close()
}
