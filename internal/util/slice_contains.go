package util

func SliceContains(slice []string, value string) bool {
	var result bool

	for _, x := range slice {
		if x == value {
			result = true
			break
		}
	}

	return result
}
