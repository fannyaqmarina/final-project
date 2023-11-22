package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Uuid     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name     string    `json:"name"`
	ImageUrl string    `json:"image_url"`
	AdminID  uint      `json:"admin_id"`
	Admin    Admin     `gorm:"foreignKey:AdminID"`
	Variants []Variant `gorm:"foreignKey:ProductID"`
}
