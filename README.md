# ana.world

A project management application designed for architects managing multiple projects. The application provides:

- Daily agenda and task management
- Comprehensive project views
- Google Calendar integration
- Supplier search functionality
- Financial and document tracking

## Technology Stack

- **Frontend**: HTMX with minimal JavaScript, styled with Tailwind CSS
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL
- **Hosting**: Netlify for frontend, with Go serverless functions or GCP for backend

## Project Structure

```
ana/
├── cmd/           # Application entry points
│   └── server/    # Main server binary
├── internal/      # Internal packages
│   ├── handlers/  # HTTP handlers
│   ├── models/    # Data models
│   └── server/    # Server configuration
├── web/           # Frontend static files
├── functions/     # Netlify serverless functions
├── db/            # Database migrations and schemas
└── config/        # Configuration files
```

## Development Setup

### Prerequisites

- Go 1.16+
- PostgreSQL
- Node.js (optional, for Netlify CLI)

### Local Development

1. Clone the repository:
   ```
   git clone https://github.com/sebae/ana.git
   cd ana
   ```

2. Install Go dependencies:
   ```
   go mod download
   ```

3. Run the server:
   ```
   go run cmd/server/main.go
   ```

4. Access the web interface at http://localhost:8080

## Deployment

### Netlify Deployment

1. Push changes to GitHub
2. Netlify will automatically build and deploy from the connected repository

## Features

- **Daily Agenda**: View and manage daily tasks
- **Project Management**: Track project details, progress, and documents
- **Google Calendar Integration**: Sync tasks with calendar
- **Supplier Search**: Find architecture material suppliers in Bogota
- **Financial Tracking**: Monitor project budgets and expenses

## License

This project is proprietary and not licensed for public use.

