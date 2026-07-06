---
task: task_01
spec: 0009-parallel-scheduling
status: completed
type: data
complexity: medium
---

# Task 01: Worktree config — concurrency, location hierarchy, slug paths

## Overview

The configuration ground for parallelism: `worktree.concurrency` and
`worktree.location` with the repo > global > builtin hierarchy, the
readable slug-based path derivation both the Run and future Task Worktrees
share, and the removal of the dead `resolve.concurrent` key. Verifiable
through config and worktree-path unit tests.

## Requirements

1. MUST add `worktree.concurrency` (int ≥ 1, builtin default 2) and
   `worktree.location` (parent directory; `~` and absolute forms; builtin
   default `~/.roundfix/worktrees`), both resolving repo Project Config >
   User Config > builtin, documented in generated config output.
2. MUST introduce one shared path-derivation helper: worktree roots are
   always `<location>/<repo-slug>/...` where `repo-slug` is the sanitized
   repository directory basename plus `-` plus 8 hex characters of the
   user-root path hash; the slug and the run-id segments are appended
   unconditionally and are never configurable.
3. MUST switch Run Worktree creation to the new derivation (replacing the
   0008 bare-hash segment) — new Runs only; existing kept worktrees stay
   discoverable through their recorded `work_dir`.
4. MUST remove `resolve.concurrent` from the config schema and generated
   output; a config file still carrying it fails validation with a named
   error pointing to `worktree.concurrency`.
5. MUST validate: concurrency < 1 rejected; location must be absolute after
   `~` expansion and must not fall inside the repository tree.

## Subtasks

- [x] Config keys, hierarchy resolution, generated output
- [x] Shared slug/path derivation helper
- [x] Run Worktree creation switched to the helper
- [x] resolve.concurrent removal with pointing error
- [x] Validation table tests

## Acceptance Criteria

- [x] Hierarchy tests: builtin only, User override, Project override
      winning, for both keys.
- [x] Path tests: slug readable + unique for two same-named repos at
      different paths; in-repo location rejected; `~` expansion covered.
- [x] A config with `resolve.concurrent` fails with the pointing message;
      generated config no longer mentions it and documents both new keys.
- [x] Existing worktree tests pass with only the deliberate path-shape
      updates.

## Verification

- `rtk go test ./internal/config/ ./internal/worktree/` — expected: all
  tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 5, 6; Core Feature 4; Decisions. `_techspec.md` →
Paths and naming, Build Order 1. ADR-0023 (location refinement).

## Result

Implemented `worktree.concurrency` and `worktree.location` with builtin,
User Config, and Project Config layering. `worktree.location` is resolved
through `~` expansion, must be absolute, and is rejected when it falls inside
the repository tree. The generated config documents both keys and no longer
emits `resolve.concurrent`; configs that still carry `resolve.concurrent`
fail with a message pointing to `worktree.concurrency`.

Added a shared worktree path derivation helper that builds
`<location>/<repo-slug>/<run-id>`, where the repo slug is the sanitized
repository basename plus 8 hex characters from the user-root path hash. Run
Worktree creation and terminal pruning now use the configured location and
that shared derivation.

Acceptance evidence:

- Hierarchy: `TestLoadAppliesWorktreeConfigHierarchy` covers builtin-only,
  User override, and Project override cases for both new keys.
- Path behavior: `TestDeriveRootPathUsesReadableUniqueRepoSlug` covers
  readable, unique slugs for same-named repositories; config tests cover
  in-repo rejection and `~` expansion.
- Deprecated key/output: `TestLoadRejectsInvalidConfig/deprecated_resolve_concurrent`
  covers the pointing error, and `TestInitCreatesUserConfig` covers generated
  output including the new keys and omitting the old key.
- Verification run during implementation: `rtk go test ./internal/config/
  ./internal/worktree/` passed with 35 tests; `rtk go test ./...` passed with
  694 tests.
