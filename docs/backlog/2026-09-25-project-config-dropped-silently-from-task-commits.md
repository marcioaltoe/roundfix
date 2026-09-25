---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# The Daemon silently leaves `.roundfixrc.yml` out of a Task commit

## Symptom

A Task authorized to change Project Config settles completed, but its commit omits `.roundfixrc.yml` with no event or reason; only the QA gate notices later.

## Where

`internal/daemon/engine.go` `diffSnapshots` skips `.roundfixrc.yml` unconditionally.

## Expected

Report the exclusion as a dropped stage path so the Task fails legibly; optionally stage it when the frozen authorization bounds the file.

## Evidence

secondbrain `inbox/roundfix/2026-09-24-daemon-silently-drops-project-config-from-task-commits.md`; Spec 0155 QA F-001.
