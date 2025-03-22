// Package sqlite provides implementations of the repository interfaces using SQLite as the storage backend.
package sqlite

import (
	"context"
	"database/sql"
	"time"

	"uni-task-manager/internal/domain/models"
)

// CourseRepository implements output.CourseRepository interface using SQLite as the storage backend.
// It handles all course-related database operations and mapping between domain models and database rows.
type CourseRepository struct {
	db *sql.DB
}

// NewCourseRepository creates a new instance of CourseRepository with the provided database connection.
func NewCourseRepository(db *sql.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// GetAll retrieves all courses from the database, ordered by name.
// It maps the database rows to domain Course objects.
func (r *CourseRepository) GetAll(ctx context.Context) ([]models.Course, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, professor, created_at, updated_at
		FROM courses
		ORDER BY name ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		var createdAt, updatedAt string

		if err := rows.Scan(
			&course.ID,
			&course.Name,
			&course.Professor,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		course.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		course.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}

// GetByID retrieves a specific course by its ID from the database.
// Returns nil if no course is found with the given ID.
func (r *CourseRepository) GetByID(ctx context.Context, id int64) (*models.Course, error) {
	var course models.Course
	var createdAt, updatedAt string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, professor, created_at, updated_at
		FROM courses
		WHERE id = ?
	`, id).Scan(
		&course.ID,
		&course.Name,
		&course.Professor,
		&createdAt,
		&updatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	course.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	course.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &course, nil
}

// Create persists a new course in the database.
// It sets the ID field of the course object with the generated ID.
func (r *CourseRepository) Create(ctx context.Context, course *models.Course) error {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO courses (name, professor, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`,
		course.Name,
		course.Professor,
		course.CreatedAt.Format(time.RFC3339),
		course.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	course.ID = id
	return nil
}

// Update modifies an existing course in the database.
// All fields except CreatedAt can be updated.
func (r *CourseRepository) Update(ctx context.Context, course *models.Course) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE courses
		SET name = ?, professor = ?, updated_at = ?
		WHERE id = ?
	`,
		course.Name,
		course.Professor,
		course.UpdatedAt.Format(time.RFC3339),
		course.ID,
	)

	return err
}

// Delete removes a course from the database by its ID.
func (r *CourseRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM courses WHERE id = ?", id)
	return err
}
