---
type: fix # feat | fix | perf | refactor
status: promoted # open | promoted | declined
created: 2026-08-08
spec: 0092-a-run-that-can-hand-back-its-work # Spec slug when status: promoted
reason: null # required when status: declined
---

# A Session that never opened is a selection failure, not a failed Batch

## Opportunity

The Fallback Chain exists so a Run survives a runtime that cannot serve it, and
it did not activate when the Codex quota was exhausted on 2026-08-08. The
repository's own Project Config declares `codex → claude` for the `backend`
category; no fallback attempt appears anywhere in the Run's console. This intent
came through
`inbox/roundfix/2026-08-08-fallback-chain-nao-dispara-em-estouro-de-cota.md` in
the Secondbrain.

## Value

A Fallback Selection may switch ACP Runtime automatically only while Agent work
has not begun. Roundfix emitted `AGENT_WORK_STARTED` and then classified the
adapter's exit as a Batch failure after work started, so the chain became
ineligible. No work had begun: the adapter printed its usage limit and exited,
and the whole Run lasted nineteen seconds. The guard is right in intent — never
swap models over partially finished work — and wrong at the boundary, treating
session opening as work.

The cost is the one the chain was configured to prevent: a Spec stopped for four
days on a quota that a configured, proven alternative could have absorbed
immediately.

A second obstacle showed up in the workaround. A one-Run override cannot express
a Fallback Chain, so when the override tuple matches the category's configured
fallback the preflight refuses with `duplicate Agent Selection ... one
additional distinct authorized and proven Agent Selection is required`, and no
flag supplies the second one. The Spec continued only by choosing an
artificially distinct reasoning effort.

## Shape

A future design could separate "the Session never opened" from "work started and
failed", making the first a selection failure that the chain covers. Recognizing
quota exhaustion as its own adapter exit reason would help on its own: today it
arrives wrapped as `agent/protocol error`, with the cause readable only in
console prose. The override could also replace the Preferred Selection while
preserving the configured chain, which is what its own warning already says it
does. This shape is non-binding.
