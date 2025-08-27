package test

import "time"

type User struct {
	Id        *int       `json:"id"`
	Name      *string    `json:"name"`
	Age       *int       `json:"age"`
	Email     *string    `json:"email"`
	CreatedAt *time.Time `json:"created_at"`
}

type PostTitle struct {
	UserId    *int    `json:"id"`
	UserName  *string `json:"name"`
	PostTitle *string `json:"title"`
}
