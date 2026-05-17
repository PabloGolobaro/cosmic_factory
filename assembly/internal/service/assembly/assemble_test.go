package assembly

import (
	"context"
	"errors"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/model"
)

func (s *AssemblySuite) TestAssemble_Success() {
	event := model.OrderPaidEvent{
		EventUUID: "event-uuid",
		OrderUUID: "order-uuid",
		UserUUID:  "user-uuid",
	}

	s.producer.EXPECT().
		PublishShipAssembled(s.ctx, mock.MatchedBy(func(e model.ShipAssembledEvent) bool {
			return e.OrderUUID == event.OrderUUID &&
				e.UserUUID == event.UserUUID &&
				e.BuildTimeSec == 0 &&
				e.EventUUID != ""
		})).
		Return(nil)

	err := s.service.Assemble(s.ctx, event)
	s.Require().NoError(err)
}

func (s *AssemblySuite) TestAssemble_ProducerError() {
	event := model.OrderPaidEvent{
		EventUUID: "event-uuid",
		OrderUUID: "order-uuid",
		UserUUID:  "user-uuid",
	}
	producerErr := errors.New("kafka недоступна")

	s.producer.EXPECT().
		PublishShipAssembled(s.ctx, mock.MatchedBy(func(e model.ShipAssembledEvent) bool {
			return e.OrderUUID == event.OrderUUID && e.UserUUID == event.UserUUID
		})).
		Return(producerErr)

	err := s.service.Assemble(s.ctx, event)
	s.Require().ErrorIs(err, producerErr)
}

func (s *AssemblySuite) TestAssemble_ContextCancelled() {
	// Большая задержка гарантирует, что select выберет ctx.Done(), а не time.After.
	s.service.buildDelay = func() time.Duration { return time.Hour }

	ctx, cancel := context.WithCancel(s.ctx)
	cancel()

	event := model.OrderPaidEvent{
		EventUUID: "event-uuid",
		OrderUUID: "order-uuid",
		UserUUID:  "user-uuid",
	}

	err := s.service.Assemble(ctx, event)
	s.Require().ErrorIs(err, context.Canceled)
}
