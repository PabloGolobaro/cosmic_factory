package assembly

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/service/assembly/mocks"
)

type AssemblySuite struct {
	suite.Suite

	ctx      context.Context
	producer *mocks.Producer
	service  *service
}

func (s *AssemblySuite) SetupTest() {
	s.ctx = context.Background()
	s.producer = mocks.NewProducer(s.T())
	s.service = NewService(s.producer)
	s.service.buildDelay = func() time.Duration { return 0 }
}

func TestAssemblySuite(t *testing.T) {
	suite.Run(t, new(AssemblySuite))
}
