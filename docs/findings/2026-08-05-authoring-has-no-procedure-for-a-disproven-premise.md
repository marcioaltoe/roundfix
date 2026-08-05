# 2026-08-05 — Spec authoring has no procedure for a disproven premise

status: pending

## What was observed

oraculum Spec 0020 needed **five gate cycles**. Three of them failed the same
finding — `RF-005`, "active Spec artifacts retain the disproven source model" —
and every one of the three was a partial fix by the Supervisor, not a defect the
implementation introduced.

The Spec's premise was disproven mid-flight. Its Story 2 promised to measure
paid late-payment interest as the difference between face value and paid value.
The maintainer confirmed the source cannot express that: `fn2_valor` holds the
amount due before settlement and the amount paid after, and never stores the
original. The difference is unobservable, and the generated SQL aliased both
ends from the same expression, so production interest could only ever be zero.

The corrected source — interest and penalty columns the ERP already writes —
had to replace the old model **everywhere**. The concept was spread across:

| Artifact | Places it lived |
| --- | --- |
| `_prd.md` | Goal, Story, Functional Requirement, intro paragraph |
| `_techspec.md` | architecture, interfaces, testing approach, decisions, risks, build order |
| `task_01.md` | overview, requirement, subtask, two acceptance criteria, `## Result` |
| `task_02.md` | overview, acceptance criterion, `## Result` |
| published | CHANGELOG, user guide, release catalog |
| code | SQL, value object, response contract |

The Supervisor corrected it three times by grepping for the terms it remembered
— `valor_original`, `face×pago`, `classificarDiferenca` — and each pass left
behind whatever used a term outside that set: a bare `diferenca` aggregate, an
`acrescimo/desconto` line in Testing Approach, a Decisions entry about summing
discounts in the same field with a sign. Each miss cost a ~35-minute gate cycle.

## Root cause

Nothing in the authoring skills covers **supersession**. `write-prd`,
`write-techspec`, and `write-tasks` all assume forward authoring: a premise is
established, then decomposed. There is no step for "this premise turned out to
be false, propagate the correction," and therefore no checklist of where a
concept can hide inside a Spec folder.

Two properties of the shape make ad-hoc grepping fail reliably:

1. **The concept has more surface than its name.** A model called "face versus
   paid" appears as `valor_original`, `valor_pago`, `diferenca`,
   `classificarDiferenca`, `acrescimo`, `desconto`, and as prose. A sweep built
   from remembered identifiers is guaranteed to be a subset.

2. **`## Result` sections count.** The QA gate audits completed-Task evidence as
   an active artifact, so a Result describing an implementation that was later
   removed reads as a live contradiction. That is not obvious from the skills:
   `implement-task` presents `## Result` as the Agent's handback record, and
   nothing says a later correction must reconcile it. The Supervisor spent a
   cycle unsure whether rewriting it would falsify history — the resolution
   that worked was a supersession note above the original text, which preserves
   the record and removes the contradiction, but that pattern is documented
   nowhere.

This is distinct from
[a static gate row reported one instance per cycle](2026-08-04-a-static-gate-row-reported-one-instance-per-cycle.md),
which is about the gate naming one instance at a time. Here the gate's cycle-4
report was comprehensive and named all three files; the cycles were still spent
because the *authoring* side had no procedure to act on it exhaustively.

## What would settle it

- Give the authoring skills a **supersession procedure**: when a premise is
  disproven, enumerate the concept's full vocabulary first — identifiers, field
  names, function names, and the prose phrase — then sweep every artifact in the
  Spec folder against that vocabulary in one pass, before any Run.
- State plainly, in `write-tasks` or `implement-task`, that **`## Result` is
  audited as active evidence**. Prescribe the supersession-note pattern: keep
  the original text, prepend a note naming what superseded it and why. Rewriting
  falsifies the handback record; leaving it untouched fails the gate.
- Consider a mechanical aid: a check that takes a retired term list and reports
  every occurrence across `docs/specs/<slug>/`, so the sweep is verifiable
  rather than remembered. The Spec Consistency Check of
  `0064-spec-artifact-consistency-gate` is the natural home — it already reads
  the whole folder, and a retired-vocabulary rule is the same traversal.
- A Spec whose premise is disproven mid-flight may deserve an explicit
  authoring state. Five cycles of partial correction is worse than one deliberate
  re-authoring pass, and nothing today prompts the Supervisor to choose the
  second.

## Evidence

- oraculum Spec 0020, QA reports `qa-report-2026-08-05.md` through
  `qa-report-2026-08-05-03.md`; `RF-005` repeats in the last three.
- The disproven premise and its source proof are recorded in
  `docs/handoffs/2026-08-05-spec-0020-congelada.md` and in the corrected
  `_prd.md` Story 2.
