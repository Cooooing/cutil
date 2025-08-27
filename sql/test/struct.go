package test

import "time"

type User struct {
	Id        *int       `json:"id"`
	Name      *string    `json:"name"`
	Age       *int       `json:"age"`
	Email     *string    `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
}
