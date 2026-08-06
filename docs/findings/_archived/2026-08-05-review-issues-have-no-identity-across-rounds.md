---
status: done
created_at: 2026-08-05
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
---

# 2026-08-05 — Review Issues have no identity across Rounds

status: pending

## What was observed

One pull request — `gesttione-solutions/vortex#123` — accumulated **114 Review
Issue files for 70 distinct review threads**. 44 files, 39% of the total, were
redundant.

The inflation came from calling `fetch` three times over the life of the pull
request. Each call created a new Round and re-imported every still-open thread
as a **new** issue file, while the previous Round's copies stayed `pending`.

| Round | Files | Statuses |
| --- | --- | --- |
| 001 | 18 | 14 resolved, 4 failed |
| 002 | 36 | 9 resolved, 17 invalid, 9 pending |
| 003 | 60 | 60 pending |

At the time of measurement GitHub reported **58 open threads**, while the
repository held **69 pending issue files** across those three Rounds.

The `resolve` Run that followed planned **five Batches of 20** for what was, in
substance, 58 findings.

### The same thread reached contradictory verdicts

Grouping every issue file by the thread id inside `source_ref` shows 18 threads
with more than one file, and some of the duplicates disagree:

```
PRRT_kwDOQLFe_s6WxoVQ → round-001:resolved, round-002:invalid, round-002:resolved
PRRT_kwDOQLFe_s6WxoVt → round-001:resolved, round-002:invalid, round-002:resolved
PRRT_kwDOQLFe_s6WxoVy → round-001:resolved, round-002:invalid, round-002:resolved
PRRT_kwDOQLFe_s6WxoV8 → round-001:resolved, round-002:resolved, round-002:resolved
```

The first three show one thread judged `invalid` in one file and `resolved` in
another **within the same Round**. That is worse than the wasted turn: the
artifact record no longer says what was decided about that finding.

## Root cause

The identity data is already modelled. Every issue file carries:

```yaml
source_ref: thread:PRRT_kwDOQLFe_s6WxoWG,comment:PRRC_kwDOQLFe_s7d8nLX
review_hash: 833a4d1c8180b2bcfc8d40768d9979607daf72656c6658089c6db434f5f34f22
duplicate_of: ""
```

`duplicate_of` exists for exactly this purpose and was empty in all 114 files.
The documented contract already reserves `duplicated` as Daemon-owned
bookkeeping that Agents must not set — so the mechanism is specified, owned,
and unused across Rounds.

Two things appear to be missing:

1. **`fetch` does not consult prior Rounds of the same pull request.** It writes
   a fresh Round from the current Review Source state without asking whether a
   thread already has an issue file.
2. **The effective key includes the comment, not just the thread.** A thread
   that accumulates comments yields several issue files, which is how the same
   thread got two files inside one Round.

## What would settle it

Key a Review Issue by its **thread**, and let `fetch` reconcile against every
Round of the same pull request:

- thread already has a file in a **terminal** status (`resolved`, `invalid`,
  `failed`) → create the new file as `duplicated` with `duplicate_of` pointing
  at it, or do not create it at all;
- thread already has a file that is **`pending`** → do not create a second file;
  update the existing one in place if the head moved;
- thread has no file → create it, as today.

`review_hash` stays useful as the content check that decides whether a
re-reported thread is the *same* finding or a genuinely new comment on an old
thread.

A cheaper partial fix, if the reconciliation is bigger than it looks: refuse
`fetch` while any Round of that pull request still holds `pending` issues, and
say so. `resolve` already picks up pending issues from every Round, so the
second `fetch` had nothing to add.

## Why it costs more than the wasted turns

The Batch count is the visible cost — five Batches where three would do. The
expensive one is quieter: an Agent that re-judges a settled finding can settle
it differently, and nothing in the artifact record marks which verdict is
current. A reviewer reading `docs/specs/_reviews/pr-123/` afterwards cannot tell
whether `PRRT_…6WxoVQ` was accepted or rejected.

## Related

- `docs/findings/2026-08-05-o-loop-devolve-controle-por-motivos-mecanicos.md`
- The measurement came from a live autonomous queue night on the `vortex`
  repository, driving Specs 0015, 0016 and 0032 through Roundfix.

## Addendum, 2026-08-06 — the same defect also hides open findings

The original measurement showed duplication: 114 files for 70 threads. Two more Specs the same
night showed the mirror failure, which is worse.

On Specs 0016, 0019 and 0017, after one `resolve` cycle settled its issues, every later `fetch`
returned `Review Issues: none` **while GitHub still reported open threads** — three, six and
twelve respectively. Roundfix considers a thread already imported and terminal, so it never
re-imports it; the Agent is never handed a finding that is demonstrably still open.

The consequence is that the normal loop cannot close those findings at all. Each of the three
Specs needed a hand-authored corrective Task naming the findings, which also invalidates the
authored gate and forces another full gate cycle.

So the missing identity costs twice, in opposite directions: without it, the same finding is
re-imported as a new issue across Rounds, and once settled, a still-open thread can never be
imported again. Keying on the thread fixes both — a thread that is open at the Review Source and
terminal in the artifacts is precisely the state that should produce a fresh issue.

## Spec pointer

None yet.
