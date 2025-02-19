package main

import (
	"mainPet/internal/kafka"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Consumer кафки, который будет получать сообщения и сохранять их в лог
func main() {
	c, err := kafka.NewConcumer("kafka:29091", "todo_list", "mainpet", "logs/kafka_events.log")
	if err != nil {
		logrus.Fatalf("failed to create consumer: %s", err.Error())
	}

	// Запускаем консюмер в горутине
	go func() {
		c.Srart()
	}()

	// реализуем graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Ожидаем сигнал
	<-sigchan
	// И останавливаем консюмер
	logrus.Fatal(c.Stop())
}
