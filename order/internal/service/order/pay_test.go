package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	orderUUID := uuid.New()
	txUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.EXPECT().PayOrder(s.ctx, orderUUID.String(), model.PaymentMethodCard).Return(txUUID.String(), nil)
	s.repo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(nil)

	txStr, err := s.service.Pay(s.ctx, orderUUID.String(), model.PaymentMethodCard)
	s.Require().NoError(err)
	s.Equal(txUUID.String(), txStr)
}

func (s *ServiceSuite) TestPayInvalidUUID() {
	_, err := s.service.Pay(s.ctx, "not-a-uuid", model.PaymentMethodCard)
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	orderUUID := uuid.New()
	repoErr := errors.New("не найдено")

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(model.Order{}, repoErr)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), model.PaymentMethodCard)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestPayAlreadyCancelled() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusCancelled,
	}

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), model.PaymentMethodCard)
	s.Require().ErrorIs(err, errs.ErrOrderCancelled)
}

func (s *ServiceSuite) TestPayAlreadyPaid() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPaid,
	}

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), model.PaymentMethodCard)
	s.Require().ErrorIs(err, errs.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestPayPaymentClientError() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}
	clientErr := errors.New("payment gateway timeout")

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.EXPECT().PayOrder(s.ctx, orderUUID.String(), model.PaymentMethodCard).Return("", clientErr)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), model.PaymentMethodCard)
	s.Require().ErrorIs(err, clientErr)
}

func (s *ServiceSuite) TestPayRepositoryUpdateError() {
	orderUUID := uuid.New()
	txUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}
	updateErr := errors.New("db error")

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.paymentClient.EXPECT().PayOrder(s.ctx, orderUUID.String(), model.PaymentMethodCard).Return(txUUID.String(), nil)
	s.repo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(updateErr)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), model.PaymentMethodCard)
	s.Require().ErrorIs(err, updateErr)
}
