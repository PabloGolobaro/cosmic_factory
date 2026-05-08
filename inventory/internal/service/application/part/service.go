package part

type service struct {
	PartRepository       PartRepository
	compatibilityChecker CompatibilityChecker
}

func NewPartService(repo PartRepository, checker CompatibilityChecker) *service {
	return &service{PartRepository: repo, compatibilityChecker: checker}
}
