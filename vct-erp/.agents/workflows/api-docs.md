---
description: Generate and publish API documentation using Swagger/OpenAPI
---

# /api-docs - API Documentation Generation

// turbo-all

## Prerequisites
- swaggo/swag installed: `go install github.com/swaggo/swag/cmd/swag@latest`

## Steps

### Step 1: Add Swagger annotations to main.go
```go
// @title           VCT Platform API
// @version         1.0
// @description     Vietnam Cycling & Triathlon Management Platform API
// @host            localhost:8080
// @BasePath        /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() { ... }
```

### Step 2: Add annotations to handlers
```go
// @Summary      Create Athlete
// @Description  Create a new athlete registration
// @Tags         athletes
// @Accept       json
// @Produce      json
// @Param        request body CreateAthleteRequest true "Athlete data"
// @Success      201 {object} AthleteResponse
// @Failure      400 {object} ErrorResponse
// @Failure      401 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Security     BearerAuth
// @Router       /athletes [post]
func (h *AthleteHandler) Create(w http.ResponseWriter, r *http.Request) { ... }
```

### Step 3: Generate Swagger docs
```bash
cd backend
swag init -g cmd/server/main.go -o docs/swagger
```

### Step 4: Verify generated files
```bash
ls -la backend/docs/swagger/
# Should contain:
# - docs.go
# - swagger.json
# - swagger.yaml
```

### Step 5: Serve Swagger UI
```go
// In router setup
import httpSwagger "github.com/swaggo/http-swagger"

r.Get("/swagger/*", httpSwagger.Handler(
    httpSwagger.URL("/swagger/doc.json"),
))
```

### Step 6: Access documentation
```bash
# Start server
cd backend && go run cmd/server/main.go

# Open in browser
echo "Swagger UI: http://localhost:8080/swagger/index.html"
```

### Step 7: Export for external use
```bash
# Copy to docs directory
cp backend/docs/swagger/swagger.json docs/api/openapi.json
cp backend/docs/swagger/swagger.yaml docs/api/openapi.yaml
```

## API Documentation Standards
- Every endpoint MUST have Swagger annotations
- Include request/response examples
- Document all error codes
- Security requirements specified
- Vietnamese field descriptions in comments
- Version all API changes
