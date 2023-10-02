package repositories

import (
	"context"
	"math/rand"
	"orders_service/internal/app/models"
	"time"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

var (
	tracer trace.Tracer
)

type OrderRepository struct {
	tracer trace.Tracer
}

func NewOrderRepository(tp *sdktrace.TracerProvider) *OrderRepository {
	tracer = tp.Tracer("order-repository")
	return &OrderRepository{
		tracer: tracer,
	}
}

func (o *OrderRepository) GetOrders(ctx context.Context) ([]models.Order, error) {
	_, span := o.tracer.Start(ctx, "order repository: get orders")
	defer span.End()

	delay := rand.Intn(500)
	time.Sleep(time.Duration(delay) * time.Millisecond)

	orders := models.MockOrders()
	return orders, nil
}

func (o *OrderRepository) CreateOrder(ctx context.Context, order models.Order) error {
	_, span := o.tracer.Start(ctx, "order repository: create order")
	defer span.End()

	slog.DebugContext(ctx, "create order", order)

	delay := rand.Intn(500)
	time.Sleep(time.Duration(delay) * time.Millisecond)
	return nil
}
