package repository

import (
	"fmt"
	"mainPet/internal/model"

	"github.com/jmoiron/sqlx"
)

// Определяем структуру для работы с postgres
type TodoListPostgres struct {
	db *sqlx.DB
}

// конструктор, который создаёт новый экземпляр структуры TodoListPostgres
func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

// Теперь реализуем все методы

// Создание одной заметки
func (r *TodoListPostgres) Create(task model.Task) (int, error) {
	var id int
	// Создаем запрос к бд
	query := `INSERT INTO todo_list (title, description) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, task.Title, task.Description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create todo list: %w", err)
	}

	return id, nil
}

// Получение всех заметок
func (r *TodoListPostgres) GetAll() ([]model.Task, error) {
	var todos []model.Task

	query := `SELECT * FROM todo_list`

	err := r.db.Select(&todos, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all todos: %w", err)
	}

	return todos, nil
}

// Удаление заметки по id
func (r *TodoListPostgres) Delete(id int) error {
	query := `DELETE FROM todo_list WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task with id %d: %w", id, err)
	}

	return nil
}

// Отмечаем заметку как сделанную
func (r *TodoListPostgres) Done(id int) error {
	query := `UPDATE todo_list SET done = true WHERE id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to mark task with id %d as done: %w", id, err)
	}

	return nil
}
