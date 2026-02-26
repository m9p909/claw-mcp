## ADDED Requirements

### Requirement: Development Docker Compose with self-signed TLS
The system SHALL provide a docker-compose.yml configuration that runs Claw and Caddy with TLS enabled for localhost/development environments using self-signed certificates.

#### Scenario: User starts development environment with TLS
- **WHEN** user runs `docker-compose up` with development environment variables
- **THEN** Caddy generates self-signed certificates for localhost
- **AND** Caddy listens on port 443 (HTTPS)
- **AND** Claw is accessible via `https://localhost`
- **AND** internal Claw port 8080 is not exposed to the host

#### Scenario: Claw container and Caddy container communicate via internal network
- **WHEN** Caddy receives a request on port 443
- **THEN** Caddy routes the request to Claw on the internal docker-compose network
- **AND** traffic between containers is unencrypted HTTP (internal only)

### Requirement: Production Docker Compose with Let's Encrypt TLS
The system SHALL provide a docker-compose.prod.yml configuration that runs Claw and Caddy with TLS enabled using Let's Encrypt certificates for production domains.

#### Scenario: User starts production environment with Let's Encrypt
- **WHEN** user runs `docker-compose -f docker-compose.prod.yml up` with DOMAIN and CLAW_TOKEN environment variables
- **THEN** Caddy requests and obtains a valid Let's Encrypt certificate for the specified domain
- **AND** Caddy automatically renews the certificate 30 days before expiration
- **AND** Claw is accessible via `https://<DOMAIN>`
- **AND** HTTP requests to port 80 are redirected to HTTPS

#### Scenario: Certificate renewal happens automatically
- **WHEN** a Let's Encrypt certificate is within 30 days of expiration
- **THEN** Caddy automatically renews the certificate without manual intervention
- **AND** Claw remains accessible during renewal (zero downtime)

### Requirement: Caddyfile template with environment variable substitution
The system SHALL provide a templated Caddyfile that uses environment variables to configure TLS mode and domain.

#### Scenario: Template substitutes DOMAIN variable
- **WHEN** Caddyfile is loaded with DOMAIN environment variable set
- **THEN** Caddy serves the specified domain with appropriate certificate configuration
- **AND** template works identically in both development and production compose files

#### Scenario: Template defaults to localhost in development
- **WHEN** DOMAIN is set to "localhost"
- **THEN** Caddy generates and uses self-signed certificates
- **AND** no Let's Encrypt requests are made

### Requirement: Port configuration prevents direct HTTP access to Claw
The system SHALL not expose Claw's internal port 8080 to the host machine.

#### Scenario: Claw port is internal-only
- **WHEN** inspecting docker-compose port mappings
- **THEN** only Caddy exposes ports to the host (80 for HTTP, 443 for HTTPS)
- **AND** Claw container has no port mappings to the host
- **AND** users cannot bypass TLS by connecting to port 8080

### Requirement: Caddy configuration via environment variables
The system SHALL support configuration of TLS behavior through environment variables passed to the Caddy container.

#### Scenario: DOMAIN variable configures domain/certificate
- **WHEN** DOMAIN environment variable is set
- **THEN** Caddy uses that domain for certificate acquisition and TLS configuration

#### Scenario: TLS_MODE variable configures certificate provider
- **WHEN** TLS_MODE="self_signed"
- **THEN** Caddy generates self-signed certificates for development
- **WHEN** TLS_MODE="letsencrypt"
- **THEN** Caddy obtains certificates from Let's Encrypt for production
