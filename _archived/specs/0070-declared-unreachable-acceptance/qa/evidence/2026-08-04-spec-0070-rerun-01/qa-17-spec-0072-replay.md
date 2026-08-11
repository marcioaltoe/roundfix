# QA-17 — Spec 0072 scoping replay

Status: pass

Historical input:
`docs/specs/_archived/0072-qa-is-a-task-not-a-flag/qa/qa-report-2026-08-03-03.md`.

- A table-only `rtk awk` count found exactly 15 rows, 7 through 21, recorded
  `blocked (finding: F-001)`.
- Row 5 is the failed protected-tooling chronology audit. Row 6's full gate
  passed. Rows 7 through 21 exercise Task Graphs, Implement, CLI flags,
  reports and events, history, legacy loading, docs, help, and Non-Goals; none
  depends on repairing row 5's commit chronology.
- Applying the current gate contract leaves row 5 failed and makes all 15
  formerly blocked rows runnable. It does not mark them passed or excuse
  F-001.

The current gate then demonstrated the same scoping rule: its transient
integration-test wait affected QA-02 while unclassified, but every independent
archive, governance, and documentation row continued and reached its own
terminal result. The unchanged-build reproduction cleared QA-02 without
changing any other result.
