package order

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.repo.EXPECT().Update(s.ctx, model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusCancelled,
	}).Return(nil)

	err := s.service.Cancel(s.ctx, orderUUID.String())
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelInvalidUUID() {
	err := s.service.Cancel(s.ctx, "not-a-uuid")
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	orderUUID := uuid.New()
	repoErr := errors.New("не найдено")

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(model.Order{}, repoErr)

	err := s.service.Cancel(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, errs.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCancelAlreadyCancelled() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusCancelled,
	}

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)

	err := s.service.Cancel(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, errs.ErrOrderCancelled)
}

func (s *ServiceSuite) TestCancelAlreadyPaid() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPaid,
	}

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)

	err := s.service.Cancel(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, errs.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestCancelRepositoryUpdateError() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPendingPayment,
	}
	updateErr := errors.New("db error")

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.repo.EXPECT().Update(s.ctx, model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusCancelled,
	}).Return(updateErr)

	err := s.service.Cancel(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, updateErr)
}
