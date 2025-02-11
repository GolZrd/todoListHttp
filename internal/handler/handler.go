package handler

import (
	"mainPet/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Group("/api")
	{
		tasks := router.Group("/tasks")
		{
			tasks.POST("/", h.CreateTask)
			tasks.GET("/", h.GetAllTasks)
			tasks.DELETE("/:id", h.DeleteTask)
			tasks.PUT("/:id/done", h.DoneTask)
		}
	}

	return router
}
