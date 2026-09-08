---
spec: 0125-repository-identity-and-run-branch-policy
status: active
created: 2026-09-08
surfaces: [backend, cli, docs]
---

# One repository keeps one identity and its chosen branch policy

The maintainer has replaced personal branch prefixes with purpose-based work branches and explicitly permits the existing Roundfix Run/Task namespace. Linked worktrees derive different repository identities from their checkout paths. Reconciliation now proves the reported deleted-target squash case, but other content and reference cases remain narrower. The maintainer needs shared repository identity and consistent branch policy without losing old Run resources or deleting unproved work.

This Spec is **in authoring**. Its protected scope and open decisions are
pending approval. The accompanying authorization is a proposal with no grant;
there is no TechSpec, Task Graph, or authority to start implementation.

## Project Constraints

- Identifier strategy: applicable — canonical repository identity must group linked worktrees while preserving access to existing Run records; keep Run/Task IDs stable and recognize old branch names during migration. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local repository/ref identity and reconciliation introduce no authentication or HTTP endpoints. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0138 preserves one commit per verified Task and the existing Clean-only opt-in push boundary. ADR-0135 requires absent diagnostic output to be reported as an explicit state, not an empty message; it does not define whether a Git reference exists. Keep branch disposition evidence-based. ADR-0150 confirms purpose-based work branches and the existing Roundfix Run namespace. ADR-0118 remains relevant for compatibility of the saved decision key, with its personal-prefix aspect replaced. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- New work branches express purpose; Run/Task branches preserve their separate Roundfix-owned namespace.
- Linked worktrees discover and reconcile the same repository's Run records.
- Existing namespaces and persisted identity mappings remain readable through migration.
- Reconciliation names absent refs truthfully and preserves any content it cannot prove delivered.

## Core Features

1. Apply the confirmed purpose-based policy to new work branches and retain the existing Roundfix Run/Task naming scheme. Keep the compatible branch decision key while replacing its personal default; no Run namespace migration is required.
2. Derive a repository identity shared by main and linked worktrees, with migration or alias evidence for prior checkout-derived records. Moving/removing a linked checkout must not make its Runs unreachable.
3. Retain the already implemented deleted-target content proof and its positive/negative tests; do not duplicate or weaken it.
4. Characterize present-target squash and generic missing-ref cases. Prove delivered content through an approved local evidence rule, or return a truthful preserved/unintegrated result with the missing proof named.
5. Update the public skill examples and recovery output so proposed branch names and commands match what the executor creates.

## Non-Goals / Out of Scope

- No mass renaming or deletion of existing branches/worktrees, force cleanup, or archive-only proof of integration.
- No merging unrelated clones merely because their remote URLs match.
- No automatic repository profile change, branch-policy exception, or new forge dependency for local reconciliation.

## Acceptance evidence

Outside-evidence row: Use the pre-existing Pantheon naming capture and real main/linked-worktree Git fixtures. Preserve a Run-only change in every negative cleanup case, including a Spec whose archive record exists but whose branch content differs.

The later Task Graph and QA must prove each Core Feature with observed
positive and negative cases. Source inspection and proposed checks do not
constitute implementation or terminal QA evidence.

## Open Questions

- The branch-policy decision is settled: the approved subset is recorded in [_authorization-2026-09-08-branch-policy.md](_authorization-2026-09-08-branch-policy.md). No branch-prefix bootstrap exception is needed.
- Approve the repository identity migration boundary and how prior checkout records remain discoverable.
- Choose the acceptable present-target squash content proof; unavailable or ambiguous evidence must preserve work.

Until answered, all proposed limits and protected mutations remain unapproved.

## Source ownership

The maintainer selected this intent for implementation. Ordinary sources now
have one primary owner and one copy under that owner's `references/` directory.
Active Rollups remain as shared archive-license roots; their dated addenda map
every remaining family to its consuming Spec. Adoption is not execution approval.

- [2026-09-08-run-branches-ignore-the-selected-prefix.md](../../history/findings/2026-09-08-run-branches-ignore-the-selected-prefix.md)
- [2026-08-06-rollup-run-lifecycle-and-branch-integrity.md](../../findings/2026-08-06-rollup-run-lifecycle-and-branch-integrity.md)
- [2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md](../0129-spec-authoring-and-gate-recovery/references/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md)

## Research basis

Secondbrain `inbox/roundfix/_triaged/2026-09-08-branches-internas-do-run-ignoram-o-prefixo-do-repositorio.md` supplies the measured naming conflict. The two triaged reconcile captures dated 2026-09-02 are already resolved for their deleted-target case, demonstrated by a fresh focused test. Exa read [Git worktree](https://git-scm.com/docs/git-worktree) and [repository layout](https://git-scm.com/docs/gitrepository-layout): linked worktrees attach to one repository while retaining distinct metadata. This supports investigating shared repository identity; the exact Roundfix migration and proof remain proposed.

The local Secondbrain index was read and its query workflow used before
authoring. The sources above affected the stated requirements and limits;
they do not approve this Spec or substitute for live capability/behavior proof.

## Authoring checkpoint

Spec 0119 owns the future authorization-reader contract; 0121 supplies skill regeneration ownership. The maintainer's confirmed branch-policy correction removes the Run bootstrap blocker. Repository identity and reconciliation work remain proposed, so the approved naming subset does not authorize the whole Spec.

Record the maintainer's bounded decision in [_authorization.md](_authorization.md),
then commit that approval record separately before its consuming tooling
changes. Only after that checkpoint may the TechSpec settle the design and a
Task Graph authorize execution. PRD-stage checks will report their actual
scope and any pending authorization findings; a green partial check is not
implementation readiness.

## Technical candidate

The [_techspec.md](_techspec.md) records the reviewable implementation map,
coverage and build order. It is a proposed candidate, not a completed authoring
gate or permission to dispatch. Exact governed grants and the named decisions
remain pending; no Task Graph or implementation result is claimed.
