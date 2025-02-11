package service

import (
	"mainPet/internal/model"
	"mainPet/internal/repository"
)

// Определяем интерфейс для работы с задачами
type TodoList interface {
	Create(task model.Task) (int, error)
	GetAll() ([]model.Task, error)
	Delete(id int) error
	Done(id int) error
}

type Service struct {
	TodoList
}

// Конструктор для создания нового сервиса
func NewService(repo *repository.Repository) *Service {
	return &Service{
		TodoList: NewTodoListService(repo.TodoList),
	}
}
