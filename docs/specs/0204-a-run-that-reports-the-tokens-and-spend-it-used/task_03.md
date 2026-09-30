---
task: task_03
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: pending
type: backend
complexity: high
---

# Task 03: `runs show`, the Implement Run summary and `deliver status` print the usage

## Overview

task_02 records usage and sums it, and no command prints it. This Task adds
one renderer and three readers: the new read-only `roundfix runs show`, a
`Tokens:` line at the end of the Implement Run summary, and a `Usage:` line in
`deliver status`, whose limits line stops saying `spend not measured`. It also
writes the user guide's account of what is measured, including the metering
gateway as an operator option. It is verifiable on its own through the three
commands on a temporary Roundfix Home.

## Requirements

1. MUST add `formatTokenTotals` in the new file `internal/cli/token_usage_render.go` with the two phrases of `_techspec.md` → Surfaces, and MUST use it for every line this Task prints, so the three surfaces cannot disagree.
2. MUST add `roundfix runs show <run-id> [--json]` in the new file `internal/cli/runs_show.go`, dispatched from `runRunsCommand`, with the text of Surface Transcripts 1 and 2 and the JSON and exit codes of API Contract 1. It MUST open the Run Database with the reader, write nothing, and refuse a missing or extra argument as a usage error with exit `2`.
3. MUST print, on every outcome of the Implement Command, the `Tokens:` line of Surface Transcript 3 right after the outcome line and before any push line, reading the Run's sums from the store after the Run ends. Exit codes MUST not change.
4. MUST make the limits line printed by `deliver start` and `deliver status` end with `tokens <n>` or `tokens none` from `DeliveryQueueLimits.MaxTokens` instead of `spend not measured`, and make `deliver status` print the `Usage:` line of `_techspec.md` → Surfaces right after it, as Surface Transcript 5 and the limits and usage lines of Surface Transcript 4 show. `deliver start` MUST print no `Usage:` line.
5. MUST NOT print an unreported prompt as `0` on any of the three surfaces. One test MUST seed a Run whose prompts are all unreported and prove that `runs show`, the Implement Run summary and `deliver status` each say `none reported` and that none of them prints `0 tokens`.
6. MUST add `roundfix runs show <run-id> [--json]` to the `runs` usage text and its own help text in `internal/cli/cli.go`, keeping every string `TestRunCommandHelp` requires, so `internal/cli/cli_test.go` does not change.
7. MUST characterize before changing: the Result MUST name every existing test in `internal/cli/implement_test.go`, `internal/cli/orphan_unix_test.go`, `internal/cli/deliver_test.go`, `internal/cli/deliver_limits_test.go` and `internal/cli/deliver_revalidate_test.go` that compares the whole Implement stdout or the limits line, and update exactly those to the new lines. No existing test is renamed or removed.
8. MUST put the new tests in `internal/cli/runs_show_test.go`, `internal/cli/implement_tokens_test.go` and `internal/cli/deliver_usage_test.go`, with a temporary Roundfix Home and repository and a fake runner that reports usage where needed.
9. MUST document `runs show`, its JSON schema `roundfix/runs-show/v1` and its exit codes in `docs/user-guide/commands/runs.md`; the `Tokens:` line in `docs/user-guide/commands/implement.md`; the limits and `Usage:` lines in `docs/user-guide/commands/deliver.md`; and add a `## Token usage` section to `docs/user-guide/usage.md` that states each adapter's basis (`request-sum` for Codex, `turn` for the others), that unreported prompts are never zero, that spend is only what an adapter reported, and how an operator routes Codex through a metering gateway: that Roundfix neither starts nor reads it, the failure measured on 2026-09-29, and that Node.js 26 offers HTTP/2 by default while forcing HTTP/1.1 takes an `allowH2: false` dispatcher in the gateway itself.
10. MUST prove each new gate can fail. The Result MUST record one sabotage of the renderer (for example printing an unreported total as `0 tokens`) and one of `runs show` (for example dropping the unreported scope line), each with the test that failed, and that the code was restored.

Spec 0194 creates the per-command guide and skill files this Task edits
(`_prd.md` → Prerequisites). They are declared under `creates:` because they do
not exist when this Spec is authored. This Task MUST NOT create one that Spec
0194 has not created: when such a file is absent, the Task stops and reports
that Spec 0194 has not landed.

## Subtasks

- [ ] Characterize the existing whole-output tests and name them.
- [ ] Add the renderer and `runs show` with its help and JSON.
- [ ] Add the Implement `Tokens:` line and the `deliver` limits and `Usage:` lines.
- [ ] Write the command guides and the token usage section.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] `roundfix runs show` reproduces Surface Transcripts 1 and 2, and `--json` reports schema `roundfix/runs-show/v1`.
- [ ] The Implement Run summary ends with the `Tokens:` line of Surface Transcript 3 on every outcome.
- [ ] `deliver start` prints Surface Transcript 5, and `deliver status` prints the limits and `Usage:` lines, `Usage: no Runs recorded` for a queue without Runs.
- [ ] No surface prints an unreported prompt as `0 tokens`.
- [ ] No non-test source under `internal/cli` still prints `spend not measured`.

## Context

- creates: `internal/cli/token_usage_render.go`
- creates: `internal/cli/runs_show.go`
- creates: `internal/cli/runs_show_test.go`
- creates: `internal/cli/implement_tokens_test.go`
- creates: `internal/cli/deliver_usage_test.go`
- interface: `internal/cli/runs.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/cli/orphan_unix_test.go`
- interface: `internal/cli/deliver_test.go`
- interface: `internal/cli/deliver_limits_test.go`
- interface: `internal/cli/deliver_revalidate_test.go`
- creates: `docs/user-guide/commands/runs.md`
- creates: `docs/user-guide/commands/implement.md`
- creates: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/usage.md`
- instruction: `internal/cli/cli_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRunsShowPrintsEachScopeAndTheTotal|TestRunsShowRefusesAnUnknownRun|TestRunsShowJSONIsSchemaV1|TestRunsShowRefusesAMissingOrExtraArgument|TestRunsShowWritesNothing|TestImplementSummaryEndsWithTheTokensLine|TestImplementSummaryPrintsRecordedUsage|TestDeliverStartPrintsTokensNone|TestDeliverStatusPrintsTheUsageLine|TestDeliverStatusPrintsNoRunsRecorded|TestUnreportedUsageIsNeverPrintedAsZero|TestRunCommandHelp)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunsShowPrintsEachScopeAndTheTotal TestRunsShowRefusesAnUnknownRun TestRunsShowJSONIsSchemaV1 TestRunsShowRefusesAMissingOrExtraArgument TestRunsShowWritesNothing TestImplementSummaryEndsWithTheTokensLine TestImplementSummaryPrintsRecordedUsage TestDeliverStartPrintsTokensNone TestDeliverStatusPrintsTheUsageLine TestDeliverStatusPrintsNoRunsRecorded TestUnreportedUsageIsNeverPrintedAsZero TestRunCommandHelp; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task eleven of the twelve named tests do not exist, so the command fails.
- `matches="$(grep -rlF --include='*.go' --exclude='*_test.go' "spend not measured" internal/cli)"; test -z "$matches" || { printf 'still prints spend not measured: %s\n' "$matches" >&2; exit 1; }` — expected: exit 0; before this Task `internal/cli/deliver.go` prints the phrase, so the command fails.
- `for pair in "docs/user-guide/commands/runs.md|roundfix runs show" "docs/user-guide/commands/runs.md|roundfix/runs-show/v1" "docs/user-guide/commands/implement.md|Tokens:" "docs/user-guide/commands/deliver.md|Usage:" "docs/user-guide/usage.md|request-sum" "docs/user-guide/usage.md|none reported" "docs/user-guide/usage.md|metering gateway" "docs/user-guide/usage.md|allowH2: false"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task no guide names `roundfix runs show` or the token usage section, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 2 and 4; User Stories 1, 3 and 5; Core Features 3, 4, 7, 8, 9 and 11; Success Metrics 2, 3 and 4; Declared breaks; Acceptance evidence
- [_techspec.md](_techspec.md) — Surfaces; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 5; API Contract 1; API Contract 2; API Contract 3; Testing Approach 4; Testing Approach 6; Build Order 3
- ADR-0198
