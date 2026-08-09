---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A runtime catalogue is read before the request that asks about it

Roundfix decided whether an Agent Model was advertised by matching the string
`did not advertise that model` in an adapter's stderr. That made the refusal the
adapter's decision, not Roundfix's, and adapters disagree. Measured on
2026-08-09, `codex-acp` refuses `gpt-9-does-not-exist` with
`model_not_advertised`, while `claude-agent-acp` accepts
`opus-9-does-not-exist`, reports it as the session's `currentValue`, and appends
it to the list of advertised options it returns. So the same profile shape passes
on one runtime and fails on another, and `roundfix profiles validate` reported
`passed` for a model that does not exist.

Reading the advertised set from a session that was ensured *with* the requested
model cannot fix this, because that read is the contaminated one. Ensuring the
same session without the override returns the honest catalogue — `default`,
`opus[1m]`, `claude-fable-5[1m]`, `sonnet`, `haiku` — with the invented model
absent.

Roundfix therefore establishes what a runtime offers before it asks about a
specific selection, and decides membership against that catalogue. The
stderr-matched refusal stays as a fast path where an adapter does emit it, but it
is no longer the only thing standing between a maintainer and a profile that
cannot run. The proof remains token-free: this is one additional ACP round trip
on a disposable session that already exists.

Trusting the adapter was rejected because `CONTEXT.md` already defines Exact
Agent Selection Proof as observing matching effective state, and an echo of the
request is not an observation. Sending a prompt to find out was rejected because
it would make preflight cost tokens, which is the property that lets readiness
run before every Run and configuration mutation.

Where an adapter is found to echo a request back into its own advertisement, that
fact is recorded in the proof evidence rather than silently worked around. A
maintainer reading a passing proof should be able to tell whose word it rests on.
