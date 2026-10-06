---
status: accepted
created_at: 2026-10-05T00:00:00Z
updated_at: 2026-10-05T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Each OpenRouter stage reads its own Roundfix key first

ADR-0201 gave Roundfix one OpenRouter key of its own,
`ROUNDFIX_OPENROUTER_API_KEY`, for every OpenRouter use, so that Roundfix's
spend stays apart from other projects'. With the Jev judge and implementation
on open models both billed to that key, OpenRouter's activity export, which
groups spend by API key, can no longer tell the two stages apart. On
2026-10-05 the maintainer chose to create two keys, each with its own monthly
limit: "Sim, eu crio".

Each OpenRouter stage now reads its own Roundfix variable first and the shared
key second. The judge reads `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, then
`ROUNDFIX_OPENROUTER_API_KEY`, then `ROUNDFIX_TYPESAFE_API_KEY` for TypeSafe
directly. Implementation on an open model through OpenCode reads
`ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY`, then `ROUNDFIX_OPENROUTER_API_KEY`.
An empty variable counts as unset. The generic `OPENROUTER_API_KEY` is never
read. One list in Roundfix names each stage's variables in that order, and
every reader takes it from there.

A stage records which variable it used, by name and never by value: every
Judge Log line and the `spec judge` report carry `key_variable`, and so does
each implementation spend record. `roundfix doctor` reports each stage's
variables as set or not set, in preference order. Each stage keeps its own
ceiling, whichever key it reads: the judge `jev.monthly_ceiling_usd` and
implementation `openrouter.implement_monthly_ceiling_usd`. A key's own
OpenRouter limit stays the hard stop on that key.

Keeping one key and splitting the cost in Roundfix's own records was
rejected: the maintainer asked for OpenRouter's export to show each stage's
cost, and only a separate key does that. Requiring the stage keys without a
fallback was rejected because it would stop both stages until the maintainer
creates the keys.

## Consequences

- ADR-0201 is superseded in part: its key scope ("reads these two names")
  becomes these per-stage names with the shared key as fallback, and the
  judge without a key names `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` first. Its
  recipients, its boundary and its Judge Log stand.
- With only the shared key set, both stages bill to it as before and the
  export cannot separate them; `key_variable` in Roundfix's own records still
  does.
- A later OpenRouter stage adds its own variable to the same list rather
  than reading the shared key directly.
