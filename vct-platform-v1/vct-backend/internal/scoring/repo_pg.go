package scoring

// ── PostgreSQL Repository ────────────────────────────────────
// Used in staging/production. Collocated within the scoring module.
// Tables: scoring_match_events, scoring_judge_scores

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"vct-platform/backend/internal/domain/scoring"
)

type pgRepository struct{ db *sql.DB }

func newPgRepository(db *sql.DB) scoring.ScoringRepository {
	return &pgRepository{db: db}
}

// ── Match Events (Event Sourcing — append-only) ──────────────

func (r *pgRepository) AppendMatchEvent(ctx context.Context, event scoring.MatchEvent) error {
	eventDataJSON, err := json.Marshal(event.EventData)
	if err != nil {
		return fmt.Errorf("marshal event_data: %w", err)
	}
	_, err = r.db.ExecContext(ctx,
		`INSERT INTO scoring_match_events
			(id, match_id, match_type, event_type, event_data, sequence_number,
			 round_number, recorded_at, recorded_by, device_id, sync_status)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		event.ID, event.MatchID, event.MatchType, string(event.EventType),
		string(eventDataJSON), event.SequenceNumber, event.RoundNumber,
		event.RecordedAt, event.RecordedBy, event.DeviceID, event.SyncStatus)
	return err
}

func (r *pgRepository) GetMatchEvents(ctx context.Context, matchID string) ([]scoring.MatchEvent, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, match_id, match_type, event_type, event_data, sequence_number,
		        round_number, recorded_at, recorded_by, device_id, sync_status
		 FROM scoring_match_events
		 WHERE match_id=$1
		 ORDER BY sequence_number ASC`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []scoring.MatchEvent
	for rows.Next() {
		var e scoring.MatchEvent
		var eventDataJSON string
		if err := rows.Scan(&e.ID, &e.MatchID, &e.MatchType, &e.EventType,
			&eventDataJSON, &e.SequenceNumber, &e.RoundNumber,
			&e.RecordedAt, &e.RecordedBy, &e.DeviceID, &e.SyncStatus); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(eventDataJSON), &e.EventData)
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *pgRepository) GetNextSequenceNumber(ctx context.Context, matchID string) (int64, error) {
	var maxSeq sql.NullInt64
	err := r.db.QueryRowContext(ctx,
		`SELECT MAX(sequence_number) FROM scoring_match_events WHERE match_id=$1`, matchID).
		Scan(&maxSeq)
	if err != nil {
		return 0, err
	}
	if !maxSeq.Valid {
		return 1, nil
	}
	return maxSeq.Int64 + 1, nil
}

// ── Judge Scores ─────────────────────────────────────────────

func (r *pgRepository) SaveJudgeScore(ctx context.Context, score scoring.JudgeScore) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO scoring_judge_scores
			(id, match_id, referee_id, athlete_id, score, penalties, is_final, submitted_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		 ON CONFLICT (id) DO UPDATE SET
			score=EXCLUDED.score, penalties=EXCLUDED.penalties,
			is_final=EXCLUDED.is_final, submitted_at=EXCLUDED.submitted_at`,
		score.ID, score.MatchID, score.RefereeID, score.AthleteID,
		score.Score, score.Penalties, score.IsFinal, score.SubmittedAt)
	return err
}

func (r *pgRepository) GetJudgeScores(ctx context.Context, matchID string) ([]scoring.JudgeScore, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, match_id, referee_id, athlete_id, score, penalties, is_final, submitted_at
		 FROM scoring_judge_scores
		 WHERE match_id=$1
		 ORDER BY submitted_at`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []scoring.JudgeScore
	for rows.Next() {
		var js scoring.JudgeScore
		if err := rows.Scan(&js.ID, &js.MatchID, &js.RefereeID, &js.AthleteID,
			&js.Score, &js.Penalties, &js.IsFinal, &js.SubmittedAt); err != nil {
			return nil, err
		}
		out = append(out, js)
	}
	return out, rows.Err()
}
