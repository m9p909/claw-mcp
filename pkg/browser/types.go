package browser

// Browser Navigation Requests/Responses

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

// Browser Snapshot Request/Response

type BrowserSnapshotRequest struct{}

type BrowserSnapshotResponse struct {
	Snapshot string `json:"snapshot" jsonschema:"description,Structured accessibility tree representation"`
}

// Browser Click Request/Response

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

// Browser Hover Request/Response

type BrowserHoverRequest struct {
	Ref     string `json:"ref" jsonschema:"description,Exact target element reference from snapshot"`
	Element string `json:"element" jsonschema:"description,Human-readable element description"`
}

type BrowserHoverResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether hover succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// Browser Type Request/Response

type BrowserTypeRequest struct {
	Text string `json:"text" jsonschema:"description,Text to type"`
	Ref  string `json:"ref" jsonschema:"description,Element reference to focus first (optional)"`
}

type BrowserTypeResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether typing succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// Browser Fill Form Request/Response

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

// Browser Select Option Request/Response

type BrowserSelectOptionRequest struct {
	Ref     string   `json:"ref" jsonschema:"description,Select element reference (required)"`
	Values  []string `json:"values" jsonschema:"description,Option values to select (required)"`
	Element string   `json:"element" jsonschema:"description,Human-readable element description"`
}

type BrowserSelectOptionResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether select succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// Browser Press Key Request/Response

type BrowserPressKeyRequest struct {
	Key string `json:"key" jsonschema:"description,Key name (Enter, Tab, Escape, etc.) or single character"`
	Ref string `json:"ref" jsonschema:"description,Element to focus first (optional)"`
}

type BrowserPressKeyResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether key press succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// Browser Wait For Request/Response

type BrowserWaitForRequest struct {
	Text     string  `json:"text" jsonschema:"description,Text to wait for (optional)"`
	TextGone string  `json:"textGone" jsonschema:"description,Text to wait to disappear (optional)"`
	Time     float64 `json:"time" jsonschema:"description,Seconds to wait (optional)"`
}

type BrowserWaitForResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether wait succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// Browser Handle Dialog Request/Response

type BrowserHandleDialogRequest struct {
	Accept     bool   `json:"accept" jsonschema:"description,true to accept/OK, false to cancel"`
	PromptText string `json:"promptText" jsonschema:"description,Text to enter in prompt dialogs (optional)"`
}

type BrowserHandleDialogResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether dialog handled successfully"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// Browser Close Request/Response

type BrowserCloseRequest struct{}

type BrowserCloseResponse struct {
	Success bool   `json:"success" jsonschema:"description,Whether close succeeded"`
	Message string `json:"message" jsonschema:"description,Status or error message"`
}

// Browser Configuration

type BrowserConfig struct {
	IdleTimeoutSecs int
	ToolTimeoutSecs int
}
