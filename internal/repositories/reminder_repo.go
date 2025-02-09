package repositories

import (
	"time"

	"github.com/PharmaKart/reminder-svc/internal/models"
	"gorm.io/gorm"
)

type ReminderRepository interface {
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
	return r.db.Create(reminder).Error
}

type ReminderWithCustomer struct {
	Reminder models.Reminder
	Email    string
	Phone    *string
	Product  string
}

func (r *reminderRepository) GetPendingReminders() ([]ReminderWithCustomer, error) {
	var results []ReminderWithCustomer

	// Fetch reminders and join with customers to get email and phone
	err := r.db.
		Table("reminders").
		Select("reminders.*, customers.email, customers.phone").
		Joins("JOIN customers ON customers.id = reminders.customer_id").
		Joins("JOIN products ON products.id = reminders.product_id").
		Where("reminder_date <= ? AND enabled = ?", time.Now(), true).
		Scan(&results).Error

	return results, err
}

func (r *reminderRepository) ListReminders(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, int32, error) {
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
		return nil, 0, err
	}

	err = r.db.Model(&models.Reminder{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return reminders, int32(total), err
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
		return nil, 0, err
	}

	err = r.db.Model(&models.Reminder{}).Where("customer_id = ?", customerID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return reminders, int32(total), err
}

func (r *reminderRepository) UpdateReminder(reminder *models.Reminder) error {
	return r.db.Save(reminder).Error
}

func (r *reminderRepository) DeleteReminder(reminderID string) error {
	return r.db.Where("id = ?", reminderID).Delete(&models.Reminder{}).Error
}

func (r *reminderRepository) ToggleReminder(reminderID string) error {
	var reminder models.Reminder
	err := r.db.Where("id = ?", reminderID).First(&reminder).Error
	if err != nil {
		return err
	}

	reminder.Enabled = !reminder.Enabled
	return r.db.Save(&reminder).Error
}
