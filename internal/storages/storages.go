package storages

import (
	"github.com/LDmitryLD/testtask/internal/db/adapter"
	"github.com/LDmitryLD/testtask/internal/modules/rates/storage"
)

type Storages struct {
	Rates storage.RatesStorager
}

func NewStorages(sqlAdapter *adapter.SQLAdapter) *Storages {
	return &Storages{
		Rates: storage.NewRatesStorage(sqlAdapter),
	}
}
