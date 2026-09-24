---
task: task_06
spec: 0155-a-verify-that-runs-what-changed
status: completed
type: backend
complexity: medium
---

# Task 06: Select every set a change can affect

## Overview

Corrective Task from the pre-PR review of 2026-09-24. Three selector gaps let a broken change pass `make verify-changed`: a non-`.go` path that is not documentation, such as a `testdata/` fixture or a `.golden` snapshot, selects no set; staged but uncommitted paths are never listed; and a Baseline-only change skips the core packages that import Baseline code (`internal/speccheck`, `internal/cli`).

## Requirements

1. MUST send a path under a package directory, including its `testdata/`, to the set that owns that package.
2. MUST select no set only for an explicit documentation allowlist (`docs/**` and Markdown outside package and embedded directories); every other unknown path fails safe to both sets.
3. MUST list staged changes as well as committed, unstaged and untracked ones.
4. MUST, when the Baseline set is selected, also run the core packages whose imports or test imports reach a Baseline package, derived from `go list` rather than a fixed list.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A fixture-only change in a core package selects the core set.
- [ ] An unknown non-documentation path selects both sets.
- [ ] A staged-only change is selected.
- [ ] A Baseline-only change also runs the core importers of Baseline packages.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/verifyselect/verifyselect.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestFixtureChangesSelectTheOwningSet|TestUnknownNonDocumentationPathsFailSafe|TestSelectListsStagedPaths|TestBaselineChangesAlsoRunCoreImporters|TestClassifyPath|TestSelectListsCommittedUnstagedAndUntrackedPaths)$" ./internal/verifyselect 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestFixtureChangesSelectTheOwningSet TestUnknownNonDocumentationPathsFailSafe TestSelectListsStagedPaths TestBaselineChangesAlsoRunCoreImporters TestClassifyPath TestSelectListsCommittedUnstagedAndUntrackedPaths; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task four of the named cases do not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

Implemented package-aware fail-safe selection in `internal/verifyselect`:

- Paths under the Baseline roots select Baseline, paths under core package trees select core regardless of extension, and only `docs/**` plus Markdown outside package and embedded directories select neither. Every other unknown path selects both sets.
- Changed-path discovery now unions the staged index with committed, unstaged, and untracked paths.
- Baseline package expansion now decodes `go list -json ./...` and walks the reverse regular/test import graph, so core packages that directly or transitively reach Baseline packages run with a Baseline selection.
- Added focused regression coverage for every acceptance criterion and extended `TestClassifyPath` for the documentation allowlist boundary.

Focused-check evidence:

- Fixture-only core/Baseline ownership, unknown non-documentation fail-safe behavior, staged-only discovery, and Baseline importer expansion: `rtk env GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 -run '^(TestFixtureChangesSelectTheOwningSet|TestUnknownNonDocumentationPathsFailSafe|TestSelectListsStagedPaths|TestBaselineChangesAlsoRunCoreImporters)$' ./internal/verifyselect` exited 0 (`ok roundfix/internal/verifyselect`).
- Existing path classification and committed/unstaged/untracked discovery: `rtk env GOCACHE=/private/tmp/roundfix-task06-gocache go test -count=1 -run '^(TestClassifyPath|TestSelectListsCommittedUnstagedAndUntrackedPaths)$' ./internal/verifyselect` exited 0 (`ok roundfix/internal/verifyselect`).
- Real repository package derivation: `rtk env GOCACHE=/private/tmp/roundfix-task06-gocache go run ./cmd/verify-select -packages baseline` exited 0 and included `./internal/cli` and `./internal/speccheck` among the derived Baseline-run packages.
- `rtk git diff --check` exited 0.

Follow-up: Task 07 owns updating the partition contract for the intentional overlap introduced when core importers also run with the Baseline set. The authored `## Verification` command was not run; Daemon Verification owns it.

## Carry-forward provenance

- Source Run: `run_20260924T174549Z_5664c74b2c5ad5f5`
- Source commit: `e4e5d4be5512f9f19abbcbe83a4e80b731117c8b`
