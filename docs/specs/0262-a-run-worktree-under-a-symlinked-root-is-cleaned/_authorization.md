---
status: approved
granted: 2026-10-12
action: resolve a recorded worktree path within its Worktree Root and compare worktree paths through their longest existing prefix, so a Run, Task or item worktree under a symbolic-linked worktree location is reconciled, released and pruned, and document the reclaim with roundfix reconcile --apply
consuming: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
paths: []
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0262

On 2026-10-12 the maintainer approved cycle AF and answered "Sim, como no AE"
to "As regras de operação do AE valem para o AF?": the same operations and the
same "Concedo" for the Governed Paths each Spec declares, without the Makefile
or the CI workflows. For this Spec, item A, the decisions are: a recorded Run
Worktree path and a registered worktree path that resolve to the same
directory are the same worktree; a path resolving outside the worktree root
stays unsafe; every comparison of recorded and registered worktree paths is
swept for the same class; and the command that
reclaims the retained terminal Run Worktrees must work on the symlinked root
and be documented. The admin-lock test source of the same cycle is fixed by
Spec 0259's corrective task_06 and leaves this Spec.

The governed set was measured on the authoring branch at `94fbf85b` with
`GovernedPath`, through a `go test -overlay` probe in `internal/speccheck`
that wrote nothing to the repository, against every file the Tasks declare.
None of them is governed, so the record carries `paths: []`.

## What is not governed

The Go sources and tests under `internal/worktree`, `internal/config`,
`internal/cli` and `internal/daemon`; `CONTEXT.md`; the user guides under
`docs/user-guide/`; `docs/references/skill-coverage.json` and
`docs/references/behavior-surfaces.json`; and
`docs/adr/0268-a-recorded-worktree-path-is-resolved-within-its-worktree-root.md`.

## Sanctioned regeneration

The Behavior Surface Record is re-recorded after the user guides change with
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
That record is not a Governed Path, and this note adds no source path.

No skill, Baseline module or derived Baseline file changes, so no skill or
module version is recorded and `make baseline-digests` rewrites nothing. No
record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, any skill, any Baseline module, or
  the `### QA settlement` section of any skill.
- No rewrite of a recorded Run path and no Run Database migration.
- No test, Verification command or QA row reads or writes the real
  `~/.roundfix`, follows the maintainer's real worktree link, or reaches the
  network.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
