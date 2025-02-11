package service

import (
	"mainPet/internal/model"
	"mainPet/internal/repository"
)

type TodoListService struct {
	repo repository.TodoList
}

func NewTodoListService(repo repository.TodoList) *TodoListService {
	return &TodoListService{repo: repo}
}

func (s *TodoListService) Create(task model.Task) (int, error) {
	return s.repo.Create(task)
}

func (s *TodoListService) GetAll() ([]model.Task, error) {
	return s.repo.GetAll()
}

func (s *TodoListService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TodoListService) Done(id int) error {
	return s.repo.Done(id)
}
