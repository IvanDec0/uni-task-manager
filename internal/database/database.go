package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/IvanDec0/uni-task-manager/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// New creates a new database connection
func New(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Check the connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

// Initialize sets up the database tables
func (db *DB) Initialize() error {
	// Create courses table
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS courses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		professor TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("error creating courses table: %w", err)
	}

	// Create tasks table with foreign key to courses
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		due_date DATETIME NOT NULL,
		priority INTEGER NOT NULL,
		status TEXT NOT NULL,
		course_id INTEGER,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (course_id) REFERENCES courses(id)
	)`)
	if err != nil {
		return fmt.Errorf("error creating tasks table: %w", err)
	}

	log.Println("Database tables initialized successfully")
	return nil
}

// GetAllTasks returns all tasks from the database
func (db *DB) GetAllTasks() ([]models.Task, error) {
	rows, err := db.Query(`
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

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&dueDate,
			&task.Priority,
			&task.Status,
			&task.CourseID,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		// Parse time strings to time.Time
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

// GetTask retrieves a single task by ID
func (db *DB) GetTask(id int64) (models.Task, error) {
	var task models.Task
	var dueDate, createdAt, updatedAt string

	err := db.QueryRow(`
		SELECT id, title, description, due_date, priority, status, course_id, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&dueDate,
		&task.Priority,
		&task.Status,
		&task.CourseID,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return models.Task{}, err
	}

	task.DueDate, _ = time.Parse(time.RFC3339, dueDate)
	task.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	task.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return task, nil
}

// CreateTask adds a new task to the database
func (db *DB) CreateTask(task models.Task) (int64, error) {
	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now

	result, err := db.Exec(`
		INSERT INTO tasks (title, description, due_date, priority, status, course_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		task.Title,
		task.Description,
		task.DueDate.Format(time.RFC3339),
		task.Priority,
		task.Status,
		task.CourseID,
		task.CreatedAt.Format(time.RFC3339),
		task.UpdatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateTask updates an existing task
func (db *DB) UpdateTask(task models.Task) error {
	task.UpdatedAt = time.Now().UTC()

	_, err := db.Exec(`
		UPDATE tasks
		SET title = ?, description = ?, due_date = ?, priority = ?, status = ?, course_id = ?, updated_at = ?
		WHERE id = ?
	`,
		task.Title,
		task.Description,
		task.DueDate.Format(time.RFC3339),
		task.Priority,
		task.Status,
		task.CourseID,
		task.UpdatedAt.Format(time.RFC3339),
		task.ID,
	)

	return err
}

// DeleteTask removes a task from the database
func (db *DB) DeleteTask(id int64) error {
	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

// GetAllCourses returns all courses from the database
func (db *DB) GetAllCourses() ([]models.Course, error) {
	rows, err := db.Query(`
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

// CreateCourse adds a new course to the database
func (db *DB) CreateCourse(course models.Course) (int64, error) {
	now := time.Now().UTC()
	course.CreatedAt = now
	course.UpdatedAt = now

	result, err := db.Exec(`
		INSERT INTO courses (name, professor, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`,
		course.Name,
		course.Professor,
		course.CreatedAt.Format(time.RFC3339),
		course.UpdatedAt.Format(time.RFC3339),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
