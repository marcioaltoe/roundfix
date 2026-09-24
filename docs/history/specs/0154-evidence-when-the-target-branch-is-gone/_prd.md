---
spec: 0154-evidence-when-the-target-branch-is-gone
status: archived
created: 2026-09-23
surfaces: [backend, cli, docs]
archived: "2026-09-24"
source_slug: 0154-evidence-when-the-target-branch-is-gone
---


# Evidence when the target branch is gone

Reconciliation asks one question about a terminal Run: did its work reach the
target branch? It answers with `merge-base --is-ancestor`, and when that misses
it looks for content evidence — a superseding QA Report on the target for the
same Spec.

Both paths need the target branch to exist. Once a pull request is squash-merged
and its branch deleted, neither can run, and the Run is preserved with the
reason `target branch "<name>" is absent`. Preserved is the safe answer; it is
also permanent.

Measured on this repository on 2026-09-22: 22 retained Run Worktrees, of which
`reconcile --apply` released zero. Every one was preserved for an absent target
branch. Only 1 of 22 Run Branches was an ancestor of `main`, because a squash
merge writes a new commit and leaves the original commits outside the default
branch's lineage — so ancestry would have proved nothing even had the branches
survived.

The work is in `main`. Every one of those Specs is archived there. Nothing in
reconciliation can see it.

This Spec is carved from Spec 0125 Core Feature 4, which asks for the present
squash and missing-ref cases to be characterized, and for delivered content to
be proven through an approved local evidence rule or returned as a truthful
preserved result with the missing proof named.

## Project Constraints

- Identifier strategy: applicable — a Run keeps its identity, its Branch and its
  Worktree path; this Spec changes what reconciliation can prove about them, not
  what they are called. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git reads only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0062 preserves operational history
  during identity-reset planning, ADR-0093 checks Spec consistency by citation,
  ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a path
  governed once bounded, ADR-0155 makes the `qa` Task declare the matrix and
  ADR-0156 makes a declared promise name a consuming Task. This Spec's own gate
  is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117, which keep an
  environment-blocked row distinct from a failure, make the gate a Task node of
  its own type, prove machine facts before spending an agent turn, carry a row
  forward only on declared unmoved evidence, and check a defect at the stage
  that can produce it. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the shipped skill documents what `reconcile`
  reports. Express maintainer authorization: the standing grant of 2026-09-18
  for the skill files and the standing grant of 2026-09-21 for governed paths a
  slice needs, both recorded in [_authorization.md](_authorization.md); bounded
  files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. The
  bounded set is the intersection of this Spec's changed paths with the literal
  set in `internal/speccheck/governed.go`, computed rather than predicted;
  `internal/worktree/worktree.go` and its test are ordinary source. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Goals

- A Run whose work reached the default branch can be proven so after its target
  branch is deleted.
- A Run whose work did not reach it stays preserved, with what is missing named.
- Nothing is released on an absence of evidence.

## Core Features

1. **A squash merge is characterized, not assumed.** Ancestry cannot prove a
   squash-merged branch reached its target, and the repository's own history is
   the proof: one Run Branch in twenty-two is an ancestor of `main`. The
   reconciliation result says which proof it tried and which it could not.
2. **An absent target falls back to the default branch.** When the target branch
   is gone, reconciliation seeks the same content evidence against the default
   branch instead of refusing to look. The evidence rule is the one already
   accepted for a missed ancestry, not a second one.
3. **Evidence, never inference.** A Run is released only on positive content
   evidence. An unreachable default branch, an unarchived Spec, or evidence that
   does not name this Run's Spec each leave it preserved, with the missing proof
   named in the reason.
4. **Preserved stays truthful.** The reason distinguishes "the target branch is
   absent and the default branch does not carry this Spec's work" from "the
   target branch is absent and nothing was checked", which is what it says
   today.

## Non-Goals / Out of Scope

- The repository identity shared by main and linked worktrees, Spec 0125 Core
  Feature 2.
- The branch-naming policy and the skill's branch examples, Core Features 1 and
  5 of Spec 0125.
- Weakening the deleted-target content proof Core Feature 3 tells this Spec to
  retain.
- Releasing a Run whose Worktree holds uncommitted changes. Dirty stays dirty.

## What this Spec does not yet guarantee

Independent pre-PR review found two defects after the corrective ceiling and the
review cap were both reached. The maintainer chose on 2026-09-24 to deliver with
them recorded rather than take a third corrective Task. Both are carried to the
restructured queue's delivery 3 (verification and settlement reliability), which
owns reconciliation.

**`reconcile --apply` can panic on a revalidated fallback candidate.** With a
deleted target and a clean Run Branch whose default branch carries an archived
superseding report, the absent-target fallback can make the branch a cleanup
candidate. On revalidation `InspectTerminalRun` compares the active report path
with the archived one and returns an unintegrated result without evidence, and
`ApplyRunBranchCandidate` still calls `cleanupTerminalRun`, which reads
`fresh.evidence.worktreePresent` through a nil evidence pointer. This path did
not exist before this Spec, so the defect is introduced here. The repair is to
carry the reconciliation evidence through that result, or refuse the candidate
before cleanup.

**An archived copy of the same QA Report is not recognised.** Archive moves
`qa-report-YYYY-MM-DD(.NN).md` from `docs/specs/<slug>/qa` to
`docs/history/specs/<slug>/qa` without renaming it, so the Run-side and
default-side reports share date and sequence. `NewestQAReportFromPaths` breaks
the tie by full path, `docs/history/...` sorts before `docs/specs/...`, and the
active Run path wins, so `proven` is false. This errs toward preserving a Run,
which is the safe direction. The repair is to compare report identity
independently of the active or archived root.

## Success Metrics

1. A fixture whose Run Branch was squash-merged into the default branch and
   whose target branch was deleted is classified releasable, naming the content
   evidence that proved it.
2. The same fixture without that evidence in the default branch stays preserved,
   and its reason names what was missing.
3. Every currently accepted classification — safe, superseded, unintegrated,
   dirty — is unchanged for a present target branch.

## Decisions

- **Ask the default branch, because that is where merged work lands.** The
  target branch is a pull request's branch and it is deleted on merge by design.
  Treating its absence as unknowable makes the common case permanent.
- **Reuse the accepted evidence rule.** A second way to prove integration would
  be a second thing that can disagree about what integration means.
- **Absence of proof is not proof of absence — and not proof of presence.** The
  release path requires positive evidence. Everything else is preserved, which
  costs disk and loses nothing.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a fallback that released on a
missing default branch, or on evidence naming another Spec, would clear the 22
Worktrees and could clear one holding work nobody has.

## Research basis

Measured on this repository, not inferred. `reconcile --apply` released zero of
22 and preserved every one for an absent target branch. `merge-base
--is-ancestor` classified 21 of 22 Run Branches as non-ancestors of `main`,
which is the expected consequence of squash merging. The ancestry check sits at
`internal/worktree/worktree.go:945`, its content-evidence fallback
`supersedingQAReport` immediately after it, and the absent-target path at
`:422`, which returns before either runs.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
