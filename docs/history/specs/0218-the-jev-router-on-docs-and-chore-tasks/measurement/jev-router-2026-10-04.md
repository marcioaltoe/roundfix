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
| 0200 task_04 | chore | opencode / roundfix-openrouter/typesafe/jev-router / "" | not started | — | — | Not started that afternoon: the OpenRouter account had about US$5.11 of credit left. Ran in the evening; see "Chore replays, 2026-10-04 evening" | — | — | — |
| 0200 task_04 (2026-10-02 default, `run_20261002T212145Z_701d3c4118ab6714`) | chore | codex / gpt-5.6-luna / max | 10m40s | 2 | 1 | Clean; Verification passed on attempt 2 | none | 5,827,127 (request-sum) | — |
| 0195 task_06 | chore | opencode / roundfix-openrouter/typesafe/jev-router / "" | not started | — | — | Not started that afternoon: same account credit. Ran in the evening; see "Chore replays, 2026-10-04 evening" | — | — | — |
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

## Chore replays, 2026-10-04 evening

The maintainer added credit to the OpenRouter account and approved the two
`chore` replays. They ran between 21:31 and 21:53 local time
(2026-10-05T00:31Z to 00:53Z), one after the other, with nothing else using
the router.

**Setup.** Same protocol as the afternoon, with two differences. The binary
was `roundfix 0.38.0` (`bin/roundfix`, built from `13894704` with a dirty tree). The scratch
`.roundfixrc.yml` was the squash commit's own Project Config with five keys
changed: `profiles.docs` and `profiles.chore` set to the router as preferred
and `codex / gpt-5.6-luna / max` as the only fallback,
`worktree.location` under the clone, `notify.enabled: false`,
`budget.max_run_duration: 60m` and `implement.auto_push: false`. The graph
declared `qa: declined` with a `qa_reason`. Each clone's push URL pointed at
a path that does not exist.

The pre-states were rebuilt from the squash commits, `bfebc0e9` (0200) and
`cad8bb25` (0195), and each Task's Verification failed on its pre-state:

- 0200 task_04: `missing pass: TestThisRepositoryHoldsEveryRequiredExternalSkill`.
  The restore brought back `.agents/skills/context7/` and removed
  `context7-cli` and the new test file.
- 0195 task_06: `TestEveryOwnedSkillVersionIsRecorded` failed with
  `open testdata/owned-skill-versions.json: no such file or directory`. As on
  2026-10-02, the restored Roundfix skill was at `0.0.2`, and the version
  record was absent at the first parent.

**Spend guard.** Before each replay the guard read
`GET /api/v1/credits` and `GET /api/v1/key`. It would not start a replay
below US$15 of account credit or above US$40 of key `usage_monthly`. During
each replay it polled the key every 30 seconds and would have interrupted the
Run past US$15. It never fired.

| Read | Account credit left | Key `usage_monthly` | Judge Log, month |
| --- | ---: | ---: | ---: |
| Before 0200 task_04 | US$34.6799 | US$10.0205 | US$9.2339 |
| Before 0195 task_06 | US$26.8121 | US$17.8882 | US$15.0431 |
| After both | US$24.3574 | US$20.3430 | US$16.9907 |

### Measured

| Task | Category | Selection | Wall time | Prompts | Verification repairs | Outcome | Fallback | Tokens | Router cost (Judge Log / key delta) |
| --- | --- | --- | --- | ---: | ---: | --- | --- | --- | --- |
| 0200 task_04 (`run_20261005T003215Z_499a8ed2e919d19a`) | chore | opencode / roundfix-openrouter/typesafe/jev-router / "" | 10m33s | 1 sent, 0 refused | 0 | Clean; Verification passed on attempt 1; Task commit `chore: this repository takes its own update` | none | 183,533 (turn); 5,930,832 (OpenCode request-sum) | US$5.8084 / US$7.8678 |
| 0200 task_04 (2026-10-02 default, `run_20261002T212145Z_701d3c4118ab6714`) | chore | codex / gpt-5.6-luna / max | 10m40s | 2 | 1 | Clean; Verification passed on attempt 2 | none | 5,827,127 (request-sum) | — |
| 0195 task_06 (`run_20261005T004451Z_4b2afc8dd9b272df`) | chore | opencode / roundfix-openrouter/typesafe/jev-router / "" | 7m32s | 1 sent, 0 refused | 0 | Clean; Verification passed on attempt 1; Task commit `chore: the Roundfix skill declares the version its merged content needs` | none | 120,646 (turn); 2,346,544 (OpenCode request-sum) | US$1.9476 / US$2.4548 |
| 0195 task_06 (2026-10-02 default, `run_20261002T212145Z_3102fd99c88c4de8`) | chore | codex / gpt-5.6-luna / max | 3m14s | 1 | 0 | Clean; Verification passed on attempt 1 | none | 687,813 (request-sum) | — |

Wall time spans the whole `roundfix implement` command, Preflight included.
The two 2026-10-02 default replays ran together; these two ran alone.

**Tokens.** The turn figure is what `roundfix runs show` reports. The
request-sum figure adds every OpenCode step of the session and of the one
`task` subagent it started, read from OpenCode's local session database:
input, cache reads, output and reasoning. It is the closest figure to Codex's
request-sum, but the two runtimes count cache differently, so treat it as an
order of magnitude.

| Session | Steps | Input | Cache read | Output | Reasoning |
| --- | ---: | ---: | ---: | ---: | ---: |
| 0200 task_04, main + subagent | 41 + 8 | 699,388 | 5,203,042 | 18,286 | 10,116 |
| 0195 task_06, main + subagent | 17 + 9 | 326,054 | 2,008,048 | 9,339 | 3,103 |

**Cost.** Each Run wrote one `router-prompt` Judge Log line, `outcome: clear`,
no error, with an empty `model` and `provider`. The line records the key's
usage change at the end of the prompt, and OpenRouter kept reporting usage
after that. The key read 60 seconds after each Run was US$2.0594 higher than
the 0200 line and US$0.5072 higher than the 0195 line. It then held still
before the next replay and after the last one. The account's `total_usage`
moved by the same US$10.3225 as the key. So the key delta is the cost of
each replay, and the Judge Log understated tonight's spend by US$2.5666, or
25 %.

| | 0200 task_04 | 0195 task_06 | Total |
| --- | ---: | ---: | ---: |
| Judge Log `cost_usd` | US$5.8084 | US$1.9476 | US$7.7560 |
| Key `usage_monthly` delta | US$7.8678 | US$2.4548 | US$10.3226 |
| Prompt latency | 583,372 ms | 389,448 ms | |

**Routed model.** It was not observable. The Judge Log line has an empty
`model` and `response_id`. OpenCode records only the configured
`roundfix-openrouter/typesafe/jev-router`, because the inline provider does
not return the routed model to it.

**Tool calls.** All tool calls ran through the router. 0200 task_04 completed
111 tool calls (88 main, 23 subagent): `bash`, `read`, `edit`, `skill`,
`grep`, `glob`, `todowrite`, `write` and `task`. 0195 task_06 completed 60
and had 2 `read` errors, both for `skills/testdata/owned-skill-versions.json`,
which is absent in that pre-state. No tool call failed in transport.

**What the routed prompts did.**

- 0200 task_04 ran the public update with an upstream source directory, the
  reconcile preview and confirm, deleted `.agents/skills/context7/`, aligned
  `skills/recommended.txt` and the digest pin, wrote
  `internal/cli/this_repository_skill_set_test.go`, and updated the three
  guides. It ran Doctor under a temporary `HOME`, because that commit's
  binary rejects the host User Config's `jev` key. The Task commit touched
  only declared files and the Task file (21 files).
- 0195 task_06 raised both versions of `.agents/skills/roundfix/SKILL.md`
  to `0.0.4`, synced the bundle and re-recorded the version file. The commit
  touched only the three declared files and the Task file.

### Comparison with the 2026-10-02 default

| | 0200 task_04 routed | 0200 task_04 default | 0195 task_06 routed | 0195 task_06 default |
| --- | --- | --- | --- | --- |
| Outcome | Clean, attempt 1 | Clean, attempt 2 | Clean, attempt 1 | Clean, attempt 1 |
| Prompts / repairs | 1 / 0 | 2 / 1 | 1 / 0 | 1 / 0 |
| Wall time | 10m33s | 10m40s | 7m32s | 3m14s |
| Request-sum tokens | 5.93 M (OpenCode) | 5.83 M (Codex) | 2.35 M (OpenCode) | 0.69 M (Codex) |
| Money | US$7.8678 | — (subscription) | US$2.4548 | — (subscription) |

On 0200 task_04 the router matched the default's wall time and needed one
prompt fewer. On 0195 task_06 it settled the same way in more than twice the
time and about three times the tokens.

## What this shows

Updated after the evening's chore replays.

- The ceiling from Spec 0225 works as configured. The gate read US$50 and
  passed all four routed prompts of the day. Nothing in Roundfix stopped a
  replay.
- The binding limit is the OpenRouter account balance, which the gate
  does not read. With US$40.31 left on the key and about US$5.11 on the
  account, the provider refused mid-Task. Roundfix reported it as an
  `agent/protocol error`, not as a spend refusal, and no fallback ran.
- The router now has `chore` evidence. It settled both `chore` Tasks on the
  first Verification attempt, with no repair, the same outcome as the
  default or better. 0200 task_04 cost US$7.87 and 0195 task_06 cost
  US$2.45.
- Across five routed attempts in two days, the router settled 3 Tasks
  (0210 task_02, 0200 task_04 and 0195 task_06) for US$10.70, or US$3.57 per
  settled Task. 0194 task_04 failed twice, for US$8.86. Counting the
  failures, the router spent US$19.55, or US$6.52 per settled Task. The
  default settled all four Tasks on the subscription, at no marginal cost.
- The router is not faster. It matched the default on 0200 task_04 and took
  more than twice as long on 0195 task_06. On 0210 and 0194 it took three to
  four times as long.
- The Judge Log understates router cost. Each line records the key's usage
  change when the prompt ends, and OpenRouter reported another 25 % in the
  following minute. The gate reads the month's spend from the key, so the
  ceiling still sees the full cost, but a reading of the Judge Log alone
  does not. The routed model remains unobservable.
- The 0194 task_04 pre-state is a poor probe. Restoring only the declared
  files leaves the version ledger at `0.0.6`. The router stops on that
  inconsistency every time and the default works around it, so this replay
  measures caution more than capability.
- No follow-up Backlog Entry. The protocol's condition covers both
  categories, and the router did not settle 0194 task_04, which the default
  settled. For `chore` alone the condition holds, but at US$2.45 to US$7.87
  per Task against a subscription that settled the same Tasks, the evidence
  does not support moving `chore` off the default.

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

## Relay confirmation, 2026-10-05

One more routed replay of 0195 task_06 (`chore`) after Spec 0229 merged. It
checks that the relay records the routed model on the Judge Log. It ran
alone, from 16:04Z to 16:12Z.

**Setup.** Same pre-state protocol and scratch `.roundfixrc.yml` as the
evening's chore replays, rebuilt from `cad8bb25` in a fresh clone. The
Verification failed on the pre-state as before, with
`open testdata/owned-skill-versions.json: no such file or directory`. The
binary was `roundfix 0.41.0` (`bin/roundfix`, built from `7919b3ff` with a
dirty tree). User Config held `jev.monthly_ceiling_usd: 50` and no
`jev.router_min_credit_usd`, so the gate's credit floor was the US$15
default. The spend guard used the same start rules, US$15 of account credit
and US$40 of key `usage_monthly`, and a stop line of US$10 for the replay.
It never fired.

| Read | Account credit left | Key `usage_monthly` | Judge Log, month |
| --- | ---: | ---: | ---: |
| Before | US$22.7644 | US$20.4128 | US$16.9907 |
| After | US$12.7830 | US$30.3942 | US$26.8548 |

**Measured.** `run_20261005T160520Z_0a44478d40386049` reached Clean with
the routed selection: 1 prompt, 0 refused, Verification passed on attempt 1,
no repair, no fallback. The Task commit
`chore: the Roundfix skill declares the version its merged content needs`
touched only the three declared files and the Task file. Wall time was
6m33s; the prompt took 352,282 ms. `roundfix runs show` reports 70,052
tokens (turn).

**The Judge Log line.** One `router-prompt` line, `outcome: clear`, empty
`error`. The fields that were empty in every earlier routed line are now
filled:

- `model`: `deepseek/deepseek-v4.1-flash, openai/gpt-6.1-sol, anthropic/claude-opus-5.5`
- `provider`: `Together, OpenAI, Google`
- `response_id`: the last response's `gen-…` id
- `requested_model`: `roundfix-openrouter/typesafe/jev-router`, as before

The two lists hold distinct values in first-seen order and are not paired,
so the line does not say which provider served which model.

**Cost.** The Judge Log `cost_usd` was US$9.8641. The key's `usage_monthly`
moved by US$9.9814, read 60 seconds after the Run and unchanged on a second
read. The gap was US$0.1173, about 1 %, against 25 % on 2026-10-04 evening.
The same Task settled the same way for US$2.4548 that evening, so this
replay cost four times as much. The routed models included
`anthropic/claude-opus-5.5`. One replay cannot say how the cost split across
the three models.

**After the replay** the account held US$12.7830, below the US$15 floor.
The next routed prompt should be refused before it starts with
`openrouter_credit_low`. This was not exercised.

## Maintainer decision, 2026-10-05

The maintainer reviewed the OpenRouter activity for 2026-10-01 to 2026-10-05
on the `roundfix_jev` key. The routed replays had billed `openai/gpt-6-astra`
(US$15.35), `anthropic/claude-opus-5.5` (US$10.12), `openai/gpt-6.1-sol`
(US$4.81) and DeepSeek (US$0.10), against US$0.014 for 456 direct Jev judge
calls. The maintainer's rule is that OpenAI and Anthropic models are used only
through the Codex and Claude subscriptions, never through OpenRouter. The
`jev-router` routes to those models, so on 2026-10-05 the maintainer decided
to retire the Jev Router ("Aposentar o router"). The direct Jev judge stays.
