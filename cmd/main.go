package main

import (
	"net"

	"github.com/PharmaKart/reminder-svc/internal/handlers"
	"github.com/PharmaKart/reminder-svc/internal/proto"
	"github.com/PharmaKart/reminder-svc/internal/repositories"
	"github.com/PharmaKart/reminder-svc/pkg/config"
	"github.com/PharmaKart/reminder-svc/pkg/utils"
	"google.golang.org/grpc"
)

func main() {
	// Initialize logger
	utils.InitLogger()

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database connection
	db, err := utils.ConnectDB(cfg)
	if err != nil {
		utils.Logger.Fatal("Failed to connect to database", map[string]interface{}{
			"error": err,
		})
	}

	// Initialize repositories
	reminderRepo := repositories.NewReminderRepository(db)
	reminderLogRepo := repositories.NewReminderLogRepository(db)

	// Initialize handlers
	reminderHandler := handlers.NewReminderHandler(reminderRepo, reminderLogRepo)

	// Cron job to send reminders
	go reminderHandler.StartReminderService(cfg)

	// Initialize gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.Port)

	if err != nil {
		utils.Logger.Fatal("Failed to listen", map[string]interface{}{
			"error": err,
		})
	}

	grpcServer := grpc.NewServer()
	proto.RegisterReminderServiceServer(grpcServer, reminderHandler)

	utils.Info("Starting reminder service", map[string]interface{}{
		"port": cfg.Port,
	})

	if err := grpcServer.Serve(lis); err != nil {
		utils.Logger.Fatal("Failed to serve", map[string]interface{}{
			"error": err,
		})
	}
}
