package repositorytest

import (
	"testing"

	"mathbattle/infrastructure"
	"mathbattle/mocks"
	"mathbattle/models/mathbattle"

	"github.com/stretchr/testify/suite"
)

type participantTs struct {
	suite.Suite

	rep     mathbattle.ParticipantRepository
	userRep mathbattle.UserRepository
}

func (s *participantTs) SetupTest() {
	container := infrastructure.NewTestContainer()
	s.rep = container.ParticipantRepository()
	s.userRep = container.UserRepository()
}

func (s *participantTs) TestSetGetUpdateDelete() {
	participants := mocks.GenParticipants(10, 11)

	for _, participant := range participants {
		u, err := s.userRep.Store(participant.User)
		s.Require().Nil(err)
		participant.User.ID = u.ID
		s.Require().Equal(participant.User, u)

		p, err := s.rep.Store(participant)
		s.Require().Nil(err)

		participant.ID = p.ID
		s.Require().Equal(participant, p)

		p, err = s.rep.GetByID(participant.ID)
		s.Require().Nil(err)
		s.Require().Equal(participant, p)
	}
}

func TestParticipantRepository(t *testing.T) {
	suite.Run(t, &participantTs{})
}
