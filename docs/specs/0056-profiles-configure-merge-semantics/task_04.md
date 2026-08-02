---
task: task_04
spec: 0056-profiles-configure-merge-semantics
status: completed
type: backend
complexity: medium
---

# Task 04: Summarize the effective change before writing

## Overview

The destructive write was mistaken for a formatting change because the
confirmation showed a 103-line diff and named only the requested category. This
Task renders the effective change per category — added, replaced, removed —
before any write, in both the interactive and machine flows, so a removal can
never hide inside reformatting churn. It also scopes proof to the categories
the operation writes.

## Requirements

1. MUST render, before any write, one line per affected category stating
   whether it is added, replaced, or removed, and nothing about untouched
   categories.
2. MUST render an explicit removal line whenever a removal was declared,
   including for a category the file did not contain.
3. MUST render the same effective change in the machine output, so an Agent
   sees what a human sees.
4. MUST prove every distinct Agent Selection tuple the operation writes before
   the confirmation, keeping proof ahead of every write.
5. MUST NOT prove tuples belonging to categories the operation does not write,
   so a stale pre-existing entry cannot block an unrelated edit — the
   deliberate deviation recorded in the TechSpec's Decisions.
6. MUST keep `--dry-run` writing nothing and reporting the same summary.
7. MUST NOT change exit codes in this Task.

## Subtasks

- [ ] Render the per-category summary for the interactive flow.
- [ ] Render the same change set in the machine output.
- [ ] Scope proof to the written categories.
- [ ] Keep `--dry-run` reporting the summary without writing.

## Acceptance Criteria

- [ ] Configuring one category prints exactly one summary line, naming it as
      added or replaced, with no line for untouched categories.
- [ ] A declared removal prints its own removal line, including when the file
      did not contain that category.
- [ ] The machine output carries the same per-category classification as the
      text output for the same operation.
- [ ] Every tuple in a written category is proven before the confirmation is
      offered.
- [ ] A configured category the operation does not write, holding a tuple that
      cannot be proven, does not block the operation.
- [ ] `--dry-run` prints the summary, writes nothing, and leaves the file
      byte-identical.
- [ ] `git status --porcelain` shows no path outside `internal/config/`,
      `internal/cli/`, and this task file.

## Context

- interface: `internal/cli/profiles_configure.go`
- interface: `internal/config/profile_config.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -run TestProfilesConfigureChangeSummary -count=1` —
  expected: exit 0; text and machine summaries agree and cover add, replace, and
  remove.
- `go test ./internal/cli -run TestProfilesConfigureProofScope -count=1` —
  expected: exit 0; written tuples are proven and an unproven untouched category
  does not block.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0.
- `go test ./internal/config ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Stories 2 and 5; Core Features 3 and 6; User Experience.
- `_techspec.md` → Coverage Map; Build Order 4; Decisions (proof-scope
  deviation).
- ADR-0037, ADR-0039.

## Result

Implemented the Task 04 slice without changing the Daemon-owned status or exit
codes. The prepared proposal now retains its ordered Effective Change Set, and
`profiles configure` uses that one value to render the text summary, populate
the additive JSON `changes` array, select only added/replaced categories for
exact tuple proof, and pass declared removals to the category-level merge. The
CLI accepts repeatable `--remove` flags, including a removal-only invocation;
an absent category remains a reported removal even when persistence is a
no-op.

Acceptance evidence:

- `TestProfilesConfigureChangeSummary/added_category_agrees_in_text_and_machine_output`
  observes exactly `added: backend` and the matching JSON
  `{category:"backend", kind:"added"}` entry.
- `TestProfilesConfigureChangeSummary/replacement_excludes_untouched_categories`
  observes exactly `replaced: backend` and proves the preview and JSON omit the
  untouched `frontend` category.
- `TestProfilesConfigureChangeSummary/declared_removals_include_present_and_absent_categories`
  observes `removed: frontend` and `removed: review` in both text and JSON,
  including the absent `review` category.
- `TestProfilesConfigureProofScope` observes all three distinct tuples from the
  written `backend` and `qa` profiles before confirmation, rejects any attempt
  to prove the stale untouched `frontend` tuple, and confirms the config bytes
  are unchanged when confirmation begins.
- `TestProfilesConfigureChangeSummary/dry_run_reports_the_same_change_and_preserves_bytes`
  observes the same added classification in text and JSON dry-run modes and
  verifies the target file remains byte-identical after each invocation.
- Changed-path inspection after implementation shows only `internal/config/`,
  `internal/cli/`, and this task file.

Focused checks run with `GOCACHE=/private/tmp/roundfix-task04-gocache`:

- Pre-change signal:
  `go test ./internal/cli -run '^TestProfilesConfigureChangeSummary$/^added_category_agrees_in_text_and_machine_output$' -count=1`
  failed to compile because the machine response had no `Changes` field.
- Each `TestProfilesConfigureChangeSummary` subtest passed independently for
  add, replace, removal, and dry-run behavior.
- `go test ./internal/cli -run '^TestProfilesConfigureProofScope$/' -count=1`
  passed.
- `go test ./internal/cli -run '^TestProfilesConfigure' -count=1` passed.
- The focused pre-existing configure regression group covering file, `--yes`,
  dry-run, proof order/failures, interactive input, fallback validation, and
  JSON rollback passed.
- `go test ./internal/config -run '^(TestEffectiveChangeSet|TestWriteProfilesConfig.*)$' -count=1`
  passed.
- `go test ./internal/cli -run '^TestProfilesDocumentationContractMatchesPublicGuidance$' -count=1`
  passed.

The commands in `## Verification` were not run; the Daemon owns that gate.
