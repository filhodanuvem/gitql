package utilities

// IsFieldPresentInArray checks the array of strings for the given field name.
func IsFieldPresentInArray(arr []string, element string) bool {
	for _, fieldInSelect := range arr {
		if fieldInSelect == element {
			return true
		}
	}
	return false
}
