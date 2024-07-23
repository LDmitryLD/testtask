package adapter

import (
	"context"
	"encoding/json"

	"github.com/LDmitryLD/testtask/internal/models"
	"github.com/jmoiron/sqlx"
)

type SQLAdapterer interface {
	InsertRates(ctx context.Context, rates models.GetRatesResponse) error
}

type SQLAdapter struct {
	db *sqlx.DB
}

func NewSQLAdapter(db *sqlx.DB) *SQLAdapter {
	return &SQLAdapter{
		db: db,
	}
}

func (s *SQLAdapter) InsertRates(ctx context.Context, rates models.GetRatesResponse) error {
	asksJSON, err := json.Marshal(rates.Asks)
	if err != nil {
		return err
	}

	bidsJSON, err := json.Marshal(rates.Bids)
	if err != nil {
		return err
	}

	q := `
	INSERT INTO rates (timestamp, asks, bids)
	VALUES ($1, $2, $3)
	`

	_, err = s.db.ExecContext(ctx, q, rates.Timestamp, asksJSON, bidsJSON)
	if err != nil {
		return err
	}

	return nil
}
