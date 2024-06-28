package repository

import (
	"congratulations_service/pkg/model"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

type Repository struct {
	mu            sync.RWMutex
	employeeStore map[int]model.Employee
	subscriptions map[string][]int
}

func NewRepository() *Repository {
	return &Repository{
		employeeStore: make(map[int]model.Employee),
		subscriptions: make(map[string][]int),
	}
}

func (r *Repository) AddEmployeeToStore(name, birthday, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Генерируем ID для нового сотрудника
	id := len(r.employeeStore) + 1

	// Хешируем пароль
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	employee := model.Employee{
		ID:           id,
		Name:         name,
		Birthday:     birthday,
		PasswordHash: string(passwordHash),
	}

	if _, exists := r.employeeStore[employee.ID]; exists {
		return fmt.Errorf("employee with id %d already exists", employee.ID)
	}

	r.employeeStore[employee.ID] = employee
	fmt.Println("Добавили сотрудника в базу")
	return nil
}

func (r *Repository) DeleteEmployeeFromStore(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exist := r.employeeStore[id]; !exist {
		return fmt.Errorf("employee with id %d does not exist", id)
	}
	delete(r.employeeStore, id)
	fmt.Println("Удалили сотрудника из базы")
	return nil
}

func (r *Repository) GetAllEmployees() ([]model.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := make([]model.Employee, 0, len(r.employeeStore))
	for _, val := range r.employeeStore {
		res = append(res, val)
	}
	return res, nil
}

func (r *Repository) GetEmployeeInfo(id int) (model.EmployeeInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var res model.EmployeeInfo

	for _, val := range r.employeeStore {
		if val.ID == id {
			res.ID = val.ID
			res.Name = val.Name
			res.Birthday = val.Birthday
			return res, nil
		}
	}
	return res, fmt.Errorf("employee with id %d not found", id)
}

func (r *Repository) AuthenticateEmployee(name, password string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, employee := range r.employeeStore {
		if employee.Name == name {
			err := bcrypt.CompareHashAndPassword([]byte(employee.PasswordHash), []byte(password))
			if err != nil {
				return false, fmt.Errorf("authentication failed: %v", err)
			}
			return true, nil
		}
	}
	return false, fmt.Errorf("employee not found")
}

func (r *Repository) AddSubscription(username string, employeeID int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subscriptions[username] = append(r.subscriptions[username], employeeID)
	fmt.Printf("%s Подписался на %s\n", username, r.employeeStore[employeeID].Name)
}

func (r *Repository) RemoveSubscription(username string, employeeID int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	subs := r.subscriptions[username]
	for i, sub := range subs {
		if sub == employeeID {
			r.subscriptions[username] = append(subs[:i], subs[i+1:]...)
			fmt.Printf("%s Отписался от %s\n", username, r.employeeStore[employeeID].Name)
			break
		}
	}
}

func (r *Repository) CheckBirthdays() error {
	checkBirthdays := func() {
		now := time.Now()
		r.mu.RLock()

		defer r.mu.RUnlock()

		for _, employee := range r.employeeStore {
			birthDate, err := time.Parse("2006-01-02", employee.Birthday)
			if err != nil {
				fmt.Printf("Could not parse birthday for employee ID %d: %v\n", employee.ID, err)
				continue
			}

			// Сравниваем только месяц и день
			if birthDate.Month() == now.Month() && birthDate.Day() == now.Day() {
				r.sendNotifications(employee)
			}
		}
	}

	// Проверяем дни рождения три раза каждую минуту чтобы наглядно увидеть оповещение при регистрации
	for i := 0; i < 3; i++ {
		checkBirthdays()
		time.Sleep(1 * time.Minute)
	}

	// Затем переходим к ежедневной проверке
	for {
		checkBirthdays()

		now := time.Now()
		// Ждем до следующего дня
		nextDay := now.AddDate(0, 0, 1)
		nextMidnight := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		time.Sleep(nextMidnight.Sub(now))
	}
}

func (r *Repository) sendNotifications(birthdayEmployee model.Employee) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for username, subs := range r.subscriptions {
		for _, subID := range subs {
			if subID == birthdayEmployee.ID {
				// Логика отправки уведомления текущему авторизованному пользователю
				fmt.Printf("%s, у %s сегодня день рождения! Поздравьте =)\n", username, birthdayEmployee.Name)
				// Дальше можно маштабировать путем отправки на почту или мессенджер и так далее
			}
		}
	}
}
