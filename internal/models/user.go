package models

import validation "github.com/go-ozzo/ozzo-validation"

type Users struct {
	ID       int    `json:"id" gorm:"autoIncrement;column:id"`
	Email    string `json:"email" gorm:"column:email"`
	Password string `json:"password" gorm:"column:password"`
	Isim     string `json:"isim" gorm:"column:isim"`
	Soyisim  string `json:"soyisim" gorm:"column:soyisim"`
	Resim    string `json:"resim" gorm:"column:resim"`
}

type CreateUsers struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Isim     string `json:"isim"`
	Soyisim  string `json:"soyisim"`
}

type UpdateUsers struct {
	Password string `json:"password"`
	Isim     string `json:"isim"`
	Soyisim  string `json:"soyisim"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m CreateUsers) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Email, validation.Required, validation.Length(5, 50)),
		validation.Field(&m.Password, validation.Required, validation.Length(5, 50)),
		validation.Field(&m.Isim, validation.Required, validation.Length(5, 50)),
		validation.Field(&m.Soyisim, validation.Required, validation.Length(5, 50)))
}

func (m UpdateUsers) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Password, validation.Required, validation.Length(5, 50)),
		validation.Field(&m.Isim, validation.Required, validation.Length(5, 50)),
		validation.Field(&m.Soyisim, validation.Required, validation.Length(5, 50)))
}
