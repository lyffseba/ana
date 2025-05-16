package models

import (
	"time"
)

// Task represents a task in the system
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
	Priority    string    `json:"priority" binding:"oneof=Low Medium High''"`
	ProjectID   int       `json:"project_id"`
	Status      string    `json:"status" binding:"oneof=To-Do In-Progress Done''"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// MockTasks provides some sample tasks for development
var MockTasks = []Task{
	{
		ID:          1,
		Title:       "Meet Client",
		Description: "Discuss project requirements for the new residential building",
		DueDate:     time.Now().AddDate(0, 0, 1), // Tomorrow
		Priority:    "High",
		ProjectID:   1,
		Status:      "To-Do",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          2,
		Title:       "Finalize Blueprints",
		Description: "Complete the final version of blueprint documents",
		DueDate:     time.Now().AddDate(0, 0, 3), // In 3 days
		Priority:    "Medium",
		ProjectID:   1,
		Status:      "To-Do",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

// GetTasksDueToday returns all tasks that are due today
func GetTasksDueToday() []Task {
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	end := start.AddDate(0, 0, 1)

	var todaysTasks []Task
	for _, task := range MockTasks {
		if task.DueDate.After(start) && task.DueDate.Before(end) {
			todaysTasks = append(todaysTasks, task)
		}
	}
	return todaysTasks
}

