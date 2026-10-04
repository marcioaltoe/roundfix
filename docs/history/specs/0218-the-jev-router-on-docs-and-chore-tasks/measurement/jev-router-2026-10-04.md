# Jev Router measurement, 2026-10-04 addendum

Date: 2026-10-04

This addendum re-runs the routed replays that the
[2026-10-02 record](jev-router.md) could not finish: 0194 task_04 (`docs`),
0200 task_04 (`chore`) and 0195 task_06 (`chore`). The month's ceiling is now
US$50, set by the maintainer in User Config (`jev.monthly_ceiling_usd: 50`,
Spec 0225), and the OpenRouter key reports `limit: 50` with
`limit_reset: monthly`. The 2026-10-02 default rows are the comparison; no
default replay ran again.

## Setup

Run on the maintainer's machine by an operator agent, outside any Run
sandbox, with `NODE_OPTIONS` unset and `ROUNDFIX_OPENROUTER_API_KEY` loaded
from the interactive shell environment. The binary was `roundfix 0.37.0`
built from `origin/main` at `f0353479`.

Each replay used a fresh clone of this repository under the operator's
temporary directory. The pre-state followed the TechSpec's measurement
protocol, as on 2026-10-02:

- The Spec's squash merge commit was checked out on a scratch branch.
- The Task's declared `interface:` and `creates:` files were restored to that
  commit's first parent. A file absent from the first parent was deleted.
- The Spec directory moved back under `docs/specs/` with `status: active`,
  without its archive stamps. The Task was set to `pending` with its
  `## Result` and Daemon-written sections removed.
- Every other Task file and `qa/` were removed. The graph held the one Task
  and declared `qa: declined`.
- The scratch `.roundfixrc.yml` named
  `opencode / roundfix-openrouter/typesafe/jev-router / ""` as the preferred
  `docs` and `chore` selection, with `codex / gpt-5.6-luna / max` as the only
  fallback. Run Worktrees went under the clone's directory, notifications
  were off, the Run budget was 60 minutes and `implement.auto_push` was false.
- One commit in the clone held that state.

Each Task's Verification failed on its pre-state, so no reserve was needed.
Each routed replay ran `roundfix implement --spec <slug> --no-input`.

A spend guard ran beside each replay. It read the month's Judge Log and the
key's `usage_monthly` through `GET https://openrouter.ai/api/v1/key` before
the replay. During the replay it polled the key every 30 seconds and would
have interrupted the Run past US$15 for one replay. It never fired. Every
clone was removed after its row was recorded.

## Measured

| Task | Category | Selection | Wall time | Prompts | Verification repairs | Outcome | Fallback and reason | Tokens | Router cost |
| --- | --- | --- | --- | ---: | ---: | --- | --- | --- | ---: |
| 0194 task_04 (`run_20261004T212840Z_25315d03ee4a1279`) | docs | opencode / roundfix-openrouter/typesafe/jev-router / "" | 7m02s | 2 sent, 0 refused by the gate | 1 attempted; died at OpenRouter | Unresolved; Task failed: `Agent Batch failed after acpx exited with code 1: agent/protocol error` | none: work had begun, so the Work Item failed and no fallback started | 79,070 (turn, 1 of 2 prompts) | US$2.0430 |
| 0194 task_04 (2026-10-02 default, `run_20261002T210520Z_d60410842cfb639f`) | docs | codex / gpt-5.6-luna / max | 4m35s | 1 | 0 | Clean; Verification passed on attempt 1 | none | 1,800,531 (request-sum) | — |
| 0200 task_04 | chore | opencode / roundfix-openrouter/typesafe/jev-router / "" | not started | — | — | Not started: the OpenRouter account had about US$5.11 of credit left | — | — | — |
| 0200 task_04 (2026-10-02 default, `run_20261002T212145Z_701d3c4118ab6714`) | chore | codex / gpt-5.6-luna / max | 10m40s | 2 | 1 | Clean; Verification passed on attempt 2 | none | 5,827,127 (request-sum) | — |
| 0195 task_06 | chore | opencode / roundfix-openrouter/typesafe/jev-router / "" | not started | — | — | Not started: same account credit | — | — | — |
| 0195 task_06 (2026-10-02 default, `run_20261002T212145Z_3102fd99c88c4de8`) | chore | codex / gpt-5.6-luna / max | 3m14s | 1 | 0 | Clean; Verification passed on attempt 1 | none | 687,813 (request-sum) | — |

The 0194 replay wrote two `router-prompt` Judge Log lines:

| Prompt | Latency | Logged tokens in/out | Cost | Outcome |
| --- | ---: | --- | ---: | --- |
| Implementation | 245,342 ms | 184 / 60 | US$0.6794 | `clear` |
| Verification repair | 118,071 ms | 0 / 0 | US$1.3637 | `skipped` |

Both lines record an empty `model` and `provider`. Their sum, US$2.0430,
equals the key's `usage_monthly` change across the replay, from US$7.6484 to
US$9.6914.

**Why the repair prompt died.** OpenCode's log shows the provider refusing a
request at 21:34:54Z with
`AI_APICallError: This request would exceed your available credits given your current in-flight requests. Retry after in-flight requests settle, or add credits.`
The Roundfix gate did not refuse it. The month's spend was US$9.69 of the
US$50 ceiling, and the key had US$40.31 of its limit left.
`GET https://openrouter.ai/api/v1/credits` then reported `total_credits: 210`
and `total_usage: 204.885`, so the account behind the key held about US$5.11.
The key limit is not the binding bound. The account balance is.

**Why the chore replays did not start.** Their pre-states were built and
their Verification failed as required. With about US$5.11 left on an account
that other keys may share, each routed prompt would likely meet the same
provider refusal and spend the rest of that balance on a failure. Neither
replay was started. Both need more account credit first.

**What the routed prompt did.** It repeated the 2026-10-02 behaviour. The
first prompt read the Task, found that the restored version ledger already
records `write-tasks` `0.0.6`, and that the Task does not declare
`skills/testdata/owned-skill-versions.json`. It wrote only a `## Result`
asking the Spec owner to reconcile, and changed no implementation file.
Verification failed on attempt 1 with
`missing pass: TestWaveCollisionAllowsDifferentCommandFiles`. The repair
prompt started inspecting the ledger history and died at OpenRouter before
it edited anything.

**Tool calls worked again.** OpenCode ran `skill`, `read`, `glob`, `grep`,
`bash`, `edit`, `todowrite` and `task` through the router: 14 completed tool
calls in the first prompt, 13 in the repair, none failed.

Month's spend after this addendum: Judge Log US$9.2339 (US$7.1909 before),
key `usage_monthly` US$9.6914, both under the US$50 ceiling and under the
US$40 stop line.

## Comparison with 2026-10-02

| | 2026-10-02 routed | 2026-10-04 routed | 2026-10-02 default |
| --- | --- | --- | --- |
| 0194 task_04 outcome | Unresolved, repair refused by the US$5 ceiling | Unresolved, repair refused by OpenRouter credit | Clean, attempt 1 |
| 0194 task_04 wall time | 15m41s | 7m02s | 4m35s |
| 0194 task_04 router cost | US$6.8120 | US$2.0430 | — (subscription) |
| Implementation prompt | 904 s, US$6.81 | 245 s, US$0.68 | — |

Two routed attempts on the same pre-state stopped at the same place, at
different prices. The router's first prompt cost ten times less on 2026-10-04
than on 2026-10-02 for the same outcome. Since neither line records the
model, the price spread cannot be attributed to a model choice.

## What this shows

- The ceiling from Spec 0225 works as configured. The gate read US$50 and
  passed both prompts. Nothing in Roundfix stopped this replay.
- The binding limit is now the OpenRouter account balance, which the gate
  does not read. With US$40.31 left on the key and about US$5.11 on the
  account, the provider refused mid-Task. Roundfix reported it as an
  `agent/protocol error`, not as a spend refusal, and no fallback ran.
- Across three routed attempts in two days, the router settled 1 Task
  (0210 task_02), at US$0.3730. 0194 task_04 failed twice, for US$8.86 in
  total. The default settled 0194 task_04 in 4m35s on the subscription.
- The 0194 task_04 pre-state is a poor probe. Restoring only the declared
  files leaves the version ledger at `0.0.6`. The router stops on that
  inconsistency every time and the default works around it, so this replay
  measures caution more than capability.
- There is still no `chore` evidence for the router.
- No follow-up Backlog Entry: the router did not settle every Task the
  default settled.

## Cost-reduction proposal

Measure-and-propose only; nothing here changes code or configuration.

1. **Read the account balance in the gate.** Have the gate take the lower of
   the key's `limit_remaining` and the account's
   `total_credits − total_usage` from `GET /api/v1/credits`. Refuse before the
   first prompt with a named reason such as `openrouter_credit_low`, instead
   of failing mid-Task as a protocol error.
2. **Treat a provider credit refusal as a spend refusal.** When OpenRouter
   answers "exceed your available credits", classify it like
   `jev_ceiling_reached` and let the category fallback (codex luna) take the
   Task when no work has begun. This avoids paying for a Task that then fails
   anyway.
3. **Bound each routed prompt.** The gate runs between prompts; one prompt
   cost US$6.81 on 2026-10-02. Use a dedicated measurement key whose limit
   equals the per-replay budget (for example US$5), so OpenRouter caps an
   in-flight prompt. Also measure whether a per-request output token cap in
   the inline provider config lowers cost without breaking tool calls.
4. **Record the routed model.** Both lines log `model: ""`. Read the model
   back per generation (OpenRouter's generation lookup by response id) so the
   US$6.81 vs US$0.68 spread can be traced. If one expensive model drives the
   cost, the next step is to restrict the router's candidates or set a price
   ceiling in the request's provider preferences, after checking the current
   OpenRouter documentation.
5. **Stop paying for a known-bad probe.** Before any further routed spend,
   fix the 0194 task_04 pre-state by also restoring
   `skills/testdata/owned-skill-versions.json` to the first parent. Or switch
   to the protocol's first reserve, 0202 task_04. Then run the two `chore`
   replays once the account holds at least the US$15 per-replay stop line.
6. **Keep `docs` and `chore` on the subscription default.** Its marginal cost
   is zero, and it settled every Task in this sample. Revisit only once the
   router has settled a `chore` Task.
