package filesearch

import (
	"os"
	"path/filepath"
	"testing"

	"awesomeProject/pkg/hash"
)

func TestHashConsistency(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := "line one\nline two with match\nline three"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Search for a match
	opts := &SearchOptions{
		Kind:   LITERAL,
		Finder: MakeStringFinder([]byte("match")),
	}

	resultChan := make(chan SearchResult, 10)
	go func() {
		Search([]string{testFile}, opts, 1, resultChan)
		close(resultChan)
	}()

	var results []SearchResult
	for result := range resultChan {
		results = append(results, result)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 match, got %d", len(results))
	}

	searchResult := results[0]

	// Verify hash matches what hash.HashLine would produce
	expectedHash := hash.HashLine("line two with match")
	if searchResult.Hash != expectedHash {
		t.Errorf("Hash mismatch: search returned %q, hash.HashLine returned %q", searchResult.Hash, expectedHash)
	}

	// Verify line number is correct
	if searchResult.Line != 2 {
		t.Errorf("Expected line 2, got %d", searchResult.Line)
	}

	// Verify content is exact
	if searchResult.Content != "line two with match" {
		t.Errorf("Content mismatch: got %q", searchResult.Content)
	}
}

func TestHashConsistencyMultipleLines(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "multi.txt")
	lines := []string{
		"package main",
		"import \"fmt\"",
		"func main() {",
		"    fmt.Println(\"test\")",
		"}",
	}
	content := ""
	for i, line := range lines {
		if i > 0 {
			content += "\n"
		}
		content += line
	}
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	// Search for "fmt"
	opts := &SearchOptions{
		Kind:   LITERAL,
		Finder: MakeStringFinder([]byte("fmt")),
	}

	resultChan := make(chan SearchResult, 10)
	go func() {
		Search([]string{testFile}, opts, 1, resultChan)
		close(resultChan)
	}()

	var results []SearchResult
	for result := range resultChan {
		results = append(results, result)
	}

	if len(results) != 2 {
		t.Fatalf("Expected 2 matches, got %d", len(results))
	}

	// Check both results have correct hashes
	for _, result := range results {
		expectedHash := hash.HashLine(result.Content)
		if result.Hash != expectedHash {
			t.Errorf("Line %d hash mismatch: got %q, expected %q", result.Line, result.Hash, expectedHash)
		}
	}
}
