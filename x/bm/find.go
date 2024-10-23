package bm

import "strings"

func FindIter(s string, b rune) bool {
	for _, c := range s {
		if c == b {
			return true
		}
	}
	return false
}

func FindContains(s string, b rune) bool {
	return strings.Contains(s, string(b))
}

