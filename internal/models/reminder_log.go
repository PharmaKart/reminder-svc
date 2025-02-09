package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReminderLog struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	ReminderID uuid.UUID `gorm:"not null"`
	OrderID    uuid.UUID `gorm:"not null"`
	Status     string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"default:now()"`
}

func (rl *ReminderLog) BeforeCreate(tx *gorm.DB) (err error) {
	rl.ID = uuid.New()
	return
}
