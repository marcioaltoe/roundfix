---
type: fix # feat | fix | perf | refactor
status: promoted
created: 2026-08-31
spec: 0126-agent-review-before-pull-request
reason: null
---

# The review Agent rewrites the contract it was asked to satisfy

## Opportunity

Two guardrails on the review-resolution Agent, split out of Spec 0105 on
2026-08-31 under that Spec's own Open Question, which named splitting as the
default if decomposition could not size seven Core Features into single-session
slices.

The split is by subsystem, not by convenience. Every other feature in 0105 lives
in the QA gate, the Spec checker, or the authoring skills. These two live in
`internal/rounds`, which those features never touch — and which today carries no
reference to Verification at all, so the guardrail is new behaviour in a
subsystem rather than an adjustment to an existing one.

**An authored Verification is not the review Agent's to rewrite.** When the
Agent resolving a Review Issue edits a Task's `## Verification`, it changes the
contract it was asked to satisfy. Both measured edits failed correct work: the
Agent read a failing gate, concluded the gate was wrong, and relaxed it. Either
the edit is forbidden, or it is surfaced as a contract change in the Round
report — 0105's Open Question left that undecided, with forbidding as the
default because both measured cases argue for it.

**A finding that needs absent infrastructure is not the review Agent's to
resolve.** Choosing a test substrate is Spec scope. The cheap trigger is a newly
introduced environment variable: an Agent that invents one to make a finding go
away has decided something the Spec owns.

## Value

Not independently measured here. The two Verification edits come from Spec
0105's own evidence set, and their cost is a correct implementation failed by a
loosened contract — the same class the repository has spent this month removing
elsewhere, where a gate that can be talked out of its verdict is not a gate.

Whoever promotes this should measure it in this repository first rather than
inheriting 0105's numbers, which were taken across five repositories for a
different claim.

## Shape

Non-binding. The undecided question travels with it: forbid the edit outright,
or permit it and report it as a contract change. Forbidding is simpler to prove
and harder to work around; reporting keeps a legitimate correction possible at
the cost of a reader who must notice it.

Note that a Verification the Daemon derives rather than an author writes — Spec
0105 Core Feature 3, for the QA Task specifically — removes this risk for that
one Task without touching the general case.

## Addendum — 2026-09-08 — Current triage

Spec 0105 explicitly excluded these resolution-Agent guardrails. Current
`internal/rounds` provides issue identity and status handling but no guard for
authored Verification or a newly invented test substrate.

Keep the entry open for provisional P7, review by Codex or Claude before the
Pull Request. Removing CodeRabbit does not remove this Agent scope hazard.
The forbid-versus-report decision and accepted-risk handling still require the
maintainer's bounded planning decision; this addendum grants neither behavior.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this remaining intent for the implementation queue.
[0126-agent-review-before-pull-request](../_prd.md) is its primary owner.
The source moves once into that Spec's reference index; other Specs link this
owned copy. Its lifecycle status records adoption, not implementation or QA
completion. Governed mutation and execution still require the consuming grant.
