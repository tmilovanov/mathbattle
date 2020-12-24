package mstd

import (
	"log"
	"strconv"
)

func IndexOf(slice []string, elem string) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == elem {
			return i
		}
	}

	return -1
}

// 0 -> A
// 1 -> B
// 2 -> C
// ...
func IndexToLetter(index int) string {
	if index > 26 {
		log.Panic("index is not supposed to be > 26")
	}
	return string(rune(int('A') + index))
}

// ValidateIndex validates stringified index
// For example, for slice ["A", "B", "C", "D"] indexes = "1", "2", "3", "4" are valid,
// But indexes "0" and "5" are invalid
func ValidateIndex(userInput string, strings []string) (int, bool) {
	index, err := strconv.Atoi(userInput)
	if err != nil {
		return -1, false
	}

	index = index - 1
	if index < 0 || index >= len(strings) {
		return -1, false
	}

	return index, true
}
