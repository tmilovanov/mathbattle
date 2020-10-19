package containers

import (
	"sort"
)

// StringUsageCounter calculates how many times each item was used
// and return on demand less used items
type StringUsageCounter struct {
	usage map[string]int
}

type StringItemUsage struct {
	Item     string
	UseCount int
}

func NewUsageCounter() StringUsageCounter {
	return StringUsageCounter{
		usage: make(map[string]int),
	}
}

func (c *StringUsageCounter) AddItems(items []string) {
	for _, item := range items {
		_, ok := c.usage[item]
		if !ok {
			c.usage[item] = 0
		}
	}
}

// GetSortedUsage returns all items sorted by usage, from less used to more used
func (c *StringUsageCounter) GetSortedByUsage() []StringItemUsage {
	result := []StringItemUsage{}

	for item, count := range c.usage {
		result = append(result, StringItemUsage{
			Item:     item,
			UseCount: count,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UseCount < result[j].UseCount
	})

	return result
}

func (c *StringUsageCounter) SortByUsage(items []string) []StringItemUsage {
	result := []StringItemUsage{}

	for _, item := range items {
		result = append(result, StringItemUsage{
			Item:     item,
			UseCount: c.usage[item],
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UseCount < result[j].UseCount
	})

	return result
}

func (c *StringUsageCounter) Use(item string) {
	c.usage[item]++
}

func (c *StringUsageCounter) UseMostUnpopular(count int) []string {
	result := []string{}

	usage := c.GetSortedByUsage()
	for i := 0; i < count; i++ {
		usedItem := usage[i].Item
		c.Use(usedItem)
		result = append(result, usedItem)
	}

	return result
}

func (c *StringUsageCounter) UseMostUnpopularFromSet(items []string, count int) []string {
	result := []string{}

	usage := c.SortByUsage(items)
	for i := 0; i < count; i++ {
		usedItem := usage[i].Item
		c.Use(usedItem)
		result = append(result, usedItem)
	}

	return result
}
