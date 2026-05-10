package order

import (
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	orderUUID := uuid.New()
	hullUUID := uuid.New()
	engineUUID := uuid.New()
	order := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		Status:     model.OrderStatusPendingPayment,
	}

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.inventoryClient.EXPECT().ReleaseParts(s.ctx, mock.MatchedBy(func(ids []string) bool {
		return len(ids) == 2 &&
			slices.Contains(ids, hullUUID.String()) &&
			slices.Contains(ids, engineUUID.String())
	})).Return(nil)
	s.repo.EXPECT().Update(s.ctx, model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		Status:     model.OrderStatusCancelled,
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
	hullUUID := uuid.New()
	engineUUID := uuid.New()
	order := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		Status:     model.OrderStatusPendingPayment,
	}
	updateErr := errors.New("db error")

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.inventoryClient.EXPECT().ReleaseParts(s.ctx, mock.Anything).Return(nil)
	s.repo.EXPECT().Update(s.ctx, model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		Status:     model.OrderStatusCancelled,
	}).Return(updateErr)

	err := s.service.Cancel(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, updateErr)
}

func (s *ServiceSuite) TestCancelReleaseError() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
		Status:     model.OrderStatusPendingPayment,
	}
	releaseErr := errors.New("inventory unavailable")

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(order, nil)
	s.inventoryClient.EXPECT().ReleaseParts(s.ctx, mock.Anything).Return(releaseErr)

	err := s.service.Cancel(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, releaseErr)
}
