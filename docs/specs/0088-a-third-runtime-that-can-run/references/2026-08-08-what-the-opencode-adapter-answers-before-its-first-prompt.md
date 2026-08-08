---
date: 2026-08-08
origin: direct measurement against a live OpenCode installation
runtime: opencode 1.18.15
adapter: acpx 0.13.0
machine: darwin 25.5.0, node v22.23.2
---

# What the OpenCode adapter answers before its first prompt

Read-only measurement of the `opencode` ACP surface, taken to explain why
Roundfix could open no Agent Session on it. Every command below was run against
a scratch working directory outside this repository. No Roundfix code was
changed to produce these numbers, and no repository file was mutated. The
Roundfix binary used for the two `roundfix` rows was built from
`ma/specs-0082-0083` at commit `2495d3cd`.

This is the External Acceptance Evidence for Spec 0088 under ADR-0104: it
originates outside the Spec's own artifacts, and it measured the runtime rather
than the design.

## The advertised catalog

```
$ opencode models | wc -l
417
$ opencode models | cut -d/ -f1 | sort | uniq -c | sort -rn
 339 openrouter
  60 opencode
  18 opencode-go
```

The eighteen `opencode-go` entries are the subscribed tier:
`deepseek-v4-flash`, `deepseek-v4-pro`, `glm-5.1`, `glm-5.2`, `gpt-5.6-luna`,
`grok-4.5`, `hy3`, `kimi-k2.6`, `kimi-k2.7-code`, `kimi-k3`, `mimo-v2.5`,
`mimo-v2.5-pro`, `minimax-m2.7`, `minimax-m3`, `qwen3.6-plus`, `qwen3.7-max`,
`qwen3.7-plus`, `qwen3.8-max`.

## What Roundfix answered

With a `data` profile added to `.roundfixrc.yml` naming
`runtime: opencode, model: opencode-go/kimi-k3, reasoning_effort: max`:

```
$ roundfix profiles validate --category data --json
"references":[{"category":"data","source":"project","role":"preferred"}]
"error":"capability evidence invalid: contradictory_response, missing_model_state, too_many_option_values"
```

The profile resolved with `source: project`. The failure was capability
evidence, not registration.

The same repository state, through the Doctor Command:

```
$ roundfix doctor
adapter: ok (claude: … | codex: …)
profiles: ok (5 distinct tuples; 10 category references)
```

Both lines are unchanged from before the `data` profile existed. Doctor reports
`ok` while a configured profile is failing, and never names `opencode`.

The probe edit was reverted; `.roundfixrc.yml` is byte-identical to its
committed state.

## The capability payload

`acpx --cwd <dir> --format json --json-strict opencode sessions show <name>`
returns schema `acpx.session.v1`. Its `acpx.config_options` carried three
entries:

| id | category | type | values |
| --- | --- | --- | --- |
| `model` | `model` | `select` | 417 |
| `effort` | `thought_level` | `select` | per-model, see below |
| `mode` | `mode` | `select` | 2 (`build`, `plan`) |

Whole-payload sizes measured: 50,352 bytes on a fresh session, 50,589 bytes
after selecting a model, 50,590 bytes on a second session. Roundfix's
`maxCapabilityResponseBytes` is 65,536, so the payload sits at 77% of the read
limit before any session history accumulates.

## The effort vocabulary is per model

Each row is `sessions ensure --model <M>` followed by `sessions show`:

| model | advertised effort values | current |
| --- | --- | --- |
| `opencode-go/kimi-k3` | `max` | `max` |
| `opencode-go/qwen3.8-max` | `high`, `max` | `high` |
| `opencode-go/glm-5.2` | `high`, `max` | `high` |
| `opencode-go/deepseek-v4-pro` | `high`, `max` | `high` |
| `opencode-go/minimax-m3` | `none`, `thinking` | `none` |
| `openrouter/anthropic/claude-opus-5` | `low`, `medium`, `high`, `xhigh`, `max` | `low` |
| `openrouter/openai/gpt-5.6-luna` | `none`, `low`, `medium`, `high`, `xhigh`, `max` | `none` |

On a session ensured with no `--model`, the current model is
`opencode/big-pickle` and **there is no `effort` option at all** — only `model`
and `mode`.

## The effort cannot be applied before the first prompt

Single scratch directory, one session, commands in the order shown:

| # | command | result |
| --- | --- | --- |
| 1 | `sessions ensure --name rf-p7 --model opencode-go/qwen3.8-max` | created |
| 2 | `sessions show rf-p7` (no `--model`) | `model=opencode-go/qwen3.8-max`, `effort=[high,max]` |
| 3 | `set effort max -s rf-p7` | `-32602 Invalid params: effort not found: max` |
| 4 | `set model opencode-go/qwen3.8-max -s rf-p7` | `{"action":"model_set","modelId":"opencode-go/qwen3.8-max"}` |
| 5 | `set effort high -s rf-p7` | `-32602 Invalid params: effort not found: high` |
| 6 | `set mode build -s rf-p7` | `config_set`; `model` currentValue is `opencode/big-pickle`, no `effort` option present |
| 7 | `prompt -s rf-p7` with `--model opencode-go/qwen3.8-max` | agent answered `opencode-go/qwen3.8-max`; `stopReason: end_turn` |
| 8 | `sessions show rf-p7` | `model=opencode-go/qwen3.8-max`, `effort=[high,max]`, current `high` |
| 9 | `set effort max -s rf-p7` | `config_set`, `configId=effort`, `value=max`, 34,803 bytes |

Step 3 was also run with `--model` on the `set` invocation, in a separate clean
directory, and failed identically.

Steps 2 and 6 disagree about the same session at the same moment. `sessions
show` reports the record acpx persisted from `--model`; `set` speaks to a
transient agent process that acpx starts per invocation, which begins on the
runtime default and therefore advertises no `effort`. The queue owner that holds
the selected model is started by the first prompt — acpx's own status line says
`session idle; queue owner will start on next prompt` — which is why step 9
succeeds where step 3 failed.

The prompt in step 7 reported `inputTokens: 6, outputTokens: 66, totalTokens:
45649, cachedWriteTokens: 45577`.

## What this establishes

1. A large advertised catalog is normal adapter output, not malformed evidence.
2. The `opencode` model selection reaches the agent through the prompt path and
   works: the model answered with its own identifier.
3. A reasoning effort cannot be applied to an `opencode` session before that
   session's first prompt, so a token-free proof cannot apply one.
4. A configured optional-category profile is invisible to the Doctor Command.
