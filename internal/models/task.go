// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents a task in the system
type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null" binding:"required"`
	Description string         `json:"description" gorm:"type:text"`
	DueDate     time.Time      `json:"due_date" gorm:"index"`
	Priority    string         `json:"priority" gorm:"size:50" binding:"oneof=Low Medium High''"`
	ProjectID   int            `json:"project_id" gorm:"index"`
	Status      string         `json:"status" gorm:"size:50;index" binding:"oneof=To-Do In-Progress Done''"`
	CreatedAt   time.Time      `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // Supports soft delete
}

// TableName specifies the table name for the Task model
func (Task) TableName() string {
	return "tasks"
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

