package order

type service struct {
	ufoRecordedProducer KafkaProducer
}

func NewService(ufoRecordedProducer KafkaProducer) *service {
	return &service{
		ufoRecordedProducer: ufoRecordedProducer,
	}
}
