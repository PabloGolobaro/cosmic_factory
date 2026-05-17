package orderproducer

type service struct {
	orderPaidProducer KafkaProducer
}

func New(orderPaidProducer KafkaProducer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}
