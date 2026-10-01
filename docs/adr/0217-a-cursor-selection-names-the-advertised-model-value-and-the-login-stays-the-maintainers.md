---
status: accepted
created_at: 2026-10-01T22:00:00Z
updated_at: 2026-10-01T22:00:00Z
deprecated_at: null
superseded_by: null
---

# A Cursor selection names the advertised model value, and the login stays the maintainer's

Cursor's agent CLI speaks ACP through its hidden `acp` subcommand, and acpx
reaches it as its built-in `cursor` agent. Two of its properties do not fit
the rules the other ACP Runtimes follow.

First, Cursor advertises each model as one value whose brackets carry its
parameters, such as `grok-4-20[thinking=true]`,
`gpt-5.4[reasoning=medium,context=272k,fast=false]` or `default[]`, and no
separate reasoning option. Roundfix's capability projection reads a trailing
bracket as a reasoning effort and refuses an empty one, so it refused a whole
Cursor session. Mapping the bracket onto Roundfix's effort vocabulary would
invent a meaning for `thinking`, `context` and `fast` that the adapter does
not state, and composing a value would let Roundfix choose the speed tier, and
with it Fast billing, on the maintainer's behalf. A Cursor selection therefore
names exactly a value the session advertises, brackets included, with an empty
`reasoning_effort`, and the projection reads every Cursor value, `default[]`
included, as one model identity. Proof sets that value and requires the
session to report it current. Other runtimes keep their parsing.

Second, `cursor-agent` authenticates with the maintainer's own Cursor login,
and acting on it needs a person in a browser. On 2026-10-01 the maintainer
decided that a missing login stops the work with a named blocker and that
Roundfix never automates a login or handles a credential. Roundfix therefore
reads only whether a login exists, through `cursor-agent status`, and refuses
a Cursor selection with `cursor_login_required` when it does not, both in the
Doctor Command and before a Run. It never runs `login` or `logout`, never
passes `--api-key` or `--auth-token`, never reads `CURSOR_API_KEY` or
`CURSOR_AUTH_TOKEN`, and never prints the account.

## Consequences

This refines ADR-0037 and ADR-0049 for one runtime: the selection is still
explicit and proved, but its model value carries what other runtimes split
into a model and an effort. `cursor` is opt-in; it enters no built-in or
Recommended Profile (ADR-0180) until a measurement supports it. If Cursor
later advertises a separate reasoning option, a new decision replaces the
empty-effort rule rather than reading both.
