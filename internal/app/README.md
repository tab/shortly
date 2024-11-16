# Internal documentation

## Overview

The `internal/app` directory contains the core logic and structure of the *shortly* application.

## Directory Structure

```sh
app/
  ├── api/        # API handlers for routing HTTP requests
  ├── config/     # Configuration management
  ├── dto/        # Data Transfer Objects (input/output models)
  ├── errors/     # Application-specific error definitions
  ├── repository/ # Interfaces and implementations for data storage
  ├── router/     # HTTP router and middleware setup
  ├── service/    # Business logic and core services
  ├── server/     # HTTP server setup
  ├── validator/  # Input validation utilities
  ├── app.go      # Application bootstrap and lifecycle management
```

## API Layer (`api/`)
- Defines handlers for each API endpoint.
- Interfaces with services to perform operations.
- Converts input and output data into user-friendly formats.

### Key Handlers
- **HealthHandler**: Provides health checks (e.g., `/ping`).
- **URLHandler**: Manages URL shortening, retrieval, and batch operations.

## Configuration (`config/`)
- Loads configuration from environment variables and `.env` files.
- Supports multiple environments (development, test, production).
- Uses flags as an additional configuration source.

## Repositories (`repository/`)
- Abstracts data storage with interfaces.
- Supports multiple implementations:
  - **InMemory**: For testing and local development.
  - **File**: For backup and persistence.
  - **Database**: For PostgreSQL storage.

### Memento Pattern
The in-memory repository supports snapshotting via the memento pattern, enabling saving and restoring states.

## Services (`service/`)
- Implements business logic.
- Examples:
  - `URLService`: Handles URL shortening and retrieval.
  - `HealthService`: Checks system health.

## Server (`server/`)
- Sets up the HTTP server with customizable timeouts and middleware.

## Router (`router/`)
- Configures the routing and middleware stack.
- Includes:
  - Logging middleware.
  - Request compression.
  - CORS support.

## Utilities
- **Validation** (`validator/`): Validates input data (e.g., URL format).
- **Logging** (`logger/`): Provides structured logging with contextual information.

## Getting Started
To understand how the application is bootstrapped, refer to:
- `app.go`: The main entry point for initializing the application.
- `api/`: Entry points for HTTP requests.
