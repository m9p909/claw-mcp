package filesearch

import (
	"os"
	"path/filepath"
	"strings"
)

type DirEntry struct {
	Name        string
	Type        string // "file" or "dir"
	Size        int64
	Permissions string
}

type FileMatch struct {
	Path     string
	Size     int64
	Modified int64
}

// ListDirectory returns all files and directories in the given path.
func ListDirectory(path string) ([]DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var result []DirEntry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		entryType := "file"
		size := info.Size()
		if entry.IsDir() {
			entryType = "dir"
			size = 0
		}

		result = append(result, DirEntry{
			Name:        entry.Name(),
			Type:        entryType,
			Size:        size,
			Permissions: info.Mode().String(),
		})
	}
	return result, nil
}

// FindFiles finds files matching a glob pattern recursively.
func FindFiles(rootPath string, pattern string) ([]FileMatch, error) {
	var matches []FileMatch

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Detect symlinks to prevent infinite loops
		if info.Mode()&os.ModeSymlink != 0 {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}

		if matched {
			matches = append(matches, FileMatch{
				Path:     path,
				Size:     info.Size(),
				Modified: info.ModTime().Unix(),
			})
		}
		return nil
	})

	return matches, err
}

// TreeDirectory generates an ASCII tree visualization of a directory structure.
func TreeDirectory(rootPath string, maxDepth int) (string, error) {
	var sb strings.Builder
	info, err := os.Stat(rootPath)
	if err != nil {
		return "", err
	}

	sb.WriteString(filepath.Base(rootPath))
	if info.IsDir() {
		sb.WriteString("/")
	}
	sb.WriteString("\n")

	if info.IsDir() {
		if err := buildTree(rootPath, "", 0, maxDepth, &sb); err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

func buildTree(path string, prefix string, depth int, maxDepth int, sb *strings.Builder) error {
	if maxDepth > 0 && depth >= maxDepth {
		return nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil // Skip unreadable directories
	}

	for i, entry := range entries {
		isLast := i == len(entries)-1
		connector := "├── "
		newPrefix := prefix + "│   "
		if isLast {
			connector = "└── "
			newPrefix = prefix + "    "
		}

		// Detect symlinks
		info, err := entry.Info()
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			continue // Skip symlinks to prevent infinite loops
		}

		name := entry.Name()
		if entry.IsDir() {
			name += "/"
		}

		sb.WriteString(prefix)
		sb.WriteString(connector)
		sb.WriteString(name)
		sb.WriteString("\n")

		if entry.IsDir() {
			subPath := filepath.Join(path, entry.Name())
			buildTree(subPath, newPrefix, depth+1, maxDepth, sb)
		}
	}

	return nil
}
