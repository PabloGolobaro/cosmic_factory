package order

import (
	"errors"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *ServiceSuite) TestCommitShipPartsSuccess() {
	orderUUID := uuid.New()
	hullUUID := uuid.New()
	engineUUID := uuid.New()
	order := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		Status:     model.OrderStatusPaid,
	}

	txPassThrough(s)
	s.repo.EXPECT().GetForUpdate(s.ctx, orderUUID).Return(order, nil)
	s.inventoryClient.EXPECT().CommitParts(s.ctx, mock.MatchedBy(func(ids []string) bool {
		return len(ids) == 2
	})).Return(nil)
	s.repo.EXPECT().Update(s.ctx, model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		Status:     model.OrderStatusAssembled,
	}).Return(nil)

	err := s.service.CommitShipParts(s.ctx, orderUUID.String())
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCommitShipPartsInvalidUUID() {
	err := s.service.CommitShipParts(s.ctx, "not-a-uuid")
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestCommitShipPartsOrderNotFound() {
	orderUUID := uuid.New()
	repoErr := errors.New("не найдено")

	txPassThrough(s)
	s.repo.EXPECT().GetForUpdate(s.ctx, orderUUID).Return(model.Order{}, repoErr)

	err := s.service.CommitShipParts(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, errs.ErrOrderNotFound)
}

func (s *ServiceSuite) TestCommitShipPartsAlreadyAssembled() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusAssembled,
	}

	txPassThrough(s)
	s.repo.EXPECT().GetForUpdate(s.ctx, orderUUID).Return(order, nil)

	err := s.service.CommitShipParts(s.ctx, orderUUID.String())
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCommitShipPartsCommitError() {
	orderUUID := uuid.New()
	order := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
		Status:     model.OrderStatusPaid,
	}
	commitErr := errors.New("inventory error")

	txPassThrough(s)
	s.repo.EXPECT().GetForUpdate(s.ctx, orderUUID).Return(order, nil)
	s.inventoryClient.EXPECT().CommitParts(s.ctx, mock.Anything).Return(commitErr)

	err := s.service.CommitShipParts(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, commitErr)
}

func (s *ServiceSuite) TestCommitShipPartsUpdateError() {
	orderUUID := uuid.New()
	hullUUID := uuid.New()
	engineUUID := uuid.New()
	order := model.Order{
		OrderUUID:  orderUUID,
		HullUUID:   hullUUID,
		EngineUUID: engineUUID,
		Status:     model.OrderStatusPaid,
	}
	updateErr := errors.New("db error")

	txPassThrough(s)
	s.repo.EXPECT().GetForUpdate(s.ctx, orderUUID).Return(order, nil)
	s.inventoryClient.EXPECT().CommitParts(s.ctx, mock.Anything).Return(nil)
	s.repo.EXPECT().Update(s.ctx, mock.AnythingOfType("model.Order")).Return(updateErr)

	err := s.service.CommitShipParts(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, updateErr)
}
