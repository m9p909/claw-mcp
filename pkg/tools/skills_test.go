package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"awesomeProject/pkg/models"
)

// Test helpers

func createTempSkill(t *testing.T, skillName string, skillmdContent string) string {
	tmpDir := t.TempDir()
	skillDir := filepath.Join(tmpDir, skillName)
	if err := os.Mkdir(skillDir, 0755); err != nil {
		t.Fatalf("failed to create skill dir: %v", err)
	}

	skillMDPath := filepath.Join(skillDir, "SKILL.md")
	if err := os.WriteFile(skillMDPath, []byte(skillmdContent), 0644); err != nil {
		t.Fatalf("failed to write SKILL.md: %v", err)
	}

	return tmpDir
}

// Helper to get skill directory for testing
func getTempSkillDir(skillsDir string, skillName string) string {
	return filepath.Join(skillsDir, skillName)
}

// Test parseSkillMetadata

func TestParseSkillMetadata_Valid(t *testing.T) {
	content := `---
name: pdf-processing
description: Extract and process PDFs
version: 1.0
license: MIT
---
# My Skill

This is the body.
`
	fm, body, err := parseSkillMetadata(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fm.Name != "pdf-processing" {
		t.Errorf("expected name 'pdf-processing', got %q", fm.Name)
	}
	if fm.Description != "Extract and process PDFs" {
		t.Errorf("expected description 'Extract and process PDFs', got %q", fm.Description)
	}
	if !strings.HasPrefix(body, "# My Skill") {
		t.Errorf("expected body to start with '# My Skill', got %q", body)
	}
}

func TestParseSkillMetadata_MissingOpeningDelimiter(t *testing.T) {
	content := `name: test
---
body
`
	_, _, err := parseSkillMetadata(content)
	if err == nil {
		t.Fatal("expected error for missing opening ---")
	}
	if !strings.Contains(err.Error(), "opening") {
		t.Errorf("expected error about opening delimiter, got: %v", err)
	}
}

func TestParseSkillMetadata_MissingClosingDelimiter(t *testing.T) {
	content := `---
name: test
description: test desc
`
	_, _, err := parseSkillMetadata(content)
	if err == nil {
		t.Fatal("expected error for missing closing ---")
	}
	if !strings.Contains(err.Error(), "closing") {
		t.Errorf("expected error about closing delimiter, got: %v", err)
	}
}

func TestParseSkillMetadata_InvalidYAML(t *testing.T) {
	content := `---
name: test
description: "unclosed string
---
body
`
	_, _, err := parseSkillMetadata(content)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
	if !strings.Contains(err.Error(), "YAML") {
		t.Errorf("expected error about YAML, got: %v", err)
	}
}

// Test validateSkillName

func TestValidateSkillName_Valid(t *testing.T) {
	tests := []string{
		"pdf-processing",
		"data-analysis",
		"test123",
		"a",
		"abc-def-ghi",
	}
	for _, name := range tests {
		if err := validateSkillName(name); err != nil {
			t.Errorf("validateSkillName(%q) returned error: %v", name, err)
		}
	}
}

func TestValidateSkillName_InvalidLeadingHyphen(t *testing.T) {
	err := validateSkillName("-test")
	if err == nil {
		t.Fatal("expected error for leading hyphen")
	}
	if !strings.Contains(err.Error(), "start") {
		t.Errorf("expected error about start, got: %v", err)
	}
}

func TestValidateSkillName_InvalidTrailingHyphen(t *testing.T) {
	err := validateSkillName("test-")
	if err == nil {
		t.Fatal("expected error for trailing hyphen")
	}
	if !strings.Contains(err.Error(), "end") {
		t.Errorf("expected error about end, got: %v", err)
	}
}

func TestValidateSkillName_InvalidConsecutiveHyphens(t *testing.T) {
	err := validateSkillName("test--name")
	if err == nil {
		t.Fatal("expected error for consecutive hyphens")
	}
	if !strings.Contains(err.Error(), "consecutive") {
		t.Errorf("expected error about consecutive, got: %v", err)
	}
}

func TestValidateSkillName_InvalidChars(t *testing.T) {
	tests := []string{
		"Test",      // uppercase
		"test name", // space
		"test_name", // underscore
		"test.name", // dot
	}
	for _, name := range tests {
		err := validateSkillName(name)
		if err == nil {
			t.Errorf("validateSkillName(%q) expected error, got nil", name)
		}
	}
}

func TestValidateSkillName_TooLong(t *testing.T) {
	longName := strings.Repeat("a", 65)
	err := validateSkillName(longName)
	if err == nil {
		t.Fatal("expected error for name > 64 chars")
	}
}

func TestValidateSkillName_Empty(t *testing.T) {
	err := validateSkillName("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

// Test validateSkillMetadata

func TestValidateSkillMetadata_Valid(t *testing.T) {
	fm := &skillFrontmatter{
		Name:        "test-skill",
		Description: "A test skill",
	}
	if err := validateSkillMetadata(fm); err != nil {
		t.Errorf("validateSkillMetadata returned error: %v", err)
	}
}

func TestValidateSkillMetadata_MissingName(t *testing.T) {
	fm := &skillFrontmatter{
		Description: "A test skill",
	}
	err := validateSkillMetadata(fm)
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestValidateSkillMetadata_MissingDescription(t *testing.T) {
	fm := &skillFrontmatter{
		Name: "test-skill",
	}
	err := validateSkillMetadata(fm)
	if err == nil {
		t.Fatal("expected error for missing description")
	}
}

func TestValidateSkillMetadata_DescriptionTooLong(t *testing.T) {
	fm := &skillFrontmatter{
		Name:        "test-skill",
		Description: strings.Repeat("a", 1025),
	}
	err := validateSkillMetadata(fm)
	if err == nil {
		t.Fatal("expected error for description > 1024 chars")
	}
}

func TestValidateSkillMetadata_CompatibilityTooLong(t *testing.T) {
	fm := &skillFrontmatter{
		Name:          "test-skill",
		Description:   "Test",
		Compatibility: strings.Repeat("a", 501),
	}
	err := validateSkillMetadata(fm)
	if err == nil {
		t.Fatal("expected error for compatibility > 500 chars")
	}
}

// Test listSkillsFromDir helper

func TestListSkillsFromDir_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	skills, err := listSkillsFromDir(context.Background(), tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 0 {
		t.Errorf("expected empty skills list, got %d skills", len(skills))
	}
}

func TestListSkillsFromDir_ValidSkills(t *testing.T) {
	skillsMD := `---
name: test-skill
description: A test skill
---
# Test Skill

Body content.
`
	tmpDir := createTempSkill(t, "test-skill", skillsMD)
	defer os.RemoveAll(tmpDir)

	skills, err := listSkillsFromDir(context.Background(), tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}
	if skills[0].Name != "test-skill" {
		t.Errorf("expected name 'test-skill', got %q", skills[0].Name)
	}
	if skills[0].Description != "A test skill" {
		t.Errorf("expected description 'A test skill', got %q", skills[0].Description)
	}
}

func TestListSkillsFromDir_SkipsInvalidSkills(t *testing.T) {
	// Create one valid and one invalid skill
	validSkillMD := `---
name: valid-skill
description: Valid skill
---
Body
`
	invalidSkillMD := `---
name: invalid-skill
description: Invalid skill
# Missing closing ---
`

	tmpDir := t.TempDir()

	// Create valid skill
	validDir := filepath.Join(tmpDir, "valid-skill")
	os.Mkdir(validDir, 0755)
	os.WriteFile(filepath.Join(validDir, "SKILL.md"), []byte(validSkillMD), 0644)

	// Create invalid skill
	invalidDir := filepath.Join(tmpDir, "invalid-skill")
	os.Mkdir(invalidDir, 0755)
	os.WriteFile(filepath.Join(invalidDir, "SKILL.md"), []byte(invalidSkillMD), 0644)

	skills, err := listSkillsFromDir(context.Background(), tmpDir)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should only return valid skill, skip invalid one with warning
	if len(skills) != 1 {
		t.Errorf("expected 1 valid skill, got %d", len(skills))
	}
	if skills[0].Name != "valid-skill" {
		t.Errorf("expected 'valid-skill', got %q", skills[0].Name)
	}
}

// Test readSkillFile helper

func TestReadSkillFile_Valid(t *testing.T) {
	skillsMD := `---
name: pdf-skill
description: PDF processing
version: 1.0
license: MIT
---
# PDF Skill

Process PDFs here.
`
	skillName := "pdf-skill"
	tmpDir := createTempSkill(t, skillName, skillsMD)
	defer os.RemoveAll(tmpDir)

	skillPath := filepath.Join(tmpDir, skillName)
	metadata, body, err := readSkillFile(skillPath)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata.Name != "pdf-skill" {
		t.Errorf("expected name 'pdf-skill', got %q", metadata.Name)
	}
	if metadata.Description != "PDF processing" {
		t.Errorf("expected description 'PDF processing', got %q", metadata.Description)
	}
	if !strings.HasPrefix(body, "# PDF Skill") {
		t.Errorf("expected body to start with '# PDF Skill', got %q", body)
	}
	if metadata.License != "MIT" {
		t.Errorf("expected license 'MIT', got %q", metadata.License)
	}
}

func TestReadSkillFile_Malformed(t *testing.T) {
	malformedMD := `---
name: bad-skill
description: Missing closing delimiter
`
	skillName := "bad-skill"
	tmpDir := createTempSkill(t, skillName, malformedMD)
	defer os.RemoveAll(tmpDir)

	skillPath := filepath.Join(tmpDir, skillName)
	_, _, err := readSkillFile(skillPath)

	if err == nil {
		t.Fatal("expected error for malformed skill")
	}
}

func TestReadSkillFile_DirectoryNameMismatch(t *testing.T) {
	skillsMD := `---
name: different-name
description: Test skill
---
Body
`
	skillName := "actual-name"
	tmpDir := createTempSkill(t, skillName, skillsMD)
	defer os.RemoveAll(tmpDir)

	skillPath := filepath.Join(tmpDir, skillName)
	_, _, err := readSkillFile(skillPath)

	if err == nil {
		t.Fatal("expected error for directory name mismatch")
	}
	if !strings.Contains(err.Error(), "does not match") {
		t.Errorf("expected error about mismatch, got: %v", err)
	}
}

// Test error handling in HandleGetSkill

func TestHandleGetSkill_EmptyName(t *testing.T) {
	input := models.GetSkillRequest{Name: ""}
	toolResult, _, err := HandleGetSkill(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if toolResult == nil {
		t.Fatal("expected error result for empty name")
	}

	if len(toolResult.Content) == 0 {
		t.Fatal("expected content in error result")
	}
}
