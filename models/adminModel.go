package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Uuid     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email    string    `gorm:"unique"`
	Password string
}
