## ADDED Requirements

### Requirement: Log foreground command execution
The system SHALL log foreground command execution with startup, completion, and timing.

#### Scenario: Successful command execution is logged
- **WHEN** HandleExecCommand executes a foreground command successfully
- **THEN** an INFO log records command and arguments (sanitized)
- **AND** a DEBUG log records execution start time
- **AND** an INFO log records exit code and execution time on completion
- **AND** a DEBUG log records stdout/stderr sizes

#### Scenario: Command failure is logged with exit code
- **WHEN** HandleExecCommand executes a command that exits with non-zero code
- **THEN** an WARN log records the non-zero exit code
- **AND** an INFO log records the stdout size for troubleshooting
- **AND** execution time is included

#### Scenario: Command execution error is logged
- **WHEN** HandleExecCommand fails to start a command (e.g., command not found)
- **THEN** an ERROR log records EXEC_FAILED
- **AND** the underlying OS error is included in full detail

#### Scenario: Stdout/stderr pipe creation failures are logged
- **WHEN** foreground execution cannot create stdout or stderr pipes
- **THEN** an ERROR log records pipe creation failure
- **AND** the underlying error is included

### Requirement: Log background command execution
The system SHALL log background command execution including session creation and process monitoring.

#### Scenario: Background command session is created and logged
- **WHEN** HandleExecCommand starts a background command
- **THEN** an INFO log records session ID
- **AND** an INFO log records the command name and arguments (sanitized)
- **AND** a DEBUG log records goroutine spawn for stdout/stderr reading

#### Scenario: Background command completion is logged
- **WHEN** a background command completes
- **THEN** an INFO log records session ID, exit code, and execution time
- **AND** a DEBUG log records final stdout/stderr sizes

#### Scenario: Background stdout/stderr capture is logged at TRACE level
- **WHEN** background pipes are being read during execution
- **THEN** TRACE logs record chunks received (with size, not content)
- **AND** this only appears if DEBUG_LEVEL=TRACE

### Requirement: Log manage_process operations
The system SHALL log process management operations (list, poll, kill, send_keys).

#### Scenario: List operation is logged
- **WHEN** HandleManageProcess is called with action="list"
- **THEN** an INFO log records the list action
- **AND** a DEBUG log records the number of sessions returned

#### Scenario: Poll operation is logged
- **WHEN** HandleManageProcess is called with action="poll" and valid session ID
- **THEN** an INFO log records the session ID being polled
- **AND** a DEBUG log records the session status snapshot

#### Scenario: Poll with invalid session ID is logged
- **WHEN** HandleManageProcess is called with action="poll" and non-existent session ID
- **THEN** an ERROR log records PROCESS_NOT_FOUND
- **AND** the session ID is included for debugging

#### Scenario: Invalid action is logged
- **WHEN** HandleManageProcess is called with unknown action
- **THEN** an ERROR log records INVALID_REQUEST with the action name
