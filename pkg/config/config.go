package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config struct
type Config struct {
	Port                  string
	DBConnString          string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	SQS_QUEUE_URL         string
	SNS_TOPIC_ARN         string
}

// LoadConfig loads the configuration from .env file
func LoadConfig() *Config {
	// Load environment variables from .env file
	if err := godotenv.Overload(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Port:                  getEnv("PORT", "50055"),
		DBConnString:          getDBConnString(),
		AWS_ACCESS_KEY_ID:     getEnv("AWS_ACCESS_KEY_ID", "aws-access-key"),
		AWS_SECRET_ACCESS_KEY: getEnv("AWS_SECRET_ACCESS_KEY", "your-aws-secret-key"),
		AWS_REGION:            getEnv("AWS_REGION", "ca-central-1"),
		SNS_TOPIC_ARN:         getEnv("SNS_TOPIC_ARN", "your-sns-topic-arn"),
		SQS_QUEUE_URL:         getEnv("SQS_QUEUE_URL", "your-sqs-queue-url"),
	}
}

func getDBConnString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "pharmakartdb"),
	)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
