package assembly

type service struct {
	producer Producer
}

func NewService(producer Producer) *service {
	return &service{producer: producer}
}
