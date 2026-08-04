# 2026-08-04 — A static gate row reported one instance per cycle

status: pending

## What was observed

Closing fluxus Spec 0012 took five Runs. Three consecutive QA gate executions
failed the **same** row — `PC-03`, the Project Constraints audit — each naming
exactly **one** missing ADR, so each fix bought exactly one more cycle.

| Gate | Verdict | Failing row | Named |
| --- | --- | --- | --- |
| #2 | fail | F-001, F-002 | ERP grain contradiction; a test mock defect |
| #3 | fail | PC-03 | ADR-0050 absent from Project Constraints |
| #4 | fail | PC-03 | ADR-0046 absent from Project Constraints |
| #5 | pass | — | — |

Gate #4 ran ten minutes to report a single missing citation. After it, the
Supervisor stopped patching what the report named and instead swept every
accepted ADR in `docs/adr/` against both Spec artifacts. That sweep found
**eleven** unnamed, not one.

Two of the eleven were load-bearing and would never have surfaced by patching
what the gate named, because the gate had not reached them yet:

- **ADR-0034** — an own-deadline expiry is an unknown outcome, not a source
  rejection. That is the entire premise of the Spec's fail-closed Coverage: a
  timed-out read must leave a Coverage Gap rather than produce a verdict.
- **ADR-0028** — transactional ERP day deletion, which bounds the manual fix
  flow the Spec deliberately does not automate.

Throughout all three cycles the behavioral score never moved: 112 rows pass,
0 finding-blocked, every Core Feature, Success Metric, Non-Goal and Task
acceptance criterion green. No code was touched after gate #2's findings were
closed. The only thing that ever failed again was prose.

A second-order trap cost a near-miss: the first sweep wrote the exclusions as
ranges — "ADR-0020 through ADR-0024". The audit matches literal identifiers,
so `ADR-0021` still read as absent. Enumerating each identifier individually
was required.

## Root cause

Two compounding causes.

**The report is not exhaustive for a class it can evaluate exhaustively.**
Nothing prevents naming every missing citation in one pass — the check is a
set difference over files. Reporting one instance turns a single defect class
into one Run per instance.

**The gap was not created by authoring the Spec, and no trigger re-examines it.**
Spec 0012 was authored 2026-07-29 and its Project Constraints were complete on
that day. ADR-0045 through ADR-0050 were accepted on 2026-08-03 while
authoring *three other Specs*, and several govern 0012 retroactively. The Spec
became stale without being touched. A per-Spec check that runs before
`implement` catches this only if something re-runs it for every active Spec
whenever an ADR is accepted.

## What would settle it

Spec `0064-spec-artifact-consistency-gate` already proposes the static
pre-`implement` check, and its Core Feature 4 is exactly this row. This
finding is fresh cross-repository evidence for it — 0064's own motivation was
drawn from Roundfix Specs 0056 and 0058, and this is the same failure in an
unrelated repository and domain, which makes it a property of the workflow
rather than of one project.

Two additions 0064 does not currently state:

1. **Exhaustiveness is the property that pays.** 0064's Core Feature 7 fixes
   the reporting *format* (file and line per side of a contradiction). What
   collapses five Runs into one is reporting *every* instance of a row in one
   pass. Worth making explicit, because a compliant implementation could
   report one finding at a time and preserve the whole cost.
2. **Accepting an ADR should invalidate the check for every active Spec**, not
   only the one being implemented. The natural trigger is ADR creation, not
   Spec authoring — that is where the staleness is introduced.

## Related

[[2026-08-04-fail-fast-verification-spends-the-single-repair-turn-on-the-first-of-n-defects]]
is the same shape one altitude down.

[[2026-08-04-what-still-needs-a-supervisor-between-a-prd-and-a-merge]] names the
manual recovery path a failed gate forces (its stop 3). This finding is the
reason that path was walked three times for one defect class.

## Spec pointer

`0064-spec-artifact-consistency-gate` (active, 9 Tasks pending).
