package main

import (
	"log"
	"mainPet/internal/handler"
	"mainPet/internal/repository"
	"mainPet/internal/service"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	// Initialize logger
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// Initialize config
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing config: %s", err.Error())
	}
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	// Инициализируем postgres
	db, err := repository.NewPostgres(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	// Инициализируем репозитории
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handlers := handler.NewHandler(service)

	server := &http.Server{
		Addr:         viper.GetString("port"),
		Handler:      handlers.InitRoutes(),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
