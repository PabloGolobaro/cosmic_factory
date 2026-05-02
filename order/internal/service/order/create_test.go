package order

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/PabloGolobaro/cosmic_factory/order/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
)

func (s *ServiceSuite) TestCreateSuccessAllParts() {
	shieldUUID := uuid.New()
	weaponUUID := uuid.New()

	order := model.Order{
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
		ShieldUUID: &shieldUUID,
		WeaponUUID: &weaponUUID,
	}

	parts := []model.Part{
		{UUID: order.HullUUID, Name: "Корпус", Price: 100, StockQuantity: 5},
		{UUID: order.EngineUUID, Name: "Двигатель", Price: 200, StockQuantity: 3},
		{UUID: shieldUUID, Name: "Щит", Price: 150, StockQuantity: 2},
		{UUID: weaponUUID, Name: "Оружие", Price: 250, StockQuantity: 7},
	}

	s.inventoryClient.EXPECT().ListParts(s.ctx, mock.MatchedBy(func(ids []string) bool {
		return slices.Contains(ids, order.HullUUID.String()) &&
			slices.Contains(ids, order.EngineUUID.String()) &&
			slices.Contains(ids, shieldUUID.String()) &&
			slices.Contains(ids, weaponUUID.String())
	})).Return(parts, nil)

	expectedOrder := order
	expectedOrder.OrderUUID = uuid.New()

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Create(s.ctx, mock.AnythingOfType("model.Order")).Return(expectedOrder, nil)
	s.orderItemRepo.EXPECT().Create(s.ctx, mock.AnythingOfType("model.OrderItem")).Return(model.OrderItem{}, nil).Times(4)

	created, err := s.service.Create(s.ctx, order)
	s.Require().NoError(err)
	s.NotEmpty(created.OrderUUID)
}

func (s *ServiceSuite) TestCreateSuccessOnlyRequired() {
	order := model.Order{
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
	}

	parts := []model.Part{
		{UUID: order.HullUUID, Name: "Корпус", Price: 100, StockQuantity: 5},
		{UUID: order.EngineUUID, Name: "Двигатель", Price: 200, StockQuantity: 3},
	}

	s.inventoryClient.EXPECT().ListParts(s.ctx, mock.MatchedBy(func(ids []string) bool {
		return len(ids) == 2 &&
			slices.Contains(ids, order.HullUUID.String()) &&
			slices.Contains(ids, order.EngineUUID.String())
	})).Return(parts, nil)

	expectedOrder := order
	expectedOrder.OrderUUID = uuid.New()

	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Create(s.ctx, mock.AnythingOfType("model.Order")).Return(expectedOrder, nil)
	s.orderItemRepo.EXPECT().Create(s.ctx, mock.AnythingOfType("model.OrderItem")).Return(model.OrderItem{}, nil).Times(2)

	created, err := s.service.Create(s.ctx, order)
	s.Require().NoError(err)
	s.NotEmpty(created.OrderUUID)
}

func (s *ServiceSuite) TestCreateWeaponMissed() {
	shieldUUID := uuid.New()
	weaponUUID := uuid.New()

	order := model.Order{
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
		ShieldUUID: &shieldUUID,
		WeaponUUID: &weaponUUID,
	}

	parts := []model.Part{
		{UUID: order.HullUUID, Name: "Корпус", Price: 100, StockQuantity: 5},
		{UUID: order.EngineUUID, Name: "Двигатель", Price: 200, StockQuantity: 3},
		{UUID: shieldUUID, Name: "Щит", Price: 150, StockQuantity: 2},
	}

	s.inventoryClient.EXPECT().ListParts(s.ctx, mock.MatchedBy(func(ids []string) bool {
		return slices.Contains(ids, order.HullUUID.String()) &&
			slices.Contains(ids, order.EngineUUID.String()) &&
			slices.Contains(ids, shieldUUID.String()) &&
			slices.Contains(ids, weaponUUID.String())
	})).Return(parts, nil)

	_, err := s.service.Create(s.ctx, order)
	s.Require().ErrorIs(err, errs.ErrPartNotFound)
}

func (s *ServiceSuite) TestCreateOutOfStock() {
	order := model.Order{
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
	}

	parts := []model.Part{
		{UUID: order.HullUUID, Name: "Корпус", Price: 100, StockQuantity: 0},
		{UUID: order.EngineUUID, Name: "Двигатель", Price: 200, StockQuantity: 3},
	}

	s.inventoryClient.EXPECT().ListParts(s.ctx, mock.Anything).Return(parts, nil)

	_, err := s.service.Create(s.ctx, order)
	s.Require().ErrorIs(err, errs.ErrOutOfStock)
}

func (s *ServiceSuite) TestCreateInventoryClientError() {
	order := model.Order{
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
	}

	clientErr := errors.New("inventory unavailable")
	s.inventoryClient.EXPECT().ListParts(s.ctx, mock.Anything).Return(nil, clientErr)

	_, err := s.service.Create(s.ctx, order)
	s.Require().ErrorIs(err, clientErr)
}

func (s *ServiceSuite) TestCreateRepositoryError() {
	order := model.Order{
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
	}

	parts := []model.Part{
		{UUID: order.HullUUID, Name: "Корпус", Price: 100, StockQuantity: 5},
		{UUID: order.EngineUUID, Name: "Двигатель", Price: 200, StockQuantity: 3},
	}

	repoErr := errors.New("db error")
	s.inventoryClient.EXPECT().ListParts(s.ctx, mock.Anything).Return(parts, nil)
	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Create(s.ctx, mock.AnythingOfType("model.Order")).Return(model.Order{}, repoErr)

	_, err := s.service.Create(s.ctx, order)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCreateOrderItemError() {
	order := model.Order{
		HullUUID:   uuid.New(),
		EngineUUID: uuid.New(),
	}

	parts := []model.Part{
		{UUID: order.HullUUID, Name: "Корпус", Price: 100, StockQuantity: 5},
		{UUID: order.EngineUUID, Name: "Двигатель", Price: 200, StockQuantity: 3},
	}

	expectedOrder := order
	expectedOrder.OrderUUID = uuid.New()
	itemErr := errors.New("order item insert failed")

	s.inventoryClient.EXPECT().ListParts(s.ctx, mock.Anything).Return(parts, nil)
	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
	s.repo.EXPECT().Create(s.ctx, mock.AnythingOfType("model.Order")).Return(expectedOrder, nil)
	s.orderItemRepo.EXPECT().Create(s.ctx, mock.AnythingOfType("model.OrderItem")).Return(model.OrderItem{}, itemErr).Once()

	_, err := s.service.Create(s.ctx, order)
	s.Require().ErrorIs(err, itemErr)
}
