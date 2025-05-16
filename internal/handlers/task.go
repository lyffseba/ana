package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sebae/ana/internal/models"
)

// Store tasks in memory for now (will be replaced with database later)
var tasks = models.MockTasks
var lastID = 2 // Since we start with 2 mock tasks

// GetTasks returns all tasks
func GetTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID returns a specific task by ID
func GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			c.JSON(http.StatusOK, task)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

// CreateTask creates a new task
func CreateTask(c *gin.Context) {
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set new task ID and timestamps
	lastID++
	newTask.ID = lastID
	newTask.CreatedAt = time.Now()
	newTask.UpdatedAt = time.Now()

	// Add to tasks slice
	tasks = append(tasks, newTask)
	c.JSON(http.StatusCreated, newTask)
}

// UpdateTask updates an existing task
func UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			// Keep the original ID and created date
			updatedTask.ID = id
			updatedTask.CreatedAt = task.CreatedAt
			updatedTask.UpdatedAt = time.Now()

			// Update the task in the slice
			tasks[i] = updatedTask
			c.JSON(http.StatusOK, updatedTask)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

// DeleteTask removes a task
func DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			// Remove the task from the slice
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

// GetTasksDueToday returns all tasks due today
func GetTasksDueToday(c *gin.Context) {
	todaysTasks := models.GetTasksDueToday()
	c.JSON(http.StatusOK, todaysTasks)
}

