---
type: fix
status: promoted
created: 2026-10-13
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
reason: null
---

# A model-at-capacity reply is accepted as a finished turn

## Symptom

codex-acp streamed "Selected model is at capacity. Please try a different
model." with `threadStatus` `systemError` and ended with `stopReason`
`end_turn` and null usage. Roundfix accepted the turn: the QA of 0265 ended
twice with a skeleton report and verdict pending, the Run was Unresolved, and
no fallback in the runtime chain fired.

## Where

The ACP turn outcome classification treats `end_turn` as success without
reading the system error the stream carried.

## Expected

A turn whose stream reports a provider capacity or system error is classified
as a runtime failure, so the configured fallback chain runs instead of the Run
settling with an empty result.

## Evidence

Runs `run_20261013T063154Z_69b987082422f8c7` and
`run_20261013T065557Z_05d88ac973ef5348`; the stream is
`~/.acpx/sessions/081b6e5c-bb31-4259-9941-598b885246eb.stream.ndjson`.
