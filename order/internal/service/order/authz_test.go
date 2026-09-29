package order

import (
	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/authz"
)

func (s *ServiceSuite) TestCreateForbidden() {
	s.expectAuthorize(authz.ActionOrderCreate, uuid.Nil, errs.ErrForbidden)

	_, err := s.service.Create(s.ctx, model.Order{HullUUID: uuid.New(), EngineUUID: uuid.New()})
	s.Require().ErrorIs(err, errs.ErrForbidden)
	// Strict mocks guarantee Inventory was not touched.
}

func (s *ServiceSuite) TestGetChecksOwner() {
	orderUUID := uuid.New()
	owner := uuid.New()

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(model.Order{OrderUUID: orderUUID, UserUUID: owner}, nil)
	s.expectAuthorize(authz.ActionOrderRead, owner, nil)
	s.orderItemRepo.EXPECT().ListByOrder(s.ctx, orderUUID).Return(nil, nil)

	got, err := s.service.Get(s.ctx, orderUUID.String())
	s.Require().NoError(err)
	s.Equal(owner, got.UserUUID)
}

func (s *ServiceSuite) TestGetForbidden() {
	orderUUID := uuid.New()
	owner := uuid.New()

	s.repo.EXPECT().Get(s.ctx, orderUUID).Return(model.Order{OrderUUID: orderUUID, UserUUID: owner}, nil)
	s.expectAuthorize(authz.ActionOrderRead, owner, errs.ErrForbidden)

	_, err := s.service.Get(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, errs.ErrForbidden)
}

func (s *ServiceSuite) TestPayForbidden() {
	orderUUID := uuid.New()
	owner := uuid.New()

	txPassThrough(s)
	s.repo.EXPECT().GetForUpdate(s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		UserUUID:  owner,
		Status:    model.OrderStatusPendingPayment,
	}, nil)
	s.expectAuthorize(authz.ActionOrderPay, owner, errs.ErrForbidden)

	_, err := s.service.Pay(s.ctx, orderUUID.String(), model.PaymentMethodCard)
	s.Require().ErrorIs(err, errs.ErrForbidden)
	// Payment must not be called: paymentClient has no expectations.
}

func (s *ServiceSuite) TestCancelForbidden() {
	orderUUID := uuid.New()
	owner := uuid.New()

	txPassThrough(s)
	s.repo.EXPECT().GetForUpdate(s.ctx, orderUUID).Return(model.Order{
		OrderUUID: orderUUID,
		UserUUID:  owner,
		Status:    model.OrderStatusPendingPayment,
	}, nil)
	s.expectAuthorize(authz.ActionOrderCancel, owner, errs.ErrForbidden)

	err := s.service.Cancel(s.ctx, orderUUID.String())
	s.Require().ErrorIs(err, errs.ErrForbidden)
	// Inventory must not release parts: inventoryClient has no expectations.
}
