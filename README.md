# mplay2-go

A Go backend application for managing media and playlists with user authentication and session management.

## Features

- **User Management**: Create and manage user accounts with external ID integration
- **Media Management**: Upload and organize media files (images/videos) with metadata
- **Playlists**: Create and manage playlists with customizable media ordering
- **Session-based Authentication**: Token-based authentication with expiry tracking
- **Search**: Search media and playlists by description

## Prerequisites

- Go 1.26 or later
- MySQL 5.7+ database
- Environment setup for database connection configuration

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd mplay2-go
   ```

2. Download dependencies:
   ```bash
   go mod download
   go mod verify
   ```

3. Initialize the database:
   ```bash
   mysql < docs/db/db.sql
   ```
   Or from MySQL CLI:
   ```sql
   source docs/db/db.sql;
   ```

## Running the Application

Build and run:
```bash
go run .
```

Or build then execute:
```bash
go build -o mplay2-go .
./mplay2-go
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with verbose output:
```bash
go test -v ./...
```

Run tests with coverage:
```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Code Quality

Format code:
```bash
go fmt ./...
```

Lint (requires golangci-lint):
```bash
golangci-lint run ./...
```

Run go vet:
```bash
go vet ./...
```

## Project Structure

```
mplay2-go/
├── main.go                    # Application entry point
├── go.mod                     # Go module definition
├── go.sum                     # Dependency checksums
├── docs/
│   └── db/
│       └── db.sql           # Database schema
├── internal/
│   ├── data/                # Domain models and value objects
│   │   └── business.go      # Business entity definitions
│   └── service/             # Service layer
│       ├── servr.go         # Service orchestrator with transaction management
│       └── dao.go           # Data access operations
└── CLAUDE.md                # Development guidelines
```

## Database Schema

The application uses MySQL with UTF8MB4 encoding. Core tables:

- **user**: User accounts with external ID mapping
- **media**: Media files with metadata (dimensions, content type, size)
- **playlist**: User playlists
- **media_playlist**: Many-to-many association with sequence ordering
- **session**: Session tokens with expiry tracking
- **seqvalues**: Sequence generation for media ordering

All primary keys use UUID v7 format.

## Architecture

### Service Layer (`Servr`)
All database access flows through the `Servr` type which:
- Opens a transaction for each operation
- Handles transaction commit/rollback
- Delegates to private DAO functions for specific database operations

### Data Access (`dao.go`)
- Private functions handle individual database operations
- Operate within a transaction context
- Pattern: `functionName(ctx context.Context, tx *sql.Tx, ...args)`

### Models
Domain entities are defined in `internal/data/business.go`:
- User
- Media
- PlayList
- ExtendedMedia (media with sequence info)
- Session

## Development

### Adding a New Feature

1. Update schema in `docs/db/db.sql` if needed
2. Add domain models in `internal/data/business.go`
3. Implement DAO function in `internal/service/dao.go`
4. Add service method in `internal/service/servr.go`
5. Write tests for each layer
6. Expose via HTTP handlers

### Common Tasks

**Creating a new value object:**
Check `internal/data/business.go` first—the type may already exist.

**Making database changes:**
Update `docs/db/db.sql`, then update domain models and DAO implementations.

**Debug mode:**
```bash
DEBUG=true go run .
```

## Dependencies

Key dependencies:
- `database/sql`: Standard Go database interface
- `github.com/go-sql-driver/mysql`: MySQL driver
- `github.com/google/uuid`: UUID v7 generation
- `log/slog`: Structured logging (Go standard library)

## License

See LICENSE file for details.
