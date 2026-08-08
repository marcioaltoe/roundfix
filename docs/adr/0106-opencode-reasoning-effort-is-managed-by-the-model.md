---
status: accepted
created_at: 2026-08-08T00:00:00Z
updated_at: 2026-08-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# OpenCode reasoning effort is managed by the model

Roundfix applies a non-empty reasoning effort as a separate ACP config option
after ensuring the Agent Session, and on the `opencode` runtime that call cannot
succeed before the Run's first prompt. Measured on 2026-08-08 against opencode
1.18.15 and acpx 0.13.0: `sessions ensure --model opencode-go/qwen3.8-max`
followed by `sessions show` reports the requested model with `effort` values
`[high, max]`, but `set effort max` answers ACP `-32602 Invalid params: effort
not found: max`, and a following `set mode build` reveals the live agent sitting
on the default `opencode/big-pickle` with no `effort` option at all. Passing
`--model` on the `set` invocation does not change it, and neither does a
preceding `set model`. The `effort` option is per-model and only exists once a
queue owner process holds the selected model, which acpx starts on the first
prompt — the same sequence run after one prompt returns a clean `config_set`.
Roundfix therefore treats `opencode` as a model-managed reasoning runtime and
refuses a non-empty `reasoning_effort` for it in configuration, naming the empty
value as the repair. "Model-managed" here means Roundfix declines to assign a
reasoning effort, not that the adapter offers no control: OpenCode does
advertise a per-model `effort` option, and an Exact Agent Selection Proof over
such a runtime must accept that option's presence rather than require its
absence, because the value it carries is the Agent Model's own and Roundfix
never assigned it. That distinction is a separate selection encoding, so the
stricter rule — an empty effort proves only against an adapter advertising no
reasoning control at all — continues to hold for every runtime Roundfix does
control. Proving the effort by observing the advertised list without
applying it was rejected because `CONTEXT.md` defines Exact Agent Selection
Proof as applying the exact model and reasoning assignment and observing
matching effective state; issuing a prompt to raise the queue owner first was
rejected because the same definition requires the proof to be token-free.
