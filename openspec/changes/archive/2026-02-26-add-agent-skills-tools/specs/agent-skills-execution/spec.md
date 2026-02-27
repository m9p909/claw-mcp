## ADDED Requirements

### Requirement: Agents can invoke skill scripts using skill directory path
The system SHALL support agents executing scripts from skill packages by providing the skill_directory path, allowing agents to compose bash commands using skill scripts.

#### Scenario: Agent invokes a skill script via bash
- **WHEN** an agent receives a skill_directory path from list_skills or get_skill, and uses that path in a bash command (e.g., `exec_command` with command "python /home/user/.mcpclaw/skills/pdf-processing/scripts/extract.py")
- **THEN** the bash execution tool runs the script and returns stdout/stderr/exit code to the agent

#### Scenario: Agent reads skill reference files
- **WHEN** an agent receives a skill_directory path and uses `read_file` to access scripts/, references/, or assets/ subdirectories
- **THEN** the file reading tool returns the contents of the referenced file, allowing agents to explore skill resources

#### Scenario: Script execution with arguments
- **WHEN** an agent uses the skill_directory path to execute a skill script with arguments (e.g., `/home/user/.mcpclaw/skills/pdf-processing/scripts/extract.py --input file.pdf`)
- **THEN** the script runs with the provided arguments and returns output to the agent

### Requirement: Skill directory path is accessible to agents
The system SHALL return skill_directory in both list_skills and get_skill responses so agents can reference it for subsequent operations.

#### Scenario: list_skills includes directory path
- **WHEN** an agent calls list_skills
- **THEN** each skill in the response includes a skill_directory field containing the absolute path to the skill directory

#### Scenario: get_skill includes directory path
- **WHEN** an agent calls get_skill with a skill name
- **THEN** the response includes a skill_directory field containing the absolute path to the skill directory

### Requirement: Agents can discover script availability
The system SHALL provide enough information for agents to determine what scripts and resources a skill provides.

#### Scenario: Agent explores skill structure
- **WHEN** an agent receives a skill_directory path, it can use read_file to list directory contents or check for scripts/references/assets subdirectories
- **THEN** the agent discovers what resources the skill provides (scripts, references, assets) and can invoke them appropriately
