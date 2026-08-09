---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Roundfix warms an OpenCode session to apply a reasoning effort

ADR-0106 made the `opencode` runtime model-managed and refused any non-empty
reasoning effort, because the `effort` config option only exists once a
queue-owner agent process holds the selected model and acpx starts that owner on
the session's first prompt — so a token-free Exact Agent Selection Proof cannot
apply one. That reasoning still holds, and its consequence turned out to be
expensive. Measured through OpenRouter on 2026-08-08, three of the four
candidate models deliver their **weakest** advertised setting by default:
`x-ai/grok-4.5` advertises `low, medium, high` and defaults to `low`;
`moonshotai/kimi-k3` and `deepseek/deepseek-v4-flash-0731` advertise
`low, high, max` and default to `low`; only `deepseek/deepseek-v4-pro`, which
advertises `high, xhigh`, opens above the floor. Accepting the default is
therefore not neutral — it silently selects the bottom of the range, and the
published benchmarks for these models describe their top.

Roundfix therefore warms an OpenCode session: after ensuring the session with
its model, it issues one minimal prompt to raise the queue owner, applies the
requested effort, observes the effective value, and only then sends work. The
cost is one round trip, not one prompt's worth of tokens — the system-prompt
cache write happens on whichever prompt comes first, so warming moves that write
rather than adding it.

The proof splits across two moments, and each is honest about what it observed.
Preflight stays token-free and proves what it can: the model is advertised and
current, and the requested effort appears among the values that model
advertises. The Run proves the rest, observing the applied effort before any
work turn. This is a distinct selection encoding — `runtime_deferred` — kept
apart from `independent`, which applies and observes in one token-free step, and
from `runtime_managed`, which is an empty effort on a runtime whose advertised
control Roundfix declines to assign. Collapsing them would let a preflight claim
an assignment it never made.

Leaving ADR-0106 in force was rejected because it spends real capability to
avoid a round trip. Applying the effort only after the session's first real
prompt was rejected because the opening turn — where an Agent reads its context
and plans — is the one that most decides a Batch, and it would be the only turn
running at the floor. This decision supersedes ADR-0106.
