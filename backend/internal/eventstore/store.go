package eventstore

import (
	"context"
	"fmt"
	"time"
)

// Event represents a single domain event persisted in the event store.
type Event struct {
	ID            uint64
	AggregateID   string
	AggregateType string
	Version       int
	EventType     string
	Payload       []byte
	RecordedBy    string
	OccurredAt    time.Time
}

// VersionConflictError indicates an optimistic concurrency violation.
// Another writer appended events to the same aggregate before this write.
type VersionConflictError struct {
	AggregateID   string
	AggregateType string
	Expected      int
}

func (e *VersionConflictError) Error() string {
	return fmt.Sprintf("version conflict on %s/%s: expected version %d",
		e.AggregateType, e.AggregateID, e.Expected)
}

// Store defines the interface for appending and loading domain events.
type Store interface {
	// Append persists events atomically. If the current version of the
	// aggregate does not match expectedVersion, a VersionConflictError is returned.
	Append(ctx context.Context, events []Event, expectedVersion int) error

	// Load returns all events for the given aggregate, ordered by version.
	Load(ctx context.Context, aggregateType, aggregateID string) ([]Event, error)
}
