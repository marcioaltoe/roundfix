---
task: task_02
spec: 0310-the-python-profile-writes-its-starter-files
status: pending
type: backend
complexity: high
---

# Task 02: Starter Files are previewed and written only where absent on a confirmed digest

## Overview

This Task adds the Baseline operation behind `roundfix baseline starters`: it
previews a profile's Starter Files as `create` or `keep` with a Plan Digest,
writes only absent paths when the current digest is confirmed, refuses a
digest that no longer matches, and removes what it created when a write fails.
Profile alignment names it in the next action of an undeclared Verification
command for a profile with Starter Files (`_prd.md` → Core Features 3-4;
ADR-0323; ADR-0068).

## Requirements

1. MUST answer the Backlog Entry Spec 0309 adopted,
   [`../0309-a-python-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md`](../0309-a-python-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md)
   of 2026-10-22, its expected item 4.
2. MUST create `internal/baseline/starters.go` with `StartersSchemaVersion`,
   `StartersRequest`, `StarterEntry`, `StartersPayload`, `ErrStartersStale`
   and `Starters` of `_techspec.md` → Interfaces, implementing API Contracts 3
   and 4 and Data Model 3: the profile resolved through `ResolveProfile`
   (built-in or repository-owned), `create` only where `os.Lstat` fails at the
   path, writes through the repository root with exclusive creation, parents
   at mode 0755 and files at 0644, refusal of a different digest with nothing
   written, and removal of every file the call created when a later write
   fails. The hazard: the repository's existing files are the destination, and
   the operation never opens one for writing.
3. MUST change `internal/baseline/profile_alignment.go` so the next action of
   `verification.command.undeclared` is the text of `_techspec.md` → API
   Contract 5 when the resolved profile's modules declare a Starter File, and
   keep its bytes otherwise; the plan characterization goldens, whose profiles
   declare none, MUST NOT change.
4. MUST create `internal/baseline/starters_test.go` with the four tests of
   `_techspec.md` → Testing Approach (Build Order 2), using `newPlanRepository`,
   `mustEmbeddedCatalog`, `ResolveProfileAlignment`, `findProfileDivergence`
   and the decisions helper the existing alignment tests use. Each top-level
   test calls `t.Parallel()` first.
5. MUST keep every existing test of `internal/baseline` passing without
   editing it.
6. MUST NOT add a CLI command, change any template, module, profile, skill,
   user guide page, `.roundfixrc.yml` or this repository's Makefile, or add a
   production hook used only by tests.

## Subtasks

- [ ] Add the preview and the confirmed write with removal on failure.
- [ ] Name the command in the undeclared-gate next action.
- [ ] Add the four tests.

## Acceptance Criteria

- [ ] A preview writes nothing; a confirmed apply creates only absent paths and leaves existing files byte-identical.
- [ ] A stale digest writes nothing, and a failed write leaves none of the files the call created.
- [ ] The next action names the command for the `python` profile and keeps its text for `rust-cli`.

## Consumers

- `Starters` and its payload — consumers: `./internal/baseline`
- checkpoint: task_03

## Context

- instruction: `docs/adr/0323-a-profile-starter-file-is-written-only-where-absent-and-only-on-a-confirmed-preview.md`
- instruction: `docs/adr/0068-baseline-command-uses-one-confirmation-gated-workflow.md`
- instruction: `internal/baseline/profile_alignment_test.go`
- creates: `internal/baseline/starters.go`
- interface: `internal/baseline/profile_alignment.go`
- creates: `internal/baseline/starters_test.go`

## Verification

- `for name in TestStartersPreviewsAndCreatesOnlyAbsentFiles TestStartersRefuseADigestThatNoLongerMatches TestStartersRemoveWhatTheyCreatedWhenAWriteFails TestTheUndeclaredGateNamesTheStartersCommand; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/baseline && GOOS=windows go vet ./internal/baseline` — expected: exit 0. Before this Task none of the four tests exists, so the command fails at its first check. Afterwards they exist and the whole package passes, including the plan characterizations, and its tests compile for Windows.

## References

- `_prd.md` → Goals; Core Features 3-4; Success Metrics 3-4; User Stories 2-4
- `_techspec.md` → Current behavior 4; Interfaces; Data Model 3; API Contracts 3-5; Invariants 1-3 and 5; Testing Approach; Build Order 2
- ADR-0323; ADR-0068
