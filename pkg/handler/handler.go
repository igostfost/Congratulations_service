package handler

import (
	"congratulations_service/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/sign-in", h.SignIn)
	router.POST("/sign-up", h.SignUp)

	protected := router.Group("/")
	protected.Use(h.AuthMiddleware())
	{
		protected.POST("/subscribe", h.Subscribe)
		protected.POST("/unsubscribe", h.Unsubscribe)
		protected.GET("/employee/:id", h.GetEmployeeInfo)
		protected.GET("/employees", h.GetEmployees)
		protected.DELETE("/del_employees/:id", h.DeleteEmployee)
	}

	return router
}
