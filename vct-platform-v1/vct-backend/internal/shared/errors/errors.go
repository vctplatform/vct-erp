// Package errors provides typed, structured errors for the VCT Platform.
// This is the canonical location for error definitions in the Hybrid Architecture.
// All errors follow GR-12, containing a machine-readable Code and a
// human-readable Vietnamese Message.
//
// Migration: This package re-exports everything from internal/apierror
// for backward compatibility. New code should import from here.
package errors

import (
	"vct-platform/backend/internal/apierror"
)

// ── Type Aliases ─────────────────────────────────────────────

// Error represents a structured, machine-readable error.
type Error = apierror.Error

// ── Constructor Aliases ──────────────────────────────────────

// New creates a new structured error.
var New = apierror.New

// Newf creates a formatted structured error.
var Newf = apierror.Newf

// Wrap wraps an existing error with a code and context.
var Wrap = apierror.Wrap

// Is checks whether err matches target using errors.Is.
var Is = apierror.Is

// ── Sentinel Errors — Data Store ─────────────────────────────

var (
	ErrNotFound       = apierror.ErrNotFound
	ErrEntityNotFound = apierror.ErrEntityNotFound
	ErrMissingID      = apierror.ErrMissingID
	ErrInvalidID      = apierror.ErrInvalidID
	ErrDuplicateID    = apierror.ErrDuplicateID
)

// ── Sentinel Errors — Auth ───────────────────────────────────

var (
	ErrUnauthorized = apierror.ErrUnauthorized
	ErrForbidden    = apierror.ErrForbidden
	ErrTokenExpired = apierror.ErrTokenExpired
	ErrTokenInvalid = apierror.ErrTokenInvalid
)

// ── Sentinel Errors — Validation ─────────────────────────────

var (
	ErrValidation  = apierror.ErrValidation
	ErrInvalidInput = apierror.ErrInvalidInput
)

// ── Sentinel Errors — Business Logic ─────────────────────────

var (
	ErrConflict        = apierror.ErrConflict
	ErrStateTransition = apierror.ErrStateTransition
	ErrQuotaExceeded   = apierror.ErrQuotaExceeded
)
