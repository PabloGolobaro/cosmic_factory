package v1

import (
	"context"
	"errors"
	"net/http"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/converter"
	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	orderv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderv1.GetOrderParams) (orderv1.GetOrderRes, error) {
	order, err := a.OrderService.Get(ctx, params.OrderUUID.String())
	if err != nil {
		if errors.Is(err, errs.ErrOrderNotFound) || errors.Is(err, errs.ErrInvalidUUID) {
			return &orderv1.GetOrderNotFound{Code: http.StatusNotFound, Message: err.Error()}, nil
		}
		return &orderv1.GetOrderInternalServerError{Code: http.StatusInternalServerError, Message: err.Error()}, nil
	}
	return converter.OrderToDto(*order), nil
}
