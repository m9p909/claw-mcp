## ADDED Requirements

### Requirement: Read file with content hashing
The system SHALL read file contents and return each line tagged with a content hash (CRC32, 2-3 character hex). Clients use these hashes as stable anchors for subsequent edits.

#### Scenario: Read existing file
- **WHEN** client calls read_file with valid file path
- **THEN** system returns all lines prefixed with line number and hash (format: `1:a3|content`)

#### Scenario: Read with offset and limit
- **WHEN** client specifies offset (line number) and limit (max lines)
- **THEN** system returns only the requested line range

#### Scenario: File not found
- **WHEN** client reads non-existent file
- **THEN** system returns FileNotFoundError with message

### Requirement: Write file with hash validation
The system SHALL write file contents, optionally validating hashes if content includes them. Hashes must match file state before overwrite.

#### Scenario: Write new file
- **WHEN** client writes to non-existent path
- **THEN** system creates file with content (hashes stripped if present)

#### Scenario: Write with hash validation
- **WHEN** client sends content with hashes matching current file state
- **THEN** system overwrites file successfully

#### Scenario: Hash mismatch on write
- **WHEN** client sends content with hashes that don't match file
- **THEN** system rejects write with HashMismatchError

#### Scenario: Permission denied
- **WHEN** client lacks write permission on path
- **THEN** system returns PermissionError

### Requirement: Edit file using hash ranges
The system SHALL edit files by replacing lines identified by hash range. Hash validation is strict: if hashes don't match, edit is rejected.

#### Scenario: Edit single line
- **WHEN** client specifies start_hash and end_hash for same line
- **THEN** system replaces that line with new_content

#### Scenario: Edit multiple lines
- **WHEN** client specifies start_hash and end_hash for different lines
- **THEN** system replaces all lines in range (inclusive) with new_content

#### Scenario: Hash mismatch
- **WHEN** start_hash or end_hash don't exist in file
- **THEN** system returns HashMismatchError, file unchanged

#### Scenario: File changed since read
- **WHEN** file was modified after client read it (hashes stale)
- **THEN** system rejects edit with HashMismatchError
