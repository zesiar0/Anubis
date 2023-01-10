package model

import "fmt"

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	IsValid  bool   `json:"is_valid"`
	IsActive bool   `json:"is_active"`
}

func (u *User) String() string {
	return fmt.Sprintf("%s(%s)", u.Name, u.Username)
}
