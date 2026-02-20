package expense

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/kikeda1102/kakei-board/backend/internal/eventstore"
)

const aggregateType = "expense"
const eventTypeRecorded = "ExpenseRecorded"

// RecordExpenseCommand holds the data needed to record a new expense.
type RecordExpenseCommand struct {
	Amount   int64  `json:"amount"`
	Category string `json:"category"`
	Memo     string `json:"memo"`
	Date     string `json:"date"`
}

// ExpenseRecordedPayload is the event payload stored in the event store.
type ExpenseRecordedPayload struct {
	Amount   int64  `json:"amount"`
	Category string `json:"category"`
	Memo     string `json:"memo"`
	Date     string `json:"date"`
}

// Validate checks that the command fields are valid.
func (c RecordExpenseCommand) Validate() error {
	var errs []error

	if c.Amount <= 0 {
		errs = append(errs, fmt.Errorf("amount must be positive"))
	}
	if c.Category == "" {
		errs = append(errs, fmt.Errorf("category is required"))
	}
	if _, err := time.Parse(time.DateOnly, c.Date); err != nil {
		errs = append(errs, fmt.Errorf("date must be in YYYY-MM-DD format"))
	}

	return errors.Join(errs...)
}

// RecordExpense creates an event for recording a new expense.
// This is a pure function that performs no I/O.
func RecordExpense(id string, cmd RecordExpenseCommand) (eventstore.Event, error) {
	if err := cmd.Validate(); err != nil {
		return eventstore.Event{}, err
	}

	payload, err := json.Marshal(ExpenseRecordedPayload{
		Amount:   cmd.Amount,
		Category: cmd.Category,
		Memo:     cmd.Memo,
		Date:     cmd.Date,
	})
	if err != nil {
		return eventstore.Event{}, fmt.Errorf("marshal payload: %w", err)
	}

	return eventstore.Event{
		AggregateID:   id,
		AggregateType: aggregateType,
		Version:       1,
		EventType:     eventTypeRecorded,
		Payload:       payload,
	}, nil
}
