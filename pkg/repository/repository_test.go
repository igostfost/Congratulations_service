package repository

import (
	"congratulations_service/pkg/model"
	"testing"
	"time"
)

func TestAddEmployeeToStore(t *testing.T) {
	repo := NewRepository()
	err := repo.AddEmployeeToStore("Alice", "1999-06-28", "password123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	employees, err := repo.GetAllEmployees()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(employees) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(employees))
	}
}

func TestDeleteEmployeeFromStore(t *testing.T) {
	repo := NewRepository()
	repo.AddEmployeeToStore("Alice", "1999-06-28", "password123")
	err := repo.DeleteEmployeeFromStore(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	employees, err := repo.GetAllEmployees()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(employees) != 0 {
		t.Fatalf("expected 0 employees, got %d", len(employees))
	}
}

func TestAuthenticateEmployee(t *testing.T) {
	repo := NewRepository()
	repo.AddEmployeeToStore("Alice", "1999-06-28", "password123")

	auth, err := repo.AuthenticateEmployee("Alice", "password123")
	if err != nil || !auth {
		t.Fatalf("expected authentication to succeed, got %v, %v", auth, err)
	}

	auth, err = repo.AuthenticateEmployee("Alice", "wrongpassword")
	if err == nil || auth {
		t.Fatalf("expected authentication to fail, got %v, %v", auth, err)
	}
}

func TestAddSubscription(t *testing.T) {
	repo := NewRepository()
	repo.AddEmployeeToStore("Alice", "1999-06-28", "password123")

	repo.AddSubscription("Bob", 1)
	if len(repo.subscriptions["Bob"]) != 1 {
		t.Fatalf("expected 1 subscription for Bob, got %d", len(repo.subscriptions["Bob"]))
	}
}

func TestRemoveSubscription(t *testing.T) {
	repo := NewRepository()
	repo.AddEmployeeToStore("Alice", "1999-06-28", "password123")

	repo.AddSubscription("Bob", 1)
	repo.RemoveSubscription("Bob", 1)
	if len(repo.subscriptions["Bob"]) != 0 {
		t.Fatalf("expected 0 subscriptions for Bob, got %d", len(repo.subscriptions["Bob"]))
	}
}

func TestCheckBirthdays(t *testing.T) {
	repo := NewRepository()
	now := time.Now().Format("2006-01-02")
	repo.AddEmployeeToStore("Alice", now, "password123")

	go func() {
		if err := repo.CheckBirthdays(); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	}()

	// Добавляем небольшую задержку, чтобы дать возможность CheckBirthdays сработать
	time.Sleep(2 * time.Second)
	// Проверка на наличие уведомлений
}

func TestSendNotifications(t *testing.T) {
	repo := NewRepository()
	now := time.Now().Format("2006-01-02")
	if err := repo.AddEmployeeToStore("Alice", now, "password123"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	repo.AddSubscription("Bob", 1)

	repo.sendNotifications(model.Employee{ID: 1, Name: "Alice", Birthday: now})

	// Здесь нужно проверить вывод, например, используя логирование или другой способ,
	// чтобы убедиться, что уведомления отправлены
}
