---
task: task_02
spec: 0132-a-grant-read-exactly-where-it-lives
status: completed
type: backend
complexity: medium
---

# Task 02: Close the frontmatter on a complete marker line

## Overview

A purported closing line such as `---evil` currently closes the frontmatter by
prefix, leaving the rest of the line discarded as body while every other field
parses cleanly, so a malformed record resolves to a grant. This slice makes the
delimiter exact: the marker and nothing else. It is verifiable alone — the
malformed record refuses and the well-formed one still grants.

## Requirements

1. MUST treat a frontmatter delimiter as a complete line: the marker alone,
   allowing only a trailing carriage return.
2. MUST classify a record whose opening or closing delimiter line carries extra
   characters as malformed, resolving to a refusal that names the offending
   line, never to a grant.
3. MUST keep a well-formed record resolving exactly as it does today, including
   its fields, paths and operations.
4. MUST keep an unreadable record and an unavailable revision reporting
   unresolved rather than refused, so the new refusal fires on a positive
   observation of malformation and never on missing evidence.
5. MUST NOT change any other parse rule, field name or reason code.
6. MUST update only the characterization row this Task intentionally moves.

## Subtasks

- [ ] Require a complete marker line at both delimiters.
- [ ] Classify a malformed delimiter as a refusal naming the line.
- [ ] Prove a well-formed record is unchanged.
- [ ] Move the one recorded characterization row.

## Acceptance Criteria

- [ ] A record whose closing line is `---evil` refuses, and the refusal names
      that line; it granted before this Task.
- [ ] A record whose opening line carries extra characters refuses the same way.
- [ ] A well-formed record grants the identical fields, paths and operations it
      granted before.
- [ ] A record with a trailing carriage return on its delimiter still parses.
- [ ] An unreadable record still reports unresolved, distinguishable from
      refused.

## Context

- interface: `internal/authorization/authorization.go`

## Verification

- `grep -q 'func TestMalformedDelimiterRefuses' internal/authorization/authorization_test.go && go test -count=1 ./internal/authorization -run '^TestMalformedDelimiterRefuses$'` — a malformed delimiter refuses and names the line; this fails today because the record grants.
- `grep -q 'func TestMalformedDelimiterRefuses' internal/authorization/authorization_test.go || exit 1; go test -count=1 ./internal/authorization -run '^TestAuthorizationParseCharacterization$'` — the recorded answers still hold with exactly the one intended row moved.
- `grep -q 'func TestMalformedDelimiterRefuses' internal/authorization/authorization_test.go || exit 1; go test -count=1 ./internal/authorization ./internal/spec ./internal/speccheck ./internal/suiteguardcontract` — every consumer of the reader still passes.

## References

- `_prd.md` → Goals 1; User Stories 1; Core Features 1; Decisions: Declared intentional breaks 1.
- `_techspec.md` → Implementation Design: Exactness at the parse boundary; Build Order 2.

## Result

Implementation evidence:

- The authorization reader now recognizes opening and closing frontmatter
  markers only when the complete observed line is `---`. A line beginning with
  that marker and carrying extra characters returns the existing
  `malformed_record` refusal and includes the offending line in its detail.
- `TestMalformedDelimiterRefuses` covers `---evil` at both delimiter positions;
  both cases refuse and name `---evil`.
- `TestAuthorizationParseCharacterization` moves only the recorded
  malformed-closing answer from granted to refused. Its well-formed row still
  asserts the same status, action, consumer, bounded path, operation list and
  `Permits` answers.
- `TestDelimiterWithTrailingCarriageReturnGrants` writes both marker lines with
  CRLF endings and observes the unchanged grant.
- `TestMissingAuthorizationEvidenceRemainsUnresolved` distinguishes a missing
  record (`unreadable_record`) and a missing Git revision
  (`unavailable_revision`) from the positive malformed-record refusal.
- The parser retains the existing field names, reason codes, frontmatter field
  validation and legacy fallback rules.

Focused-check evidence:

- Before the parser edit,
  `rtk go test -count=1 -v ./internal/authorization -run TestMalformedDelimiterRefuses`
  failed both cases: the opening refusal omitted `---evil`, and the malformed
  closing resolved to `granted`. The first sandboxed attempt could not access
  the standard Go build cache; the unchanged cache-enabled retry produced the
  behavioral failures.
- After the final Go edits,
  `rtk go test -count=1 -v ./internal/authorization -run 'Test(AuthorizationParseCharacterization|DelimiterWithTrailingCarriageReturnGrants|MissingAuthorizationEvidenceRemainsUnresolved|MalformedDelimiterRefuses)'`
  passed 10 tests.
- `rtk go test -count=1 ./internal/authorization` passed all 11 package tests.
- The host Go 1.27.1 formatter rejected two unchanged CLI test files, so the
  repository incremental gate was rerun with its compatible Go 1.26.7
  toolchain. The network-restricted attempt reached the skill-sync stage and
  was blocked from `api.github.com`; the unchanged access-enabled retry of
  `rtk make verify-incremental` passed all tests, skill checks and the build.
- A later compatible-toolchain incremental run after the final test cleanup
  timed out in three unchanged daemon concurrency tests while the
  authorization package and its consumers passed. The same three tests passed
  together in isolation with
  `rtk go test -count=1 -v ./internal/daemon -run '^(TestTaskCycleIntegratedVerificationCapacityOneBoundsConcurrentTaskWorktrees|TestTaskCycleSchedulesIndependentWaveWithConcurrencyCap|TestTaskCycleParallelTaskPromptUsesTaskWorktreeContextBase)$'`;
  one unchanged incremental retry then passed all tests, skill checks and the
  build.
- The Task's declared `## Verification` commands were not run; Daemon
  Verification owns them after this handoff.

Follow-up note:

- The full-suite-load daemon timeouts are outside this parser slice. This Task
  did not change their timing, fixtures or assertions.
