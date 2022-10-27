package models

import "time"

type User struct {
	Id        uint   	`gorm:"primaryKey;" json:"id"`
	Email     string 	`validate:"required,email" gorm:"type:varchar(255);uniqueIndex;not null;" json:"email"`
	Username  string 	`validate:"required" gorm:"type:varchar(100);uniqueIndex;not null;" json:"username"`
	Password  string 	`validate:"required,min=6" gorm:"not null;" json:"-"`
	Age       int    	`validate:"required,gte=8" gorm:"not null;" json:"age"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}