package rates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	pb "github.com/LDmitryLD/testtask/grpc/proto/api"
	"github.com/LDmitryLD/testtask/internal/models"
	"github.com/LDmitryLD/testtask/internal/modules/rates/service"
	"go.opentelemetry.io/otel"
)

const garantexURL = "https://garantex.org/api/v2/depth"

type RatesServiceGRPC struct {
	rateService service.RatesServicer
	pb.UnimplementedRatesServiceServer
}

func NewRatesService(ratesService service.RatesServicer) *RatesServiceGRPC {
	return &RatesServiceGRPC{
		rateService: ratesService,
	}
}

func (s *RatesServiceGRPC) GetRates(ctx context.Context, in *pb.GetRatesRequest) (*pb.GetRatesResponse, error) {
	tracer := otel.Tracer("rates_service")
	ctx, span := tracer.Start(ctx, "GetRates")
	defer span.End()

	rates, err := getRates(in.Market)
	if err != nil {
		return nil, err
	}

	err = s.rateService.AddRates(ctx, rates)
	if err != nil {
		return nil, err
	}

	return models.RatesResponseFromDomain(rates), nil
}

func getRates(market string) (models.GetRatesResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s?market=%s", garantexURL, market))
	if err != nil {
		return models.GetRatesResponse{}, err
	}
	defer resp.Body.Close()

	var getRatesResponse models.GetRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&getRatesResponse); err != nil {
		return models.GetRatesResponse{}, err
	}

	return getRatesResponse, nil
}
