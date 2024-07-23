package modules

import (
	"github.com/LDmitryLD/testtask/internal/modules/rates/service"
	"github.com/LDmitryLD/testtask/internal/storages"
	"go.uber.org/zap"
)

type Services struct {
	Rates service.RatesServicer
}

func NewServices(storages *storages.Storages, logger *zap.Logger) *Services {
	return &Services{
		Rates: service.NewGarantexRates(storages.Rates, logger),
	}
}
