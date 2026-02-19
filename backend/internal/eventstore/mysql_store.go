package eventstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// mysqlDuplicateEntryCode is the MySQL error code for duplicate key violations.
const mysqlDuplicateEntryCode = 1062

// MySQLStore implements Store backed by a MySQL events table.
type MySQLStore struct {
	db *sql.DB
}

// NewMySQLStore creates a new MySQLStore.
func NewMySQLStore(db *sql.DB) *MySQLStore {
	return &MySQLStore{db: db}
}

// Append persists events in a single transaction.
// Returns VersionConflictError when the UNIQUE constraint on
// (aggregate_id, aggregate_type, version) is violated.
func (s *MySQLStore) Append(ctx context.Context, events []Event, expectedVersion int) error {
	if len(events) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO events (aggregate_id, aggregate_type, version, event_type, payload, recorded_by)
		 VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer stmt.Close()

	for _, e := range events {
		_, err := stmt.ExecContext(ctx,
			e.AggregateID, e.AggregateType, e.Version, e.EventType, e.Payload, e.RecordedBy)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == mysqlDuplicateEntryCode {
				return &VersionConflictError{
					AggregateID:   e.AggregateID,
					AggregateType: e.AggregateType,
					Expected:      expectedVersion,
				}
			}
			return fmt.Errorf("insert event: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

// Load returns all events for the given aggregate ordered by version.
func (s *MySQLStore) Load(ctx context.Context, aggregateType, aggregateID string) ([]Event, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, aggregate_id, aggregate_type, version, event_type, payload, recorded_by, occurred_at
		 FROM events
		 WHERE aggregate_type = ? AND aggregate_id = ?
		 ORDER BY version ASC`,
		aggregateType, aggregateID)
	if err != nil {
		return nil, fmt.Errorf("query events: %w", err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.AggregateID, &e.AggregateType, &e.Version,
			&e.EventType, &e.Payload, &e.RecordedBy, &e.OccurredAt); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate events: %w", err)
	}
	return events, nil
}
