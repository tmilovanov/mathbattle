package mock

import (
	"time"

	mathbattle "mathbattle/models"

	"github.com/pkg/errors"
)

type RoundRepository struct {
	impl map[string]mathbattle.Round
}

func NewRoundRepository() RoundRepository {
	return RoundRepository{make(map[string]mathbattle.Round)}
}

func getID(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

func (r *RoundRepository) Store(round mathbattle.Round) error {
	r.impl[getID(round.StartDate)] = round
	return nil
}

func (r *RoundRepository) GetAll() ([]mathbattle.Round, error) {
	return []mathbattle.Round{}, errors.Errorf("Not implemented")
}

func (r *RoundRepository) GetDistributionForRound(roundID string) (mathbattle.RoundDistribution, error) {
	return mathbattle.RoundDistribution{}, nil
}

func (r *RoundRepository) GetRunning() (mathbattle.Round, error) {
	for item := range r.impl {
		emptyTime := time.Time{}
		if r.impl[item].EndDate == emptyTime {
			return r.impl[item], nil
		}
	}

	return mathbattle.Round{}, mathbattle.ErrNotFound
}
