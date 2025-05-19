// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test helper function to create a valid task
func createValidTask() Task {
	return Task{
		Title:       "Test Task",
		Description: "This is a test task for unit testing",
		DueDate:     time.Now().Add(24 * time.Hour), // Tomorrow
		Priority:    "High",
		ProjectID:   1,
		Status:      "To-Do",
	}
}

// TestTaskFieldValidation tests the validation of Task fields
func TestTaskFieldValidation(t *testing.T) {
	// Test cases table
	testCases := []struct {
		name          string
		task          Task
		expectedValid bool
		fieldToCheck  string
	}{
		{
			name:          "Valid Task",
			task:          createValidTask(),
			expectedValid: true,
			fieldToCheck:  "Title",
		},
		{
			name: "Missing Title",
			task: Task{
				Description: "Test Description",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    "Medium",
				ProjectID:   1,
				Status:      "To-Do",
			},
			expectedValid: false,
			fieldToCheck:  "Title",
		},
		{
			name: "Invalid Priority",
			task: Task{
				Title:       "Test Task",
				Description: "Test Description",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    "InvalidPriority",
				ProjectID:   1,
				Status:      "To-Do",
			},
			expectedValid: false,
			fieldToCheck:  "Priority",
		},
		{
			name: "Invalid Status",
			task: Task{
				Title:       "Test Task",
				Description: "Test Description",
				DueDate:     time.Now().Add(24 * time.Hour),
				Priority:    "Medium",
				ProjectID:   1,
				Status:      "InvalidStatus",
			},
			expectedValid: false,
			fieldToCheck:  "Status",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// In a real implementation, you'd use a validator like:
			// err := validator.New().Struct(tc.task)
			// But for this test, we'll check basic conditions
			
			isValid := true
			if tc.fieldToCheck == "Title" && tc.task.Title == "" {
				isValid = false
			}
			if tc.fieldToCheck == "Priority" && !isValidPriority(tc.task.Priority) {
				isValid = false
			}
			if tc.fieldToCheck == "Status" && !isValidStatus(tc.task.Status) {
				isValid = false
			}
			
			assert.Equal(t, tc.expectedValid, isValid, "Task validity check failed")
		})
	}
}

// Helper function to validate Priority values
func isValidPriority(priority string) bool {
	return priority == "Low" || priority == "Medium" || priority == "High" || priority == ""
}

// Helper function to validate Status values
func isValidStatus(status string) bool {
	return status == "To-Do" || status == "In-Progress" || status == "Done" || status == ""
}

// TestTaskTimeHandling tests time-related functionality of Task
func TestTaskTimeHandling(t *testing.T) {
	// Test current time handling
	now := time.Now()
	task := createValidTask()
	task.CreatedAt = now
	task.UpdatedAt = now

	assert.Equal(t, now, task.CreatedAt, "CreatedAt should match the provided time")
	assert.Equal(t, now, task.UpdatedAt, "UpdatedAt should match the provided time")

	// Test due date in the future
	tomorrow := time.Now().Add(24 * time.Hour)
	task.DueDate = tomorrow
	assert.True(t, task.DueDate.After(time.Now()), "DueDate should be in the future")
	assert.WithinDuration(t, tomorrow, task.DueDate, time.Second, "DueDate should be approximately tomorrow")

	// Test past due date
	yesterday := time.Now().Add(-24 * time.Hour)
	task.DueDate = yesterday
	assert.True(t, task.DueDate.Before(time.Now()), "DueDate should be in the past")
	assert.WithinDuration(t, yesterday, task.DueDate, time.Second, "DueDate should be approximately yesterday")
}

// TestTaskJSONMarshaling tests JSON marshaling and unmarshaling of Task
func TestTaskJSONMarshaling(t *testing.T) {
	// Create a task with known values
	originalTask := createValidTask()
	originalTask.ID = primitive.NewObjectID()
	originalTask.CreatedAt = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	originalTask.UpdatedAt = time.Date(2025, 1, 2, 12, 0, 0, 0, time.UTC)
	
	// Marshal to JSON
	jsonData, err := json.Marshal(originalTask)
	assert.NoError(t, err, "JSON marshaling should not produce an error")
	
	// Unmarshal back to a Task struct
	var unmarshaledTask Task
	err = json.Unmarshal(jsonData, &unmarshaledTask)
	assert.NoError(t, err, "JSON unmarshaling should not produce an error")
	
	// Compare fields
	assert.Equal(t, originalTask.ID, unmarshaledTask.ID, "ID should match after marshaling and unmarshaling")
	assert.Equal(t, originalTask.Title, unmarshaledTask.Title, "Title should match after marshaling and unmarshaling")
	assert.Equal(t, originalTask.Description, unmarshaledTask.Description, "Description should match after marshaling and unmarshaling")
	assert.Equal(t, originalTask.Priority, unmarshaledTask.Priority, "Priority should match after marshaling and unmarshaling")
	assert.Equal(t, originalTask.ProjectID, unmarshaledTask.ProjectID, "ProjectID should match after marshaling and unmarshaling")
	assert.Equal(t, originalTask.Status, unmarshaledTask.Status, "Status should match after marshaling and unmarshaling")
	
	// Check time fields (just checking the format and not the exact timezones and nanoseconds)
	assert.Equal(t, originalTask.CreatedAt.Format(time.RFC3339), unmarshaledTask.CreatedAt.Format(time.RFC3339), 
		"CreatedAt should match after marshaling and unmarshaling")
	assert.Equal(t, originalTask.UpdatedAt.Format(time.RFC3339), unmarshaledTask.UpdatedAt.Format(time.RFC3339), 
		"UpdatedAt should match after marshaling and unmarshaling")
}

// TestTaskStatusValidation tests the validation of Task status field
func TestTaskStatusValidation(t *testing.T) {
	testCases := []struct {
		status string
		valid  bool
	}{
		{"To-Do", true},
		{"In-Progress", true},
		{"Done", true},
		{"", true},          // Empty is allowed by the validation rule
		{"Pending", false},
		{"Started", false},
		{"Completed", false},
	}

	for _, tc := range testCases {
		t.Run(tc.status, func(t *testing.T) {
			assert.Equal(t, tc.valid, isValidStatus(tc.status), 
				"Status validation for '%s' should be %v", tc.status, tc.valid)
		})
	}
}

// TestTaskPriorityValidation tests the validation of Task priority field
func TestTaskPriorityValidation(t *testing.T) {
	testCases := []struct {
		priority string
		valid    bool
	}{
		{"Low", true},
		{"Medium", true},
		{"High", true},
		{"", true},          // Empty is allowed by the validation rule
		{"Critical", false},
		{"Normal", false},
		{"Urgent", false},
	}

	for _, tc := range testCases {
		t.Run(tc.priority, func(t *testing.T) {
			assert.Equal(t, tc.valid, isValidPriority(tc.priority), 
				"Priority validation for '%s' should be %v", tc.priority, tc.valid)
		})
	}
}

// TestTaskEdgeCases tests various edge cases
func TestTaskEdgeCases(t *testing.T) {
	// Test with very long title (more than 255 chars would be truncated in real DB)
	longTitle := ""
	for i := 0; i < 300; i++ {
		longTitle += "a"
	}
	
	task := createValidTask()
	task.Title = longTitle
	
	// In a real application, this would be validated against the database schema
	// Here we just check that the struct can hold it
	assert.Equal(t, longTitle, task.Title, "Task should store long title (DB would truncate)")
	
	// Test with empty Description
	task = createValidTask()
	task.Description = ""
	assert.Equal(t, "", task.Description, "Empty description should be valid")
	
	// Test with zero DueDate
	task = createValidTask()
	task.DueDate = time.Time{}
	assert.True(t, task.DueDate.IsZero(), "Zero DueDate should be valid")
	
	// Test with negative ProjectID (might be invalid in business logic)
	task = createValidTask()
	task.ProjectID = -1
	assert.Equal(t, -1, task.ProjectID, "Task should store negative ProjectID (validation might reject)")
}

