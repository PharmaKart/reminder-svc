package repositories

import (
	"github.com/PharmaKart/reminder-svc/internal/models"
	"github.com/PharmaKart/reminder-svc/pkg/errors"
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
	if err := r.db.Create(reminderLog).Error; err != nil {
		return errors.NewInternalError(err)
	}
	return nil
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

	// Base query that filters by reminderID
	query := r.db.Where("reminder_id = ?", reminderID)

	// Add additional filters if provided
	if filter != "" && filterValue != "" {
		query = query.Where(filter+" = ?", filterValue)
	}

	// Add ordering
	if sortBy != "" {
		if sortOrder == "" {
			sortOrder = "asc"
		}
		query = query.Order(sortBy + " " + sortOrder)
	}

	// Execute the paginated query
	err := query.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&reminderLogs).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	// Get total count for pagination
	err = query.Model(&models.ReminderLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	// Return empty slice with 0 total instead of an error if no logs found
	if len(reminderLogs) == 0 {
		return []models.ReminderLog{}, 0, nil
	}

	return reminderLogs, int32(total), nil
}
