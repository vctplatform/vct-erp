// Package scoring implements the VCT Scoring Module using Hybrid Architecture.
//
// This module is self-contained: it owns its routes, handlers, domain logic,
// and repository implementations. It communicates with other modules
// exclusively through the shared event bus.
//
// Architecture:
//
//	scoring/
//	├── features/          ← Feature slices (one folder per use-case)
//	│   ├── start_match/   ← Medium: validate + write event
//	│   ├── submit_score/  ← Complex: validate + rules + event + broadcast
//	│   ├── record_penalty/← Medium: validate + event
//	│   ├── end_match/     ← Complex: calculate result + event
//	│   ├── get_state/     ← Simple: query only
//	│   ├── submit_forms/  ← Medium: validate + event
//	│   └── finalize_forms/← Complex: calculation + result
//	├── shared/            ← Module-level shared (models, config, calculator)
//	├── module.go          ← Module struct + constructor (this file)
//	└── routes.go          ← Self-registering HTTP routes
package scoring
