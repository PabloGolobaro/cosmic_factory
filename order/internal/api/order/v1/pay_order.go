package v1

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	orderv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (orderv1.PayOrderRes, error) {
	txUUID, err := a.OrderService.Pay(ctx, params.OrderUUID.String(), string(req.GetPaymentMethod()))
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound) || errors.Is(err, errs.ErrInvalidUUID):
			return &orderv1.PayOrderNotFound{Code: http.StatusNotFound, Message: err.Error()}, nil
		case errors.Is(err, errs.ErrOrderCancelled) || errors.Is(err, errs.ErrOrderAlreadyPaid):
			return &orderv1.PayOrderConflict{Code: http.StatusConflict, Message: err.Error()}, nil
		case errors.Is(err, errs.ErrInvalidPaymentMethod):
			return &orderv1.PayOrderBadRequest{Code: http.StatusBadRequest, Message: err.Error()}, nil
		default:
			return &orderv1.PayOrderInternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
		}
	}
	return &orderv1.PayOrderResponse{TransactionUUID: txUUID}, nil
}
