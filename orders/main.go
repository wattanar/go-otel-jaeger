package main

import (
	"context"
	"log"
	"net/http"
	"orders_service/internal/app/handlers"
	"orders_service/internal/app/repositories"
	"orders_service/internal/app/services"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"golang.org/x/exp/slog"
)

var (
	JAEGER_ENDPOINT string = "localhost:4318"
	SERVICE_NAME    string = "order-service"
)

func newExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	// Your preferred exporter: console, jaeger, zipkin, OTLP, etc.
	tp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint(JAEGER_ENDPOINT))

	if err != nil {
		slog.ErrorCtx(ctx, "setup exporter failed got error %v", err)
		return nil, err
	}

	return tp, nil
}

func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(SERVICE_NAME),
		),
	)

	if err != nil {
		panic(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func main() {
	ctx := context.Background()

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	// Create a new tracer provider with a batch span processor and the given exporter.
	tp := newTraceProvider(exp)

	// Handle shutdown properly so nothing leaks.
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)

	// create handlers
	orderHanler := createOrderHandler(tp)

	// routes
	http.HandleFunc("/orders", orderHanler.GetOrders)
	http.HandleFunc("/orders/create", orderHanler.CreateOrder)

	http.ListenAndServe(":8080", nil)
}

func createOrderHandler(tp *sdktrace.TracerProvider) *handlers.OrderHandler {

	orderRepository := repositories.NewOrderRepository(tp)
	orderService := services.NewOrderService(tp, *orderRepository)
	orderHandler := handlers.NewOrderHandler(tp, *orderService)
	return orderHandler
}
