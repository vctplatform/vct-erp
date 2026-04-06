// Package events provides the domain event bus for the VCT Platform.
// This is the canonical location for event types and bus in the Hybrid Architecture.
//
// Migration: Re-exports from internal/events for backward compatibility.
// New code should import from shared/events.
package events

import (
	original "vct-platform/backend/internal/events"
)

// ── Type Aliases ─────────────────────────────────────────────

// EventType categorizes domain events.
type EventType = original.EventType

// DomainEvent represents a domain event emitted by services.
type DomainEvent = original.DomainEvent

// Handler is a function that processes domain events.
type Handler = original.Handler

// Bus is the event bus interface.
type Bus = original.Bus

// BusStats holds event bus statistics.
type BusStats = original.BusStats

// InMemoryBus is the in-memory event bus implementation.
type InMemoryBus = original.InMemoryBus

// ── Constructor Aliases ──────────────────────────────────────

// NewBus creates a new in-memory event bus.
var NewBus = original.NewBus

// NewEvent creates a DomainEvent with common fields filled.
var NewEvent = original.NewEvent

// ── Event Type Constants ─────────────────────────────────────
// Approval Events
const (
	EventApprovalSubmitted = original.EventApprovalSubmitted
	EventApprovalApproved  = original.EventApprovalApproved
	EventApprovalRejected  = original.EventApprovalRejected
	EventApprovalReturned  = original.EventApprovalReturned
	EventApprovalCancelled = original.EventApprovalCancelled
)

// Document Events
const (
	EventDocumentDrafted   = original.EventDocumentDrafted
	EventDocumentSubmitted = original.EventDocumentSubmitted
	EventDocumentApproved  = original.EventDocumentApproved
	EventDocumentPublished = original.EventDocumentPublished
	EventDocumentRevoked   = original.EventDocumentRevoked
)

// Certification Events
const (
	EventCertIssued   = original.EventCertIssued
	EventCertRenewed  = original.EventCertRenewed
	EventCertRevoked  = original.EventCertRevoked
	EventCertExpiring = original.EventCertExpiring
)

// Discipline Events
const (
	EventCaseReported     = original.EventCaseReported
	EventCaseInvestigated = original.EventCaseInvestigated
	EventHearingScheduled = original.EventHearingScheduled
	EventCaseDismissed    = original.EventCaseDismissed
)

// International Events
const (
	EventPartnerCreated   = original.EventPartnerCreated
	EventEventCreated     = original.EventEventCreated
	EventDelegationFormed = original.EventDelegationFormed
)

// Federation Events
const (
	EventProvinceCreated   = original.EventProvinceCreated
	EventUnitCreated       = original.EventUnitCreated
	EventPersonnelAssigned = original.EventPersonnelAssigned
)
