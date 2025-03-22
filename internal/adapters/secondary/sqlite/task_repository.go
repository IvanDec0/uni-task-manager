// Package sqlite provides implementations of the repository interfaces using SQLite as the storage backend.
package sqlite

import (
	"context"
	"database/sql"
	"time"

	"uni-task-manager/internal/domain/models"
)

// TaskRepository implements output.TaskRepository interface using SQLite as the storage backend.
// It handles all task-related database operations and mapping between domain models and database rows.
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new instance of TaskRepository with the provided database connection.
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// GetAll retrieves all tasks from the database, ordered by due date.
// It maps the database rows to domain Task objects.
func (r *TaskRepository) GetAll(ctx context.Context) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, due_date, priority, status, course_id, created_at, updated_at
		FROM tasks
		ORDER BY due_date ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var dueDate, createdAt, updatedAt string
		var status string

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&dueDate,
			&task.Priority,
			&status,
			&task.CourseID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		task.Status = models.TaskStatus(status)
		task.DueDate, _ = time.Parse(time.RFC3339, dueDate)
		task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetByID retrieves a specific task by its ID from the database.
// Returns nil if no task is found with the given ID.
func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	var task models.Task
	var dueDate, createdAt, updatedAt string
	var status string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, description, due_date, priority, status, course_id, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&dueDate,
		&task.Priority,
		&status,
		&task.CourseID,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	task.Status = models.TaskStatus(status)
	task.DueDate, _ = time.Parse(time.RFC3339, dueDate)
	task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &task, nil
}

// Create persists a new task in the database.
// It sets the ID field of the task object with the generated ID.
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO tasks (title, description, due_date, priority, status, course_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.Title,
		task.Description,
		task.DueDate.Format(time.RFC3339),
		task.Priority,
		string(task.Status),
		task.CourseID,
		task.CreatedAt.Format(time.RFC3339),
		task.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	task.ID = id
	return nil
}

// Update modifies an existing task in the database.
// All fields except CreatedAt can be updated.
func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tasks
		SET title = ?, description = ?, due_date = ?, priority = ?, status = ?, course_id = ?, updated_at = ?
		WHERE id = ?
	`,
		task.Title,
		task.Description,
		task.DueDate.Format(time.RFC3339),
		task.Priority,
		string(task.Status),
		task.CourseID,
		task.UpdatedAt.Format(time.RFC3339),
		task.ID,
	)

	return err
}

// Delete removes a task from the database by its ID.
func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = ?", id)
	return err
}

// GetByCourseID retrieves all tasks associated with a specific course.
// Tasks are ordered by due date.
func (r *TaskRepository) GetByCourseID(ctx context.Context, courseID int64) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, description, due_date, priority, status, course_id, created_at, updated_at
		FROM tasks
		WHERE course_id = ?
		ORDER BY due_date ASC
	`, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var dueDate, createdAt, updatedAt string
		var status string

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&dueDate,
			&task.Priority,
			&status,
			&task.CourseID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		task.Status = models.TaskStatus(status)
		task.DueDate, _ = time.Parse(time.RFC3339, dueDate)
		task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
