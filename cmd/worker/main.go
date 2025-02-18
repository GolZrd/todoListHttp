package main

import (
	"mainPet/internal/kafka"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Consumer кафки, который будет получать сообщения и сохранять их в лог
func main() {
	c, err := kafka.NewConcumer(viper.GetStringSlice("kafka.broker"), os.Getenv("KAFKA_TOPIC"), viper.GetString("kafka.consumerGroup"), "logs/kafka_events.log")
	if err != nil {
		logrus.Fatalf("failed to create consumer: %s", err.Error())
	}

	// Запускаем консюмер в горутине
	go c.Srart()

	// реализуем graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнал
	<-sigchan
	// И останавливаем консюмер
	logrus.Fatal(c.Stop())
}
