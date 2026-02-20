package expense

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kikeda1102/kakei-board/backend/internal/eventstore"
)

// Projector applies expense events to the read model (expenses table).
type Projector struct {
	db *sql.DB
}

// NewProjector creates a new Projector.
func NewProjector(db *sql.DB) *Projector {
	return &Projector{db: db}
}

// Apply processes an event and updates the read model accordingly.
func (p *Projector) Apply(ctx context.Context, event eventstore.Event) error {
	switch event.EventType {
	case eventTypeRecorded:
		return p.applyRecorded(ctx, event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (p *Projector) applyRecorded(ctx context.Context, event eventstore.Event) error {
	var payload ExpenseRecordedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("unmarshal payload: %w", err)
	}

	_, err := p.db.ExecContext(ctx,
		`INSERT INTO expenses (id, amount, category, memo, date) VALUES (?, ?, ?, ?, ?)`,
		event.AggregateID, payload.Amount, payload.Category, payload.Memo, payload.Date,
	)
	if err != nil {
		return fmt.Errorf("insert expense: %w", err)
	}
	return nil
}
