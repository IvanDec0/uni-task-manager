// Package models contains the core domain entities for the University Task Manager application.
package models

import "time"

// Task represents a university task entity that tracks assignments, homework, or other academic work.
// Each task can be associated with a course and includes metadata like priority and status.
type Task struct {
	// ID uniquely identifies the task
	ID int64

	// Title is the name or short description of the task
	Title string

	// Description provides detailed information about what needs to be done
	Description string

	// DueDate specifies when the task needs to be completed
	DueDate time.Time

	// Priority ranges from 1 (lowest) to 5 (highest) to indicate task importance
	Priority int

	// Status indicates the current state of the task (pending, in progress, completed)
	Status TaskStatus

	// CourseID references the associated course (optional)
	CourseID int64

	// CreatedAt tracks when the task was created
	CreatedAt time.Time

	// UpdatedAt tracks the last modification time
	UpdatedAt time.Time
}

// TaskStatus represents the current status of a task as an enumerated type
type TaskStatus string

// Task status constants define the possible states of a task
const (
	// TaskStatusPending indicates a task that hasn't been started
	TaskStatusPending TaskStatus = "pending"

	// TaskStatusInProgress indicates a task that is currently being worked on
	TaskStatusInProgress TaskStatus = "in_progress"

	// TaskStatusCompleted indicates a task that has been finished
	TaskStatusCompleted TaskStatus = "completed"
)

// Course represents a university course that can have multiple associated tasks.
// It tracks basic course information including the professor teaching it.
type Course struct {
	// ID uniquely identifies the course
	ID int64

	// Name is the title of the course (e.g., "Computer Science 101")
	Name string

	// Professor is the name of the instructor teaching the course
	Professor string

	// CreatedAt tracks when the course was added to the system
	CreatedAt time.Time

	// UpdatedAt tracks the last modification time
	UpdatedAt time.Time
}
