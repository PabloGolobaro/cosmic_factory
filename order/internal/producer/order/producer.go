package order

type service struct {
	orderPaidProducer KafkaProducer
}

func NewService(orderPaidProducer KafkaProducer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}
