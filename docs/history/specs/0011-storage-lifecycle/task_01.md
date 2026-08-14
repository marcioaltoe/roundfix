---
task: task_01
spec: 0011-storage-lifecycle
status: completed
type: backend
complexity: medium
---

# Task 01: Review artifact location resolver and spec association

## Overview

Move review-artifact placement behind a resolver so Round and Review Issue
artifacts default into the repository's spec tree instead of loose in Roundfix
Home. The resolver chooses explicit config first, then a spec-associated folder,
then an in-repo review root — a self-contained, table-testable slice wired into
the resolve and watch artifact paths.

## Requirements

1. MUST add a resolver that returns the directory under which
   `reviews/pr-<n>/round-*` is written, choosing in order: an explicitly
   configured Artifact Directory; else `docs/specs/<slug>/reviews/` when the PR
   is associated with a Spec; else `docs/specs/_reviews/pr-<n>/`.
2. MUST derive the Spec association from an explicit `--spec <slug>` selector
   (wins) or the newest `Roundfix-Spec` commit trailer on the PR head; a
   discovered slug MUST name an existing spec folder or be treated as no
   association.
3. MUST wire the resolver into the resolve and watch artifact paths, replacing
   the current `<artifact_dir>/reviews/pr-<n>` derivation, and MUST NOT commit
   or gitignore the artifacts.
4. MUST leave behavior byte-stable for a user who set `artifact_dir` explicitly.

## Subtasks

- [x] Resolver with the three-branch hierarchy
- [x] Spec association: `--spec` selector plus `Roundfix-Spec` trailer discovery with folder validation
- [x] Wire resolver into resolve and watch artifact roots
- [x] Table tests over all branches and the explicit-config passthrough

## Acceptance Criteria

- [x] A spec-associated PR (via `--spec` or trailer) writes Round artifacts under `docs/specs/<slug>/reviews/`.
- [x] A spec-less PR writes under `docs/specs/_reviews/pr-<n>/`.
- [x] An explicitly configured `artifact_dir` still wins and its layout is unchanged.
- [x] An unknown/invalid trailer slug falls back to the spec-less default rather than writing to a non-existent folder.

## Verification

- `rtk go test ./internal/config/ ./internal/cli/` — expected: all tests pass, including the new resolver table tests.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1-3; Core Feature 1. `_techspec.md` → Review artifact
location, Build Order 1, Interfaces: `ResolveReviewRoot`. ADR-0029 (supersedes
ADR-0003). Work-plan finding R1-17.

## Result

Implemented the Review artifact root resolver and wired operational review
commands to use it for Round and Review Issue artifact placement. `rounds` now
accepts a resolved Review root while preserving the legacy Artifact Directory
fallback for existing callers. `fetch`, `resolve`, and `watch` resolve the
Review root after preflight; the daemon cycle uses the same root when counting
remaining Review Issues.

Evidence:

- `rtk go test ./internal/config/ ./internal/cli/` — passed, 304 tests.
- `rtk go test ./...` — passed, 747 tests across 17 packages.
- `rtk make verify` — passed: full Go suite, `roundfix skills check`, and build.

Acceptance evidence:

- Spec-associated placement: `TestRunFetchWritesReviewArtifactsUnderSpecSelector`
  and `TestRunFetchWritesReviewArtifactsUnderTrailerSpec` passed, asserting
  `docs/specs/<slug>/reviews/round-001/issue_001.md`.
- Spec-less placement: `TestRunFetchWritesReviewArtifactsUnderSpeclessRoot`
  passed, asserting `docs/specs/_reviews/pr-123/round-001/issue_001.md`.
- Explicit `artifact_dir` passthrough: `TestResolveReviewRoot/explicit artifact
  directory keeps existing layout` passed, asserting
  `<artifact_dir>/reviews/pr-123`.
- Invalid trailer fallback: `TestReviewSpecSlugAssociation/unknown trailer slug
  is no association` and `TestResolveReviewRoot/unknown spec falls back to
  spec-less root` passed.

Follow-up:

- Docs and shipped Roundfix skill updates for the new `--spec` operational flag
  remain in task_05, per this Spec's build order.
