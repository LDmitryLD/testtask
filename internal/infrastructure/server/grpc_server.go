package server

import (
	"context"
	"fmt"
	"net"

	"github.com/LDmitryLD/testtask/config"
	"github.com/LDmitryLD/testtask/grpc/rates"

	pb "github.com/LDmitryLD/testtask/grpc/proto/api"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server interface {
	Serve(ctx context.Context) error
}

type GRPCServer struct {
	conf   config.GRPCServer
	srv    *grpc.Server
	rates  *rates.RatesServiceGRPC
	logger *zap.Logger
}

func NewGRPCServer(conf config.GRPCServer, rates *rates.RatesServiceGRPC, logger *zap.Logger) Server {
	return &GRPCServer{
		conf:   conf,
		srv:    grpc.NewServer(),
		rates:  rates,
		logger: logger,
	}
}

func (s *GRPCServer) Serve(ctx context.Context) error {
	var err error

	chErr := make(chan error)
	go func() {
		var l net.Listener
		l, err = net.Listen("tcp", fmt.Sprintf(":%s", s.conf.Port))
		if err != nil {
			s.logger.Error("gRPC server register error:", zap.Error(err))
			chErr <- err
		}

		s.logger.Info("gRPC server started", zap.String("addr", l.Addr().String()))

		pb.RegisterRatesServiceServer(s.srv, s.rates)

		if err = s.srv.Serve(l); err != nil {
			s.logger.Error("grpc server error: ", zap.Error(err))
			chErr <- err
		}
	}()

	select {
	case <-chErr:
		return err
	case <-ctx.Done():
		s.srv.GracefulStop()
	}
	return nil
}
