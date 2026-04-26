package order

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	orderUUID := uuid.New()
	expected := model.Order{
		OrderUUID: orderUUID,
		HullUUID:  uuid.New(),
		Status:    model.OrderStatusPendingPayment,
	}

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(expected, nil)

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
