## Why

The `mcpclaw` binary is not easily distributed to users. Currently, users must clone the repo, build locally, or download pre-built binaries manually. By automating binary builds and releases on every push to main, users can quickly curl and run the latest version without friction.

## What Changes

- Create GitHub Actions workflow that triggers on every push to main
- Build stripped `mcpclaw` binary for Linux x86-64
- Automatically create/update a GitHub Release with version tag `v{YYYY.MM.DD}-{git-sha-short}`
- Attach binary as release asset `mcpclaw-linux-amd64` for easy download

## Capabilities

### New Capabilities
- `binary-release-automation`: Automated building and publishing of mcpclaw binary to GitHub Releases on every push to main, with date-based versioning and git hash

### Modified Capabilities
<!-- No existing capabilities require changes -->

## Impact

- Users can now install mcpclaw directly: `curl -L https://github.com/.../releases/download/v{version}/mcpclaw-linux-amd64 -o mcpclaw && chmod +x mcpclaw`
- CI/CD system: Adds new GitHub Actions workflow file
- Release cadence: Every push to main generates a timestamped release
