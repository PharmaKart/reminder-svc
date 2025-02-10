package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/PharmaKart/reminder-svc/internal/models"
	"github.com/PharmaKart/reminder-svc/internal/repositories"
	"github.com/PharmaKart/reminder-svc/pkg/config"
	"github.com/PharmaKart/reminder-svc/pkg/utils"
	"github.com/google/uuid"
)

type ReminderService interface {
	ScheduleReminder(customerID, orderID string, productID string, reminderDate string) error
	GetPendingReminders() ([]repositories.ReminderWithCustomer, error)
	ListReminders(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, int32, error)
	ListCustomerReminders(customerID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, int32, error)
	ListReminderLogs(reminderID string, customerId string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.ReminderLog, int32, error)
	UpdateReminder(reminderID string, customerId string, orderID string, reminderDate string) error
	DeleteReminder(reminderID string, customerId string) error
	ToggleReminder(reminderID string, customerId string) error
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
	order_id, err := uuid.Parse(orderID)
	product_id, err := uuid.Parse(productID)
	reminder_date, err := time.Parse(time.RFC3339, reminderDate)
	if err != nil {
		return err
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

func (s *reminderService) ListReminders(page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, int32, error) {
	return s.reminderRepo.ListReminders(page, limit, sortBy, sortOrder, filter, filterValue)
}

func (s *reminderService) ListCustomerReminders(customerID string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.Reminder, int32, error) {
	return s.reminderRepo.ListCustomerReminders(customerID, page, limit, sortBy, sortOrder, filter, filterValue)
}

func (s *reminderService) UpdateReminder(reminderID string, customerId string, orderID string, reminderDate string) error {
	customerID, err := s.reminderRepo.GetReminderCustomer(reminderID)
	if err != nil {
		return err
	}

	if customerId != customerID {
		return errors.New("Access denied")
	}
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

func (s *reminderService) DeleteReminder(reminderID string, customerId string) error {
	customerID, err := s.reminderRepo.GetReminderCustomer(reminderID)
	if err != nil {
		return err
	}

	if customerId != customerID {
		return errors.New("Access denied")
	}

	return s.reminderRepo.DeleteReminder(reminderID)
}

func (s *reminderService) ToggleReminder(reminderID string, customerId string) error {
	customerID, err := s.reminderRepo.GetReminderCustomer(reminderID)
	if err != nil {
		return err
	}

	if customerId != customerID {
		return errors.New("Access denied")
	}

	return s.reminderRepo.ToggleReminder(reminderID)
}

func (s *reminderService) ListReminderLogs(reminderID string, customerId string, page int32, limit int32, sortBy string, sortOrder string, filter string, filterValue string) ([]models.ReminderLog, int32, error) {
	customerID, err := s.reminderRepo.GetReminderCustomer(reminderID)
	if err != nil {
		return nil, 0, err
	}

	if customerId != customerID {
		return nil, 0, errors.New("Access denied")
	}

	return s.reminderLogRepo.ListReminderLogs(reminderID, page, limit, sortBy, sortOrder, filter, filterValue)
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

func (s *reminderService) SendReminderAlert(cfg *config.Config) {
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
	// 	Credentials: credentials.NewStaticCredentials(cfg.AWS_ACCESS_KEY_ID, cfg.AWS_SECRET_ACCESS_KEY, ""),
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
