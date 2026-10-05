---
type: feat
status: open
created: 2026-10-05
spec: null
---

# A judge-assigned model tier per Task

## Opportunity

The Jev Router (Spec 0218) chose a model for every request inside a Task, and
through OpenRouter it chose `openai/gpt-6-astra`, `openai/gpt-6.1-sol` and
`anthropic/claude-opus-5.5`, which the maintainer's subscriptions already
cover. On 2026-10-05 the maintainer set the rule that OpenAI and Anthropic
models run only through the Codex and Claude subscriptions, and the router was
retired (Spec 0230, ADR-0235). The maintainer found the idea of choosing a
model dynamically worth keeping in another form, for later analysis:

- **Choose per Task, not per request.** A subscription is reachable only
  through its own agent (codex-acp, claude-agent-acp), so one agent session
  cannot mix subscription and OpenRouter requests. The choice has to be made
  before a Task is dispatched.
- **The judge assigns the tier when the Tasks are written.** `write-tasks`
  asks the Jev judge for a `model_tier` (`light`, `standard`, `heavy`) for each
  Task and records it in the Task's front matter, where the maintainer reviews
  and can change it in the planning Pull Request. Dispatch only maps the tier
  to a profile candidate, with no Jev call at run time:
  - `light`: an open model through OpenRouter (DeepSeek first, from a
    configurable allow-list);
  - `standard`: `codex / gpt-5.6-luna`;
  - `heavy`: `codex / gpt-6.1-sol` or `claude / opus`.
- **The graph is unchanged.** The tier is metadata. Dependencies and waves are
  written as today. Without Jev the tier defaults to `standard`, which is
  today's behavior. QA Tasks and Tasks that touch Governed Paths always stay on
  the subscriptions. A `light` Task that fails Verification escalates one tier
  once.
- **Separate OpenRouter keys per stage.** The judge and the implementation use
  their own keys (for example `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` and
  `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY`), each with its own monthly limit, so
  the OpenRouter activity export shows the cost of each stage. The maintainer
  would create the keys.
- **Gate kept.** Any OpenRouter model id under `openai/` or `anthropic/` stays
  refused (ADR-0235).

## Value

Offload simple Tasks from the subscription quota at very low cost, with a
decision the maintainer can see and override before anything runs.

Measured inputs (OpenRouter export `docs/_inbox/openrouter_activity_2026-10-05.csv`
and the archived Task files, 2026-10-05):

- Effective prices per million tokens, cache included:
  `deepseek-v4.1-flash` US$0.05, `typesafe/jev-1.13` US$0.04,
  `gpt-6.1-sol` US$0.55, `gpt-6-astra` US$2.17, `claude-opus-5.5` US$4.23.
- A Task in the measured Runs used 1 to 6 million tokens, so a `light` Task on
  DeepSeek would cost about US$0.05 to US$0.30. The judge's tier call costs
  about US$0.0001 per Task.
- Of 1,247 archived Tasks, about 230 (18%) are natural `light` candidates:
  docs low 91, backend low 96, test low 26, chore low 18. Because they use
  fewer tokens than heavy Tasks, the subscription quota saved would be roughly
  10 to 15%.

## Risks

- DeepSeek has never run a whole Roundfix Task on its own. A failed `light`
  Task adds repair attempts and an escalation, which costs more time than
  sending it to Codex directly. The idea pays off only if `light` Tasks pass
  Verification on the first attempt in the large majority of cases (about 80%
  or more).
- Documentation Tasks in this repository are precise (required phrases, skill
  versions, mirrors), which raises the failure risk for a weaker model.
- Routed OpenCode Runs took 6 to 10 minutes against 3 to 10 minutes on Codex.
  There is no wall-time data for DeepSeek alone.

## Shape

Non-binding. Measure before building:

1. Replay six archived `light` Tasks (two docs low, two backend low, one test
   low, one chore) directly on `opencode` with DeepSeek, no router, and compare
   first-attempt Verification, wall time and cost with their original Codex
   Runs. Estimated cost under US$2.
2. Ask the Jev judge for the tier of about 50 archived Tasks offline and
   compare with the repairs each Task actually needed. Estimated cost about
   US$0.01.
3. Proceed only if first-attempt Verification is at least 80% and wall time at
   most 1.5 times Codex. Otherwise keep only the separate keys.

Any measurement sends only this repository's code and Spec artifacts, uses
only open models and the Jev judge, and stops at US$5.

## Measurement, 2026-10-05

Run on the maintainer's machine by an operator agent, outside any Run
sandbox, with `NODE_OPTIONS` unset and `bin/roundfix` 0.43.0 (`69462a0c`).
OpenRouter was reached only on `ROUNDFIX_OPENROUTER_API_KEY`, for
`deepseek/deepseek-v4.1-flash` and the direct Jev judge. Nothing in Roundfix
was built or changed.

### Setup

Each replay used a fresh clone of this repository and the pre-state protocol
of the Spec 0218 measurement: the Spec's squash merge commit on a scratch
branch, the Task's `interface:` and `creates:` files restored to its first
parent (a file absent there deleted), the Spec moved back under `docs/specs/`
as `active` without archive stamps, the Task `pending` without its `## Result`
and Daemon-written sections, every other Task and `qa/` removed, and a graph
of the one Task with `qa: declined`. Two changes from that protocol: the
clone's local `main` was set to the squash commit's first parent, as the base
`make verify-changed` diffs against, and only Tasks whose declared files no
sibling Task declares were used, because restoring a shared file also undid a
sibling's work (0221 task_03 did not compile that way and was replaced). The
scratch `.roundfixrc.yml` was the squash commit's own with the worktree
location, notifications, a 60-minute budget and `implement.auto_push: false`
changed. Every Task's Verification failed on its pre-state, and every
pre-state compiled.

Each replay ran
`roundfix implement --spec <slug> --no-input --agent opencode --model openrouter/deepseek/deepseek-v4.1-flash --reasoning-effort ""`,
alone, one after the other. OpenCode got the key through
`OPENCODE_CONFIG_CONTENT` as `{env:ROUNDFIX_OPENROUTER_API_KEY}`, with the
generic `OPENROUTER_API_KEY` unset, its MCP servers disabled and `webfetch`
off, so only this repository's files reached the model. A spend guard polled
`GET /api/v1/key` every 30 seconds and would have stopped a Run past US$1.50
per replay or US$34.50 of key usage. It never fired.

Codex rows come from the Run Database: each Task's original Run in its Spec's
delivery, and for 0195 task_06 also the default replays of 2026-10-02 and
2026-10-03. Agent time runs from the Task's start to its first Verification;
Task time runs to settlement.

### Measured

| Task | DeepSeek Run | First attempt | Repairs | Command wall time | Agent / Task time | Tokens (turn; OpenCode request-sum) | Cost | Codex: Run, model, attempts, agent / Task time, tokens |
| --- | --- | --- | ---: | --- | --- | --- | ---: | --- |
| 0195 task_06 (`chore`) | `run_20261005T213554Z_d594ef57185a5418` | passed | 0 | 4m01s | 3m12s / 3m13s | 72,007; 1,677,105 | US$0.0559 | `run_20260930T205335Z_76492e3075ab72f0`, `gpt-5.6-luna`/max, passed on attempt 1, 2m15s / 2m16s, not recorded; default replays 2m41s and 2m39s, 687,813 and 910,160 |
| 0225 task_01 (`docs`) | `run_20261005T214045Z_c3791553ac8aa122` | passed | 0 | 3m01s | 2m03s / 2m04s | 98,922; 4,019,769 | US$0.0741 | `run_20261004T194332Z_e41f982109e26c75`, `gpt-5.6-luna`/max, passed on attempt 1, 3m38s / 8m50s, 915,505; two earlier Runs failed both attempts on the same command |
| 0226 task_01 (`docs`) | `run_20261005T214448Z_c51c76af72c74dae` | passed | 0 | 2m31s | 1m40s / 1m41s | 53,262; 742,466 | US$0.0574 | `run_20261004T181758Z_aba276cda07403d2`, `gpt-5.6-luna`/max, passed on attempt 1, 5m02s / 10m21s, 1,716,418 |
| 0224 task_02 (`backend`) | `run_20261005T220104Z_4d80cc83d23b4efd` | passed | 0 | 5m33s | 4m08s / 4m10s | 78,679; 1,165,608 | US$0.0390 | `run_20261004T162613Z_0bb1b015eb38156b`, `gpt-6.1-sol`/high, passed on attempt 1, 13m56s / 18m43s, 4,678,910 |
| 0190 task_07 (`backend`) | `run_20261005T215049Z_58b3f5952aefe65a` | passed | 0 | 2m31s | 1m29s / 1m33s | 82,671; 1,837,296 | US$0.0471 | `run_20261001T035640Z_f7cf1b9d014c1c3e`, `gpt-6.1-sol`/high, passed on attempt 1, 8m57s / 9m03s, not recorded |
| 0231 task_01 (`test`) | `run_20261005T215435Z_b7fa1b68ea76f19a` | passed | 0 | 5m02s | 2m28s / 3m24s | 61,055; 1,401,045 | US$0.0431 | `run_20261005T182524Z_6b3e6994ac77c127`, `gpt-6.1-sol`/high, passed on attempt 1, 2m29s / 10m02s, 661,689 |

All six Runs reached Clean with one prompt and no repair. Each Task commit
touched only the Task's declared files and the Task file. Cost is OpenCode's
per-session figure; over the six Runs it summed to US$0.3166, within about
2 % of the key's movement for them. The 0224 task_02 replay was retried once:
the first start failed Preflight when the Project Config's `claude` / `opus`
fallback proof timed out under machine load, before any Run existed.

**The gate was weaker than in delivery.** With `qa: declined` the Daemon runs
only the Task's own Verification commands: settlement checks and the
repository command `make verify-changed` run only when the graph has a QA
Task. The 0218 replays had the same property. So `make verify-changed` was run
afterwards on each DeepSeek result in its clone. It passed on all six (4 to 7
minutes each). `roundfix spec check` reported only errors the protocol itself
causes (`SC-COVERAGE-UNTASKED` for the removed Tasks, and older-format codes
on 0190), none on a file the Tasks changed.

**Wall time.** Codex Task time includes `make verify-changed` and the
settlement checks, which the replays did not run, and Codex shared Task and
Verification capacity with sibling Tasks. Agent time is the comparable
figure. DeepSeek's agent time over the six Tasks was 15m00s against Codex's
36m17s, 0.41 times. Per Task the ratio ran from 0.17 (0190) to 1.42 (0195
against its original Run; 1.21 against the two default replays). No Task
exceeded 1.5 times.

**A warm-up prompt did the work at the lowest effort.** A first 0195 task_06
replay asked `--reasoning-effort max`
(`run_20261005T213028Z_fc7faa0a3769cade`, Clean, attempt 1, 4m02s, 120,015
turn tokens, US$0.1302). For an OpenCode model Roundfix first sends the
deferred-effort warm-up prompt `Session setup.` at the model's default
variant, `low`, then applies `max`. OpenCode's session shows DeepSeek
answering `Session setup.` with 38 steps that read the Task and implemented
it, so the `max` prompt only re-validated work done at `low`. The six
measured replays therefore used an empty effort, which sends no warm-up.

### The judge's tier offline

50 archived non-QA Tasks with a recorded original Run were sampled with a
fixed seed: 17 `low` (6 that needed a Verification repair, 11 that passed on
attempt 1), 17 `medium` (5 and 12), 16 `high` (8 and 8). Repairs are
oversampled: 19 of 50 against 50 of 262 in the Run Database. For each, one
request to `https://openrouter.ai/api/v1/systemone` with `model: jev-1.13`,
the client `roundfix spec judge` uses, asked a Choice "the cheapest model
tier that would most likely pass Verification on the first attempt"
(`light`, `standard`, `heavy`, each defined) and a four-level difficulty
Score. The state was the Task file text without `status`, `complexity`,
`## Result` and Daemon-written sections. Every answer came from
`typesafe/jev-1.13-20260917`.

| `complexity` | judged `light` | judged `standard` | judged `heavy` |
| --- | ---: | ---: | ---: |
| `low` (17) | 11 | 5 | 1 |
| `medium` (17) | 0 | 4 | 13 |
| `high` (16) | 0 | 0 | 16 |

- Agreement with the authored `complexity` (`low` → `light`, `medium` →
  `standard`, `high` → `heavy`): 31 of 50. The judge never sent a `medium` or
  `high` Task to `light`, and sent most `medium` Tasks to `heavy`.
- Repairs are not predicted. Judged `light`: 7 of 11 passed on attempt 1;
  `standard`: 6 of 9; `heavy`: 18 of 30. Within `low`, the judge chose
  `light` for 4 of the 6 Tasks that needed a repair and for 7 of the 11 that
  did not. Reweighted to the Run Database's `low` population (43 passed, 9
  repaired), judged-`light` Tasks pass on attempt 1 about 82 % of the time,
  the same as every `low` Task on Codex (83 %).
- The difficulty Score separates `complexity` (means 1.78, 2.79 and 2.93) but
  not outcome: 2.49 for Tasks that passed on attempt 1, 2.50 for those that
  needed a repair.
- On the six replayed Tasks the judge chose `light` for the three `docs` and
  `chore` Tasks only, `standard` for the two `backend` Tasks and `heavy`
  (confidence 0.27) for the `test` Task. DeepSeek passed all six.

The 54 judge calls (50 sampled, 4 replayed Tasks outside the sample) read
122,550 input tokens in the sample and cost US$0.0054.

### Spend

`GET /api/v1/key` read `usage_monthly` US$30.3943 before the first call and
US$30.8123 after the last: **US$0.4181** for the whole measurement, seven
DeepSeek replays and 54 judge calls. The account held US$12.16 of credit
afterwards.

### Decision against the criterion

The criterion holds on this sample: 6 of 6 `light` Tasks passed Verification
on the first attempt (at least 80 % required) and DeepSeek's agent time was
0.41 times Codex's, at most 1.42 times on any Task (at most 1.5 times
required), for about US$0.05 per Task. **Go** for the `light` tier, with four
conditions taken from this measurement:

- Six Tasks is a small sample, all with a single, well-specified Verification.
  Count the `light` escalations from the first real Runs before widening the
  allow-list.
- A `light` selection uses an empty reasoning effort, or the warm-up prompt
  has to become inert first: with an effort set, the warm-up does the Task at
  the lowest effort.
- The replay protocol gates less than delivery does. A `light` Task in a real
  Run also faces `make verify-changed` and the settlement checks; the six
  results passed them afterwards.
- The judge's tier adds nothing measurable over `complexity: low`: its
  `light` Tasks pass on the first attempt as often as any `low` Task, and it
  keeps Go code Tasks off `light` that DeepSeek settled. Map `light` from
  `complexity: low` (QA and Governed Paths excluded, as the Opportunity
  says) and keep the judge's tier advisory in the planning Pull Request until
  it is measured against `light` outcomes rather than Codex repairs.
