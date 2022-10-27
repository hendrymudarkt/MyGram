package models

import "time"

type SocialMedia struct {
	Id        		uint   		`gorm:"primaryKey;" json:"id"`
	Name     		string 		`gorm:"not null;" validate:"required" json:"name"`
	SocialMediaUrl  string 		`gorm:"not null;" validate:"required" json:"social_media_url"`
	User_id			int			`json:"user_id"`
	User	  		User 		`gorm:"foreignKey:User_id;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt		time.Time 	`json:"created_at"`
	UpdatedAt		time.Time 	`json:"updated_at"`
}