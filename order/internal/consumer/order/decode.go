package order

//func decodeUFORecorded(data []byte) (model.UFORecordedEvent, error) {
//	var pb eventsv1.UFORecorded
//	if err := proto.Unmarshal(data, &pb); err != nil {
//		return model.UFORecordedEvent{}, fmt.Errorf("не удалось десериализовать protobuf: %w", err)
//	}
//
//	var observedAt *time.Time
//	if pb.ObservedAt != nil {
//		observedAt = new(pb.ObservedAt.AsTime())
//	}
//
//	return model.UFORecordedEvent{
//		UUID:        pb.Uuid,
//		ObservedAt:  observedAt,
//		Location:    pb.Location,
//		Description: pb.Description,
//	}, nil
//}
