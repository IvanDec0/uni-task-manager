// Package services implements the core business logic of the application.
// It contains the implementation of the input ports (service interfaces) and encapsulates
// all domain rules and validation logic.
package services

import (
	"context"
	"errors"
	"time"

	"uni-task-manager/internal/domain/models"
	"uni-task-manager/internal/ports/input"
	"uni-task-manager/internal/ports/output"
)

// Domain-specific errors that can be returned by the TaskService
var (
	// ErrInvalidTaskPriority indicates that the task priority is outside the valid range (1-5)
	ErrInvalidTaskPriority = errors.New("task priority must be between 1 and 5")

	// ErrInvalidDueDate indicates that the task's due date is in the past
	ErrInvalidDueDate = errors.New("due date must be in the future")

	// ErrTaskNotFound indicates that the requested task does not exist
	ErrTaskNotFound = errors.New("task not found")
)

// Verify TaskService implements input.TaskService interface at compile time
var _ input.TaskService = (*TaskService)(nil)

// TaskService implements the task-related business logic and orchestrates
// interactions between the domain model and storage layer.
type TaskService struct {
	taskRepo   output.TaskRepository
	courseRepo output.CourseRepository
}

// NewTaskService creates a new instance of TaskService with the required dependencies.
func NewTaskService(taskRepo output.TaskRepository, courseRepo output.CourseRepository) *TaskService {
	return &TaskService{
		taskRepo:   taskRepo,
		courseRepo: courseRepo,
	}
}

// CreateTask implements input.TaskService.CreateTask.
// It validates the task data and ensures any referenced course exists before creation.
func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	if err := s.validateTask(task); err != nil {
		return err
	}

	if task.CourseID != 0 {
		// Verify course exists if specified
		_, err := s.courseRepo.GetByID(ctx, task.CourseID)
		if err != nil {
			return err
		}
	}

	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now
	if task.Status == "" {
		task.Status = models.TaskStatusPending
	}

	return s.taskRepo.Create(ctx, task)
}

// UpdateTask implements input.TaskService.UpdateTask.
// It validates the updated task data and ensures the task exists before updating.
func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	if err := s.validateTask(task); err != nil {
		return err
	}

	existing, err := s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTaskNotFound
	}

	if task.CourseID != 0 {
		// Verify course exists if specified
		_, err := s.courseRepo.GetByID(ctx, task.CourseID)
		if err != nil {
			return err
		}
	}

	task.CreatedAt = existing.CreatedAt
	task.UpdatedAt = time.Now().UTC()

	return s.taskRepo.Update(ctx, task)
}

// GetTask implements input.TaskService.GetTask.
// It retrieves a specific task by its ID.
func (s *TaskService) GetTask(ctx context.Context, id int64) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

// GetAllTasks implements input.TaskService.GetAllTasks.
// It retrieves all tasks from the repository.
func (s *TaskService) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	return s.taskRepo.GetAll(ctx)
}

// DeleteTask implements input.TaskService.DeleteTask.
// It ensures the task exists before deletion.
func (s *TaskService) DeleteTask(ctx context.Context, id int64) error {
	existing, err := s.taskRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTaskNotFound
	}
	return s.taskRepo.Delete(ctx, id)
}

// validateTask performs validation of task data according to business rules.
// It checks priority range and ensures the due date is in the future.
func (s *TaskService) validateTask(task *models.Task) error {
	if task.Priority < 1 || task.Priority > 5 {
		return ErrInvalidTaskPriority
	}

	if !task.DueDate.IsZero() && task.DueDate.Before(time.Now()) {
		return ErrInvalidDueDate
	}

	return nil
}
