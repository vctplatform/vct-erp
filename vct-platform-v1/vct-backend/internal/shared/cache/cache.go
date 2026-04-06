// Package cache provides caching infrastructure for the VCT Platform.
// This is the canonical location for cache implementations in the Hybrid Architecture.
//
// Migration: Re-exports from internal/cache for backward compatibility.
// New code should import from shared/cache.
package cache

import (
	original "vct-platform/backend/internal/cache"
)

// ── Type Aliases ─────────────────────────────────────────────

// RedisCache wraps a Redis client with cache operations.
type RedisCache = original.RedisCache

// ── Constructor Aliases ──────────────────────────────────────

// NewRedisCache creates a new Redis cache connection.
var NewRedisCache = original.NewRedisCache
