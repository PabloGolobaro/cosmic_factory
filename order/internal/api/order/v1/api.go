package v1

import (
	"fmt"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	orderv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/openapi/order/v1"
)

const (
	middlewareTimeout = 10 * time.Second
)

type api struct {
	OrderService OrderService
	orderv1.UnimplementedHandler
}

func NewApi(orderService OrderService) *api {
	return &api{OrderService: orderService}
}

func (a *api) SetupRouter() (chi.Router, error) {
	// Создать OpenAPI сервер
	orderServer, err := orderv1.NewServer(a)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(middlewareTimeout))

	r.Handle("/api/*", orderServer)

	return r, nil
}
