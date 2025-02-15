package main

import (
	"mainPet/internal/handler"
	"mainPet/internal/logger"
	"mainPet/internal/repository"
	"mainPet/internal/service"
	"mainPet/pkg/migrator"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	// Инициализируем логгеры
	requestLogger := logger.NewRequestLogger("logs/request.log")
	responseLogger := logger.NewResponseLogger("logs/response.log")
	errorLogger := logger.NewErrorLogger("logs/error.log")

	// Initialize config
	if err := initConfig(); err != nil {
		errorLogger.Fatalf("error initializing config: %s", err.Error())
	}
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		errorLogger.Fatalf("error loading env variables: %s", err.Error())
	}

	// Инициализируем postgres
	db, err := repository.NewPostgres(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		errorLogger.Fatalf("failed to initialize db: %s", err.Error())
	}

	// Запускаем миграции
	if err := migrator.Migrate(db); err != nil {
		errorLogger.Fatalf("failed to migrate db: %s", err.Error())
	}

	// Инициализируем репозитории
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	handlers := handler.NewHandler(service, requestLogger, responseLogger, errorLogger)

	server := &http.Server{
		Addr:         ":" + viper.GetString("srv.port"),
		Handler:      handlers.InitRoutes(),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		errorLogger.Fatalf("error occured while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
