package combinator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCombinator(t *testing.T) {
	req := require.New(t)

	req.Equal([][]int{
		{0},
	}, GetRest(New(1, 0)))

	req.Equal([][]int{
		{0}, {1},
	}, GetRest(New(1, 1)))

	req.Equal([][]int{
		{0}, {1}, {2},
	}, GetRest(New(1, 2)))

	req.Equal([][]int{
		{0, 0},
	}, GetRest(New(2, 0)))

	req.Equal([][]int{
		{0, 0}, {0, 1},
		{1, 0}, {1, 1},
	}, GetRest(New(2, 1)))

	req.Equal([][]int{
		{0, 0}, {0, 1}, {0, 2},
		{1, 0}, {1, 1}, {1, 2},
		{2, 0}, {2, 1}, {2, 2},
	}, GetRest(New(2, 2)))

	req.Equal([][]int{
		{0, 0}, {0, 1}, {0, 2}, {0, 3},
		{1, 0}, {1, 1}, {1, 2}, {1, 3},
		{2, 0}, {2, 1}, {2, 2}, {2, 3},
		{3, 0}, {3, 1}, {3, 2}, {3, 3},
	}, GetRest(New(2, 3)))
}
