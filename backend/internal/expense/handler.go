package expense

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/kikeda1102/kakei-board/backend/internal/eventstore"
)

// Handler handles HTTP requests for the expense domain.
type Handler struct {
	store     eventstore.Store
	projector *Projector
	repo      *Repository
}

// NewHandler creates a new Handler.
func NewHandler(store eventstore.Store, projector *Projector, repo *Repository) *Handler {
	return &Handler{
		store:     store,
		projector: projector,
		repo:      repo,
	}
}

// Register adds expense routes to the given mux.
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /expenses", h.RecordExpense)
	mux.HandleFunc("GET /expenses", h.ListExpenses)
}

type recordExpenseResponse struct {
	ID string `json:"id"`
}

// RecordExpense handles POST /expenses.
func (h *Handler) RecordExpense(w http.ResponseWriter, r *http.Request) {
	var cmd RecordExpenseCommand
	if err := json.NewDecoder(r.Body).Decode(&cmd); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	id := uuid.New().String()
	event, err := RecordExpense(id, cmd)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	ctx := r.Context()
	if err := h.store.Append(ctx, []eventstore.Event{event}, 0); err != nil {
		log.Printf("append event: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	if err := h.projector.Apply(ctx, event); err != nil {
		log.Printf("apply projection: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusCreated, recordExpenseResponse{ID: id})
}

// ListExpenses handles GET /expenses.
func (h *Handler) ListExpenses(w http.ResponseWriter, r *http.Request) {
	limit := queryInt(r.Context(), r, "limit", defaultLimit)
	offset := queryInt(r.Context(), r, "offset", 0)

	expenses, err := h.repo.List(r.Context(), limit, offset)
	if err != nil {
		log.Printf("list expenses: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	// Return empty array instead of null
	if expenses == nil {
		expenses = []ExpenseRow{}
	}

	writeJSON(w, http.StatusOK, expenses)
}

func queryInt(_ context.Context, r *http.Request, key string, defaultVal int) int {
	s := r.URL.Query().Get(key)
	if s == "" {
		return defaultVal
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return v
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("write response: %v", err)
	}
}
