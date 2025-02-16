package handler

import (
	"mainPet/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	services    *service.Service
	requestLog  *logrus.Logger
	responseLog *logrus.Logger
	errorLog    *logrus.Logger
}

func NewHandler(services *service.Service, requestLog, responseLog, errorLog *logrus.Logger) *Handler {
	return &Handler{
		services:    services,
		requestLog:  requestLog,
		responseLog: responseLog,
		errorLog:    errorLog,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Group("/api")
	{
		tasks := router.Group("/tasks")
		{
			tasks.POST("/", h.CreateTask)
			tasks.GET("/", h.GetAllTasks)
			tasks.GET("/:id", h.GetTaskById)
			tasks.DELETE("/:id", h.DeleteTask)
			tasks.PUT("/:id/done", h.DoneTask)
		}
	}

	return router
}
