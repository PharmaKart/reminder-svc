package repositories

import (
	"fmt"
	"strings"
	"time"

	"github.com/PharmaKart/reminder-svc/internal/models"
	"github.com/PharmaKart/reminder-svc/pkg/errors"
	"github.com/PharmaKart/reminder-svc/pkg/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReminderRepository interface {
	GetReminderCustomer(reminderID string) (string, error)
	ScheduleReminder(reminder *models.Reminder) error
	GetPendingReminders() ([]ReminderWithCustomer, error)
	ListReminders(filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Reminder, int32, error)
	ListCustomerReminders(customerID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Reminder, int32, error)
	UpdateReminder(reminder *models.Reminder) error
	DeleteReminder(reminderID string) error
	ToggleReminder(reminderID string) error
}

type reminderRepository struct {
	db *gorm.DB
}

func NewReminderRepository(db *gorm.DB) ReminderRepository {
	return &reminderRepository{db}
}

func (r *reminderRepository) ScheduleReminder(reminder *models.Reminder) error {
	if err := r.db.Create(reminder).Error; err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

type ReminderWithCustomer struct {
	Reminder models.Reminder
	Email    string
	Phone    *string
	Product  string
}

func (r *reminderRepository) GetReminderCustomer(reminderID string) (string, error) {
	var customerID uuid.UUID
	err := r.db.Table("reminders").Select("customer_id").Where("id = ?", reminderID).Row().Scan(&customerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.NewNotFoundError(fmt.Sprintf("Reminder with ID '%s' not found", reminderID))
		}
		return "", errors.NewInternalError(err)
	}
	return customerID.String(), nil
}

func (r *reminderRepository) GetPendingReminders() ([]ReminderWithCustomer, error) {
	var results []ReminderWithCustomer

	err := r.db.
		Table("reminders").
		Select("reminders.*, customers.email, customers.phone, products.name as product").
		Joins("JOIN customers ON customers.id = reminders.customer_id").
		Joins("JOIN products ON products.id = reminders.product_id").
		Where("reminder_date <= ? AND enabled = ?", time.Now(), true).
		Scan(&results).Error

	if err != nil {
		return nil, errors.NewInternalError(err)
	}
	return results, nil
}

func (r *reminderRepository) ListReminders(filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Reminder, int32, error) {
	var reminders []models.Reminder
	var total int64

	allowedColumns := utils.GetModelColumns(&models.Reminder{})

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

	query := r.db.Model(&models.Reminder{})

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

	err = query.Find(&reminders).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	return reminders, int32(total), nil
}

func (r *reminderRepository) ListCustomerReminders(customerID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Reminder, int32, error) {
	var reminders []models.Reminder
	var total int64

	allowedColumns := utils.GetModelColumns(&models.Reminder{})

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

	query := r.db.Model(&models.Reminder{}).Where("customer_id = ?", customerID)

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

	err = query.Find(&reminders).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	return reminders, int32(total), nil
}

func (r *reminderRepository) UpdateReminder(reminder *models.Reminder) error {
	if err := r.db.Save(reminder).Error; err != nil {
		return errors.NewInternalError(err)
	}
	return nil
}

func (r *reminderRepository) DeleteReminder(reminderID string) error {
	result := r.db.Where("id = ?", reminderID).Delete(&models.Reminder{})
	if result.Error != nil {
		return errors.NewInternalError(result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewNotFoundError(fmt.Sprintf("Reminder with ID '%s' not found", reminderID))
	}
	return nil
}

func (r *reminderRepository) ToggleReminder(reminderID string) error {
	var reminder models.Reminder
	if err := r.db.Where("id = ?", reminderID).First(&reminder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError(fmt.Sprintf("Reminder with ID '%s' not found", reminderID))
		}
		return errors.NewInternalError(err)
	}

	reminder.Enabled = !reminder.Enabled
	if err := r.db.Save(&reminder).Error; err != nil {
		return errors.NewInternalError(err)
	}

	return nil
}
