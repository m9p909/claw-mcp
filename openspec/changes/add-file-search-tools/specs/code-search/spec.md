## ADDED Requirements

### Requirement: Search with literal string
The search_code tool SHALL support literal substring matching across files using Boyer-Moore algorithm.

#### Scenario: Match found with literal query
- **WHEN** user searches for literal string "func main" in a directory
- **THEN** system returns all matching lines with file paths, line numbers, and CRC32 hashes

#### Scenario: No matches found
- **WHEN** user searches for string that doesn't exist
- **THEN** system returns empty results array

#### Scenario: Case-sensitive matching
- **WHEN** user searches for "Error" (capital E)
- **THEN** system MUST NOT match lines containing "error" (lowercase e)

### Requirement: Search with regex pattern
The search_code tool SHALL support regular expression pattern matching.

#### Scenario: Regex pattern matches multiple variants
- **WHEN** user searches with regex "func \w+\("
- **THEN** system returns all function definitions matching the pattern

#### Scenario: Invalid regex pattern
- **WHEN** user provides malformed regex like "func[("
- **THEN** system returns error with code INVALID_REGEX

### Requirement: Hash-compatible output format
Search results SHALL include CRC32 line hashes matching read_file format for edit_file compatibility.

#### Scenario: Hash matches read_file output
- **WHEN** search_code returns line "42:a3f|package main"
- **THEN** read_file for same file MUST return identical hash "a3f" for line 42

#### Scenario: Multi-line match preserves all hashes
- **WHEN** search matches 3 consecutive lines
- **THEN** each line SHALL have independent hash calculated from actual file content

### Requirement: Binary file detection
The search_code tool SHALL detect binary files and report matches without printing unprintable content.

#### Scenario: Binary file contains match
- **WHEN** search finds match in file with NUL byte in first buffer
- **THEN** system returns "Binary file <path> matches" and stops scanning that file

#### Scenario: Text file processed normally
- **WHEN** search finds match in file with no NUL bytes
- **THEN** system returns formatted line with hash and content

### Requirement: Concurrent search performance
The search_code tool SHALL use worker pool pattern for parallel file processing.

#### Scenario: Large directory scanned concurrently
- **WHEN** user searches directory with 1000+ files
- **THEN** system processes files concurrently using configurable worker pool (default 128)

#### Scenario: Workers complete without leaks
- **WHEN** search completes on large directory
- **THEN** all goroutines SHALL terminate and results SHALL be complete

### Requirement: Path validation
The search_code tool SHALL validate and resolve paths before searching.

#### Scenario: Absolute path provided
- **WHEN** user provides "/home/user/project"
- **THEN** system searches that exact path

#### Scenario: Relative path provided
- **WHEN** user provides "./src"
- **THEN** system resolves to absolute path before searching

#### Scenario: Invalid path provided
- **WHEN** user provides non-existent path
- **THEN** system returns error with code PATH_NOT_FOUND

### Requirement: Recursive directory traversal
The search_code tool SHALL recursively traverse directories to find all matching files.

#### Scenario: Nested directory structure
- **WHEN** user searches "/project" containing subdirectories /src, /src/utils, /tests
- **THEN** system searches all files in all subdirectories

#### Scenario: Single file provided
- **WHEN** user provides path to single file
- **THEN** system searches only that file without directory traversal
