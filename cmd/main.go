package main

import (
	"congratulations_service"
	"congratulations_service/pkg/handler"
	"congratulations_service/pkg/model"
	"congratulations_service/pkg/repository"
	"congratulations_service/pkg/service"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	repo := repository.NewRepository()
	// Инициализируем репозиторий с начальными фейковыми данными (для наглядности)
	initFakeEmployees(repo)

	services := service.NewService(repo)
	handlers := handler.NewHandler(services)
	srv := new(congratulations_service.Server)

	go func() {
		if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
			log.Fatalf("error running http server: %s", err)
		}
		fmt.Println("Congratulations service is running on port 8080")
	}()

	go func() {
		if err := services.CheckBirthdays(); err != nil {
			log.Println(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	fmt.Println("APP shutting down")
	if err := srv.ShutDown(context.Background()); err != nil {
		fmt.Printf("error occurred on server shutting down: %s", err)
	}
}

func initFakeEmployees(repo *repository.Repository) {
	now := time.Now()
	today := now.Format("2006-01-02")

	// Примеры начальных значений сотрудников
	initialEmployees := []model.EmployeeRegistration{
		{Name: "Alice", Birthday: "1999-06-28"},
		{Name: "Bob", Birthday: "1992-03-13"},
		{Name: "Charlie", Birthday: today},
	}

	for _, emp := range initialEmployees {
		if err := repo.AddEmployeeToStore(emp.Name, emp.Birthday, emp.Password); err != nil {
			log.Printf("Error adding employee %v: %v", emp, err)
		}
	}
}
