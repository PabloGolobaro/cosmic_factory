package part

type service struct {
	PartRepository       PartRepository
	compatibilityChecker CompatibilityChecker
	txManager            TxManager
}

func NewPartService(repo PartRepository, checker CompatibilityChecker, txManager TxManager) *service {
	return &service{PartRepository: repo, compatibilityChecker: checker, txManager: txManager}
}
