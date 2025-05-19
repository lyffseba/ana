// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task represents a task in the system (MongoDB)
type Task struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Description string             `bson:"description" json:"description"`
	DueDate     time.Time          `bson:"due_date" json:"due_date"`
	Priority    string             `bson:"priority" json:"priority" binding:"oneof=Low Medium High"`
	ProjectID   int                `bson:"project_id" json:"project_id"`
	Status      string             `bson:"status" json:"status" binding:"oneof=To-Do In-Progress Done"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}


// MockTasks provides some sample tasks for development
var MockTasks = []Task{
	{
		ID:          primitive.NewObjectID(),
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
		ID:          primitive.NewObjectID(),
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

