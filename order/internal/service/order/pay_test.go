package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.EXPECT().PayOrder(s.ctx, orderUUID.String(), model.PaymentMethodCard).Return(nil)
	s.repo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(nil)

	txUUID, err := s.service.Pay(s.ctx, orderUUID.String(), "CARD")
	s.Require().NoError(err)
	s.NotEqual(uuid.Nil, txUUID)
}

func (s *ServiceSuite) TestPayInvalidUUID() {
	_, err := s.service.Pay(s.ctx, "not-a-uuid", "CARD")
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestPayInvalidPaymentMethod() {
	orderUUID := uuid.New()

	_, err := s.service.Pay(s.ctx, orderUUID.String(), "BITCOIN")
	s.Require().ErrorIs(err, errs.ErrInvalidPaymentMethod)
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	orderUUID := uuid.New()
	repoErr := errors.New("не найдено")

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(model.Order{}, repoErr)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), "CARD")
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestPayAlreadyCancelled() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusCancelled,
	}

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), "CARD")
	s.Require().ErrorIs(err, errs.ErrOrderCancelled)
}

func (s *ServiceSuite) TestPayAlreadyPaid() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPaid,
	}

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), "CARD")
	s.Require().ErrorIs(err, errs.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestPayPaymentClientError() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}
	clientErr := errors.New("payment gateway timeout")

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.EXPECT().PayOrder(s.ctx, orderUUID.String(), model.PaymentMethodCard).Return(clientErr)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), "CARD")
	s.Require().ErrorIs(err, clientErr)
}

func (s *ServiceSuite) TestPayRepositoryUpdateError() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}
	updateErr := errors.New("db error")

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.EXPECT().PayOrder(s.ctx, orderUUID.String(), model.PaymentMethodCard).Return(nil)
	s.repo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(updateErr)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), "CARD")
	s.Require().ErrorIs(err, updateErr)
}
