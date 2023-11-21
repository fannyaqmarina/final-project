package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Uuid     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name     string
	ImageUrl string
	AdminID  uint
	Variants []Variant `gorm:"foreignKey:ProductID"`
}
