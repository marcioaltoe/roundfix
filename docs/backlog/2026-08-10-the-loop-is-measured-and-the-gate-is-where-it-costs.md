---
type: perf # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The loop is measured, and the gate is where it costs

## Opportunity

Measured across five repositories on 2026-08-10, from `run_events` since
2026-07-01. Recorded here because the numbers took a full session to gather and
answer questions that opinion had been answering badly.

**Task implementation is not the problem, and roundfix is not an outlier.**

| repository | tasks | implementation failures | QA gate failures |
| --- | --- | --- | --- |
| roundfix | 328 | 12.2% | 35 |
| vortex | 155 | 10.3% | 24 |
| oraculum | 138 | 6.5% | 31 |
| fluxus | 124 | 13.7% | 17 |
| fiscus | 77 | 11.7% | 13 |

**Of 211 failed tasks since 2026-08-03, 123 are the QA gate returning a
verdict, not code breaking:** 104 `fail`, 18 `partial`, 1 unreadable. Only 46
are a declared Verification failing — about 7% of 667 settled tasks. Twenty are
agent or session infrastructure, seven are vacuous Verification, five are
timeouts and worktree bootstrap.

**The gate is not the bottleneck per task.** With the model held constant at
`codex/gpt-5.6-sol` across the whole period, the agent's median is 623–782s per
task and did not degrade; the declared Verification's median is 2–110s, 0.3% to
32% of the total. The pre-work probe costs 5–40 seconds and spends zero tokens:
a refused task opens no Agent Session at all, confirmed by a Run that recorded
zero `agent.*` events.

**Specs mostly land first time.** Of 78 archived Specs, 51 passed with a single
gate; 16 needed two, 5 needed three, 4 needed four, 2 needed five. First-gate
verdicts were 54 `pass` / 16 `fail` / 7 `partial`. Final verdicts carry *more*
failures than first ones — 23 vs 16 — because some Specs archive on
`qa_override`, which is a maintainer decision rather than a gate defect.

## Value

The expensive tail is the 27 Specs that needed a gate rerun, and its cause is
one family: **the Spec assumes the world behaves like its fakes.** Spec 0091
demonstrated it end to end on 2026-08-10 — its characterization corpus ran
against a fake harness, so a design premise survived authoring, implementation
and every unit test, and died at the QA gate four Runs later:

- F-001: the PRD promised Roundfix's membership verdict would own the refusal.
  Against live adapters, all three refuse first. Not a bug — a premise that
  execution does not support, now recorded in ADR-0119.
- F-002: Claude normalizes a requested `opus` to the advertised `opus[1m]`. No
  fake did that, so proof compared raw values and rejected a valid selection.
  Proved a regression by running the same command with a binary built from
  `main`.

`oraculum` shows the same shape without any adapter: the lowest implementation
failure rate of all five repositories, 6.5%, and the second-highest count of QA
gate failures. Code passes its own tests; the gate still finds things. There the
"real world" is a database with data rather than an external process, which is
why arranging it is harder and matters more.

## Shape

Two directions, both non-binding.

**Characterization should touch the real boundary.** When a Spec crosses an
external surface — an ACP adapter, an HTTP contract, a database — its
characterization Task should record what the real thing does, not what a fake
does. On 0091 that single change would have surfaced F-001 on day one instead
of after four Runs. For data-shaped projects the equivalent needs a prepared
mass, which is infrastructure rather than method: the content-addressed golden
fixture pattern used for this repository's test performance on the same day
applies directly — build the mass once, name it by the hash of its inputs, copy
it per test.

**The gate's own economics deserve the same measurement this entry applied to
the loop.** The QA gate on Spec 0091 took 27 minutes per run and three runs to
converge. Nothing here measured how much of that is evidence-gathering versus
re-reading artifacts a cheap check could settle, which is the question ADR-0117
already answered for authoring defects and has not been asked of the gate.

Worth settling in the same work: `SC-VERIFY-VACUOUS-COMMAND` was shipped on
2026-08-10 and immediately disabled by maintainer decision until this work
happens. It found ten vacuous commands across Specs 0080, 0081 and 0085 on its
first run, each of which would have burned a daemon cycle. While it is off, the
Daemon's pre-work probe still refuses those commands — the discovery just moves
back to Run time, at 5–40 seconds and zero tokens per refusal instead of 0.04
seconds at authoring. Re-enabling it is one line in the staged detector
registry.
