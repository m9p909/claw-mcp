## ADDED Requirements

### Requirement: Logger initialization with DEBUG_LEVEL support
The system SHALL initialize a structured logger on startup that respects the DEBUG_LEVEL environment variable to control verbosity.

#### Scenario: Logger starts with DEBUG_LEVEL=INFO
- **WHEN** application starts with DEBUG_LEVEL=INFO (or no DEBUG_LEVEL set)
- **THEN** logger is configured to output INFO, WARN, and ERROR levels only

#### Scenario: Logger starts with DEBUG_LEVEL=DEBUG
- **WHEN** application starts with DEBUG_LEVEL=DEBUG
- **THEN** logger is configured to output DEBUG, INFO, WARN, and ERROR levels

#### Scenario: Logger starts with DEBUG_LEVEL=TRACE
- **WHEN** application starts with DEBUG_LEVEL=TRACE
- **THEN** logger is configured to output all levels including TRACE

### Requirement: Request ID generation and context propagation
The system SHALL generate a unique UUID request ID for each incoming HTTP request and propagate it through context to all tool handlers.

#### Scenario: Request ID is generated and available in logs
- **WHEN** an HTTP request arrives at /mcp
- **THEN** a UUID is generated and passed via context
- **AND** all subsequent log statements include this request ID

#### Scenario: Request ID correlates logs across tool calls
- **WHEN** a tool handler calls multiple internal functions
- **THEN** all log statements from those calls include the same request ID
- **AND** failures from the same request can be traced together

### Requirement: Session ID extraction from HTTP headers
The system SHALL extract the Mcp-Session-Id header from HTTP requests and propagate it through context to all tool handlers.

#### Scenario: Session ID is extracted and included in logs
- **WHEN** an HTTP request includes Mcp-Session-Id header
- **THEN** the session ID is extracted and passed via context
- **AND** all subsequent log statements include this session ID

#### Scenario: Multiple requests in same session are correlated
- **WHEN** tool handlers are called with the same Mcp-Session-Id across multiple requests
- **THEN** all logs from those requests share the same session ID
- **AND** an observer can trace the full session timeline

### Requirement: Structured logging functions with request/session context
The system SHALL provide logging functions (Info, Warn, Error) that automatically include request ID and session ID from context.

#### Scenario: Info log includes context
- **WHEN** a tool handler calls logger.Info("operation started")
- **THEN** the log entry includes request ID and session ID
- **AND** it is human-readable text format

#### Scenario: Error log includes context and error details
- **WHEN** a tool handler calls logger.Error("operation failed", err)
- **THEN** the log entry includes request ID, session ID, error message, and error type

#### Scenario: Logger handles missing context gracefully
- **WHEN** logging is called without request ID or session ID in context
- **THEN** logging still works with empty/nil fields instead of panicking
