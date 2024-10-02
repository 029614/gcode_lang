package array

// ContainsSingleThreaded checks if a given item is in the slice using a single thread.
func ContainsSingleThreaded(slice []any, item int) bool {
	var v any
	for _, v = range slice {
		if v == item {
			return true
		}
	}
	return false
}
