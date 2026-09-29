package order

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/service/order/mocks"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/auth"
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
	iamClient         *mocks.IAMClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = auth.WithUserUUID(context.Background(), uuid.New())

	s.txManager = mocks.NewTxManager(s.T())
	s.repo = mocks.NewRepository(s.T())
	s.inventoryClient = mocks.NewInventoryClient(s.T())
	s.paymentClient = mocks.NewPaymentClient(s.T())
	s.orderItemRepo = mocks.NewOrderItemRepository(s.T())
	s.orderPaidProducer = mocks.NewOrderPaidProducer(s.T())
	s.iamClient = mocks.NewIAMClient(s.T())

	// By default the policy allows everything; authorization-specific tests
	// override it with expectAuthorize.
	s.iamClient.EXPECT().Authorize(mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	var err error
	s.service, err = NewService(s.txManager, s.repo, s.inventoryClient, s.paymentClient, s.orderItemRepo, s.orderPaidProducer, s.iamClient)
	s.Require().NoError(err)
}

// expectAuthorize replaces the default allow-all expectation with a strict one.
func (s *ServiceSuite) expectAuthorize(action string, owner uuid.UUID, err error) {
	s.iamClient.ExpectedCalls = nil
	s.iamClient.EXPECT().Authorize(mock.Anything, action, owner).Return(err).Once()
}

func (s *ServiceSuite) TearDownTest() {
	s.T().Log("TearDownTest: очистка после", s.T().Name())
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
