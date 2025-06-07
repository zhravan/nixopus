# Contributing to Nixopus Backend

This guide provides detailed instructions for contributing to the Nixopus backend codebase.

## Setup for Backend Development

1. **Prerequisites**
   - Go version 1.23.6 or newer
   - PostgreSQL
   - Docker and Docker Compose (recommended)

2. **Environment Setup**

   ```bash
   # Clone the repository
   git clone https://github.com/raghavyuva/nixopus.git
   cd nixopus
   
   # Set up PostgreSQL database
   createdb nixopus -U postgres
   createdb nixopus_test -U postgres
   
   # Copy environment template
   cp api/.env.sample api/.env
   
   # Install dependencies
   cd api
   go mod download
   ```

3. **Database Migrations**

   ```bash
   # Run migrations
   cd api
   go run migrations/main.go
   ```

## Project Structure

The backend follows a clean architecture approach:

```
api/
├── internal/
│   ├── features/      # Feature modules
│   ├── middleware/    # HTTP middleware
│   ├── config/        # Application configuration
│   ├── storage/       # Data storage implementation
│   ├── utils/         # Utility functions
│   └── types/         # Type definitions
├── migrations/        # Database migrations
└── tests/             # Test utilities
```

## Adding a New Feature

1. **Create a New Branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Implement Your Feature**

   Create a new directory under `api/internal/features/` with the following structure:

   ```
   api/internal/features/your-feature/
   ├── controller.go   # HTTP handlers
   ├── service.go      # Business logic
   ├── storage.go      # Data access
   └── types.go        # Type definitions
   ```

   Here's a sample implementation:

   ```go
   // types.go
   package yourfeature
   
   type YourEntity struct {
       ID        string `json:"id" db:"id"`
       Name      string `json:"name" db:"name"`
       CreatedAt string `json:"created_at" db:"created_at"`
       UpdatedAt string `json:"updated_at" db:"updated_at"`
   }
   
   // controller.go
   package yourfeature
   
   import (
       "net/http"
       
       "github.com/gin-gonic/gin"
   )
   
   type Controller struct {
       service *Service
   }
   
   func NewController(service *Service) *Controller {
       return &Controller{service: service}
   }
   
   func (c *Controller) GetEntity(ctx *gin.Context) {
       // Implementation
   }
   
   // service.go
   package yourfeature
   
   type Service struct {
       storage *Storage
   }
   
   func NewService(storage *Storage) *Service {
       return &Service{storage: storage}
   }
   
   // storage.go
   package yourfeature
   
   import (
       "database/sql"
   )
   
   type Storage struct {
       db *sql.DB
   }
   
   func NewStorage(db *sql.DB) *Storage {
       return &Storage{db: db}
   }
   ```

3. **Register Routes**

   Update `api/internal/routes.go` to include your new endpoints:

   ```go
   // Register your feature routes
   yourFeatureStorage := yourfeature.NewStorage(db)
   yourFeatureService := yourfeature.NewService(yourFeatureStorage)
   yourFeatureController := yourfeature.NewController(yourFeatureService)
   
   api := router.Group("/api")
   {
       // Your feature endpoints
       featureGroup := api.Group("/your-feature")
       {
           featureGroup.GET("/", middleware.Authorize(), yourFeatureController.GetEntities)
           featureGroup.POST("/", middleware.Authorize(), yourFeatureController.CreateEntity)
           // Add more routes as needed
       }
   }
   ```

4. **Add Database Migrations**

   Create migration files in `api/migrations/your-feature/`:

   ```sql
   -- 20250607000000_create_your_feature_table.up.sql
   CREATE TABLE your_feature (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       name TEXT NOT NULL,
       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
   );
   
   -- 20250607000000_create_your_feature_table.down.sql
   DROP TABLE IF EXISTS your_feature;
   ```

5. **Write Tests**

   Create tests in the same directory:

   ```go
   // controller_test.go
   package yourfeature
   
   import (
       "testing"
       // Import other necessary packages
   )
   
   func TestGetEntity(t *testing.T) {
       // Test implementation
   }
   
   // service_test.go
   // storage_test.go
   ```

6. **Update API Documentation**

   Update the OpenAPI specification in `api/doc/openapi.json` to include your new endpoints.

## Testing

1. **Run Unit Tests**

   ```bash
   cd api
   go test ./internal/features/your-feature/...
   ```

2. **Run Integration Tests**

   ```bash
   cd api
   go test ./internal/features/your-feature/... -tags=integration
   ```

3. **Run All Tests**

   ```bash
   cd api
   make test
   ```

## Optimizing Performance

1. **Use Database Indices** for frequently queried columns
2. **Implement Caching** for expensive operations
3. **Optimize SQL Queries** for better performance
4. **Add Proper Error Handling** and logging

## Code Style and Guidelines

1. Follow Go's [Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
2. Use meaningful variable and function names
3. Add comments for complex logic
4. Structure code for readability and maintainability
5. Follow the existing project patterns

## Common Pitfalls

1. Forgetting to update migrations
2. Not handling database transactions properly
3. Missing error handling
4. Inadequate test coverage
5. Not considering performance implications

## Submitting Your Contribution

1. **Commit Changes**

   ```bash
   git add .
   git commit -m "feat: add your feature"
   ```

2. **Push and Create a Pull Request**

   ```bash
   git push origin feature/your-feature-name
   ```

3. Follow the PR template and provide detailed information about your changes.

## Need Help?

If you need assistance, feel free to:

- Create an issue on GitHub
- Reach out on the project's Discord channel
- Contact the maintainers directly

Thank you for contributing to Nixopus!
