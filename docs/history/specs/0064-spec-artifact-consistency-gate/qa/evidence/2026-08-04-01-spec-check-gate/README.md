# Evidence — 2026-08-04 Spec Consistency Check gate rerun 01

This file indexes concise fresh evidence captured against source commit
`45f6eb6f9aa96bdeeed422829546c7f6693a0730`. The QA report contains the exact
matrix and verdict; entries below record command-level observations without
copying full test logs.

## E-01 — Graph, constraints, and tooling authority

- Graph: `_tasks.md` names terminal `task_09`; every transitive dependency is
  `completed`, and task 09 is `pending` under Daemon ownership.
- Constraints: `_prd.md:25-51` and `_techspec.md:33-68` carry all four rows,
  operative `docs/agents/` sources, and the same governing ADR-0089 budget
  obligation.
- Authorization chronology: `git merge-base --is-ancestor` exited 0 for
  `2e560ce → 02c35e2`, `c4e7613 → cab533f`, and `cab533f → 121fabd`.
- Tooling paths: `git diff-tree --no-commit-id --name-only -r` shows only
  `Makefile` and the assigned Task file in task 08 and task 11.
- Pre-flow status: only this untracked report and evidence directory;
  `git diff --check` exited 0.

## E-02–E-14 — Public flows blocked by F-002

The QA contract stops flow testing after a code-caused static failure. These
rows did not run against the existing `bin/roundfix`, because `make verify`
failed before its build stage and that artifact was not rebuilt from the
current source commit.

## E-15 — Mandatory gate and deterministic helper reproduction

- `rtk make verify` — exit 2. The test stage reported 3197 passed, 1 failed,
  and 3 skipped; `TestOwnerProcessControllerForceKillExitProof` failed because
  the readiness pipe was already closed. No later gate stage ran.
- Direct helper invocation — exit 2 after printing `ready`, with
  `fatal error: all goroutines are asleep - deadlock!` at
  `internal/store/process_unix_test.go:473`.
- Focused parent/helper pair, `-count=50 -parallel=16` — exit 0, but every
  child printed the same fatal deadlock. This proves the focused success can
  race past a helper that exits before the controller terminates it.
- `git diff origin/main..HEAD -- internal/store` — empty. The defect predates
  this Spec branch, but it is code-caused and blocks the current repository
  gate rather than qualifying as an environment constraint.
