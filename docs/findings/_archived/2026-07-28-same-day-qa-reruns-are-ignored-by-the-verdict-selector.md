---
status: done
created_at: 2026-07-28
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# QA verdict — every same-day rerun was ignored, so a passing Spec still reported fail (2026-07-28)

`NewestQAReport` picked the day's **first** QA Report, not its latest. A Spec
that failed its gate and then passed on a rerun the same day kept reporting the
stale failure, and `roundfix archive` refused it.

The selector globbed `qa/qa-report-*.md`, sorted the paths as raw strings, and
took the last, documenting the assumption that "Report names embed the date as
YYYY-MM-DD, so the lexicographically greatest path is the newest report". Same-day
reruns get a sequence suffix, and that assumption inverts for them:

```
qa-report-2026-07-28.md        first run of the day
qa-report-2026-07-28-02.md     second run of the day — newer
```

After the shared prefix, the suffixed name has `-` (0x2D) and the plain one has
`.` (0x2E). `-` sorts first, so `reports[len(reports)-1]` returns the *older*
file. Every same-day rerun was invisible.

## Evidence

Spec `0042-verification-capacity-and-daemon-task-settlement`, 2026-07-28. The
gate failed on build `8593002` with two findings, both were fixed, and the rerun
on build `ffd6852` recorded `verdict: pass` with "No current-build product
finding". Roundfix reported:

```text
qa fail — docs/specs/0042-…/qa/qa-report-2026-07-28.md
```

Both files were present: `qa-report-2026-07-28.md` (`verdict: fail`, build
`8593002`) and `qa-report-2026-07-28-02.md` (`verdict: pass`, build `ffd6852`).
The Daemon read the first.

The defect is not cosmetic. `QAVerdict`, `spec.ArchiveSpec`, and the Daemon's
`settleQAVerdict` all funnel through this one selector, so a same-day rerun could
not end a Run Clean and could not archive its Spec — the loop had no way forward
except waiting for the next calendar day.

## Fixed

`internal/spec/qa.go` now orders by real recency instead of raw bytes: parsed
date, then sequenced-before-unsequenced within a date, then numeric sequence
(unsuffixed = 1, so `-10` beats `-02` beats unsuffixed), then path as a
deterministic tiebreak. A name with a valid date and a non-numeric suffix keeps
its date and loses only to sequenced reports of that same date; a name with no
valid date ranks below every dated report. Regression tests cover the exact
failing case and each ordering rule; four of them fail against the previous
implementation.

## Related defect, not fixed here

The Daemon's Agent prompt and the qa-gate Skill disagree about report naming.
`internal/agent/spec_prompt.go:53` instructs the Agent to write
`qa-report-YYYY-MM-DD.md` with no mention of a suffix, while
`.agents/skills/qa-gate/SKILL.md:59` specifies
`qa-report-YYYY-MM-DD-<scope-or-build>.md` or `-NN` for collision safety. So the
Daemon-driven QA step can be told to overwrite the day's report while the
skill-driven one creates a suffixed sibling — the same input producing two
different artifacts depending on which surface ran it. Reconciling that requires
touching a protected Skill or the prompt plus its exact-string test, so it is
left for an authorized change.

## Suggested acceptance checks

- A Spec that fails its gate and passes on a same-day rerun reports `pass` and
  archives.
- `-10` is treated as newer than `-02`; a later date beats any sequence of an
  earlier one.
- The Daemon prompt and the qa-gate Skill agree on the report filename.

## What worked — keep

- Writing a suffixed sibling instead of overwriting the day's report is the right
  behavior; it preserved both verdicts and made the defect diagnosable at all.
  The bug was entirely in how they were ordered afterwards.

## Addendum — 2026-07-28 — Fixed directly on `main`

The recency-ordering fix and its regression tests merged to `main` through
Pull Request #40 (squash commit `ed4abec`). Spec
`0042-verification-capacity-and-daemon-task-settlement` archived through the
corrected selector the same day, satisfying the first acceptance check. The
selector needed no implementation Spec; this status records the direct fix.
The related report-naming disagreement between the Daemon prompt and the
qa-gate Skill remains open and travels with the QA-gate reachability work in
[2026-07-28 — QA gate cannot reach Pull Request journeys](2026-07-28-qa-gate-cannot-reach-pull-request-journeys.md).
