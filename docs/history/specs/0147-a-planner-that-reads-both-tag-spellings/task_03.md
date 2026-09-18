---
task: task_03
spec: 0147-a-planner-that-reads-both-tag-spellings
status: completed
type: backend
complexity: medium
---

# Task 03: Refuse the ambiguity and keep the spelling

## Overview

The command turns the ambiguity report into the preflight refusal it already
uses, and renders every proposal in the spelling the selection resolved to.

## Requirements

1. MUST refuse in preflight when the highest version is ambiguous, naming both
   refs and the `--from` selector that resolves it.
2. MUST propose no version and approve nothing in that refusal, and MUST use the
   exit code preflight refusals already use.
3. MUST render the proposed version, the approval question and reset planning in
   the spelling of the selected tag.
4. MUST keep reset planning's inventory matching Releases by exact tag name, and
   MUST keep each digest entry's spelling, ref and commit.
5. MUST leave every existing state, exit code and message unchanged for a
   repository with one spelling.

## Subtasks

- [ ] Wire the ambiguity report to the preflight refusal.
- [ ] Render proposals in the selected spelling.
- [ ] Keep the reset inventory and digest entries distinct per ref.
- [ ] Cover the refusal, a bare-tag plan and an unchanged prefixed plan.

## Acceptance Criteria

- [ ] An ambiguous repository refuses with both refs and the selector named, and
      no proposed version.
- [ ] A bare-tagged repository plans, with base and proposal both bare.
- [ ] A prefixed repository plans exactly as it does today.
- [ ] Reset digest entries keep spelling, ref and commit.

## Context

- interface: `internal/cli/releaseplan_command.go`
- interface: `internal/cli/releaseplan_git_source.go`

## Verification

- `out="$(go test -count=1 -run "^TestReleasePlanRefusesAmbiguousHighestVersion$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.
- `out="$(go test -count=1 -run "^TestReleasePlanReadsBareTags$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 3-6; User Stories 3-4; Goals 3-4;
Success Metrics 2-3; Regression locks;
`_techspec.md` → Implementation Design: Refusing the ambiguity, Keeping the
spelling; API Contracts 2-3; Build Order 3.

## Result

Implementation:

- Default range discovery now inventories every reachable tag, parses both
  stable spellings, and delegates highest-version choice to the shared
  selector. An ambiguous highest version remains a typed preflight error and
  adds the existing `--from` selector as its next action.
- Version increments preserve the parsed prefix bit, so range proposals,
  approval data and approval questions use the selected base tag's spelling.
- Reset planning accepts and renders a bare target. Its command coverage now
  inventories bare and prefixed aliases separately, matches Releases by exact
  tag name, and observes distinct target commits. Digest coverage proves that
  changing a tag's spelling, ref or target commit changes the plan digest.

Focused checks and acceptance evidence:

- Ambiguous refusal: `GOCACHE=/tmp/roundfix-task03-go-cache rtk go test
  ./internal/cli -run '^TestReleasePlan'` passed 59 tests, including
  `TestReleasePlanRefusesAmbiguousHighestVersion`; it asserts the preflight exit
  code, both refs, `--from`, empty stdout and no proposal or approval text.
- Bare plan and unchanged prefixed plan: the same 59-test check includes
  `TestReleasePlanReadsBareTags`, which observes base `1.4.2`, proposal `1.5.0`
  and its bare approval question, plus the existing prefixed range and outcome
  cases unchanged.
- Reset spelling and exact inventory matching: the same check includes
  `TestReleasePlanResetTextAndJSONInventoryMatchThroughRunBoundary`; it observes
  the bare reset target and approval, retains both `0.2.0` and `v0.2.0`, and
  resolves their Releases to their distinct exact-tag commits.
- Digest entries: `GOCACHE=/tmp/roundfix-task03-go-cache rtk go test
  ./internal/releaseplan` passed 127 tests, including independent digest
  mutations for tag spelling, ref and target commit.
- `rtk git diff --check` passed.

Broader check note:

- `GOCACHE=/tmp/roundfix-task03-go-cache rtk go test ./internal/cli` exercised
  1,150 tests; 1,148 passed and two unrelated force-stop integration tests in
  `orphan_unix_test.go` failed. This Task does not change that surface. The
  Daemon-owned Verification commands were not run in this turn.
