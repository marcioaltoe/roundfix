---
spec: 0154-evidence-when-the-target-branch-is-gone
status: active
created: 2026-09-23
surfaces: [backend, cli, docs]
---

# Evidence when the target branch is gone

## Executive Summary

When a Run's target branch no longer exists, reconciliation seeks the same
content evidence it already accepts for a missed ancestry, against the default
branch. A Run is released only on positive evidence; everything else stays
preserved with the missing proof named.

## Project Constraints

- Identifier strategy: applicable — Runs, Branches and Worktree paths keep their
  identities. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git reads only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0062, ADR-0093, ADR-0104, ADR-0130,
  ADR-0155 and ADR-0156 all hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — `reconcile` reporting is public CLI behavior
  the skill documents. Express maintainer authorization: recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Where the answer is lost today

Three paths decide whether a terminal Run is releasable:

- `internal/worktree/worktree.go:422` returns early when the target branch is
  absent, with `reconciliationReasonTargetBranchAbsent`. Nothing is checked.
- `:945` asks `merge-base --is-ancestor` when the target branch is present.
- The miss branch immediately after it calls `supersedingQAReport`, which proves
  delivery by content: the target carries a newer QA Report for this Spec and
  differs from the Run only by QA Report commits.

The third is an accepted content-evidence rule. It never runs when the branch is
gone, which is the state every squash-merged pull request leaves behind.

## The fallback

When the target branch is absent, reconciliation resolves the default branch and
applies the same content-evidence rule against it. If the rule proves delivery,
the Run is classified as it would have been with the target present. If it does
not, the Run is preserved.

Reusing `supersedingQAReport` rather than writing a second rule is the point: a
second way to prove integration is a second thing that can disagree about what
integration means.

## What still preserves

- The default branch cannot be resolved.
- The evidence rule finds nothing for this Run's Spec.
- The evidence names a different Spec.
- The Worktree holds uncommitted changes.

Each keeps its own reason, and each reason names the proof that was missing.
`target branch "<name>" is absent` alone is no longer a terminal answer: it is
the condition that sends the check elsewhere.

## API Contracts

1. A Run whose Spec's work reached the default branch, with its target branch
   deleted, is classified releasable and names the evidence.
2. A Run without that evidence stays preserved, and its reason names what was
   missing.
3. Every classification for a present target branch is unchanged.

## Coverage Map

- Goal 1 → The fallback; API Contract 1.
- Goal 2 → What still preserves; API Contract 2.
- Goal 3 → What still preserves.
- Core Feature 1 → Where the answer is lost today; Testing Approach 1.
- Core Feature 2 → The fallback.
- Core Feature 3 → What still preserves; API Contract 2.
- Core Feature 4 → What still preserves.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 4.
- API Contracts 1-3 → The fallback, What still preserves.

## Integration Points

- **`supersedingQAReport`.** Reused, not modified.
- **`reconciliationReasonTargetBranchAbsent`.** Becomes one outcome among
  several rather than the only one.
- **Spec 0125.** Repository identity, branch-naming policy and the skill's
  branch examples stay there.

## Testing Approach

1. **The squash case is characterized.** A fixture whose Run Branch was
   squash-merged is shown not to be an ancestor of the default branch, so the
   test states why ancestry cannot answer.
2. **Absent target, evidence present.** That fixture with its target branch
   deleted is classified releasable, naming the evidence. Fails on the tree as
   it stands, where the absent target returns before any check.
3. **Absent target, evidence missing.** The same fixture without the Spec's work
   in the default branch stays preserved, and the reason names the missing
   proof. Covers an unreachable default branch and evidence naming another Spec.
4. **Present target unchanged.** Safe, superseded, unintegrated and dirty
   classifications are asserted unchanged when the target branch exists.
5. **The skill is true.** The shipped skill and its mirror describe the fallback
   and what still preserves, and the mirror matches the canonical file.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The default-branch fallback and its evidence call, with unit tests for the
   present and missing cases (depends on: none).
2. The preserved reasons that name the missing proof (depends on: 1).
3. The shipped skill (depends on: 1, 2).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **A fallback that releases on absence.** Clearing 22 Worktrees is the visible
  reward, and a rule that released without evidence would clear them all —
  including any holding work nobody else has. Testing Approach 3 is the control,
  and it is a Success Metric rather than a nicety.
- **Two rules drifting.** Reusing `supersedingQAReport` avoids it; writing a
  default-branch variant of the same idea would reintroduce it.
- **A default branch that is itself behind.** The evidence proves the Spec's
  work is present in the resolved default branch, not that the branch is current
  with a remote. A stale local default branch yields preserved, which is the
  safe direction.
