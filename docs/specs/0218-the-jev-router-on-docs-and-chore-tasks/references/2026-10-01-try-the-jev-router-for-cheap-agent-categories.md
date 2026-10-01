---
type: feat
status: promoted
created: 2026-10-01
spec: 0218-the-jev-router-on-docs-and-chore-tasks
reason: null
---

# Try the Jev Router for cheap agent categories

## Opportunity

OpenRouter publishes `typesafe/jev-router` (<https://openrouter.ai/typesafe/jev-router>), an OpenAI-compatible chat router that uses Jev to pick a model and a reasoning effort for each request. Roundfix pays the same strong model and effort for every request of a work category, even where most requests are routine, such as in the `docs` and `chore` categories.

## Value

If the router keeps quality on routine work while it picks cheaper models or lower effort, those categories cost less per Task. The hypothesis is that `docs` and `chore` Tasks run through the router cost less per settled Task without raising the corrective-Task rate.

## Shape

This is intent only, not a commitment.

- Today the router does not fit implementation categories. Roundfix Runs drive ACP agent runtimes (`codex`, `claude`, `opencode`) logged in to subscription accounts, not a chat-completions endpoint billed per token, so the router has no place in those selections.
- A possible experiment: the `opencode` runtime configured with an OpenRouter provider whose model is `typesafe/jev-router`, selected only for the `docs` and `chore` categories, measured against the current defaults on cost per settled Task and corrective-Task rate over a fixed set of Specs.
- The router's spend would go to an OpenRouter key of its own, as the advisory judge's does (`ROUNDFIX_OPENROUTER_API_KEY`, Spec 0205), so the experiment's cost is separable.
- The Node.js dependency of a Run comes from `acpx` and the ACP adapters, not from Jev: the advisory judge calls Jev from Go over HTTPS, and this experiment would add none.

Evidence: OpenRouter's Jev hub (<https://openrouter.ai/docs/guides/community/jev>) lists the Jev Router as a router that uses Jev to pick a model and reasoning effort per request; read 2026-10-01.
