---
type: fix
status: promoted
created: 2026-11-12
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
reason: null
---

# A retry after resolving a pull request conflict finds no Specs Root and no pushed head

## Symptom

Reported by quaestor in ai-memory message `01a1d231-03da-79f1-9834-51487a593986` of 2026-11-12 (20:56 BRT), with Roundfix 0.141.1 and 0.142.0. Its item `0101-aviso-do-baseline-e-referencias-sem-uso` archived its Spec and opened pull request #318 while `main` received an amendment of the same Spec (#317), so the pull request had a modify/delete conflict. The Queue Owner stopped with `delivery-error: resolve pull request conflict: git restore --source=HEAD --staged --worktree -- .: error: path 'docs/specs/<slug>/_authorization.md' is unmerged`. The operator resolved it by hand in the item worktree (`git merge origin/main`, `git rm -r -f docs/specs/<slug>` to keep the archive, commit `1a87865`). Two defects followed:

1. The `git rm -r` removed the last file under `docs/specs/`, and Git keeps no empty directory. The retry stopped at the gate with `delivery-error: read delivery authorization: ... resolved to "<worktree>/docs/specs", which does not exist`. The operator's workaround was `mkdir -p <worktree>/docs/specs`.
2. Two consecutive `deliver retry` calls parked the item as `review-stale`. The review record (`ceiling-closed`, head `1a87865`) was right, but pull request #318 still had the old head `1699788`: the retry does not push the local resolution commit before comparing. Neither the status nor the deliver reference says what to do; `review-stale` appears only in the class table. The operator's workaround was a manual `git push origin HEAD:<item branch>`, after which the next retry reached checking and merged.

The same message reports the Queue Owner panic at the gate again with 0.141.1, holding only its own `waiting` entry; that panic is Spec 0385's subject and is not part of this entry.

## Where

`roundfix deliver retry` after a pull request conflict that the operator resolves in the item worktree: reading the delivery authorization of an archived Spec, and the `review-stale` park.

## Expected

Reading the delivery authorization of an already archived Spec does not require the Specs Root to exist (it reads the Archive Record, or treats an absent root as empty). A retry whose item head is ahead of the pull request's head pushes the item branch before comparing, or the `review-stale` park says `push the item branch, then retry`.

## Evidence

none yet — reported by the sender, not reproduced here. Sender's evidence: Quaestor on the mini, 2026-11-12 ~20:00–21:00 BRT; pull request #318 of gesttione-solutions/quaestor; console logs `~/.roundfix/artifacts/85a8e4c6a0c58989/delivery/delivery-85a8e4c6a0c5-{12366,38057,17157,19634,21360}/console.log`.
