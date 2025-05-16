// Package repositories provides data access layer implementations
package repositories

import (
	"log"
	"time"

	"github.com/sebae/ana/internal/database"
	"github.com/sebae/ana/internal/models"
)

// TaskRepository handles database operations for tasks
type TaskRepository struct{}

// NewTaskRepository creates a new instance of TaskRepository
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// FindAll retrieves all tasks from the database
func (r *TaskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	result := database.DB.Find(&tasks)
	return tasks, result.Error
}

// FindByID retrieves a task by its ID
func (r *TaskRepository) FindByID(id uint) (models.Task, error) {
	var task models.Task
	result := database.DB.First(&task, id)
	return task, result.Error
}

// Create adds a new task to the database
func (r *TaskRepository) Create(task *models.Task) error {
	return database.DB.Create(task).Error
}

// Update modifies an existing task in the database
func (r *TaskRepository) Update(task *models.Task) error {
	return database.DB.Save(task).Error
}

// Delete removes a task from the database
func (r *TaskRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Task{}, id).Error
}

// FindTasksDueToday retrieves all tasks due on the current day
func (r *TaskRepository) FindTasksDueToday() ([]models.Task, error) {
	var tasks []models.Task
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	end := start.AddDate(0, 0, 1)
	
	result := database.DB.Where("due_date >= ? AND due_date < ?", start, end).Find(&tasks)
	return tasks, result.Error
}

