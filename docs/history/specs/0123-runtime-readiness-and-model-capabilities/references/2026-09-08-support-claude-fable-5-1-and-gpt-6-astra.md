---
type: feat
status: promoted
created: 2026-09-08
spec: 0123-runtime-readiness-and-model-capabilities
reason: null
---

# Support Claude Fable 5.1 and GPT-6 Astra as proven Agent Selections

## Opportunity

The maintainer requested support for Claude Fable 5.1 and Codex GPT-6 Astra on
2026-09-08. Make these models available for deliberate selection once the
configured ACP adapters advertise and successfully prove their identifiers and
controls.

## Value

Model choice can use the refreshed version and cost reference alongside actual
runtime capability evidence. A public API identifier alone does not establish
that the same string works through Claude Code, Codex, or OpenCode's adapter.

## Shape

Provisional group P4, runtime readiness and model capabilities, owns the support
work. The session's Exa research checked primary vendor documentation and
confirmed API IDs `gpt-6-astra` and `claude-fable-5-1`. OpenRouter uses
`openai/gpt-6-astra` and `anthropic/claude-fable-5.1`; those provider-qualified
identifiers are a different interface from ACP model choices.

Discover the configured runtime's advertised model ID, reasoning controls,
capability bounds, and selection result before publishing usable choices.
Update the catalog and its reference with the proved mapping. Keep effective
profiles unchanged until a separate explicit selection authorizes a switch;
this Backlog Entry does not perform that switch.

Research sources: [OpenAI GPT-6 Astra model documentation](https://developers.openai.com/api/docs/models/gpt-6-astra),
[Anthropic model overview](https://docs.anthropic.com/en/docs/about-claude/models/overview),
and [Anthropic API release notes](https://docs.anthropic.com/en/release-notes/api).
See the project [model selection reference](../../../references/model-selection.md)
for the dated price and version evidence. No model price or ACP availability
is inferred in this entry.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this remaining intent for the implementation queue.
[0123-runtime-readiness-and-model-capabilities](../_prd.md) is its primary owner.
The source moves once into that Spec's reference index; other Specs link this
owned copy. Its lifecycle status records adoption, not implementation or QA
completion. Governed mutation and execution still require the consuming grant.
