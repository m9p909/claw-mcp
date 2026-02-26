## 1. Create GitHub Actions Workflow

- [x] 1.1 Create `.github/workflows/` directory if it doesn't exist
- [x] 1.2 Create `.github/workflows/build-release.yml` workflow file
- [x] 1.3 Configure workflow to trigger on push to main branch
- [x] 1.4 Add checkout action to get code
- [x] 1.5 Add setup Go action (Go 1.25)
- [x] 1.6 Add build step with `-ldflags="-s -w"` for binary stripping
- [x] 1.7 Add step to compute version tag from date and git hash

## 2. Release Management

- [x] 2.1 Add step to delete existing release from current date (if exists)
- [x] 2.2 Add step to create new GitHub Release with computed version tag
- [x] 2.3 Add step to upload binary as release asset named `mcpclaw-linux-amd64`
- [x] 2.4 Configure workflow to use GitHub token from secrets for authentication

## 3. Testing and Validation

- [x] 3.1 Push code to main and verify workflow triggers
- [x] 3.2 Verify release is created with correct version tag format
- [x] 3.3 Verify binary asset is attached to release
- [x] 3.4 Download binary using curl and verify it's executable
- [x] 3.5 Test that subsequent push updates release with new binary

## 4. Documentation

- [x] 4.1 Add installation instructions to README or documentation
- [x] 4.2 Document how to download and run mcpclaw binary from releases
