// Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
// Last Updated: Sat May 17 07:34:44 AM CEST 2025

// Package repositories provides data access layer implementations
package repositories

import (
	"context"
	"time"

	"github.com/lyffseba/ana/internal/database"
	"github.com/lyffseba/ana/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskRepository handles database operations for tasks
// (MongoDB implementation)
type TaskRepository struct{}

// NewTaskRepository creates a new instance of TaskRepository
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

// FindAll retrieves all tasks from MongoDB
func (r *TaskRepository) FindAll() ([]models.Task, error) {
	coll := database.GetCollection("", "tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var tasks []models.Task
	for cur.Next(ctx) {
		var t models.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, cur.Err()
}

// FindByID retrieves a task by its ObjectID
func (r *TaskRepository) FindByID(id primitive.ObjectID) (models.Task, error) {
	coll := database.GetCollection("", "tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var task models.Task
	err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	return task, err
}

// Create adds a new task to MongoDB
func (r *TaskRepository) Create(task *models.Task) error {
	coll := database.GetCollection("", "tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if task.ID.IsZero() {
		task.ID = primitive.NewObjectID()
	}
	_, err := coll.InsertOne(ctx, task)
	return err
}

// Update modifies an existing task in MongoDB
func (r *TaskRepository) Update(task *models.Task) error {
	coll := database.GetCollection("", "tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.UpdateOne(ctx, bson.M{"_id": task.ID}, bson.M{"$set": task})
	return err
}

// Delete removes a task from MongoDB by ObjectID
func (r *TaskRepository) Delete(id primitive.ObjectID) error {
	coll := database.GetCollection("", "tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// FindTasksDueToday retrieves all tasks due on the current day
func (r *TaskRepository) FindTasksDueToday() ([]models.Task, error) {
	coll := database.GetCollection("", "tasks")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	end := start.AddDate(0, 0, 1)
	filter := bson.M{"due_date": bson.M{"$gte": start, "$lt": end}}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var tasks []models.Task
	for cur.Next(ctx) {
		var t models.Task
		if err := cur.Decode(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, cur.Err()
}

