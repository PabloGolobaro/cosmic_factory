package orderproducer

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/platform/pkg/kafka"
	eventsv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/events/v1"
)

func (s *service) PublishOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	pb := &eventsv1.OrderPaid{
		EventUuid: event.EventUUID,
		OrderUuid: event.OrderUUID,
		UserUuid:  event.UserUUID,
	}

	data, err := proto.Marshal(pb)
	if err != nil {
		return fmt.Errorf("сериализация OrderPaid: %w", err)
	}

	return s.orderPaidProducer.Send(ctx, &kafka.Message{
		Key:       []byte(event.OrderUUID),
		Value:     data,
		Timestamp: time.Now(),
	})
}
