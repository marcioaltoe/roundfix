---
type: perf # feat | fix | perf | refactor
status: promoted # open | promoted | declined
created: 2026-08-03
spec: 0080-cheap-detectors-run-before-the-gate
reason: null # required when status: declined
---

# Verification performance contract for adopting repositories

## Slow

Adopting repositories lack a per-profile contract that keeps the local
implementation loop fast while requiring CI to judge the complete tree from a
fresh run. Without that split, a local gate can make every edit pay the full
suite cost, or a fast local loop can become the only pre-merge verdict.

## Measured

Roundfix's [archived verification-cost report](../specs/_archived/0071-verification-cost/baseline/2026-08-03-after.md)
measured `make verify` at 4.9s on an unchanged tree, about 5s after a
light-package edit, 48.5s after an `internal/cli` edit, and 54.5s after an
`internal/baseline` edit. A complete fresh run took 88.9s. The local tier kept
every test and used Go's test-result cache; CI disabled that cache and ran the
complete suite.

## Target

Add a per-profile Baseline clause that requires a two-tier gate: local
incremental verification completes within 10s on an unchanged tree and within
60s after a typical change, while CI always runs a fresh, complete gate.

## Related

`2026-08-06-two-stage-qa-gate-economics.md` needs exactly this split from the
other side: its mechanical stage is the fast tier, and its audit stage is what
should never re-prove what the tier already proved. The two belong in one Spec.
`2026-08-06-event-journal-payload-economics.md` is downstream of both — agent
turns are what produce journal bytes — but lands on a disjoint surface and can
ship on its own.
