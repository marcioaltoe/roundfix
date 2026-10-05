---
status: superseded
created_at: 2026-10-04T00:00:00Z
updated_at: 2026-10-05T00:00:00Z
deprecated_at: null
superseded_by: ADR-0235
---

# The Jev Router gate reads the account credit, and a relay names the routed model

On 2026-10-04 a routed repair prompt passed the Jev Router gate (US$9.69 of
the US$50 ceiling, US$40.31 left on the key) and was refused by OpenRouter
with HTTP 402 because the account behind the key held about US$5.11. Roundfix
reported `agent/protocol error`, and the Judge Log could not say which model
the router had chosen. The maintainer approved cost controls the same day and
had asked earlier for the work to be "implementado e testado sem
limitações". Three decisions follow.

**The gate reads the account credit.** Before each routed prompt, after the
key-limit and ceiling checks of ADR-0218 and ADR-0231, the gate reads
OpenRouter's account credit (total credits less total usage) and takes the
lower of it and the key's remaining limit when the key reports one. Below a
floor the prompt is not sent, with the reason `openrouter_credit_low`. The
floor is the User Config value `jev.router_min_credit_usd`, a finite number
greater than zero, US$15 when unset; Project Config is ignored with the
warning Roundfix gives for a User Config-only setting, for the reason
ADR-0231 gives for the ceiling: the account is the machine's, not a
repository's. An unreadable credit answer refuses the prompt as
`jev_spend_unreadable`, because the credit above the floor cannot be shown.
US$15 is the per-replay stop line of the 2026-10-04 protocol and about twice
the largest routed prompt measured (US$6.81); OpenRouter holds a prompt's
estimated cost against only a fraction of a low balance, so a floor near one
prompt's price would still be refused in flight.

**A credit refusal has a name.** When OpenRouter answers a routed request
with HTTP 402, the prompt's failure carries the reason
`openrouter_credit_refused` and OpenRouter's `limit_source` when it gives
one. The boundary of ADR-0114 and ADR-0218 is unchanged: before Agent work
began the refusal is a failed selection and the Fallback Chain takes the
Task; after it, the Work Item fails with that reason, because switching
models over modified state is what ADR-0050 forbids.

**A relay sees what the router reports.** Neither the ACP stream nor
OpenCode's records carry the model the router chose; only OpenRouter's
responses do, in each response's `model`, `provider` and `id`. For each
routed Agent Session, Roundfix therefore points the inline provider's base
URL at a relay it owns on the loopback interface, under a path token unique
to that session. The relay forwards each request and response unchanged to
OpenRouter's API, refuses any path without a live token, and notes the
response fields above and any 402 refusal. It never logs, stores or rewrites
a header or a body; the key still reaches OpenCode only through the
`{env:ROUNDFIX_OPENROUTER_API_KEY}` placeholder and now crosses Roundfix's
memory on its way to OpenRouter, never its arguments, files or logs. Each
`router-prompt` Judge Log line records the models and providers the relay
noted for that prompt, in first-seen order, and the last response id. The
relay lives while a routed session is open and stops when the last one
closes.

Reading OpenCode's local database, OpenRouter's daily activity report and a
generation lookup by id were rejected: the first is another program's private
storage and does not hold the routed model, the second needs a management key
and reports completed days only, and the third needs the response ids that
only the responses carry.

## Consequences

- ADR-0218's "Roundfix gives OpenCode an inline provider definition" now
  points that definition at the relay; its key placeholder, its project-only
  selection and its ceiling stand.
- Each routed prompt adds one read of OpenRouter before it (the credit read,
  beside the existing key read) and one loopback hop per model request.
- A routed session that outlives the Roundfix process that opened it loses
  its relay, so its next request fails; Roundfix does not resume routed
  sessions across processes.
- OpenRouter documents `GET /api/v1/credits` as needing a management key; on
  2026-10-04 an ordinary key read it. If OpenRouter enforces that, every
  routed prompt is refused as `jev_spend_unreadable` until the decision is
  revisited.
- An older Roundfix binary refuses a User Config that carries
  `jev.router_min_credit_usd` (ADR-0027).

**Retired (2026-10-05).** Superseded by ADR-0235: the gate, the credit floor `jev.router_min_credit_usd`, the named credit refusals and the relay existed only for the Jev Router, which ADR-0235 retires. The floor key is now a deprecated key that loads with a warning (ADR-0027).
