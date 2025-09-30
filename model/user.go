package model

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	// UUID unic identifier
	UUID string
	// FristName of user
	FristName string
	// LastName of user
	LastName string
}

func NewUser(fn, ln string) *User {
	return &User{
		UUID:      uuid.NewString(),
		FristName: fn,
		LastName:  ln,
	}
}

func (u User) SayHello() string {
	return fmt.Sprintln(u.FristName, "Hello")
}

func (u *User) UpdateName(fn, ln string) {
	u.FristName = fn
	u.LastName = ln
}
