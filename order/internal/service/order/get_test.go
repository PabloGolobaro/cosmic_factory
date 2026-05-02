package order

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	orderUUID := uuid.New()
	hullUUID := uuid.New()
	engineUUID := uuid.New()

	items := []model.OrderItem{
		{UUID: uuid.New(), OrderUUID: orderUUID, PartUUID: hullUUID, PartType: model.PartTypeHull, Price: 500},
		{UUID: uuid.New(), OrderUUID: orderUUID, PartUUID: engineUUID, PartType: model.PartTypeEngine, Price: 300},
	}

	repoOrder := model.Order{
		OrderUUID:  orderUUID,
		Status:     model.OrderStatusPendingPayment,
		TotalPrice: 800,
	}
	expected := repoOrder
	expected.HullUUID = hullUUID
	expected.EngineUUID = engineUUID

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(repoOrder, nil)
	s.orderItemRepo.EXPECT().ListByOrder(s.ctx, orderUUID).Return(items, nil)

	got, err := s.service.Get(s.ctx, orderUUID.String())
	s.Require().NoError(err)
	s.Equal(expected, *got)
}

func (s *ServiceSuite) TestGetInvalidUUID() {
	_, err := s.service.Get(s.ctx, "not-a-uuid")
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestGetNotFound() {
	orderUUID := uuid.New()
	repoErr := errors.New("запись не найдена")

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(model.Order{}, repoErr)

	_, err := s.service.Get(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, repoErr)
}
