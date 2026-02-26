## ADDED Requirements

### Requirement: GitHub Actions workflow builds and publishes binary on push to main
The system SHALL automatically build the mcpclaw binary and create a GitHub Release whenever code is pushed to the main branch. The release tag SHALL be formatted as `v{YYYY.MM.DD}-{git-commit-sha-short}` and the binary artifact SHALL be named `mcpclaw-linux-amd64`.

#### Scenario: First push to main creates release
- **WHEN** code is pushed to main branch
- **THEN** GitHub Actions workflow triggers and creates a new GitHub Release with version tag `v{YYYY.MM.DD}-{git-commit-sha}` and attaches the binary `mcpclaw-linux-amd64` as an asset

#### Scenario: Subsequent push to main updates release
- **WHEN** code is pushed to main branch after a previous release
- **THEN** the existing release from the current date is deleted and recreated with the new binary built from the latest commit

### Requirement: Binary is stripped and optimized for distribution
The built mcpclaw binary SHALL have debug symbols removed to reduce size, suitable for end-user distribution.

#### Scenario: Binary is stripped
- **WHEN** the build process compiles mcpclaw with `-ldflags="-s -w"`
- **THEN** the resulting binary contains no debug symbols or DWARF information

### Requirement: Users can download and execute binary via curl
Users SHALL be able to download the compiled binary from a GitHub Release and execute it directly without requiring a local Go build environment.

#### Scenario: User downloads and runs latest binary
- **WHEN** a user executes `curl -L https://github.com/{owner}/awesomeProject/releases/download/v{YYYY.MM.DD}-{hash}/mcpclaw-linux-amd64 -o mcpclaw && chmod +x mcpclaw && ./mcpclaw`
- **THEN** the binary downloads successfully and executes without errors

### Requirement: Release versioning uses date and commit hash
The version tag for each release SHALL encode the date and git commit hash, enabling traceability and ensuring each push generates a unique, sortable version identifier.

#### Scenario: Version tag is correctly formatted
- **WHEN** a release is created on the 2026-02-26 from commit abc1234
- **THEN** the release tag is exactly `v2026.02.26-abc1234`

#### Scenario: Version tags sort chronologically
- **WHEN** releases are created on 2026-02-25, 2026-02-26, and 2026-02-27
- **THEN** sorting version tags alphabetically produces the correct chronological order
