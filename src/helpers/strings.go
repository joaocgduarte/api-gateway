package helpers

func InArray(needle string, haystack []string) bool {
	for _, val := range haystack {
		if needle == val {
			return true
		}
	}
	return false
}
