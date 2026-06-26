package models

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    Name      string         `gorm:"not null" json:"name" validate:"required"`
    Email     string         `gorm:"unique;not null" json:"email" validate:"required,email"`
    Password  string         `gorm:"not null" json:"-"` // "-" hides from JSON
    Role      string         `gorm:"default:driver" json:"role" validate:"oneof=driver admin"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}