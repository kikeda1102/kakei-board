package expense

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	defaultLimit = 50
	maxLimit     = 200
)

// ExpenseRow represents a row from the expenses read model.
type ExpenseRow struct {
	ID        string    `json:"id"`
	Amount    int64     `json:"amount"`
	Category  string    `json:"category"`
	Memo      string    `json:"memo"`
	Date      string    `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

// Repository reads from the expenses read model.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// List returns expenses ordered by date descending with pagination.
func (r *Repository) List(ctx context.Context, limit, offset int) ([]ExpenseRow, error) {
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, amount, category, memo, date, created_at
		 FROM expenses
		 ORDER BY date DESC, created_at DESC
		 LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("query expenses: %w", err)
	}
	defer rows.Close()

	var expenses []ExpenseRow
	for rows.Next() {
		var e ExpenseRow
		if err := rows.Scan(&e.ID, &e.Amount, &e.Category, &e.Memo, &e.Date, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan expense: %w", err)
		}
		expenses = append(expenses, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate expenses: %w", err)
	}
	return expenses, nil
}
