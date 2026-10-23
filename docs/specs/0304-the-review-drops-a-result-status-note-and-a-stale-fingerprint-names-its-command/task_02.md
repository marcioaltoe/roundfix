---
task: task_02
spec: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
status: pending
type: test
complexity: low
---

# Task 02: A stale Behavior Surface fingerprint names the command that re-records it

## Overview

`TestTheSkillCoverageMapIsCurrent` reports a stale Behavior Surface Record
entry as `<id>: stale fingerprint` and nothing else, so on 2026-10-22 a Spec
0298 item parked `gate-failed` and a Spec 0299 Run ended Unresolved after a
guide edit, with no pointer to the fix. This Task makes the problem name the
re-record command and holds that command equal to the regeneration
`.roundfixrc.yml` declares for the record (`_prd.md` → Core Feature 4).

## Requirements

1. MUST answer the adopted Backlog Entry
   [`references/2026-10-22-a-stale-behavior-surface-fingerprint-does-not-name-its-re-record-command.md`](references/2026-10-22-a-stale-behavior-surface-fingerprint-does-not-name-its-re-record-command.md)
   of 2026-10-22.
2. MUST add to `internal/docscontract/skill_coverage_test.go` the constant
   `behaviorSurfaceRecordCommand` with the value of `_techspec.md` →
   Interfaces, and make `skillCoverageProblems` report a stale entry as the
   text of API Contract 5: `<id>: stale fingerprint; re-record
   docs/references/behavior-surfaces.json after the last edit of its sources
   with: ` followed by the constant, the record path written through
   `skillcoverage.RecordPath`. Every other problem the function reports keeps
   its bytes.
3. MUST change the `stale` case of `TestAStaleBehaviorSurfaceIsReported` to
   expect that text and leave its `missing` and `extra` cases byte-identical
   (the declared break of this Task).
4. MUST add `TestTheStaleFingerprintProblemNamesTheProjectRegeneration` to
   the same file, calling `t.Parallel()` first, which holds `_techspec.md` →
   API Contract 6: it reads `.roundfixrc.yml` at the repository root found by
   `baselineDocumentationRepoRoot`, parses it with `gopkg.in/yaml.v3`, and
   requires the `regenerate` values of the `delivery.derived_paths` entries
   whose `paths` hold `skillcoverage.RecordPath` to be exactly
   `[behaviorSurfaceRecordCommand]`; it also requires the flag
   `record-skill-coverage` to be registered and the command to contain `
   -record-skill-coverage` and `-run '^TestTheSkillCoverageMapIsCurrent$'`.
5. MUST NOT add a file to `internal/docscontract`, change `.roundfixrc.yml`,
   `go.mod`, `go.sum`, `docs/references/behavior-surfaces.json` or any other
   test of the package.

## Subtasks

- [ ] Add the command constant and the hint to the stale problem.
- [ ] Update the stale case and add the contract test against `.roundfixrc.yml`.

## Acceptance Criteria

- [ ] A stale fingerprint problem names the record path and the re-record command; the missing and extra problems keep their bytes.
- [ ] The command in the problem equals the `.roundfixrc.yml` regeneration for the record, and a different constant fails the contract test.

## Consumers

None. The change lives in test-only code of `internal/docscontract`; no production code or other suite reads the problem text.

## Context

- interface: `internal/docscontract/skill_coverage_test.go`

## Verification

- `go test -tags docscontract -list '^TestTheStaleFingerprintProblemNamesTheProjectRegeneration$' ./internal/docscontract | grep -qx 'TestTheStaleFingerprintProblemNamesTheProjectRegeneration' || { echo 'missing contract test' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^(TestAStaleBehaviorSurfaceIsReported|TestTheStaleFingerprintProblemNamesTheProjectRegeneration|TestTheSkillCoverageMapIsCurrent)$' ./internal/docscontract` — expected: exit 0. Before this Task the contract test does not exist, so the command fails at its first check. Afterwards the stale case expects the named command, the contract test holds it equal to `.roundfixrc.yml`, and the Skill Coverage Map stays current.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 5
- `_techspec.md` → Current behavior; API Contracts 5-6; Invariant 6; Testing Approach; Build Order 2
