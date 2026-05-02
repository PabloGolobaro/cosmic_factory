package order

type service struct {
	txManager           TxManager
	Repository          OrderRepository
	InventoryClient     InventoryClient
	PaymentClient       PaymentClient
	OrderItemRepository OrderItemRepository
}

func NewService(txManager TxManager, repository OrderRepository, inventoryClient InventoryClient, paymentClient PaymentClient, orderItemRepository OrderItemRepository) *service {
	return &service{txManager: txManager, Repository: repository, InventoryClient: inventoryClient, PaymentClient: paymentClient, OrderItemRepository: orderItemRepository}
}
