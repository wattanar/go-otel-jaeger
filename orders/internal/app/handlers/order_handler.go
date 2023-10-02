package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"orders_service/internal/app/models"
	"orders_service/internal/app/services"

	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

var (
	tracer trace.Tracer
)

type OrderHandler struct {
	tracer       trace.Tracer
	orderService services.OrderService
}

func NewOrderHandler(tp *sdktrace.TracerProvider, orderService services.OrderService) *OrderHandler {

	tracer = tp.Tracer("order-handler")

	return &OrderHandler{
		tracer:       tracer,
		orderService: orderService,
	}
}

func (o *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, span := o.tracer.Start(ctx, "order handler: get_orders")
	defer span.End()

	w.Header().Set("Content-Type", "application-json")

	n := rand.Intn(2)

	orders, err := o.orderService.GetOrders(ctx)
	if err != nil || n == 1 {
		slog.ErrorContext(ctx, "get orders failed got error", err)
		span.SetStatus(codes.Error, "get orders failed")
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, span := o.tracer.Start(ctx, "order handler: create_order")
	defer span.End()

	w.Header().Set("Content-Type", "application-json")

	order := models.Order{
		OrderID: "f43adfba-8cde-4cf9-9526-e5b470766ec2",
		Product: "Oppo A33",
		Price:   7000,
		Qty:     1,
	}

	err := o.orderService.CreateOrder(ctx, order)
	if err != nil {
		slog.ErrorContext(ctx, "create order failed got error", err)
		span.SetStatus(codes.Error, "create order failed")
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}
