# ana

A project repository with standardized monitoring, documentation, and AI integration.

## Features

- Standardized repository structure
- Integrated monitoring system
- Health checks and status reporting
- Documentation cross-references
- CI/CD integration
- Cerebras AI integration for intelligent assistance

## Technology Stack

- **Backend**: Go
- **Frontend**: HTML/CSS/JavaScript
- **Database**: MongoDB
- **AI Integration**: Cerebras API
- **Documentation**: Markdown
- **Monitoring**: Prometheus/Grafana
- **Testing**: Go testing framework

## Project Structure

```
ana/
├── .github/                  # GitHub workflows and configuration
├── cmd/                      # Application entry points
│   └── server/               # Main server application
├── config/                   # Configuration files
├── db/                       # Database migrations and schemas
├── internal/                 # Internal packages
│   ├── ai/                   # AI integration (Cerebras)
│   ├── api/                  # API handlers and routes
│   ├── config/               # Configuration handling
│   ├── database/             # Database access
│   ├── handlers/             # Request handlers
│   ├── logging/              # Logging utilities
│   ├── metrics/              # Metrics collection
│   ├── middleware/           # HTTP middleware
│   ├── models/               # Data models
│   ├── monitoring/           # Monitoring utilities
│   ├── repositories/         # Data repositories
│   └── server/               # Server setup
├── monitoring/               # Monitoring configuration
│   ├── grafana/              # Grafana dashboards
│   └── prometheus/           # Prometheus configuration
├── scripts/                  # Utility scripts
├── web/                      # Web assets
├── .env.example              # Example environment variables
├── .gitignore                # Git ignore file
├── cerebras_api.md           # Cerebras API reference
├── cerebras_integration.md   # Cerebras integration documentation
├── cerebras_models.md        # Cerebras models documentation
├── DEVELOPMENT.md            # Development guidelines
├── docker-compose.yml        # Docker Compose configuration
├── Dockerfile                # Docker configuration
├── go.mod                    # Go module file
├── go.sum                    # Go dependencies checksum
├── MONITORING.md             # Monitoring documentation
└── README.md                 # This file
```

## Monitoring

This repository includes a standardized monitoring system that tracks:

- Repository health status
- Code quality metrics
- Build and deployment status
- API endpoint availability
- Service health checks

For more details, see [MONITORING.md](MONITORING.md).

## Cerebras AI Integration

The project integrates with the Cerebras AI API to provide intelligent assistance for architectural queries. Key features include:

- Support for multiple Cerebras AI models
- Response caching for improved performance
- Circuit breaker pattern for fault tolerance
- Detailed metrics and monitoring

For more information, see:
- [cerebras_integration.md](cerebras_integration.md) - Overview of the integration
- [cerebras_models.md](cerebras_models.md) - Details about available models
- [cerebras_api.md](cerebras_api.md) - API reference for developers

## Getting Started

1. Clone the repository
   ```bash
   git clone https://github.com/lyffseba/ana.git
   cd ana
   ```

2. Set up environment variables
   ```bash
   cp .env.example .env
   # Edit .env to add your Cerebras API key and other configuration
   ```

3. Run the server
   ```bash
   go run cmd/server/main.go
   ```

4. Test the Cerebras AI integration
   ```bash
   # Using the provided test script
   ./test_cerebras.sh
   
   # Or using curl directly
   curl -X POST http://localhost:8081/api/cerebras/assistant \
     -F "query=What are the key considerations for sustainable architecture?" \
     -F "model_type=qwen-3-32b"
   ```

5. Access the web interface
   ```bash
   # Open in your browser
   open http://localhost:8081
   ```

## Contact

For questions or support, please contact:
- Email: juansg777@gmail.com
- GitHub: [@lyffseba](https://github.com/lyffseba)

## License

This project is proprietary and intended for specific use.

## Development Session
Reference: https://app.warp.dev/session/b660fd8a-f765-449c-a70c-f8c7b971e3c4?pwd=e9ccd7cb-d8be-494e-a2f2-35469f726896
Last Updated: Sat May 17 07:34:44 AM CEST 2025


## Development Sessions
- Warp Session: https://app.warp.dev/session/7b7b92ff-9036-49e7-bef1-8dfa1e2d2e39?pwd=170af286-10c6-4dbd-a043-329cd22b8d4c
- Grok Session: https://grok.com/workspace/aa440a48-7734-41b1-9cd1-46994ad7d6d2

Last updated: 2025-05-28 09:17:33
