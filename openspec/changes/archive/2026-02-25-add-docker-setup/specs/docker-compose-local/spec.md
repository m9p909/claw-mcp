## ADDED Requirements

### Requirement: Docker Compose file defines multi-container development environment
A docker-compose.yml file SHALL define a local development environment with Claw server and optional supporting services.

#### Scenario: Compose start brings up local server
- **WHEN** developer runs `docker-compose up`
- **THEN** the Claw server starts and is accessible on localhost:8080

### Requirement: Docker Compose defines Claw data volume
The Compose file SHALL define a named or bind volume for `~/.mcpclaw`, allowing developers to interact with workspace and database files from the host machine.

#### Scenario: Developer accesses workspace from host
- **WHEN** an agent writes a file to `~/.mcpclaw/workspace` in the container
- **THEN** the file is accessible on the host machine at the mapped path

#### Scenario: Memory persists across compose restart
- **WHEN** the developer runs `docker-compose down` and `docker-compose up`
- **THEN** previously written memory data in `~/.mcpclaw/data` is still accessible

### Requirement: Port mapping for local access
The Compose service SHALL map the container port 8080 to a host port (default 8080), allowing local clients to connect to the server.

#### Scenario: Local health check succeeds
- **WHEN** developer runs `curl http://localhost:8080/health`
- **THEN** the response is `{"status":"ok"}`

### Requirement: Environment variable configuration in Compose
The Compose service definition SHALL expose environment variables for customization (e.g., PORT, LOG_LEVEL), allowing developers to adjust server behavior without editing the Dockerfile.

#### Scenario: Developer overrides port via Compose environment
- **WHEN** Compose is configured with `PORT=9000` in environment section
- **THEN** the server listens on port 9000 within the container
