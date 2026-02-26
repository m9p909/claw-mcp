## Context

Currently, mcpclaw binary distribution relies on manual local builds or pre-built artifacts. The project uses Go 1.25 with a single `main.go` entry point that compiles to the `mcpclaw` binary. GitHub Actions is already available as the CI/CD platform (no additional infrastructure needed).

## Goals / Non-Goals

**Goals:**
- Automate binary builds and releases on every push to main
- Use explicit versioning (v{YYYY.MM.DD}-{git-sha-short}) for traceability
- Strip debug symbols to reduce binary size
- Target Linux x86-64 platform
- Support curl-based installation workflow

**Non-Goals:**
- Multi-platform builds (macOS, Windows) - Linux x86-64 only for now
- Signed releases or provenance attestation
- Automated changelog generation from commits
- Pre-release or draft releases - all releases are final

## Decisions

**Decision: Use explicit version tags instead of "latest" release**
- Rationale: Users can always reference the exact version they're running, enabling reproducibility. Explicit tags prevent accidental breakage if old scripts reference "latest".
- Alternative: Use "latest" release to simplify for casual users. Rejected: less transparent, harder to roll back.

**Decision: Overwrite release on each push to main (not accumulate releases)**
- Rationale: Each push generates a unique version tag (date + git hash), so we always have history in git. No need to keep stale releases. Keeps GitHub release list clean.
- Alternative: Create new release for each push. Rejected: release list becomes cluttered, wastes storage.

**Decision: Use `-ldflags="-s -w"` for binary stripping**
- Rationale: Removes debug symbols and DWARF info, reducing binary size by ~40-50% with minimal trade-off (no debugging in production, but symbols still in git).
- Alternative: Keep debug info. Rejected: larger download, not needed for users.

**Decision: Workflow triggers on `push to main` only**
- Rationale: Main branch represents stable, deployable code. Prevents release spam from feature branches.
- Alternative: Manual dispatch or tag-triggered. Rejected: less automation, users must remember to trigger.

**Decision: Version format is `v{YYYY.MM.DD}-{git-sha-short}`**
- Rationale: Human-readable date ordering + unique git commit hash. Ensures builds from same date are distinguishable.
- Alternative: Semantic versioning (v1.0.0). Rejected: requires manual version file updates, more overhead.

**Decision: Use GitHub CLI (gh) in workflow for release management**
- Rationale: Built-in to GitHub Actions environment, simpler than REST API calls, supports JSON output for parsing.

## Risks / Trade-offs

**[Risk: High release churn]** → Each push creates a new release. Mitigated by: Releases use content-addressable versions (git hash), so any push can be recovered. Users expected to use explicit version tags, not scroll release list.

**[Risk: No cross-platform builds]** → Linux x86-64 only, but project may eventually need macOS/Windows binaries. Mitigated by: Design is extensible; adding new platforms is just adding build matrix entries to workflow.

**[Risk: No signature verification]** → Users can't verify binary authenticity. Acceptable for now as project has limited security-critical use case, but signature verification could be added later.

**[Trade-off: Stripped binary vs debuggability]** → Users can't debug production binary, but debug symbols are always available in git. Acceptable for distribution use case where size matters.

## Migration Plan

1. Create `.github/workflows/build-release.yml` with build and release steps
2. Test workflow on first push to main - will create first release with version `v2026.02.26-{current-sha}`
3. Subsequent pushes will update the release with new binary

No rollback needed - workflow is additive, only creates releases. Can disable by removing workflow file if needed.
