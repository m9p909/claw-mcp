## ADDED Requirements

### Requirement: Boyer-Moore string finder initialization
The system SHALL implement Boyer-Moore algorithm with badCharSkip and goodSuffixSkip tables.

#### Scenario: Pattern compiled into finder
- **WHEN** creating string finder for pattern "search"
- **THEN** system initializes badCharSkip table with 256 entries and goodSuffixSkip table with pattern length entries

#### Scenario: Skip tables correctly calculated
- **WHEN** pattern is "example"
- **THEN** badCharSkip['x'] SHALL equal distance from 'x' to end of pattern

### Requirement: Substring search with skip optimization
The Boyer-Moore finder SHALL use skip tables to avoid unnecessary character comparisons.

#### Scenario: Mismatch triggers skip
- **WHEN** searching for "test" in "this is a best case"
- **THEN** system skips ahead using max(badCharSkip, goodSuffixSkip) on mismatch

#### Scenario: Match found and position returned
- **WHEN** searching for "hello" in "say hello world"
- **THEN** system returns index 4 (start position of match)

#### Scenario: No match returns -1
- **WHEN** searching for "missing" in "some text here"
- **THEN** system returns -1

### Requirement: Byte slice optimization
The string finder SHALL operate on []byte instead of string to reduce allocations.

#### Scenario: No string-to-bytes conversion in hot path
- **WHEN** searching file content already in []byte form
- **THEN** system performs search without allocating new string copies

#### Scenario: Pattern stored as bytes
- **WHEN** initializing finder with pattern
- **THEN** pattern is stored as []byte internally

### Requirement: Longest common suffix calculation
The system SHALL correctly compute longest common suffix for goodSuffixSkip table.

#### Scenario: Suffix exists earlier in pattern
- **WHEN** pattern contains repeated suffix "iss" in "mississi"
- **THEN** goodSuffixSkip correctly identifies earlier occurrence

#### Scenario: Suffix shares prefix
- **WHEN** pattern suffix shares characters with pattern prefix
- **THEN** goodSuffixSkip allows minimal shift aligning prefix with suffix end

### Requirement: Character skip table coverage
The badCharSkip table SHALL handle all 256 possible byte values.

#### Scenario: Character not in pattern
- **WHEN** search encounters character absent from pattern
- **THEN** badCharSkip for that character equals full pattern length

#### Scenario: Last character has no skip
- **WHEN** pattern's last character found in text
- **THEN** badCharSkip does NOT allow zero-distance skip (prevents missing matches)

### Requirement: Search performance characteristics
The Boyer-Moore implementation SHALL achieve sublinear performance on average.

#### Scenario: Large text with rare pattern
- **WHEN** searching for uncommon pattern in large text
- **THEN** system skips large portions of text using badCharSkip table

#### Scenario: Best case skip behavior
- **WHEN** pattern length is 10 and no characters match
- **THEN** system advances by 10 positions per iteration (not 1)

### Requirement: Correctness on edge cases
The string finder SHALL handle boundary conditions correctly.

#### Scenario: Pattern longer than text
- **WHEN** searching for "longpattern" in "tiny"
- **THEN** system returns -1 without out-of-bounds access

#### Scenario: Empty text
- **WHEN** searching for any pattern in empty []byte
- **THEN** system returns -1

#### Scenario: Single character pattern
- **WHEN** pattern is single byte 'x'
- **THEN** search degrades to simple scan with correct behavior

### Requirement: Thread-safe finder usage
A single stringFinder instance SHALL be safely usable by multiple goroutines for searching different texts.

#### Scenario: Concurrent searches with same pattern
- **WHEN** 128 workers share one stringFinder instance
- **THEN** all searches produce correct results without data races

#### Scenario: Read-only access to skip tables
- **WHEN** multiple goroutines call next() method
- **THEN** skip tables are never modified (read-only after initialization)
