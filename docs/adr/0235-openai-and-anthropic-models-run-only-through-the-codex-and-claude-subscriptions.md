---
status: accepted
created_at: 2026-10-05T00:00:00Z
updated_at: 2026-10-05T00:00:00Z
deprecated_at: null
superseded_by: null
---

# OpenAI and Anthropic models run only through the codex and claude subscriptions

ADR-0218 let Project Config select the Jev Router, an OpenRouter model that
picks a model for each request, through the OpenCode runtime. OpenRouter's
activity export for Roundfix's key from 2026-10-01 to 2026-10-05 shows where
that went: routed sessions were billed for `openai/gpt-6-astra` (US$15.35),
`anthropic/claude-opus-5.5` (US$10.12) and `openai/gpt-6.1-sol` (US$4.81),
pay-per-use, for models this repository already reaches through the codex and
claude subscriptions, against US$0.014 for 456 direct `typesafe/jev-1.13`
judge calls. On 2026-10-05 the maintainer set the rule "os modelos da openai
e anthropic devem ser utilizado exclusivamente pela assinatura que o codex e
claude fornecem", decided "Aposentar o router", and, after weighing a
restricted router, confirmed "Vamos seguir com a remoção".

Roundfix therefore retires the Jev Router entirely: its OpenCode selection,
the inline `roundfix-openrouter` provider, the loopback relay, the spend and
credit gate, the credit floor `jev.router_min_credit_usd` and the
`router-prompt` Judge Log lines. Configuration validation refuses an
`opencode` or `opencode-custom` Agent Selection whose model is under the
retired `roundfix-openrouter` provider, or under OpenCode's `openrouter`
provider with an OpenRouter model whose author is `openai` or `anthropic`
(a leading `~` alias marker ignored), or whose author is `openrouter` or
`typesafe` or a saved preset (`@` prefix), because those name routers that
pick a model, OpenAI and Anthropic included, per request. The refusal ends
with the rule, "OpenAI and Anthropic models run only through the codex and
claude subscriptions". It holds in every scope that names a selection,
Project Config, User Config, a one-Run override and `roundfix profiles
configure`, and the ACPX runner applies the same check before it prepares a
session, so a selection that never passed configuration validation is still
refused before any prompt is sent.

Other OpenRouter models stay selectable through OpenCode: the rule is about
OpenAI and Anthropic, and the repository's model reference measures other
vendors there. Refusing the whole provider was rejected as wider than the
rule. Keeping the router with an exclusion list was weighed and declined by
the maintainer's confirmation of the removal.

## Consequences

- ADR-0218 and ADR-0234 are superseded: everything they decide exists only
  for the router. ADR-0231 is superseded in its router parts only; the
  ceiling `jev.monthly_ceiling_usd` stays a User Config value that the judge
  reads.
- The direct judge is unchanged: `roundfix spec judge`, its keys, its
  recipients, its ceiling and its Judge Log lines. `router-prompt` lines
  written before the retirement stay in their month's file and still count
  toward that month's ceiling.
- `jev.router_min_credit_usd` becomes a deprecated key: a configuration that
  still carries it loads with one warning and the value is ignored
  (ADR-0027).
- A third-party router under another author, or an OpenCode provider the
  operator defines outside Roundfix that points at OpenRouter, is not
  recognized; OpenRouter's account-level model restrictions are the
  operator's backstop for those.
- OpenCode's own `openai`, `anthropic`, `opencode` and `opencode-go`
  providers and the Cursor runtime can also serve OpenAI or Anthropic models.
  Some of them bill through a subscription and the measured spend was
  OpenRouter's, so they are left to a later decision.
- Choosing a cheaper or stronger model per Task, which the router was
  measured for, returns as a separate intent:
  `docs/backlog/2026-10-05-a-judge-assigned-model-tier-per-task.md`.
