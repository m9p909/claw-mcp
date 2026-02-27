package filesearch

import "bytes"

// stringFinder efficiently finds strings in a source text using Boyer-Moore algorithm.
// Adapted from Go's internal implementation (go1.18.1/src/strings/search.go)
// Modified to use []byte instead of string to reduce allocations.
type stringFinder struct {
	pattern        []byte
	badCharSkip    [256]int
	goodSuffixSkip []int
}

// MakeStringFinder creates a stringFinder for the given pattern.
func MakeStringFinder(pattern []byte) *stringFinder {
	f := &stringFinder{
		pattern:        pattern,
		goodSuffixSkip: make([]int, len(pattern)),
	}
	last := len(pattern) - 1

	// Build bad character table
	for i := range f.badCharSkip {
		f.badCharSkip[i] = len(pattern)
	}
	for i := 0; i < last; i++ {
		f.badCharSkip[pattern[i]] = last - i
	}

	// Build good suffix table
	lastPrefix := last
	for i := last; i >= 0; i-- {
		if bytes.HasPrefix(pattern, pattern[i+1:]) {
			lastPrefix = i + 1
		}
		f.goodSuffixSkip[i] = lastPrefix + last - i
	}
	for i := 0; i < last; i++ {
		lenSuffix := longestCommonSuffix(pattern, pattern[1:i+1])
		if pattern[i-lenSuffix] != pattern[last-lenSuffix] {
			f.goodSuffixSkip[last-lenSuffix] = lenSuffix + last - i
		}
	}

	return f
}

// longestCommonSuffix returns the length of the longest common suffix between a and b.
func longestCommonSuffix(a, b []byte) (i int) {
	for ; i < len(a) && i < len(b); i++ {
		if a[len(a)-1-i] != b[len(b)-1-i] {
			break
		}
	}
	return
}

// next returns the index in text of the first occurrence of the pattern.
// If the pattern is not found, it returns -1.
func (f *stringFinder) next(text []byte) int {
	i := len(f.pattern) - 1
	for i < len(text) {
		j := len(f.pattern) - 1
		for j >= 0 && text[i] == f.pattern[j] {
			i--
			j--
		}
		if j < 0 {
			return i + 1
		}
		i += max(f.badCharSkip[text[i]], f.goodSuffixSkip[j])
	}
	return -1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
