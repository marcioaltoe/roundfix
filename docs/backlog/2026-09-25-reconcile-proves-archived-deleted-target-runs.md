---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# Reconcile can never release a Run whose squash-merged target branch was deleted and whose Spec was archived

## Symptom

For a deleted target branch, reconcile compares the Run head's whole tree with the current `main`. Archiving moves the Spec folder the Run wrote, and `main` keeps moving, so the proof always fails (`11 Run-only files, 35 differing shared files`), yet the same check requires the Spec to be archived. Such Runs are retained forever.

## Where

`internal/worktree/worktree.go` `inspectDeletedTargetRunByContent` / `compareRunContentToDefault`.

## Expected

Prove against the merge commit delivery recorded, or restrict the comparison to the paths the Run changed with the archive rename mapped; keep negative fixtures for unrepresented work.

## Evidence

secondbrain `inbox/roundfix/2026-09-10-reconcile-nao-classifica-branch-cujo-pr-foi-squash.md`; live repro `run_20260924T223414Z_54e6f85ae98b6550` (Spec 0164, merged in #246).
