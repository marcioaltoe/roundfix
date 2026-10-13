---
type: fix
status: open
created: 2026-10-13
spec: null
reason: null
---

# QA evidence records environment values, secrets included

## Symptom

The QA Runs of Spec 0259 ran `make -p` and committed its output as
`qa/evidence/*/make-database-stdout.txt`. The make database prints every
inherited environment variable with its value, so the evidence of three Runs
held a third-party API key. GitHub push protection refused the push; the key
never left the machine, and the evidence folder was dropped from the delivery.

## Where

The QA protocol (`qa-execution` skill and the QA prompt) lets a QA Run persist
any command output as evidence, and nothing between the QA commit and the push
scans it.

## Expected

QA evidence records environment variable names, never their values, and a
command whose output carries the environment is not persisted. Before a QA
report or a Delivery Queue item is committed or pushed, evidence is scanned for
secret-shaped values and the step refuses with the file and line.

## Evidence

Delivery Queue item `0259-proven-signals-and-a-shared-build-cache` on
2026-10-13; the push of PR #512 was blocked by push protection, and the branch
was rebuilt without `qa/evidence`.
