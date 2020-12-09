package mstd

func IndexOf(slice []string, elem string) int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == elem {
			return i
		}
	}

	return -1
}
