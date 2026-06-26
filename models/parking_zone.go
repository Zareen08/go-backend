package models

import (
    "time"
    "gorm.io/gorm"
)

type ParkingZone struct {
    ID            uint           `gorm:"primaryKey" json:"id"`
    Name          string         `gorm:"not null" json:"name" validate:"required"`
    Type          string         `gorm:"not null" json:"type" validate:"oneof=general ev_charging covered"`
    TotalCapacity int            `gorm:"not null" json:"total_capacity" validate:"required,gt=0"`
    PricePerHour  float64        `gorm:"type:decimal(10,2);not null" json:"price_per_hour" validate:"required,gt=0"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}