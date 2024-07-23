package storage

import (
	"context"

	"github.com/LDmitryLD/testtask/internal/db/adapter"
	"github.com/LDmitryLD/testtask/internal/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=RatesStorager
type RatesStorager interface {
	InsertRates(ctx context.Context, rates models.GetRatesResponse) error
}

type RatesStorage struct {
	adapter adapter.SQLAdapterer
}

func NewRatesStorage(sqlAdapter adapter.SQLAdapterer) RatesStorager {
	return &RatesStorage{
		adapter: sqlAdapter,
	}
}

func (s *RatesStorage) InsertRates(ctx context.Context, rates models.GetRatesResponse) error {
	return s.adapter.InsertRates(ctx, rates)
}
