package order

type service struct {
	Repository      Repository
	InventoryClient InventoryClient
	PaymentClient   PaymentClient
}

func NewService(repository Repository, inventoryClient InventoryClient, paymentClient PaymentClient) *service {
	return &service{
		Repository:      repository,
		InventoryClient: inventoryClient,
		PaymentClient:   paymentClient,
	}
}
