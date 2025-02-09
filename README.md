# Reminder Service

The **Reminder Service** is a critical component of the Pharmakart platform, responsible for sending automated refill reminders to customers. It ensures that users receive timely notifications before their prescriptions run out, helping them maintain uninterrupted medication schedules.

---

## Table of Contents
1. [Overview](#overview)
2. [Features](#features)
3. [Prerequisites](#prerequisites)
4. [Setup and Installation](#setup-and-installation)
5. [Running the Service](#running-the-service)
6. [Environment Variables](#environment-variables)
7. [Contributing](#contributing)
8. [License](#license)

---

## Overview

The Reminder Service handles:
- Scheduling and sending refill reminders via email or SMS.
- Tracking reminder statuses (sent or pending).
- Role-based access control (only customers receive reminders).

It is built using **gRPC** for communication and **PostgreSQL** for storing reminders.

---

## Features

- **Automated Refill Reminders**:
  - Sends reminders via email or SMS before a prescription runs out.
  - Uses AWS Simple Notification Service (SNS) or AWS Simple Queue Service (SQS) for message handling.
- **Reminder Tracking**:
  - Stores reminders in the database with timestamps.
  - Ensures reminders are sent only when necessary.
- **Integration with Order Service**:
  - Automatically schedules reminders when a prescription-based order is placed.
- **Role-Based Access Control**:
  - Customers receive reminders, while admins can monitor logs.

---

## Prerequisites

Before setting up the service, ensure you have the following installed:
- **Docker**
- **Go** (for building and running the service)
- **Protobuf Compiler** (`protoc`) for generating gRPC/protobuf files
- **AWS CLI** (if using SNS/SQS for notifications)

---

## Setup and Installation

### 1. Clone the Repository
Clone the repository and navigate to the reminder service directory:
```bash
git clone https://github.com/PharmaKart/reminder-svc.git
cd reminder-svc
```

### 2. Generate Protobuf Files
Generate the protobuf files using the provided `Makefile`:
```bash
make proto
```

### 3. Install Dependencies
Run the following command to ensure all dependencies are installed:
```bash
go mod tidy
```

### 4. Build the Service
To build the service, run:
```bash
make build
```

---

## Running the Service

### Option 1: Run Using Docker
To run the service using Docker, execute:
```bash
docker run -p 50055:50055 pharmakart/reminder-svc
```

### Option 2: Run Using Makefile
To run the service directly using Go, execute:
```bash
make run
```

The service will be available at:
- **gRPC**: `localhost:50055`

---

## Environment Variables

The service requires the following environment variables. Create a `.env` file in the `reminder-svc` directory with the following:

```env
PORT=50055
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=pharmakartdb
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
AWS_REGION=ca-central-1
SNS_TOPIC_ARN=your-sns-topic-arn
SQS_QUEUE_URL=your-sqs-queue-url
```

---

## Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Submit a pull request with a detailed description of your changes.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Support

For any questions or issues, please open an issue in the repository or contact the maintainers.

