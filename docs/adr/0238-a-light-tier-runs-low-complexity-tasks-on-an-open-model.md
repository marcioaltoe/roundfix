---
status: accepted
created_at: 2026-10-05T00:00:00Z
updated_at: 2026-10-05T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A light tier runs low-complexity Tasks on an open model

ADR-0235 retired the Jev Router and reserved OpenAI and Anthropic models for
the codex and claude subscriptions, and left choosing a cheaper model per Task
as a separate intent. An operator measurement on 2026-10-05 replayed six
archived `complexity: low` Tasks on `opencode` with OpenRouter's
`deepseek/deepseek-v4.1-flash`: all six passed Verification on the first
attempt, DeepSeek's agent time was 0.41 times Codex's, and the six cost
US$0.32 together. The same measurement asked the Jev judge for a tier on 50
archived Tasks and found that its `light` choices passed on the first attempt
as often as any `low` Task (about 82 % against 83 %). On 2026-10-05 the
maintainer chose to build the tier, on by default for every project ("Ligado
para todos"), capped by its own monthly ceiling, falling back to the default
selection when no key is present ("Cai para o padrão").

Roundfix therefore gives each Task a tier when it dispatches it. A Task is
`light` when its `complexity` is `low`, its `type` is not `qa`, and it
declares no Governed Path to create, change or delete; every other Task is
`standard` and runs exactly as before. The tier is read from the Task file,
not chosen by the judge: the measurement found the judge no better than the
authored `complexity`, so `roundfix spec judge --stage tasks` may suggest a
tier while the Tasks are written, and that suggestion stays advisory and is
never read by dispatch.

A `light` Task runs on an Agent Selection Profile derived at dispatch and
recorded with the profile source `light-tier`: its Preferred Selection is the
first light model, and its Fallback Chain is the other light models followed by
the category's own Preferred Selection and Fallback Chain, so the Run Database
and its Run Events need no new selection role. A light model is `opencode`
with an open OpenRouter model from a User Config allow-list, `openrouter.light_models`, whose default is
`deepseek/deepseek-v4.1-flash` and whose entries are OpenRouter model ids that
Roundfix prefixes with OpenCode's `openrouter/` selection namespace. Every
entry must pass ADR-0235's rule, and an empty list turns the tier off on that
machine. The selection carries no reasoning effort. With an effort, ADR-0108
sends a `Session setup.` warm-up prompt before applying it, and the
measurement saw DeepSeek do a whole Task in that prompt at the lowest effort:
acpx 0.19.4 forwards its tool restriction only to the Claude Code and Qoder
adapters, so OpenCode receives no restriction. An empty effort sends no
warm-up, so the light tier does not depend on that restriction.

The light session reads its OpenRouter key from the one variable a single
helper names, `ROUNDFIX_OPENROUTER_API_KEY` today. Roundfix passes OpenCode an
inline provider option that references that variable by name, never its value,
and removes the generic `OPENROUTER_API_KEY` from the session's environment, so
the spend lands on Roundfix's key. With no key, the Task runs on its
category's profile unchanged and the Run prints a warning. A light session that
fails to start falls to the next candidate, and finally to the category's
Preferred Selection, as any selection failure before Agent work does under
ADR-0114, so the allow-list is not proved by Preflight.

Spend is capped by `openrouter.implement_monthly_ceiling_usd`, a User Config
value with a US$10 default, separate from the judge's
`jev.monthly_ceiling_usd`. Cost is the figure OpenCode reports for the
session, which the measurement found within about 2 % of the key's own usage;
Roundfix appends each increase to a Light Spend Log in Roundfix Home, one file
per UTC month, summed across every repository. A month that has reached the
ceiling, or whose log cannot be read, runs light Tasks on their category's
profile unchanged, with a warning. Reading the key's usage from OpenRouter was rejected
because it adds a network call and a second credential use for a figure the
session already reports.

A `light` Task whose first Verification fails escalates once: its one
Verification Feedback turn runs on the category's Preferred Selection, or its
Fallback Chain if that fails to start, in a new Agent Session, which receives the Task prompt, the diagnostics, and the
statement that the working tree holds another model's attempt. The light
attempt's changes are kept because Verification has just described their
state, and the Task still commits only after Verification passes. This is not
a Fallback Chain step: ADR-0050's rule that a fallback never switches models
after Agent work begins stays as it is. Discarding the attempt and starting
again was rejected because Roundfix has no safe partial restore of a Task's
worktree and the attempt is usually close.

## Consequences

- Every repository's `low` Tasks may leave the machine for OpenRouter and the
  open model's provider: the Task prompt, the files the agent reads and the
  Verification diagnostics. The key value never reaches a file, a log or a
  Run Event.
- A one-Run `--agent` override turns the tier off for that Run, because the
  operator named the selection.
- The ceiling is checked before each light Task, so Tasks already running can
  pass it by at most their own cost. OpenRouter's monthly limit on the key
  stays the operator's backstop, and a session whose cost OpenCode does not
  report records zero with an `unreported` source.
- ADR-0108's warm-up remains for OpenCode selections with an effort; that its
  inert prompt is not inert under acpx 0.19.4 is left to a later decision.
