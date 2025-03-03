package repositories

import (
	"github.com/PharmaKart/reminder-svc/internal/models"
	"gorm.io/gorm"
)

type ReminderLogRepository interface {
	CreateReminderLog(reminderLog *models.ReminderLog) error
	ListReminderLogs(reminderID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.ReminderLog, int32, error)
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

func (r *reminderLogRepository) ListReminderLogs(reminderID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.ReminderLog, int32, error) {
	var reminderLogs []models.ReminderLog
	var total int64

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := r.db
	if filter != "" && filterValue != "" {
		query = query.Where(filter+" = ?", filterValue)
	}

	if sortBy != "" {
		if sortOrder == "" {
			sortOrder = "asc"
		}
		query = query.Order(sortBy + " " + sortOrder)
	}

	err := query.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&reminderLogs).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Model(&models.ReminderLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return reminderLogs, int32(total), err
}
