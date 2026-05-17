package shipassembled

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	eventsv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/events/v1"
)

func decodeShipAssembled(data []byte) (model.ShipAssembledEvent, error) {
	var pb eventsv1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("десериализация ShipAssembled: %w", err)
	}

	return model.ShipAssembledEvent{
		EventUUID:    pb.GetEventUuid(),
		OrderUUID:    pb.GetOrderUuid(),
		UserUUID:     pb.GetUserUuid(),
		BuildTimeSec: pb.GetBuildTimeSec(),
	}, nil
}
