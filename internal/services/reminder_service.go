package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/PharmaKart/reminder-svc/internal/models"
	"github.com/PharmaKart/reminder-svc/internal/repositories"
	"github.com/PharmaKart/reminder-svc/pkg/config"
	"github.com/PharmaKart/reminder-svc/pkg/errors"
	"github.com/PharmaKart/reminder-svc/pkg/utils"
	"github.com/google/uuid"
)

type ReminderService interface {
	ScheduleReminder(customerID, orderID string, productID string, reminderDate string) error
	GetPendingReminders() ([]repositories.ReminderWithCustomer, error)
	ListReminders(filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Reminder, int32, error)
	ListCustomerReminders(customerID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Reminder, int32, error)
	ListReminderLogs(reminderID string, customerId string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.ReminderLog, int32, error)
	UpdateReminder(reminderID string, customerId string, orderID string, reminderDate string) error
	DeleteReminder(reminderID string, customerId string) error
	ToggleReminder(reminderID string, customerId string) error
	StartReminderService(cfg *config.Config)
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

func (s *reminderService) ScheduleReminder(customerID, orderID string, productID string, reminderDate string) error {
	customer_id, err := uuid.Parse(customerID)
	if err != nil {
		return errors.NewInternalError(err)
	}

	order_id, err := uuid.Parse(orderID)
	if err != nil {
		return errors.NewInternalError(err)
	}

	product_id, err := uuid.Parse(productID)
	if err != nil {
		return errors.NewInternalError(err)
	}

	reminder_date, err := time.Parse(time.RFC3339, reminderDate)
	if err != nil {
		return errors.NewInternalError(err)
	}

	// Check if reminder already exists with same product and customer
	reminderExists, err := s.reminderRepo.ReminderExists(productID, customerID)
	if err != nil {
		return errors.NewInternalError(err)
	}

	if reminderExists {
		return errors.NewConflictError("Reminder already exists for this product")
	}

	reminder := &models.Reminder{
		CustomerID:   customer_id,
		OrderID:      order_id,
		ProductID:    product_id,
		ReminderDate: reminder_date,
	}
	return s.reminderRepo.ScheduleReminder(reminder)
}

func (s *reminderService) GetPendingReminders() ([]repositories.ReminderWithCustomer, error) {
	return s.reminderRepo.GetPendingReminders()
}

func (s *reminderService) ListReminders(filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Reminder, int32, error) {
	return s.reminderRepo.ListReminders(filter, sortBy, sortOrder, page, limit)
}

func (s *reminderService) ListCustomerReminders(customerID string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.Reminder, int32, error) {
	return s.reminderRepo.ListCustomerReminders(customerID, filter, sortBy, sortOrder, page, limit)
}

func (s *reminderService) UpdateReminder(reminderID string, customerId string, orderID string, reminderDate string) error {
	customerID, err := s.reminderRepo.GetReminderCustomer(reminderID)
	if err != nil {
		return err
	}

	if customerId != customerID {
		return errors.NewAuthError("Access denied")
	}
	reminder_id, err := uuid.Parse(reminderID)
	if err != nil {
		return errors.NewInternalError(err)
	}

	order_id, err := uuid.Parse(orderID)
	if err != nil {
		return errors.NewInternalError(err)
	}

	reminder_date, err := time.Parse(time.RFC3339, reminderDate)
	if err != nil {
		return errors.NewInternalError(err)
	}

	reminder := &models.Reminder{
		ID:           reminder_id,
		OrderID:      order_id,
		ReminderDate: reminder_date,
	}
	return s.reminderRepo.UpdateReminder(reminder)
}

func (s *reminderService) DeleteReminder(reminderID string, customerId string) error {
	customerID, err := s.reminderRepo.GetReminderCustomer(reminderID)
	if err != nil {
		return err
	}

	if customerId != customerID {
		return errors.NewAuthError("Access denied")
	}

	return s.reminderRepo.DeleteReminder(reminderID)
}

func (s *reminderService) ToggleReminder(reminderID string, customerId string) error {
	customerID, err := s.reminderRepo.GetReminderCustomer(reminderID)
	if err != nil {
		return err
	}

	if customerId != customerID {
		return errors.NewAuthError("Access denied")
	}

	return s.reminderRepo.ToggleReminder(reminderID)
}

func (s *reminderService) ListReminderLogs(reminderID string, customerId string, filter models.Filter, sortBy string, sortOrder string, page, limit int32) ([]models.ReminderLog, int32, error) {
	customerID, err := s.reminderRepo.GetReminderCustomer(reminderID)
	if err != nil {
		return nil, 0, err
	}

	if customerId != customerID {
		return nil, 0, errors.NewAuthError("Access denied")
	}

	return s.reminderLogRepo.ListReminderLogs(reminderID, filter, sortBy, sortOrder, page, limit)
}

type ReminderMessage struct {
	ReminderID   string `json:"reminder_id"`
	CustomerID   string `json:"customer_id"`
	OrderID      string `json:"order_id"`
	ProductID    string `json:"product_id"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	ReminderDate string `json:"reminder_date"`
}

func (s *reminderService) StartReminderService(cfg *config.Config) {
	// Get pending reminders
	reminders, err := s.GetPendingReminders()
	if err != nil {
		utils.Error("Failed to get pending reminders", map[string]interface{}{
			"error": err,
		})
		return
	}

	// // Initialize AWS session
	// sess, err := session.NewSession(&aws.Config{
	// 	Region: aws.String(cfg.AWS_REGION),
	// })
	// if err != nil {
	// 	utils.Error("Failed to initialize AWS session", map[string]interface{}{
	// 		"error": err,
	// 	})
	// 	return
	// }

	// // Initialize SQS client
	// sqsClient := sqs.New(sess)

	// Iterate over reminders and queue messages
	for _, reminder := range reminders {
		message := ReminderMessage{
			ReminderID:   reminder.Reminder.ID.String(),
			CustomerID:   reminder.Reminder.CustomerID.String(),
			OrderID:      reminder.Reminder.OrderID.String(),
			ProductID:    reminder.Reminder.ProductID.String(),
			Email:        reminder.Email,
			Phone:        "",
			ReminderDate: reminder.Reminder.ReminderDate.Format(time.RFC3339),
		}

		// Include phone number if available
		if reminder.Phone != nil {
			message.Phone = *reminder.Phone
		}

		// Serialize message to JSON
		messageBody, err := json.Marshal(message)
		if err != nil {
			utils.Error("Failed to marshal reminder message", map[string]interface{}{
				"error": err,
			})
			continue
		}

		// // Send message to SQS queue
		// _, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		// 	QueueUrl:    aws.String(cfg.SQS_QUEUE_URL),
		// 	MessageBody: aws.String(string(messageBody)),
		// })
		// if err != nil {
		// 	utils.Error("Failed to send message to SQS", map[string]interface{}{
		// 		"error": err,
		// 	})
		// 	continue
		// }

		utils.Info("Reminder queue", map[string]interface{}{
			"message": string(messageBody),
		})

		utils.Info(fmt.Sprintf("Reminder queued successfully for customer %s", reminder.Email), nil)
	}
}
