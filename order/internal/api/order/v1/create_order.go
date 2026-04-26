package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/converter"
	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	orderv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (orderv1.CreateOrderRes, error) {
	created, err := a.OrderService.Create(ctx, converter.OrderFromCreateRequest(req))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrInvalidUUID):
			return &orderv1.CreateOrderBadRequest{Code: http.StatusBadRequest, Message: err.Error()}, nil
		case errors.Is(err, errs.ErrPartNotFound):
			return &orderv1.CreateOrderNotFound{Code: http.StatusNotFound, Message: err.Error()}, nil
		case errors.Is(err, errs.ErrOutOfStock):
			return &orderv1.CreateOrderConflict{Code: http.StatusConflict, Message: err.Error()}, nil
		default:
			return &orderv1.CreateOrderInternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		}
	}

	return &orderv1.CreateOrderResponse{
		OrderUUID:  created.OrderUUID,
		TotalPrice: created.TotalPrice,
	}, nil
}
