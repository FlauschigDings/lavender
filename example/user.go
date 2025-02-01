package example

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Email    string
	Password string
}

func NewUser(email, password string) *User {
	return &User{
		Id:       uuid.New(),
		Email:    email,
		Password: password,
	}
}
