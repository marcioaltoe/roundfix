---
task: task_02
spec: 0174-operator-surfaces-that-tell-the-truth
status: pending
type: backend
complexity: medium
---

# Task 02: A help token means help only where it is an argument

## Overview

`commandWantsHelp` in `internal/cli/cli.go` returns true when `-h`, `--help` or a bare `help` appears anywhere in a command's arguments, and 47 call sites print usage and exit `0` on it. So `roundfix archive --bogus help` exits `0` where the same line without `help` exits `2`, `roundfix implement --spec help` prints usage, and no flag can take the value `help`. The arguments come from the operator's or an Agent's command line; the exit code and stdout are read by scripts and Supervisors that treat `0` as success.

## Requirements

1. MUST change `commandWantsHelp` so that `help` counts only as the first argument it receives, `-h` and `--help` count at any position before the first `--`, and no token after `--` counts. Its callers, their usage text and the top-level dispatch in `run` stay unchanged.
2. MUST keep subcommand help working through the existing `args[1:]` hand-off: `roundfix runs list help` and `roundfix skills check help` still print usage and exit `0`.
3. MUST add to `internal/cli/cli_test.go` one test per case, asserting the exit code and whether stdout carries the command's usage: `roundfix archive --bogus help` exits `2` with no usage; `roundfix implement --spec help` prints no usage and exits with the same non-zero code as `roundfix implement --spec no-such-spec` in the same workspace; `roundfix archive some-slug help` exits `2` with no usage; `roundfix archive -- --help` exits `2` with no usage; `roundfix archive help` exits `0` with the archive usage; `roundfix runs list help` exits `0` with the runs usage; `roundfix archive some-slug --help` exits `0` with the archive usage; and a table test of `commandWantsHelp` over those argument lists.
4. MUST state in the Global contract of `docs/user-guide/commands.md` that a bare `help` requests usage only as a command's first argument and that `-h` or `--help` requests it anywhere before `--`, using the phrase `only as a command's first argument`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A `help` that is a flag value, a trailing positional or follows `--` is not a help request, and the command fails as it would without it.
- [ ] A leading `help`, a subcommand's leading `help` and a `--help` after an argument still print usage and exit `0`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `docs/user-guide/commands.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCommandWantsHelpReadsOnlyArgumentPositions|TestHelpTokenAsAFlagValueIsNotHelp|TestHelpTokenAsASpecValueIsNotHelp|TestTrailingHelpPositionalIsNotHelp|TestHelpAfterTheTerminatorIsNotHelp|TestLeadingHelpTokenPrintsUsage|TestSubcommandLeadingHelpTokenPrintsUsage|TestDashHelpAfterAnArgumentPrintsUsage|TestRunCommandHelp|TestRunHelp)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCommandWantsHelpReadsOnlyArgumentPositions TestHelpTokenAsAFlagValueIsNotHelp TestHelpTokenAsASpecValueIsNotHelp TestTrailingHelpPositionalIsNotHelp TestHelpAfterTheTerminatorIsNotHelp TestLeadingHelpTokenPrintsUsage TestSubcommandLeadingHelpTokenPrintsUsage TestDashHelpAfterAnArgumentPrintsUsage TestRunCommandHelp TestRunHelp; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && grep -q "only as a command's first argument" docs/user-guide/commands.md` — expected: exit 0; before this Task none of the new named tests exists and the guide states no help-token rule, so the command fails.

## References

- [_techspec.md](_techspec.md) — The help predicate
- `_prd.md` → Goal 3; Core Feature 2; Success Metric 2
- `_techspec.md` → API Contract 5; Testing Approach 2
