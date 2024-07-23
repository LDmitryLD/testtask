package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/LDmitryLD/testtask/config"
	"github.com/LDmitryLD/testtask/internal/infrastructure/logs"
	"github.com/LDmitryLD/testtask/run"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/zap"
)

func initTracer() func() {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint("jaeger:4317"),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(10*time.Second),
		otlptracegrpc.WithRetry(otlptracegrpc.RetryConfig{
			Enabled:         true,
			InitialInterval: 1 * time.Second,
			MaxInterval:     5 * time.Second,
			MaxElapsedTime:  30 * time.Second,
		}),
	)

	if err != nil {
		log.Fatal("failed to create exporter: ", err.Error())
	}

	res, err := resource.New(
		ctx, resource.WithAttributes(semconv.ServiceNameKey.String("rares_service")),
	)
	if err != nil {
		log.Fatal("failed to create resource: ", err.Error())
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter), trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal("failed to shutdown TracerProvider: ", err.Error())
		}
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed load env", err)
	}

	cleanup := initTracer()
	defer cleanup()

	conf := config.NewAppConf()
	logger := logs.NewLogger(conf, os.Stdout)

	conf.Init(logger)

	tracer := otel.Tracer("rates_service")

	_, span := tracer.Start(context.Background(), "main-function")
	defer span.End()

	app := run.NewApp(conf, logger)

	if err := app.Bootstrap().Run(); err != nil {
		logger.Error("app run error", zap.Error(err))
		os.Exit(2)
	}

	span.AddEvent("rates_service started successfully")
}
