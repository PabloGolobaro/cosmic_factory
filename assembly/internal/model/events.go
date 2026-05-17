package model

import "time"

type OrderPaidEvent struct {
	EventUUID string
	OrderUUID string
	UserUUID  string
}

type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
	AssembledAt  time.Time
}
