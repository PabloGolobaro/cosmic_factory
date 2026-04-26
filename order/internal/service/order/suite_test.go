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

	inventoryClient *mocks.InventoryClient
	paymentClient   *mocks.PaymentClient
	repo            *mocks.Repository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.repo = mocks.NewRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())

	s.service = NewService(s.repo, s.inventoryClient, s.paymentClient)
}

func (s *ServiceSuite) TearDownTest() {
	s.T().Log("TearDownTest: очистка после", s.T().Name())
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
