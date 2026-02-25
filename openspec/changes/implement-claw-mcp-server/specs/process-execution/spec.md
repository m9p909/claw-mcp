

### Requirement: Execute commands synchronously
The system SHALL execute shell commands synchronously, capturing stdout/stderr and return exit code with full output.

#### Scenario: Successful execution
- **WHEN** client calls exec_command with valid command and background=false
- **THEN** system runs command, waits for completion, returns exit_code and output

#### Scenario: Execution timeout
- **WHEN** command exceeds timeout duration (default 60s)
- **THEN** system kills process, returns TimeoutError with partial output captured so far

#### Scenario: Command not found
- **WHEN** executable doesn't exist
- **THEN** system returns error describing process launch failure

#### Scenario: Custom environment variables
- **WHEN** client provides env map
- **THEN** system merges with inherited environment and executes

### Requirement: Execute commands in background
The system SHALL spawn background processes and return immediately with session ID for polling.

#### Scenario: Background process launch
- **WHEN** client calls exec_command with background=true
- **THEN** system returns session_id immediately while process runs detached

#### Scenario: Process continues after request
- **WHEN** background process is spawned
- **THEN** process continues running independent of client connection

### Requirement: Manage running processes
The system SHALL support listing, polling, sending input, and killing processes by session ID.

#### Scenario: List all sessions
- **WHEN** client calls manage_process with action=list
- **THEN** system returns all sessions (running and completed) with status, PID, timestamps

#### Scenario: Poll running process
- **WHEN** client polls a running session
- **THEN** system returns status=running, PID, elapsed time

#### Scenario: Poll completed process
- **WHEN** client polls a completed session
- **THEN** system returns status=completed, exit_code, total runtime

#### Scenario: Send input to process
- **WHEN** client calls manage_process with action=send_keys and input
- **THEN** system writes input to process stdin, returns bytes_sent

#### Scenario: Kill process
- **WHEN** client calls manage_process with action=kill
- **THEN** system sends SIGTERM, updates session status to killed

#### Scenario: Session not found
- **WHEN** client references invalid session_id
- **THEN** system returns SessionNotFoundError

### Requirement: Capture process output
The system SHALL buffer all stdout/stderr from executed processes with no size limits, available via poll action.

#### Scenario: Retrieve output
- **WHEN** client calls manage_process with action=log and session_id
- **THEN** system returns stdout and stderr as arrays of lines, with line counts
