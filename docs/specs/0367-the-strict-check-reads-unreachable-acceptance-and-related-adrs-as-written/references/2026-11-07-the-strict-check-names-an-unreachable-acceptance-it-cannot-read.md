---
type: fix
status: promoted
created: 2026-11-07
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
reason: null
---

# The strict check names an Unreachable Acceptance entry it cannot read

## Symptom

A PRD whose `## Unreachable Acceptance` section holds an entry the declaration parser does not read passes `roundfix spec check <slug> --strict` with no findings. The mismatch surfaces only at the end of the Run, when the QA gate's `roundfix qa-report accept` refuses the report with `rows_blocked_declared is 1, but Spec declares 0 unreachable acceptances; shortfall is 1`, and the item parks `run-unresolved` after every implementation Task has already passed. Spec 0364 lost its first QA Run this way on 2026-11-06: its PRD wrote the one entry as a numbered prose paragraph (`1. A real rerun of a failed release job ...`) instead of the fielded form.

## Where

`internal/spec/spec.go` `parseUnreachableDeclarations` starts an entry only at a top-level `- ` line and reads the `criterion`, `reason` and `satisfied-by` fields; any other non-empty text in the section is ignored, so a numbered or prose entry yields zero declarations without an error. `internal/speccheck` does not call `spec.Unreachable`, so the strict check never compares the section with what the parser reads. Callers today: `internal/spec/archive.go`, `internal/spec/qa.go`, `internal/cli/baseline_update.go`.

## Expected

The strict check reports a finding, naming the PRD line, when the `## Unreachable Acceptance` section holds non-empty content that does not parse as a fielded declaration (for example a numbered item, a prose paragraph, or a `- ` line without `criterion:`), and the finding shows the fielded form to use. A section that parses fully, or holds no entries, keeps passing. Whether the parser itself should refuse unread content instead of ignoring it is a design choice for the Spec.

## Evidence

The parked QA Report of Run `run_20261107T002718Z_f5b7975ba1530c9f` (finding F-01, rows 13 and 15); the amendment that unblocked it, which rewrote the entry in the fielded form, is in the history of PR #783; a probe on 2026-11-06 read 0 declarations from the prose PRD and 1, matching QA row 13, from the fielded one, while `spec check --strict` reported no findings on both. Requested by the maintainer on 2026-11-07, relayed through jarvis.
