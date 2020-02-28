package util

func Contains(slice []string, item string) bool {
	for _, val := range slice {
		if val == item {
			return true
		}
	}

	return false
}
