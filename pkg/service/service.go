package service

import (
	"congratulations_service/pkg/model"
	"congratulations_service/pkg/repository"
)

type Service struct {
	Repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Repo: repo,
	}
}

func (s *Service) SignUpEmployee(name, birthday, password string) error {
	return s.Repo.AddEmployeeToStore(name, birthday, password)
}

func (s *Service) SignInEmployee(name, password string) (bool, error) {
	return s.Repo.AuthenticateEmployee(name, password)
}

func (s *Service) GetAllEmployees() ([]model.Employee, error) {
	return s.Repo.GetAllEmployees()
}

func (s *Service) GetEmployeeInfo(id int) (model.EmployeeInfo, error) {
	return s.Repo.GetEmployeeInfo(id)
}

func (s *Service) DeleteEmployee(id int) error {
	return s.Repo.DeleteEmployeeFromStore(id)
}

func (s *Service) Subscribe(username string, id int) {
	s.Repo.AddSubscription(username, id)
}

func (s *Service) Unsubscribe(username string, id int) {
	s.Repo.RemoveSubscription(username, id)
}

func (s *Service) CheckBirthdays() error {
	return s.Repo.CheckBirthdays()
}
