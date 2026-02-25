package hash

import (
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
)

// HashLine returns a CRC32 hash of the line content, formatted as 2-3 char hex
func HashLine(content string) string {
	crc := crc32.ChecksumIEEE([]byte(content))
	// Use 2-3 character hex representation
	hex := fmt.Sprintf("%x", crc)
	if len(hex) > 3 {
		hex = hex[len(hex)-3:]
	}
	// Pad to at least 2 characters
	if len(hex) < 2 {
		hex = "0" + hex
	}
	return hex
}

// FormatLineWithHash returns the line in format: "hash|content"
func FormatLineWithHash(lineNum int, content string) string {
	h := HashLine(content)
	return fmt.Sprintf("%d:%s|%s", lineNum, h, content)
}

// ExtractHashFromLine extracts hash from "linenum:hash|content" format
func ExtractHashFromLine(formatted string) (string, error) {
	parts := strings.SplitN(formatted, "|", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format")
	}

	hashPart := parts[0]
	colonIdx := strings.LastIndex(hashPart, ":")
	if colonIdx == -1 {
		return "", fmt.Errorf("missing hash")
	}

	return hashPart[colonIdx+1:], nil
}

// ValidateHash checks if the provided hash matches the content
func ValidateHash(content string, expectedHash string) bool {
	return HashLine(content) == expectedHash
}

// ExtractLineNumber extracts line number from "linenum:hash|content" format
func ExtractLineNumber(formatted string) (int, error) {
	parts := strings.SplitN(formatted, "|", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid format")
	}

	hashPart := parts[0]
	colonIdx := strings.LastIndex(hashPart, ":")
	if colonIdx == -1 {
		return 0, fmt.Errorf("missing line number")
	}

	numStr := hashPart[:colonIdx]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, fmt.Errorf("invalid line number: %w", err)
	}

	return num, nil
}
