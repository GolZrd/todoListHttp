package service

import (
	"context"
	"encoding/json"
	"mainPet/internal/model"
	"mainPet/internal/repository"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type TodoListService struct {
	repo  repository.TodoList
	cache *redis.Client
}

func NewTodoListService(repo repository.TodoList, cache *redis.Client) *TodoListService {
	return &TodoListService{
		repo:  repo,
		cache: cache,
	}
}

func (s *TodoListService) Create(task model.Task) (int, error) {
	return s.repo.Create(task)
}

func (s *TodoListService) GetAll() ([]model.Task, error) {
	return s.repo.GetAll()
}

func (s *TodoListService) GetById(id int) (model.Task, error) {
	// Пытаемся получить задачу из кэша
	cachedTask, err := s.cache.Get(context.Background(), strconv.Itoa(id)).Result()
	if err == nil {
		var task model.Task
		err := json.Unmarshal([]byte(cachedTask), &task)
		if err != nil {
			logrus.Warn("Failed to unmarshal task from cache: ", err)
		}
		return task, nil
	} else if err == redis.Nil {
		logrus.Info("Task not found in cache")
	} else {
		logrus.Warn("Failed to get task from cache: ", err)
	}

	task, err := s.repo.GetById(id)
	if err != nil {
		return model.Task{}, err
	}
	taskJson, err := json.Marshal(task)
	if err != nil {
		logrus.Warn("Failed to marshal task: ", err)
	} else {
		s.cache.Set(context.Background(), strconv.Itoa(id), taskJson, 20*time.Second)
	}

	return task, err
}

func (s *TodoListService) Delete(id int) error {
	s.cache.Del(context.Background(), strconv.Itoa(id))
	return s.repo.Delete(id)
}

func (s *TodoListService) Done(id int) error {
	s.cache.Del(context.Background(), strconv.Itoa(id))
	return s.repo.Done(id)
}
