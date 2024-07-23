package service

import (
	"context"
	"os"
	"testing"

	"github.com/LDmitryLD/testtask/config"
	"github.com/LDmitryLD/testtask/internal/infrastructure/logs"
	"github.com/LDmitryLD/testtask/internal/models"
	"github.com/LDmitryLD/testtask/internal/modules/rates/storage/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRGarantexRates_AddRates(t *testing.T) {
	logger := logs.NewLogger(config.NewAppConf(), os.Stdout)
	mockStorage := mocks.NewRatesStorager(t)
	mockStorage.On("InsertRates", mock.Anything, models.GetRatesResponse{}).Return(nil)

	rates := NewGarantexRates(mockStorage, logger)

	err := rates.AddRates(context.Background(), models.GetRatesResponse{})

	assert.NoError(t, err)
}
