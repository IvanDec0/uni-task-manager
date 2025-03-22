// Package input defines the input ports (driving ports) for the hexagonal architecture.
// These interfaces define how external actors interact with the domain layer.
package input

import (
	"context"

	"uni-task-manager/internal/domain/models"
)

// TaskService defines the primary port for task-related business operations.
// This interface represents the API through which the application core can be used.
type TaskService interface {
	// CreateTask creates a new task with input validation and business rules
	// Returns an error if the task data is invalid or if the referenced course doesn't exist
	CreateTask(ctx context.Context, task *models.Task) error

	// UpdateTask modifies an existing task with input validation and business rules
	// Returns an error if the task doesn't exist or if the updated data is invalid
	UpdateTask(ctx context.Context, task *models.Task) error

	// GetTask retrieves a specific task by its ID
	// Returns ErrTaskNotFound if the task doesn't exist
	GetTask(ctx context.Context, id int64) (*models.Task, error)

	// GetAllTasks retrieves all tasks in the system
	GetAllTasks(ctx context.Context) ([]models.Task, error)

	// DeleteTask removes a task from the system
	// Returns ErrTaskNotFound if the task doesn't exist
	DeleteTask(ctx context.Context, id int64) error
}

// CourseService defines the primary port for course-related business operations.
// This interface represents the API through which the application core can be used.
type CourseService interface {
	// CreateCourse creates a new course with input validation and business rules
	// Returns an error if the course data is invalid
	CreateCourse(ctx context.Context, course *models.Course) error

	// UpdateCourse modifies an existing course with input validation and business rules
	// Returns an error if the course doesn't exist or if the updated data is invalid
	UpdateCourse(ctx context.Context, course *models.Course) error

	// GetCourse retrieves a specific course by its ID
	// Returns ErrCourseNotFound if the course doesn't exist
	GetCourse(ctx context.Context, id int64) (*models.Course, error)

	// GetAllCourses retrieves all courses in the system
	GetAllCourses(ctx context.Context) ([]models.Course, error)

	// DeleteCourse removes a course from the system
	// Returns ErrCourseNotFound if the course doesn't exist
	DeleteCourse(ctx context.Context, id int64) error
}
