package models

import "time"

type Comment struct {
	Id        	uint   		`gorm:"primaryKey;" json:"id"`
	Message     string 		`gorm:"not null;" validate:"required" json:"message"`
	User_id	  	int 		`json:"user_id"`
	User	  	User 		`gorm:"foreignKey:User_id;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Photo_id	int 		`json:"photo_id"`
	Photo		Photo 		`gorm:"foreignKey:Photo_id;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
}