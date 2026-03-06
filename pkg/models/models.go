package models

// Filesystem Tool Models

type ReadFileRequest struct {
	Path string `json:"path" jsonschema:"description,Absolute file path to read"`
}

type ReadFileResponse struct {
	Content string `json:"content" jsonschema:"description,File content with line numbers (format: linenum|content)"`
}

type WriteFileRequest struct {
	Path    string `json:"path" jsonschema:"description,Absolute file path to write"`
	Content string `json:"content" jsonschema:"description,File content to write verbatim"`
}

type WriteFileResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether write succeeded"`
	Message string `json:"message" jsonschema:"description,Status message"`
}

type EditFileRequest struct {
	Path       string `json:"path" jsonschema:"description,Absolute file path to edit"`
	StartLine  int    `json:"start_line" jsonschema:"description,1-based line number of first line to replace"`
	EndLine    int    `json:"end_line" jsonschema:"description,1-based line number of last line to replace"`
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
	SessionID string `json:"session_id" jsonschema:"description,Session ID for background processes"`
	Stdout    string `json:"stdout" jsonschema:"description,Standard output"`
	Stderr    string `json:"stderr" jsonschema:"description,Standard error"`
	ExitCode  int    `json:"exit_code" jsonschema:"description,Process exit code"`
	Status    string `json:"status" jsonschema:"description,running or completed"`
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

// Browser Automation Tool Models

type BrowserNavigateRequest struct {
	URL     string `json:"url" jsonschema:"description,URL to navigate to"`
	Timeout int    `json:"timeout" jsonschema:"description,Timeout in seconds (default: 30)"`
}

type BrowserNavigateResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether navigation succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserNavigateBackRequest struct {
	Timeout int `json:"timeout" jsonschema:"description,Timeout in seconds (default: 30)"`
}

type BrowserNavigateBackResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether navigation back succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserSnapshotRequest struct{}

type BrowserSnapshotResponse struct {
	Snapshot string `json:"snapshot" jsonschema:"description,Structured accessibility tree representation"`
}

type BrowserClickRequest struct {
	Ref         string   `json:"ref" jsonschema:"description,Exact target element reference from snapshot"`
	Element     string   `json:"element" jsonschema:"description,Human-readable element description"`
	Button      string   `json:"button" jsonschema:"description,left|right|middle (default: left)"`
	DoubleClick bool     `json:"doubleClick" jsonschema:"description,Whether to perform double-click"`
	Modifiers   []string `json:"modifiers" jsonschema:"description,Modifier keys: Alt, Control, Meta, Shift"`
}

type BrowserClickResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether click succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserHoverRequest struct {
	Ref     string `json:"ref" jsonschema:"description,Exact target element reference from snapshot"`
	Element string `json:"element" jsonschema:"description,Human-readable element description"`
}

type BrowserHoverResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether hover succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserTypeRequest struct {
	Text string `json:"text" jsonschema:"description,Text to type"`
	Ref  string `json:"ref" jsonschema:"description,Element reference to focus first (optional)"`
}

type BrowserTypeResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether typing succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserFormField struct {
	Ref   string `json:"ref" jsonschema:"description,Element reference (required)"`
	Value string `json:"value" jsonschema:"description,Value to fill (required)"`
	Name  string `json:"name" jsonschema:"description,Human-readable field name (optional)"`
}

type BrowserFillFormRequest struct {
	Fields []BrowserFormField `json:"fields" jsonschema:"description,Array of form fields to fill"`
}

type BrowserFillFormResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether fill succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserSelectOptionRequest struct {
	Ref     string   `json:"ref" jsonschema:"description,Select element reference (required)"`
	Values  []string `json:"values" jsonschema:"description,Option values to select (required)"`
	Element string   `json:"element" jsonschema:"description,Human-readable element description"`
}

type BrowserSelectOptionResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether select succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserPressKeyRequest struct {
	Key string `json:"key" jsonschema:"description,Key name (Enter, Tab, Escape, etc.) or single character"`
	Ref string `json:"ref" jsonschema:"description,Element to focus first (optional)"`
}

type BrowserPressKeyResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether key press succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserWaitForRequest struct {
	Text     string  `json:"text" jsonschema:"description,Text to wait for (optional)"`
	TextGone string  `json:"textGone" jsonschema:"description,Text to wait to disappear (optional)"`
	Time     float64 `json:"time" jsonschema:"description,Seconds to wait (optional)"`
}

type BrowserWaitForResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether wait succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserHandleDialogRequest struct {
	Accept     bool   `json:"accept" jsonschema:"description,true to accept/OK, false to cancel"`
	PromptText string `json:"promptText" jsonschema:"description,Text to enter in prompt dialogs (optional)"`
}

type BrowserHandleDialogResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether dialog handled successfully"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

type BrowserCloseRequest struct{}

type BrowserCloseResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether close succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// File Search Tool Models

type SearchFileRequest struct {
	Path    string `json:"path" jsonschema:"description,Absolute or relative path to search"`
	Query   string `json:"query" jsonschema:"description,Search query (literal string or regex pattern)"`
	Regex   bool   `json:"regex" jsonschema:"description,Treat query as regex pattern"`
	Workers int    `json:"workers" jsonschema:"description,Number of worker goroutines (default: 128)"`
}

type SearchFileResult struct {
	File    string `json:"file" jsonschema:"description,File path containing match"`
	Line    int    `json:"line" jsonschema:"description,Line number (0 for binary files)"`
	Content string `json:"content" jsonschema:"description,Matched line content"`
}

type SearchFileResponse struct {
	Results []SearchFileResult `json:"results" jsonschema:"description,Search results with hashes"`
	Message string             `json:"message,omitempty" jsonschema:"description,Status or error message"`
}

type FindFilesRequest struct {
	Path    string `json:"path" jsonschema:"description,Root directory to search"`
	Pattern string `json:"pattern" jsonschema:"description,Glob pattern to match (e.g., *.go, test_*)"`
}

type FindFilesResult struct {
	Path     string `json:"path" jsonschema:"description,Absolute file path"`
	Size     int64  `json:"size" jsonschema:"description,File size in bytes"`
	Modified int64  `json:"modified" jsonschema:"description,Unix timestamp of last modification"`
}

type FindFilesResponse struct {
	Files   []FindFilesResult `json:"files" jsonschema:"description,Matched files"`
	Message string            `json:"message,omitempty" jsonschema:"description,Status or error message"`
}

type ListDirectoryRequest struct {
	Path string `json:"path" jsonschema:"description,Directory path to list"`
}

type ListDirectoryEntry struct {
	Name        string `json:"name" jsonschema:"description,File or directory name"`
	Type        string `json:"type" jsonschema:"description,Entry type: file or dir"`
	Size        int64  `json:"size" jsonschema:"description,Size in bytes (0 for directories)"`
	Permissions string `json:"permissions" jsonschema:"description,File permissions string"`
}

type ListDirectoryResponse struct {
	Entries []ListDirectoryEntry `json:"entries" jsonschema:"description,Directory contents"`
	Message string               `json:"message,omitempty" jsonschema:"description,Status or error message"`
}

type TreeDirectoryRequest struct {
	Path     string `json:"path" jsonschema:"description,Root directory for tree"`
	MaxDepth int    `json:"max_depth" jsonschema:"description,Maximum depth (0 = unlimited)"`
}

type TreeDirectoryResponse struct {
	Tree    string `json:"tree" jsonschema:"description,ASCII tree visualization"`
	Message string `json:"message,omitempty" jsonschema:"description,Status or error message"`
}

// Agent Skills Models

type SkillMetadata struct {
	Name           string            `json:"name" jsonschema:"description,Skill name (lowercase alphanumeric and hyphens)"`
	Description    string            `json:"description" jsonschema:"description,Skill description (1-1024 characters)"`
	License        string            `json:"license,omitempty" jsonschema:"description,License name or reference"`
	Compatibility  string            `json:"compatibility,omitempty" jsonschema:"description,Environment requirements"`
	AllowedTools   string            `json:"allowed_tools,omitempty" jsonschema:"description,Space-delimited list of pre-approved tools"`
	Metadata       map[string]string `json:"metadata,omitempty" jsonschema:"description,Arbitrary key-value metadata"`
	SkillDirectory string            `json:"skill_directory" jsonschema:"description,Absolute path to skill directory"`
}

type ListSkillsRequest struct{}

type ListSkillsResponse struct {
	Skills  []SkillMetadata `json:"skills" jsonschema:"description,List of valid skills"`
	Message string          `json:"message,omitempty" jsonschema:"description,Status message"`
}

type GetSkillRequest struct {
	Name string `json:"name" jsonschema:"description,Name of skill to retrieve"`
}

type GetSkillResponse struct {
	Name           string            `json:"name" jsonschema:"description,Skill name"`
	Description    string            `json:"description" jsonschema:"description,Skill description"`
	License        string            `json:"license,omitempty" jsonschema:"description,License"`
	Compatibility  string            `json:"compatibility,omitempty" jsonschema:"description,Environment requirements"`
	AllowedTools   string            `json:"allowed_tools,omitempty" jsonschema:"description,Space-delimited list of pre-approved tools"`
	Metadata       map[string]string `json:"metadata,omitempty" jsonschema:"description,Arbitrary key-value metadata"`
	SkillDirectory string            `json:"skill_directory" jsonschema:"description,Absolute path to skill directory"`
	Body           string            `json:"body" jsonschema:"description,Markdown body of SKILL.md after frontmatter"`
	Message        string            `json:"message,omitempty" jsonschema:"description,Status or error message"`
}

// Agent Context Models

type GetAgentContextRequest struct{}

type GetAgentContextResponse struct {
	Content string `json:"content" jsonschema:"description,Full agent context documentation"`
	Message string `json:"message,omitempty" jsonschema:"description,Status or error message"`
}

// Error Response

type ErrorResponse struct {
	Code    string `json:"code" jsonschema:"description,Error code"`
	Message string `json:"message" jsonschema:"description,Error message"`
}
