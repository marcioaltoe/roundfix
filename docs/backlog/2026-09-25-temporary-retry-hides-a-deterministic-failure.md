---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# A temporary failure on the exclusive retry hides the command's first-run deterministic failure

## Symptom

In independent Verification mode, when a command fails deterministically in the
first run and then fails temporarily on the exclusive retry, the retry's
temporary outcome replaces the deterministic failure. The Task still settles
failed, but its reason names only a temporary failure, so the operator is led to
re-run instead of fixing an observed defect.

## Where

`internal/daemon/task_engine.go` `retainCollectedVerificationFailures`: the
retry verdict set excludes the unknown command but not the temporary one.

## Expected

Only a passed or failed retry replaces a first-run verdict; a command that ended
temporarily on the retry keeps its first-run deterministic failure and
diagnostic.

## Evidence

Recorded limit of Spec 0167 (second pre-PR review of 2026-09-25), reproduced
with a direct-merge test: first run A deterministic + B temporary, retry A
temporary → only A's temporary failure is reported.
