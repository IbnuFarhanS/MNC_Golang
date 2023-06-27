package models

type Customer struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    int    `json:"phone"`
}
