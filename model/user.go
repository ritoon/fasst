package model

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	// UUID unic identifier
	UUID string `json:"uuid"`
	// FristName of user
	FristName string `json:"first_name"`
	// LastName of user
	LastName string `json:"last_name"`
	// Email
	Email string `json:"email"`
	// Password
	Password string `json:"password"`
}

type UserForUpdate struct {
	// FristName of user
	FristName string `json:"first_name"`
	// LastName of user
	LastName string `json:"last_name"`
	// Email
	Email string `json:"email"`
	// Password
	Password string `json:"password"`
}

func NewUser(fn, ln string) *User {
	return &User{
		UUID:      uuid.NewString(),
		FristName: fn,
		LastName:  ln,
	}
}

func (u User) ValidateForCreate() error {
	if len(u.Email) < 0 {
		return errors.New("no email")
	}
	return nil
}

func (u User) ValidateForUpdate() error {
	return nil
}

func (u *User) UpdateName(fn, ln string) {
	u.FristName = fn
	u.LastName = ln
}

func (u *User) UpdateFromMap(data map[string]interface{}) error {
	if u == nil || data == nil {
		return nil
	}

	if v, ok := data["first_name"]; ok {
		fn, ok := v.(string)
		if ok {
			u.FristName = fn
		}
	}

	if v, ok := data["last_name"]; ok {
		fn, ok := v.(string)
		if ok {
			u.LastName = fn
		}
	}

	if v, ok := data["email"]; ok {
		fn, ok := v.(string)
		if ok {
			u.Email = fn
		}
	}

	if v, ok := data["password"]; ok {
		fn, ok := v.(string)
		if ok {
			u.Password = fn
		}
	}

	return nil
}

type UserForPatch struct {
	// FristName of user
	FristName *string `json:"first_name"`
	// LastName of user
	LastName *string `json:"last_name"`
	// Email
	Email *string `json:"email"`
	// Password
	Password *string `json:"password"`
}

func (u UserForPatch) Validate() error {
	return nil
}
