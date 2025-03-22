package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/IvanDec0/uni-task-manager/internal/database"
	"github.com/IvanDec0/uni-task-manager/internal/models"
	"github.com/gorilla/mux"
)

// Handler contains dependencies for the handlers
type Handler struct {
	DB        *database.DB
	Templates *template.Template
}

// NewHandler creates a new handler with dependencies
func NewHandler(db *database.DB, templates *template.Template) *Handler {
	return &Handler{
		DB:        db,
		Templates: templates,
	}
}

// Index handles the home page
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.DB.GetAllTasks()
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	courses, err := h.DB.GetAllCourses()
	if err != nil {
		http.Error(w, "Error fetching courses", http.StatusInternalServerError)
		return
	}

	// Create a map from course ID to course name for displaying in the template
	courseMap := make(map[int64]string)
	for _, course := range courses {
		courseMap[course.ID] = course.Name
	}

	data := struct {
		Tasks     []models.Task
		Courses   []models.Course
		CourseMap map[int64]string
	}{
		Tasks:     tasks,
		Courses:   courses,
		CourseMap: courseMap,
	}

	h.Templates.ExecuteTemplate(w, "index.html", data)
}

// CreateTaskForm handles the task creation form
func (h *Handler) CreateTaskForm(w http.ResponseWriter, r *http.Request) {
	courses, err := h.DB.GetAllCourses()
	if err != nil {
		http.Error(w, "Error fetching courses", http.StatusInternalServerError)
		return
	}

	h.Templates.ExecuteTemplate(w, "create-task.html", courses)
}

// CreateTask handles the task creation POST request
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Parse course ID
	courseID, err := strconv.ParseInt(r.FormValue("course_id"), 10, 64)
	if err != nil {
		courseID = 0 // No course selected
	}

	// Parse priority
	priority, err := strconv.Atoi(r.FormValue("priority"))
	if err != nil || priority < 1 || priority > 5 {
		priority = 3 // Default to medium priority
	}

	// Parse due date
	dueDate, err := time.Parse("2006-01-02T15:04", r.FormValue("due_date"))
	if err != nil {
		http.Error(w, "Invalid due date format", http.StatusBadRequest)
		return
	}

	task := models.Task{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		DueDate:     dueDate,
		Priority:    priority,
		Status:      "pending",
		CourseID:    courseID,
	}

	_, err = h.DB.CreateTask(task)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// EditTaskForm handles the task edit form
func (h *Handler) EditTaskForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.DB.GetTask(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	courses, err := h.DB.GetAllCourses()
	if err != nil {
		http.Error(w, "Error fetching courses", http.StatusInternalServerError)
		return
	}

	data := struct {
		Task    models.Task
		Courses []models.Course
	}{
		Task:    task,
		Courses: courses,
	}

	h.Templates.ExecuteTemplate(w, "edit-task.html", data)
}

// UpdateTask handles the task update POST request
func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	// Get existing task to preserve creation time
	existingTask, err := h.DB.GetTask(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Parse course ID
	courseID, err := strconv.ParseInt(r.FormValue("course_id"), 10, 64)
	if err != nil {
		courseID = 0 // No course selected
	}

	// Parse priority
	priority, err := strconv.Atoi(r.FormValue("priority"))
	if err != nil || priority < 1 || priority > 5 {
		priority = 3 // Default to medium priority
	}

	// Parse due date
	dueDate, err := time.Parse("2006-01-02T15:04", r.FormValue("due_date"))
	if err != nil {
		http.Error(w, "Invalid due date format", http.StatusBadRequest)
		return
	}

	task := models.Task{
		ID:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		DueDate:     dueDate,
		Priority:    priority,
		Status:      r.FormValue("status"),
		CourseID:    courseID,
		CreatedAt:   existingTask.CreatedAt,
	}

	err = h.DB.UpdateTask(task)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteTask handles the task deletion
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.DB.DeleteTask(id)
	if err != nil {
		http.Error(w, "Error deleting task", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CreateCourseForm handles the course creation form
func (h *Handler) CreateCourseForm(w http.ResponseWriter, r *http.Request) {
	h.Templates.ExecuteTemplate(w, "create-course.html", nil)
}

// CreateCourse handles the course creation POST request
func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	course := models.Course{
		Name:      r.FormValue("name"),
		Professor: r.FormValue("professor"),
	}

	_, err := h.DB.CreateCourse(course)
	if err != nil {
		http.Error(w, "Error creating course", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/courses", http.StatusSeeOther)
}

// Courses handles the courses listing page
func (h *Handler) Courses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.DB.GetAllCourses()
	if err != nil {
		http.Error(w, "Error fetching courses", http.StatusInternalServerError)
		return
	}

	h.Templates.ExecuteTemplate(w, "courses.html", courses)
}

// API handlers for JSON responses

// APIGetTasks returns all tasks as JSON
func (h *Handler) APIGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.DB.GetAllTasks()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching tasks: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// APIGetTask returns a single task as JSON
func (h *Handler) APIGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.DB.GetTask(id)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// APICreateTask creates a task from JSON
func (h *Handler) APICreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.DB.CreateTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating task: %v", err), http.StatusInternalServerError)
		return
	}

	task.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// APIUpdateTask updates a task from JSON
func (h *Handler) APIUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task.ID = id
	err = h.DB.UpdateTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error updating task: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// APIDeleteTask deletes a task
func (h *Handler) APIDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.DB.DeleteTask(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting task: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
