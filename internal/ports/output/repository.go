// Package output defines the output ports (driven ports) for the hexagonal architecture.
// These interfaces define how the domain layer interacts with external systems like databases.
package output

import (
	"context"

	"uni-task-manager/internal/domain/models"
)

// TaskRepository defines the interface for task storage operations.
// This is an output port that must be implemented by storage adapters (e.g., SQLite, PostgreSQL).
type TaskRepository interface {
	// GetAll retrieves all tasks from the storage, ordered by due date
	GetAll(ctx context.Context) ([]models.Task, error)

	// GetByID retrieves a specific task by its unique identifier
	// Returns nil if the task is not found
	GetByID(ctx context.Context, id int64) (*models.Task, error)

	// Create persists a new task in the storage
	// The ID field will be populated with the generated identifier
	Create(ctx context.Context, task *models.Task) error

	// Update modifies an existing task in the storage
	// Returns an error if the task doesn't exist
	Update(ctx context.Context, task *models.Task) error

	// Delete removes a task from the storage
	// Returns an error if the task doesn't exist
	Delete(ctx context.Context, id int64) error

	// GetByCourseID retrieves all tasks associated with a specific course
	GetByCourseID(ctx context.Context, courseID int64) ([]models.Task, error)
}

// CourseRepository defines the interface for course storage operations.
// This is an output port that must be implemented by storage adapters (e.g., SQLite, PostgreSQL).
type CourseRepository interface {
	// GetAll retrieves all courses from the storage, ordered by name
	GetAll(ctx context.Context) ([]models.Course, error)

	// GetByID retrieves a specific course by its unique identifier
	// Returns nil if the course is not found
	GetByID(ctx context.Context, id int64) (*models.Course, error)

	// Create persists a new course in the storage
	// The ID field will be populated with the generated identifier
	Create(ctx context.Context, course *models.Course) error

	// Update modifies an existing course in the storage
	// Returns an error if the course doesn't exist
	Update(ctx context.Context, course *models.Course) error

	// Delete removes a course from the storage
	// Returns an error if the course doesn't exist
	Delete(ctx context.Context, id int64) error
}
