---
type: fix
status: promoted
created: 2026-10-26
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
reason: null
---

# Settle and a tracked mode change commit an executable the Task commit refuses

## Symptom

A Task commit refuses an untracked executable file unless the Spec's
`_authorization.md` lists it under `executables:` (ADR-0157, ADR-0340). Two
paths still commit one: `roundfix settle` stages every change of its surface
with `git add -A` and filters nothing, so settling a Task commits an
executable the Daemon refused for the same Task; and a Task that adds an
execute bit to a file Git already tracks passes the filter, because a tracked
path is work output whatever its mode.

## Where

`internal/cli/settle.go` `settleTaskAndCommit`; `internal/daemon/task_engine.go`
`FilterStageablePaths` (the tracked-path branch).

## Expected

Decide whether settle applies the Task commit's executable rule, including
the `executables:` list, and whether a mode change of a tracked file to
executable needs the same grant; then make both paths agree with the decision.

## Evidence

Found while planning Spec
`0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate`
on 2026-10-26, from Quaestor's report (ai-memory message
`01a17abd-20fe-7213-ab5c-f1a1bdb043fe`, 2026-10-26 21:22 BRT): Quaestor
recovered its refused launchd scripts with `roundfix settle` after making them
`0644`; settle would have committed them at `0755` too.
