package expense_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kikeda1102/kakei-board/backend/internal/eventstore"
	"github.com/kikeda1102/kakei-board/backend/internal/expense"
	"github.com/kikeda1102/kakei-board/backend/internal/testhelper"
	"github.com/kikeda1102/kakei-board/backend/migrations"
)

func setupHandler(t *testing.T) http.Handler {
	t.Helper()

	db := testhelper.OpenTestDB(t)
	if err := migrations.Run(db); err != nil {
		t.Fatalf("run migrations: %v", err)
	}

	store := eventstore.NewMySQLStore(db)
	projector := expense.NewProjector(db)
	repo := expense.NewRepository(db)
	h := expense.NewHandler(store, projector, repo)

	mux := http.NewServeMux()
	h.Register(mux)
	return mux
}

func TestRecordAndListExpenses(t *testing.T) {
	handler := setupHandler(t)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	// POST: record an expense
	body := `{"amount":1500,"category":"食費","memo":"コンビニ","date":"2026-02-20"}`
	resp, err := http.Post(srv.URL+"/expenses", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST /expenses: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST status = %d, want %d", resp.StatusCode, http.StatusCreated)
	}

	var created struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty ID")
	}

	// GET: list expenses
	resp2, err := http.Get(srv.URL + "/expenses")
	if err != nil {
		t.Fatalf("GET /expenses: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("GET status = %d, want %d", resp2.StatusCode, http.StatusOK)
	}

	var expenses []expense.ExpenseRow
	if err := json.NewDecoder(resp2.Body).Decode(&expenses); err != nil {
		t.Fatalf("decode expenses: %v", err)
	}

	if len(expenses) != 1 {
		t.Fatalf("len(expenses) = %d, want 1", len(expenses))
	}

	got := expenses[0]
	if got.ID != created.ID {
		t.Errorf("ID = %q, want %q", got.ID, created.ID)
	}
	if got.Amount != 1500 {
		t.Errorf("Amount = %d, want 1500", got.Amount)
	}
	if got.Category != "食費" {
		t.Errorf("Category = %q, want %q", got.Category, "食費")
	}
	if got.Memo != "コンビニ" {
		t.Errorf("Memo = %q, want %q", got.Memo, "コンビニ")
	}
	if got.Date != "2026-02-20" {
		t.Errorf("Date = %q, want %q", got.Date, "2026-02-20")
	}
}

func TestRecordExpense_ValidationError(t *testing.T) {
	handler := setupHandler(t)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	tests := []struct {
		name string
		body string
	}{
		{"missing amount", `{"category":"食費","date":"2026-02-20"}`},
		{"zero amount", `{"amount":0,"category":"食費","date":"2026-02-20"}`},
		{"missing category", `{"amount":1000,"date":"2026-02-20"}`},
		{"invalid date", `{"amount":1000,"category":"食費","date":"invalid"}`},
		{"invalid json", `{invalid}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.Post(srv.URL+"/expenses", "application/json", strings.NewReader(tt.body))
			if err != nil {
				t.Fatalf("POST /expenses: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
			}
		})
	}
}

func TestListExpenses_Empty(t *testing.T) {
	handler := setupHandler(t)
	srv := httptest.NewServer(handler)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/expenses")
	if err != nil {
		t.Fatalf("GET /expenses: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var expenses []expense.ExpenseRow
	if err := json.NewDecoder(resp.Body).Decode(&expenses); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(expenses) != 0 {
		t.Errorf("len(expenses) = %d, want 0", len(expenses))
	}
}
