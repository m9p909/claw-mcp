package filesearch

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestSearchLiteral(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	content := "line one\nline two with pattern\nline three\nline two with pattern again"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	opts := &SearchOptions{
		Kind:   LITERAL,
		Finder: MakeStringFinder([]byte("pattern")),
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
		t.Errorf("Expected 2 matches, got %d", len(results))
	}

	if results[0].Line != 2 {
		t.Errorf("First match should be line 2, got %d", results[0].Line)
	}
	if results[1].Line != 4 {
		t.Errorf("Second match should be line 4, got %d", results[1].Line)
	}
}

func TestSearchRegex(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")
	content := "package main\nfunc foo() {}\nfunc bar() {}\nvar x = 1"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	re := regexp.MustCompile(`func \w+\(`)
	opts := &SearchOptions{
		Kind:  REGEX,
		Regex: re,
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
		t.Errorf("Expected 2 function matches, got %d", len(results))
	}
}

func TestSearchBinaryFile(t *testing.T) {
	tempDir := t.TempDir()
	binFile := filepath.Join(tempDir, "binary.dat")
	binaryContent := []byte{0x00, 0xFF, 'p', 'a', 't', 't', 'e', 'r', 'n', 0xAB}
	if err := os.WriteFile(binFile, binaryContent, 0644); err != nil {
		t.Fatal(err)
	}

	opts := &SearchOptions{
		Kind:   LITERAL,
		Finder: MakeStringFinder([]byte("pattern")),
	}

	resultChan := make(chan SearchResult, 10)
	go func() {
		Search([]string{binFile}, opts, 1, resultChan)
		close(resultChan)
	}()

	var results []SearchResult
	for result := range resultChan {
		results = append(results, result)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 binary file match, got %d", len(results))
	}

	if results[0].Line != 0 {
		t.Errorf("Binary file match should have line 0, got %d", results[0].Line)
	}
}

func TestSearchDirectory(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")

	if err := os.WriteFile(file1, []byte("match here\nno match"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file2, []byte("another match here"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := &SearchOptions{
		Kind:   LITERAL,
		Finder: MakeStringFinder([]byte("match")),
	}

	resultChan := make(chan SearchResult, 10)
	go func() {
		Search([]string{tempDir}, opts, 2, resultChan)
		close(resultChan)
	}()

	var results []SearchResult
	for result := range resultChan {
		results = append(results, result)
	}

	// Should find "match" in file1.txt line 1 and file2.txt line 1
	// (file1.txt has "match" in first line and "no match" in second - both have "match")
	if len(results) != 3 {
		t.Errorf("Expected 3 matches across files (2 in file1, 1 in file2), got %d", len(results))
	}
}

func TestSearchNoMatch(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("no pattern here"), 0644); err != nil {
		t.Fatal(err)
	}

	opts := &SearchOptions{
		Kind:   LITERAL,
		Finder: MakeStringFinder([]byte("missing")),
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

	if len(results) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(results))
	}
}
