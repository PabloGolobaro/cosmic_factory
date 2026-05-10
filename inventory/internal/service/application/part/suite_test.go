package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/service/application/part/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx       context.Context
	repo      *mocks.PartRepository
	txManager *mocks.TxManager
	svc       *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.repo = mocks.NewPartRepository(s.T())
	s.txManager = mocks.NewTxManager(s.T())
	s.svc = NewPartService(s.repo, nil, s.txManager)
}

func (s *ServiceSuite) TearDownTest() {
	s.T().Log("TearDownTest: очистка после", s.T().Name())
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
