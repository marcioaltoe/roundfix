---
type: fix
status: promoted
created: 2026-10-26
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
reason: null
---

# Settle run inside the Run Worktree takes it for the user checkout

## Symptom

`roundfix settle --spec <slug> --task <id>` run with the working directory
inside the kept Run Worktree finds the Run, commits in the Run Worktree, and
then integrates into the Run Worktree itself as if it were the user checkout:
it prints `integration pending — git merge --ff-only roundfix/run-<id>` and
exits `1`, and the Run Worktree is kept. Run from the user checkout, the same
settle integrates, removes the Run Worktree and exits `0`.

## Where

`internal/cli/settle.go` `preflightSettle`: `runworktree.Ref.UserRoot` is the
repository root of the command's working directory, which is the Run Worktree
itself when settle runs there.

## Expected

Settle resolves the user checkout from the Run's record whatever directory it
runs in, or refuses to run inside a Run Worktree and names the user checkout
to run it from.

## Evidence

Measured on 2026-10-26 at `6a2630e5` while planning Spec
`0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate`,
with the built binary on a kept Run created by the `internal/cli` implement
test fixture: from the Run Worktree, `integration pending` and exit `1`; from
the user checkout, `settled task_01 completed` and exit `0`.
