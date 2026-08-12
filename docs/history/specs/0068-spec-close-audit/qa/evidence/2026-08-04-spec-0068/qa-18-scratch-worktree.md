# QA-18 — Pushed-and-merged Supervisor scratch worktree

Status: fail. Finding: F-001 (Blocks-Completion).

Public reproduction on build
`1346d83d4213e10b73a89bae6796d6d95dda6c18`:

1. Initialized a disposable `main` repository and bare `origin`.
2. Created worktree `/private/tmp/roundfix-qa0068-scratch.D31PJZ-wt` on
   `ma/scratch-close`.
3. Committed `scratch.txt` with `Roundfix-Spec: 0068-spec-close-audit`, pushed
   the scratch branch, squash-merged its content to `main`, and pushed `main`.
4. Kept the scratch worktree registered, matching the PRD's surviving
   Supervisor scratch-worktree state.
5. Ran the built command in text and JSON:

   `bin/roundfix spec audit 0068-spec-close-audit`

Observed exit 1 and these survivors:

```text
- residue branch ma/scratch-close
  reclaim: git branch -d -- 'ma/scratch-close'
- residue branch origin/ma/scratch-close
  reclaim: git push --delete 'origin' 'ma/scratch-close'
- preserved worktree /private/tmp/roundfix-qa0068-scratch.D31PJZ-wt
  evidence: worktree ... has no matching Run in the Run Database
```

Expected: PRD Core Feature 4 and Task 02's coverage statement require the
pushed-and-merged scratch worktree to classify `residue` with an exact
`git worktree remove --` operator command. The audit still must execute
nothing.

Independent confirmation: JSON reported the worktree with
`"kind":"preserved"` and no `reclaim`. A fresh public
`git worktree list --porcelain` still showed the worktree on
`ma/scratch-close`; `git branch --list` showed the local branch as checked out
there. Repeating the audit left both surfaces intact.

The audit's local branch reclaim command cannot complete while that branch is
checked out in the preserved worktree. The operator receives no supported
command for removing the worktree first, so the product's close-cleanup
journey cannot finish from its own report.
