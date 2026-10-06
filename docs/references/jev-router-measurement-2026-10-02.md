<!-- Moved upstream on 2026-10-06 from docs/history/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router.md, so the measurement outlives the archived Spec (docs/agents/docs-layout.md). -->

# Jev Router measurement

Date: 2026-10-02

Run on the maintainer's machine by the operator, with `bin/roundfix` 0.29.0
(`139e488f`, the Task 01–03 tree), `NODE_OPTIONS` unset and
`ROUNDFIX_OPENROUTER_API_KEY` present in the environment. Each replay ran in a
fresh scratch clone of this repository under `/private/tmp/claude-501/`, built
as the TechSpec's measurement protocol describes. The Spec's squash merge
commit was checked out on a scratch branch. The Task's declared `interface:`
and `creates:` files were restored to the commit's first parent. The Spec
directory was moved back under `docs/specs/` with that Task `pending`, its
`## Result` removed, every other Task and the QA Task removed from the graph,
and `qa/` deleted. One commit inside the clone held that state and the scratch
`.roundfixrc.yml`. That config sends Run Worktrees under the clone's temporary
directory, disables notifications and bounds a Run at 60 minutes. Each
Verification failed on its pre-state, so no reserve was needed.

The routed `.roundfixrc.yml` named
`opencode / roundfix-openrouter/typesafe/jev-router / ""` as the category's
preferred selection, with the current default
`codex / gpt-5.6-luna / max` as its only fallback. Each default replay ran
`roundfix implement --spec <slug> --no-input --agent codex --model gpt-5.6-luna --reasoning-effort max`.
Every clone was removed after its row was recorded.

## Measured

| Task | Category | Selection | Wall time | Prompts | Verification repairs | Outcome | Fallback and reason | Tokens | Router cost |
| --- | --- | --- | --- | ---: | ---: | --- | --- | --- | ---: |
| 0210 task_02 (`run_20261002T205607Z_b68a9687429eb1e3`) | docs | opencode / roundfix-openrouter/typesafe/jev-router / "" | 8m28s | 1 | 0 | Clean; Verification passed on attempt 1 | none | 98,418 (turn) | US$0.3730 |
| 0210 task_02 (`run_20261002T205709Z_b76342c69208cb31`) | docs | codex / gpt-5.6-luna / max | 2m40s | 1 | 0 | Clean; Verification passed on attempt 1 | none | 442,780 (request-sum) | — |
| 0194 task_04 (`run_20261002T210509Z_23b172a783e9d8e1`) | docs | opencode / roundfix-openrouter/typesafe/jev-router / "" | 15m41s | 1 sent, 1 refused | 1 attempted, refused by the gate | Unresolved; Task failed: `jev_ceiling_reached: month's Jev spend US$7.4408 of US$5.0000` | none: the gate refused the repair prompt after work had begun, so the Work Item failed and no fallback started | 151,936 (turn) | US$6.8120 |
| 0194 task_04 (`run_20261002T210520Z_d60410842cfb639f`) | docs | codex / gpt-5.6-luna / max | 4m35s | 1 | 0 | Clean; Verification passed on attempt 1 | none | 1,800,531 (request-sum) | — |
| 0200 task_04 (`run_20261002T212145Z_701d3c4118ab6714`) | chore | codex / gpt-5.6-luna / max | 10m40s | 2 | 1 | Clean; Verification passed on attempt 2 | none | 5,827,127 (request-sum) | — |
| 0195 task_06 (`run_20261002T212145Z_3102fd99c88c4de8`) | chore | codex / gpt-5.6-luna / max | 3m14s | 1 | 0 | Clean; Verification passed on attempt 1 | none | 687,813 (request-sum) | — |

Routed replays of 0200 task_04 and 0195 task_06 were not started. The gate
refused a routed prompt in the 0194 task_04 replay with
`jev_ceiling_reached`, and the protocol starts no further routed replay after
a refusal. Their default replays ran, because they send nothing to
OpenRouter.

Router cost is the sum of each Run's `router-prompt` Judge Log lines. Each
routed Run wrote exactly one line, `outcome: clear`, with no error. Prompts,
repairs and tokens come from `roundfix runs show <run>` and the Run's event
output. Wall time spans the whole `roundfix implement` command, Preflight
included. Two replays ran at a time (each Task's routed replay, then its
default replay; the two chore defaults together), so wall times include some
contention.

Tokens are not comparable across runtimes. OpenCode reports a `turn` basis,
and Codex reports a `request-sum` across every request of the prompt. The
router lines record only 966 in and 51 out tokens for 0210, and 3,693 in and
75 out for 0194, so the router's token figures cover only the last turn.

Month's Jev spend after the measurement, as the gate computed it: US$7.4408
of the US$5.0000 ceiling. The two `router-prompt` lines log US$7.1850. Before
the measurement, the month's Judge Log held about US$0.0059.

## Reading

**Cost per settled Task on the router: US$0.3730**, from one settled Task. The
router settled 1 of the 2 Tasks it attempted, and the two routed prompts cost
US$7.1850 in total. Counting the failed Task, the cost is US$7.19 per settled
Task. Both Tasks settled on the default with no router spend.

**Tool calls worked through the router.** In both routed sessions OpenCode
ran its tools through `roundfix-openrouter/typesafe/jev-router`: `skill`,
`read`, `glob`, `grep`, `bash`, `edit` and `todowrite`. The 0210 session had
31 completed tool calls and the 0194 session had 52, with no failed tool
call. The 0210 edits passed Verification on the first attempt.

**The router is slower and far more expensive on this sample.** The 0210 docs
Task took 8m28s on the router against 2m40s on the default, with one
474-second prompt. The 0194 docs Task took 15m41s on the router against 4m35s
on the default, with one 904-second prompt costing US$6.81. That prompt did
not implement the Task. It found that the replay's version ledger already
recorded `0.0.6` for the `write-tasks` skill, stopped to ask whether to use
`0.0.7`, and wrote only the Task's `## Result`. On the same pre-state, the
default implemented the Task and passed Verification on its first attempt.

**The ceiling held only after the fact.** The gate checks spend before each
prompt and found US$0.379 of US$5 before the 0194 prompt. That single
in-flight prompt then carried the month to US$7.44, past the ceiling. This is
the "A prompt in flight" risk in the TechSpec, and it was much larger than
expected. The gate refused the next prompt, the repair, as designed. Because
work had begun, the Work Item failed and no fallback ran.

**What the sample can say.** Two routed replays, both `docs`, are enough to
show three things. The routed selection works end to end through the inline
provider: Preflight proof, session, tools, the Judge Log line and the gate's
refusal. One routed Task can cost more than the whole monthly Jev ceiling.
On these two Tasks the router settled fewer Tasks than the default, and took
three to four times its wall time.

**What it cannot say.** It gives no `chore` evidence for the router. It cannot
say which model the router picked on each turn, because the line records an
empty `model`. It cannot give a stable per-Task cost: n = 1 settled Task, and
OpenRouter's usage lag can move cost between adjacent lines, although their
sum stays right. Its tokens cannot be compared with Codex's. It also says
nothing about the default's money cost, because the Codex subscription
reports none.

Two pre-state caveats come from the protocol's reset of declared files to the
squash commit's first parent:

- 0194 task_04 started with a version ledger that already recorded a later
  `write-tasks` version. The router stopped on that inconsistency; the
  default worked around it.
- 0195 task_06's restored skill was at `0.0.2`, not the `0.0.3` the Task
  describes, and its version record was absent. The default settled the Task
  anyway.

**No follow-up Backlog Entry.** The condition did not hold. The default
settled 0194 task_04 and the router did not, and the two chore Tasks have no
routed replay. On this evidence the router should not be proposed for `docs`
or `chore`. The month's Jev spend also already exceeds the shared US$5
ceiling, so every routed prompt and every Jev judgment is refused for the
rest of October 2026. Any further measurement needs a new month or a
maintainer decision on the ceiling. It also needs a per-prompt spend bound,
because the gate cannot stop a prompt that is already running.
