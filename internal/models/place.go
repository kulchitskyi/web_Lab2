package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,latitude"`
	Longitude float64 `json:"longitude" validate:"required,longitude"`
}

// Value implements driver.Valuer interface for JSONB storage
func (l Location) Value() (driver.Value, error) {
	return json.Marshal(l)
}

// Scan implements sql.Scanner interface for JSONB retrieval
func (l *Location) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	return json.Unmarshal(bytes, l)
}

type PlaceCreateRequest struct {
	Name        string      `json:"name" validate:"required,min=5,max=100"`
	Description string      `json:"description" validate:"required,min=10,max=1000"`
	Location    Location    `json:"location" validate:"required"`
	Address     string      `json:"address" validate:"required,max=255"`
	Rating      int         `json:"averageRating" validate:"required,min=1,max=5"`
}

type Place struct {
	gorm.Model
	Id 			string 		`gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name 		string 		`gorm:"type:varchar(255);not null" json:"name"`
	Description string 		`gorm:"type:text" json:"description"`
	Location	Location	`gorm:"type:jsonb" json:"location"`
	Address		string 		`gorm:"type:varchar(255)" json:"address"`
	Rating		int 		`gorm:"type:numeric" json:"rating"`
	CreatedAt 	time.Time	`json:"createdAt"`
}

type PlaceUpdateRequest struct {
	Name        *string     `json:"name,omitempty" validate:"omitempty,min=5,max=100"`
	Description *string     `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Location    *Location   `json:"location,omitempty" validate:"omitempty"`
	Address     *string     `json:"address,omitempty" validate:"omitempty,max=255"`
	Rating      *int        `json:"averageRating,omitempty" validate:"omitempty,min=1,max=5"`
}