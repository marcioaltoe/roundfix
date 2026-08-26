---
type: fix # feat | fix | perf | refactor
status: done
created: 2026-08-13
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A refused gate writes a report its own contract rejects

## Symptom

When the QA gate refuses at its authoring precondition, it writes a report whose
Results table is empty — correctly, because it stopped before building the
matrix. The mechanical stage of every later run then reads that report and
refuses:

```text
QA-REPORT-SHAPE
location: qa/qa-report-2026-08-12.md:1
detail:   Results table has no report rows
fix:      Materialize every planned QA row with one terminal status.
```

So a gate that refused for a good reason leaves behind an artifact that blocks
every future gate run on the same Spec, and the prescribed fix — materialize
every planned row — is impossible for a run that never built a matrix.

Measured on Spec 0094 on 2026-08-12/13. The 2026-08-12 run stopped at a strict
Spec check failure, which is the behaviour ADR-0096 asks for: prove machine facts
before spending an Agent turn. It wrote a structurally invalid report while doing
the right thing. The 2026-08-13 run then failed with one
`rows_blocked_finding`, pointing at the previous day's file rather than at
anything in the work under test.

The only exit was deleting the empty report, which is evidence removal — the one
move the repository's rules single out as forbidden when it makes a failure
disappear. It took a maintainer decision to do it.

## Where

The QA gate's refusal path, where it writes a report after stopping at a
precondition, and the mechanical stage that reads every report in the Spec's
evidence directory rather than the one whose verdict is current.

## Expected

A gate that refuses before building its matrix either writes no report, or writes
one whose shape its own contract accepts — for instance a single terminal row
recording the precondition refusal, which is what actually happened and is
materially more useful than an empty table.

Failing that, the mechanical stage reads the newest report rather than every
report, so a superseded refusal cannot block the run that supersedes it.

Worth settling in the same work: whether a QA report that records a precondition
refusal is evidence worth keeping at all, or whether the refusal belongs in the
Run's own record and not in the Spec's evidence directory.

## Evidence

`docs/specs/0094-one-history-root-under-docs/qa/qa-report-2026-08-12.md` before
its removal, and `qa-report-2026-08-13.md`, whose sole finding is about that
file. Related: `docs/history/findings/2026-07-28-same-day-qa-reruns-are-ignored-by-the-verdict-selector.md`
records an adjacent defect in how reports are selected.

---

Triage 2026-08-26: delivered by Spec 0113. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
