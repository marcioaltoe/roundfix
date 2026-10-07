---
task: task_03
spec: 0243-a-history-that-holds-only-records
status: completed
type: backend
complexity: high
---

# Task 03: roundfix history sanitize plans the history and applies one batch after the History Full Tag

## Overview

This Task adds the History Sanitize Command, `roundfix history sanitize`. It
orders the units of the existing history, prints the plan, adds Jev's advice
on request, and with `--apply --batch <n>` applies exactly the next `n`
units through task_01's conversion and task_02's kinds. It refuses to apply
on a dirty tree, or before the annotated `history-full` tag holds every path
the batch removes or rewrites. It never commits, tags or pushes; the
operator turns each batch into one Pull Request. It answers the Backlog
Entry "History keeps only what the Secondbrain needs" of 2026-10-06,
Expected 3, and the maintainer's decision "Aplicar em lotes via PR
(Recommended)".

## Requirements

1. MUST add `internal/cli/history.go` with `roundfix history sanitize` and
   its usage text, and dispatch `history` from `internal/cli/cli.go`. The
   top-level usage gains the line
   `roundfix history sanitize [--batch <n>] [--advise] [--apply] [--promote <path> ...]`
   and the command list gains `history`. `roundfix history` without a
   subcommand, and `--help`, print the usage and exit 0; an unknown
   subcommand exits 2.
2. MUST order units per Invariant 2 and print the plan of Surface
   Transcripts 1 and 4 per Invariants 9-11. A Legacy Archive Folder line
   names its files, bytes, record path and size, disposition and delivery
   (`delivery <12-hex> #<n>`, `delivery <12-hex>` or `delivery unknown`).
   Candidate lines follow the archive plan's grouping. A plan, with or
   without `--advise`, writes nothing under the repository.
3. MUST call `judge.AdviseArchive` for `--advise` exactly as
   `roundfix archive <slug> --plan` does: the same key variables from the
   process environment, the configured monthly ceiling, the injected
   transport and clock, and the Roundfix Home's judge log. The advice is
   printed and never changes which units apply or whether the command
   succeeds.
4. MUST implement `--apply` per Invariants 12 and 13 and API Contract 3,
   with every check before any write, and print the confirmation line of
   Surface Transcript 2. The source revision and the reduction revision are
   `HEAD`. `--promote` paths are repository-relative and must fall inside a
   folder of the batch.
5. MUST refuse an external Spec Root per Invariant 10 and the usage errors
   of Invariant 14, each through Preflight Validation with exit 2.
6. MUST add `docs/user-guide/commands/history.md`. It describes the plan,
   `--batch`, `--advise`, `--apply`, `--promote`, every refusal, the
   History Full Tag and the `no-qa` disposition. It shows the plan's first
   line, `history sanitize plan:`, the confirmation, `history sanitize
   applied`, and the operator batch procedure of `_techspec.md` →
   Operator batch procedure, including `git tag -a history-full`. Add a
   `history` row to `docs/user-guide/commands.md`.
7. MUST add the tests named in Verification to
   `internal/cli/history_test.go`. Each uses a temporary repository built
   with `internal/gittest`, a temporary Roundfix Home, and for advice a fake
   `http.RoundTripper`. None reads this repository's `docs/history` or
   reaches a provider:
   - the plan of a fixture with three Legacy Archive Folders and every kind
     matches Surface Transcript 1 and leaves `git status --porcelain
     --untracked-files=all` and every file byte-identical;
   - `--apply --batch 2` converts exactly the first two folders, the third
     stays, and a second run continues from the third; removed files read
     back from the tag;
   - each tag refusal (missing, lightweight, not an ancestor of `HEAD`,
     lacking a batch path) and a dirty tree exit 2 and change nothing;
   - `--advise` without a key matches Surface Transcript 4 and sends no
     request; with a fake judge it prints each choice; neither writes under
     the repository;
   - a promotion copies the file to `docs/references/` and the record lists
     it; a promotion outside the batch exits 2;
   - every usage error of Invariant 14, and an external Spec Root;
   - with nothing pending, the plan says so and `--apply` exits 0 without
     writing.
8. MUST NOT run the command against this repository, create a tag, or
   change any file under `docs/history` or `.secondbrain-export`.

## Subtasks

- [ ] Add the command, its dispatch and its usage.
- [ ] Add the plan, the advice and the apply path with its refusals.
- [ ] Add the command guide and its index row.
- [ ] Add the command tests.

## Acceptance Criteria

- [ ] A plan never writes in the repository.
- [ ] A batch applies exactly the next units, only on a clean tree covered
      by an annotated ancestor `history-full` tag.
- [ ] The command guide documents the procedure and both emitted patterns.

## Context

- creates: `internal/cli/history.go`
- creates: `internal/cli/history_test.go`
- interface: `internal/cli/cli.go`
- creates: `docs/user-guide/commands/history.md`
- interface: `docs/user-guide/commands.md`
- instruction: `internal/cli/archive.go`
- instruction: `internal/cli/archive_plan_test.go`
- instruction: `internal/judge/archive_advice.go`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- instruction: `docs/agents/cli.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestHistorySanitizePlanWritesNothing|TestHistorySanitizeAppliesTheNextBatch|TestHistorySanitizeRefusesWithoutAnAnnotatedAncestorTag|TestHistorySanitizeRefusesATagWithoutTheBatchPaths|TestHistorySanitizeRefusesADirtyTree|TestHistorySanitizeAdviseFailsOpenWithoutAKey|TestHistorySanitizeAdviseWithAFakeJudge|TestHistorySanitizePromotesIntoReferences|TestHistorySanitizeUsageErrors|TestHistorySanitizeRefusesAnExternalSpecRoot|TestHistorySanitizeWithNothingPending)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistorySanitizePlanWritesNothing TestHistorySanitizeAppliesTheNextBatch TestHistorySanitizeRefusesWithoutAnAnnotatedAncestorTag TestHistorySanitizeRefusesATagWithoutTheBatchPaths TestHistorySanitizeRefusesADirtyTree TestHistorySanitizeAdviseFailsOpenWithoutAKey TestHistorySanitizeAdviseWithAFakeJudge TestHistorySanitizePromotesIntoReferences TestHistorySanitizeUsageErrors TestHistorySanitizeRefusesAnExternalSpecRoot TestHistorySanitizeWithNothingPending; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the eleven tests exists, so the command fails.
- `tmp="$(mktemp -d)" && trap 'rm -rf "$tmp"' EXIT && go build -buildvcs=false -o "$tmp/roundfix" ./cmd/roundfix && help="$("$tmp/roundfix" --help)" && printf '%s\n' "$help" | grep -qF -- 'roundfix history sanitize [--batch <n>]' && "$tmp/roundfix" history --help | grep -qF -- '--apply' && for phrase in 'history sanitize plan:' 'history sanitize applied' 'git tag -a history-full' '--promote' 'no-qa'; do tr -s '[:space:]' ' ' < docs/user-guide/commands/history.md | grep -qF -- "$phrase" || { printf 'guide lacks: %s\n' "$phrase" >&2; exit 1; }; done && grep -qF -- 'commands/history.md' docs/user-guide/commands.md` — expected: exit 0; before this Task the binary has no history command and the guide does not exist, so the command fails.

## References

- `_prd.md` → Goals 1 and 2; User Stories 2 and 3; Core Features 1 and 2; Success Metrics 1, 2 and 3
- `_techspec.md` → Invariants 2 and 9-14; API Contracts 1-4; Surface Transcripts 1-4; Operator batch procedure; Build Order 3
- ADR-0248; ADR-0247; ADR-0215

## Result

Implemented the History Sanitize Command's dispatch, usage, ordered plan,
optional advice and explicit batch apply. The CLI inventories pending units,
checks clean-tree and annotated ancestor tag coverage, and plans all selected
units and promotions before the first write. It reuses the existing legacy
conversion and history-kind APIs with `HEAD` provenance. Advice uses the archive
judge's environment keys, User Config ceiling, injected transport and clock,
and Roundfix Home log; advice failures remain advisory. Plans report delivery,
candidates, measured sizes and file or folder citations. Apply reports one
confirmation and names a unit on write failure.

Added the command guide and command-index row, including both emitted patterns,
all flags and refusals, `no-qa`, the annotated tag and the operator's six-batch
procedure. Existing skill guidance already describes this command from the
preceding Task; no skill or Baseline changes were needed in this slice.

Acceptance evidence from focused implementation checks:

- Plan writes nothing: `TestHistorySanitizePlanWritesNothing` compares the
  complete seven-unit transcript for three folders and every kind, repository
  file bytes and porcelain status. Both advice tests also compare bytes and
  status, exclude core/evidence candidates and use fake HTTP. The additional
  advice checks cover service failure, the configured monthly ceiling and an
  unknown delivery; the fake judge log records its returned versioned model.
- Exactly the next covered units apply: `TestHistorySanitizeAppliesTheNextBatch`
  converts two folders, preserves the third, reads each removed file back
  byte-identically from `history-full`, commits only in the disposable fixture,
  continues at the third and then applies all four kinds. It checks `HEAD`
  provenance, reduction, removal and unchanged retired ADRs. Tag, missing-path,
  dirty-tree, usage, external-root and promotion refusal tests assert exit 2
  and unchanged bytes/status. `TestHistorySanitizePreflightsWholeBatch` proves
  that a malformed second unit prevents the first write, while that unselected
  folder does not prevent a one-unit batch. Additional coverage distinguishes
  a folder named `findings` from the kind and covers folder-only citations.
- Guide describes procedure and output: inspected the guide against the
  TechSpec's Operator batch procedure, checked its output patterns, tag,
  promotion and `no-qa` text, command-index link and Secondbrain link. The
  usage-error test exercises help and subcommand dispatch through the CLI.

Focused commands run (each final invocation exited 0):

- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestHistorySanitizePlanWritesNothing|TestHistorySanitizeAdvise' -count=1`
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestHistorySanitizeApplies|TestHistorySanitizeRefuses|TestHistorySanitizePromotes|TestHistorySanitizeUsage|TestHistorySanitizeWithNothing' -count=1`
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestHistorySanitizePreflights|TestHistorySanitizeFolder|TestHistorySanitizeAdviceIs' -count=1`
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestHistorySanitizeApplies|TestHistorySanitizePromotes|TestHistorySanitizeAdviceWithout' -count=1`
- Python read-only guide assertions: output patterns, operator procedure, index
  and link present.
- `rtk git -c core.fsmonitor=false diff --check`: no whitespace errors.

The task file was already modified at entry. Its status and the Task Graph were
left untouched. No declared Verification command or repository-wide gate was
run; the Daemon owns those checks and settlement. No sanitize command ran
against this repository, no tag was created here, and no `docs/history` or
`.secondbrain-export` file was changed. No commit, push or Pull Request was
made. Every command test uses gittest repositories and temporary Roundfix Home;
no provider was reached.

Consultation limitation: the attempted read of
`https://docs.typesafe.ai/llms.txt` was blocked by the sandbox allowlist. The
implementation uses the existing local `judge.AdviseArchive` contract and does
not change the provider API, questions, model or thresholds.

### Verification feedback repair — attempt 1

Inspected the Daemon diagnostic at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261007T040343Z_bb67f44effbee11e/verification/batch-004-attempt-1.log`.
The affected tests failed while snapshotting Git metadata: the helper read
`.git` before filtering it, and the Daemon environment had an fsmonitor socket
there. This was a snapshot-boundary error, not a sanitize write.

Replaced the snapshot helper with a working-tree walk that skips `.git` before
traversing or opening it, preserving the comparison of all working-tree file
bytes and the separate porcelain-status assertion. Added
`TestHistorySnapshotSkipsGitMetadataBeforeReading`: an explicitly unreadable
metadata entry must not prevent a complete working-tree snapshot. A real-socket
regression fixture was initially attempted, but sandbox socket binding was
denied; the final fixture uses a dangling metadata symlink to exercise the same
read boundary without requiring socket privileges. No production code changed
in this repair.

Fresh focused checks after the repair (exit 0):

- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestHistorySnapshotSkipsGitMetadataBeforeReading|TestHistorySanitizePlanWritesNothing|TestHistorySanitizeAdvise|TestHistorySanitizePromotes' -count=1`
- `GOCACHE=/private/tmp/roundfix-task03-gocache rtk proxy go test ./internal/cli -run 'TestHistorySanitizeRefuses|TestHistorySanitizeUsage|TestHistorySanitizeApplies|TestHistorySanitizeWithNothing|TestHistorySanitizePreflights|TestHistorySanitizeFolder|TestHistorySanitizeAdvice' -count=1`
- `rtk git -c core.fsmonitor=false diff --check`: no whitespace errors.

These checks refresh the plan/no-write, batch/refusal and help evidence above;
the command guide is unchanged. The declared Verification commands were not
rerun. Task status, Task Graph, repository history data and export configuration
remain untouched; no commit, push or Pull Request was created. The Daemon owns
the next full Verification and settlement.

## Carry-forward provenance

- Source Run: `run_20261007T040343Z_bb67f44effbee11e`
- Source commit: `1048ca2f2b2bfdc6ef6bfecd93277bade741860e`
