package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Name         string   `json:"name"`
	Email        string   `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string   `gorm:"not null" json:"-"`
	Role         UserRole `gorm:"type:varchar(20);not null;default:'student'" json:"role"`
}
