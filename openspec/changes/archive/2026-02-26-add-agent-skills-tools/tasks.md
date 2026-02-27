## 1. Dependencies & Setup

- [x] 1.1 Add gopkg.in/yaml.v3 dependency to go.mod
- [x] 1.2 Run go mod tidy to update go.sum

## 2. Models & Data Types

- [x] 2.1 Add SkillMetadata struct to pkg/models/models.go
- [x] 2.2 Add ListSkillsRequest struct to pkg/models/models.go
- [x] 2.3 Add ListSkillsResponse struct to pkg/models/models.go
- [x] 2.4 Add GetSkillRequest struct to pkg/models/models.go
- [x] 2.5 Add GetSkillResponse struct to pkg/models/models.go

## 3. Core Implementation - Skills Tool

- [x] 3.1 Create pkg/tools/skills.go with skill discovery helpers
- [x] 3.2 Implement parseSkillMetadata() to extract YAML frontmatter with strict validation
- [x] 3.3 Implement validateSkillName() to check name format per spec (lowercase, hyphens, length, no leading/trailing/consecutive hyphens)
- [x] 3.4 Implement validateSkillMetadata() to check required/optional field constraints
- [x] 3.5 Implement extractBody() to extract markdown content after closing ---
- [x] 3.6 Implement HandleListSkills() to scan ~/.mcpclaw/skills/, validate each skill, return valid ones with warning logs for invalid ones
- [x] 3.7 Implement HandleGetSkill() to retrieve specific skill, return error if not found or malformed, include full body + metadata

## 4. Integration

- [x] 4.1 Register list_skills tool in internal/server.go registerTools()
- [x] 4.2 Register get_skill tool in internal/server.go registerTools()
- [x] 4.3 Update tool count log message in NewServer() from 20 to 22

## 5. Testing

- [x] 5.1 Create pkg/tools/skills_test.go with unit tests for parseSkillMetadata()
- [x] 5.2 Create unit tests for validateSkillName() covering all format rules
- [x] 5.3 Create unit tests for validateSkillMetadata() covering required/optional constraints
- [x] 5.4 Create integration test for ListSkillsFromDir() with valid and invalid skills
- [x] 5.5 Create integration test for ListSkillsFromDir() with empty directory
- [x] 5.6 Create unit test for readSkillFile() with valid skill
- [x] 5.7 Create unit test for readSkillFile() with malformed skill
- [x] 5.8 Create unit test for readSkillFile() with directory name mismatch
- [x] 5.9 Create unit test for HandleGetSkill() with empty name
- [x] 5.10 Verify all tests pass with go test ./...

## 6. Verification

- [x] 6.1 Verify list_skills returns skill_directory path for each skill (confirmed in tests)
- [x] 6.2 Verify get_skill returns skill_directory path (confirmed in tests)
- [x] 6.3 Verify invalid skills are skipped and warning logged (confirmed in tests)
- [x] 6.4 Verify error messages are descriptive for get_skill with missing/malformed skills (confirmed in code)
- [x] 6.5 Manual test: create sample skill in ~/.mcpclaw/skills/ and verify tools can discover/retrieve it
