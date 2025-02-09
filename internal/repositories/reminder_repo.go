package repositories

import (
	"github.com/PharmaKart/reminder-svc/internal/models"
	"gorm.io/gorm"
)

type ReminderRepository interface {
	ScheduleReminder(reminder *models.Reminder) error
	GetPendingReminders() ([]models.Reminder, error)
	ListReminders(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, error)
	ListCustomerReminders(customerID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, error)
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

func (r *reminderRepository) GetPendingReminders() ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Where("sent = ?", false).Find(&reminders).Error
	return reminders, err
}

func (r *reminderRepository) ListReminders(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, error) {
	var reminders []models.Reminder
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
	return reminders, err
}

func (r *reminderRepository) ListCustomerReminders(customerID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, error) {
	var reminders []models.Reminder
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
	return reminders, err
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
