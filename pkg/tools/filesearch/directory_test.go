package filesearch

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListDirectory(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files and directories
	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("test"), 0644)
	os.WriteFile(filepath.Join(tempDir, "file2.go"), []byte("package main"), 0644)
	os.Mkdir(filepath.Join(tempDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(tempDir, ".hidden"), []byte("hidden"), 0644)

	entries, err := ListDirectory(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(entries))
	}

	// Check for specific entries
	foundFile := false
	foundDir := false
	foundHidden := false
	for _, entry := range entries {
		if entry.Name == "file1.txt" && entry.Type == "file" {
			foundFile = true
			if entry.Size == 0 {
				t.Error("File should have non-zero size")
			}
		}
		if entry.Name == "subdir" && entry.Type == "dir" {
			foundDir = true
			if entry.Size != 0 {
				t.Error("Directory should have size 0")
			}
		}
		if entry.Name == ".hidden" {
			foundHidden = true
		}
	}

	if !foundFile {
		t.Error("file1.txt not found")
	}
	if !foundDir {
		t.Error("subdir not found")
	}
	if !foundHidden {
		t.Error("Hidden file should be included")
	}
}

func TestListDirectoryEmpty(t *testing.T) {
	tempDir := t.TempDir()

	entries, err := ListDirectory(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries in empty dir, got %d", len(entries))
	}
}

func TestFindFiles(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "sub")
	os.Mkdir(subDir, 0755)

	// Create test files
	os.WriteFile(filepath.Join(tempDir, "test1.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tempDir, "test2.go"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tempDir, "readme.md"), []byte(""), 0644)
	os.WriteFile(filepath.Join(subDir, "sub_test.go"), []byte(""), 0644)

	// Find all .go files
	matches, err := FindFiles(tempDir, "*.go")
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 3 {
		t.Errorf("Expected 3 .go files, got %d", len(matches))
	}

	// Find test_*.go pattern
	matches, err = FindFiles(tempDir, "test*.go")
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 2 {
		t.Errorf("Expected 2 test*.go files, got %d", len(matches))
	}
}

func TestFindFilesNoMatch(t *testing.T) {
	tempDir := t.TempDir()
	os.WriteFile(filepath.Join(tempDir, "file.txt"), []byte(""), 0644)

	matches, err := FindFiles(tempDir, "*.go")
	if err != nil {
		t.Fatal(err)
	}

	if len(matches) != 0 {
		t.Errorf("Expected 0 matches, got %d", len(matches))
	}
}

func TestTreeDirectory(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	os.Mkdir(subDir, 0755)

	os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tempDir, "file2.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte(""), 0644)

	tree, err := TreeDirectory(tempDir, 0)
	if err != nil {
		t.Fatal(err)
	}

	// Check for ASCII tree characters
	if !strings.Contains(tree, "├──") && !strings.Contains(tree, "└──") {
		t.Error("Tree should contain box-drawing characters")
	}

	// Check for directory marker
	if !strings.Contains(tree, "subdir/") {
		t.Error("Tree should mark directories with /")
	}

	// Check for files
	if !strings.Contains(tree, "file1.txt") {
		t.Error("Tree should contain file1.txt")
	}
	if !strings.Contains(tree, "nested.txt") {
		t.Error("Tree should contain nested file")
	}
}

func TestTreeDirectoryDepthLimit(t *testing.T) {
	tempDir := t.TempDir()
	level1 := filepath.Join(tempDir, "level1")
	level2 := filepath.Join(level1, "level2")
	os.MkdirAll(level2, 0755)

	os.WriteFile(filepath.Join(tempDir, "root.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(level1, "l1.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(level2, "l2.txt"), []byte(""), 0644)

	// Max depth 1: shows immediate children only, not contents of subdirectories
	tree, err := TreeDirectory(tempDir, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Should show level1/ directory but not l1.txt inside it
	if !strings.Contains(tree, "level1/") {
		t.Error("Should include level1 directory")
	}
	if strings.Contains(tree, "l1.txt") {
		t.Error("Should NOT include files inside subdirectories with max_depth=1")
	}
	if strings.Contains(tree, "l2.txt") {
		t.Error("Should NOT include level 2 file with max_depth=1")
	}

	// Max depth 2: shows 2 levels down
	tree2, err := TreeDirectory(tempDir, 2)
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(tree2, "l1.txt") {
		t.Error("Should include level 1 file with max_depth=2")
	}
	if strings.Contains(tree2, "l2.txt") {
		t.Error("Should NOT include level 2 file with max_depth=2")
	}
}

func TestTreeDirectoryFormattingLastItem(t *testing.T) {
	tempDir := t.TempDir()
	os.WriteFile(filepath.Join(tempDir, "aaa.txt"), []byte(""), 0644)
	os.WriteFile(filepath.Join(tempDir, "zzz.txt"), []byte(""), 0644)

	tree, err := TreeDirectory(tempDir, 0)
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(tree, "\n")
	// Last item should use └── not ├──
	foundLastMarker := false
	for _, line := range lines {
		if strings.Contains(line, "└──") {
			foundLastMarker = true
			break
		}
	}

	if !foundLastMarker {
		t.Error("Tree should use └── for last items")
	}
}

func TestSymlinkSkipping(t *testing.T) {
	tempDir := t.TempDir()
	realFile := filepath.Join(tempDir, "real.txt")
	os.WriteFile(realFile, []byte("test"), 0644)

	// Create symlink
	linkPath := filepath.Join(tempDir, "link.txt")
	if err := os.Symlink(realFile, linkPath); err != nil {
		t.Skip("Symlinks not supported on this system")
	}

	// FindFiles should skip symlinks
	matches, err := FindFiles(tempDir, "*.txt")
	if err != nil {
		t.Fatal(err)
	}

	// Should only find the real file, not the symlink
	if len(matches) != 1 {
		t.Errorf("Expected 1 file (symlink should be skipped), got %d", len(matches))
	}
}
