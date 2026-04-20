package utils

import "regexp"

func IsMatchRegexp(s, pattern string) bool {
	if s == "" || pattern == "" {
		return false
	}

	regexp := regexp.MustCompile(pattern)
	return regexp.MatchString(s)
}
