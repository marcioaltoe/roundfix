---
status: accepted
created_at: 2026-10-06T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A blocked review names the protocol step and retries once before the prompt

On 2026-10-02 a repository reviewing through `codex / gpt-5.6-luna / max`
blocked twice with `review runtime failure: Agent Selection failed for runtime
"codex": agent/protocol error` and an empty answer file, while the Doctor
Command and the Codex CLI answered. Three days later the same review ran. The
reason said neither where the protocol failed nor whether repeating was safe,
and the maintainer skipped the review.

A measurement on 2026-10-06 with acpx 0.19.4 and a fake ACP adapter showed
why the reason is empty. Each `acpx prompt` process repeats the protocol:
`initialize`, then `session/load` or `session/new`, then `session/set_model`
when a model is passed, then `session/prompt`. In JSON mode acpx echoes every
request on stdout, and a failure at any of those steps arrives as a JSON-RPC
error line on stdout carrying the adapter's own message and the failing
request's id, with an empty stderr and exit code `1`. The ACPX Runner parsed
that line and then discarded it in favor of the exit code, so the three steps
produced the same byte-identical reason as the incident.

Now the runner keeps the first JSON-RPC error of a prompt process that exited
without a parsed result. It names the step by the method of the request the
error answers, keeps only the error's `message`, bounded to one line of at
most 512 bytes, and records whether a `session/prompt` request had been sent.
The review command names that step and message in its reason and record.

A failure during Agent Selection, before the review prompt was sent, is
retried once on the same selection, and the record lists the retry. That
covers a failed session preparation and a prompt process that failed before
it sent `session/prompt`. A failure after the prompt was sent, or one the
runner cannot place before it, is not retried, so no candidate is reviewed twice. A
configured fallback still activates only before the prompt and only after the
retry also failed. The policy never falls back to `none`.

Retrying on any failure was rejected because a failure after the prompt may
already have spent a review. Trying the fallback before the retry was rejected
because the incident was transient on the preferred reviewer the repository
chose. Matching adapter error texts to decide what is transient was rejected:
the adapter's wording is not a contract, and the position in the protocol is.

## Consequences

An operator reads which step failed and what the adapter said, and a transient
failure before the prompt no longer needs a manual rerun. A review whose
preferred selection fails twice before the prompt takes one more Agent Session
preparation before its fallback. Runs are unchanged: their selection failures
before Agent work keep activating the next fallback directly (ADR-0050,
ADR-0114). The detection rests on acpx echoing outbound requests in JSON mode;
if a later acpx stops echoing them, no `session/prompt` is seen and a failure
that names no step is not retried.
