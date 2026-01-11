# Manabu Service - Japanese Learning Application

## Description

API backend for Manabu, a comprehensive Japanese language learning application. This service provides user management, authentication, and will support vocabulary, lessons, exercises, and progress tracking features.

## Features

- ✅ User authentication & authorization (JWT-based)
- ✅ Role-based access control
- ✅ RESTful API architecture
- ✅ Interactive API documentation (Swagger/OpenAPI)
- ✅ Rate limiting
- ✅ CORS support
- ✅ PostgreSQL database
- ✅ Clean architecture (layered structure)
- ✅ **Vocabulary management** (CRUD operations)
- ✅ **JLPT level classification** (N5-N1)
- ✅ **Category management** (Vocabulary categorization)
- ✅ **Tags system** (Vocabulary tagging)
- ✅ **Progress tracking** (Simple learning status tracking)
- 🚧 Courses & lessons (coming soon)
- 🚧 Quiz & exercises (coming soon)

## Tech Stack

- **Language:** Go 1.24
- **Framework:** Gin (Web framework)
- **Database:** PostgreSQL
- **ORM:** GORM
- **Authentication:** JWT (golang-jwt/jwt)
- **Documentation:** Swagger/OpenAPI (swaggo)
- **Rate Limiting:** Tollbooth
- **Configuration:** Viper
- **Hot Reload:** Air

## Directory Structure

```
manabu-service
    ├── cmd                          → Main entry point and application configuration
    ├── common                       → Common utilities and helper functions
    ├── config                       → Application configuration (database, env vars)
    ├── constants                    → Global constant values
    ├── controllers                  → HTTP request handlers
    ├── database
    │   └── seeders                  → Database seed data scripts
    ├── docs                         → API documentation & ERD
    ├── domain
    │   ├── dto                      → Data Transfer Objects
    │   └── models                   → Database entity models
    ├── middlewares                  → HTTP middlewares (auth, rate limit, etc.)
    ├── repositories                 → Data access layer
    ├── routes                       → API route definitions
    └── services                     → Business logic layer
```

## Prerequisites

- Go 1.24 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Setup

### 1. Clone the repository

```bash
git clone <repository-url>
cd manabu-service
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Configure the application

```bash
# Copy configuration file
cp config.json.example config.json

# Edit config.json with your settings
# - Database credentials
# - JWT secret key
# - Port (default: 8001)
```

### 4. Setup database

Create a PostgreSQL database:

```sql
CREATE DATABASE manabu_service;
```

Update `config.json` with your database credentials:

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "manabu_service",
    "username": "postgres",
    "password": "your_password"
  }
}
```

## Running the Application

### Development Mode (with hot reload)

```bash
# First time only - install air for hot reload
make watch-prepare

# Run with hot reload (auto-generates Swagger docs)
make watch
```

### Development Mode (without hot reload)

```bash
# Run directly (auto-generates Swagger docs)
make run

# Or manually
go generate  # Generates Swagger docs
go run main.go serve
```

### Production Mode

```bash
# Build the application (auto-generates Swagger docs)
make build

# Run the binary
./manabu-service serve
```

### With Docker

```bash
docker-compose up -d --build --force-recreate
```

The API will be available at `http://localhost:8001`

## API Documentation

### Swagger UI

Interactive API documentation is available via Swagger UI:

```
http://localhost:8001/swagger/index.html
```

### Generate Swagger Documentation

**Automatic Generation:**
Swagger docs are automatically generated when you:
- Run the app: `make run` or `make watch`
- Build the app: `make build`
- Execute: `go generate`

**Manual Generation:**
If you need to generate manually:

```bash
make swagger
# or
go generate
```

This will regenerate the Swagger documentation files in the `docs/` folder.

**Note:** The `//go:generate` directive in `main.go` ensures docs are always up-to-date.

### API Endpoints

#### Authentication

- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `GET /api/v1/auth/user` - Get logged in user (requires authentication)
- `GET /api/v1/auth/{uuid}` - Get user by UUID (requires authentication)
- `PUT /api/v1/auth/{uuid}` - Update user profile (requires authentication)

#### Vocabulary Management

- `GET /api/v1/vocabularies` - Get all vocabularies (with pagination & filters)
- `GET /api/v1/vocabularies/{id}` - Get vocabulary by ID
- `POST /api/v1/vocabularies` - Create new vocabulary (admin only)
- `PUT /api/v1/vocabularies/{id}` - Update vocabulary (admin only)
- `DELETE /api/v1/vocabularies/{id}` - Delete vocabulary (admin only)

#### JLPT Levels

- `GET /api/v1/jlpt-levels` - Get all JLPT levels (N5-N1)
- `GET /api/v1/jlpt-levels/{id}` - Get JLPT level by ID
- `POST /api/v1/jlpt-levels` - Create JLPT level (admin only)
- `PUT /api/v1/jlpt-levels/{id}` - Update JLPT level (admin only)
- `DELETE /api/v1/jlpt-levels/{id}` - Delete JLPT level (admin only)

#### Categories

- `GET /api/v1/categories` - Get all categories
- `GET /api/v1/categories/{id}` - Get category by ID
- `POST /api/v1/categories` - Create category (admin only)
- `PUT /api/v1/categories/{id}` - Update category (admin only)
- `DELETE /api/v1/categories/{id}` - Delete category (admin only)

#### Tags

- `GET /api/v1/tags` - Get all tags
- `GET /api/v1/tags/{id}` - Get tag by ID
- `GET /api/v1/tags/search` - Search tags by name
- `POST /api/v1/tags` - Create tag (admin only)
- `PUT /api/v1/tags/{id}` - Update tag (admin only)
- `DELETE /api/v1/tags/{id}` - Delete tag (admin only)

#### User Vocabulary Status (Progress Tracking)

- `POST /api/v1/user-vocabulary-status` - Start learning a vocabulary
- `GET /api/v1/user-vocabulary-status` - Get all learning progress (paginated)
- `GET /api/v1/user-vocabulary-status/due` - Get vocabularies due for review
- `GET /api/v1/user-vocabulary-status/{id}` - Get specific progress by ID
- `POST /api/v1/user-vocabulary-status/{vocabulary_id}/review` - Review a vocabulary

📖 **Detailed API Documentation:** [User Vocabulary Status API Guide](docs/USER_VOCABULARY_STATUS.md)

#### Testing with Swagger

1. **Login** to get JWT token via `/auth/login`
2. Click **"Authorize"** button in Swagger UI
3. Enter: `Bearer <your-token-here>`
4. Click **"Authorize"** then **"Close"**
5. Now you can test protected endpoints

## Development

### Project Commands

```bash
# Install hot reload tool (first time only)
make watch-prepare

# Run with hot reload (auto-generates Swagger)
make watch

# Run without hot reload (auto-generates Swagger)
make run

# Build the application (auto-generates Swagger)
make build

# Generate Swagger docs only (manual)
make swagger
# or
go generate

# Run with Docker
make docker-compose

# Build Docker image with tag
make docker-build tag=1.0.0

# Push Docker image
make docker-push tag=1.0.0
```

### Database Migrations

The application uses a **hybrid approach** for database migrations:

1. **GORM AutoMigrate** - For automatic schema creation on first run
2. **SQL Migrations** - For schema modifications after initial deployment

#### Automatic Migrations (GORM)

GORM AutoMigrate runs automatically on application start:

```go
// In cmd/main.go
db.AutoMigrate(
    &models.Role{},
    &models.User{},
    &models.JlptLevel{},
    &models.Category{},
    &models.Vocabulary{},
    &models.Tag{},
    &models.UserVocabularyStatus{},
)
```

#### Manual Migrations (SQL)

For schema modifications after deployment, use SQL migration files in the `migrations/` folder.

**Run a migration:**

```bash
# Method 1: Using Go tool (recommended)
go run tools/migrate.go <migration_file>

# Example
go run tools/migrate.go 002_rename_users_uuid_constraint.sql

# Method 2: Using Windows batch script
migrate.bat 002_rename_users_uuid_constraint.sql

# Method 3: Using psql directly
psql -U postgres -d manabu -f migrations/<migration_file>
```

**Check migration status:**

```bash
# Verify database constraints
go run tools/check_constraints.go
```

**Available Migrations:**
- `001_alter_user_vocabulary_status_user_id_to_uuid.sql` - Changes user_id to UUID type
- `002_rename_users_uuid_constraint.sql` - Fixes constraint naming for GORM
- `003_update_user_vocabulary_status_check_constraint.sql` - Updates status constraint from SM-2 algorithm to simple tracking

📖 **Full documentation:** See [migrations/README.md](migrations/README.md) for:
- Creating new migrations
- Migration best practices
- Rollback strategies
- Troubleshooting guide

### Seed Data

Seed data is automatically populated on first run:

- Default roles (admin, student)
- Admin user (if configured)

## Configuration

### config.json

```json
{
  "port": 8001,
  "appName": "manabu-service",
  "appEnv": "local",
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "manabu_service",
    "username": "postgres",
    "password": "your_password",
    "maxOpenConnection": 10,
    "maxLifetimeConnection": 10,
    "maxIdleConnection": 10,
    "maxIdleTime": 10
  },
  "rateLimiterMaxRequest": 1000,
  "rateLimiterTimeSecond": 60,
  "jwtSecretKey": "your-secret-key-here",
  "jwtExpirationTime": 1440
}
```

### Configuration Options

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `port` | int | Server port | 8001 |
| `appName` | string | Application name | manabu-service |
| `appEnv` | string | Environment (local, dev, prod) | local |
| `jwtSecretKey` | string | Secret key for JWT signing | - |
| `jwtExpirationTime` | int | JWT expiration in minutes | 1440 (24 hours) |
| `rateLimiterMaxRequest` | float64 | Max requests per time window | 1000 |
| `rateLimiterTimeSecond` | int | Rate limiter time window (seconds) | 60 |

## Documentation

- [ERD (Entity Relationship Diagram)](docs/ERD.md)
- [API Gaps Analysis](docs/API_GAPS.md)
- [API Registration Documentation](docs/API_REGISTRATION.md)
- [Swagger Guide](docs/SWAGGER_GUIDE.md)
- **[Database Migrations](MIGRATIONS.md)**
  - [Quick Reference Guide](MIGRATIONS.md)
  - [Complete Documentation](migrations/README.md)
- **[User Vocabulary Status API](docs/USER_VOCABULARY_STATUS.md)** 🆕
  - Complete guide for progress tracking feature
  - API usage examples and workflows

## Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
Controllers (HTTP Layer)
    ↓
Services (Business Logic)
    ↓
Repositories (Data Access)
    ↓
Database (PostgreSQL)
```

## Security

- ✅ JWT-based authentication
- ✅ Password hashing with bcrypt
- ✅ Rate limiting to prevent abuse
- ✅ CORS configured
- ✅ Input validation
- ✅ SQL injection protection (via GORM)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Apache 2.0

## Support

For support, email support@manabu.com or open an issue in the repository.

