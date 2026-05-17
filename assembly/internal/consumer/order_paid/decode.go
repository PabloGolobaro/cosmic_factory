package orderpaid

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/PabloGolobaro/cosmic_factory/assembly/internal/model"
	eventsv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/events/v1"
)

func decodeOrderPaid(data []byte) (model.OrderPaidEvent, error) {
	var pb eventsv1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("десериализация OrderPaid: %w", err)
	}

	return model.OrderPaidEvent{
		EventUUID: pb.GetEventUuid(),
		OrderUUID: pb.GetOrderUuid(),
		UserUUID:  pb.GetUserUuid(),
	}, nil
}
