---
type: fix
status: open
created: 2026-10-13
spec: null
reason: null
---

# A QA Run can send Spec content to the judge provider

## Symptom

The QA of 0259 ran `roundfix archive --plan`, which asks the judge for Archive
Advice. Because the QA environment inherited the judge's keys, the command
reached the external judge provider during QA with the Spec's files, without
an authorization for that data to leave the machine.

## Where

The Run environment passes the judge's credentials through to QA commands, and
`roundfix archive --plan` calls `judge.AdviseArchive` whenever keys are present.

## Expected

A Run does not reach an external judge unless its Spec authorizes it: QA
commands run without the judge's keys, or `archive --plan` inside a Run skips
Archive Advice and says so.

## Evidence

0259 QA finding F4 (operator log entry 287, 2026-10-12); the corrective added a
Spec-local requirement that QA run without judge keys.
