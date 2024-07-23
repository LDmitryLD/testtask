package models

import (
	pb "github.com/LDmitryLD/testtask/grpc/proto/api"
)

func RatesResponseFromDomain(in GetRatesResponse) *pb.GetRatesResponse {
	return &pb.GetRatesResponse{
		Timestamp: in.Timestamp,
		Asks:      OrdersFromDomain(in.Asks),
		Bids:      OrdersFromDomain(in.Bids),
	}
}

func OrdersFromDomain(in []Order) []*pb.Order {
	results := make([]*pb.Order, len(in))
	for i, order := range in {
		results[i] = OrderFromDomain(order)
	}

	return results
}

func OrderFromDomain(in Order) *pb.Order {
	return &pb.Order{
		Price:  in.Price,
		Volume: in.Volume,
		Amount: in.Amount,
		Factor: in.Factor,
		Type:   in.Type,
	}
}
