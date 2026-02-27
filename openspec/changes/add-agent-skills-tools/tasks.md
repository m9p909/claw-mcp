## 1. Dependencies & Setup

- [ ] 1.1 Add gopkg.in/yaml.v3 dependency to go.mod
- [ ] 1.2 Run go mod tidy to update go.sum

## 2. Models & Data Types

- [ ] 2.1 Add SkillMetadata struct to pkg/models/models.go
- [ ] 2.2 Add ListSkillsRequest struct to pkg/models/models.go
- [ ] 2.3 Add ListSkillsResponse struct to pkg/models/models.go
- [ ] 2.4 Add GetSkillRequest struct to pkg/models/models.go
- [ ] 2.5 Add GetSkillResponse struct to pkg/models/models.go

## 3. Core Implementation - Skills Tool

- [ ] 3.1 Create pkg/tools/skills.go with skill discovery helpers
- [ ] 3.2 Implement parseSkillMetadata() to extract YAML frontmatter with strict validation
- [ ] 3.3 Implement validateSkillName() to check name format per spec (lowercase, hyphens, length, no leading/trailing/consecutive hyphens)
- [ ] 3.4 Implement validateSkillMetadata() to check required/optional field constraints
- [ ] 3.5 Implement extractBody() to extract markdown content after closing ---
- [ ] 3.6 Implement HandleListSkills() to scan ~/.mcpclaw/skills/, validate each skill, return valid ones with warning logs for invalid ones
- [ ] 3.7 Implement HandleGetSkill() to retrieve specific skill, return error if not found or malformed, include full body + metadata

## 4. Integration

- [ ] 4.1 Register list_skills tool in internal/server.go registerTools()
- [ ] 4.2 Register get_skill tool in internal/server.go registerTools()
- [ ] 4.3 Update tool count log message in NewServer() from 20 to 22

## 5. Testing

- [ ] 5.1 Create pkg/tools/skills_test.go with unit tests for parseSkillMetadata()
- [ ] 5.2 Create unit tests for validateSkillName() covering all format rules
- [ ] 5.3 Create unit tests for validateSkillMetadata() covering required/optional constraints
- [ ] 5.4 Create integration test for HandleListSkills() with valid and invalid skills
- [ ] 5.5 Create integration test for HandleListSkills() with empty directory
- [ ] 5.6 Create integration test for HandleGetSkill() with valid skill
- [ ] 5.7 Create integration test for HandleGetSkill() with missing skill
- [ ] 5.8 Create integration test for HandleGetSkill() with malformed skill
- [ ] 5.9 Create test fixtures (sample SKILL.md files) in test directory
- [ ] 5.10 Verify all tests pass with go test ./...

## 6. Verification

- [ ] 6.1 Verify list_skills returns skill_directory path for each skill
- [ ] 6.2 Verify get_skill returns skill_directory path
- [ ] 6.3 Verify invalid skills are skipped and warning logged
- [ ] 6.4 Verify error messages are descriptive for get_skill with missing/malformed skills
- [ ] 6.5 Manual test: create sample skill in ~/.mcpclaw/skills/ and verify tools can discover/retrieve it
