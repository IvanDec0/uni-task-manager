// Package http provides the primary adapters (driving adapters) for HTTP-based interaction.
// It implements both web interface handlers and REST API endpoints.
package http

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"uni-task-manager/internal/domain/models"
	"uni-task-manager/internal/ports/input"

	"github.com/gorilla/mux"
)

// Handler encapsulates the dependencies required for HTTP request handling.
// It serves as a primary adapter in the hexagonal architecture.
type Handler struct {
	taskService   input.TaskService
	courseService input.CourseService
	templates     *template.Template
}

// NewHandler creates a new instance of Handler with the required dependencies.
func NewHandler(taskService input.TaskService, courseService input.CourseService, templates *template.Template) *Handler {
	return &Handler{
		taskService:   taskService,
		courseService: courseService,
		templates:     templates,
	}
}

// Web Interface Handlers

// Index handles the home page request, displaying all tasks and courses.
func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tasks, err := h.taskService.GetAllTasks(ctx)
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	courses, err := h.courseService.GetAllCourses(ctx)
	if err != nil {
		http.Error(w, "Error fetching courses", http.StatusInternalServerError)
		return
	}

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

	h.templates.ExecuteTemplate(w, "index.html", data)
}

// CreateTaskForm displays the form for creating a new task.
func (h *Handler) CreateTaskForm(w http.ResponseWriter, r *http.Request) {
	courses, err := h.courseService.GetAllCourses(r.Context())
	if err != nil {
		http.Error(w, "Error fetching courses", http.StatusInternalServerError)
		return
	}

	h.templates.ExecuteTemplate(w, "create-task.html", courses)
}

// CreateTask handles the submission of a new task from the web form.
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	courseID, _ := strconv.ParseInt(r.FormValue("course_id"), 10, 64)
	priority, _ := strconv.Atoi(r.FormValue("priority"))
	dueDate, err := time.Parse("2006-01-02T15:04", r.FormValue("due_date"))
	if err != nil {
		http.Error(w, "Invalid due date format", http.StatusBadRequest)
		return
	}

	task := &models.Task{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		DueDate:     dueDate,
		Priority:    priority,
		Status:      models.TaskStatusPending,
		CourseID:    courseID,
	}

	err = h.taskService.CreateTask(r.Context(), task)
	if err != nil {
		http.Error(w, "Error creating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// EditTaskForm displays the form for editing an existing task.
func (h *Handler) EditTaskForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	task, err := h.taskService.GetTask(ctx, id)
	if err != nil {
		http.Error(w, "Error fetching task", http.StatusInternalServerError)
		return
	}

	courses, err := h.courseService.GetAllCourses(ctx)
	if err != nil {
		http.Error(w, "Error fetching courses", http.StatusInternalServerError)
		return
	}

	data := struct {
		Task    *models.Task
		Courses []models.Course
	}{
		Task:    task,
		Courses: courses,
	}

	h.templates.ExecuteTemplate(w, "edit-task.html", data)
}

// UpdateTask handles the submission of task updates from the web form.
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

	courseID, _ := strconv.ParseInt(r.FormValue("course_id"), 10, 64)
	priority, _ := strconv.Atoi(r.FormValue("priority"))
	dueDate, err := time.Parse("2006-01-02T15:04", r.FormValue("due_date"))
	if err != nil {
		http.Error(w, "Invalid due date format", http.StatusBadRequest)
		return
	}

	task := &models.Task{
		ID:          id,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		DueDate:     dueDate,
		Priority:    priority,
		Status:      models.TaskStatus(r.FormValue("status")),
		CourseID:    courseID,
	}

	err = h.taskService.UpdateTask(r.Context(), task)
	if err != nil {
		http.Error(w, "Error updating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// DeleteTask handles the deletion of a task from the web interface.
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.taskService.DeleteTask(r.Context(), id)
	if err != nil {
		http.Error(w, "Error deleting task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Course Management Handlers

// CreateCourseForm displays the form for creating a new course.
func (h *Handler) CreateCourseForm(w http.ResponseWriter, r *http.Request) {
	h.templates.ExecuteTemplate(w, "create-course.html", nil)
}

// CreateCourse handles the submission of a new course from the web form.
func (h *Handler) CreateCourse(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	course := &models.Course{
		Name:      r.FormValue("name"),
		Professor: r.FormValue("professor"),
	}

	err := h.courseService.CreateCourse(r.Context(), course)
	if err != nil {
		http.Error(w, "Error creating course: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/courses", http.StatusSeeOther)
}

// ListCourses displays the list of all courses.
func (h *Handler) ListCourses(w http.ResponseWriter, r *http.Request) {
	courses, err := h.courseService.GetAllCourses(r.Context())
	if err != nil {
		http.Error(w, "Error fetching courses", http.StatusInternalServerError)
		return
	}

	h.templates.ExecuteTemplate(w, "courses.html", courses)
}

// REST API Handlers

// APIGetTasks handles GET requests to retrieve all tasks.
// Returns a JSON array of tasks.
func (h *Handler) APIGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.taskService.GetAllTasks(r.Context())
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// APIGetTask handles GET requests to retrieve a specific task.
// Returns a JSON object containing task details or 404 if not found.
func (h *Handler) APIGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	task, err := h.taskService.GetTask(r.Context(), id)
	if err != nil {
		http.Error(w, "Error fetching task", http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// APICreateTask handles POST requests to create a new task.
// Accepts a JSON task object in the request body.
func (h *Handler) APICreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.taskService.CreateTask(r.Context(), &task)
	if err != nil {
		http.Error(w, "Error creating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// APIUpdateTask handles PUT requests to update an existing task.
// Accepts a JSON task object in the request body.
func (h *Handler) APIUpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task.ID = id
	err = h.taskService.UpdateTask(r.Context(), &task)
	if err != nil {
		http.Error(w, "Error updating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// APIDeleteTask handles DELETE requests to remove a task.
// Returns 204 No Content on success or appropriate error status.
func (h *Handler) APIDeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	err = h.taskService.DeleteTask(r.Context(), id)
	if err != nil {
		http.Error(w, "Error deleting task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
