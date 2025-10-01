package model_test

import (
	"testing"

	"formation/model"
)

func TestUserSayHello(t *testing.T) {
	u := model.NewUser("toto", "titi")
	u.ValidateForCreate()

}
