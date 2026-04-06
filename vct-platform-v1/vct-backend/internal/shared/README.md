# Shared Kernel — VCT Platform

Global infrastructure code shared across all domain modules.
This is the foundation for the Hybrid Architecture migration.

## Packages

| Package | Purpose | Migrated From |
|---------|---------|---------------|
| `errors/` | Structured API errors | `internal/apierror/` |
| `events/` | Domain event bus | `internal/events/` |
| `httputil/` | HTTP helpers, envelope, response | `internal/httpapi/helpers.go` |
| `database/` | DB connection pool | `internal/store/` (pool only) |
| `cache/` | Redis cache | `internal/cache/` |
| `observability/` | Logging, metrics, tracing | `internal/logging/` + `metrics/` + `tracing/` |

## Migration Strategy

Each package uses **re-export aliases** during transition:
- New code imports from `internal/shared/X`
- Old code continues importing from `internal/X` (alias re-exports)
- After full migration, old packages are deprecated
