package v1

import (
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/static"
	cosmicapi "github.com/PabloGolobaro/cosmic_factory/shared/api"
	orderv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/openapi/order/v1"
)

const (
	middlewareTimeout = 10 * time.Second
	swaggerUIFile     = "files/swagger-ui.html"
)

type api struct {
	OrderService OrderService
	orderv1.UnimplementedHandler
}

func NewApi(orderService OrderService) *api {
	return &api{OrderService: orderService}
}

func (a *api) SetupRouter() (chi.Router, error) {
	orderServer, err := orderv1.NewServer(a)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания сервера OpenAPI: %w", err)
	}

	specFS, err := fs.Sub(cosmicapi.FS, "order/v1")
	if err != nil {
		return nil, fmt.Errorf("инициализация spec FS: %w", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(middlewareTimeout))

	r.Handle("/api/*", orderServer)

	r.Handle("/spec/*", http.StripPrefix("/spec", http.FileServer(http.FS(specFS))))

	r.Get("/swagger-ui.html", func(w http.ResponseWriter, _ *http.Request) {
		data, readErr := static.FS.ReadFile(swaggerUIFile)
		if readErr != nil {
			http.Error(w, "swagger-ui.html не найден", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if _, writeErr := w.Write(data); writeErr != nil {
			slog.Error("ошибка записи swagger-ui", "error", writeErr)
		}
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
	})

	return r, nil
}
