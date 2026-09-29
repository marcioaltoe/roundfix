---
type: fix
status: open
created: 2026-09-29
spec: null
reason: null
---

# The citation projection strips sections from any file named like a Task

## Symptom

Spec 0181's task_06 makes the Spec citation walk skip the Agent-owned `## Result` and the Daemon-owned `## Recorded paths` and `## Carry-forward provenance` sections. It skips them in every file whose basename matches `task_*.md`, including a file under `references/`. An adopted source with such a name would have those sections left out of `SC-ADR-UNLISTED`, although `references/` is authored and should be read in full. No adopted reference in the repository has that name today.

## Where

`readSpecCitations` in `internal/speccheck/citations.go`, the Task file test at about line 1149.

## Expected

Only the Spec folder's own top-level Task files have their Agent- and Daemon-owned sections skipped. Every file under `references/` is read in full.

## Evidence

The second pre-PR review of Spec 0181's candidate `7efe65f2` on 2026-09-29. It was dismissed at the review ceiling with this entry as the follow-up.
