---
status: superseded
created_at: 2026-10-01T22:30:00Z
updated_at: 2026-10-05T00:00:00Z
deprecated_at: null
superseded_by: ADR-0235
---

# The Jev Router is a project-selected OpenCode model under the shared Jev ceiling

OpenRouter publishes `typesafe/jev-router`, an OpenAI-compatible model that
uses Jev to pick a model and a reasoning effort for each request, billed to
the caller's OpenRouter account at a price its model listing gives as
variable (`-1`), so what a routed prompt costs is known only after it ran.
Roundfix drives
ACP Runtimes, not chat endpoints, so the router can only be reached through a
runtime that speaks to OpenAI-compatible providers: OpenCode. On 2026-10-01
the maintainer allowed agent prompts that carry this repository's code and
diffs to go to OpenRouter, never adopters' or the Secondbrain's content, on
Roundfix's own key `ROUNDFIX_OPENROUTER_API_KEY` and never the generic
`OPENROUTER_API_KEY`, under the US$5 monthly ceiling that ADR-0201 sets for
the judge and the maintainer extended to all Jev use, and only as an
experiment that changes no default.

Roundfix therefore reaches the router as exactly one OpenCode Agent
Selection, `roundfix-openrouter/typesafe/jev-router` with an empty reasoning
effort, because the router chooses the effort. For that selection's sessions
only, Roundfix gives OpenCode an inline provider definition whose key is the
placeholder `{env:ROUNDFIX_OPENROUTER_API_KEY}`, so the key travels from the
environment to OpenCode and never through Roundfix's arguments, files or
logs. Roundfix checks only that the variable is set.

Only Project Config may select the router. A repository opts in by a
committed decision that names it; User Config, which spans every repository
on the machine, and a one-Run override cannot, so the maintainer's key never
follows a configuration into a repository that did not choose it.

Every routed prompt is gated by the shared ceiling. Before the prompt,
Roundfix computes the month's Jev spend as the Judge Log's TypeSafe cost plus
the larger of the Judge Log's OpenRouter cost and the key's `usage_monthly`
that OpenRouter reports. At or above the ceiling, or when either source
cannot be read, the prompt is not sent. Before Agent work has begun the
refusal is a failed selection and the Fallback Chain moves on (ADR-0050,
ADR-0114); after work has begun the Work Item fails, because switching models
over modified state is what ADR-0050 forbids. After each routed prompt,
Roundfix appends the change in the key's usage to the Judge Log as a
`router-prompt` line in its own schema, so the judge's ceiling counts the
router's spend and the router's gate counts the judge's.

## Consequences

The router enters no built-in or Recommended Profile (ADR-0180); whether it
should is a decision for after the measurement. The ceiling can be overrun by
at most the prompt in flight, and by any lag in OpenRouter's usage report. A
repository other than this one that opts in spends the maintainer's key only
because its own committed Project Config says so.

**Ceiling value (2026-10-04).** Superseded in part by ADR-0231: the shared ceiling, and the key limit it bounds, is the User Config value `jev.monthly_ceiling_usd`, US$5 when unset. Every other part of this decision stands.

**Relay and credit (2026-10-04).** Superseded in part by ADR-0234: the inline provider now points at the loopback relay, and the gate also refuses below the account credit floor. Every other part of this decision stands.

**Retired (2026-10-05).** Superseded by ADR-0235: the maintainer retired the Jev Router because its routed sessions billed OpenAI and Anthropic models through OpenRouter, which the subscription rule forbids. No part of this decision stands; the shared ceiling it relied on stays the judge's, as ADR-0231 decides.
