---
task: task_03
spec: 0204-a-run-that-reports-the-tokens-and-spend-it-used
status: completed
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

## Result

Implemented the Task 03 reader and output slice. Task status remains
Daemon-owned; no declared Verification command was run, and no commit, push
or Pull Request was made. The initial worktree change was the Daemon's
`status: pending` to `status: in_progress` change in this Task file.

### Implementation and acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| `runs show` reproduces the scope/total and unknown-Run transcripts, with schema v1 JSON | Added `runs_show.go`, dispatched from `runRunsCommand`, using `store.OpenReader` through the existing reader helper. `TestRunsShowPrintsEachScopeAndTheTotal` checks the complete text including the unreported QA scope; `TestRunsShowRefusesAnUnknownRun` checks the complete stderr and exit 2; `TestRunsShowJSONIsSchemaV1` checks identity, counts, sums and explicit null split fields; `TestRunsShowRefusesAMissingOrExtraArgument` checks usage errors. `TestRunsShowWritesNothing` snapshots durable Home files and rejects non-empty WAL writes, allowing SQLite's read-only shared-memory coordination and empty WAL sidecars. `TestRunsShowPrintsNoPromptsRecorded` covers Runs with no rows. |
| Every terminal Implement outcome has `Tokens:` immediately after its outcome | Added reads of `RunTokenUsage` after completion, using an uncancelled context and the shared renderer. Existing outcome paths print the line before any push line; a terminal failure fallback covers setup, infrastructure and push failures. Its counts come from the task cycle. `TestImplementSummaryEndsWithTheTokensLine` exercises Clean, Unresolved and Stopped. `TestImplementTokenSummaryCoversSetupIntegrationAndPush` exercises Failed before prompts, Integration Pending, Budget Exceeded after the cycle, push failure and successful push, including line order and unchanged exit codes. `TestImplementSummaryPrintsRecordedUsage` uses a fake runner reporting two 5,639,755-token prompts and compares the 11,279,510-token Implement total with `runs show`. No-Run preflight, already-completed and detached startup reports retain their contracts. |
| Delivery limits and status usage are printed, including an empty queue's usage | `printDeliveryLimits` reads `MaxTokens` and prints `tokens none` or the ceiling. Only status reads `DeliveryQueueTokenUsage`, immediately after limits. `TestDeliverStartPrintsTokensNone` checks the entire start output and absence of a Usage line. `TestDeliverStatusPrintsTheUsageLine` seeds a 5,000,000 ceiling and replaces an item's current Run ID, proving both Run links still contribute: 9,659,850 tokens from 3 of 4 prompts across 2 Runs. `TestDeliverStatusPrintsNoRunsRecorded` checks the entire empty-usage status output. |
| Unreported usage is never printed as zero | `formatTokenTotals` handles no rows, no reports and reported totals distinctly, and supplies token/cost phrases for all three surfaces. `TestUnreportedUsageIsNeverPrintedAsZero` creates one real CLI Implement Run with two unreported prompts, links that Run to a temporary queue, then reads the same Run through Implement, `runs show` and `deliver status`: every surface says `none reported by 2 prompt(s)` and none says `0 tokens`. `TestTokenTotalsRenderReportedZeroAndCurrencies` separately preserves a genuinely reported zero and checks two-currency rendering/rounding. |
| No production CLI source prints `spend not measured` | A focused `rtk proxy rg -n 'spend not measured' internal/cli -g '*.go' -g '!**/*_test.go'` returned no matches (rg exit 1). |

The Spec 0194 command guides were present before editing. Updated the runs
guide with text output, schema `roundfix/runs-show/v1`, fields, null semantics
and exit codes; the Implement guide with `Tokens:`; the delivery guide with
limits and `Usage:`; and `usage.md` with adapter bases, missing reports,
adapter-only cost, and the operator-owned metering gateway. The gateway
failure is the authored 2026-09-29 measurement in the Spec/ADR, not a new
measurement. The transport caveat was checked against Node.js 26's release
page, Undici PR 4828 and its Connector documentation, linked in the guide.
A Python local-link check passed for all four edited guides.

### Characterization of existing output assertions

Read the existing tests before updating each affected expectation; no
existing test was renamed or removed. The intentional contract now prints
usage for failed Runs too, replacing the three old empty-stdout expectations.
These tests in `implement_test.go` compare complete terminal stdout and were
updated only for the new outcome/usage lines:

- `TestRunImplementExecutesSpecEndToEnd`
- `TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow`
- `TestRunImplementTemporaryVerificationFlowRetriesOnceWithoutAgentRepair`
- `TestRunImplementTemporaryVerificationFlowPreservesDeterministicRepair`
- `TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks`
- `TestRunImplementBootstrapFailureEndsFailedBeforeAgentWork`
- `TestRunImplementAutoPushFailureEndsFailedAndJournalsPush`
- `TestRunImplementStopRequestEndsStoppedWithInterruptMapping`
- `TestRunImplementDatabaseStopRequestAfterTaskCommitEndsStoppedAndReleasesLock`
- `TestRunImplementQAVerdictMatrix`
- `TestRunImplementQAOnlyRunSettlesOutcomeFromVerdict`
- `TestRunImplementInfrastructureFailureEndsFailed`

`TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree`
compares the terminal suffix rather than complete stdout; its expected suffix
also now includes the Tokens line. Complete detached-startup/no-Run/preflight
assertions in `TestRunImplementDetachPrintsReportAndCompletesRun`,
`TestRunImplementDetachReportsAndRelaysPreflightFailure`,
`TestRunImplementAllTasksCompletedReportsWithoutRun` and the refusing
Implement case in `TestAgentSelectionProfilesMacro` are unchanged because
those paths create no foreground terminal summary.

Other unchanged empty-stdout preflight assertions in `implement_test.go`
were also characterized; these create no Run and retain their expectations:
`TestRunImplementRemovedQAFlagExplainsTaskGraph`,
`TestRunImplementValidationFailures`,
`TestRunImplementRejectsInvalidVerificationCapacityBeforeRunCreation`,
`TestRunImplementPreflightFailures`, `TestImplementRunWindow`,
`TestImplementRefusesWhenCarryForwardIsAvailable`,
`TestImplementRejectsInvalidTaskTypeBeforeSideEffects`,
`TestRunImplementPreflightRejectsActiveRunInWorkingTree`,
`TestRunCeilingRefusalNamesTheWayOut`,
`TestImplementProfilePreflightFailureCreatesNoRunWorktreeOrAgentPrompt`,
`TestProfilePreflightRefusesUnhonouredAccessPolicy`,
`TestRunImplementSelectionFailureReportsProfileRemediationWithoutCreatingRun`,
`TestRunImplementSelectionFailureDoesNotPromptForDynamicFallback`,
`TestRunImplementPassesOneRunSelectionOverridesToPreflight`,
`TestRunImplementRejectsExplicitEmptySelectionOverrides`, and
`TestRunImplementSelectionOverrideRejectsPartialBeforeConfigLoad`.

The named delivery files had these complete-output/limits assertions, all
updated to the new limits and, for status, Usage line:

- `deliver_test.go`: `TestDeliverStatusPrintsTheItemWorktree` (both outputs).
- `deliver_limits_test.go`: `TestDeliverStartRecordsAndPrintsItsLimits`,
  `TestDeliverStartRecordsNoneForOmittedLimits`.
- `deliver_revalidate_test.go`: `TestDeliverStatusPrintsAnItemWarning`,
  `TestDeliverStatusPrintsNoWarningLineWithoutAWarning`.

`orphan_unix_test.go` has no whole Implement-summary or limits comparison.
Its Implement tests assert substrings; the complete empty-output assertion
in `TestReopenRefusesAnActiveRunWithoutReclaimingIt` is a Reopen refusal and
stays unchanged. No change to `orphan_unix_test.go` or `cli_test.go` was needed.
The additional ordinary path `deliver_park_status_test.go` also compares
whole status output: `TestDeliverStatusPrintsAParkLinePerParkedItem` and
`TestDeliverStatusReproducesSurfaceTranscriptOne` were characterized and their
limits/Usage expectations updated for this same contract change. Other
existing tests and expectations remain intact.

### Focused checks and mutation evidence

- Before production edits, `rtk proxy go test ./internal/cli -run
  '^TestRunsShow' -count=1` failed on unknown `show` dispatch and absent JSON,
  establishing the red starting point.
- Focused `go test` selections covering the new surface tests, existing help,
  and characterized Implement/delivery expectations passed. Cache access
  later failed in the sandbox, so subsequent commands used
  `GOCACHE=/private/tmp/roundfix-task03-go-cache`.
- Renderer sabotage: temporarily replaced the unreported token phrase with
  `0 tokens`. `go test ./internal/cli -run
  '^TestUnreportedUsageIsNeverPrintedAsZero$' -count=1` exited 1 with that named
  test failing on the fabricated zero. Restored the renderer byte-for-byte.
- `runs show` sabotage: temporarily skipped scopes whose Tokens were nil.
  `go test ./internal/cli -run '^TestRunsShowPrintsEachScopeAndTheTotal$'
  -count=1` exited 1 because the unreported QA row disappeared. Restored
  `runs_show.go` byte-for-byte.
- After both restorations, `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk
  proxy go test ./internal/cli -run
  'Test(RunsShow|TokenTotalsRender|ImplementSummary|ImplementTokenSummary|UnreportedUsage|DeliverStartPrints|DeliverStatusPrints|RunCommandHelp)'
  -count=1` passed.
- First `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk make
  verify-incremental` attempt reached the full tests and exited 2: three
  failure-path stdout assertions still expected empty output, and two
  force-stop integration tests could not read the sandbox's process table.
  Updated only the three intentionally changed assertions, and made the
  failure fallback preserve task-cycle counts. The two force-stop tests
  passed when rerun with the required process-table permission; their source
  remains unchanged.
- Final focused and incremental check outcomes are recorded below after the
  retry. The Task's declared Verification remains for the Daemon.

No follow-up implementation was folded into this diff. Token-ceiling flag
parsing, enforcement and skill updates remain Task 04's slice.

### Final local check outcomes

- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk proxy go test
  ./internal/cli -run
  'Test(ImplementTokenSummary|RunImplementBootstrapFailure|RunImplementAutoPushFailure|RunImplementInfrastructureFailure)'
  -count=1`: exit 0 after updating the failure expectations and preserving
  task-cycle counts.
- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk make verify-incremental`:
  exit 0 with the required process-table permission. Formatting, vet, the
  repository tests (CLI ran freshly in 110.796s), skill synchronization and
  checks, and the binary build passed. Other unchanged package checks reused
  the incremental cache where applicable.
- Guide content assertions passed for command syntax, schema/exit-code
  documentation, Tokens/Usage output and the required measurement/gateway
  phrases. Local-link checks passed for the four guides.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0.
- Handoff scope contains only this Task's reader/renderer, CLI output/help,
  new tests, intentionally changed existing output expectations, four guides
  and this Result. The pre-existing Daemon-owned frontmatter is preserved;
  `_tasks.md`, other Task files and governed tooling/skills were not edited.

Implementation is handed back for Daemon Verification and settlement; this
Result makes no terminal Task verdict.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/deliver_park_status_test.go`

## Carry-forward provenance

- Source Run: `run_20261001T094428Z_f4d3047c47b08ac0`
- Source commit: `28be35fcd7212a7a65c3289e74d18d112e350f24`
