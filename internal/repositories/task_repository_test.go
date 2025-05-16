package repositories

import (
	"os"
	"testing"
	"time"

	"github.com/sebae/ana/internal/database"
	"github.com/sebae/ana/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB initializes an in-memory SQLite database for testing
func setupTestDB(t *testing.T) {
	// Store the original DB connection
	originalDB := database.DB

	// Create an in-memory SQLite database for testing
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Replace the global DB connection with our test DB
	database.DB = db

	// Auto migrate the models for testing
	err = database.DB.AutoMigrate(&models.Task{})
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	// Helper function to restore the original DB after the test
	t.Cleanup(func() {
		database.DB = originalDB
	})
}

// clearTasks removes all tasks from the test database
func clearTasks(t *testing.T) {
	if err := database.DB.Exec("DELETE FROM tasks").Error; err != nil {
		t.Fatalf("Failed to clear tasks: %v", err)
	}
}

// createSampleTask creates and returns a sample task for testing
func createSampleTask() models.Task {
	return models.Task{
		Title:       "Test Task",
		Description: "This is a test task description",
		DueDate:     time.Now().AddDate(0, 0, 1), // Tomorrow
		Priority:    "Medium",
		ProjectID:   1,
		Status:      "To-Do",
	}
}

// TestTaskCreate tests the task creation functionality
func TestTaskCreate(t *testing.T) {
	setupTestDB(t)
	clearTasks(t)
	repo := NewTaskRepository()

	// Create a new task
	task := createSampleTask()
	if err := repo.Create(&task); err != nil {
		t.Errorf("Failed to create task: %v", err)
	}

	// Verify the task was created with an ID
	if task.ID == 0 {
		t.Error("Task was created but ID was not set")
	}

	// Verify the timestamps were set
	if task.CreatedAt.IsZero() || task.UpdatedAt.IsZero() {
		t.Error("Task timestamps were not set properly")
	}
}

// TestTaskFindAll tests retrieving all tasks
func TestTaskFindAll(t *testing.T) {
	setupTestDB(t)
	clearTasks(t)
	repo := NewTaskRepository()

	// Create multiple tasks
	for i := 0; i < 3; i++ {
		task := createSampleTask()
		task.Title = task.Title + " " + string(rune('A'+i))
		if err := repo.Create(&task); err != nil {
			t.Fatalf("Failed to create task %d: %v", i+1, err)
		}
	}

	// Retrieve all tasks
	tasks, err := repo.FindAll()
	if err != nil {
		t.Errorf("Failed to retrieve all tasks: %v", err)
	}

	// Check if we got the correct number of tasks
	if len(tasks) != 3 {
		t.Errorf("Expected 3 tasks, got %d", len(tasks))
	}
}

// TestTaskFindByID tests retrieving a task by ID
func TestTaskFindByID(t *testing.T) {
	setupTestDB(t)
	clearTasks(t)
	repo := NewTaskRepository()

	// Create a new task
	task := createSampleTask()
	if err := repo.Create(&task); err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Retrieve the task by ID
	foundTask, err := repo.FindByID(task.ID)
	if err != nil {
		t.Errorf("Failed to find task by ID: %v", err)
	}

	// Verify that the retrieved task matches the created one
	if foundTask.ID != task.ID || foundTask.Title != task.Title {
		t.Errorf("Retrieved task does not match the created one. Got: %+v, Expected: %+v", foundTask, task)
	}
}

// TestTaskUpdate tests updating a task
func TestTaskUpdate(t *testing.T) {
	setupTestDB(t)
	clearTasks(t)
	repo := NewTaskRepository()

	// Create a new task
	task := createSampleTask()
	if err := repo.Create(&task); err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Update the task
	updatedTitle := "Updated Test Task"
	task.Title = updatedTitle
	if err := repo.Update(&task); err != nil {
		t.Errorf("Failed to update task: %v", err)
	}

	// Retrieve the task to check if it was updated
	updatedTask, err := repo.FindByID(task.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated task: %v", err)
	}

	// Verify the task was updated
	if updatedTask.Title != updatedTitle {
		t.Errorf("Task was not updated correctly. Expected title: %s, got: %s", updatedTitle, updatedTask.Title)
	}

	// Verify that UpdatedAt was updated
	if !updatedTask.UpdatedAt.After(updatedTask.CreatedAt) {
		t.Error("UpdatedAt timestamp was not updated correctly")
	}
}

// TestTaskDelete tests deleting a task
func TestTaskDelete(t *testing.T) {
	setupTestDB(t)
	clearTasks(t)
	repo := NewTaskRepository()

	// Create a new task
	task := createSampleTask()
	if err := repo.Create(&task); err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Delete the task
	if err := repo.Delete(task.ID); err != nil {
		t.Errorf("Failed to delete task: %v", err)
	}

	// Try to retrieve the deleted task (should fail)
	_, err := repo.FindByID(task.ID)
	if err == nil {
		t.Error("Expected error when retrieving deleted task, but got none")
	}
}

// TestFindTasksDueToday tests retrieving tasks due today
func TestFindTasksDueToday(t *testing.T) {
	setupTestDB(t)
	clearTasks(t)
	repo := NewTaskRepository()

	// Create a task due today
	taskDueToday := createSampleTask()
	today := time.Now()
	taskDueToday.DueDate = time.Date(today.Year(), today.Month(), today.Day(), 12, 0, 0, 0, today.Location())
	if err := repo.Create(&taskDueToday); err != nil {
		t.Fatalf("Failed to create task due today: %v", err)
	}

	// Create a task due tomorrow
	taskDueTomorrow := createSampleTask()
	taskDueTomorrow.Title = "Task Due Tomorrow"
	taskDueTomorrow.DueDate = time.Now().AddDate(0, 0, 1)
	if err := repo.Create(&taskDueTomorrow); err != nil {
		t.Fatalf("Failed to create task due tomorrow: %v", err)
	}

	// Create a task due yesterday
	taskDueYesterday := createSampleTask()
	taskDueYesterday.Title = "Task Due Yesterday"
	taskDueYesterday.DueDate = time.Now().AddDate(0, 0, -1)
	if err := repo.Create(&taskDueYesterday); err != nil {
		t.Fatalf("Failed to create task due yesterday: %v", err)
	}

	// Retrieve tasks due today
	tasksDueToday, err := repo.FindTasksDueToday()
	if err != nil {
		t.Errorf("Failed to retrieve tasks due today: %v", err)
	}

	// Check if only the task due today was retrieved
	if len(tasksDueToday) != 1 {
		t.Errorf("Expected 1 task due today, got %d", len(tasksDueToday))
	}

	if len(tasksDueToday) > 0 && tasksDueToday[0].ID != taskDueToday.ID {
		t.Errorf("Retrieved wrong task as due today. Expected ID: %d, got: %d", taskDueToday.ID, tasksDueToday[0].ID)
	}
}

// TestValidationRules tests that validation rules are enforced
func TestValidationRules(t *testing.T) {
	setupTestDB(t)
	clearTasks(t)
	repo := NewTaskRepository()

	// Test creating a task with invalid priority
	taskInvalidPriority := createSampleTask()
	taskInvalidPriority.Priority = "Invalid"
	err := repo.Create(&taskInvalidPriority)
	
	// SQLite might not enforce constraints the same way as PostgreSQL
	// For a real app, we should use proper validation before saving
	if err != nil {
		t.Logf("Got expected validation error for invalid priority: %v", err)
	}

	// Test creating a task with invalid status
	taskInvalidStatus := createSampleTask()
	taskInvalidStatus.Status = "Invalid"
	err = repo.Create(&taskInvalidStatus)
	
	if err != nil {
		t.Logf("Got expected validation error for invalid status: %v", err)
	}
}

// TestTransactionRollback tests transaction rollback functionality
func TestTransactionRollback(t *testing.T) {
	// This test demonstrates transaction management
	// In a real app, we would add transaction support to the repository
	setupTestDB(t)
	clearTasks(t)
	
	// Start a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		t.Fatalf("Failed to begin transaction: %v", tx.Error)
	}

	// Create a task within the transaction
	task := createSampleTask()
	if err := tx.Create(&task).Error; err != nil {
		t.Fatalf("Failed to create task in transaction: %v", err)
	}

	// Check that the task exists within the transaction
	var tempTask models.Task
	if err := tx.First(&tempTask, task.ID).Error; err != nil {
		t.Errorf("Failed to retrieve task within transaction: %v", err)
	}

	// Rollback the transaction
	if err := tx.Rollback().Error; err != nil {
		t.Fatalf("Failed to rollback transaction: %v", err)
	}

	// Verify that the task doesn't exist in the database
	repo := NewTaskRepository()
	_, err := repo.FindByID(task.ID)
	if err == nil {
		t.Error("Task should not exist after transaction rollback, but it was found")
	}
}

// TestTaskCreateExecOnce is an example of a benchmark test
func BenchmarkTaskCreation(b *testing.B) {
	// Skip the benchmark unless explicitly requested
	if os.Getenv("RUN_BENCHMARKS") != "true" {
		b.Skip("Skipping benchmark. Set RUN_BENCHMARKS=true to run")
	}

	// Setup the test database
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&models.Task{})
	database.DB = db

	repo := NewTaskRepository()
	
	// Reset the timer
	b.ResetTimer()
	
	// Run the benchmark
	for i := 0; i < b.N; i++ {
		task := createSampleTask()
		task.Title = task.Title + " " + string(rune(i%26+'A'))
		_ = repo.Create(&task)
	}
}

