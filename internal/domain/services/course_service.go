// Package services implements the core business logic for courses
package services

import (
	"context"
	"errors"
	"time"

	"uni-task-manager/internal/ports/output"

	"uni-task-manager/internal/domain/models"
	"uni-task-manager/internal/ports/input"
)

// Domain-specific errors that can be returned by the CourseService
var (
	// ErrCourseNotFound indicates that the requested course does not exist
	ErrCourseNotFound = errors.New("course not found")

	// ErrEmptyName indicates that the course name is empty, which is not allowed
	ErrEmptyName = errors.New("course name cannot be empty")
)

// Verify CourseService implements input.CourseService interface at compile time
var _ input.CourseService = (*CourseService)(nil)

// CourseService implements the course-related business logic and orchestrates
// interactions between the domain model and storage layer.
type CourseService struct {
	courseRepo output.CourseRepository
}

// NewCourseService creates a new instance of CourseService with the required dependencies.
func NewCourseService(courseRepo output.CourseRepository) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
	}
}

// CreateCourse implements input.CourseService.CreateCourse.
// It validates the course data before creation and sets metadata fields.
func (s *CourseService) CreateCourse(ctx context.Context, course *models.Course) error {
	if err := s.validateCourse(course); err != nil {
		return err
	}

	now := time.Now().UTC()
	course.CreatedAt = now
	course.UpdatedAt = now

	return s.courseRepo.Create(ctx, course)
}

// UpdateCourse implements input.CourseService.UpdateCourse.
// It validates the updated course data and ensures the course exists before updating.
func (s *CourseService) UpdateCourse(ctx context.Context, course *models.Course) error {
	if err := s.validateCourse(course); err != nil {
		return err
	}

	existing, err := s.courseRepo.GetByID(ctx, course.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCourseNotFound
	}

	course.CreatedAt = existing.CreatedAt
	course.UpdatedAt = time.Now().UTC()

	return s.courseRepo.Update(ctx, course)
}

// GetCourse implements input.CourseService.GetCourse.
// It retrieves a specific course by its ID.
func (s *CourseService) GetCourse(ctx context.Context, id int64) (*models.Course, error) {
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, ErrCourseNotFound
	}
	return course, nil
}

// GetAllCourses implements input.CourseService.GetAllCourses.
// It retrieves all courses from the repository.
func (s *CourseService) GetAllCourses(ctx context.Context) ([]models.Course, error) {
	return s.courseRepo.GetAll(ctx)
}

// DeleteCourse implements input.CourseService.DeleteCourse.
// It ensures the course exists before deletion.
func (s *CourseService) DeleteCourse(ctx context.Context, id int64) error {
	existing, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrCourseNotFound
	}
	return s.courseRepo.Delete(ctx, id)
}

// validateCourse performs validation of course data according to business rules.
// Currently, it only checks that the course name is not empty.
func (s *CourseService) validateCourse(course *models.Course) error {
	if course.Name == "" {
		return ErrEmptyName
	}
	return nil
}
