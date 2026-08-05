---
task: task_02
spec: 0068-spec-close-audit
status: completed
type: backend
complexity: high
---

# Task 02: Classify every survivor with its evidence

## Overview

The audit's core: enumerate the branches and worktrees a Spec's cycle produced
and classify each survivor as backing an open Pull Request, holding
unintegrated work, reclaimable residue, or unclassifiable and therefore
preserved. Every classification carries the evidence that produced it.

Verifiable on its own through fixture repositories, one per kind.

## Requirements

1. MUST enumerate the branches and worktrees associated with a Spec, resolving
   its Runs from the Run Database rather than requiring the caller to know them.
2. MUST classify each survivor as `pull-request`, `pending`, `residue`, or
   `preserved`, and MUST attach the evidence string that produced it.
3. MUST classify a worktree with no matching Run as `preserved`, never as
   residue. Scratch worktrees live outside the Run Database, and assuming
   residue there is how work gets deleted.
4. MUST attach the exact reclaim command to a `residue` survivor and MUST NOT
   execute it. The audit reports; the operator reclaims.
5. MUST NOT touch, lock, or reclaim anything, and MUST never classify a branch
   or worktree belonging to an Active Run as residue, per ADR-0052.
6. MUST open no network connection and write no file.

## Subtasks

- [ ] Add the package with its kinds, survivor model, and result.
- [ ] Enumerate branches, worktrees, and the Spec's Runs.
- [ ] Classify each survivor and attach its evidence.
- [ ] Add the Active Run guard before the classifier that could violate it.
- [ ] Add one fixture repository per kind.

## Acceptance Criteria

- [ ] A branch backing an open Pull Request classifies `pull-request` with
      evidence naming the Pull Request.
- [ ] A branch holding unintegrated commits classifies `pending` with evidence
      naming what is unintegrated.
- [ ] A merged branch with no Pull Request classifies `residue` and carries the
      exact reclaim command as a string.
- [ ] A worktree with no matching Run classifies `preserved`, never `residue`.
- [ ] An Active Run's branch and worktree are never classified `residue`,
      proven by an injected active fixture.
- [ ] Every survivor in every fixture carries a non-empty evidence string.
- [ ] The package opens no transport and writes no file, proven by a
      repository-wide check rather than by inspection.

## Context

- interface: `internal/store`
- interface: `internal/gittest`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/specaudit -count=1 -run 'Classif|Survivor|Active|Evidence' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the classification tests ran and passed.
- `go test ./internal/specaudit -count=1` — expected: exit 0.
- `if grep -rn "net/http\|os.Create\|os.WriteFile\|os.RemoveAll" internal/specaudit --include="*.go" | grep -v "_test.go" | grep -q .; then exit 1; fi`
  — expected: exit 0; the audit opens no transport and removes nothing.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 1, 4, 5 and 6; Decisions. Core Feature 4 is
  satisfied here by classifying a pushed-and-merged scratch worktree as
  `residue` with its reclaim command; the reclaiming stays the operator action
  the Decisions require.
- `_techspec.md` → Interfaces; Build Order 2.
- ADR-0052.

## Result

Implemented the read-only `internal/specaudit` package. `Audit` resolves every
Run for the requested repository and Spec through `store.OpenReader`, enumerates
local branches and worktrees with optional Git locks disabled, discovers
Spec-authored branch tips through the `Roundfix-Spec` trailer, and returns a
stable survivor list with one evidence string per classification. Reclaiming is
represented only by shell-safe command strings on `residue` survivors.

Focused checks run during implementation (the declared `## Verification`
commands remain Daemon-owned and were not run):

- Red signal: `rtk env GOCACHE=/private/tmp/roundfix-specaudit-gocache go test
  ./internal/specaudit -run '^TestAuditClassifiesPullRequestBranch$' -count=1`
  failed because the new package models and `Audit` did not exist.
- `rtk env GOCACHE=/private/tmp/roundfix-specaudit-gocache go test
  ./internal/specaudit -run
  '^TestAudit(ClassifiesPullRequestBranch|ClassifiesPendingBranch|ClassifiesResidueBranch|PreservesUnmatchedWorktree|PreservesActiveRunSurvivors)$'
  -count=1 -v` passed all five real Git and Run Database fixtures.
- `rtk env GOCACHE=/private/tmp/roundfix-specaudit-gocache go vet
  ./internal/specaudit` exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-specaudit-gocache go list -f
  '{{join .Imports " "}}' ./internal/specaudit` reported only `bytes context
  errors fmt os/exec path/filepath roundfix/internal/store sort strconv
  strings`; no transport package is imported.
- `rtk rg -n 'net/http|os\.(Create|WriteFile|RemoveAll)'
  internal/specaudit -g '*.go' -g '!*_test.go'` returned no matches (ripgrep's
  no-match exit status is 1).

Acceptance evidence:

1. `TestAuditClassifiesPullRequestBranch` passed and asserted
   `KindPullRequest` with evidence containing `Pull Request #42` from the
   Spec-associated Run.
2. `TestAuditClassifiesPendingBranch` passed and asserted `KindPending` with
   evidence naming the one commit and branch-only file not represented on the
   default branch.
3. `TestAuditClassifiesResidueBranch` passed against a squash-merged branch,
   asserted content-based `KindResidue`, and matched the exact command
   `git branch -d -- 'ma/spec-close-residue'` without executing it.
4. `TestAuditPreservesUnmatchedWorktree` passed and asserted `KindPreserved`
   with `no matching Run` evidence even though the worktree's branch had Pull
   Request evidence; its reclaim string stayed empty.
5. `TestAuditPreservesActiveRunSurvivors` passed for both the injected Active
   Run branch and worktree. Both evidence strings name the Active Run ID,
   neither survivor is `residue`, and neither carries a reclaim command.
6. Every classification fixture calls `assertEverySurvivorHasEvidence`; all
   five fixtures passed with no empty evidence string.
7. The focused direct-import and forbidden-call sweeps above found no network
   transport or file-writing API in production package files. The broader
   repository-wide absence command remains for Daemon Verification.

Follow-ups: none discovered inside this Task's slice.
