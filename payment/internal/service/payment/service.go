package payment

type paymentService struct{}

func NewPaymentService() *paymentService {
	return &paymentService{}
}
