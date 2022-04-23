package simplerouter

// cleanPattern
// Makes sure the patterns are in the right format
// That is, always a leading slash, and no trailing slashes
func cleanPattern(pattern string) string {
	if len(pattern) == 0 {
		return "/"
	}

	// Add leading slash
	if pattern[0] != '/' {
		pattern = "/" + pattern
	}

	// Remove trailing slash
	size := len(pattern)
	if pattern[size-1] == '/' {
		pattern = pattern[:size-1]
	}

	return pattern
}
