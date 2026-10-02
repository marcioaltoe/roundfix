---
task: task_04
spec: 0214-measure-before-changing
status: completed
type: docs
complexity: medium
---

# Task 04: Record why corrective Tasks happen and what the Task acceptance question is worth

## Overview

task_01 can classify causes and task_02 can re-score the Task acceptance
question; neither has been run on Roundfix's own records. This Task runs both
against the maintainer's Run Database and the judge's transports, saves the
two records, writes two reference documents whose verdicts follow
`_techspec.md` → Decision rules, and writes a Backlog Entry for each rule that
fires. It changes no code, gate, judgment or signature. It is verifiable on
its own: the two record consistency tests recompute every figure and verdict
from the committed records and documents.

## Requirements

1. MUST build the binary with `make build` and run `bin/roundfix runs causes --since 2026-09-17 --until <the UTC date this Task runs> --format json` from this Task's worktree with `NODE_OPTIONS` unset, saving its standard output unchanged as `docs/references/corrective-causes.json`. The command opens the Run Database read-only; this Task MUST NOT open it any other way, and MUST NOT run `roundfix gc` or anything else that writes Roundfix Home.
2. MUST NOT change `internal/runcause/signatures.json` or any other code, test or fixture. When more than a third of the items are `unclassified`, the document lists the ten most frequent unclassified checks instead.
3. MUST write `docs/references/corrective-causes-measurement.md` with the window, the table digest, the count per class and per trigger, the repository-knowledge share, the ten most frequent signatures, a sample of five items per class with their signature, the rule of `_techspec.md` → Decision rules, rule 1, applied to the counts, and exactly one line `Verdict: reopen`, `Verdict: keep closed` or `Verdict: inconclusive`. It MUST compare the result with the published evidence of `_prd.md` → Acceptance evidence without letting that evidence change the verdict.
4. MUST run `go test -count=1 -timeout 60m -run '^TestMeasureTaskAcceptance$' ./internal/judge -args -measure-task-acceptance -measure-labels=<absolute path of docs/references/corrective-causes.json> -measure-out=<absolute path of docs/references/task-acceptance-remeasurement.json>` with the keys the environment holds, and MUST NOT set, print or write a key. Requests carry only this repository's Task files, through OpenRouter first and TypeSafe second, under the Judge Log's US$5 monthly ceiling. When neither `ROUNDFIX_OPENROUTER_API_KEY` nor `ROUNDFIX_TYPESAFE_API_KEY` is set, the harness writes a `blocked` record, and the document says so.
5. MUST write `docs/references/task-acceptance-measurement.md` with the transport and the model each answer reported, the calls, input tokens and cost, the counts answered, excluded and repaired, AUROC with its interval, the random and length baselines, precision and recall at the threshold, a comparison with the 2026-09-30 figures (AUROC 0.62 [0.55, 0.70]) that names the question as new, the rule of `_techspec.md` → Decision rules, rule 2, applied to the figures, and exactly one line `Verdict: adopt`, `Verdict: do not adopt` or `Verdict: inconclusive`.
6. MUST write `docs/backlog/<date>-reopen-the-memory-question-for-agent-sessions.md` (`feat`, `open`) only when the cause verdict is `reopen`, and `docs/backlog/<date>-adopt-the-task-acceptance-judgment.md` (`feat`, `open`) only when the acceptance verdict is `adopt`, each citing its reference document. Their file names carry the date the Task runs, so they are not declared in Context and the Daemon records them under `## Recorded paths`.

## Subtasks

- [ ] Build the binary and save the cause record over the window.
- [ ] Write the cause document and apply rule 1.
- [ ] Run the harness with its flag and save the Task acceptance record.
- [ ] Write the Task acceptance document and apply rule 2.
- [ ] Write a Backlog Entry for each rule that fired.

## Acceptance Criteria

- [ ] `TestCausesRecordIsConsistent` passes on the committed cause record and document.
- [ ] `TestTaskAcceptanceRecordIsConsistent` passes on the committed Task acceptance record, its document and the cause record as labels.
- [ ] A reopen or adopt Backlog Entry exists exactly when its verdict says so.
- [ ] No code, test, fixture or signature file changed.

## Result

The read-only causes measurement was built with `env -u NODE_OPTIONS make build`
and run for `2026-09-17..2026-10-02`, producing 177 Runs and 132 items. The
record reports 23 repository-knowledge items and 100 unclassified items; the
cause document applies Decision rule 1 and records `Verdict: inconclusive`.

The acceptance harness ran with the required live flag and saved a measured
record: OpenRouter, `jev-1.13`, 237 calls, 237 answered, 31 repaired, 7
excluded, AUROC 0.466881 with interval [0.358034, 0.574594], and no discarded
bootstrap draws. The acceptance document applies Decision rule 2 and records
`Verdict: do not adopt`.

No rule fired a Backlog Entry. The only changed implementation-scope paths are
the two JSON records, the two measurement documents, and this Result section;
no code, test, fixture, or signature file was changed. The Daemon still owns
the declared Verification commands and Task status.

Focused checks passed: `go test -count=1 -run
'^TestCausesRecordDecisionBoundaries$' ./internal/runcause`,
`go test -count=1 -run '^TestTaskAcceptanceVerdictRules$' ./internal/judge`,
and independent `jq` arithmetic checks over both saved records. The two
declared record-consistency tests were not run because they are Daemon-owned
Verification commands; their fresh outcomes remain pending for settlement.

## Context

- creates: `docs/references/corrective-causes.json`
- creates: `docs/references/corrective-causes-measurement.md`
- creates: `docs/references/task-acceptance-remeasurement.json`
- creates: `docs/references/task-acceptance-measurement.md`

## Verification

- `root="$(git rev-parse --show-toplevel)" || exit 1; out="$(go test -count=1 -v -run '^TestCausesRecordIsConsistent$' ./internal/runcause -args -causes-record="$root/docs/references/corrective-causes.json" -causes-document="$root/docs/references/corrective-causes-measurement.md" 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestCausesRecordIsConsistent" || { printf 'missing pass: TestCausesRecordIsConsistent\n' >&2; exit 1; }` — expected: exit 0; before this Task the cause record does not exist, so the test fails.
- `root="$(git rev-parse --show-toplevel)" || exit 1; out="$(go test -count=1 -v -run '^TestTaskAcceptanceRecordIsConsistent$' ./internal/judge -args -task-acceptance-record="$root/docs/references/task-acceptance-remeasurement.json" -task-acceptance-document="$root/docs/references/task-acceptance-measurement.md" -task-acceptance-labels="$root/docs/references/corrective-causes.json" 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestTaskAcceptanceRecordIsConsistent" || { printf 'missing pass: TestTaskAcceptanceRecordIsConsistent\n' >&2; exit 1; }` — expected: exit 0; before this Task the Task acceptance record does not exist, so the test fails.
- `for spec in "docs/references/corrective-causes-measurement.md|reopen|reopen-the-memory-question-for-agent-sessions" "docs/references/task-acceptance-measurement.md|adopt|adopt-the-task-acceptance-judgment"; do file="${spec%%|*}"; rest="${spec#*|}"; fires="${rest%%|*}"; name="${rest#*|}"; test -f "$file" || { printf 'missing %s\n' "$file" >&2; exit 1; }; verdict="$(sed -n 's/^Verdict: //p' "$file")"; test -n "$verdict" || { printf 'no Verdict line in %s\n' "$file" >&2; exit 1; }; count=0; for entry in docs/backlog/*-"$name".md; do test -f "$entry" && count=$((count + 1)); done; if test "$verdict" = "$fires"; then test "$count" -eq 1 || { printf 'verdict %s needs one Backlog Entry %s\n' "$verdict" "$name" >&2; exit 1; }; else test "$count" -eq 0 || { printf 'verdict %s needs no Backlog Entry %s\n' "$verdict" "$name" >&2; exit 1; }; fi; done` — expected: exit 0; before this Task neither document exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 2, 3 and 5; User Stories 1 and 3; Core Features 8 and 11; Success Metrics 2, 3 and 5; Acceptance evidence; Recorded limits
- [_techspec.md](_techspec.md) — Decision rules; The harness; Data Models; API Contract 4; Measured outside evidence; Testing Approach 5; Build Order 4
- [references/2026-09-30-measure-why-corrective-tasks-happen.md](references/2026-09-30-measure-why-corrective-tasks-happen.md)
- [references/2026-09-30-re-measure-the-task-lint-judgment-on-a-cleaner-label.md](references/2026-09-30-re-measure-the-task-lint-judgment-on-a-cleaner-label.md)
- ADR-0214; ADR-0215; ADR-0200; ADR-0201

## Carry-forward provenance

- Source Run: `run_20261002T171516Z_94169d908e48d811`
- Source commit: `41942907d4ca9640549a8831ebfd693f44e3fdd6`
