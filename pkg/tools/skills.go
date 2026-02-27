package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gopkg.in/yaml.v3"

	pkglog "awesomeProject/pkg/log"
	"awesomeProject/pkg/models"
)

// skillFrontmatter represents the parsed YAML frontmatter from SKILL.md
type skillFrontmatter struct {
	Name          string            `yaml:"name"`
	Description   string            `yaml:"description"`
	License       string            `yaml:"license,omitempty"`
	Compatibility string            `yaml:"compatibility,omitempty"`
	AllowedTools  string            `yaml:"allowed-tools,omitempty"`
	Metadata      map[string]string `yaml:"metadata,omitempty"`
}

// parseSkillMetadata extracts and validates YAML frontmatter from SKILL.md content.
// Returns frontmatter, body, and error. On validation error, returns nil, "", error.
func parseSkillMetadata(content string) (*skillFrontmatter, string, error) {
	// Check for opening ---
	if !strings.HasPrefix(content, "---\n") {
		return nil, "", fmt.Errorf("missing opening --- delimiter")
	}

	// Find closing ---
	rest := content[4:] // Skip opening "---\n"
	closingIdx := strings.Index(rest, "\n---\n")
	if closingIdx == -1 {
		return nil, "", fmt.Errorf("missing closing --- delimiter")
	}

	yamlContent := rest[:closingIdx]
	body := rest[closingIdx+5:] // Skip "\n---\n"

	// Parse YAML
	var fm skillFrontmatter
	if err := yaml.Unmarshal([]byte(yamlContent), &fm); err != nil {
		return nil, "", fmt.Errorf("invalid YAML: %w", err)
	}

	return &fm, body, nil
}

// validateSkillName checks name format per Agent Skills spec.
// Name must be: 1-64 chars, lowercase alphanumeric + hyphens, no leading/trailing/consecutive hyphens.
func validateSkillName(name string) error {
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 64 {
		return fmt.Errorf("name exceeds 64 characters")
	}
	if strings.HasPrefix(name, "-") || strings.HasSuffix(name, "-") {
		return fmt.Errorf("name cannot start or end with hyphen")
	}
	if strings.Contains(name, "--") {
		return fmt.Errorf("name cannot contain consecutive hyphens")
	}
	if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(name) {
		return fmt.Errorf("name must contain only lowercase letters, numbers, and hyphens")
	}
	return nil
}

// validateSkillMetadata checks required and optional field constraints.
func validateSkillMetadata(fm *skillFrontmatter) error {
	if fm.Name == "" {
		return fmt.Errorf("name field is required")
	}
	if err := validateSkillName(fm.Name); err != nil {
		return err
	}
	if fm.Description == "" {
		return fmt.Errorf("description field is required")
	}
	if len(fm.Description) > 1024 {
		return fmt.Errorf("description exceeds 1024 characters")
	}
	if fm.Compatibility != "" && len(fm.Compatibility) > 500 {
		return fmt.Errorf("compatibility exceeds 500 characters")
	}
	return nil
}

// getSkillsDir returns the absolute path to ~/.mcpclaw/skills/
func getSkillsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".mcpclaw", "skills"), nil
}

// readSkillFile reads and parses a single SKILL.md file.
// Returns SkillMetadata and body, or error on failure.
func readSkillFile(skillPath string) (*models.SkillMetadata, string, error) {
	skillMDPath := filepath.Join(skillPath, "SKILL.md")
	content, err := os.ReadFile(skillMDPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read SKILL.md: %w", err)
	}

	fm, body, err := parseSkillMetadata(string(content))
	if err != nil {
		return nil, "", err
	}

	if err := validateSkillMetadata(fm); err != nil {
		return nil, "", err
	}

	// Verify skill directory name matches skill name
	skillDirName := filepath.Base(skillPath)
	if skillDirName != fm.Name {
		return nil, "", fmt.Errorf("skill name '%s' does not match directory name '%s'", fm.Name, skillDirName)
	}

	metadata := &models.SkillMetadata{
		Name:           fm.Name,
		Description:    fm.Description,
		License:        fm.License,
		Compatibility:  fm.Compatibility,
		AllowedTools:   fm.AllowedTools,
		Metadata:       fm.Metadata,
		SkillDirectory: skillPath,
	}

	return metadata, body, nil
}

// successResultSkills converts a response to MCP CallToolResult
func successResultSkills(ctx context.Context, resp interface{}) *mcp.CallToolResult {
	jsonData, err := json.Marshal(resp)
	if err != nil {
		return errorResult(ctx, "INTERNAL_ERROR", "failed to marshal response: "+err.Error())
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonData)},
		},
	}
}

// listSkillsFromDir is the internal implementation that accepts a skills directory path.
// This allows testing with temporary directories.
func listSkillsFromDir(ctx context.Context, skillsDir string) ([]models.SkillMetadata, error) {
	logger := pkglog.NewLogger()

	// Check if directory exists
	if _, err := os.Stat(skillsDir); os.IsNotExist(err) {
		// Directory doesn't exist, return empty list (graceful)
		return []models.SkillMetadata{}, nil
	}

	// Read directory
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skills directory: %v", err)
	}

	var validSkills []models.SkillMetadata

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // Skip non-directories
		}

		skillPath := filepath.Join(skillsDir, entry.Name())
		metadata, _, err := readSkillFile(skillPath)
		if err != nil {
			logger.Warn(ctx, fmt.Sprintf("skipped invalid skill '%s': %v", entry.Name(), err))
			continue
		}

		validSkills = append(validSkills, *metadata)
	}

	return validSkills, nil
}

// HandleListSkills scans ~/.mcpclaw/skills/ and returns metadata for all valid skills.
// Invalid skills are skipped with warning logs.
func HandleListSkills(ctx context.Context, req *mcp.CallToolRequest, input models.ListSkillsRequest) (*mcp.CallToolResult, models.ListSkillsResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	skillsDir, err := getSkillsDir()
	if err != nil {
		return errorResult(ctx, "SKILLS_ERROR", fmt.Sprintf("failed to determine skills directory: %v", err)), models.ListSkillsResponse{}, nil
	}

	logger.Info(ctx, "Listing skills", "directory", skillsDir)

	validSkills, err := listSkillsFromDir(ctx, skillsDir)
	if err != nil {
		return errorResult(ctx, "SKILLS_ERROR", fmt.Sprintf("failed to list skills: %v", err)), models.ListSkillsResponse{}, nil
	}

	resp := models.ListSkillsResponse{
		Skills: validSkills,
	}

	logger.Info(ctx, "Listed skills",
		"total_found", len(validSkills),
		pkglog.Duration(time.Since(start)))

	return nil, resp, nil
}

// HandleGetSkill retrieves a specific skill by name and returns its full content.
func HandleGetSkill(ctx context.Context, req *mcp.CallToolRequest, input models.GetSkillRequest) (*mcp.CallToolResult, models.GetSkillResponse, error) {
	logger := pkglog.NewLogger()
	start := time.Now()

	if input.Name == "" {
		return errorResult(ctx, "INVALID_REQUEST", "skill name is required"), models.GetSkillResponse{}, nil
	}

	logger.Info(ctx, "Getting skill", "name", input.Name)

	skillsDir, err := getSkillsDir()
	if err != nil {
		return errorResult(ctx, "SKILLS_ERROR", fmt.Sprintf("failed to determine skills directory: %v", err)), models.GetSkillResponse{}, nil
	}

	skillPath := filepath.Join(skillsDir, input.Name)

	// Check if directory exists
	if _, err := os.Stat(skillPath); os.IsNotExist(err) {
		msg := fmt.Sprintf("skill '%s' not found in %s", input.Name, skillsDir)
		return errorResult(ctx, "SKILL_NOT_FOUND", msg), models.GetSkillResponse{}, nil
	}

	// Read and parse skill
	metadata, body, err := readSkillFile(skillPath)
	if err != nil {
		msg := fmt.Sprintf("skill '%s' is malformed: %v", input.Name, err)
		return errorResult(ctx, "INVALID_SKILL", msg), models.GetSkillResponse{}, nil
	}

	resp := models.GetSkillResponse{
		Name:           metadata.Name,
		Description:    metadata.Description,
		License:        metadata.License,
		Compatibility:  metadata.Compatibility,
		AllowedTools:   metadata.AllowedTools,
		Metadata:       metadata.Metadata,
		SkillDirectory: metadata.SkillDirectory,
		Body:           body,
	}

	logger.Info(ctx, "Retrieved skill",
		"name", input.Name,
		"body_size", len(body),
		pkglog.Duration(time.Since(start)))

	return nil, resp, nil
}
