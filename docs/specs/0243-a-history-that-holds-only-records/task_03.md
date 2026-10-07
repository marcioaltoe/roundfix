---
task: task_03
spec: 0243-a-history-that-holds-only-records
status: pending
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
