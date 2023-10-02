package services

import (
	"context"
	"math/rand"
	"net/http"
	"orders_service/internal/app/models"
	"orders_service/internal/app/repositories"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

var (
	tracer trace.Tracer
)

type OrderService struct {
	tracer          trace.Tracer
	orderRepository repositories.OrderRepository
}

func NewOrderService(tp *sdktrace.TracerProvider, orderRepository repositories.OrderRepository) *OrderService {
	tracer = tp.Tracer("order-service")

	return &OrderService{
		tracer:          tracer,
		orderRepository: orderRepository,
	}
}

func (o *OrderService) GetOrders(ctx context.Context) ([]models.Order, error) {
	ctx, span := o.tracer.Start(ctx, "order service: get orders")
	defer span.End()

	delay := rand.Intn(300)
	time.Sleep(time.Duration(delay) * time.Millisecond)

	orders, err := o.orderRepository.GetOrders(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "get orders failed got error", err)
		span.SetStatus(codes.Error, "get orders failed")
		span.RecordError(err)
		return []models.Order{}, err
	}

	return orders, nil
}

func (o *OrderService) CreateOrder(ctx context.Context, order models.Order) error {
	ctx, span := o.tracer.Start(ctx, "order service: create order")
	defer span.End()

	delay := rand.Intn(300)
	time.Sleep(time.Duration(delay) * time.Millisecond)

	if err := o.orderRepository.CreateOrder(ctx, order); err != nil {
		slog.ErrorContext(ctx, "create order failed got error", err)
		span.SetStatus(codes.Error, "create order failed")
		span.RecordError(err)
		return err
	}

	ctx, spanInventory := o.tracer.Start(ctx, "call inventory service: update inventory")
	defer spanInventory.End()

	client := http.Client{}
	req, _ := http.NewRequest("POST", "http://localhost:8081/inventory/update", nil)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	res, err := client.Do(req)

	if err != nil || res.StatusCode != 200 {
		slog.ErrorContext(ctx, "update inventory failed got error", err)
		spanInventory.SetStatus(codes.Error, "update inventory failed")
		spanInventory.RecordError(err)
		return err
	}

	defer res.Body.Close()

	slog.InfoContext(ctx, "status code:", res.StatusCode)

	return nil
}
