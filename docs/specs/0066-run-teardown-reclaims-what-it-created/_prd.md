---
spec: 0066-run-teardown-reclaims-what-it-created
status: active
created: 2026-08-01
surfaces: [backend, cli, docs]
---

# Run teardown reclaims what it created

A Run that ends without passing leaves two kinds of debris behind, and nothing
reclaims either.

Four QA cycles on Spec 0053 left four Run Branches that neither
`reconcile --apply` nor Branch Integrity Preflight could clear, and the
accumulation then *refused* a `roundfix watch` on the Pull Request those cycles
were trying to make Clean. Spec 0053 shipped a `superseded` classification for
a Run Branch holding nothing but an obsolete QA report; the case its design
does not cover is several such branches at once, none passing. Evidence:
[failed QA runs accumulate Run Branches nothing can release](../../findings/2026-07-30-failed-qa-runs-accumulate-unreleasable-run-branches.md).

Separately, terminating a Run does not reach the `acpx` child it started. Four
processes from Spec 0037's QA fixture were found still running after **three
days and six hours**, reparented to init, pointing at worktrees that no longer
existed. Nothing had ever told them to stop. Evidence:
[run termination does not reach the acpx child](../../findings/2026-07-30-run-termination-does-not-reach-the-acpx-child.md).

Both are the same contract failure: a Run creates a branch and a process tree,
and its teardown reclaims neither reliably.

## Project Constraints

- Identifier strategy: not applicable — Run IDs, Run Branch names, and process
  identities keep their existing contracts; no project-owned Internal
  Identifier is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the governing clause prohibits
  reading, printing, committing, or generating secrets and forbids inventing
  authentication, authorization, transport, or deployment policy; all behavior
  is local process and Git state management. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0052 protects terminal completion,
  so no reclamation may touch an Active Run or a live writer; ADR-0044 relies
  on the process-gone signal to reclaim orphaned locks, and any change to
  process termination must preserve how that signal is distinguished from a
  host that cannot answer. ADR-0053 governs proof-based Run Worktree
  reconciliation, which this Spec extends to accumulated failed-cycle Run
  Branches. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is
  proposed or authorized; the work is product code, CLI surface, and
  documentation. Source: `docs/agents/agent-instructions.md`.

## Goals

- Terminating a Run terminates the process tree it created, including the
  adapter child, and proves it rather than assuming it.
- A set of Run Branches left by failed QA cycles can be classified and cleared
  without hand-deleting refs.
- Accumulated Run Branches stop blocking the review work that would resolve
  the Spec they belong to.
- Nothing reclaimable is reclaimed while it could still be in use.

## Core Features

1. Run termination reaches the whole process tree the Run started, including a
   child that outlives its immediate parent, and reports which processes it
   terminated and which it could not.
2. A process the Run cannot prove terminated is reported as such, distinctly
   from one proven gone, so a host that cannot answer never reads as success.
3. Run Branches from failed cycles on the same target are classifiable as a
   set: the newest carries the current evidence, and the older ones holding
   only superseded reports become releasable.
4. Branch Integrity Preflight distinguishes Run Branch work that must block a
   review Run from accumulated superseded work that must not, so a Pull
   Request can be made Clean while its failed cycles are still on disk.
5. A reclamation pass is dry-run first, names what it would remove and the
   proof that makes it removable, and preserves anything ambiguous.
6. An orphaned process or branch older than the Run that created it is
   discoverable without knowing the Run ID.

## Non-Goals / Out of Scope

- Changing Journal Retention, the Run Event Journal, or artifact retention.
- Deleting Run Branches that carry unintegrated implementation work.
- Terminating processes Roundfix did not start.
- Storage compaction and global artifact sanitation, owned by Spec 0059.

## Success Metrics

- Stopping a Run leaves no descendant of that Run running, verified against a
  fixture that starts an adapter child which outlives its immediate parent.
- A repository with four failed-cycle Run Branches on one target reaches a
  state where `watch` proceeds, without any branch being hand-deleted.
- A dry-run reclamation pass on a machine with orphaned branches and processes
  names every candidate with its proof, and a second pass after applying is a
  no-op.
- An Active Run is never touched by any reclamation path, proven by an injected
  active fixture.

## Decisions

- Teardown proves termination rather than assuming it; an unprovable
  termination is reported, never silently treated as done.
- Reclamation is dry-run first and preserves ambiguity, matching the posture
  the GC and reconcile paths already take.
- This Spec evolves Run lifecycle and never regresses it. The declared break is
  that Branch Integrity Preflight stops blocking on superseded Run Branch work;
  every other refusal it makes today it still makes.

## Open Questions

None.
