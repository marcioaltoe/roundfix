---
spec: 0142-an-absent-ref-is-named-absent
status: active
created: 2026-09-17
surfaces: [backend, cli]
---

# An absent ref is named absent

Reconciliation tells a maintainer why it cannot release a Run's debris. On this
machine it told them the wrong thing: fourteen candidates were refused with
`resolve local branch "ma/0118-a-task-proved-once-does-not-run-twice": short ref
is ambiguous`, for a branch that does not exist at all. Git lists no ref by that
name.

The reason is a counting rule. The resolver asks how many candidate refs share
the short name and refuses when the count is not one, so an absent ref — count
zero — is reported with the message written for a name that matches several. The
existence check that would tell them apart already exists beside it and is not
used here; the validity check that runs first only asks whether the name is
well formed.

The refusal then costs more than a wrong sentence. A Run set whose recorded
target branch is absent fails classification as a whole, so every candidate in
it is preserved with a reason that names neither the missing ref nor what would
recover it. The supported cleanup path is unavailable and the maintainer is
pointed at an ambiguity that does not exist.

This Spec is the first slice carved from Spec 0125, which keeps repository
identity across linked worktrees and the delivered-content evidence rule.

## Project Constraints

- Identifier strategy: applicable — Run identifiers, Run Branch names and the recorded target branch keep their spelling; this Spec reads refs by their fully qualified names and coins no identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — reconciliation reads the local Git repository and the Run Database, and opens no credential store or network transport. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the reconciliation outcome vocabulary and its preservation rules are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0057 applies: the Daemon owns Run and Task status, and this Spec changes no status, only what reconciliation reports and classifies.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the adopted finding.
  ADR-0127 applies: process residue is a readiness observation, not authority to settle work, so naming an absent ref creates no Run and settles nothing.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/worktree`, which is not a Governed Path, and its test files are ordinary too. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A ref that does not exist is reported as absent, and a name that matches
  several refs is reported as ambiguous.
- A Run whose recorded target branch is absent is classified and preserved with
  a reason that names the missing ref, instead of failing the whole set.
- Every release, preservation and content rule that applies when the target
  branch exists keeps its current behavior.

## User Stories

1. As a maintainer running reconciliation after old branches were deleted, I
   want the refusal to say the branch is gone, so that I stop looking for an
   ambiguity that does not exist.
2. As a maintainer with one absent target branch among many Runs, I want the
   other Runs still classified, so that one deleted branch does not hide every
   candidate's state.
3. As a maintainer reading a preserved candidate, I want the missing ref named
   in its reason, so that I know exactly what would have to exist to prove
   delivery.

## Core Features

1. **Absent and ambiguous are different answers.** The resolver distinguishes a
   short name that matches no ref from one that matches several, and each
   carries its own message naming the branch.
2. **An absent target branch preserves, and does not refuse.** When a Run set's
   recorded target branch does not exist locally, each candidate in the set is
   classified and preserved with a reason naming the absent branch. The set's
   classification no longer fails as a whole.
3. **Everything provable stays provable.** When the target branch exists, the
   existing release, supersession, content-proof and preservation rules reach
   exactly the outcomes they reach today.

## User Experience

A maintainer runs reconciliation in a repository whose old target branches were
deleted. Each affected candidate is preserved with a reason that names the
branch that is gone. Candidates whose target branch still exists are classified
as before, so one absent branch no longer decides what the whole report says.

## Non-Goals / Out of Scope

- Proving delivered content when the target branch is absent, or under a squash
  merge. That evidence rule stays with Spec 0125.
- Repository identity shared between the main checkout and linked worktrees,
  and the migration of checkout-derived Run records. Also Spec 0125.
- Releasing anything that is preserved today. This Spec changes what
  reconciliation says and classifies, never what it deletes.
- Changing Run Branch naming, the Run Database schema, or the reconciliation
  outcome vocabulary.
- Repairing the machine state that produced the observation; the deleted
  branches stay deleted.

## Declared intentional breaks

- A refusal that reads `short ref is ambiguous` for a branch that does not exist
  no longer appears. Any tooling or test pinned to that sentence for an absent
  ref reads the absent-ref sentence instead.

## Regression locks

- A name that genuinely matches several refs keeps the ambiguity refusal.
- A candidate released today is still released; a candidate preserved for a
  reason other than an absent target branch keeps that reason.
- The Run Database, Run Branch names and worktree paths are untouched.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the Run
Database on this machine, where fifteen Runs recorded a target branch that has
since been deleted, and the reconciliation output that refused fourteen of them
with the ambiguity sentence. The replay reproduces that shape in a fixture
repository — a recorded target branch with no ref — and shows the absent-ref
reason in place of the ambiguity one. Where the recorded observation cannot be
read, the row records that reason and does not block.

## Success Metrics

1. A fixture whose recorded target branch has no ref reports the absent-ref
   reason, and the same fixture on the unchanged tree reports the ambiguity one.
2. A fixture whose short name matches both a tag and a branch still reports the
   ambiguity reason.
3. A fixture mixing one absent target branch with one present target branch
   classifies the present one exactly as it does today.

## Decisions

- **Two answers, not one.** Counting refs cannot distinguish absence from
  ambiguity, and the existence check that can already lives beside the resolver.
- **Preserve, do not release.** An absent target branch proves nothing about
  delivery, so the candidate stays preserved; only its reason changes.
- **Slice, not portfolio.** Spec 0125 carries five Core Features; this Spec takes
  the truthful-refusal rule and leaves repository identity and the delivered
  content evidence there.

## Research basis

**Secondbrain.** Consulted before authoring, index first, then a query on
truthful failure reporting in repository tooling. The results were dominated by
mirrors of this repository, which are references rather than independent
knowledge, and no source changed the design.

The repository's pending Inbox Entries were read before authoring. One was the
source for this Spec: the 2026-09-16 capture that recorded the ambiguity refusal
over fourteen candidates, written while reconciling this machine's Run debris.
It was promoted to a finding, adopted by this Spec, and its entry triaged to the
adopted copy.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read, so no external validation is claimed.

**Local measurement.** The counting rule, the unused existence check, the
name-shape validity check and the set-level classification failure were each
read in this repository's source before this PRD was written.

## Open Questions

None.
