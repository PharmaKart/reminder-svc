package handlers

import (
	"context"

	"github.com/PharmaKart/reminder-svc/internal/proto"
	"github.com/PharmaKart/reminder-svc/internal/repositories"
	"github.com/PharmaKart/reminder-svc/internal/services"
	"github.com/PharmaKart/reminder-svc/pkg/utils"
	"github.com/robfig/cron/v3"
)

type ReminderHandler interface {
	ScheduleReminder(ctx context.Context, req *proto.ScheduleReminderRequest) (*proto.ScheduleReminderResponse, error)
	ListReminders(ctx context.Context, req *proto.ListRemindersRequest) (*proto.ListRemindersResponse, error)
	ListCustomerReminders(ctx context.Context, req *proto.ListCustomerRemindersRequest) (*proto.ListRemindersResponse, error)
	UpdateReminder(ctx context.Context, req *proto.UpdateReminderRequest) (*proto.UpdateReminderResponse, error)
	DeleteReminder(ctx context.Context, req *proto.DeleteReminderRequest) (*proto.DeleteReminderResponse, error)
	ToggleReminder(ctx context.Context, req *proto.ToggleReminderRequest) (*proto.ToggleReminderResponse, error)
	ListReminderLogs(ctx context.Context, req *proto.ListReminderLogsRequest) (*proto.ListReminderLogsResponse, error)
}

type reminderHandler struct {
	proto.UnimplementedReminderServiceServer
	reminderService services.ReminderService
}

func NewReminderHandler(reminderRepo repositories.ReminderRepository, reminderLogRepo repositories.ReminderLogRepository) *reminderHandler {
	return &reminderHandler{
		reminderService: services.NewReminderService(reminderRepo, reminderLogRepo),
	}
}

func (h *reminderHandler) ScheduleReminder(ctx context.Context, req *proto.ScheduleReminderRequest) (*proto.ScheduleReminderResponse, error) {
	err := h.reminderService.ScheduleReminder(req.CustomerId, req.OrderId, req.ProductId, req.ReminderDate)
	if err != nil {
		return nil, err
	}

	return &proto.ScheduleReminderResponse{}, nil
}

func (h *reminderHandler) ListReminders(ctx context.Context, req *proto.ListRemindersRequest) (*proto.ListRemindersResponse, error) {
	reminders, total, err := h.reminderService.ListReminders(req.Page, req.Limit, req.SortBy, req.SortOrder, req.Filter, req.FilterValue)
	if err != nil {
		return nil, err
	}

	protoReminders := make([]*proto.Reminder, len(reminders))
	for i, reminder := range reminders {
		protoReminders[i] = &proto.Reminder{
			Id:           reminder.ID.String(),
			CustomerId:   reminder.CustomerID.String(),
			OrderId:      reminder.OrderID.String(),
			ReminderDate: reminder.ReminderDate.Format("2006-01-02"),
			LastSentAt:   reminder.LastSentAt.Format("2006-01-02"),
			Enabled:      reminder.Enabled,
		}
	}

	return &proto.ListRemindersResponse{
		Reminders: protoReminders,
		Total:     total,
		Page:      req.Page,
		Limit:     req.Limit,
	}, nil
}

func (h *reminderHandler) ListCustomerReminders(ctx context.Context, req *proto.ListCustomerRemindersRequest) (*proto.ListRemindersResponse, error) {
	reminders, total, err := h.reminderService.ListCustomerReminders(req.CustomerId, req.Page, req.Limit, req.SortBy, req.SortOrder, req.Filter, req.FilterValue)
	if err != nil {
		return nil, err
	}

	protoReminders := make([]*proto.Reminder, len(reminders))
	for i, reminder := range reminders {
		protoReminders[i] = &proto.Reminder{
			Id:           reminder.ID.String(),
			CustomerId:   reminder.CustomerID.String(),
			OrderId:      reminder.OrderID.String(),
			ReminderDate: reminder.ReminderDate.Format("2006-01-02"),
			LastSentAt:   reminder.LastSentAt.Format("2006-01-02"),
			Enabled:      reminder.Enabled,
		}
	}

	return &proto.ListRemindersResponse{
		Reminders: protoReminders,
		Total:     total,
		Page:      req.Page,
		Limit:     req.Limit,
	}, nil
}

func (h *reminderHandler) UpdateReminder(ctx context.Context, req *proto.UpdateReminderRequest) (*proto.UpdateReminderResponse, error) {
	err := h.reminderService.UpdateReminder(req.ReminderId, req.OrderId, req.ReminderDate)
	if err != nil {
		return nil, err
	}

	return &proto.UpdateReminderResponse{}, nil
}

func (h *reminderHandler) DeleteReminder(ctx context.Context, req *proto.DeleteReminderRequest) (*proto.DeleteReminderResponse, error) {
	err := h.reminderService.DeleteReminder(req.ReminderId)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteReminderResponse{}, nil
}

func (h *reminderHandler) ToggleReminder(ctx context.Context, req *proto.ToggleReminderRequest) (*proto.ToggleReminderResponse, error) {
	err := h.reminderService.ToggleReminder(req.ReminderId)
	if err != nil {
		return nil, err
	}

	return &proto.ToggleReminderResponse{}, nil
}

func (h *reminderHandler) ListReminderLogs(ctx context.Context, req *proto.ListReminderLogsRequest) (*proto.ListReminderLogsResponse, error) {
	reminderLogs, total, err := h.reminderService.ListReminderLogs(req.ReminderId, req.Page, req.Limit, req.SortBy, req.SortOrder, req.Filter, req.FilterValue)
	if err != nil {
		return nil, err
	}

	protoReminderLogs := make([]*proto.ReminderLog, len(reminderLogs))
	for i, reminderLog := range reminderLogs {
		protoReminderLogs[i] = &proto.ReminderLog{
			Id:         reminderLog.ID.String(),
			ReminderId: reminderLog.ReminderID.String(),
			OrderId:    reminderLog.OrderID.String(),
			Status:     reminderLog.Status,
			CreatedAt:  reminderLog.CreatedAt.Format("2006-01-02"),
		}
	}

	return &proto.ListReminderLogsResponse{
		Logs:  protoReminderLogs,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}, nil
}

func (h *reminderHandler) SendReminders() {
	c := cron.New()

	_, err := c.AddFunc("0 0 * * *", func() {
		h.reminderService.GetPendingReminders()
	})

	if err != nil {
		utils.Error("Failed to schedule reminder job", map[string]interface{}{
			"error": err,
		})
	}

	c.Start()
}
