---
type: fix
status: open
created: 2026-10-04
spec: null
---

# Queued Specs raise the same owned skill version and conflict at merge

## Problem

On 2026-10-04 four queued Specs each raised the Roundfix Skill version
"from the value on the tree the Task starts from" (0223, 0226, 0225, 0224).
The queue started every Run from the item branch it created at queue start
(9a053bee), not from main as it moved. So 0226 and 0225 both recorded
0.1.26, and 0224 recorded 0.1.25 over main's 0.1.27. Two items parked
(`gate-failed`, `pull-request-conflict`), and the operator resolved each by
merging main and raising to the next patch version by hand (intervention
log entries 157 and 159).

The same Specs hit a related defect. The Codex Task agent for 0225 twice
wrote the version digest by hand instead of running
`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
because it took that command for part of the Verification. The record flag
never replaces a recorded digest, so the stale entry failed every later
attempt (entries 150 and 154).

## Expected

Two queued Specs that raise the same owned skill never record the same
version. Either the queue starts or refreshes an item from current main
before its Run, or the version raise is a derived change the queue
regenerates at merge, like `delivery.derived_paths`. A Task that raises an
owned skill version records it with the record command, and a stale record
for an unreleased version can be replaced.

## Notes

Non-binding. ADR-0189 ties a version to content. A derived-path
regeneration (ADR-0192) of `skills/testdata/owned-skill-versions.json` plus
the two front-matter fields may be the smallest fix.
