## ADDED Requirements

### Requirement: List all valid skills from the skills directory
The system SHALL scan the `~/.mcpclaw/skills/` directory, validate each skill's SKILL.md against the Agent Skills specification, and return a list of valid skill metadata.

#### Scenario: User lists skills when directory is empty
- **WHEN** an agent calls `list_skills` and the `~/.mcpclaw/skills/` directory is empty or does not exist
- **THEN** the system returns an empty skills list (no error, graceful handling)

#### Scenario: User lists skills with valid skills present
- **WHEN** an agent calls `list_skills` and the directory contains valid SKILL.md files
- **THEN** the system returns metadata for each valid skill including name, description, version, license, compatibility, metadata, allowed-tools, and skill_directory path

#### Scenario: User lists skills with malformed skills present
- **WHEN** an agent calls `list_skills` and the directory contains invalid SKILL.md files (missing closing ---,  invalid YAML, missing required fields)
- **THEN** the system skips invalid skills, logs a warning per invalid skill, and returns only valid skills

#### Scenario: Skill name validation - invalid names excluded
- **WHEN** a skill directory has an invalid name (uppercase, leading hyphen, consecutive hyphens, spaces, special chars) or SKILL.md contains mismatched name
- **THEN** the skill is excluded from results with a warning log indicating the validation failure

#### Scenario: Skill description validation - missing or too long
- **WHEN** a skill SKILL.md is missing the required description field or description exceeds 1024 characters
- **THEN** the skill is excluded from results with a warning log

### Requirement: Retrieve a specific skill's full content and metadata
The system SHALL retrieve a single skill by name, validate it, and return complete metadata plus markdown body.

#### Scenario: Get skill that exists and is valid
- **WHEN** an agent calls `get_skill` with a valid skill name
- **THEN** the system returns the skill's name, description, version, license, compatibility, metadata, allowed-tools, skill_directory, and full markdown body (everything after closing ---)

#### Scenario: Get skill that does not exist
- **WHEN** an agent calls `get_skill` with a skill name that does not exist in `~/.mcpclaw/skills/`
- **THEN** the system returns an error result with a message like "skill 'pdf-processing' not found in ~/.mcpclaw/skills/"

#### Scenario: Get skill that exists but is malformed
- **WHEN** an agent calls `get_skill` with a skill name whose SKILL.md is invalid (missing closing ---, invalid YAML, missing required fields)
- **THEN** the system returns an error result describing the validation failure (e.g., "invalid SKILL.md format: missing closing ---")

### Requirement: Validate Agent Skills format per specification
The system SHALL validate all SKILL.md files strictly according to the Agent Skills specification.

#### Scenario: Name field constraints
- **WHEN** validating a skill name
- **THEN** the system checks name is 1-64 chars, contains only lowercase alphanumeric and hyphens, does not start/end/have consecutive hyphens, and matches the parent directory name

#### Scenario: Description field constraints
- **WHEN** validating a skill description
- **THEN** the system checks description is present, 1-1024 chars, and non-empty

#### Scenario: Optional field constraints
- **WHEN** validating optional fields (license, compatibility, metadata, allowed-tools)
- **THEN** the system checks compatibility does not exceed 500 chars (if present) and metadata is a valid map (if present)

#### Scenario: Frontmatter structure
- **WHEN** validating frontmatter
- **THEN** the system requires exact `---\nYAML\n---\n` structure; missing closing --- causes validation failure

### Requirement: Return skill directory path for script access
The system SHALL include the absolute filesystem path to the skill directory in all responses.

#### Scenario: Agents can locate skill scripts
- **WHEN** an agent receives a response from list_skills or get_skill
- **THEN** the skill_directory field contains the absolute path (e.g., /home/user/.mcpclaw/skills/pdf-processing) allowing agents to use bash commands to invoke scripts in scripts/ subdirectory

### Requirement: Log warnings for invalid skills
The system SHALL log a warning message for each invalid skill found during scanning.

#### Scenario: Warning logged for malformed skill
- **WHEN** a skill directory is found but its SKILL.md fails validation
- **THEN** the system logs a warning message indicating the skill name and reason for failure (e.g., "WARNING: skill 'bad-skill' invalid: missing closing ---")
