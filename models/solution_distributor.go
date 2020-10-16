package models

// SolutionDistributor распределяет решения участников на ревью после заврешения этапа решения
type SolutionDistributor interface {
	// Распределить все решения, сданные в текущем раунде, на ревью между участниками.
	// Каждое решение будет отправлено нескольким другим участникам (reviewerCount)
	Get(allRoundSolutions []Solution, reviewerCount uint) ReviewDistribution
}
