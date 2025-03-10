package service

import (
	"mainPet/internal/kafka"
	"mainPet/internal/model"
	"mainPet/internal/repository"

	"github.com/redis/go-redis/v9"
)

// Определяем интерфейс для работы с задачами
type TodoList interface {
	Create(task model.Task) (int, error)
	GetAll() ([]model.Task, error)
	GetById(id int) (model.Task, error)
	Delete(id int) error
	Done(id int) error
}

type Service struct {
	TodoList
	Cache *redis.Client
	Kafka kafka.Producer
}

// Конструктор для создания нового сервиса
func NewService(repo *repository.Repository, cache *redis.Client, kafka kafka.Producer) *Service {
	return &Service{
		TodoList: NewTodoListService(repo.TodoList, cache, kafka),
	}
}
