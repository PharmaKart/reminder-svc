package repositories

import (
	"fmt"
	"time"

	"github.com/PharmaKart/reminder-svc/internal/models"
	"github.com/PharmaKart/reminder-svc/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReminderRepository interface {
	GetReminderCustomer(reminderID string) (string, error)
	ScheduleReminder(reminder *models.Reminder) error
	GetPendingReminders() ([]ReminderWithCustomer, error)
	ListReminders(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, int32, error)
	ListCustomerReminders(customerID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, int32, error)
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

func (r *reminderRepository) ListReminders(page, limit int32, sortBy, sortOrder, filter, filterValue string) ([]models.Reminder, int32, error) {
	var reminders []models.Reminder
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

	err := query.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&reminders).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	err = query.Model(&models.Reminder{}).Count(&total).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	return reminders, int32(total), nil
}

func (r *reminderRepository) ListCustomerReminders(customerID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, int32, error) {
	var reminders []models.Reminder
	var total int64

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	query := r.db.Where("customer_id = ?", customerID)
	if filter != "" && filterValue != "" {
		query = query.Where(filter+" = ?", filterValue)
	}

	if sortBy != "" {
		if sortOrder == "" {
			sortOrder = "asc"
		}
		query = query.Order(sortBy + " " + sortOrder)
	}

	err := query.Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&reminders).Error
	if err != nil {
		return nil, 0, errors.NewInternalError(err)
	}

	err = r.db.Model(&models.Reminder{}).Where("customer_id = ?", customerID).Count(&total).Error
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
