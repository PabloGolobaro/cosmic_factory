package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/service/order/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	txManager         *mocks.TxManager
	inventoryClient   *mocks.InventoryClient
	paymentClient     *mocks.PaymentClient
	repo              *mocks.Repository
	orderItemRepo     *mocks.OrderItemRepository
	orderPaidProducer *mocks.OrderPaidProducer

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.txManager = mocks.NewTxManager(s.T())
	s.repo = mocks.NewRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())
	s.orderItemRepo = mocks.NewOrderItemRepository(s.T())
	s.orderPaidProducer = mocks.NewOrderPaidProducer(s.T())

	s.service = NewService(s.txManager, s.repo, s.inventoryClient, s.paymentClient, s.orderItemRepo, s.orderPaidProducer)
}

func (s *ServiceSuite) TearDownTest() {
	s.T().Log("TearDownTest: очистка после", s.T().Name())
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
