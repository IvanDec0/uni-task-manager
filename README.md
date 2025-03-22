# University Task Manager

A Go application for managing university tasks and courses using SQLite for persistence.

## Features

- Create, view, edit, and delete university tasks
- Organize tasks by courses
- Set priorities and deadlines for tasks
- Track task status (pending, in progress, completed)
- Responsive web interface using Bootstrap
- REST API endpoints for programmatic access

## Requirements

- Go 1.16+
- SQLite3

## Installation

1. Clone the repository:

   ```
   git clone https://github.com/IvanDec0/uni-task-manager.git
   cd uni-task-manager
   ```

2. Build and run:

   ```
   go build -o uni-task-manager ./cmd/api
   ./uni-task-manager
   ```

3. Access the application:
   Open your web browser and navigate to `http://localhost:8080`

## Project Structure

- `cmd/api`: Main application entry point
- `internal/database`: SQLite database operations
- `internal/handlers`: HTTP request handlers
- `internal/models`: Data models
- `web/templates`: HTML templates for the web interface

## API Endpoints

The application provides the following RESTful API endpoints:

- `GET /api/tasks`: Get all tasks
- `GET /api/tasks/{id}`: Get a specific task
- `POST /api/tasks`: Create a new task
- `PUT /api/tasks/{id}`: Update an existing task
- `DELETE /api/tasks/{id}`: Delete a task

## License

MIT
