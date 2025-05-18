// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lyffseba/ana/internal/models"
	"github.com/lyffseba/ana/internal/repositories"
)

// taskRepo is the repository for task operations
var taskRepo = repositories.NewTaskRepository()

// GetTasks returns all tasks
func GetTasks(c *gin.Context) {
	tasks, err := taskRepo.FindAll()
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}
	
	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID returns a specific task by ID
func GetTaskByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	task, err := taskRepo.FindByID(uint(id))
	if err != nil {
		log.Printf("Error fetching task with ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save to database using repository
	if err := taskRepo.Create(&newTask); err != nil {
		log.Printf("Error creating task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}

// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// First find the existing task
	existingTask, err := taskRepo.FindByID(uint(id))
	if err != nil {
		log.Printf("Error finding task to update with ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Bind JSON to the existing task
	if err := c.ShouldBindJSON(&existingTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID remains the same
	existingTask.ID = uint(id)

	// Update in the database
	if err := taskRepo.Update(&existingTask); err != nil {
		log.Printf("Error updating task with ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, existingTask)
}

// DeleteTask removes a task
func DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Check if task exists
	_, err = taskRepo.FindByID(uint(id))
	if err != nil {
		log.Printf("Error finding task to delete with ID %d: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Delete from database
	if err := taskRepo.Delete(uint(id)); err != nil {
		log.Printf("Error deleting task with ID %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// GetTasksDueToday returns all tasks due today
func GetTasksDueToday(c *gin.Context) {
	todaysTasks, err := taskRepo.FindTasksDueToday()
	if err != nil {
		log.Printf("Error fetching today's tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve today's tasks"})
		return
	}
	
	c.JSON(http.StatusOK, todaysTasks)
}

