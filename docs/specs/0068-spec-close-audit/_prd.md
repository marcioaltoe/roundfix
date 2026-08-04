---
spec: 0068-spec-close-audit
status: active
created: 2026-08-02
surfaces: [backend, cli, docs]
---

# Spec close audit

The per-Spec loop ends at "squash merge and reconcile" and never audits what
the cycle created against what survived it. In one session that left four
kinds of residue at once: two Supervisor scratch worktrees never removed after
their push, a Run Worktree orphaned because the squash merge deleted the target
branch `reconcile` resolves by name, a remote backup branch whose purpose ended
when its work merged, and two Pull Requests opened and left unreviewed.

The last one is the expensive kind. Spec 0058's archive and five newly queued
Specs existed only on unmerged branches, so the default branch still showed
0058 active and showed none of the new Specs — work reported as delivered was
not where a reader would look for it. The maintainer found all four by running
`git branch -l` and `git worktree list` by hand. Evidence:
[a Spec cycle leaves branches and worktrees nobody audits](../../findings/2026-08-02-a-spec-cycle-leaves-branches-and-worktrees-nobody-audits.md).

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier is
  created; Run IDs, branch names, and worktree paths keep their existing
  contracts. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the governing clause prohibits
  reading, printing, committing, or generating secrets and forbids inventing
  authentication, authorization, transport, or deployment policy; the audit
  reads local Git state and reports. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0052 protects terminal completion,
  so no audit path may touch an Active Run or a live writer; the audit reports
  and never reclaims on its own. ADR-0053 governs proof-based Run Worktree
  reconciliation, whose missing-target case this Spec resolves by content.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized; the work is product code, CLI surface, and documentation.
  Source: `docs/agents/agent-instructions.md`.

## Goals

- Closing a Spec reports every branch and worktree the cycle created, which
  survive, and why each survivor is legitimate.
- Work a Spec claims to have delivered is verified present on the synced
  default branch, so an unmerged Pull Request cannot read as delivered.
- A Run whose target branch was deleted by its own merge is resolved by
  content instead of degrading to `unknown` forever.
- Nothing is reclaimed on ambiguity.

## Core Features

1. A read-only close audit, runnable at merge-and-sync, that enumerates the
   branches and worktrees associated with a Spec and classifies each survivor:
   backing an open Pull Request, awaiting integration, or residue.
2. Delivery verification: every artifact the Spec claims — its implementation,
   its archive, any Spec folder it created — is checked present on the synced
   default branch, and anything absent is reported with the branch that holds
   it.
3. Content-based Run integration resolution: when a Run's target branch no
   longer exists, integration is decided by comparing the Run Branch against
   the default branch — no file unique to the Run Branch and identical
   implementation files proves integration — replacing today's `unknown`
   outcome.
4. Supervisor scratch worktrees are reclaimed once their branch is pushed, so
   the branch remains and the working copy does not.
5. Every classification carries the evidence that produced it, and anything the
   audit cannot classify is reported as preserved rather than removed.
6. The audit is discoverable without knowing which Runs or branches a Spec
   produced.

## Non-Goals / Out of Scope

- Deleting branches that back an open Pull Request, or any branch carrying
  unintegrated work.
- Automatically merging, reviewing, or closing Pull Requests.
- Terminating processes or reclaiming Run process trees, owned by Spec 0066.
- Storage compaction and artifact sanitation, owned by Spec 0059.

## Success Metrics

- Replayed against this session's end state, the audit reports all four residue
  kinds: two scratch worktrees, one orphaned Run Worktree, one stale remote
  backup branch, and two branches held by unmerged Pull Requests.
- A Run whose target branch was deleted by `--delete-branch` resolves as
  integrated by content instead of `unknown`, proven against that exact case.
- A Spec whose archive lives only on an unmerged branch is reported as
  not-yet-delivered, naming the branch that holds it.
- An Active Run is never touched by the audit, proven by an injected active
  fixture.

## Decisions

- The audit reports; it does not reclaim. Naming residue is the deliverable,
  and an aggressive cleanup that guesses is worse than debris.
- Integration is proven by content when the branch name is gone, because the
  name is the handle that `--delete-branch` destroys and the content is what
  actually matters.
- This Spec evolves the close of the loop and never regresses it: no existing
  command changes behavior, and `reconcile` keeps every refusal it makes today
  except the one the content check now resolves.

## Open Questions

None.
