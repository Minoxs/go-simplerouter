package simplerouter

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
