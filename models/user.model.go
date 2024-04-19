package models

import "gorm.io/gorm"

//json untuk membaca di postman
//gorm digunakan untuk membuat constraintnya
type User struct {
	gorm.Model
	Email string `gorm:"unique" json:"email"`
	Password string
}