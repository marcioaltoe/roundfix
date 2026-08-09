---
status: done
created_at: 2026-08-09
updated_at: 2026-08-09
kind: finding
origin: direct measurement against a live OpenCode installation
runtime: opencode 1.18.15
adapter: acpx 0.13.0
machine: darwin 25.5.0, node v22.23.2
---

# The OpenCode runtime hands back the floor of every range

Read-only measurement of what reasoning effort each candidate model actually
starts at when Roundfix accepts the model-managed default that ADR-0106
mandated. Taken 2026-08-08 in a scratch working directory outside any
repository, before any code changed.

## Method

For each identifier, one disposable session, then read the record:

```bash
acpx --cwd <scratch> --model <M> opencode sessions ensure --name <probe>
acpx --cwd <scratch> --format json --json-strict opencode sessions show <probe>
```

The `acpx.config_options` entry with `id: effort` carries both the advertised
values for that model and its `currentValue`, which is the level Roundfix
inherits when it assigns none.

## Result

| model | advertised effort values | default Roundfix inherits |
| --- | --- | --- |
| `openrouter/x-ai/grok-4.5` | `low`, `medium`, `high` | **`low`** |
| `openrouter/moonshotai/kimi-k3` | `low`, `high`, `max` | **`low`** |
| `openrouter/deepseek/deepseek-v4-flash-0731` | `low`, `high`, `max` | **`low`** |
| `openrouter/deepseek/deepseek-v4-pro` | `high`, `xhigh` | `high` |

Three of four open at the bottom of their own range. The fourth opens above the
floor only because it advertises no floor.

The subscribed `opencode-go/` tier behaves differently again:
`opencode-go/kimi-k3` advertises `max` alone, so accepting the default there
happens to give the top. That coincidence does not generalize — it is the reason
the earlier measurement read as reassuring.

## Why this matters against published numbers

The comparison the maintainer ran on OpenRouter on 2026-08-08 benchmarks these
models at named efforts: the Artificial Analysis variants are `Kimi K3 (max)`,
`Grok 4.5 (high)`, and `DeepSeek V4 Pro (Reasoning, Max Eff…)`. Design Arena
puts Kimi K3 at the 98th–99th percentile and Grok 4.5 at the 89th–95th across
categories. Every one of those figures describes a configuration Roundfix could
not request while ADR-0106 held. Routing to a model on its published score and
then running it at `low` is not a conservative choice; it is a different model
than the one that was chosen.

## What it does not establish

Nothing about whether a higher effort improves Roundfix's own outcomes. That
requires running the same work at two levels and comparing, which no measurement
here does. This finding establishes only the gap between the level a maintainer
selects and the level the runtime delivers.
