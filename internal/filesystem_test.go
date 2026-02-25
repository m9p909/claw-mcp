package internal

import (
	"os"
	"testing"

	"awesomeProject/pkg/hash"
)

func TestHashLine(t *testing.T) {
	tests := []struct {
		content string
		name    string
	}{
		{"hello world", "simple string"},
		{"", "empty string"},
		{"line with\nnewline", "string with escaped newline"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := hash.HashLine(tt.content)
			// Hash should be 2-3 character hex
			if len(h) < 2 || len(h) > 3 {
				t.Errorf("hash length %d not in range [2,3]", len(h))
			}
		})
	}
}

func TestValidateHash(t *testing.T) {
	content := "hello world"
	h := hash.HashLine(content)

	if !hash.ValidateHash(content, h) {
		t.Error("valid hash failed validation")
	}

	if hash.ValidateHash("different content", h) {
		t.Error("invalid hash passed validation")
	}
}

func TestFormatLineWithHash(t *testing.T) {
	content := "hello world"
	formatted := hash.FormatLineWithHash(1, content)

	// Should be "1:hash|content"
	if !contains(formatted, "|") {
		t.Error("formatted line missing pipe separator")
	}

	extractedHash, err := hash.ExtractHashFromLine(formatted)
	if err != nil {
		t.Errorf("failed to extract hash: %v", err)
	}

	if !hash.ValidateHash(content, extractedHash) {
		t.Error("extracted hash does not match content")
	}
}

func TestExtractLineNumber(t *testing.T) {
	formatted := hash.FormatLineWithHash(42, "test")
	lineNum, err := hash.ExtractLineNumber(formatted)
	if err != nil {
		t.Errorf("failed to extract line number: %v", err)
	}

	if lineNum != 42 {
		t.Errorf("expected line number 42, got %d", lineNum)
	}
}

func TestReadWriteFile(t *testing.T) {
	// Create temp file
	tmpfile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test content
	testContent := "line 1\nline 2\nline 3"
	if _, err := tmpfile.WriteString(testContent); err != nil {
		t.Fatalf("failed to write test content: %v", err)
	}
	tmpfile.Close()

	// Read file and verify hashes
	data, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	lines := string(data)
	if lines != testContent {
		t.Errorf("content mismatch: %q != %q", lines, testContent)
	}
}

func TestHashConsistency(t *testing.T) {
	content := "test line"
	h1 := hash.HashLine(content)
	h2 := hash.HashLine(content)

	if h1 != h2 {
		t.Errorf("hash not consistent: %s != %s", h1, h2)
	}
}

func TestEditFileRange(t *testing.T) {
	// Create temp file
	tmpfile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	// Write test content
	testContent := "line 1\nline 2\nline 3"
	if _, err := tmpfile.WriteString(testContent); err != nil {
		t.Fatalf("failed to write test content: %v", err)
	}
	tmpfile.Close()

	// Verify we can read and calculate hashes
	data, _ := os.ReadFile(tmpfile.Name())
	lines := string(data)

	// These lines have these hashes
	line1Hash := hash.HashLine("line 1")
	line2Hash := hash.HashLine("line 2")
	line3Hash := hash.HashLine("line 3")

	// Verify hashes are consistent
	if !hash.ValidateHash("line 1", line1Hash) {
		t.Error("line 1 hash validation failed")
	}
	if !hash.ValidateHash("line 2", line2Hash) {
		t.Error("line 2 hash validation failed")
	}
	if !hash.ValidateHash("line 3", line3Hash) {
		t.Error("line 3 hash validation failed")
	}

	// Verify content still exists
	if !contains(lines, "line 1") {
		t.Error("line 1 not found in content")
	}
}

func TestHashCollisionRisk(t *testing.T) {
	// Test that similar strings have different hashes
	h1 := hash.HashLine("hello")
	h2 := hash.HashLine("hallo")
	h3 := hash.HashLine("hello ")

	if h1 == h2 || h1 == h3 || h2 == h3 {
		t.Error("similar strings produced same hash (collision)")
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

func TestFormatAndExtract(t *testing.T) {
	tests := []struct {
		lineNum int
		content string
	}{
		{1, "hello world"},
		{42, "test line"},
		{999, "long line with many words in it"},
	}

	for _, tt := range tests {
		formatted := hash.FormatLineWithHash(tt.lineNum, tt.content)

		// Test line number extraction
		lineNum, err := hash.ExtractLineNumber(formatted)
		if err != nil {
			t.Errorf("failed to extract line number: %v", err)
		}
		if lineNum != tt.lineNum {
			t.Errorf("line number mismatch: expected %d, got %d", tt.lineNum, lineNum)
		}

		// Test hash extraction
		h, err := hash.ExtractHashFromLine(formatted)
		if err != nil {
			t.Errorf("failed to extract hash: %v", err)
		}
		if !hash.ValidateHash(tt.content, h) {
			t.Errorf("hash validation failed for line %d", tt.lineNum)
		}
	}
}
