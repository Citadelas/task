
# Citadelas/task

**Task Management Service** implemented in Go using gRPC and PostgreSQL.

## Overview

This repository contains a simple task management service with the following features:
- Create, retrieve, update, and delete tasks
- Change task status
- Store data in PostgreSQL
- Manage database migrations via golang-migrate
- Containerized deployment using Docker and Docker Compose

## Technologies

- **Language:** Go 1.24
- **API:** gRPC with Protocol Buffers
- **Database:** PostgreSQL
- **Migrations:** golang-migrate
- **Validation:** go-playground/validator
- **Logging:** Go’s slog package
- **Containerization:** Docker, Docker Compose

## Quick Start

1. **Clone the repository**

git clone https://github.com/Citadelas/task.git
cd task

2. **Create configuration file**  
   Place a file named `local.yaml` in the `config` directory or set the `CONFIG_PATH` environment variable. Example `local.yaml`:

env: "local"
storage_path: "postgres://postgres:postgres@postgres-db:5432/task?sslmode=disable"
grpc:
port: 44044
timeout: 5s

3. **Run with Docker Compose**

docker-compose up --build
- **postgres-db:** PostgreSQL database
- **db-migrate:** Applies database migrations
- **task-app:** Task management service

4. **Access the service**  
   Use a gRPC client to connect to `localhost:44044`.

## Configuration

Configuration parameters are loaded from a YAML file or the `CONFIG_PATH` environment variable:

- `env` – Environment identifier: `local`, `dev`, or `prod`
- `storage_path` – PostgreSQL connection string
- `grpc.port` – Service port
- `grpc.timeout` – gRPC request timeout

## Project Structure


.
├── cmd  
│   └── task            # Entry point (main.go)  
├── internal  
│   ├── config         # Configuration loader  
│   ├── domain  
│   │   └── models     # Domain models  
│   ├── grpc           # gRPC server implementation and validation  
│   ├── services       # Business logic for tasks  
│   ├── storage        # PostgreSQL storage implementation  
│   └── lib/logger     # Logging utilities  
├── migrations         # Database migration scripts  
├── Dockerfile         # Docker build instructions  
├── docker-compose.yml # Docker Compose configuration  
├── go.mod             # Go module file  
└── README.md          # This file


## API Methods

- **CreateTask**  
  Create a new task with title, description, and priority.
- **GetTask**  
  Retrieve a task by its ID.
- **UpdateTask**  
  Update one or more fields of an existing task.
- **DeleteTask**  
  Remove a task by its ID.
- **UpdateStatus**  
  Change the status of an existing task.

### Allowed Values

- **Priority:** `LOW`, `MEDIUM`, `HIGH`
- **Status:** `BACKLOG`, `NEW`, `IN_PROGRESS`, `DONE`

## Error Handling and Validation

- All input fields are validated (length constraints, required fields, enum values)
- Custom errors: `ErrTaskNotFound`, `ErrInputTooLong`
- Errors are mapped to appropriate gRPC status codes

## Recommended Enhancements

- Add authentication and authorization
- Implement unit and integration tests
- Introduce pagination and filtering for task listings
- Provide health checks and metrics
- Expand documentation (e.g., OpenAPI definitions, usage examples)