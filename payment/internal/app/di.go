package app

import (
	apiPayment "github.com/PabloGolobaro/cosmic_factory/payment/internal/api/payment/v1"
	"github.com/PabloGolobaro/cosmic_factory/payment/internal/config"
	paymentservice "github.com/PabloGolobaro/cosmic_factory/payment/internal/service/payment"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

// diContainer — контейнер зависимостей (Composition Root) приложения.
//
// Payment — stateless сервис: нет базы данных, нет исходящих соединений.
// Все геттеры инфаллибельны и возвращают единственное значение без ошибки.
type diContainer struct {
	conf config.Config

	// Сервисный слой (интерфейс из api/payment/v1/deps.go)
	paymentSvc apiPayment.PaymentService

	// gRPC handler
	paymentHandler paymentv1.PaymentServiceServer
}

func newDIContainer(conf config.Config) *diContainer {
	return &diContainer{conf: conf}
}

// PaymentSvc возвращает сервис бизнес-логики платежей.
func (d *diContainer) PaymentSvc() apiPayment.PaymentService {
	if d.paymentSvc == nil {
		d.paymentSvc = paymentservice.NewPaymentService()
	}

	return d.paymentSvc
}

// PaymentHandler возвращает gRPC-обработчик сервиса платежей.
func (d *diContainer) PaymentHandler() paymentv1.PaymentServiceServer {
	if d.paymentHandler == nil {
		d.paymentHandler = apiPayment.NewApi(d.PaymentSvc())
	}

	return d.paymentHandler
}
