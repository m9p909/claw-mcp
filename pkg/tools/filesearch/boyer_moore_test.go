package filesearch

import "testing"

func TestBoyerMoore_BasicMatching(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		text    string
		want    int
	}{
		{"simple match", "hello", "say hello world", 4},
		{"no match", "missing", "some text here", -1},
		{"match at start", "test", "test this", 0},
		{"match at end", "end", "this is the end", 12},
		{"multiple occurrences returns first", "the", "the cat and the dog", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := MakeStringFinder([]byte(tt.pattern))
			got := finder.next([]byte(tt.text))
			if got != tt.want {
				t.Errorf("next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoyerMoore_EdgeCases(t *testing.T) {
	t.Run("empty text", func(t *testing.T) {
		finder := MakeStringFinder([]byte("pattern"))
		got := finder.next([]byte(""))
		if got != -1 {
			t.Errorf("next() = %v, want -1 for empty text", got)
		}
	})

	t.Run("pattern longer than text", func(t *testing.T) {
		finder := MakeStringFinder([]byte("longpattern"))
		got := finder.next([]byte("tiny"))
		if got != -1 {
			t.Errorf("next() = %v, want -1 for pattern longer than text", got)
		}
	})

	t.Run("single character pattern", func(t *testing.T) {
		finder := MakeStringFinder([]byte("x"))
		got := finder.next([]byte("abxcd"))
		if got != 2 {
			t.Errorf("next() = %v, want 2 for single char pattern", got)
		}
	})

	t.Run("single character text and pattern match", func(t *testing.T) {
		finder := MakeStringFinder([]byte("a"))
		got := finder.next([]byte("a"))
		if got != 0 {
			t.Errorf("next() = %v, want 0 for single char match", got)
		}
	})

	t.Run("single character text and pattern no match", func(t *testing.T) {
		finder := MakeStringFinder([]byte("a"))
		got := finder.next([]byte("b"))
		if got != -1 {
			t.Errorf("next() = %v, want -1 for single char no match", got)
		}
	})
}

func TestBoyerMoore_ByteSliceOptimization(t *testing.T) {
	t.Run("handles binary data", func(t *testing.T) {
		pattern := []byte{0x00, 0xFF, 0xAB}
		text := []byte{0x12, 0x34, 0x00, 0xFF, 0xAB, 0x56}
		finder := MakeStringFinder(pattern)
		got := finder.next(text)
		if got != 2 {
			t.Errorf("next() = %v, want 2 for binary pattern", got)
		}
	})
}

func TestLongestCommonSuffix(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want int
	}{
		{"no common suffix", "abc", "def", 0},
		{"partial suffix", "testing", "running", 3}, // "ing"
		{"full suffix", "abc", "xyzabc", 3},
		{"identical", "test", "test", 4},
		{"one empty", "test", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := longestCommonSuffix([]byte(tt.a), []byte(tt.b))
			if got != tt.want {
				t.Errorf("longestCommonSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}
