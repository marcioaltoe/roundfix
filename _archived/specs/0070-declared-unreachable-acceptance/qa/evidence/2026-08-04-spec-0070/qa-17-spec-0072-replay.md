# QA-17 — Spec 0072 scoping replay

Status: pass

Historical input:
`docs/specs/_archived/0072-qa-is-a-task-not-a-flag/qa/qa-report-2026-08-03-03.md`.

- A table-only `awk` count found exactly 15 rows, 7 through 21, recorded
  `blocked (finding: F-001)`.
- Row 5 is the failed tooling-authority chronology audit. Its entry point and
  evidence are authorization records, Git ancestry, skill mirrors, and derived
  pins. Row 6's full `make verify` passed. Rows 7–21 exercise Task Graphs,
  Implement, CLI flags, reports/events, history, legacy loading, docs, help,
  and Non-Goals; their entry points and observables do not depend on repairing
  the row-5 commit chronology.
- Applying the current QA contract therefore leaves row 5 failed and makes all
  fifteen formerly blocked rows runnable. It does not mark them passed or
  excuse F-001.

The current gate then exercised the rule rather than stopping at artifact
analysis. This static contract check exited 1:

```text
rtk rg -q 'partial verdict whose blocked rows' \
  docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md
```

Its missing text is the QA-16 documentation finding. The check implicates only
QA-16: it does not build, execute, persist, or confirm any archive behavior.
QA-01 through QA-15 remained valid and reported their own terminal results;
the positive control `roundfix archive --help` exited 0 and exposed the
missing contract. This is the live scoping behavior Core Feature 7 requires.
