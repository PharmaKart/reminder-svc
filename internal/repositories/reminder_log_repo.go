package repositories

import (
	"github.com/PharmaKart/reminder-svc/internal/models"
	"gorm.io/gorm"
)

type ReminderLogRepository interface {
	CreateReminderLog(reminderLog *models.ReminderLog) error
	GetReminderLogsByReminderID(reminderID string) (*[]models.ReminderLog, error)
}

type reminderLogRepository struct {
	db *gorm.DB
}

func NewReminderLogRepository(db *gorm.DB) ReminderLogRepository {
	return &reminderLogRepository{db}
}

func (r *reminderLogRepository) CreateReminderLog(reminderLog *models.ReminderLog) error {
	return r.db.Create(reminderLog).Error
}

func (r *reminderLogRepository) GetReminderLogsByReminderID(reminderID string) (*[]models.ReminderLog, error) {
	var reminderLogs []models.ReminderLog
	err := r.db.Where("reminder_id = ?", reminderID).Find(&reminderLogs).Error
	return &reminderLogs, err
}
