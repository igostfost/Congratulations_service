package model

type Employee struct {
	ID           int    `json:"id"`
	Name         string `json:"username"`
	Birthday     string `json:"birthday"` // "YYYY-MM-DD"
	PasswordHash string `json:"-"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type EmployeeRegistration struct {
	Name     string `json:"username"`
	Birthday string `json:"birthday"` // "YYYY-MM-DD"
	Password string `json:"password"`
}

type EmployeeInfo struct {
	ID       int    `json:"id"`
	Name     string `json:"username"`
	Birthday string `json:"birthday"` // "YYYY-MM-DD"
}
