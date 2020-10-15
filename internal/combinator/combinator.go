package combinator

import "errors"

var ErrCombinationsEnd = errors.New("End of combinations")

type Combinator struct {
	itemsCombination []int
	itemMax          int
}

func New(itemsCount, itemsMax int) Combinator {
	impl := []int{}
	for i := 0; i < itemsCount-1; i++ {
		impl = append(impl, 0)
	}
	impl = append(impl, -1)

	return Combinator{
		itemsCombination: impl,
		itemMax:          itemsMax,
	}
}

func (c *Combinator) Next() ([]int, error) {
	for i := len(c.itemsCombination) - 1; i >= 0; i-- {
		if c.itemsCombination[i] != c.itemMax {
			c.itemsCombination[i]++
			return c.itemsCombination, nil
		} else {
			c.itemsCombination[i] = 0
		}
	}
	return []int{}, ErrCombinationsEnd
}

func GetRest(combinator Combinator) [][]int {
	result := [][]int{}

	c, err := combinator.Next()
	for err != ErrCombinationsEnd {
		curCombination := make([]int, len(c))
		copy(curCombination, c)
		result = append(result, curCombination)
		c, err = combinator.Next()
	}

	return result
}

func GetAll(itemsCount, itemsMax int) [][]int {
	return GetRest(New(itemsCount, itemsMax))
}
