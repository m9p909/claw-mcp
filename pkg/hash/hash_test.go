package hash

import (
	"testing"
)

func TestHashLine(t *testing.T) {
	tests := []struct {
		content string
		name    string
	}{
		{"hello world", "simple string"},
		{"", "empty string"},
		{"line with special chars !@#$%", "special characters"},
		{"function main() {", "code line"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := HashLine(tt.content)
			// Hash should be 2-3 character hex
			if len(h) < 2 || len(h) > 3 {
				t.Errorf("hash length %d not in range [2,3]", len(h))
			}
			// Should be valid hex
			for _, c := range h {
				if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
					t.Errorf("invalid hex character: %c", c)
				}
			}
		})
	}
}

func TestValidateHash(t *testing.T) {
	content := "hello world"
	h := HashLine(content)

	if !ValidateHash(content, h) {
		t.Error("valid hash failed validation")
	}

	if ValidateHash("different content", h) {
		t.Error("invalid hash passed validation")
	}
}

func TestFormatLineWithHash(t *testing.T) {
	content := "test line"
	lineNum := 5
	formatted := FormatLineWithHash(lineNum, content)

	// Should contain pipe separator
	if !contains(formatted, "|") {
		t.Error("formatted line missing pipe separator")
	}

	// Should start with line number
	if !contains(formatted, "5:") {
		t.Error("formatted line missing line number")
	}

	// Extract and validate hash
	extractedHash, err := ExtractHashFromLine(formatted)
	if err != nil {
		t.Errorf("failed to extract hash: %v", err)
	}

	if !ValidateHash(content, extractedHash) {
		t.Error("extracted hash does not match content")
	}
}

func TestExtractLineNumber(t *testing.T) {
	tests := []struct {
		lineNum int
		content string
	}{
		{1, "first line"},
		{42, "middle line"},
		{999, "large number"},
	}

	for _, tt := range tests {
		formatted := FormatLineWithHash(tt.lineNum, tt.content)
		extractedLineNum, err := ExtractLineNumber(formatted)
		if err != nil {
			t.Errorf("failed to extract line number: %v", err)
		}

		if extractedLineNum != tt.lineNum {
			t.Errorf("line number mismatch: expected %d, got %d", tt.lineNum, extractedLineNum)
		}
	}
}

func TestExtractHashFromLine(t *testing.T) {
	tests := []struct {
		content string
	}{
		{"simple"},
		{"with spaces here"},
		{"123 456"},
	}

	for _, tt := range tests {
		formatted := FormatLineWithHash(1, tt.content)
		extractedHash, err := ExtractHashFromLine(formatted)
		if err != nil {
			t.Errorf("failed to extract hash: %v", err)
		}

		expectedHash := HashLine(tt.content)
		if extractedHash != expectedHash {
			t.Errorf("extracted hash %s != expected %s", extractedHash, expectedHash)
		}
	}
}

func TestHashConsistency(t *testing.T) {
	content := "consistency test"
	h1 := HashLine(content)
	h2 := HashLine(content)
	h3 := HashLine(content)

	if h1 != h2 || h2 != h3 {
		t.Errorf("hash not consistent: %s, %s, %s", h1, h2, h3)
	}
}

func TestHashCollisionRisk(t *testing.T) {
	// Test that similar strings have different hashes (most of the time)
	testPairs := []struct {
		a string
		b string
	}{
		{"hello", "hallo"},
		{"test", "tests"},
		{"func", "funcs"},
		{"foo", "foo "},
	}

	collisions := 0
	for _, pair := range testPairs {
		h1 := HashLine(pair.a)
		h2 := HashLine(pair.b)
		if h1 == h2 {
			collisions++
		}
	}

	// Allow up to 1 collision in this small test (CRC32 truncation has collision risk)
	if collisions > 1 {
		t.Errorf("too many hash collisions: %d", collisions)
	}
}

func TestFormatAndExtractRoundtrip(t *testing.T) {
	tests := []struct {
		lineNum int
		content string
	}{
		{1, "hello world"},
		{42, "test line"},
		{999, "long line with many words in it for testing purposes"},
	}

	for _, tt := range tests {
		t.Run(tt.content, func(t *testing.T) {
			// Format
			formatted := FormatLineWithHash(tt.lineNum, tt.content)

			// Extract line number
			extractedLineNum, err := ExtractLineNumber(formatted)
			if err != nil {
				t.Fatalf("failed to extract line number: %v", err)
			}
			if extractedLineNum != tt.lineNum {
				t.Errorf("line number mismatch: %d != %d", extractedLineNum, tt.lineNum)
			}

			// Extract hash
			extractedHash, err := ExtractHashFromLine(formatted)
			if err != nil {
				t.Fatalf("failed to extract hash: %v", err)
			}

			// Validate hash
			if !ValidateHash(tt.content, extractedHash) {
				t.Errorf("hash validation failed for content: %q", tt.content)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	edgeCases := []string{
		"",
		"\n",
		"\t",
		"   ",
		string([]byte{0, 1, 2, 3}), // non-printable characters
		"你好世界",                      // unicode
		"emoji: 😀",
	}

	for _, content := range edgeCases {
		h := HashLine(content)
		// Should produce a valid hash
		if len(h) < 2 || len(h) > 3 {
			t.Errorf("invalid hash length for content %q: %s", content, h)
		}

		// Should be able to validate
		if !ValidateHash(content, h) {
			t.Errorf("hash validation failed for content %q", content)
		}
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
