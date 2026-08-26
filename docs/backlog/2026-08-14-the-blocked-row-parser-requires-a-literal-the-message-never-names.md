---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The blocked-row parser requires a literal the message never names

## Symptom

A QA gate cycle executed no journey. It stopped at the mechanical detector, which
evaluated the *previous* report and refused it on shape, producing a new report
whose entire content is that refusal:

```text
### QA-REPORT-SHAPE
- location: .../qa-report-2026-08-12-01.md:79
- detail: row NG-01 has a blocked cause outside environment, finding, or declared
- fix: Use blocked (environment: ...), blocked (finding: ...), or blocked (declared: ...).
```

The refused row was:

```text
| NG-01 | ... | blocked (finding: F-06 — the declared runtime canary aborts on its
stale function-count assertion) | ... |
```

It **is** typed as `finding`. What it lacks is the literal `waits on`: the
`qa-gate` contract declares the form `blocked (finding: <id> — waits on <named
failure>)`, and the parser does not recognise the row without it. Because neither
row matched, the detector also reported `rows_blocked_finding is 2 but the
Results table contains 0 matching rows` — two symptoms of one cause.

## Where

The QA report's blocked-row parser and the `QA-REPORT-SHAPE` diagnostic it emits.

## Expected

The message names the literal it requires. A `fix:` line that repeats the three
type prefixes, when the row already carries the right prefix, sends a reader
looking for the wrong thing — and the second symptom sends them looking for a
counting bug that does not exist.

Worth settling in the same work: whether a row typed correctly but phrased
differently should refuse at all, or whether the parser should accept the type
and treat the prose after it as free text.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-12-gate-gasta-um-ciclo-validando-a-forma-do-relatorio-anterior.md`
in the Secondbrain. Related:
`docs/backlog/2026-08-13-a-refused-gate-writes-a-report-its-own-contract-rejects.md`
records the adjacent defect that a refused gate leaves a report blocking its own
successor; this entry names why one such refusal was unactionable.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
