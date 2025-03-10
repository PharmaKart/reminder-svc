package repositories

import (
	"strings"

	"github.com/PharmaKart/reminder-svc/internal/models"
	"github.com/PharmaKart/reminder-svc/pkg/errors"
	"github.com/PharmaKart/reminder-svc/pkg/utils"
	"gorm.io/gorm"
)

type ReminderLogRepository interface {
	CreateReminderLog(reminderLog *models.ReminderLog) error
	ListReminderLogs(reminderID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.ReminderLog, int32, error)
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

func (r *reminderLogRepository) ListReminderLogs(reminderID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.ReminderLog, int32, error) {
	var reminderLogs []models.ReminderLog
	var total int64

	allowedColumns := utils.GetModelColumns(&models.ReminderLog{})

	allowedOperators := map[string]string{
		"eq":      "=",           // Equal to
		"neq":     "!=",          // Not equal to
		"gt":      ">",           // Greater than
		"gte":     ">=",          // Greater than or equal to
		"lt":      "<",           // Less than
		"lte":     "<=",          // Less than or equal to
		"like":    "LIKE",        // LIKE for pattern matching
		"ilike":   "ILIKE",       // Case insensitive LIKE (for PostgreSQL)
		"in":      "IN",          // IN for multiple values
		"null":    "IS NULL",     // IS NULL check
		"notnull": "IS NOT NULL", // IS NOT NULL check
	}

	query := r.db.Model(&models.ReminderLog{}).Where("reminder_id = ?", reminderID)

	if filter != (models.Filter{}) {
		if _, allowed := allowedColumns[filter.Column]; !allowed {
			return nil, 0, errors.NewBadRequestError("invalid filter column: " + filter.Column)
		}

		op, allowed := allowedOperators[filter.Operator]
		if !allowed {
			return nil, 0, errors.NewBadRequestError("invalid filter operator: " + filter.Operator)
		}

		switch filter.Operator {
		case "like", "ilike":
			query = query.Where(filter.Column+" "+op+" ?", "%"+filter.Value+"%")
		case "in":
			values := strings.Split(filter.Value, ",")
			query = query.Where(filter.Column+" "+op+" (?)", values)
		case "null", "notnull":
			query = query.Where(filter.Column + " " + op)
		default:
			query = query.Where(filter.Column+" "+op+" ?", filter.Value)
		}
	}

	if sortBy != "" {
		if _, allowed := allowedColumns[sortBy]; !allowed {
			return nil, 0, errors.NewBadRequestError("invalid sort column: " + sortBy)
		}

		sortOrder = strings.ToLower(sortOrder)
		if sortOrder != "asc" && sortOrder != "desc" {
			sortOrder = "asc"
		}

		query = query.Order(sortBy + " " + sortOrder)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	if limit > 0 {
		offset := max(int((page-1)*limit), 0)
		query = query.Offset(offset).Limit(int(limit))
	}

	err = query.Find(&reminderLogs).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	return reminderLogs, int32(total), nil
}
