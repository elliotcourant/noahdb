package util

// Pluralize returns a single character 's' unless n == 1.
func Pluralize(n int64) string {
	if n == 1 {
		return ""
	}
	return "s"
}
