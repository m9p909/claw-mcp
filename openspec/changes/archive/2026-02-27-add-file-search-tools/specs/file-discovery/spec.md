## ADDED Requirements

### Requirement: List directory contents
The list_directory tool SHALL return all files and subdirectories in a given path.

#### Scenario: Directory with mixed contents
- **WHEN** user lists "/project" containing 3 files and 2 directories
- **THEN** system returns array with 5 entries showing name, type (file|dir), size, and permissions

#### Scenario: Empty directory
- **WHEN** user lists empty directory
- **THEN** system returns empty array with success status

#### Scenario: Permission denied
- **WHEN** user lists directory without read permissions
- **THEN** system returns error with code PERMISSION_DENIED

### Requirement: Find files by name pattern
The find_files tool SHALL locate files matching glob-style name patterns.

#### Scenario: Extension pattern match
- **WHEN** user searches for "*.go" in /project
- **THEN** system returns all .go files in directory tree

#### Scenario: Wildcard pattern match
- **WHEN** user searches for "test_*.py"
- **THEN** system returns all Python test files matching pattern

#### Scenario: No matches found
- **WHEN** user searches for pattern with no matches
- **THEN** system returns empty results array

### Requirement: Generate ASCII directory tree
The tree_directory tool SHALL produce ASCII art visualization of directory structure.

#### Scenario: Nested directory visualization
- **WHEN** user requests tree for /project with subdirectories
- **THEN** system returns ASCII tree using ├──, └──, and │ characters

#### Scenario: Single level directory
- **WHEN** user requests tree for flat directory
- **THEN** system returns simple list with └── for each entry

#### Scenario: Depth limit respected
- **WHEN** user specifies max depth of 2 levels
- **THEN** system MUST NOT traverse beyond 2 levels from root

### Requirement: Path resolution for all tools
All file discovery tools SHALL accept both absolute and relative paths.

#### Scenario: Absolute path used directly
- **WHEN** user provides "/home/user/project"
- **THEN** system uses path without modification

#### Scenario: Relative path resolved
- **WHEN** user provides "../sibling" from /home/user/project
- **THEN** system resolves to /home/user/sibling

#### Scenario: Path does not exist
- **WHEN** user provides non-existent path to any discovery tool
- **THEN** system returns error with code PATH_NOT_FOUND

### Requirement: File metadata in listings
The list_directory tool SHALL include standard file metadata for each entry.

#### Scenario: File entry includes size
- **WHEN** user lists directory containing files
- **THEN** each file entry includes size in bytes

#### Scenario: Directory entry marked distinctly
- **WHEN** user lists directory containing subdirectories
- **THEN** each directory entry has type "dir" and size 0 or omitted

#### Scenario: Hidden files included
- **WHEN** user lists directory containing .hidden files
- **THEN** system includes hidden files in results (no filtering)

### Requirement: Recursive find traversal
The find_files tool SHALL recursively search all subdirectories by default.

#### Scenario: Deep nested match found
- **WHEN** user searches for "config.json" in /project with 5 levels of nesting
- **THEN** system finds all matching files regardless of depth

#### Scenario: Symlink handling
- **WHEN** user searches directory containing symlinks
- **THEN** system MUST NOT follow symlinks to avoid infinite loops

### Requirement: Tree formatting consistency
The tree_directory tool SHALL use consistent ASCII box-drawing characters.

#### Scenario: Last item uses correct marker
- **WHEN** rendering last file in directory
- **THEN** system uses └── (not ├──)

#### Scenario: Intermediate items use tree connector
- **WHEN** rendering non-last items
- **THEN** system uses ├── with │ vertical lines for parent levels

#### Scenario: File vs directory distinction
- **WHEN** rendering tree
- **THEN** directories SHALL be marked distinctly (e.g., trailing "/" or different symbol)
