---
spec: 0101-a-terminal-branch-with-one-disposition
status: active
created: 2026-08-12
surfaces: [backend, cli]
---

# A terminal branch with one disposition

A Run that ends leaves commits on a Run Branch that only a human can dispose of.
Integration is proven by ancestry, which fails after a rebase even when the
content is already integrated, so branches that carry nothing new are classified
as unintegrated and block the next Run. Preflight then prescribes a recovery —
merge the branch — that would in one measured case have regressed the repository
with work deliberately rejected. A worktree of the same repository gets its own
identity, so each side sees the other's Run Branches as pending work. In two
sessions the only exit was disabling the guard entirely, which also disables it
for work that genuinely is pending. A guard whose normal path is impassable
teaches its own bypass.

## Project Constraints

- Identifier strategy: applicable — Run Branch, Run Worktree, Task Worktree, and
  the repository identity a Run is scoped by are glossary terms this Spec changes
  the resolution of, and a new terminal classification is vocabulary the glossary
  must own. The closing node checks whether the work introduced or changed a term.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work reads local Git refs, trees, and the Run
  Database. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0111 makes an unobserved Verification
  unknown rather than a verdict, and this Spec applies the same discipline to
  integration: an ancestry check that cannot see the content is unknown, not
  unintegrated. No accepted ADR governs Run Branch classification or repository
  identity, which is why both rules this Spec adds are new. ADR-0135 makes an absent diagnostic a reported state in the repair prompt an Agent reads; this Spec gives a terminal branch one disposition and writes no Verification feedback, so it does not apply. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the reconcile, preflight, and store
  packages plus their tests. Source: `docs/agents/agent-instructions.md`.

## Goals

1. A branch whose content is already integrated is released without a bypass.
2. A branch that carries work nobody wants has a recorded terminal disposition.
3. A repository and its worktrees share one identity, so each sees one set of Runs.
4. Disabling the integrity guard stops being the routine way to deliver.

## User Stories

1. As a Supervisor whose branch was rebased before its Pull Request, I want a Run
   Branch whose content is already in the target to be released, so that I do not
   disable the guard to deliver.
2. As a Supervisor holding a Run Branch with deliberately rejected work, I want to
   record that it is superseded with a reason, so that reconciliation can release
   it without me merging work I refused.
3. As a Supervisor whose Run's target branch was merged and deleted, I want that
   branch classified as releasable, so that a finished Spec does not leave a
   permanent unresolvable entry.
4. As a Supervisor working in a Git worktree of a repository, I want Runs and
   branch checks scoped to the repository rather than to the checkout, so that
   each side does not report the other's work as pending.

## Core Features

1. **Integration proven by content.** A Run Branch is classified as integrated
   when its content is present in the target — an identical tree, or the absence
   of any file it holds exclusively — not only when ancestry proves it. Ancestry
   remains sufficient where it holds.
2. **A recorded supersession.** A terminal Run Branch can be recorded as
   superseded with a reason, which reconciliation reads as a disposition rather
   than requiring a merge it cannot justify.
3. **A deleted target is a releasable branch.** When a terminal Run's target
   branch no longer exists, the Run Branch is classified releasable rather than
   left permanently unintegrated.
4. **One identity per repository.** A repository and its Git worktrees resolve to
   one identity, so Run listing, branch integrity preflight, and reconciliation
   see one consistent set.
5. **A preflight that never prescribes a regression.** A prescribed recovery must
   be one that cannot lose or revert work; where no such recovery exists, the
   preflight names the state and the dispositions available instead of naming a
   merge.

## User Experience

Reconciliation reports each terminal Run Branch with one classification and one
available disposition. A preflight refusal names the branches, why each is
unresolved, and the disposition each needs — never a merge that would regress the
tree. Recording a supersession takes a reason and reports what it released.

## Non-Goals / Out of Scope

- Removing the branch integrity guard or its bypass; this Spec makes the normal
  path passable so the bypass stops being routine.
- Automatic integration of pending work. A branch carrying real work still
  requires a human decision.
- Changing Run termination, the Run Database schema beyond what identity requires,
  or garbage collection policy.
- Cross-repository branch operations.

## Success Metrics

- A rebased branch whose content is integrated is released without the bypass,
  measured against a session in a repository this Spec did not build where five
  terminal branches forced the bypass on 2026-08-08.
- A superseded branch is released with its reason recorded, without a merge.
- A Run whose target branch was deleted reaches a terminal disposition.
- A worktree and its repository report one set of Runs.
- A preflight refusal in the measured superseded-branch case no longer prescribes
  a merge.

## Decisions

- Content proof supplements ancestry rather than replacing it; ancestry is cheap
  and correct where it holds, and the content path exists for where it does not.

## Open Questions

- Which content proof is authoritative — an identical tree, or the absence of any
  exclusively-held file. The default until answered is to accept either, since a
  branch satisfying one is integrated under both readings.
- Whether recording a supersession is a command a maintainer runs or a
  classification reconciliation offers to apply. The default is an explicit
  maintainer action, because releasing work is not reversible by re-running.
