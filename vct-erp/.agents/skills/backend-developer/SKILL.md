---
name: backend-developer
description: Backend Developer role - Go 1.26 development, REST API implementation, Clean Architecture patterns, Supabase/Neon integration, middleware, testing, and database patterns for VCT Platform.
---

# Backend Developer - VCT Platform

## Role Overview
Implements server-side logic using Go 1.26, following Clean Architecture principles. Responsible for API endpoints, business logic, database interactions via Neon (serverless PostgreSQL) and Supabase, and backend testing.

## Technology Stack
- **Language**: Go 1.26 (range-over-func, enhanced generics, improved stdlib, iterator patterns)
- **HTTP Router**: Chi v5 / stdlib `net/http` (enhanced ServeMux with patterns)
- **Database Driver**: pgx v5 (PostgreSQL 18+)
- **Database Provider**: Neon (serverless PostgreSQL, branching, autoscaling)
- **BaaS**: Supabase (Auth verification, Realtime, Storage)
- **Connection Pool**: pgxpool / Neon serverless driver
- **Logger**: zerolog / slog (Go stdlib structured logging)
- **Validator**: go-playground/validator v10
- **Testing**: Go standard testing + testify
- **API Docs**: swaggo/swag (Swagger / OpenAPI 3.1)
- **Config**: Viper or envconfig

## Core Patterns

### 1. Project Layout
```
cmd/
├── server/
│   └── main.go              # Entry point
internal/
├── modules/
│   └── athlete/
│       ├── domain/
│       │   ├── entity.go     # Athlete struct
│       │   ├── repository.go # Interface
│       │   └── errors.go     # AthleteNotFound, etc.
│       ├── usecase/
│       │   ├── create_athlete.go
│       │   ├── get_athlete.go
│       │   └── list_athletes.go
│       └── adapter/
│           ├── postgres/
│           │   └── repository.go   # Neon PostgreSQL implementation
│           ├── supabase/
│           │   └── auth.go         # Supabase Auth verification
│           └── http/
│               ├── handler.go      # HTTP handlers
│               ├── request.go      # CreateAthleteRequest
│               ├── response.go     # AthleteResponse
│               └── router.go       # Route registration
├── shared/
│   ├── middleware/
│   │   ├── auth.go           # Supabase JWT verification
│   │   ├── cors.go
│   │   ├── logging.go
│   │   └── ratelimit.go
│   ├── errors/
│   ├── pagination/
│   └── config/
├── pkg/                    # Reusable packages
│   ├── httputil/
│   ├── dbutil/
│   │   └── neon.go         # Neon connection helpers
│   └── testutil/
go.mod
go.sum
```

### 2. Entity Pattern (Go 1.26 features)
```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type Athlete struct {
    ID          uuid.UUID  `json:"id"`
    FirstName   string     `json:"first_name" validate:"required,min=2,max=100"`
    LastName    string     `json:"last_name" validate:"required,min=2,max=100"`
    DateOfBirth time.Time  `json:"date_of_birth" validate:"required"`
    Gender      Gender     `json:"gender" validate:"required,oneof=male female"`
    Email       string     `json:"email" validate:"required,email"`
    Phone       string     `json:"phone" validate:"omitempty,e164"`
    ClubID      *uuid.UUID `json:"club_id"`
    Status      Status     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}

type Gender string
const (
    GenderMale   Gender = "male"
    GenderFemale Gender = "female"
)

type Status string
const (
    StatusActive    Status = "active"
    StatusInactive  Status = "inactive"
    StatusSuspended Status = "suspended"
)
```

### 3. Repository Interface Pattern
```go
package domain

import (
    "context"
    "iter"
    "github.com/google/uuid"
)

type AthleteRepository interface {
    Create(ctx context.Context, athlete *Athlete) error
    GetByID(ctx context.Context, id uuid.UUID) (*Athlete, error)
    Update(ctx context.Context, athlete *Athlete) error
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context, filter AthleteFilter) ([]Athlete, int64, error)
    // Go 1.26 iterator pattern for streaming large datasets
    Stream(ctx context.Context, filter AthleteFilter) iter.Seq2[Athlete, error]
}

type AthleteFilter struct {
    ClubID  *uuid.UUID
    Gender  *Gender
    Status  *Status
    Search  string
    Page    int
    PerPage int
    SortBy  string
    SortDir string
}
```

### 4. Use Case Pattern
```go
package usecase

import (
    "context"
    "fmt"
    "module/internal/modules/athlete/domain"
)

type CreateAthleteUseCase struct {
    repo domain.AthleteRepository
}

func NewCreateAthleteUseCase(repo domain.AthleteRepository) *CreateAthleteUseCase {
    return &CreateAthleteUseCase{repo: repo}
}

func (uc *CreateAthleteUseCase) Execute(ctx context.Context, input CreateAthleteInput) (*domain.Athlete, error) {
    athlete := &domain.Athlete{
        FirstName:   input.FirstName,
        LastName:    input.LastName,
        DateOfBirth: input.DateOfBirth,
        Gender:      input.Gender,
        Email:       input.Email,
        Status:      domain.StatusActive,
    }

    if err := uc.repo.Create(ctx, athlete); err != nil {
        return nil, fmt.Errorf("create athlete: %w", err)
    }

    return athlete, nil
}
```

### 5. HTTP Handler Pattern
```go
package http

import (
    "encoding/json"
    "net/http"
    "module/internal/shared/errors"
)

type AthleteHandler struct {
    createUC *usecase.CreateAthleteUseCase
    getUC    *usecase.GetAthleteUseCase
    listUC   *usecase.ListAthletesUseCase
}

func (h *AthleteHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req CreateAthleteRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        errors.WriteJSON(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
        return
    }

    if err := req.Validate(); err != nil {
        errors.WriteValidationError(w, err)
        return
    }

    athlete, err := h.createUC.Execute(r.Context(), req.ToInput())
    if err != nil {
        errors.WriteInternalError(w, err)
        return
    }

    writeJSON(w, http.StatusCreated, NewAthleteResponse(athlete))
}
```

### 6. Neon PostgreSQL Repository Pattern
```go
package postgres

import (
    "context"
    "fmt"
    "iter"

    "github.com/jackc/pgx/v5/pgxpool"
)

// NewNeonPool creates a connection pool optimized for Neon serverless
func NewNeonPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, fmt.Errorf("parse neon config: %w", err)
    }

    // Neon-optimized settings
    config.MaxConns = 10                   // Neon handles scaling
    config.MinConns = 0                    // Allow scale-to-zero
    config.MaxConnLifetime = 30 * time.Minute
    config.MaxConnIdleTime = 5 * time.Minute

    return pgxpool.NewWithConfig(ctx, config)
}

const createAthleteQuery = `
INSERT INTO athletes (first_name, last_name, date_of_birth, gender, email, phone, club_id, status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, created_at, updated_at`

func (r *AthleteRepository) Create(ctx context.Context, a *domain.Athlete) error {
    return r.pool.QueryRow(ctx, createAthleteQuery,
        a.FirstName, a.LastName, a.DateOfBirth, a.Gender,
        a.Email, a.Phone, a.ClubID, a.Status,
    ).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt)
}

// Go 1.26 iterator for streaming results from Neon
func (r *AthleteRepository) Stream(ctx context.Context, filter AthleteFilter) iter.Seq2[domain.Athlete, error] {
    return func(yield func(domain.Athlete, error) bool) {
        rows, err := r.pool.Query(ctx, listAthletesQuery, filter.toArgs()...)
        if err != nil {
            yield(domain.Athlete{}, err)
            return
        }
        defer rows.Close()

        for rows.Next() {
            var a domain.Athlete
            if err := rows.Scan(&a.ID, &a.FirstName, &a.LastName); err != nil {
                if !yield(domain.Athlete{}, err) {
                    return
                }
                continue
            }
            if !yield(a, nil) {
                return
            }
        }
    }
}
```

### 7. Supabase Auth Middleware
```go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

// VerifySupabaseToken validates Supabase JWT tokens
func VerifySupabaseToken(supabaseJWTSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }

            tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
            token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
                return []byte(supabaseJWTSecret), nil
            })
            if err != nil || !token.Valid {
                http.Error(w, "invalid token", http.StatusUnauthorized)
                return
            }

            claims := token.Claims.(jwt.MapClaims)
            ctx := context.WithValue(r.Context(), "user_id", claims["sub"])
            ctx = context.WithValue(ctx, "role", claims["role"])
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

### 8. Testing Standards

#### Unit Test Pattern
```go
func TestCreateAthleteUseCase_Execute(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateAthleteInput
        mock    func(*MockRepository)
        want    *domain.Athlete
        wantErr bool
    }{
        {
            name:  "success",
            input: validInput(),
            mock: func(m *MockRepository) {
                m.On("Create", mock.Anything, mock.Anything).Return(nil)
            },
            wantErr: false,
        },
        {
            name:  "duplicate email",
            input: validInput(),
            mock: func(m *MockRepository) {
                m.On("Create", mock.Anything, mock.Anything).Return(domain.ErrDuplicateEmail)
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockRepository)
            tt.mock(repo)
            uc := NewCreateAthleteUseCase(repo)
            _, err := uc.Execute(context.Background(), tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 9. Error Handling Standards
```go
// Domain errors
var (
    ErrAthleteNotFound = errors.New("athlete not found")
    ErrDuplicateEmail  = errors.New("email already exists")
    ErrInvalidAge      = errors.New("athlete must be at least 6 years old")
    ErrClubNotFound    = errors.New("club not found")
)

// HTTP error mapping
func mapDomainError(err error) (int, string) {
    switch {
    case errors.Is(err, domain.ErrAthleteNotFound):
        return http.StatusNotFound, "ATHLETE_NOT_FOUND"
    case errors.Is(err, domain.ErrDuplicateEmail):
        return http.StatusConflict, "DUPLICATE_EMAIL"
    default:
        return http.StatusInternalServerError, "INTERNAL_ERROR"
    }
}
```

### 10. Development Checklist
- [ ] Follow Clean Architecture boundaries strictly
- [ ] All exported functions have doc comments
- [ ] Tests for use cases (mock repository)
- [ ] Tests for handlers (httptest)
- [ ] Integration tests for repository (Neon branch)
- [ ] `golangci-lint` passes with zero errors
- [ ] `go vet` passes
- [ ] SQL queries use parameterized inputs (no SQL injection)
- [ ] Context propagation in all DB calls
- [ ] Proper error wrapping with `fmt.Errorf("context: %w", err)`
- [ ] Supabase JWT verification on all protected endpoints
- [ ] Neon connection pooling properly configured
