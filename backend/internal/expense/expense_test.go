package expense

import (
	"encoding/json"
	"testing"
)

func TestRecordExpense_Success(t *testing.T) {
	cmd := RecordExpenseCommand{
		Amount:   1500,
		Category: "食費",
		Memo:     "コンビニ",
		Date:     "2026-02-20",
	}

	event, err := RecordExpense("test-id", cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if event.AggregateID != "test-id" {
		t.Errorf("AggregateID = %q, want %q", event.AggregateID, "test-id")
	}
	if event.AggregateType != "expense" {
		t.Errorf("AggregateType = %q, want %q", event.AggregateType, "expense")
	}
	if event.Version != 1 {
		t.Errorf("Version = %d, want 1", event.Version)
	}
	if event.EventType != "ExpenseRecorded" {
		t.Errorf("EventType = %q, want %q", event.EventType, "ExpenseRecorded")
	}

	var payload ExpenseRecordedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload.Amount != 1500 {
		t.Errorf("payload.Amount = %d, want 1500", payload.Amount)
	}
	if payload.Category != "食費" {
		t.Errorf("payload.Category = %q, want %q", payload.Category, "食費")
	}
	if payload.Memo != "コンビニ" {
		t.Errorf("payload.Memo = %q, want %q", payload.Memo, "コンビニ")
	}
	if payload.Date != "2026-02-20" {
		t.Errorf("payload.Date = %q, want %q", payload.Date, "2026-02-20")
	}
}

func TestRecordExpense_EmptyMemo(t *testing.T) {
	cmd := RecordExpenseCommand{
		Amount:   500,
		Category: "交通費",
		Date:     "2026-01-15",
	}

	event, err := RecordExpense("test-id", cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var payload ExpenseRecordedPayload
	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}
	if payload.Memo != "" {
		t.Errorf("payload.Memo = %q, want empty", payload.Memo)
	}
}

func TestRecordExpense_InvalidAmount(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
	}{
		{"zero", 0},
		{"negative", -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := RecordExpenseCommand{
				Amount:   tt.amount,
				Category: "食費",
				Date:     "2026-02-20",
			}

			_, err := RecordExpense("test-id", cmd)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestRecordExpense_MissingCategory(t *testing.T) {
	cmd := RecordExpenseCommand{
		Amount: 1000,
		Date:   "2026-02-20",
	}

	_, err := RecordExpense("test-id", cmd)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRecordExpense_InvalidDate(t *testing.T) {
	tests := []struct {
		name string
		date string
	}{
		{"empty", ""},
		{"wrong format", "20260220"},
		{"invalid date", "2026-13-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := RecordExpenseCommand{
				Amount:   1000,
				Category: "食費",
				Date:     tt.date,
			}

			_, err := RecordExpense("test-id", cmd)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}

func TestRecordExpense_MultipleErrors(t *testing.T) {
	cmd := RecordExpenseCommand{
		Amount: -1,
		Date:   "invalid",
	}

	_, err := RecordExpense("test-id", cmd)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
