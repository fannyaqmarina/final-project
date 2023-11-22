package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Variant struct {
	gorm.Model
	Uuid        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	VariantName string
	Quantity    int
	ProductID   uint
	Products    Product `gorm:"foreignKey:ProductID"`
}
