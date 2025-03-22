package models

import (
	"time"
)

// Task represents a university task with deadlines and priority
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    int       `json:"priority"` // 1-5 scale: 1 (lowest) to 5 (highest)
	Status      string    `json:"status"`   // "pending", "in_progress", "completed"
	CourseID    int64     `json:"course_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Course represents a university course
type Course struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Professor string    `json:"professor"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
