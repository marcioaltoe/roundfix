---
date: 2026-08-04
surface: docs/findings, .agents/skills
status: open
---

# A finding does not prevent the recurrence it describes

At 23:45 on 2026-08-03 this repository recorded
`2026-08-03-a-200ms-attach-budget-fails-under-ci-load`: a wall-clock budget
sized against an unloaded machine, failing on a loaded runner. The finding
named the class, named the fix shape, and named the three packages Spec 0071
had already repaired for the same reason.

Four hours later the same session's Supervisor authored this line into Spec
0064's TechSpec:

> A budget assertion keeps Core Feature 1 honest: the full corpus sweep
> completes well inside a second, measured in the test rather than claimed in
> prose.

That is a wall-clock assertion inside `go test -parallel 16 ./...`. It measured
0.64 s alone and 3.0–4.0 s under the suite. It failed the authored QA gate as
F-001, blocked twelve of fifteen QA rows, and cost a full gate cycle plus two
corrective Tasks.

The agent that wrote the defect had written the finding describing it, in the
same session, using the same context.

## What this says

A finding is a record, not a control. Nothing reads `docs/findings/` at
authoring time. The knowledge lands in a file whose only consumer is a human
who happens to open it, so its effect on the next artifact is whatever the
authoring agent happens to recall — which is empirically not enough at a
four-hour distance inside one context window.

The second brain's loop-engineering material states the missing piece
directly (`wiki/concepts/agent-workflows-e-loop-engineering.md`): a loop needs
a **verifier** and **persistent state outside the model's context**, and a
**bilevel or meta-loop** is what lets an outer loop improve the inner loop's
own search. This repository has the inner loop — Spec, Task Graph, Daemon
Verification, authored QA gate — and no outer loop. Findings are the state that
no verifier reads.

The same source names the reason it matters here specifically: *"sem estado,
repete erros"*. Three occurrences of one class — Spec 0071's three packages,
the attach reader, and now the corpus sweep — is that sentence measured.

## The shape of a control

A finding becomes a control when something fails on it. Two candidate
mechanisms, both cheap relative to a gate cycle:

- **A check.** A wall-clock assertion inside a package that runs under the
  parallel sweep is mechanically detectable — `time.Since` compared against a
  duration literal in a `_test.go` file reachable by `go test ./...`. This is
  the same detection class Spec 0064 already builds for Spec artifacts, applied
  to test code.
- **An authoring rule the skill carries.** `write-techspec` and `write-tasks`
  author Testing Approach and Verification sections. A rule stating that a
  timing budget declares its execution conditions would have blocked the line
  above at the moment it was written.

The second is upstream-managed and needs its own authorization; the first does
not. Neither is proposed here as a decision — this finding records why the
current arrangement cannot hold, not which repair to buy.

## Evidence

- The prior finding:
  `docs/findings/2026-08-03-a-200ms-attach-budget-fails-under-ci-load.md`,
  committed in PR #99 at 2026-08-03.
- The recurrence: Spec 0064 `_techspec.md` Testing Approach, PR #98, authored
  the same session.
- The cost: `docs/specs/0064-spec-artifact-consistency-gate/qa/qa-report-2026-08-03.md`,
  F-001, verdict `fail`, 12 of 15 rows blocked, `rows_blocked_environment: 0`.
- Measurements: sweep 0.64 s isolated; 3.036 s under `make verify`; 4.024 s
  under `go test -parallel 16 ./...`.
- Prior art for the class: Spec 0071 repaired load-sensitive wait budgets in
  three packages on 2026-08-03.
- Second brain: `wiki/concepts/agent-workflows-e-loop-engineering.md`;
  `wiki/sources/x-codila-loop-engineering-graph-engineering-2026.md`.
