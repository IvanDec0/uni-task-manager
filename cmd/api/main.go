package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/IvanDec0/uni-task-manager/internal/database"
	"github.com/IvanDec0/uni-task-manager/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting University Task Manager app...")

	// Create database directory if it doesn't exist
	dbDir := filepath.Join(".", "data")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory: %v", err)
	}

	// Initialize database connection
	dbPath := filepath.Join(dbDir, "uni-tasks.db")
	db, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Initialize database schema
	if err := db.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	// Parse templates
	templatesDir := filepath.Join(".", "web", "templates")
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	// Initialize handlers
	h := handlers.NewHandler(db, templates)

	// Create router
	r := mux.NewRouter()

	// Web routes
	r.HandleFunc("/", h.Index).Methods("GET")
	r.HandleFunc("/tasks/new", h.CreateTaskForm).Methods("GET")
	r.HandleFunc("/tasks", h.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id:[0-9]+}/edit", h.EditTaskForm).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", h.UpdateTask).Methods("POST")
	r.HandleFunc("/tasks/{id:[0-9]+}/delete", h.DeleteTask).Methods("POST")
	r.HandleFunc("/courses", h.Courses).Methods("GET")
	r.HandleFunc("/courses/new", h.CreateCourseForm).Methods("GET")
	r.HandleFunc("/courses", h.CreateCourse).Methods("POST")

	// API routes
	r.HandleFunc("/api/tasks", h.APIGetTasks).Methods("GET")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", h.APIGetTask).Methods("GET")
	r.HandleFunc("/api/tasks", h.APICreateTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", h.APIUpdateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", h.APIDeleteTask).Methods("DELETE")

	// Start server
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server started on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
