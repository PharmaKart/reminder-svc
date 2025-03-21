package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reminder struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CustomerID   uuid.UUID `gorm:"not null"`
	OrderID      uuid.UUID `gorm:"not null"`
	ProductID    uuid.UUID `gorm:"not null"`
	ReminderDate time.Time `gorm:"type:timestamptz;not null"`
	LastSentAt   time.Time `gorm:"type:timestamptz;default:null"`
	Enabled      bool      `gorm:"default:true"`
	CreatedAt    time.Time `gorm:"type:timestamptz;default:now()"`
}

func (r *Reminder) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
