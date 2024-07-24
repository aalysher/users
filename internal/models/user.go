package models

import "time"

type User struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Age       uint      `json:"age"`
	Email     string    `json:"email"`
	Created   time.Time `json:"created"`
}

type UserUpdate struct {
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	Age       *uint   `json:"age,omitempty"`
	Email     *string `json:"email,omitempty"`
}
