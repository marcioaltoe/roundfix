---
type: feat # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-08
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# Record what each Agent Session consumed, next to which selection ran it

## Opportunity

Roundfix persists the effective Agent Selection for every Task — ACP Runtime,
Agent Model, and reasoning effort — and records nothing about what that Session
consumed. Run Events cover `task-status`, `agent-selection`, `verification`, and
`outcome`; no production path carries a token count. This intent came through
`inbox/roundfix/2026-08-08-uso-por-task-com-agent-modelo-e-reasoning.md` in the
Secondbrain.

## Value

Concurrency and reasoning effort are different knobs and a provider's daily
total cannot separate them: Task Capacity changes how fast a quota window
empties, while reasoning effort changes how much one unit of work costs. The
number that discriminates is output per completed Task, grouped by selection —
which no artifact in the repository can produce today.

The gap was measured on 2026-08-08, when the Codex quota was exhausted mid-Spec
and the question "was it concurrency or reasoning?" could only be answered by
inference over billing aggregates: output over input moved from 11.9% in the
week of 26/07–01/08 to 13.4% in 02/08–08/08, while total work nearly doubled.
The conclusion — mostly volume, with a small reasoning rise consistent with a
`max`-effort model entering the review profile — remains unverified, because the
row linking a Task Type to its cost does not exist.

Recording it turns four recurring questions into queries: cost per Task Type,
cost of the same model at different reasoning efforts, cost of Verification
retries, and how much of a Spec was spent on Tasks that failed and were redone.
That last one is invisible today and is not small: four Runs of Spec 0084 ran on
2026-08-08 and three failed on task-authoring defects, one of them discarding a
Task Worktree with twenty-two finished files.

## Shape

A future design could attach consumption to each Agent Session alongside the
identity already persisted: input and output tokens with cache reads separated
when the adapter reports them, the owning Task, Task Type and Batch, whether the
Session ran a Preferred Selection or a Fallback Selection, and whether it was
the initial turn or a Verification retry. A missing measurement must stay
observably absent rather than becoming a zero, which would read as "spent
nothing".

Whether ACP adapters expose usage at all is the open question the owning Spec
must settle first: `codex-acp` and `claude-agent-acp` need checking before any
event shape is fixed. This shape is non-binding.
