package maps

// Remaps map[Key][]Value -> map[Value][]Key
func Remap(input map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for key, values := range input {
		for _, val := range values {
			result[val] = append(result[val], key)
		}
	}

	return result
}
