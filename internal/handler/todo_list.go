package handler

import (
	"mainPet/internal/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Создание задачи
func (h *Handler) CreateTask(c *gin.Context) {
	// Логирование входящего запроса
	h.requestLog.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")

	// Определяем структуру нашей задачи и заполняем ее данными из запроса
	var input model.Task

	if err := c.BindJSON(&input); err != nil {
		h.errorLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to bind JSON")
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Создаем задачу
	id, err := h.services.TodoList.Create(input)
	if err != nil {
		h.errorLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create task")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Логирование исходящего ответа
	h.responseLog.WithFields(logrus.Fields{
		"status": http.StatusOK,
		"id":     id,
	}).Info("Outgoing response")

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// Получение всех задач
func (h *Handler) GetAllTasks(c *gin.Context) {
	h.requestLog.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")

	tasks, err := h.services.TodoList.GetAll()
	if err != nil {
		h.errorLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to get all tasks")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.responseLog.WithFields(logrus.Fields{
		"status": http.StatusOK,
	}).Info("Outgoing response")

	c.JSON(http.StatusOK, tasks)
}

// Отмечаем задачу как сделанную
func (h *Handler) DoneTask(c *gin.Context) {
	h.requestLog.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to convert id to int")
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoList.Done(id)
	if err != nil {
		h.errorLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to mark task as done")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.responseLog.WithFields(logrus.Fields{
		"status": http.StatusOK,
	}).Info("Outgoing response")

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// Удаление задачи
func (h *Handler) DeleteTask(c *gin.Context) {
	h.requestLog.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"path":   c.Request.URL.Path,
	}).Info("Incoming request")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to convert id to int")
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoList.Delete(id)
	if err != nil {
		h.errorLog.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to delete task")
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	h.responseLog.WithFields(logrus.Fields{
		"status": http.StatusOK,
	}).Info("Outgoing response")

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
