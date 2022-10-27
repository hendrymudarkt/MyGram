package models

import "time"

type Photo struct {
	Id        	uint   		`gorm:"primaryKey;" json:"id"`
	Title  		string 		`gorm:"not null;" validate:"required" json:"title"`
	Caption     string 		`json:"caption"`
	Photo_url  	string 		`gorm:"not null;" validate:"required" json:"photo_url"`
	User_id		int			`json:"user_id"`
	User	  	User 		`gorm:"foreignKey:User_id;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	CreatedAt	time.Time 	`json:"created_at"`
	UpdatedAt	time.Time 	`json:"updated_at"`
}
