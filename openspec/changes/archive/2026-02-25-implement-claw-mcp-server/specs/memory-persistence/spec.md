## ADDED Requirements

### Requirement: Write memories to persistent storage
The system SHALL store memories in SQLite with category, content, and timestamp. Categories are constrained to predefined types.

#### Scenario: Write fact memory
- **WHEN** client calls write_memory with category=fact and content
- **THEN** system stores memory, returns id, category, created_at timestamp

#### Scenario: Write todo memory
- **WHEN** client calls write_memory with category=todo
- **THEN** system stores memory with todo category

#### Scenario: Write decision memory
- **WHEN** client calls write_memory with category=decision
- **THEN** system stores memory with decision category

#### Scenario: Write preference memory
- **WHEN** client calls write_memory with category=preference
- **THEN** system stores memory with preference category

#### Scenario: Invalid category
- **WHEN** client uses category not in {fact, todo, decision, preference}
- **THEN** system returns InvalidInputError

#### Scenario: Missing content
- **WHEN** client omits content field
- **THEN** system returns InvalidInputError

### Requirement: Query memories with SQL
The system SHALL allow clients to execute SELECT queries against the memories table.

#### Scenario: Query all memories
- **WHEN** client calls query_memory with sql="SELECT * FROM memories"
- **THEN** system returns all memories as array of objects with id, category, content, created_at

#### Scenario: Filter by category
- **WHEN** client queries with WHERE category='fact'
- **THEN** system returns only memories matching that category

#### Scenario: Apply limit
- **WHEN** client specifies limit parameter
- **THEN** system returns at most that many results

#### Scenario: SQL syntax error
- **WHEN** SQL query is malformed
- **THEN** system returns SQLError with error message

#### Scenario: Prevent data mutation
- **WHEN** client attempts INSERT, UPDATE, or DELETE
- **THEN** system only allows SELECT, rejects mutation attempts

### Requirement: Search memories by substring
The system SHALL search all memories for substring matches (case-insensitive), returning top N results.

#### Scenario: Substring match
- **WHEN** client calls memory_search with query string
- **THEN** system returns all memories containing query (case-insensitive)

#### Scenario: Exact match priority
- **WHEN** multiple memories match
- **THEN** exact matches appear before partial matches in results

#### Scenario: Apply search limit
- **WHEN** client specifies limit
- **THEN** system returns at most that many results

#### Scenario: No matches
- **WHEN** query doesn't match any memory
- **THEN** system returns empty results with success=true
