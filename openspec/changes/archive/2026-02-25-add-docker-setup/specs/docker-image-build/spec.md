## ADDED Requirements

### Requirement: Docker image contains compiled Claw server
The Docker image SHALL contain the compiled Claw MCP Server binary built from source, ready to execute on container startup.

#### Scenario: Image starts server successfully
- **WHEN** a container is launched from the image
- **THEN** the Claw server starts on port 8080 and responds to health checks

### Requirement: Image includes batteries-included Unix tools
The image SHALL include core Unix utilities required for agent execution, including git, curl, wget, grep, sed, awk, and find.

#### Scenario: Agent can execute standard Unix commands
- **WHEN** an agent executes a command like `git clone` or `grep` in the workspace
- **THEN** the command succeeds without requiring additional installation

### Requirement: Image includes Node.js runtime
The image SHALL include Node.js with npm package manager, allowing agents to run JavaScript/Node code in the workspace.

#### Scenario: Agent runs Node script
- **WHEN** an agent executes a `.js` file with the Node runtime
- **THEN** the script runs successfully with access to npm packages

### Requirement: Image includes Python runtime
The image SHALL include Python 3 with pip package manager, allowing agents to run Python scripts in the workspace.

#### Scenario: Agent runs Python script
- **WHEN** an agent executes a `.py` file with the Python runtime
- **THEN** the script runs successfully with access to pip packages

### Requirement: Image uses multi-stage build
The Docker build process SHALL use a multi-stage approach with separate builder and runtime stages to minimize final image size and reduce attack surface.

#### Scenario: Runtime image excludes build dependencies
- **WHEN** the image is built
- **THEN** the final image contains only the compiled binary and runtime dependencies, not Go SDK or build tools

### Requirement: Entrypoint script handles initialization
The image SHALL use an entrypoint script that initializes the server, handles environment configuration, and manages volume mounts.

#### Scenario: Server starts with configured port
- **WHEN** the container starts with a PORT environment variable
- **THEN** the server listens on the specified port instead of the default 8080
