// Package main is the entry point for the University Task Manager application.
// It initializes the application components and starts the HTTP server.
package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	// Primary adapters (driving adapters)
	httpHandlers "uni-task-manager/internal/adapters/primary/http"
	// Secondary adapters (driven adapters)
	"uni-task-manager/internal/adapters/secondary/sqlite"
	// Domain services
	"uni-task-manager/internal/domain/services"

	"github.com/gorilla/mux"
	_ "modernc.org/sqlite"
)

// main is the application entry point that initializes and starts the server
func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

// run encapsulates the core application startup logic and error handling
func run() error {
	// Configure logging with file and line number for better debugging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting University Task Manager app...")

	// Initialize infrastructure components (directories, etc.)
	if err := setupInfrastructure(); err != nil {
		return err
	}

	// Initialize database connection and schema
	db, err := initializeDatabase()
	if err != nil {
		return err
	}
	defer db.Close()

	// Initialize application components (services, handlers)
	app, err := initializeApplication(db)
	if err != nil {
		return err
	}

	// Start HTTP server
	return startServer(app)
}

// setupInfrastructure ensures required directories exist
func setupInfrastructure() error {
	dbDir := filepath.Join(".", "data")
	return os.MkdirAll(dbDir, 0755)
}

// initializeDatabase sets up the SQLite database connection and schema
func initializeDatabase() (*sql.DB, error) {
	dbPath := filepath.Join(".", "data", "uni-tasks.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Initialize schema
	if err := createDatabaseSchema(db); err != nil {
		return nil, err
	}

	return db, nil
}

// application holds the initialized components of the application
type application struct {
	handler   *httpHandlers.Handler
	templates *template.Template
}

// initializeApplication sets up all application components following hexagonal architecture
func initializeApplication(db *sql.DB) (*application, error) {
	// Initialize repositories (secondary/driven adapters)
	taskRepo := sqlite.NewTaskRepository(db)
	courseRepo := sqlite.NewCourseRepository(db)

	// Initialize domain services
	taskService := services.NewTaskService(taskRepo, courseRepo)
	courseService := services.NewCourseService(courseRepo)

	// Load HTML templates
	templatesDir := filepath.Join(".", "web", "templates")
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		return nil, err
	}

	// Initialize HTTP handlers (primary/driving adapters)
	handler := httpHandlers.NewHandler(taskService, courseService, templates)

	return &application{
		handler:   handler,
		templates: templates,
	}, nil
}

// startServer configures and starts the HTTP server
func startServer(app *application) error {
	// Create router and configure routes
	r := mux.NewRouter()

	// Web routes for the user interface
	r.HandleFunc("/", app.handler.Index).Methods("GET")
	r.HandleFunc("/tasks/new", app.handler.CreateTaskForm).Methods("GET")
	r.HandleFunc("/tasks", app.handler.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id:[0-9]+}/edit", app.handler.EditTaskForm).Methods("GET")
	r.HandleFunc("/tasks/{id:[0-9]+}", app.handler.UpdateTask).Methods("POST")
	r.HandleFunc("/tasks/{id:[0-9]+}/delete", app.handler.DeleteTask).Methods("POST")
	r.HandleFunc("/courses", app.handler.ListCourses).Methods("GET")
	r.HandleFunc("/courses/new", app.handler.CreateCourseForm).Methods("GET")
	r.HandleFunc("/courses", app.handler.CreateCourse).Methods("POST")

	// API routes for programmatic access
	r.HandleFunc("/api/tasks", app.handler.APIGetTasks).Methods("GET")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", app.handler.APIGetTask).Methods("GET")
	r.HandleFunc("/api/tasks", app.handler.APICreateTask).Methods("POST")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", app.handler.APIUpdateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id:[0-9]+}", app.handler.APIDeleteTask).Methods("DELETE")

	// Configure server with timeouts for security and reliability
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Server started on http://localhost:8080")
	return srv.ListenAndServe()
}

// createDatabaseSchema initializes the SQLite database schema
func createDatabaseSchema(db *sql.DB) error {
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
		return err
	}

	// Create tasks table with foreign key to courses
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		due_date DATETIME NOT NULL,
		priority INTEGER NOT NULL CHECK (priority BETWEEN 1 AND 5),
		status TEXT NOT NULL CHECK (status IN ('pending', 'in_progress', 'completed')),
		course_id INTEGER,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (course_id) REFERENCES courses(id)
	)`)
	if err != nil {
		return err
	}

	log.Println("Database tables initialized successfully")
	return nil
}
