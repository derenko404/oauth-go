# oauth-go

A Golang OAuth 2.0 server implementation designed to facilitate secure authorization for web, mobile, and desktop applications.

## Features

- Implements the OAuth 2.0 protocol for secure authorization.
- Structured with modular components for scalability and maintainability.
- Includes database migrations for setting up the necessary schema.
- Provides a Docker Compose configuration for easy deployment.

## API Documentation

Interactive API documentation is available via Swagger UI. After starting the application, navigate to: http://{host}/swagger/index.html

## Getting Started

To get started with the `oauth-go` project, follow these steps:

### Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/oauth-go.git
   cd oauth-go
   ```

2. Create a `.env` file in the root directory or update the existing one with your PostgreSQL credentials:

   ```ini
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=oauth_db
   ```

3. To start the application in development mode, use the following command to start the Docker containers:

   ```bash
   make start-dev
   ```

4. Build the application:

   ```bash
   make build
   ```

5. Run the application:

   ```bash
   make run
   ```

6. Access the Swagger API documentation by navigating to:  
   [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

## Makefile Commands

The Makefile provides several commands to help with common tasks:

- **all**: Builds the project.

  ```bash
  make all
  ```

- **run**: Runs the application in the current environment.

  ```bash
  make run
  ```

- **start-dev**: Starts the application with Docker Compose in development mode.

  ```bash
  make start-dev
  ```

- **build**: Builds the application binary.

  ```bash
  make build
  ```

- **test**: Runs tests for the project.

  ```bash
  make test
  ```

- **clean**: Cleans the build directory.

  ```bash
  make clean
  ```

- **format**: Formats the code using `gofmt`.

  ```bash
  make format
  ```

- **migrate-create**: Creates a new migration file.

  ```bash
  make migrate-create name=create_table_name
  ```

- **migrate-up**: Applies the database migrations up to the specified version.

  ```bash
  make migrate-up
  ```

- **migrate-down**: Rolls back the database migrations by the specified version.

  ```bash
  make migrate-down
  ```

- **swagger**: Generates the Swagger documentation for the API.
  ```bash
  make swagger
  ```

## Database Migrations

This project uses the [migrate](https://github.com/golang-migrate/migrate) tool for managing database schema changes.

### Creating a Migration

To create a new migration, use the `migrate-create` command:

```bash
make migrate-create name=create_users_table
```
