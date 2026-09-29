---
spec: 0175-cleanup-after-a-squash-merge
status: archived
created: 2026-09-28
surfaces: [backend, cli, docs]
archived: "2026-09-28"
source_slug: 0175-cleanup-after-a-squash-merge
---


# Cleanup after a squash merge

When a Spec's pull request is squash-merged, the Run Worktrees, Run Branches
and staging worktrees its Runs left behind stay in the repository. The
maintainer asked for every finished squash merge to be followed by a cleanup
of completed Runs and branches. On 2026-09-25 this repository held 16
`roundfix/run-*` branches, 16 Run Worktrees and a locked carry-forward staging
worktree, most of them for Specs already merged and archived. Roundfix
released almost none of them. Six defects explain why:

- `roundfix deliver` records the merge, removes the item worktree and item
  branch, and stops. Nothing releases the Spec's Runs after a merge.
- A Run whose target branch was deleted is compared, whole tree, with the
  current default branch. Archiving moves the Spec folder and `main` keeps
  moving, so the proof cannot pass, yet the same check requires the Spec to be
  archived. After Spec 0172 merged (#259), `roundfix reconcile
  run_20260925T200712Z_0929550d586b4ba2` kept that Run as `unintegrated` ("48
  Run-only files, 42 differing shared files against default branch main";
  "default branch has no superseding QA Report after target branch
  disappeared"). Its Task 01–04 commits were ancestors of the merged item
  branch. Its only unique commits were a Task 05 redone after an authorization
  change and a failed QA Report. The operator removed it by hand.
- `--apply` never acts on a Run whose Tasks a later Run delivered; only
  `--discard-superseded` does, so merged Specs leave those Runs behind.
- Runs created from a since-removed linked worktree have no repository key:
  `runs list`, `reconcile` and the implement-preflight prune never see them,
  although their Run Worktrees are registered in this repository.
- A carry-forward killed during `git worktree add` leaves a `locked
  initializing` staging worktree that one `--force` cannot remove and nothing
  sweeps.
- Cleaning a half-removed item directory deletes it whole, even when another
  registered worktree lives inside it.

This Spec makes Roundfix prove each of those surfaces integrated or superseded
and release it, with no manual `git worktree remove`.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; Run IDs, Run
  Branch names, item branch names and the `roundfix-reconcile/v1` schema name
  stay as they are. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git, local files and the Run
  Database only. The merged pull request head comes from the Delivery Queue
  item the delivery owner already recorded; no credential and no network call
  is added. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0053 makes reconciliation explicit
  and proof-based and read-only without `--apply`; the new proof is one more
  proof, and dry-run stays read-only. ADR-0115 keeps the Branch Disposition a
  separate, named act with its evidence written first, so `--apply` does not
  absorb `--discard-superseded`, and the Run it names is released only when the
  merged-head proof covers it. ADR-0023 gives each spec Run its own worktree
  and branch, which is what gets released. ADR-0052 keeps the Integration
  Pending promotion compare-and-set. ADR-0004 keeps the Run Database the one
  record of Runs, and the repository key is backfilled there. ADR-0014 and
  ADR-0057 keep Task status Daemon-owned; the proof only reads committed Task
  files. ADR-0044 reclaims a resource only on proven owner death, which the
  staging sweep follows. ADR-0158 keeps a BudgetExceeded Run recoverable, which
  the proof respects by releasing it only when its work is represented. ADR-0163
  already accepts a recorded squash merge with its merge commit as a receipt
  for Review Artifacts; this Spec applies the same principle to Runs. ADR-0020,
  ADR-0022, ADR-0038, ADR-0056, ADR-0127, ADR-0141, ADR-0159 and ADR-0160 cite
  listed ADRs but govern the Agent prompt result, Stop Requests, the
  Verification repair, Verification capacity, process residue, review Runs in
  the user checkout, independent Verification and the red repository gate,
  which this Spec does not touch, so they do not apply. This Spec adds
  ADR-0161. This Spec's gate is bound by ADR-0080, ADR-0091, ADR-0093,
  ADR-0096, ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156. All
  hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer approved Onda 2 of the
  efficiency sequence in chat on 2026-09-28; the Go test file rides the
  standing grant of 2026-09-21 for governed source, and the skill files ride
  the standing grant of 2026-09-18 for keeping the shipped skills true to the
  CLI, recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`, `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A squash-merged Spec's Runs are released without a manual `git worktree
  remove`, and only on proof.
- `roundfix deliver` releases a merged Spec's Runs by itself, right after the
  merge.
- `roundfix reconcile --apply` releases the same Runs for merges made outside
  the delivery loop.
- Legacy Runs whose checkout vanished are visible to every repository-scoped
  cleanup.
- Staging worktrees and half-removed item directories are cleaned without
  losing anything live.

## Core Features

1. **The merged-head proof.** A terminal Run of a merged Spec is proven when
   every commit its Run Branch holds and the merged head lacks is represented
   there: a Task commit of that Spec whose Task is `completed` at the merged
   head, a QA Report commit whose report the merged head supersedes, or a
   commit whose changed files the merged head holds byte for byte. The merged
   head is the pull request head the Delivery Queue item recorded at merge, or,
   without a record, the default branch head when the Spec is archived there.
   Only the files the Run changed are compared, so a moving `main` and an
   archive rename no longer defeat the proof. Anything else keeps the Run and
   names the first commit that is not represented.
2. **Automatic release after a delivery merge.** After the delivery owner
   records a merge, and before it removes the item worktree and item branch,
   it releases every terminal Run of that Spec the merged-head proof covers. A
   Run it cannot prove stays, and its reason is recorded as the item's cleanup
   warning in `roundfix deliver status`.
3. **Reconcile reads the merge record.** `roundfix reconcile` applies the
   merged-head proof, using the Delivery Queue's merge records when they exist.
   Dry-run reports the Run as `safe` or `superseded` with the proof as
   evidence, and `--apply` releases it through the existing revalidated path.
4. **Legacy Runs get their repository key.** Opening the Run Database for
   writing keys an unkeyed Run from its Run Worktree when its recorded checkout
   is gone, and the
   implement-preflight prune matches Runs by repository key, not by checkout
   path.
5. **Staging worktrees are owned and swept.** A carry-forward staging worktree
   is created and removed through the worktree administration lock with an
   owner record. A stale one is reported by dry-run and released by `--apply`
   and `--carry-forward`: its owner is proven dead, or, with no owner record, it
   is still `locked initializing` while Roundfix holds the lock. Every other
   production `git worktree` add, remove, prune or move goes through the same
   lock.
6. **A half-removed item never takes a nested worktree with it.** Cleaning an
   item directory that is no longer registered is refused, with the cleanup
   warning, when any registered worktree lives under it.

## Non-Goals / Out of Scope

- Integrating unique commits, repairing dirty work or choosing another target
  branch; a Run that cannot be proven stays.
- A force flag, an age threshold or a manual assertion that bypasses a proof.
- Folding `--discard-superseded` into `--apply`, or changing the Branch
  Disposition record.
- Pruning Task Worktrees and Task Branches in the post-merge release; the
  implement-preflight prune keeps that job.
- Specs under an external Specs Root; their Task files are not in the merged
  tree, and their Runs stay preserved.
- Git garbage collection, journal retention, or the GC Command.

## Success Metrics

1. The reproduction of the Spec 0172 Run — Task 01–04 commits contained in the
   merged head, a unique redone Task 05 commit and a unique failed QA Report
   commit — is classified `superseded` against both a Delivery Queue merge
   record and, without one, a default branch carrying the archived Spec, and
   `--apply` removes its Run Worktree and Run Branch.
2. The reproduction of the Spec 0164 Run — an archived Spec whose default
   branch later changed other files and files the Run also touched — is
   released, while a Run holding one commit that is not represented is kept
   with a reason naming that commit.
3. After a delivery merge, every provable terminal Run of the merged Spec is
   released before the item worktree and item branch are removed, and a Run
   that cannot be proven appears in the item's cleanup warning.
4. An unkeyed legacy Run whose checkout was removed is listed by
   `roundfix runs list` and `roundfix reconcile` from the main checkout once
   the Run Database has been opened for writing.
5. A staging worktree left `locked initializing`, or owned by a dead process,
   is reported by dry-run and released by `--apply`. One owned by a live
   process is kept. No production file outside `internal/worktree` runs `git
   worktree` add, remove, prune or move.

## Recorded limits

- A merge record lives only as long as its Delivery Queue: a later `roundfix
  deliver start` replaces a terminal queue. After that, the default-branch
  proof applies, and it needs the Spec archived on the default branch.
- The merged pull request head must still exist locally. Roundfix never runs
  Git garbage collection, but an operator-run `git gc --prune` can remove an
  unreachable head, and then the record proof falls back to the default branch.
- Command-level tests that exercise carry-forward keep staging under the
  process temporary directory; a leaked one is reclaimed by the owner proof,
  because the test process is gone.

## Decisions

- **The proof is per commit, not per tree.** A squash merge rewrites every
  commit, and the default branch keeps moving, so no tree equality survives.
  A Daemon Task commit carries its Spec and Task in trailers, and the merged
  head records whether that Task completed. That record, a superseding QA
  Report, or byte-identical content settles each commit on its own.
- **The delivery record comes first.** The Delivery Queue item holds the pull
  request head GitHub reported as merged. It is the strongest evidence
  Roundfix has and needs no network call. The archived Spec on the default
  branch is the fallback when the record is gone.
- **Release automatically, but only after a merge.** ADR-0115 rejected
  discarding a branch to get past a refusal. A recorded merge is a different
  trigger: the Spec's work has landed, and the maintainer asked for this
  cleanup to follow every merge.
- **`--apply` keeps its meaning.** It releases only `safe` and `superseded`
  Runs. The merged-head proof makes more Runs reach those states;
  `--discard-superseded` stays the separate, recorded act ADR-0115 requires.
- **A staging worktree is stale only on proof.** An owner record with a dead
  owner, or a `locked initializing` entry observed while the administration
  lock is held, is proof. Age is not.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The negative cases carry the weight: an unrepresented commit
released, a Task not completed at the merged head treated as done, another
Spec's merge record applied, a live staging worktree removed, or a nested
worktree deleted would each pass a happy-path test.

The outside-evidence row rests on state this Spec did not create: the built
binary's read-only `roundfix reconcile --format json` in the maintainer's own
checkout at `/Users/marcio/dev/roundfix`, whose retained Runs the Finding
names (Specs 0155–0164). The row records how each Run that still exists is
classified, and that row is blocked with its reason when none remains.

## Research basis

The Finding and the five Backlog Entries adopted under
[references/](references/_index.md) record the retained Runs, the proofs that
failed and their live reproductions. The Spec 0172 reproduction was observed
on 2026-09-25, after #259 merged. Spec 0171 serialized Git worktree
administration per repository in `internal/worktree`, and this Spec routes the
remaining production callers through it.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
