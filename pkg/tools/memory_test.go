package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"awesomeProject/pkg/models"
	"awesomeProject/pkg/storage"
)

// Test WriteMemory: Valid category and content
func TestHandleWriteMemory_ValidWrite(t *testing.T) {
	// Clear memory for clean test
	storage.ClearMemory()

	input := models.WriteMemoryRequest{
		Category: "fact",
		Content:  "PostgreSQL is a relational database",
	}
	toolResult, resp, err := HandleWriteMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if !resp.Success {
		t.Errorf("expected success: %s", resp.Message)
	}
}

// Test WriteMemory: Multiple categories
func TestHandleWriteMemory_MultipleCategories(t *testing.T) {
	storage.ClearMemory()

	categories := []string{"fact", "todo", "decision", "preference"}
	contents := []string{"Memory content 1", "Memory content 2", "Memory content 3", "Memory content 4"}

	for i, cat := range categories {
		input := models.WriteMemoryRequest{
			Category: cat,
			Content:  contents[i],
		}
		toolResult, resp, err := HandleWriteMemory(context.Background(), nil, input)

		if err != nil {
			t.Fatalf("failed to write category %s: %v", cat, err)
		}
		if toolResult != nil {
			t.Fatalf("expected nil error for category %s", cat)
		}
		if !resp.Success {
			t.Errorf("write failed for category %s", cat)
		}
	}
}

// Test WriteMemory: Empty category
func TestHandleWriteMemory_EmptyCategory(t *testing.T) {
	input := models.WriteMemoryRequest{
		Category: "",
		Content:  "Some content",
	}
	toolResult, _, err := HandleWriteMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for empty category")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_REQUEST") {
		t.Errorf("expected INVALID_REQUEST error")
	}
}

// Test WriteMemory: Empty content
func TestHandleWriteMemory_EmptyContent(t *testing.T) {
	input := models.WriteMemoryRequest{
		Category: "fact",
		Content:  "",
	}
	toolResult, _, err := HandleWriteMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for empty content")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_REQUEST") {
		t.Errorf("expected INVALID_REQUEST error")
	}
}

// Test WriteMemory: Invalid category
func TestHandleWriteMemory_InvalidCategory(t *testing.T) {
	input := models.WriteMemoryRequest{
		Category: "invalid_category",
		Content:  "Some content",
	}
	toolResult, _, err := HandleWriteMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for invalid category")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INTERNAL_ERROR") {
		t.Errorf("expected INTERNAL_ERROR for invalid category, got: %s", errorText)
	}
}

// Test QueryMemory: Simple SELECT query
func TestHandleQueryMemory_SimpleSelect(t *testing.T) {
	storage.ClearMemory()

	// Write test data
	storage.WriteMemory("fact", "Database fact 1")
	storage.WriteMemory("fact", "Database fact 2")

	input := models.QueryMemoryRequest{
		Query: "SELECT * FROM memories WHERE category = 'fact'",
	}
	toolResult, resp, err := HandleQueryMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Results) < 2 {
		t.Errorf("expected at least 2 results, got %d", len(resp.Results))
	}
}

// Test QueryMemory: COUNT query
func TestHandleQueryMemory_CountQuery(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("todo", "Task 1")
	storage.WriteMemory("todo", "Task 2")
	storage.WriteMemory("todo", "Task 3")

	input := models.QueryMemoryRequest{
		Query: "SELECT COUNT(*) as count FROM memories WHERE category = 'todo'",
	}
	toolResult, resp, err := HandleQueryMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp.Results))
	}
}

// Test QueryMemory: Empty query
func TestHandleQueryMemory_EmptyQuery(t *testing.T) {
	input := models.QueryMemoryRequest{
		Query: "",
	}
	toolResult, _, err := HandleQueryMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for empty query")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_REQUEST") {
		t.Errorf("expected INVALID_REQUEST error")
	}
}

// Test QueryMemory: Non-SELECT query (should fail)
func TestHandleQueryMemory_MutationQuery(t *testing.T) {
	input := models.QueryMemoryRequest{
		Query: "INSERT INTO memories (category, content) VALUES ('fact', 'hack')",
	}
	toolResult, _, err := HandleQueryMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for non-SELECT query")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "QUERY_FAILED") {
		t.Errorf("expected QUERY_FAILED error")
	}
}

// Test QueryMemory: UPDATE query (should fail)
func TestHandleQueryMemory_UpdateQuery(t *testing.T) {
	input := models.QueryMemoryRequest{
		Query: "UPDATE memories SET content = 'hacked' WHERE id = 1",
	}
	toolResult, _, err := HandleQueryMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for UPDATE query")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "QUERY_FAILED") {
		t.Errorf("expected QUERY_FAILED error")
	}
}

// Test SearchMemory: Substring match (case-insensitive)
func TestHandleSearchMemory_SubstringMatch(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("fact", "PostgreSQL database")
	storage.WriteMemory("fact", "MySQL database")
	storage.WriteMemory("todo", "Learn Redis")

	input := models.SearchMemoryRequest{
		Query: "database",
		Limit: 0,
	}
	toolResult, resp, err := HandleMemorySearch(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Results) < 2 {
		t.Errorf("expected at least 2 results, got %d", len(resp.Results))
	}
}

// Test SearchMemory: Case-insensitive search
func TestHandleSearchMemory_CaseInsensitive(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("fact", "POSTGRESQL is great")
	storage.WriteMemory("fact", "postgresql is reliable")

	input := models.SearchMemoryRequest{
		Query: "postgresql",
		Limit: 0,
	}
	toolResult, resp, err := HandleMemorySearch(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(resp.Results))
	}
}

// Test SearchMemory: With limit
func TestHandleSearchMemory_WithLimit(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("decision", "Decision 1")
	storage.WriteMemory("decision", "Decision 2")
	storage.WriteMemory("decision", "Decision 3")
	storage.WriteMemory("decision", "Decision 4")

	input := models.SearchMemoryRequest{
		Query: "Decision",
		Limit: 2,
	}
	toolResult, resp, err := HandleMemorySearch(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Results) != 2 {
		t.Errorf("expected 2 results (limit), got %d", len(resp.Results))
	}
}

// Test SearchMemory: No matches
func TestHandleSearchMemory_NoMatches(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("fact", "PostgreSQL")
	storage.WriteMemory("fact", "MySQL")

	input := models.SearchMemoryRequest{
		Query: "NonExistentTerm",
		Limit: 0,
	}
	toolResult, resp, err := HandleMemorySearch(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Results) != 0 {
		t.Errorf("expected 0 results, got %d", len(resp.Results))
	}
}

// Test SearchMemory: Empty query
func TestHandleSearchMemory_EmptyQuery(t *testing.T) {
	input := models.SearchMemoryRequest{
		Query: "",
		Limit: 0,
	}
	toolResult, _, err := HandleMemorySearch(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult == nil {
		t.Fatalf("expected error result for empty query")
	}

	errorText := toolResult.Content[0].(*mcp.TextContent).Text
	if !strings.Contains(errorText, "INVALID_REQUEST") {
		t.Errorf("expected INVALID_REQUEST error")
	}
}

// Integration Test: Write → Query → Search workflow
func TestIntegration_MemoryWorkflow(t *testing.T) {
	storage.ClearMemory()

	// Write multiple memories
	categories := []string{"fact", "fact", "todo", "decision", "preference"}
	contents := []string{
		"PostgreSQL supports JSON",
		"MySQL is widely used",
		"Implement caching layer",
		"Use async processing",
		"Prefer structured logging",
	}

	for i, cat := range categories {
		input := models.WriteMemoryRequest{
			Category: cat,
			Content:  contents[i],
		}
		toolResult, _, _ := HandleWriteMemory(context.Background(), nil, input)
		if toolResult != nil {
			t.Fatalf("failed to write memory %d", i)
		}
	}

	// Query for all facts
	queryInput := models.QueryMemoryRequest{
		Query: "SELECT category, content FROM memories WHERE category = 'fact'",
	}
	toolResult, queryResp, _ := HandleQueryMemory(context.Background(), nil, queryInput)
	if toolResult != nil {
		t.Fatalf("query failed")
	}
	if len(queryResp.Results) != 2 {
		t.Errorf("expected 2 facts, got %d", len(queryResp.Results))
	}

	// Search for 'processing'
	searchInput := models.SearchMemoryRequest{
		Query: "processing",
		Limit: 0,
	}
	toolResult, searchResp, _ := HandleMemorySearch(context.Background(), nil, searchInput)
	if toolResult != nil {
		t.Fatalf("search failed")
	}
	if len(searchResp.Results) != 1 {
		t.Errorf("expected 1 match for 'processing', got %d", len(searchResp.Results))
	}

	// Search for 'SQL' (case-insensitive)
	searchInput.Query = "SQL"
	toolResult, searchResp, _ = HandleMemorySearch(context.Background(), nil, searchInput)
	if len(searchResp.Results) != 2 {
		t.Errorf("expected 2 matches for 'SQL' (case-insensitive), got %d", len(searchResp.Results))
	}
}

// Test SearchMemory: Special characters in query
func TestHandleSearchMemory_SpecialCharacters(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("fact", "Use @mention for notifications")
	storage.WriteMemory("fact", "Email: test@example.com")

	input := models.SearchMemoryRequest{
		Query: "@",
		Limit: 0,
	}
	toolResult, resp, err := HandleMemorySearch(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Results) < 2 {
		t.Errorf("expected at least 2 results for '@', got %d", len(resp.Results))
	}
}

// Test QueryMemory: Multiple WHERE conditions
func TestHandleQueryMemory_MultipleConditions(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("fact", "PostgreSQL")
	storage.WriteMemory("fact", "MySQL")
	storage.WriteMemory("todo", "Learn PostgreSQL")

	input := models.QueryMemoryRequest{
		Query: "SELECT * FROM memories WHERE category = 'fact' AND content LIKE '%SQL%'",
	}
	toolResult, resp, err := HandleQueryMemory(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}
	if len(resp.Results) < 2 {
		t.Errorf("expected at least 2 results, got %d", len(resp.Results))
	}
}

// Test SearchMemory: Exact vs case-insensitive match
func TestHandleSearchMemory_ExactMatchPriority(t *testing.T) {
	storage.ClearMemory()

	storage.WriteMemory("fact", "PostgreSQL is great")
	storage.WriteMemory("fact", "postgresql is reliable")

	input := models.SearchMemoryRequest{
		Query: "postgresql",
		Limit: 10,
	}
	toolResult, resp, err := HandleMemorySearch(context.Background(), nil, input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if toolResult != nil {
		t.Fatalf("expected nil error result")
	}

	// Results should include both exact and case-insensitive matches
	for _, result := range resp.Results {
		if result.Match == "" {
			t.Errorf("result should have a match field")
		}
	}
}
