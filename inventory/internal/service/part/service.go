package part

type service struct {
	PartRepository PartRepository
}

func NewPartService(partRepository PartRepository) *service {
	return &service{PartRepository: partRepository}
}
