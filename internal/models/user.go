package models

import "github.com/go-playground/validator/v10"

type User struct {
	Id           int    `json:"id"`
	Login        string `json:"login" validate:"required"`
	Pwd          string `json:"password,omitempty" validate:"required"`
	EncryptedPwd string `json:"-"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
