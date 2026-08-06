---
type: perf # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-06
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Split the QA gate into a mechanical stage and an audit stage

## Opportunity

The QA gate is the loop's slowest step and the one every fleet project pays
for, because Roundfix is the shared harness. Today it is a single agent audit
that rebuilds its whole criterion matrix from scratch on every round, so a
one-line correction costs a full audit.

Measured on Spec 0079 (2026-08-06): the first Run took 92 minutes for seven
Tasks plus the gate; the second Run, whose only work was moving one declared
constant, took 29 minutes; the third took 30. Inside those rounds the
authoritative repository gate — `make verify` — takes about 90 seconds, and
the third round wrote eleven evidence files from scratch for rows whose
inputs had not moved.

The three defects those rounds found separate cleanly by what could have
caught them:

- A stale declared expectation in a characterization fixture: `make verify`
  alone catches it, in 90 seconds, if it runs right after the Task commit
  instead of at the end of the graph.
- A commit-shape violation (a consequent fix folded into a protected-tooling
  commit): mechanically checkable with `git diff-tree` against the recorded
  authorization, in seconds.
- A PRD promise no clause delivered: genuinely needs a reader comparing
  intent with delivery. This is what the audit stage is for.

Two of the three never needed an agent.

## Fleet evidence

The pattern is not local to this repository. Findings already versioned in
other projects' repositories record the same shape:

- `fluxus` — three consecutive gate rounds refused one Spec; another Spec's
  gate blocked ten of twenty-two rows; gates elsewhere ran 113 rows green.
- `oraculum` — one delivery accumulated twelve QA findings beside
  twenty-three review findings, and records that deferring only transfers
  the cost to the gate, "which finds it by construction".
- `vortex` — records the inversion that produced today's shape: each Task
  proves its own effect and the repository gate runs at Run level.
- `fiscus`, `tax-poc` — gate and verification cost appear in their autonomous
  delivery retrospectives.

This repository's own evidence is consolidated in the rollup
`2026-08-06-rollup-qa-gates-and-verification-evidence.md`, whose nineteen
members include `qa-cycle-cost-is-cold-environments-and-agent-turns`,
`qa-cycle-latency-and-detector-placement`, and `qa-gate-round-economics`.
Nineteen observations over three weeks, none of which became a Spec.

## Shape

A two-stage gate:

1. **Mechanical stage** — the repository gate plus the policy checks that are
   greps and git commands today performed inside agent turns: authorization
   versus changed paths, consequent-fix commit shape, report row shape,
   evidence-path resolution. Runs first, fails fast, costs no agent turn.
2. **Audit stage** — the agent observing behavior, but only over rows whose
   evidence inputs moved since the last report, with unchanged rows carried
   forward by reference to the report that established them.

Open questions the Spec must answer: what makes a row's inputs "moved"
(commit range, changed paths, or declared dependency), how carried-forward
rows stay honest so a stale pass cannot ride along forever, and whether the
mechanical stage belongs in `roundfix` as a command or in the gate skill.

## Priority

Raised to priority on 2026-08-06 by the maintainer, on the condition — since
confirmed above — that the cost is fleet-wide. Next Spec candidate once
Spec 0079 closes.
