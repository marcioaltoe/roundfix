---
type: fix
status: open
created: 2026-10-13
spec: null
reason: null
---

# A Claude QA turn ends while its work still runs in the background

## Symptom

An Opus QA Run started its load driver as a background shell job and ended its
turn before the job finished. The QA report stayed `pending` and the Run was
Unresolved, although the work it waited for was still running.

## Where

The QA prompt and the Claude runtime profile allow background shell jobs, and
the turn outcome does not check for work the agent left running.

## Expected

A QA turn either runs its commands in the foreground or waits for every
background job before it ends; a turn that ends with its own jobs still running
is reported as unfinished with the job named, not settled as a pending report.

## Evidence

Operator log entry 302 (2026-10-13 05:30 BRT), Delivery Queue item 0265.
