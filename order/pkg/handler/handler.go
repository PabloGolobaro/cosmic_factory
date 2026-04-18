package handler

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	orderv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

// Order представляет заказ на постройку космического корабля.
type Order struct {
	OrderUUID       uuid.UUID
	HullUUID        uuid.UUID
	EngineUUID      uuid.UUID
	ShieldUUID      *uuid.UUID // опциональный
	WeaponUUID      *uuid.UUID // опциональный
	TotalPrice      int64      // в копейках
	TransactionUUID *uuid.UUID
	PaymentMethod   *string
	Status          string // PENDING_PAYMENT, PAID, CANCELLED
	CreatedAt       time.Time
}

// OrderStore — хранилище заказов (in-memory).
type OrderStore struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]Order
}

// NewOrderStore создаёт новое пустое хранилище заказов.
func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: make(map[uuid.UUID]Order),
	}
}

// OrderHandler реализует интерфейс orderv1.Handler, сгенерированный ogen.
type OrderHandler struct {
	orderv1.UnimplementedHandler
	inventoryClient inventoryv1.InventoryServiceClient
	paymentClient   paymentv1.PaymentServiceClient
	store           *OrderStore
}

// NewOrderHandler создаёт новый обработчик заказов.
func NewOrderHandler(
	inventoryClient inventoryv1.InventoryServiceClient,
	paymentClient paymentv1.PaymentServiceClient,
	store *OrderStore,
) *OrderHandler {
	return &OrderHandler{
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
		store:           store,
	}
}

// SetupServer создаёт OpenAPI сервер на основе обработчика.
func SetupServer(h *OrderHandler) (*orderv1.Server, error) {
	return orderv1.NewServer(h)
}

// GetOrder реализует операцию getOrder (пример реализации).
// GET /api/v1/orders/{order_uuid}.
func (h *OrderHandler) GetOrder(_ context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	// 1. Найти заказ в store (с блокировкой для thread-safety)
	h.store.mu.RLock()
	order, ok := h.store.orders[params.OrderUUID]
	h.store.mu.RUnlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.GetOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	// 3. Преобразовать в DTO и вернуть
	var shieldUUID orderv1.OptNilUUID
	if order.ShieldUUID != nil {
		shieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}

	var weaponUUID orderv1.OptNilUUID
	if order.WeaponUUID != nil {
		weaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}

	var transactionUUID orderv1.OptNilUUID
	if order.TransactionUUID != nil {
		transactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}

	var paymentMethod orderv1.OptNilPaymentMethod
	if order.PaymentMethod != nil {
		paymentMethod = orderv1.NewOptNilPaymentMethod(orderv1.PaymentMethod(*order.PaymentMethod))
	}

	return &orderv1.OrderDto{
		OrderUUID:       order.OrderUUID,
		HullUUID:        order.HullUUID,
		EngineUUID:      order.EngineUUID,
		ShieldUUID:      shieldUUID,
		WeaponUUID:      weaponUUID,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: transactionUUID,
		PaymentMethod:   paymentMethod,
		Status:          orderv1.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
	}, nil
}

// CancelOrder реализует операцию cancelOrder.
// POST /api/v1/orders/{order_uuid}/cancel.
func (h *OrderHandler) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	h.store.mu.Lock()
	order, ok := h.store.orders[params.OrderUUID]
	defer h.store.mu.Unlock()

	// 2. Если не найден — вернуть 404
	if !ok {
		return &orderv1.CancelOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}

	switch order.Status {
	case string(orderv1.OrderStatusCANCELLED):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже отменен",
		}, nil
	case string(orderv1.OrderStatusPAID):
		return &orderv1.CancelOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен",
		}, nil
	default:
	}

	order.Status = string(orderv1.OrderStatusCANCELLED)

	h.store.orders[params.OrderUUID] = order

	slog.Info("Ордер успешно отменен", slog.String("order", order.OrderUUID.String()))

	return &orderv1.CancelOrderResponse{}, nil
}

// CreateOrder реализует операцию createOrder.
// POST /api/v1/orders.
func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	uuids := []string{req.HullUUID.String(), req.EngineUUID.String()}
	var shieldUUID, weaponUUID *uuid.UUID
	if v, ok := req.ShieldUUID.Get(); ok {
		uuids = append(uuids, v.String())
		shieldUUID = &v
	}
	if v, ok := req.WeaponUUID.Get(); ok {
		uuids = append(uuids, v.String())
		weaponUUID = &v
	}

	resp, err := h.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{Uuids: uuids})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				return &orderv1.CreateOrderNotFound{Code: http.StatusNotFound, Message: st.Message()}, nil
			case codes.InvalidArgument:
				return &orderv1.CreateOrderBadRequest{Code: http.StatusBadRequest, Message: st.Message()}, nil
			}
		}
		return &orderv1.CreateOrderInternalServerError{Code: http.StatusInternalServerError, Message: "ошибка инвентарного сервиса"}, nil
	}

	partsMap := make(map[string]*inventoryv1.Part, len(resp.GetParts()))
	for _, p := range resp.GetParts() {
		partsMap[p.GetUuid()] = p
	}

	var totalPrice int64
	for _, id := range uuids {
		p, ok := partsMap[id]
		if !ok {
			return &orderv1.CreateOrderNotFound{Code: http.StatusNotFound, Message: "деталь не найдена: " + id}, nil
		}
		if p.GetStockQuantity() <= 0 {
			return &orderv1.CreateOrderConflict{Code: http.StatusConflict, Message: "деталь отсутствует на складе: " + p.GetName()}, nil
		}
		totalPrice += p.GetPrice()
	}

	orderUUID := uuid.New()
	h.store.mu.Lock()
	h.store.orders[orderUUID] = Order{
		OrderUUID:  orderUUID,
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
		ShieldUUID: shieldUUID,
		WeaponUUID: weaponUUID,
		TotalPrice: totalPrice,
		Status:     string(orderv1.OrderStatusPENDINGPAYMENT),
		CreatedAt:  time.Now(),
	}
	h.store.mu.Unlock()

	slog.Info("заказ создан", slog.String("order_uuid", orderUUID.String()))

	return &orderv1.CreateOrderResponse{OrderUUID: orderUUID, TotalPrice: totalPrice}, nil
}

// PayOrder реализует операцию payOrder.
// POST /api/v1/orders/{order_uuid}/pay.
func (h *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	id := params.OrderUUID

	// 1. Найти заказ в store
	h.store.mu.Lock()
	order, ok := h.store.orders[id]
	defer h.store.mu.Unlock()

	if !ok {
		return &orderv1.PayOrderNotFound{
			Code:    http.StatusNotFound,
			Message: "заказ не найден",
		}, nil
	}
	// 2. Проверить статус == PENDING_PAYMENT
	switch order.Status {
	case string(orderv1.OrderStatusCANCELLED):
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже отменен",
		}, nil
	case string(orderv1.OrderStatusPAID):
		return &orderv1.PayOrderConflict{
			Code:    http.StatusConflict,
			Message: "заказ уже оплачен",
		}, nil
	default:
	}

	paymentMethod := req.GetPaymentMethod()

	// 3. Вызвать h.paymentClient.PayOrder для обработки платежа
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	payOrderResponse, err := h.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{OrderUuid: id.String(), PaymentMethod: PaymentMethodConvert(paymentMethod)})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				return &orderv1.PayOrderBadRequest{
					Code:    http.StatusBadRequest,
					Message: st.Message(),
				}, nil
			case codes.NotFound:
				return &orderv1.PayOrderNotFound{
					Code:    http.StatusNotFound,
					Message: st.Message(),
				}, nil
			}
		}
		return &orderv1.PayOrderInternalServerError{
			Code:    http.StatusInternalServerError,
			Message: "внутренняя ошибка платёжного сервиса",
		}, nil
	}

	// 4. Обновить статус на PAID и сохранить transaction_uuid
	txUUID, err := uuid.Parse(payOrderResponse.GetTransactionUuid())
	if err != nil {
		return nil, err
	}

	order.PaymentMethod = new(string(paymentMethod))
	order.TransactionUUID = new(txUUID)
	order.Status = string(orderv1.OrderStatusPAID)
	h.store.orders[id] = order

	// 5. Вернуть transaction_uuid
	return &orderv1.PayOrderResponse{TransactionUUID: txUUID}, nil
}

func PaymentMethodConvert(method orderv1.PaymentMethod) paymentv1.PaymentMethod {
	switch method {
	case orderv1.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
