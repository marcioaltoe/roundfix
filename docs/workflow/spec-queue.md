# Spec queue

The one dependency-and-risk-ordered queue of approved Specs that
`docs/agents/autonomous-work.md` requires the Supervisor to maintain. The
Supervisor works this list top to bottom. A Spec leaves the list when it
archives.

Rewritten on 2026-08-26 at the triage recorded in
`2026-08-26-triage-the-queue-earns-its-tokens.md`. The previous order was
written on 2026-08-09 and every Spec on it has since shipped or been removed at
that triage. The criterion this queue is ordered by: reduce Run errors, reduce
rework of finished Tasks, reduce token consumption in the implement loop.

## Order

| # | Spec | State | Why here |
| --- | --- | --- | --- |
| 1 | Unresolved-Run work reuse | to author (lean: techspec-direct) | The largest measured waste of the 2026-08-25/26 session: an Unresolved Run leaves completed Task commits unintegrated on its Run Branch, the next implement reads the checkout, sees `pending`, and redoes them — about 18 redundant Task executions across three Runs of one Spec. Adopts the finding `2026-08-12-five-unresolved-runs-to-deliver-one-spec`, which measured the same shape two weeks earlier. Everything below runs cheaper once finished work survives an Unresolved outcome. |
| 2 | `0116-a-verdict-that-states-its-own-scope` | PRD ready, source adopted, tooling authorized | Five of eight authored Verification commands in Spec 0098 and six of six in Spec 0113 were vacuous or non-hermetic past a clean `spec check`, costing three Unresolved Runs. ADR-0124's shared prober exists; no authoring skill names it. Every Spec authored after this one is checked at authoring instead of one lost Run per bad command. |
| 3 | `0097-a-wave-that-cannot-collide` | PRD ready, needs techspec + graph | Two Tasks a graph declares independent that edit the same file die at integration, discarding a passing Task. Measured first-hand on 2026-08-26: 0113's task_05 completed, then failed on `integration conflict` with sibling task_07. Also carries the authoring rule that wave independence is stated in shared edit targets. |
| 4 | `0105-the-gates-own-economics` | PRD ready, needs techspec + graph | The gate is where the loop's cost concentrates: 201 failed Tasks across five repositories in its PRD, and five gate executions in the 2026-08-25/26 session at 20–40 minutes each. Absorbs the whole of former Spec 0104 (the test-cache target — Task-sized, same authorization umbrella). Deliberately after 1–3: optimising the cost of a verdict is the wrong order while verdicts still discard work. |

## After the queue

The strongest candidate to join once reproduced in this repository is the
Daemon postcondition gap (finding
`2026-08-26-the-daemon-has-no-postcondition-so-a-task-settles-clean-breaking-the-tree`,
deferred): a false Clean is rework material by definition, but its measurement
is from fluxus and this queue only carries locally measured work.

Everything else that was on the 2026-08-12 queue is deferred with its evidence
preserved; the mapping lives in the 2026-08-26 triage record beside this file.
