# University Task Manager 📚✅

A sophisticated Go application for managing university tasks and courses, built using clean architecture principles. This application helps students and professors organize academic work efficiently with a modern web interface and REST API.

## ✨ Features

- **Task Management**

  - Create, view, edit, and delete academic tasks
  - Set priorities (1-5 scale)
  - Track deadlines with due dates
  - Monitor task status (pending, in progress, completed)
  - Associate tasks with specific courses

- **Course Management**

  - Manage course information
  - Track professor details
  - Organize tasks by course

- **User Interface**

  - Clean, responsive web interface using Bootstrap
  - Intuitive task and course management
  - Visual indicators for task priority and status

- **API Support**
  - RESTful API for programmatic access
  - JSON-based data exchange
  - Complete CRUD operations for tasks

## 🏗️ Architecture

The application follows hexagonal architecture (ports and adapters) principles:

```
internal/
├── domain/          # Core business logic and entities
│   ├── models/      # Domain models (Task, Course)
│   └── services/    # Business logic implementation
├── ports/           # Interface definitions
│   ├── input/       # Primary ports (service interfaces)
│   └── output/      # Secondary ports (repository interfaces)
└── adapters/        # Interface implementations
    ├── primary/     # Driving adapters (HTTP handlers)
    └── secondary/   # Driven adapters (SQLite repositories)
```

### Key Components

- **Domain Layer**: Contains core business logic and entities
- **Ports Layer**: Defines interfaces for the application core
- **Adapters Layer**: Implements interfaces for external interactions
- **Web Interface**: Bootstrap-based responsive UI
- **Database**: SQLite for simple, file-based storage

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or higher
- SQLite3
- Git

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/IvanDec0/uni-task-manager.git
   cd uni-task-manager
   ```

2. Build the application:

   ```bash
   go build -o uni-task-manager ./cmd/api
   ```

3. Run the server:

   ```bash
   ./uni-task-manager
   ```

4. Access the application:
   - Web Interface: Open http://localhost:8080 in your browser
   - API: Send requests to http://localhost:8080/api/...

## 🔧 API Endpoints

### Tasks

- `GET /api/tasks` - List all tasks
- `GET /api/tasks/{id}` - Get task details
- `POST /api/tasks` - Create a new task
- `PUT /api/tasks/{id}` - Update a task
- `DELETE /api/tasks/{id}` - Delete a task

### Example Request (Create Task)

```json
POST /api/tasks
{
  "title": "Final Project",
  "description": "Complete the semester project",
  "due_date": "2025-04-15T23:59:59Z",
  "priority": 4,
  "course_id": 1
}
```

## 🧪 Testing

Run the test suite:

```bash
go test ./...
```

## 🔍 Design Decisions

- **Hexagonal Architecture**: Ensures separation of concerns and testability
- **SQLite**: Chosen for simplicity and zero-configuration setup
- **Bootstrap**: Provides responsive design with minimal custom CSS
- **Go Modules**: Modern dependency management
- **RESTful API**: Standard-compliant web API design
