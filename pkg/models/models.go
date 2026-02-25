package models

// Filesystem Tool Models

type ReadFileRequest struct {
	Path string `json:"path" jsonschema:"description,Absolute file path to read"`
}

type ReadFileResponse struct {
	Content string `json:"content" jsonschema:"description,File content with hashes (format: linenum:hash|content)"`
}

type WriteFileRequest struct {
	Path    string `json:"path" jsonschema:"description,Absolute file path to write"`
	Content string `json:"content" jsonschema:"description,File content (can include hashes from previous read)"`
}

type WriteFileResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether write succeeded"`
	Message string `json:"message" jsonschema:"description,Status message"`
}

type EditFileRequest struct {
	Path       string `json:"path" jsonschema:"description,Absolute file path to edit"`
	StartHash  string `json:"start_hash" jsonschema:"description,Hash of first line to replace"`
	EndHash    string `json:"end_hash" jsonschema:"description,Hash of last line to replace"`
	NewContent string `json:"new_content" jsonschema:"description,Replacement content (can span multiple lines)"`
}

type EditFileResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether edit succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// Process Execution Models

type ExecCommandRequest struct {
	Command    string            `json:"command" jsonschema:"description,Command to execute"`
	Args       []string          `json:"args" jsonschema:"description,Command arguments"`
	Background bool              `json:"background" jsonschema:"description,Run in background and return immediately"`
	Timeout    int               `json:"timeout" jsonschema:"description,Timeout in seconds (0 = no timeout)"`
	Env        map[string]string `json:"env" jsonschema:"description,Environment variables to set"`
}

type ExecCommandResponse struct {
	SessionID  string `json:"session_id" jsonschema:"description,Session ID for background processes"`
	Stdout     string `json:"stdout" jsonschema:"description,Standard output"`
	Stderr     string `json:"stderr" jsonschema:"description,Standard error"`
	ExitCode   int    `json:"exit_code" jsonschema:"description,Process exit code"`
	Status     string `json:"status" jsonschema:"description,running or completed"`
}

type ManageProcessRequest struct {
	Action    string `json:"action" jsonschema:"description,list, poll, send_keys, or kill"`
	SessionID string `json:"session_id" jsonschema:"description,Session ID (required for poll/send_keys/kill)"`
	Keys      string `json:"keys" jsonschema:"description,Keys to send (for send_keys action)"`
}

type ManageProcessResponse struct {
	Sessions []ProcessSession `json:"sessions" jsonschema:"description,List of process sessions"`
	Message  string           `json:"message" jsonschema:"description,Status message"`
}

type ProcessSession struct {
	SessionID string `json:"session_id" jsonschema:"description,Unique session ID"`
	Command   string `json:"command" jsonschema:"description,Command executed"`
	Status    string `json:"status" jsonschema:"description,running or completed"`
	ExitCode  int    `json:"exit_code" jsonschema:"description,Exit code (0 if still running)"`
	Stdout    string `json:"stdout" jsonschema:"description,Captured stdout"`
	Stderr    string `json:"stderr" jsonschema:"description,Captured stderr"`
}

// Memory Persistence Models

type WriteMemoryRequest struct {
	Category string `json:"category" jsonschema:"description,Category: fact, todo, decision, or preference"`
	Content  string `json:"content" jsonschema:"description,Memory content to store"`
}

type WriteMemoryResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether write succeeded"`
	Message string `json:"message" jsonschema:"description,Status message"`
}

type QueryMemoryRequest struct {
	Query string `json:"query" jsonschema:"description,SQL SELECT query (no mutations allowed)"`
}

type QueryMemoryResponse struct {
	Results []map[string]interface{} `json:"results" jsonschema:"description,Query results"`
	Message string                   `json:"message" jsonschema:"description,Status or error message"`
}

type SearchMemoryRequest struct {
	Query string `json:"query" jsonschema:"description,Substring to search for (case-insensitive)"`
	Limit int    `json:"limit" jsonschema:"description,Max results to return (0 = no limit)"`
}

type SearchMemoryResponse struct {
	Results []MemoryResult `json:"results" jsonschema:"description,Search results"`
	Message string         `json:"message" jsonschema:"description,Status message"`
}

type MemoryResult struct {
	ID        int    `json:"id" jsonschema:"description,Memory record ID"`
	Category  string `json:"category" jsonschema:"description,Memory category"`
	Content   string `json:"content" jsonschema:"description,Memory content"`
	CreatedAt string `json:"created_at" jsonschema:"description,Creation timestamp"`
	Match     string `json:"match" jsonschema:"description,Matched substring (exact or partial)"`
}

// Error Response

type ErrorResponse struct {
	Code    string `json:"code" jsonschema:"description,Error code"`
	Message string `json:"message" jsonschema:"description,Error message"`
}
