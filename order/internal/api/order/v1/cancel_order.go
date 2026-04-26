package v1

import (
	"context"
	"errors"
	"net/http"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	orderv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (orderv1.CancelOrderRes, error) {
	if err := a.OrderService.Cancel(ctx, params.OrderUUID.String()); err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderNotFound):
			return &orderv1.CancelOrderNotFound{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}, nil
		case errors.Is(err, errs.ErrOrderCancelled) || errors.Is(err, errs.ErrOrderAlreadyPaid):
			return &orderv1.CancelOrderConflict{
				Code:    http.StatusConflict,
				Message: err.Error(),
			}, nil
		default:
			return &orderv1.CancelOrderInternalServerError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}, nil
		}
	}

	return &orderv1.CancelOrderResponse{}, nil
}
