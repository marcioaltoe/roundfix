---
task: task_01
spec: 0011-storage-lifecycle
status: pending
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

- [ ] Resolver with the three-branch hierarchy
- [ ] Spec association: `--spec` selector plus `Roundfix-Spec` trailer discovery with folder validation
- [ ] Wire resolver into resolve and watch artifact roots
- [ ] Table tests over all branches and the explicit-config passthrough

## Acceptance Criteria

- [ ] A spec-associated PR (via `--spec` or trailer) writes Round artifacts under `docs/specs/<slug>/reviews/`.
- [ ] A spec-less PR writes under `docs/specs/_reviews/pr-<n>/`.
- [ ] An explicitly configured `artifact_dir` still wins and its layout is unchanged.
- [ ] An unknown/invalid trailer slug falls back to the spec-less default rather than writing to a non-existent folder.

## Verification

- `rtk go test ./internal/config/ ./internal/cli/` — expected: all tests pass, including the new resolver table tests.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1-3; Core Feature 1. `_techspec.md` → Review artifact
location, Build Order 1, Interfaces: `ResolveReviewRoot`. ADR-0029 (supersedes
ADR-0003). Work-plan finding R1-17.
