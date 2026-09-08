---
spec: 0122-verified-content-and-terminal-settlement
status: active
created: 2026-09-08
surfaces: [backend, cli, docs]
---

# Clean proves the content that reached the destination

A Task can verify executable source, omit it from the commit, and leave a Clean Run whose destination fails the same assertion. Task checks also do not establish a general postcondition over committed output. Separately, archive accepts a properly declared-only partial QA result while settlement leaves its QA Task failed. The maintainer needs verification, committed content, and terminal disposition to agree without crediting unobserved acceptance.

This Spec is **in authoring**. Its protected scope and open decisions are
pending approval. The accompanying authorization is a proposal with no grant;
there is no TechSpec, Task Graph, or authority to start implementation.

## Project Constraints

- Identifier strategy: applicable — preserve Run IDs, Task IDs, Spec slugs, and the identity of evidence attached to each settlement; any new diagnostic identity requires an explicit contract. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — this is local Git, Task settlement, and QA evidence; no authentication or HTTP contract changes. Source: `docs/agents/cli.md` and `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the following decisions constrain settlement and its evidence. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0014 is applicable: preserve Daemon-run Verification and final settlement; Agents do not establish completion by narrating success.
  ADR-0020 is applicable: preserve parsed prompt-result precedence over a later nonzero acpx teardown exit, followed by actual Daemon Verification; teardown noise never supplies missing content proof.
  ADR-0038 is applicable: retain the single same-Session Verification repair and bounded failure policy; the proposed postcondition does not authorize retries until success.
  ADR-0056 is applicable: preserve separate Task and Verification Capacity, cancellation-aware acquisition, and the existing single exit-75 retry; this Spec adds no machine-wide scheduler or retry budget.
  ADR-0057 is applicable: the Daemon remains the exclusive writer of Implement Task status and normalizes Agent-authored terminal values before judging the work.
  ADR-0080 is applicable: preserve typed blocked causes and equivalent observed evidence; aligning the declared-only archive case cannot credit failed or unobserved acceptance.
  ADR-0091 is applicable: keep QA as the terminal Task node depending on every leaf, with explicit gate inclusion or decline and the existing graph invalidation rules.
  ADR-0093 is applicable: consistency checks report written declarations and citation gaps; this Spec's settlement policy still requires explicit design and behavioral proof rather than inference from a checker pass.
  ADR-0096 is applicable: retain the Daemon-owned mechanical stage and its machine-fact evidence; no settlement change may make QA verdicts more permissive through that stage.
  ADR-0097 is applicable: carry a QA row only from a prior pass with declared, unchanged repository evidence; blocked or partial acceptance cannot become passed evidence by carry-forward.
  ADR-0104 is applicable: use the pre-existing Pantheon omission and declared-partial observations as outside acceptance evidence, preserving provenance and explicit blocked status when evidence cannot be obtained.
  ADR-0117 is applicable: place artifact checks at the stage producing their defect; committed-content checks need the commit evidence that pre-work authoring cannot supply, and public behavior still belongs to QA.
  ADR-0127 is not applicable to this change: machine process-residue inventory remains a readiness fact and this Spec introduces no residue command or synthetic Run record.
  ADR-0138 is applicable: preserve one commit per verified Task and opt-in push only at Clean; the proposed postcondition supplies additional completion evidence without granting new delivery actions.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `internal/spec/archive.go`, `internal/spec/archive_test.go`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/qa-gate/SKILL.md`, `skills/archive-spec/SKILL.md`, `internal/baseline/assets/modules/spec-workflow.json`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `skills/write-tasks/SKILL.md`, `skills/write-tasks/references/task-template.md`, `.agents/skills/implement-task/SKILL.md`, `skills/implement-task/SKILL.md`, `internal/speccheck/coherence.go`, `docs/agents/spec-routing.md`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A required output omitted from a commit prevents a misleading Clean outcome and remains recoverable.
- Intentional executable repository source can be delivered without treating its permission bit as build-artifact proof.
- The committed and integrated candidate satisfies the approved repository postcondition.
- The same declared-only QA evidence yields compatible settlement and archive decisions.

## Core Features

1. Distinguish intentional repository source from disposable output using evidence beyond executable mode. At minimum, preserve changes to already tracked executable source.
2. When required output cannot be staged or committed, report its paths and cause in the Task/Run outcome and preserve the recovery surface; a console warning alone cannot establish completion.
3. Proposed: run the selected repository Verification against the integrated candidate before declaring Clean, with Daemon ownership and an explicit failure/recovery result. The design must decide its exact commit boundary and avoid redundant full gates per Task.
4. Use one declared-acceptance eligibility policy for the newest QA Report across settlement, the derived QA command, and archive. A qualifying partial report retains its partial verdict and records unproven actions; failed, missing, undeclared, or unobserved acceptance stays blocking.
5. Make the approved settlement semantics consistent in the QA, archive, and public command guidance.

6. An explicitly authorized Task may enter to repair its named known-red repository precondition under frozen Spec/Task/source authority. Normal Tasks remain blocked by red prerequisites. The same required gate and focused repair assertion must pass before settlement; shell rephrasing, an Agent-edited policy or an unapproved failure cannot bypass entry checks.
7. Collect failures from explicitly independent Verification groups before spending the existing single repair turn, while preserving dependency ordering and skipping checks whose setup prerequisite failed. Do not infer independence from arbitrary shell text or simply continue every command after failure. Preserve all diagnostics and the existing retry ceiling.

## Non-Goals / Out of Scope

- No blanket acceptance of every executable or untracked file, no dropped-path suppression, and no deleting evidence to obtain Clean.
- No relaxation of failed rows, fabricated equivalent evidence, or automatic acceptance of unknown risk.
- No change to release/deployment authority or permission for QA Agents to commit and push.

## Acceptance evidence

Outside-evidence row: Replay the Pantheon tracked-executable omission and the existing declared-only partial report shape from the Inbox/Findings. The candidate must preserve the executable change, reject omitted required work, and never turn an ineligible partial into success.

The later Task Graph and QA must prove each Core Feature with observed
positive and negative cases. Source inspection and proposed checks do not
constitute implementation or terminal QA evidence.

## Open Questions

- Approve the proposed postcondition boundary and failure disposition before the TechSpec fixes transaction order.
- Approve the declared-only partial settlement policy with explicit unproven actions, or keep it blocked and revise archive to match.
- Choose the evidence required for newly created executable source; already tracked source is the measured defect.

Until answered, all proposed limits and protected mutations remain unapproved.

## Source ownership

The maintainer selected this intent for implementation. Ordinary sources now
have one primary owner and one copy under that owner's `references/` directory.
Active Rollups remain as shared archive-license roots; their dated addenda map
every remaining family to its consuming Spec. Adoption is not execution approval.

- [2026-09-08-verified-executable-source-is-dropped-before-settlement.md](references/2026-09-08-verified-executable-source-is-dropped-before-settlement.md)
- [2026-08-06-rollup-qa-gates-and-verification-evidence.md](../../findings/2026-08-06-rollup-qa-gates-and-verification-evidence.md)
- [2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md](../0129-spec-authoring-and-gate-recovery/references/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md)

## Research basis

Secondbrain `inbox/roundfix/_triaged/2026-09-02-commit-boundary-descarta-arquivo-executavel-e-o-run-reporta-clean.md` supplies the original mismatch between verified and delivered content; `wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md` requires the checker to observe the property it claims. Exa read [Git diff documentation](https://git-scm.com/docs/git-diff), which distinguishes comparisons of working tree, index, and committed trees. This informed separate evidence for each boundary, not a claim that working-tree verification proves the commit.

The local Secondbrain index was read and its query workflow used before
authoring. The sources above affected the stated requirements and limits;
they do not approve this Spec or substitute for live capability/behavior proof.

## Authoring checkpoint

Spec 0119 supplies the approval contract; 0121 supplies skill regeneration ownership. This integrity work is a prerequisite to unattended delivery claims in Spec 0127.

Record the maintainer's bounded decision in [_authorization.md](_authorization.md),
then commit that approval record separately before its consuming tooling
changes. Only after that checkpoint may the TechSpec settle the design and a
Task Graph authorize execution. PRD-stage checks will report their actual
scope and any pending authorization findings; a green partial check is not
implementation readiness.

The [source ownership index](references/_index.md) records the pre-adoption
path, type, primary owner and current owned copy. Secondary consumers link
that copy; lifecycle completion means routing, not verified implementation.

## Technical candidate

The [_techspec.md](_techspec.md) records the reviewable implementation map,
coverage and build order. It is a proposed candidate, not a completed authoring
gate or permission to dispatch. Exact governed grants and the named decisions
remain pending; no Task Graph or implementation result is claimed.
