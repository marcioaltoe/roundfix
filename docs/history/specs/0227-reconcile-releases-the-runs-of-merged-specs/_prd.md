---
spec: 0227-reconcile-releases-the-runs-of-merged-specs
status: archived
created: 2026-10-04
surfaces: [backend, cli, docs]
archived: "2026-10-05"
source_slug: 0227-reconcile-releases-the-runs-of-merged-specs
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: row 07 needs historical Run Database items of the hand-cleaned Specs, which no longer exist (intervention log entries 144-145 record them), and row 10 has no Pull Request yet; the queue opens it. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 03454c09d9f23c0ec5c69c1d981267fc673a1be7
---


# Reconcile releases the Runs of merged Specs

On 2026-10-04 the maintainer found twelve `roundfix/run-…` and item branches
in the main checkout. `roundfix reconcile --apply` released only the four it
classified `superseded`. It kept six Runs of Specs 0181, 0184 (two), 0200,
0204 and 0207, every one archived and merged on main, as `unintegrated` or
`dirty`, and it never reported the item branches of 0205 and 0217 that the
Delivery Queue had recreated. The operator removed them by hand
(`~/.roundfix-operator/queue-interventions.md`, entry 144). Every queue merge
of the day also ended with `warning: cleanup failed: release merged Spec
"<slug>" Runs: candidate head "<sha>" is not represented by merge commit
"<sha>"` (0222, 0224, 0226), so a delivered Spec's Runs waited for a manual
reconcile ([references/2026-10-04-reconcile-keeps-runs-of-merged-specs.md](references/2026-10-04-reconcile-keeps-runs-of-merged-specs.md)).

This is a bug fix. Its three causes were measured during authoring, with the
shipped binary on disposable repositories and with Git on this repository's
own merges:

1. When no Delivery Queue merge record remains, reconcile proves a Run of a
   merged Spec against the current default branch. A commit the Run inherited
   that is neither a Task commit nor Spec-directory work, such as an operator
   commit on the item branch, is compared file by file with the default
   branch. A fixture whose later Spec edited one such file and renamed
   another flipped from `superseded` to `unintegrated`, "1 Run-only file, 1
   differing shared file against default branch "main"".
2. The post-merge cleanup accepts a candidate head only when it is an
   ancestor of the merge commit or has the same tree. A squash merge onto a
   default branch that moved has neither: for 0224 and 0226 one other Spec's
   merge landed during the item, and the candidate's tree differs from the
   merge commit's, while `git merge-tree --write-tree <merge>^1 <candidate>`
   reproduces the merge commit's tree exactly for 0223, 0224 and 0226.
3. Reconcile inspects only Run Worktrees and Run Branches. A
   `roundfix/deliver-<slug>-<16 hex>` branch left behind when the queue
   recreated an item is never reported, whatever its Spec's state.

## Prerequisites

None. Specs 0228 and 0229 are authored in the same cycle; if either also
raises the Roundfix Skill's version, the operator orders the queue so that the
later Spec raises it from the earlier one's value.

## Project Constraints

- Identifier strategy: not applicable — no identifier changes. Runs keep their
  Run identifiers, item branches keep the `roundfix/deliver-<slug>-<16 hex>`
  name `deliver start` gives them, and the `roundfix-reconcile/v1` schema
  name is unchanged because the report only gains a field. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — every proof reads local Git and
  the Run Database; no credential, forge read or network call is added, and
  no test reaches GitHub. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0232 (this Spec) decides the merge
  evidence, ADR-0232: "Without merge evidence the content proof of ADR-0161
  and ADR-0212 stays exactly as it is". ADR-0161 releases a merged Spec's
  Runs on the merged head, ADR-0161: "The delivery owner releases the Spec's
  Runs after it records the merge", and ADR-0212 supersedes Spec-directory
  work and declared leftovers, ADR-0212: "An uncommitted change outside that
  scope, such as a hand edit to an unrelated file, still keeps the Run".
  ADR-0053 keeps reconciliation explicit and read-only by default, ADR-0053:
  "is read-only unless `--apply` is supplied", and ADR-0115 discards a superseded
  Run Branch through Roundfix with its evidence first, ADR-0115: "The disposition is a
  separate, named act rather than a widening of `reconcile`". ADR-0187 and
  ADR-0189 govern the Roundfix Skill edit and its version, and ADR-0184 has
  the TechSpec state changed command surfaces, ADR-0184: "A TechSpec now
  declares numbered Surface Transcripts". The gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0104, ADR-0155 and ADR-0156, and ADR-0093,
  ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency by
  citation and receipt. ADR-0178 and ADR-0179 decide the grant each Task
  commit runs under, and ADR-0166 records undeclared paths; every Task
  declares its paths. ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry, ADR-0167 cites ADR-0080 but decides the pre-PR row of a
  qualifying partial, ADR-0170 cites ADR-0053 but decides Task Carry-Forward,
  which this Spec leaves unchanged, ADR-0182 cites ADR-0117 but settles a Task on
  the facts its gate will check, and ADR-0192 cites ADR-0178 but decides
  derived-path conflicts. ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097
  but decide what a QA row records, when it is observed again and its
  evidence snapshot, and ADR-0229 cites ADR-0167 but decides the Delivery
  Retry of an operator archive; this Spec changes none of them, so none
  applies. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and
  its `SKILL.md` mirror are Governed Paths, and the maintainer authorized
  skill edits ("considere autorizado a ajustar todas as skills se
  necessário") and this Spec's cycle (2026-10-04, "Aprovado"). No other
  Governed Path changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0227-reconcile-releases-the-runs-of-merged-specs/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A terminal Run of a Spec that carries merge evidence is released by
  `roundfix reconcile --apply` and by the delivery owner after a merge, even
  when files it changed diverged on the default branch later, while a Run
  of a Spec without merge evidence keeps today's content proof.
- The post-merge cleanup proves a merge commit that holds the candidate's
  changes plus default-branch commits that landed during the item, so the
  queue no longer warns on such merges.
- `roundfix reconcile` reports, and `--apply` releases, an item branch of a
  merged Spec that no live Delivery Queue item records.
- Uncommitted work outside a merged Spec's declared scope still keeps its Run
  or item branch.

## Core Features

1. **Merge evidence.** The default-branch head holds the Spec's archived
   `_prd.md`, and the commit that added it, the delivery commit, is not
   reachable from the branch being proven. With it, every Task commit must
   name the Spec and its Task must be completed in the archived Spec; every
   other commit is superseded by the delivery (ADR-0232).
2. **A cleanup proof for a moved default branch.** A merge commit whose tree
   equals the merge of the candidate head into its first parent represents
   the candidate (ADR-0232).
3. **Item branches in reconcile.** A new `itemBranchCandidates` list in the
   reconcile report names each item branch of a merged Spec, its head, its
   proof and its action; `--apply` deletes it, with its worktree when that
   worktree is clean; anything unproven goes to `preservedCandidates`.
4. **The skill and the guide say so.** The Roundfix Skill's `reconcile`
   reference and the `reconcile` command guide describe merge evidence, the
   item branch candidates and the unchanged dirty rule.

## Non-Goals / Out of Scope

- Releasing uncommitted changes outside a merged Spec's declared or recorded
  scope, or a dirty item worktree.
- Changing the content proof, QA supersession or Task status proof for a Spec
  that is not archived on the default branch.
- Reading GitHub to prove a merge; every proof stays local.
- Changing `--discard-superseded`, `--carry-forward`, the Branch Disposition
  record, `deliver start`'s item-branch selection, or the queue's own
  cleanup of the item it just merged.
- Reclaiming Run artifacts or the Run Database; that stays with `roundfix gc`.

## Success Metrics

1. Success Metric: in a disposable repository, a terminal Run of a merged
   Spec whose inherited operator commit touched a file a later Spec edited
   and another it renamed is `superseded` and released by `--apply`, where
   the shipped binary reports `unintegrated` with "1 Run-only file, 1
   differing shared file".
2. Success Metric: the same Run with a Task commit whose Task is pending in
   the archived Spec, a Run of a Spec that is not archived, and a Run that
   already holds the delivery commit stay `unintegrated`; a dirty path
   outside the declared scope stays `dirty`.
3. Success Metric: `ReleaseMergedRuns` releases the Runs of an item whose
   squash merge landed after another commit on the default branch, and still
   refuses an unrelated candidate with "not represented by merge commit".
4. Success Metric: an item branch of a merged Spec that no live item records
   is listed in `itemBranchCandidates` by the dry-run and is gone after
   `--apply`; the branch of a live item, a dirty item worktree and an item
   branch of an unmerged Spec stay.
5. Success Metric: every existing reconcile, merged-head and release test
   passes, with one declared change: the archived-Spec case of
   `TestArchivedSpecDirectoryProofStillComparesOtherCommitPaths` now expects
   `superseded`. The Run-ID report keeps its JSON keys; only a full scan
   carries `itemBranchCandidates`.

## Acceptance evidence

The outside-evidence row rests on records this Spec did not produce:

- The operator's intervention log
  (`~/.roundfix-operator/queue-interventions.md`), entry 144 (2026-10-04):
  "removed 6 preserved Run worktrees+branches (0181, 0184 x2, 0200, 0204,
  0207; all Specs archived and merged) and 2 superseded item branches (0205
  old, 0217 old)".
- The live Run Database, read with `sqlite3 -readonly` during authoring on
  2026-10-04: the merged 0224 and 0226 items record candidate heads
  `88e2fe3c` and `e62588bd` and merge commits `f0353479` and `69ace45d`;
  neither candidate is an ancestor of its merge commit nor has its tree, and
  `git merge-tree --write-tree` of each merge commit's first parent with its
  candidate equals the merge commit's tree, as it does for 0223.
- This repository's default branch: for every Spec the operator cleaned by
  hand (0181, 0184, 0200, 0204, 0205, 0207 and 0217) the commit that added
  its archived `_prd.md` is its delivery's squash merge (`122098e7`,
  `404cf0c8`, `bfebc0e9`, `58e6d57f`, `14e049b3`, `18ef15eb`, `23000806`).
- An adopter's independent report, the Fiscus repository's Backlog Entry
  `2026-09-01-o-squash-merge-impede-o-reconcile-de-provar-integracao`
  (Secondbrain mirror): after squash merges reconcile reported "116 Run-only
  files, 54 differing shared files", and its hand audit found the
  divergences were the default branch moving on, not lost content.
- Git's documentation of `git merge-tree`
  (<https://git-scm.com/docs/git-merge-tree>): the `--write-tree` mode
  "Performs a merge, but does not make any new commits and does not read from
  or write to either the working tree or index".
- Published squash-merge cleanup practice
  (<https://www.git-automation.com/conflict-resolution-safe-merge-operations/squash-fixup-strategies/detecting-squash-merged-branches/>):
  "Delete a branch only when it is squash-merged by content and its tip
  matches a merged pull request's head, or when it is an ancestor of main",
  the same pairing of a merge record with Git evidence this Spec uses.

## Decisions

- Merge evidence replaces the content comparison only for Specs archived on
  the default branch by a delivery commit the branch does not hold. See
  ADR-0232.
- Dirty work keeps ADR-0212's scope rule; a dirty item worktree is never
  released. See ADR-0232.
- The cleanup proves a moved default branch with `git merge-tree`, failing
  closed on a Git without that mode. See ADR-0232.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
