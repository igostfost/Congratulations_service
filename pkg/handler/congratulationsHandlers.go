package handler

import (
	"congratulations_service/pkg/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var jwtKey = []byte("secret_key")
var subscriptions = make(map[string][]int)
var mu sync.RWMutex

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func (h *Handler) SignIn(c *gin.Context) {
	var creds model.Credentials
	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Аутентификация сотрудника
	authenticated, err := h.service.SignInEmployee(creds.Username, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !authenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged in"})
}

func (h *Handler) SignUp(c *gin.Context) {
	var reg model.EmployeeRegistration
	if err := c.BindJSON(&reg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SignUpEmployee(reg.Name, reg.Birthday, reg.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee registered"})
}

func (h *Handler) GetEmployees(c *gin.Context) {
	employees, err := h.service.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Преобразуем значение в строку
	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found or invalid type"})
		return
	}
	fmt.Println("Текущий пользователь - ", username)
}

func (h *Handler) GetEmployeeInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	employeeInfo, err := h.service.GetEmployeeInfo(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employeeInfo)
}

func (h *Handler) DeleteEmployee(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.service.DeleteEmployee(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee deleted"})
}

func (h *Handler) Subscribe(c *gin.Context) {
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Преобразуем значение в строку
	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found or invalid type"})
		return
	}

	// Структура для привязки JSON
	var requestBody struct {
		EmployeeID int `json:"follow_id"`
	}

	// Привязываем JSON тело запроса к структуре requestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body"})
		return
	}

	// Получаем id сотрудника из тела запроса
	id := requestBody.EmployeeID

	h.service.Subscribe(username, id)
	c.JSON(http.StatusOK, gin.H{"message": "subscribed"})
}

func (h *Handler) Unsubscribe(c *gin.Context) {
	usernameInterface, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Преобразуем значение в строку
	username, ok := usernameInterface.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "username not found or invalid type"})
		return
	}

	idStr := c.Query("employee_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
		return
	}

	h.service.Unsubscribe(username, id)
	c.JSON(http.StatusOK, gin.H{"message": "unsubscribed"})
}
