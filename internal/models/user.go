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
