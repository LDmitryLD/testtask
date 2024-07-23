package service

import (
	"context"

	"github.com/LDmitryLD/testtask/internal/models"
	"github.com/LDmitryLD/testtask/internal/modules/rates/storage"

	"go.uber.org/zap"
)

type RatesServicer interface {
	AddRates(ctx context.Context, rates models.GetRatesResponse) error
}

type GarantexRates struct {
	storage storage.RatesStorager
	logger  *zap.Logger
}

func NewGarantexRates(storage storage.RatesStorager, logger *zap.Logger) RatesServicer {
	return &GarantexRates{
		storage: storage,
		logger:  logger,
	}
}

func (s *GarantexRates) AddRates(ctx context.Context, rates models.GetRatesResponse) error {
	if err := s.storage.InsertRates(ctx, rates); err != nil {
		s.logger.Error("failed to insert rates into database", zap.Error(err))
		return err
	}
	s.logger.Info("rates inserted into database")

	return nil
}
