package service

import (
	"context"
	"encoding/json"
	"mainPet/internal/kafka"
	"mainPet/internal/model"
	"mainPet/internal/repository"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type TodoListService struct {
	repo  repository.TodoList
	cache *redis.Client
	kafka kafka.Producer
}

func NewTodoListService(repo repository.TodoList, cache *redis.Client, kafka kafka.Producer) *TodoListService {
	return &TodoListService{
		repo:  repo,
		cache: cache,
		kafka: kafka,
	}
}

func (s *TodoListService) Create(task model.Task) (int, error) {
	if err := s.kafka.Produce("Create", os.Getenv("KAFKA_TOPIC")); err != nil {
		logrus.Warn("Failed to produce message: ", err)
	}

	return s.repo.Create(task)
}

func (s *TodoListService) GetAll() ([]model.Task, error) {
	if err := s.kafka.Produce("GetAll", os.Getenv("KAFKA_TOPIC")); err != nil {
		logrus.Warn("Failed to produce message: ", err)
	}
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

	if err := s.kafka.Produce("GetById", os.Getenv("KAFKA_TOPIC")); err != nil {
		logrus.Warn("Failed to produce message: ", err)
	}

	return task, err
}

func (s *TodoListService) Delete(id int) error {
	s.cache.Del(context.Background(), strconv.Itoa(id))
	if err := s.kafka.Produce("Delete", os.Getenv("KAFKA_TOPIC")); err != nil {
		logrus.Warn("Failed to produce message: ", err)
	}

	return s.repo.Delete(id)
}

func (s *TodoListService) Done(id int) error {
	s.cache.Del(context.Background(), strconv.Itoa(id))
	if err := s.kafka.Produce("Done", os.Getenv("KAFKA_TOPIC")); err != nil {
		logrus.Warn("Failed to produce message: ", err)
	}
	return s.repo.Done(id)
}
