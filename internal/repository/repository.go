package repository

import (
	"mainPet/internal/model"

	"github.com/jmoiron/sqlx"
)

// В этом файле мы будем создавать наш репозиторий, и использовать его в нашем сервисе
// Для репозитория определим интерфейс, и реализуем его в файле todo_list_postgres.go

// Определяем интерфейс для нашего репозитория, который будет содержать все методы для работы с нашей базой данных
type TodoList interface {
	Create(task model.Task) (int, error) //Создание одной заметки
	GetAll() ([]model.Task, error)       //Получение всех заметок
	Delete(id int) error                 //Удаление одной заметки
	Done(id int) error                   //Отметка заметки как выполненной
}

// Определяем структуру для нашего репозитория
// Внутри структуры будем наш интерфейс для работы с нашей базой данных
type Repository struct {
	TodoList
}

// Определяем функцию для создания нового репозитория
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TodoList: NewTodoListPostgres(db),
	}
}
