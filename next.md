# Implementation Plan for ana.world

## Current Status Overview

We have completed the initial project setup as outlined in Phase 1 of our plan:

- ✅ Created GitHub repository for version control
- ✅ Set up Go development environment with Gin framework
- ✅ Implemented basic API endpoints for task management
- ✅ Created frontend static HTML files with HTMX integration
- ✅ Designed UI with terra colors using Tailwind CSS
- ✅ Implemented basic routing with proper separation of API and static files

## Next Implementation Steps

### 1. Google Calendar and Gmail Integration (Priority: Highest)

- [ ] **Set up Google Cloud Project**
  - Create a new Google Cloud project for ana.world
  - Enable Google Calendar API and Gmail API
  - Configure OAuth 2.0 consent screen for anais.villamarinj@gmail.com
  - Create API keys and client secrets
  - Refer to [google_calendar_gmail.md](google_calendar_gmail.md) for detailed setup instructions

- [ ] **Implement Gmail integration for anais.villamarinj@gmail.com**
  - Set up OAuth 2.0 flow with proper scopes
  - Create email notification system for project updates
  - Implement email sending for task reminders
  - Add support for email-based task creation
  - Design bilingual email templates (English/Spanish)

- [ ] **Implement Google Calendar integration**
  - Sync project timelines and milestones to Google Calendar
  - Create calendar events for task deadlines
  - Set up bidirectional sync between app tasks and calendar events
  - Configure recurring events for regular project check-ins
  - Use color coding for task priorities (Red=High, Yellow=Medium, Blue=Low)

- [ ] **Build UI for Google integrations**
  - Add OAuth authorization flow in the frontend
  - Create settings page for managing Google integrations
  - Implement calendar viewing interface embedded in the app
  - Add email preference configuration options

- [ ] **Ensure data security and privacy**
  - Implement token storage with proper encryption
  - Set up token refresh mechanism
  - Create permission scopes with least privilege principle
  - Add clear documentation on data access and usage

See [google_calendar_gmail.md](google_calendar_gmail.md) for detailed implementation plan and timeline.

### 2. PostgreSQL Database Integration (Priority: High)

- [ ] **Install and configure PostgreSQL**
  - Install PostgreSQL 14+ locally for development
  - Create dedicated database and user for the application
  - Set up connection pooling for optimal performance

- [ ] **Design database schema**
  - Create initial migration for tasks table with fields:
    ```sql
    CREATE TABLE tasks (
      id SERIAL PRIMARY KEY,
      title VARCHAR(255) NOT NULL,
      description TEXT,
      due_date TIMESTAMP,
      priority VARCHAR(50),
      project_id INTEGER,
      status VARCHAR(50),
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    ```
  - Set up proper indexes for frequently queried fields

- [ ] **Implement database connection**
  - Add PostgreSQL driver (`github.com/lib/pq`) and GORM ORM (`gorm.io/gorm`)
  - Create database package with connection management
  - Implement environment-based configuration
  - Add connection pooling with sensible defaults

- [ ] **Create data access layer**
  - Refactor models to use GORM tags
  - Create repository pattern implementation for tasks
  - Implement proper error handling and logging

### 2. Authentication System (Priority: Medium)

- [ ] **Research authentication options**
  - Evaluate Google OAuth vs. simple authentication
  - Consider token-based authentication with JWT

- [ ] **Implement authentication middleware**
  - Create middleware for auth token validation
  - Set up protected routes
  - Implement CSRF protection

- [ ] **Set up user model and management**
  - Create users table with appropriate fields
  - Implement registration and login flows
  - Add user-task relationships

### 3. Enhanced Task Management (Priority: Medium)

- [ ] **Complete task CRUD operations**
  - Update handlers to use database repositories
  - Implement filtering and sorting options
  - Add pagination for performance

- [ ] **Improve error handling**
  - Create standardized error responses
  - Add validation for inputs
  - Implement proper logging

- [ ] **Implement task completion functionality**
  - Add status update endpoints
  - Create frontend for marking tasks complete

### 4. AI Assistant with Cerebras API (Priority: Medium)

- [ ] **Set up Cerebras API integration**
  - Register for Cerebras API access (https://inference.cerebras.ai/)
  - Create API key and configure environment variables
  - Implement API client in Go with proper error handling

- [ ] **Design AI assistant interface**
  - Create conversation model structure with user and AI messages
  - Implement context-aware prompts with architectural domain knowledge
  - Add message history management for coherent conversations

- [ ] **Build assistant frontend**
  - Design chat interface with HTMX for real-time updates
  - Implement user input validation and response formatting
  - Create responsive UI with terra color theme
  - Add loading indicators for API requests

- [ ] **Implement domain-specific features**
  - Project context awareness (reference current projects in conversations)
  - Support for architectural terminology in both English and Spanish
  - Ability to suggest task creation from conversation
  - Task scheduling assistant capabilities

### 5. Deployment Preparation (Priority: Low)

- [ ] **Configure Netlify deployment**
  - Refine netlify.toml configuration
  - Set up build scripts
  - Configure environment variables

- [ ] **Prepare for backend deployment**
  - Research Netlify Functions vs. standalone backend
  - Document deployment process
  - Create CI/CD pipeline

## Best Practices and Security Considerations

### Database Security

1. **Connection Security**
   - Use environment variables for sensitive credentials
   - Never hardcode database credentials in application code
   - Use TLS/SSL for database connections in production

2. **Query Security**
   - Use parameterized queries to prevent SQL injection
   - Implement proper input validation
   - Use the principle of least privilege for database users

3. **Data Protection**
   - Consider encryption for sensitive data
   - Implement proper backup strategies
   - Use database migrations for schema changes

### Application Security

1. **Authentication Best Practices**
   - Implement proper password hashing (if using password auth)
   - Use secure tokens with reasonable expiration
   - Implement rate limiting for auth endpoints

2. **API Security**
   - Validate all inputs
   - Implement proper CORS configuration
   - Use HTTPS in all environments

3. **Dependency Management**
   - Regularly update dependencies
   - Use vulnerability scanning tools
   - Follow security advisories for key packages

## Implementation Approach for Key Features

### Implementation Approach for Google Calendar & Gmail Integration

#### Week 1: OAuth and Configuration Setup

```bash
# Install required Go packages
go get -u google.golang.org/api/calendar/v3
go get -u google.golang.org/api/gmail/v1
go get -u golang.org/x/oauth2/google
```

1. Create a basic OAuth flow for Google account:

```go
// internal/google/auth.go
package google

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GetConfig returns a config for OAuth2 authentication
func GetConfig(scopes []string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, fmt.Errorf("unable to read client credentials file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client credentials file: %v", err)
	}

	return config, nil
}

// GetTokenFromWeb gets a token from the web flow
func GetTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following URL and enter the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %v", err)
	}
	return tok, nil
}
```

2. Set up environment variables for the project:

```bash
# Add to .env file
GOOGLE_CLIENT_ID=your_client_id_here
GOOGLE_CLIENT_SECRET=your_client_secret_here
GOOGLE_REDIRECT_URI=http://localhost:8080/api/auth/google/callback
GOOGLE_PRIMARY_EMAIL=anais.villamarinj@gmail.com
```

3. Create handlers for OAuth authentication flow:

```go
// internal/handlers/google_auth_handler.go
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sebae/ana/internal/google"
)

// GoogleAuthHandler initiates the OAuth flow
func GoogleAuthHandler(c *gin.Context) {
	// Implement OAuth flow initiation
}

// GoogleCallbackHandler handles the OAuth callback
func GoogleCallbackHandler(c *gin.Context) {
	// Implement OAuth callback handling
}
```

See [google_calendar_gmail.md](google_calendar_gmail.md) for complete implementation details.

### Implementation Approach for PostgreSQL Setup

### Step 1: Install PostgreSQL

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start and enable the service
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

### Step 2: Configure Database

```bash
# Create database and user
sudo -u postgres psql

postgres=# CREATE DATABASE ana_world;
postgres=# CREATE USER ana_user WITH ENCRYPTED PASSWORD 'secure_password';
postgres=# GRANT ALL PRIVILEGES ON DATABASE ana_world TO ana_user;
postgres=# \q
```

### Step 3: Set Up Database Connection in Go

1. Create a new package for database operations:

```go
// internal/database/db.go
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "ana_user"),
		getEnv("DB_PASSWORD", "secure_password"),
		getEnv("DB_NAME", "ana_world"),
		getEnv("DB_PORT", "5432"),
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
```

2. Update models to use GORM:

```go
// internal/models/task.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null" binding:"required"`
	Description string         `json:"description" gorm:"type:text"`
	DueDate     time.Time      `json:"due_date"`
	Priority    string         `json:"priority" gorm:"size:50" binding:"oneof=Low Medium High''"`
	ProjectID   int            `json:"project_id"`
	Status      string         `json:"status" gorm:"size:50" binding:"oneof=To-Do In-Progress Done''"`
	CreatedAt   time.Time      `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
```

3. Create a task repository:

```go
// internal/repositories/task_repository.go
package repositories

import (
	"github.com/sebae/ana/internal/database"
	"github.com/sebae/ana/internal/models"
)

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) FindAll() ([]models.Task, error) {
	var tasks []models.Task
	result := database.DB.Find(&tasks)
	return tasks, result.Error
}

func (r *TaskRepository) FindByID(id uint) (models.Task, error) {
	var task models.Task
	result := database.DB.First(&task, id)
	return task, result.Error
}

func (r *TaskRepository) Create(task *models.Task) error {
	return database.DB.Create(task).Error
}

func (r *TaskRepository) Update(task *models.Task) error {
	return database.DB.Save(task).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Task{}, id).Error
}

func (r *TaskRepository) FindTasksDueToday() ([]models.Task, error) {
	var tasks []models.Task
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	end := start.AddDate(0, 0, 1)
	
	result := database.DB.Where("due_date >= ? AND due_date < ?", start, end).Find(&tasks)
	return tasks, result.Error
}
```

## References

### PostgreSQL with Go
- [GORM Documentation](https://gorm.io/docs/index.html) - The official GORM ORM documentation
- [PostgreSQL Documentation](https://www.postgresql.org/docs/) - Official PostgreSQL documentation
- [Go Database/SQL Tutorial](https://go.dev/doc/tutorial/database-access) - Go's official database access tutorial

### Authentication
- [OAuth 2.0 with Go](https://developers.google.com/identity/protocols/oauth2) - Google's OAuth 2.0 documentation
- [OWASP Authentication Best Practices](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html) - Security best practices
- [JWT Authentication in Go](https://www.sohamkamani.com/golang/jwt-authentication/) - Guide to implementing JWT authentication

### HTMX + Go Integration
- [HTMX Go Examples](https://htmx.org/examples/) - Official HTMX examples that can be adapted for Go
- [go-htmx Package](https://github.com/donseba/go-htmx) - Go package for HTMX integration

### Google APIs Integration
- [Google Calendar API Go Quickstart](https://developers.google.com/calendar/api/quickstart/go) - Official quickstart guide
- [Gmail API Go Quickstart](https://developers.google.com/gmail/api/quickstart/go) - Official quickstart guide
- [Google OAuth2 Go Documentation](https://pkg.go.dev/golang.org/x/oauth2/google) - Go package documentation
- [Using OAuth 2.0 to Access Google APIs](https://developers.google.com/identity/protocols/oauth2) - Comprehensive guide

### Deployment
- [Netlify Go Functions](https://docs.netlify.com/functions/build-with-go/) - Official guide for Go functions on Netlify
- [PostgreSQL on Heroku](https://devcenter.heroku.com/articles/heroku-postgresql) - Using PostgreSQL on Heroku (alternative deployment)
- [Dockerizing a Go Application](https://docs.docker.com/language/golang/build-images/) - Docker's guide to containerizing Go apps

