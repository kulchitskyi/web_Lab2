package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id 			string 		`gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name 		string 		`gorm:"type:varchar(255);not null" json:"username"`
	Email 		string 		`gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	CreatedAt 	time.Time 	`json:"createdAt"`
}

type UserCreateRequest struct {
	Name        string  `json:"username" validate:"required,min=3,max=50"`
	Email       string  `json:"email" validate:"required,email"`
}

type UserUpdateRequest struct {
	Name        *string     `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email       *string     `json:"email,omitempty" validate:"omitempty,email"`
}

type UserPlace struct {
	UserID    string `gorm:"primaryKey;type:uuid" json:"user_id" validate:"required,uuid"`
	PlaceID   string `gorm:"primaryKey;type:uuid" json:"place_id" validate:"required,uuid"`
}