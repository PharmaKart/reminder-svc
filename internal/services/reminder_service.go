package services

import (
	"time"

	"github.com/PharmaKart/reminder-svc/internal/models"
	"github.com/PharmaKart/reminder-svc/internal/repositories"
	"github.com/google/uuid"
)

type ReminderService interface {
	ScheduleReminder(customerID, orderID string, reminderDate string) error
	GetPendingReminders() ([]models.Reminder, error)
	ListReminders(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, error)
	ListCustomerReminders(customerID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, error)
	ListReminderLogs(reminderID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.ReminderLog, error)
	UpdateReminder(reminderID string, orderID string, reminderDate string) error
	DeleteReminder(reminderID string) error
	ToggleReminder(reminderID string) error
}

type reminderService struct {
	reminderRepo    repositories.ReminderRepository
	reminderLogRepo repositories.ReminderLogRepository
}

func NewReminderService(reminderRepo repositories.ReminderRepository, reminderLogRepo repositories.ReminderLogRepository) ReminderService {
	return &reminderService{
		reminderRepo:    reminderRepo,
		reminderLogRepo: reminderLogRepo,
	}
}

func (s *reminderService) ScheduleReminder(customerID, orderID string, reminderDate string) error {
	customer_id, err := uuid.Parse(customerID)
	order_id, err := uuid.Parse(orderID)
	reminder_date, err := time.Parse(time.RFC3339, reminderDate)
	if err != nil {
		return err
	}

	reminder := &models.Reminder{
		CustomerID:   customer_id,
		OrderID:      order_id,
		ReminderDate: reminder_date,
	}
	return s.reminderRepo.ScheduleReminder(reminder)
}

func (s *reminderService) GetPendingReminders() ([]models.Reminder, error) {
	return s.reminderRepo.GetPendingReminders()
}

func (s *reminderService) ListReminders(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, error) {
	return s.reminderRepo.ListReminders(page, limit, sortBy, sortOrder, filter, filterValue)
}

func (s *reminderService) ListCustomerReminders(customerID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, error) {
	return s.reminderRepo.ListCustomerReminders(customerID, page, limit, sortBy, sortOrder, filter, filterValue)
}

func (s *reminderService) UpdateReminder(reminderID string, orderID string, reminderDate string) error {
	reminder_id, err := uuid.Parse(reminderID)
	order_id, err := uuid.Parse(orderID)
	reminder_date, err := time.Parse(time.RFC3339, reminderDate)
	if err != nil {
		return err
	}

	reminder := &models.Reminder{
		ID:           reminder_id,
		OrderID:      order_id,
		ReminderDate: reminder_date,
	}
	return s.reminderRepo.UpdateReminder(reminder)
}

func (s *reminderService) DeleteReminder(reminderID string) error {
	return s.reminderRepo.DeleteReminder(reminderID)
}

func (s *reminderService) ToggleReminder(reminderID string) error {
	return s.reminderRepo.ToggleReminder(reminderID)
}

func (s *reminderService) ListReminderLogs(reminderID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.ReminderLog, error) {
	return s.reminderLogRepo.ListReminderLogs(reminderID, page, limit, sortBy, sortOrder, filter, filterValue)
}
