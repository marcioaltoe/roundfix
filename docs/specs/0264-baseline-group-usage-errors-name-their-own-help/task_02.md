---
task: task_02
spec: 0264-baseline-group-usage-errors-name-their-own-help
status: pending
type: backend
complexity: low
---

# Task 02: Baseline usage errors name the group once and their own help

## Overview

A bare `baseline capabilities`, `baseline skills` or `baseline assets`, and an
unknown `baseline` command, print a usage error that doubles the word
`baseline` and points at `baseline profile --help` (`_prd.md`; the adopted
entry
[2026-10-12-baseline-group-usage-errors-name-the-wrong-command.md](references/2026-10-12-baseline-group-usage-errors-name-the-wrong-command.md),
recorded 2026-10-12). This Task gives the `baseline` dispatcher a usage-error
printer of its own, asserts the six Surface Transcripts, and adds a class
guard over the usage error of every group of the root usage. It is
implemented on the tree after Spec 0260 merges and is verifiable on its own
through the named tests.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12,
   `2026-10-12-baseline-group-usage-errors-name-the-wrong-command.md`: the
   three groups follow `_techspec.md` → API Contract 1 and Surface
   Transcripts 1-5, and an unknown `baseline` command follows API Contract 2
   and Surface Transcript 6.
2. MUST add `printBaselineUsageError` with the signature in `_techspec.md` →
   Interfaces and call it at the five usage-error sites of
   `runBaselineCommand`, below the help branches Spec 0260 adds; the help
   branches MUST stay as Spec 0260 left them.
3. MUST keep `printBaselineProfileFailure`, its text and every
   `baseline profile` caller unchanged, write nothing to standard output on a
   `baseline` usage error, and keep exit 2 (Invariants 1, 2 and 5).
4. MUST add `internal/cli/baseline_group_usage_test.go` with
   `TestBaselineGroupUsageTranscripts`, which asserts Surface Transcripts 1-6
   in subtests named `transcript_1` to `transcript_6` as `_techspec.md` →
   Testing Approach states, and `TestEveryGroupUsageErrorPointsAtItsCommands`,
   which follows API Contract 3 and reuses `helpCommandPaths` from
   `internal/cli/command_help_test.go`. Each calls `t.Parallel()` as its
   first statement (ADR-0259).
5. MUST change in `internal/cli/command_help_test.go` only the expected
   standard error of the bare and unknown `capabilities`, `skills` and
   `assets` rows of `TestCommandHelpArgumentBoundaries`, which Spec 0260
   pinned to today's
   `roundfix: baseline baseline <group> failed: ... Run 'roundfix baseline profile --help' for usage.`
   text, so they expect API Contract 1. Every other row and the
   `helpCommandPaths` helper keep their bytes.
6. MUST NOT change `internal/cli/cli_test.go`, any help or usage text, the
   user guide, `docs/references/behavior-surfaces.json`, or any other
   command's behavior.

## Subtasks

- [ ] Add the printer and route the five usage-error sites through it.
- [ ] Add the transcript test with its pointer checks.
- [ ] Add the class guard over every group of the root usage.

## Acceptance Criteria

- [ ] `TestBaselineGroupUsageTranscripts` passes with all six transcripts and
      their pointers resolving to a help that exits 0.
- [ ] `TestEveryGroupUsageErrorPointsAtItsCommands` passes over at least 12
      groups and names the group, word or pointer on failure.
- [ ] The existing `baseline_profile_test.go` cases pass unchanged.

## Consumers

- Standard error of a `baseline` group usage error — consumers: `./internal/cli`

## Context

- instruction: `docs/agents/cli.md`
- instruction: `docs/adr/0184-a-techspec-states-a-command-surface-as-a-transcript.md`
- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- interface: `internal/cli/baseline_profile.go`
- interface: `internal/cli/command_help_test.go`
- creates: `internal/cli/baseline_group_usage_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestBaselineGroupUsageTranscripts|TestEveryGroupUsageErrorPointsAtItsCommands)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineGroupUsageTranscripts TestBaselineGroupUsageTranscripts/transcript_1 TestBaselineGroupUsageTranscripts/transcript_2 TestBaselineGroupUsageTranscripts/transcript_3 TestBaselineGroupUsageTranscripts/transcript_4 TestBaselineGroupUsageTranscripts/transcript_5 TestBaselineGroupUsageTranscripts/transcript_6 TestEveryGroupUsageErrorPointsAtItsCommands; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two tests do not exist, so no PASS line is printed.
- `grep -qF -- 'func printBaselineUsageError(' internal/cli/baseline_profile.go || { printf 'the baseline usage-error printer is missing\n' >&2; exit 1; }; go test -count=1 -run '^(TestCapabilityRecheck|TestBaselineProfileCommandInitShowAndValidate|TestBaselineProfileCommandRejectsInvalidAndNonRepositoryProfiles|TestBaselineProfileHelpContract|TestCommandHelpArgumentBoundaries)$' ./internal/cli` — expected: exit 0; before this Task the printer does not exist.

## References

- [_prd.md](_prd.md) — Goals; User Stories 1-3; Core Features 1-3; Success Metric 1; Success Metric 2
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Invariants 1-5; Testing Approach; Build Order 2
- [2026-10-12-baseline-group-usage-errors-name-the-wrong-command.md](references/2026-10-12-baseline-group-usage-errors-name-the-wrong-command.md)
- ADR-0184; ADR-0259
