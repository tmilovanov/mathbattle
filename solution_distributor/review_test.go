package solutiondistributor

import (
	"fmt"
	"sort"
	"strconv"
	"testing"

	"mathbattle/mocks"
	mathbattle "mathbattle/models"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Sort1(d map[string][]string) {
	for _, participantIDs := range d {
		sort.Strings(participantIDs)
	}
}

func Sort2(d []mathbattle.Solution) {
	sort.Slice(d, func(i, j int) bool {
		return d[i].ID < d[j].ID
	})
}

func SortAll(d mathbattle.ReviewDistribution) {
	Sort1(d.BetweenParticipants)
	Sort2(d.ToOrganizers)
}

func IsEachParticipantGotKSolutions(d map[string][]string, k uint) bool {
	participantsToSolutions := make(map[string][]string)
	for solutionID, participantIDs := range d {
		for _, pID := range participantIDs {
			participantsToSolutions[pID] = append(participantsToSolutions[pID], solutionID)
		}
	}

	for _, solutionIDs := range participantsToSolutions {
		if uint(len(solutionIDs)) != k {
			return false
		}
	}

	return true
}

func IsEachSolutionGoesToKParticiapnts(d map[string][]string, k uint) bool {
	for _, v := range d {
		if uint(len(v)) != k {
			return false
		}
	}
	return true
}

func helperTestExpect(req *require.Assertions, distributor SolutionDistributor, count uint,
	allRoundSolutions []mathbattle.Solution, expected mathbattle.ReviewDistribution) {

	r := distributor.Get(allRoundSolutions, count)

	SortAll(expected)
	SortAll(r)

	req.Equal(expected.BetweenParticipants, r.BetweenParticipants)
	req.Equal(len(expected.ToOrganizers), len(r.ToOrganizers))
	req.Equal(expected.ToOrganizers, r.ToOrganizers)
}

// One problem for all participants
type oneProblemToAll struct {
	suite.Suite

	distributor SolutionDistributor
	k           uint
}

func (s *oneProblemToAll) TestNoSolutions() {
	helperTestExpect(s.Require(), s.distributor, s.k,
		[]mathbattle.Solution{},
		mathbattle.ReviewDistribution{
			BetweenParticipants: make(map[string][]string),
			ToOrganizers:        []mathbattle.Solution{},
		})
}

func (s *oneProblemToAll) TestOneSolution() {
	helperTestExpect(s.Require(), s.distributor, s.k,
		[]mathbattle.Solution{
			{ID: "s1", ProblemID: "A"},
		},
		mathbattle.ReviewDistribution{
			BetweenParticipants: make(map[string][]string),
			ToOrganizers: []mathbattle.Solution{
				{ID: "s1", ProblemID: "A"},
			},
		})
}

func (s *oneProblemToAll) TestLessThanNeedSolutions() {
	helperTestExpect(s.Require(), s.distributor, s.k,
		[]mathbattle.Solution{
			{ID: "s1", ProblemID: "A", ParticipantID: "p1"},
			{ID: "s2", ProblemID: "A", ParticipantID: "p2"},
		},
		mathbattle.ReviewDistribution{
			BetweenParticipants: map[string][]string{
				"s1": {"p2"},
				"s2": {"p1"},
			},
			ToOrganizers: []mathbattle.Solution{},
		})
}

func (s *oneProblemToAll) TestEnoughSolutions() {
	helperTestExpect(s.Require(), s.distributor, s.k,
		[]mathbattle.Solution{
			{ID: "s1", ProblemID: "A", ParticipantID: "p1"},
			{ID: "s2", ProblemID: "A", ParticipantID: "p2"},
			{ID: "s3", ProblemID: "A", ParticipantID: "p3"},
		}, mathbattle.ReviewDistribution{
			BetweenParticipants: map[string][]string{
				"s1": {"p2", "p3"},
				"s2": {"p1", "p3"},
				"s3": {"p2", "p1"},
			},
			ToOrganizers: []mathbattle.Solution{},
		})

	helperTestExpect(s.Require(), s.distributor, s.k,
		[]mathbattle.Solution{
			{ID: "s1", ProblemID: "A", ParticipantID: "p1"},
			{ID: "s2", ProblemID: "A", ParticipantID: "p2"},
			{ID: "s3", ProblemID: "A", ParticipantID: "p3"},
			{ID: "s4", ProblemID: "A", ParticipantID: "p4"},
		}, mathbattle.ReviewDistribution{
			BetweenParticipants: map[string][]string{
				"s1": {"p2", "p3"},
				"s2": {"p3", "p4"},
				"s3": {"p4", "p1"},
				"s4": {"p1", "p2"},
			},
			ToOrganizers: []mathbattle.Solution{},
		})
}

func TestMap(t *testing.T) {
	solutions := []mathbattle.Solution{}
	for i := 0; i < 4; i++ {
		solutions = append(solutions, mathbattle.Solution{
			ID:            strconv.Itoa(i),
			ParticipantID: strconv.Itoa(i),
		})
	}

	result := MapSolutionsToParticipants(solutions, 3)
	for sID := range result {
		fmt.Printf("%s: ", sID)
		for _, pID := range result[sID] {
			fmt.Printf("%s,", pID)
		}
		fmt.Print("\n")
	}
}

// One problem for all participants
type basicTestSuite struct {
	suite.Suite

	distributor      SolutionDistributor
	problemCount     int
	participantCount int
	k                uint
}

func newBasicTestSuite(problemCount int, participantCount, k int) basicTestSuite {
	return basicTestSuite{
		problemCount:     problemCount,
		participantCount: participantCount,
		k:                uint(k),
	}
}

func (s *basicTestSuite) TestAll() {
	combinations := mocks.GenAllSolutionsCombinations(s.problemCount, s.participantCount)
	for _, c := range combinations {
		r := s.distributor.Get(c, s.k)
		if !IsEachParticipantGotKSolutions(r.BetweenParticipants, s.k) {
			fmt.Println(mathbattle.RoundSolutionsToString(c))
			fmt.Println(r.ToString())
			s.Require().FailNow("Each participant should got equal count of soltuions for review")
		}
		if !IsEachSolutionGoesToKParticiapnts(r.BetweenParticipants, s.k) {
			fmt.Println(r.ToString())
			s.Require().FailNow("Each solution should be reviewed by the same count of participants")
		}
	}
}

func TestSplit(t *testing.T) {
	problemIDs := []string{"A", "B", "C"}

	solutions := []mathbattle.Solution{}
	for i := 0; i < 10; i++ {
		solutions = append(solutions, mathbattle.Solution{
			ID:            strconv.Itoa(i),
			ParticipantID: strconv.Itoa(i),
			ProblemID:     problemIDs[i%len(problemIDs)],
		})
	}

	for _, s := range solutions {
		fmt.Printf("%s:%s, ", s.ID, s.ProblemID)
	}
	fmt.Print("\n")

	result := mathbattle.SplitInGroupsByProblem(solutions)
	for pID := range result {
		fmt.Printf("%s: ", pID)
		for _, s := range result[pID] {
			fmt.Printf("%s,", s.ID)
		}
		fmt.Print("\n")
	}
}

func TestAll(t *testing.T) {
	suite.Run(t, &oneProblemToAll{k: 2})
	//ts := newBasicTestSuite(1, 5, 2)
	//suite.Run(t, &ts)
}
