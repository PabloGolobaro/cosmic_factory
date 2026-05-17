package assembly

import (
	"math/rand/v2"
	"time"
)

type service struct {
	producer   Producer
	buildDelay func() time.Duration // nil → random [5,15]s; overridable in tests
}

func NewService(producer Producer) *service {
	return &service{producer: producer}
}

func (s *service) getBuildDelay() time.Duration {
	if s.buildDelay != nil {
		return s.buildDelay()
	}

	return time.Duration(5+rand.IntN(11)) * time.Second //nolint:mnd
}
