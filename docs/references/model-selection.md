# Model selection reference

Source snapshot: 2026-08-07; OpenCode runtime facts remeasured 2026-08-08
against opencode 1.18.15 through acpx 0.13.0.
Status: recommendation input, not routing policy.

Sources, each answering a different question:

- [DeepSWE v1.1](https://deepswe.datacurve.ai/) — result, cost, output tokens,
  and steps per agentic coding task.
- [OpenRouter models API](https://openrouter.ai/api/v1/models) — token price,
  cache price, and context window.
- [whatllm.org](https://whatllm.org/explore) — a quality index and per-response
  latency.

No single source answers "which model should this profile use". Task cost hides
latency, token price hides how many tokens a model spends, and a quality index
hides both.

This is the live reference. The 2026-07-16 snapshot inside Spec 0035 is that
Spec's historical artifact and is not maintained; an archived Spec may be deleted
at any time, so the durable table lives here.

Roundfix never routes automatically from this table. It uses only configured
Agent Selection Profiles; a ranking can help a maintainer fill a profile but can
never select or modify one.

## Benchmark name to Agent Selection

The benchmark publishes display names. Roundfix accepts only identifiers the ACP
adapter advertises. They are not the same vocabulary, and a row that cannot be
selected is worth knowing before someone tries.

The codex list below was read from the adapter itself, by offering an invalid
model and capturing its refusal:

```
advertised Agent Models: gpt-5.6-sol, gpt-5.6-terra, gpt-5.6-luna,
                         gpt-5.5, gpt-5.4, gpt-5.4-mini, gpt-5.3-codex-spark
```

The claude list was read from `@agentclientprotocol/claude-agent-acp` 0.63.0 by
opening a session and reading the `model` entry of its advertised
`configOptions`:

```
default              Opus 5 with 1M context · Best for everyday, complex tasks
opus[1m]             Opus 5 with 1M context
claude-fable-5[1m]   Fable 5
sonnet               Sonnet 5 · Efficient for routine tasks
haiku                Haiku 4.5 · Fastest for quick answers
```

That adapter also advertises `effort` as its own option — `default`, `low`,
`medium`, `high`, `xhigh`, `max` — under category `thought_level`. Because a
separate reasoning control exists, roundfix parses advertised model identifiers
as opaque wholes and takes the part before `[` as the canonical model. The
canonical name for `opus[1m]` is therefore `opus`, and reasoning effort comes
from the separate option rather than from the bracket.

Do **not** take model identifiers from `internal/agent/catalog.go`. That
hardcoded list is stale against this adapter: it offers `claude-opus-5` and
`claude-opus-4-8`, neither of which the adapter advertises. The adapter is the
source; the catalog is a copy nothing checks. See
[the finding](../findings/2026-08-07-claude-agent-selections-are-never-proven.md).

| Benchmark name | Runtime | Agent Selection model | Selectable |
| --- | --- | --- | --- |
| `gpt-5.6-sol` | codex | `gpt-5.6-sol` | yes — adapter-advertised |
| `gpt-5.6-terra` | codex | `gpt-5.6-terra` | yes — adapter-advertised |
| `gpt-5.6-luna` | codex | `gpt-5.6-luna` | yes — adapter-advertised |
| `gpt-5.5` | codex | `gpt-5.5` | yes — adapter-advertised |
| `gpt-5.4` | codex | `gpt-5.4` | yes — adapter-advertised |
| `claude-opus-5` | claude | `opus` (advertised `opus[1m]`) | yes — adapter-advertised |
| `claude-fable-5` | claude | `claude-fable-5` (advertised `claude-fable-5[1m]`) | yes — adapter-advertised |
| `claude-sonnet-5` | claude | `sonnet` | yes — adapter-advertised |
| — | claude | `haiku` | yes — adapter-advertised; not on the benchmark board |
| — | claude | `default` | yes — resolves to Opus 5 with 1M context |
| `claude-opus-4.8` | — | — | not via the claude adapter; reachable through `opencode` |
| `kimi-k3`, `qwen3.8-max`, `glm-5.2`, `grok-4.5`, `deepseek-v4-flash`, `muse-spark-1.1`, `gemini-3.6-flash` | opencode | `openrouter/<vendor>/<model>` | yes — see below |

### The opencode runtime reaches everything else

Roundfix supports three runtimes — `codex`, `claude`, `opencode` — and the third
changes what is reachable. Measured 2026-08-08 against opencode 1.18.15 through
acpx 0.13.0, `opencode models` advertises **417** identifiers: 339 under
`openrouter/`, 60 under `opencode/`, and 18 under `opencode-go/`. The 2026-08-07
reading of 431 was taken before the account's provider set changed; the count
moves, so read it rather than quote it.

Every model on the DeepSWE board that neither first-party adapter offers is
selectable this way:

| Requested | Agent Selection under `opencode` |
| --- | --- |
| `anthropic/claude-opus-5` | `openrouter/anthropic/claude-opus-5` |
| `openai/gpt-5.6-sol` | `openrouter/openai/gpt-5.6-sol` |
| `openai/gpt-5.6-terra` | `openrouter/openai/gpt-5.6-terra` |
| `openai/gpt-5.6-luna` | `openrouter/openai/gpt-5.6-luna` |
| `deepseek/deepseek-v4-pro` | `openrouter/deepseek/deepseek-v4-pro` |
| `deepseek/deepseek-v4-flash-0731` | `openrouter/deepseek/deepseek-v4-flash-0731` |
| `z-ai/glm-5.2` | `openrouter/z-ai/glm-5.2` |
| `xiaomi/mimo-v2.5` | `openrouter/xiaomi/mimo-v2.5` |

The same identifiers also appear without the `openrouter/` prefix and under an
`opencode/` prefix for some models; the prefixed OpenRouter form is the one that
names its provider explicitly and is preferred here.

Roundfix offers **no interactive catalog for `opencode`** — `ModelCatalog`
returns nothing for that runtime — so these must be written into a profile by
hand. Combined with the adapter accepting unknown identifiers, a typo here is
silent until a Run.

Reasoning-effort names (`low`, `medium`, `high`, `xhigh`, `max`) match the
benchmark's bracketed suffix and the **codex and claude** adapters' advertised
effort values. OpenCode does not share that vocabulary — see below.

#### The `opencode-go` subscription tier — 2026-08-08

Three tiers hide behind the one runtime, and only the middle one is the
subscription. `openrouter/` is pay-per-use through an OpenRouter key.
`opencode/` — the entries the picker labels *OpenCode Zen*, such as
`gpt-5.6-sol`, `gpt-5.5-pro`, and `grok-build-0.1` — is a separate pay-per-use
tier. `opencode-go/` is what the subscription grants, eighteen models:

`deepseek-v4-flash`, `deepseek-v4-pro`, `glm-5.1`, `glm-5.2`, `gpt-5.6-luna`,
`grok-4.5`, `hy3`, `kimi-k2.6`, `kimi-k2.7-code`, `kimi-k3`, `mimo-v2.5`,
`mimo-v2.5-pro`, `minimax-m2.7`, `minimax-m3`, `qwen3.6-plus`, `qwen3.7-max`,
`qwen3.7-plus`, `qwen3.8-max`.

Two of them bill at **2x usage**: `gpt-5.6-luna` and `deepseek-v4-flash`. That
inverts the obvious guess, because `deepseek-v4-pro` carries no multiplier.

#### OpenCode reasoning effort is per model, and Roundfix does not set it

The `effort` option OpenCode advertises is model-dependent and does not use one
vocabulary. Measured 2026-08-08, each row from `sessions ensure --model <M>`
followed by `sessions show`:

| model | advertised effort values | default |
| --- | --- | --- |
| `opencode-go/kimi-k3` | `max` | `max` |
| `opencode-go/qwen3.8-max` | `high`, `max` | `high` |
| `opencode-go/glm-5.2` | `high`, `max` | `high` |
| `opencode-go/deepseek-v4-pro` | `high`, `max` | `high` |
| `opencode-go/minimax-m3` | `none`, `thinking` | `none` |
| `openrouter/anthropic/claude-opus-5` | `low`, `medium`, `high`, `xhigh`, `max` | `low` |
| `openrouter/openai/gpt-5.6-luna` | `none`, `low`, `medium`, `high`, `xhigh`, `max` | `none` |

A session ensured with no `--model` sits on `opencode/big-pickle` and advertises
**no `effort` option at all**.

Roundfix therefore treats `opencode` as a model-managed reasoning runtime and
refuses any non-empty `reasoning_effort` for it — see ADR-0106. Write
`reasoning_effort: ""` and the model runs at its own default. The reason is
mechanical rather than stylistic: the option only exists once a queue-owner
agent process holds the selected model, which acpx starts on the first prompt,
so a config set issued during a token-free Exact Agent Selection Proof reaches a
transient process on the runtime default and answers ACP `-32602`.

**The adapter refuses nothing.** An unknown model string comes back labelled
`Custom model` rather than rejected, so a typo in a claude Selection survives
configuration and fails later inside a Run. Read identifiers from the adapter
before writing them into a profile; the preview will not catch a mistake.

## Snapshot — 2026-08-07

Sorted by result. Cost is average cost per task; steps and output tokens matter
for wall clock, which cost alone does not express.

| Agent Selection | Result | Cost | Out tok | Steps | Cost per point |
| --- | ---: | ---: | ---: | ---: | ---: |
| `claude / claude-opus-5 / max` | 74% | $11.84 | 118k | 99 | $0.160 |
| `claude / claude-opus-5 / xhigh` | 73% | $9.07 | 92k | 89 | $0.124 |
| `claude / claude-opus-5 / high` | 73% | $6.08 | 64k | 73 | $0.083 |
| `codex / gpt-5.6-sol / max` | 73% | $8.39 | 60k | 61 | $0.115 |
| `codex / gpt-5.6-sol / xhigh` | 71% | $4.70 | 41k | 44 | $0.066 |
| `claude / claude-fable-5 / xhigh` | 70% | $13.41 | 80k | 68 | $0.192 |
| `codex / gpt-5.6-terra / max` | 70% | $3.96 | 72k | 76 | $0.057 |
| `codex / gpt-5.6-sol / high` | 69% | $3.47 | 28k | 37 | $0.050 |
| `claude / claude-opus-5 / medium` | 69% | $3.29 | 37k | 52 | $0.048 |
| `claude / claude-fable-5 / high` | 69% | $9.18 | 57k | 59 | $0.133 |
| **`codex / gpt-5.6-luna / max`** | **67%** | **$0.61** | 73k | 102 | **$0.009** |
| `codex / gpt-5.5 / xhigh` | 67% | $7.23 | 46k | 82 | $0.108 |
| `codex / gpt-5.6-terra / xhigh` | 60% | $1.70 | 40k | 43 | $0.028 |
| `claude / claude-opus-5 / low` | 58% | $1.66 | 20k | 36 | $0.029 |
| `codex / gpt-5.6-luna / xhigh` | 57% | $0.31 | 45k | 71 | $0.005 |
| `codex / gpt-5.6-sol / low` | 45% | $1.07 | 11k | 23 | $0.024 |
| `codex / gpt-5.6-luna / high` | 44% | $0.16 | 26k | 49 | $0.004 |

## What changed since 2026-07-16

| Agent Selection | Cost then | Cost now | Change |
| --- | ---: | ---: | --- |
| `codex / gpt-5.6-luna / max` | $3.03 | $0.61 | −80% |
| `codex / gpt-5.6-terra / max` | $4.95 | $3.96 | −20% |
| `codex / gpt-5.6-sol / high` | $3.47 | $3.47 | unchanged |
| `claude / claude-fable-5 / high` | $9.18 | $9.18 | unchanged |
| `codex / gpt-5.5 / xhigh` | $7.23 | $7.23 | unchanged |

The previous snapshot predates `claude-opus-5` entirely and offered
`claude-opus-4-8` as the Claude option at 52%. Opus 5 now leads the board at
74%, and its `high` setting reaches 73% for half the cost of its `max`.

## Token price and latency — 2026-08-07

DeepSWE reports cost per task. Two other sources report the inputs that produce
it, and they answer questions the task cost hides.

Token price, from the OpenRouter models API, in dollars per million tokens:

| Model | Input | Output | Cache read | Context |
| --- | ---: | ---: | ---: | ---: |
| `openai/gpt-5.6-luna` | $0.10 | $0.60 | $0.01 | 1.05M |
| `openai/gpt-5.6-terra` | $1.00 | $6.00 | $0.10 | 1.05M |
| `openai/gpt-5.6-sol` | $5.00 | $30.00 | $0.50 | 1.05M |
| `anthropic/claude-opus-5` | $5.00 | $25.00 | $0.50 | 1M |
| `anthropic/claude-sonnet-5` | $2.00 | $10.00 | $0.20 | 1M |
| `anthropic/claude-fable-5` | $10.00 | $50.00 | $1.00 | 1M |
| `deepseek/deepseek-v4-pro` | $0.43 | $0.87 | $0.004 | 1.05M |
| `deepseek/deepseek-v4-flash-0731` | $0.09 | $0.18 | $0.018 | 1.05M |
| `z-ai/glm-5.2` | $0.53 | $1.67 | $0.099 | 1.05M |
| `xiaomi/mimo-v2.5` | $0.14 | $0.28 | $0.003 | 1.05M |

`:batch` variants price at half. `-pro` variants exist for luna, terra, and sol;
luna-pro matches luna's price, while sol-pro and terra-pro match their base
models. `claude-opus-5-fast` doubles Opus 5's price.

Luna is **50× cheaper per token** than sol, which is far more than the 5.7×
gap in task cost. Both are consistent: luna is 50× cheaper per token and emits
about 2.6× the output tokens, landing at roughly one-sixth the task cost. Its
cache-read price of $0.01 per million is close to free, which matters
disproportionately in a Run, where the same repository context is re-read across
many turns.

Latency and a quality index, from whatllm.org:

| Model | Intelligence | Response time | Task cost |
| --- | ---: | ---: | ---: |
| Claude Opus 5 | 63.1 | 14s | $0.425 |
| Claude Fable 5 | 62.1 | **1.7 min** | $3.14 |
| GPT-5.6 Sol | 60.9 | 9.4s | $0.237 |
| GPT-5.6 Terra | 56.6 | 5.7s | $0.102 |
| Claude Sonnet 5 | 55.3 | 10s | $0.417 |
| GPT-5.6 Luna | 52.3 | **4.1s** | $0.012 |

That source's task cost uses a different task mix from DeepSWE's and is not
comparable to the table above; only the ordering and the response times carry
across.

## Reading these tables honestly

**Step count alone does not measure wall clock.** `gpt-5.6-luna / max` takes 102
steps against `gpt-5.6-sol / high`'s 37, which invites the conclusion that it is
roughly three times slower. Its per-response latency is 4.1s against sol's 9.4s,
so the two effects nearly cancel:

```
sol / high:   37 steps × 9.4s ≈ 348s
luna / max:  102 steps × 4.1s ≈ 418s        ≈ 20% slower, not 180%
```

Treat that as an estimate, not a measurement. A Roundfix Run inserts Verification
commands between Agent turns, and those gates can dominate the total. Only a
measured Run on this repository settles it, which is why the section below
exists.

**Fable 5's response time is a category apart.** At 1.7 minutes per response it
is roughly 25× sol's latency, on top of being the most expensive model per token
on this list. A frontend profile that prefers it is choosing quality at a wall
clock cost worth being deliberate about.

**Opus 5 is cheaper per output token than Sol** — $25 against $30 — while
scoring higher on both boards. Its cost per task is higher because it emits more
tokens and takes more steps, not because its tokens cost more.

**Cheap does not stop at luna.** Through `opencode`, `deepseek-v4-flash-0731`
costs $0.09/$0.18 per million and `mimo-v2.5` costs $0.14/$0.28 — below luna's
$0.10/$0.60 on output. DeepSWE places `deepseek-v4-flash / max` at 53% for $0.10
per task, against luna's 67% for $0.61. That is the real shape of the low end: a
sixth of luna's task cost for fourteen points less. Whether fourteen points is
affordable depends entirely on what the Task is; a bounded, well-specified slice
with a Verification gate behind it tolerates a weaker model far better than an
exploratory one does, because the gate catches what the model misses.

No model in that tier has been measured on this repository. Their numbers here
are price and benchmark only.

**The benchmark is general coding-agent work.** It publishes no frontend,
backend, QA, or review slice, so any category ordering derived from it combines a
general result with a qualitative task-fit judgment. Every CLI row must keep
reporting `category_specific: false` until a category-specific evaluation exists.

**A result is a snapshot, not a guarantee.** Display the source date with the
numbers, and never let ranking data change a configured profile on its own.

## Local measurements

Benchmark numbers do not predict this repository's Runs: its tasks are larger
than benchmark tasks, its Verification commands gate every commit, and its wall
clock includes those gates. Measurements taken on this repository belong here as
they are made.

| Date | Agent Selection | Spec | Tasks | Wall clock | Outcome | Feedback rounds |
| --- | --- | --- | ---: | ---: | --- | ---: |
| 2026-08-07 | `codex / gpt-5.6-sol / high` | 0082 | 7 of 8 | 1h12m to task_06 | Unresolved at task_07 | 1 |
