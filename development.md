# Development Guidelines for ana.world

This document outlines the standards, patterns, and practices to maintain consistency throughout the development of the ana.world project.

## Documentation Standards

### General Documentation

- **Bilingual Requirement**: All user-facing documentation must be in both English and Spanish.
- **Documentation Files**: 
  - `README.md` - Project overview and setup instructions
  - `plan.md` - Comprehensive project plan 
  - `next.md` - Current implementation steps
  - `flow.md` - Project evolution and development process
  - `development.md` - Development standards (this document)

### Code Documentation

- **Function Comments**: All exported functions must have comments in the following format:
  ```go
  // FunctionName does X when given Y.
  // It returns Z or an error if something goes wrong.
  func FunctionName() {}
  ```

- **Package Documentation**: Each package should have a package comment:
  ```go
  // Package name provides functionality for...
  package name
  ```

- **Complex Logic**: Any complex logic should be commented with explanations.

### API Documentation 

- **Endpoint Documentation**: All API endpoints should be documented with:
  - HTTP method 
  - Route
  - Request parameters/body
  - Response format
  - Example request/response
  - Possible errors

## Development Workflow Patterns

### Feature Development Process

1. **Planning**:
   - Review `plan.md` for feature requirements
   - Update `next.md` with detailed implementation steps if needed

2. **Branch Creation**:
   - Create a feature branch: `git checkout -b feature/feature-name`

3. **Implementation**:
   - Backend implementation first (API endpoints)
   - Frontend implementation second
   - Tests implementation third

4. **Documentation**:
   - Update relevant documentation
   - Add bilingual user documentation if needed

5. **Testing**:
   - Run unit tests
   - Manual testing of the feature
   - Verify both Spanish and English interfaces

6. **Pull Request**:
   - Submit a pull request with a comprehensive description
   - Ensure all CI checks pass

7. **Review & Merge**:
   - Address review comments
   - Merge to main branch once approved

### Issue Management

- Create detailed issue descriptions including:
  - Problem statement
  - Expected behavior
  - Steps to reproduce
  - Suggested implementation (if applicable)

## Commit Message Conventions

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation changes
- `style`: Changes that don't affect code functionality (formatting, etc.)
- `refactor`: Code changes that neither fix bugs nor add features
- `perf`: Performance improvements
- `test`: Adding or correcting tests
- `chore`: Changes to build process, tools, etc.

### Examples

```
feat(tasks): add due date filtering functionality

fix(auth): resolve token refresh issue

docs: update PostgreSQL setup instructions in Spanish and English
```

## File Organization Rules

### Directory Structure

Maintain the established directory structure as outlined in the README:

```
ana/
├── cmd/                    # Application entry points
│   └── server/             # Main server binary
├── internal/               # Internal packages
│   ├── handlers/           # HTTP handlers
│   ├── models/             # Data models
│   ├── server/             # Server configuration
│   ├── database/           # Database connection
│   └── repositories/       # Data access layer
├── web/                    # Frontend static files
├── functions/              # Netlify serverless functions
├── db/                     # Database migrations and schemas
│   └── migrations/         # SQL migration files
└── config/                 # Configuration files
```

### Naming Conventions

- **Files**: Use snake_case for filenames (e.g., `task_handler.go`)
- **Go Packages**: Use lowercase, single-word names (e.g., `models`)
- **Go Structs/Interfaces**: Use PascalCase (e.g., `TaskRepository`)
- **Frontend Components**: Use kebab-case for file names (e.g., `task-list.html`)
- **Database Tables**: Use snake_case (e.g., `user_projects`)

## Quality Assurance Steps

### Code Review Checklist

- [ ] Code follows established patterns and naming conventions
- [ ] New functions/methods have appropriate documentation
- [ ] Error handling is implemented correctly
- [ ] Tests are included for new functionality
- [ ] Security considerations have been addressed
- [ ] Performance impact has been considered
- [ ] Bilingual support is maintained where applicable

### Testing Requirements

- **Unit Tests**: All backend packages should have unit tests with >70% coverage
- **Integration Tests**: API endpoints should have integration tests
- **Frontend Testing**: Manual testing of HTMX interactions
- **Database Testing**: Test migration scripts in a test database before applying to development

### Security Checks

- Input validation on all API endpoints
- Parameterized SQL queries
- Authentication/authorization checks
- Secure handling of sensitive data

## Project Consistency Guidelines

### Coding Style

#### Go

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `go fmt` before committing code
- Prefer standard library packages when possible
- Keep functions small and focused on a single responsibility

#### SQL

- Use uppercase for SQL keywords
- Use parameterized queries
- Document complex queries

#### Frontend

- Follow HTMX patterns for dynamic content loading
- Maintain consistent Tailwind CSS usage for styling
- Use Spanish language in user-facing elements with English comments

### Configuration Management

- Use environment variables for configuration
- Document all configuration options in both languages
- Provide sensible defaults for development

### Version Control

- Make small, focused commits
- Use feature branches for development
- Rebase feature branches before merging
- Squash commits when appropriate

## Implementation Phases

Follow the implementation phases as defined in the project plan:

1. **Phase 1**: Setup and Core Features
2. **Phase 2**: Project Management and Agenda
3. **Phase 3**: Google Calendar Integration
4. **Phase 4**: Supplier Search
5. **Phase 5**: Dashboards and Enhancements
6. **Phase 6**: Deployment and Maintenance

Document phase transitions in the `flow.md` file.

## Deployment Practices

- Test in development environment before deploying to production
- Document deployment steps in both languages
- Maintain a deployment checklist
- Use environment-specific configuration

---

By following these guidelines, we ensure a consistent, high-quality, and well-documented development process for ana.world. These standards will evolve as the project grows, with updates documented in this file.

