package models

import "time"

type Task struct {
	ID          uint    `gorm:"primaryKey"`
	Judul       string  `gorm:"type:varchar(255);not null"`
	Status      string  `gorm:"type:varchar(50);not null;default:'todo'"`
	Tipe        string  `gorm:"type:varchar(50);not null"`
	PathProject *string `gorm:"type:text"`
	LinkWebsite *string `gorm:"type:text"`
	Tags        string  `gorm:"type:text"`
	Catatan     string  `gorm:"type:text"`
	Cover       string  `gorm:"type:varchar(255)"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
