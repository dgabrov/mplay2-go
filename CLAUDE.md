# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

You are a senior programmer who does what I tell you. Thanks for that. 

## Project Overview

**mplay2-go** is a Go backend application for managing media and playlists. The project is a Go implementation of a media management system that supports users, media files, playlists, and sessions/authentication.

### Key Features (from schema)
- **User Management**: User accounts with authentication via sessions
- **Media Management**: Upload and manage media files (images/videos) with metadata
- **Playlists**: Organize media into playlists with sequencing
- **Session-based Auth**: Token-based session management with expiry

for log you use slog of course

## Build, Run, and Test Commands

### Build
```bash
# Build the executable
go build -o mplay2-go .

# Build with specific output
go build -o bin/mplay2-go .
```

### Run
```bash
# Run directly from source
go run .

# Run the built executable
./mplay2-go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run a specific test
go test -run TestName ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Quality
```bash
# Format code
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run ./...

# Run go vet
go vet ./...
```

### Dependencies
```bash
# Download and verify dependencies
go mod download
go mod verify

# Tidy up dependencies
go mod tidy

# View dependencies
go list -m all
```

## Project Structure

The project is currently in early development with a minimal structure:

```
mplay2-go/
├── main.go                 # Application entry point
├── go.mod                  # Go module definition
├── docs/
│   └── db/
│       └── db.sql         # Database schema
└── CLAUDE.md              # This file
```

### Future Expected Structure

As the project develops, expect:
- `cmd/` - Command-line application entry points
- `internal/` - Internal packages (handlers, services, models, repository)
- `pkg/` - Public packages (if shared as a library)
- `test/` - Test utilities and fixtures
- `migrations/` - Database migration files

## Database Schema

The application uses a MySQL database with the following core tables:

### Users
- `user_id` (PK): System-generated unique user identifier
- `provided_user_id` (UNIQUE): External user identifier
- `login`: User login name
- `name`: User display name

### Media
- `media_id` (PK): System-generated unique media identifier
- `user_id` (FK): Owner of the media
- `description`: Media description
- `content_type`: MIME type of the media
- `size`, `width`, `height`: Media metadata

### Playlists
- `playlist_id` (PK): System-generated unique playlist identifier
- `user_id` (FK): Owner of the playlist
- `description`: Playlist description

### Media-Playlist Association
- `media_playlist_id` (PK): Unique association record ID
- `playlist_id`, `media_id` (FKs): Many-to-many association
- `seq_no`: Sequence number for ordering within playlist

### Sessions
- `session_id` (PK): Unique session identifier
- `user_id` (FK): Associated user
- `token`: Session token (128 char)
- `expired_ind`: Expiration flag ('Y'/'N')
- `expiry_dt`: Expiration datetime

### Sequence Values
- Simple table for sequence generation (legacy pattern, may be replaced by database native features)

Database is UTF8MB4 encoded for full Unicode support.

## Development Notes

### Go Version
The project targets **Go 1.26**. Ensure your local environment runs this version or later.

### Database Setup
Before running the application, initialize the database using `docs/db/db.sql`:
```bash
mysql < docs/db/db.sql
# or in MySQL CLI:
# mysql> source docs/db/db.sql;
```

### Configuration
Currently, the application has minimal configuration. Future instances should look for:
- Environment variables for database connection
- Configuration files (likely YAML or TOML)
- Command-line flags

### Authentication Pattern
The application uses token-based session authentication stored in the `session` table. Sessions include expiry tracking.

## Architecture Notes

### Current State
The application is in early development with only the data model defined. Implementation of HTTP handlers, services, and repositories is ahead.

### Expected Architecture
Based on the schema and typical Go patterns:
1. **Models/Domain**: User, Media, Playlist, Session entities
2. **Repository Layer**: Database access for each domain entity
3. **Service Layer**: Business logic (playlist management, media handling, auth)
4. **HTTP Handlers**: RESTful API endpoints
5. **Middleware**: Authentication, error handling, logging

### ID Generation
The schema uses string-based IDs (likely UUIDs). Consider using a library like `github.com/google/uuid` for consistent ID generation across the application.

## Common Development Tasks

### Adding a New Feature
1. Define database schema changes in `docs/db/db.sql` if needed
2. Create domain models in `internal/models/`
3. Implement repository interfaces and SQLx/database/sql implementations
4. Add service layer business logic
5. Expose via HTTP handlers
6. Write tests at each layer

### primary key
All varchar primary keys are to be v7 uuid.

### Making Database Changes
1. Update schema in `docs/db/db.sql`
2. Plan migration strategy (consider migration tools like `flyway`, `migrate`, or `golang-migrate`)
3. Update domain models
4. Update repository implementations

### Debugging
```bash
# Run with debug logging (if logger supports DEBUG env var)
DEBUG=true go run .

# Run tests with more detail
go test -v -race ./...
```

## Dependencies to Consider

As the project develops, evaluate these common Go libraries:
- **HTTP**: `github.com/gorilla/mux` or `chi` for routing
- **Database**: `database/sql` with `github.com/go-sql-driver/mysql` for MySQL
- **Validation**: `github.com/go-playground/validator`
- **Logging**: `github.com/sirupsen/logrus` or `go.uber.org/zap`
- **Config**: `github.com/spf13/viper`
- **UUID**: `github.com/google/uuid`
- **Testing**: `github.com/stretchr/testify` for assertions

# Database Access

- All database access points are done through the type Servr defined in @internal/service/servr.go
- The top level function you add opens a db transaction and then takes care of having it committed / rolled back
- The top level function you add has first parameter a context that is created on topmost level
- Any direct database access - one "hit", one sql, you do in separate private function created preferrably in file @internal/service/dao.go that takes transaction as one of the parameters
- before you want to create a new value object, look in package data under @internal/data maybe the type is already there. If not there, create it there in business.go
