package storage

import (
	"fmt"
	"strings"

	"awesomeProject/pkg/models"
)

const (
	CategoryFact       = "fact"
	CategoryTodo       = "todo"
	CategoryDecision   = "decision"
	CategoryPreference = "preference"
)

var validCategories = map[string]bool{
	CategoryFact:       true,
	CategoryTodo:       true,
	CategoryDecision:   true,
	CategoryPreference: true,
}

func WriteMemory(category, content string) error {
	if !validCategories[category] {
		return fmt.Errorf("invalid category: %s", category)
	}

	query := `INSERT INTO memories (category, content) VALUES (?, ?)`
	_, err := globalDB.Exec(query, category, content)
	return err
}

func QueryMemory(sqlQuery string) ([]map[string]interface{}, error) {
	// Ensure query is SELECT only (prevent mutations)
	trimmed := strings.TrimSpace(strings.ToUpper(sqlQuery))
	if !strings.HasPrefix(trimmed, "SELECT") {
		return nil, fmt.Errorf("only SELECT queries allowed")
	}

	rows, err := globalDB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(cols))
		valuePtrs := make([]interface{}, len(cols))
		for i := range cols {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range cols {
			entry[col] = values[i]
		}
		results = append(results, entry)
	}

	return results, rows.Err()
}

func SearchMemory(query string, limit int) ([]models.MemoryResult, error) {
	searchLower := strings.ToLower(query)

	sqlQuery := `SELECT id, category, content, created_at FROM memories ORDER BY created_at DESC`
	if limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := globalDB.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.MemoryResult
	for rows.Next() {
		var id int
		var category, content, createdAt string
		if err := rows.Scan(&id, &category, &content, &createdAt); err != nil {
			return nil, err
		}

		contentLower := strings.ToLower(content)
		if strings.Contains(contentLower, searchLower) {
			// Exact match takes precedence
			match := query
			if !strings.Contains(content, query) {
				// Fall back to case-insensitive match
				match = strings.ToLower(query)
			}

			results = append(results, models.MemoryResult{
				ID:        id,
				Category:  category,
				Content:   content,
				CreatedAt: createdAt,
				Match:     match,
			})
		}
	}

	return results, rows.Err()
}
