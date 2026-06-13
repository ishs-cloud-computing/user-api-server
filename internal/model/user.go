package model

import "time"

type User struct {
	ID        uint32    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Email     string    `db:"email" json:"email"`
	Age       uint8     `db:"age" json:"age"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// CreateUserRequest - POST /v1/users
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   uint8  `json:"age"`
}

// UpdateUserRequest - PATCH /v1/users/:id
type UpdateUserRequest struct {
	Name string `json:"name"`
	Age  uint8  `json:"age"`
}
