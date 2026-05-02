package record

import "time"

type OrderItemRecord struct {
	UUID      string
	OrderUUID string
	PartUUID  string
	PartType  string
	Price     int64
	CreatedAt time.Time
}
