# Go Backend

A Go-based backend application with user management, room functionality, and WebSocket support.

## Setup

### Prerequisites

- Go (1.16 or later)
- MySQL
- Make (for running commands)

### Environment Configuration

Create a `.env` file in the root directory with the following content:

```
PUBLIC_HOST=http://localhost
APP_PORT=3000
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_HOST=127.0.0.1
DB_NAME=go_backend
```

Adjust the values according to your environment setup.

### Installation

Install all required dependencies:

```bash
make init
```

## Database Migrations

This project uses the `golang-migrate` library for database migrations.

### Available Migration Commands

- Create a new migration:
  ```bash
  make migration <migration_name>
  ```
  This creates two files in `cmd/migrate/migrations/`: `<timestamp>_<migration_name>.up.sql` and `<timestamp>_<migration_name>.down.sql`

- Apply all pending migrations:
  ```bash
  make migrate-up
  ```

- Rollback the last migration:
  ```bash
  make migrate-down
  ```

- Check migration status:
  ```bash
  make migrate-status
  ```
  This shows the current migration version and whether the database is in a "dirty" state.

- Force migration version (useful for fixing issues):
  ```bash
  make migrate-force VERSION=<version_number>
  ```

### Troubleshooting Migrations

If you encounter issues with migrations:

1. Check the current migration status:
   ```bash
   make migrate-status
   ```

2. If the database is in a "dirty" state or migrations are out of sync:
   - Identify the last successful migration version
   - Force the migration to that version:
     ```bash
     make migrate-force VERSION=<last_successful_version>
     ```
   - Run migrations again:
     ```bash
     make migrate-up
     ```

## Running the Application

Start the server:

```bash
make run
```

The server will start on the port specified in your `.env` file (default: 3000).

## Project Structure

- `cmd/`: Contains all entry points for the application
  - `main.go`: Main application entry point
  - `migrate/`: Database migration tools
- `config/`: Configuration management
- `db/`: Database connection and utilities
- `service/`: Business logic and API handlers
  - `user/`: User management functionality
  - `rooms/`: Room management functionality
  - `websocket/`: WebSocket implementation

## Testing

Run tests:

```bash
make test
```

