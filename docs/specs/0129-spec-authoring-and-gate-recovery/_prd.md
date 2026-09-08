---
spec: 0129-spec-authoring-and-gate-recovery
status: active
created: 2026-09-08
surfaces: [backend, cli, docs]
---

# A Spec exposes its executable promises and recovery path

Historical queue Findings identify promises and API contracts with no Task owner, unsupported recovery after QA becomes stale, and authoring assumptions that CI or later prerequisites invalidate. Give the existing authoring/checking workflow explicit coverage, semantic review and a supported recovery contract before the durable delivery owner consumes it.

This Spec is in authoring. Its implementation intent is selected by the maintainer's
complete-queue request; the exact governed mutations and execution limits remain
proposed until their operative grant is approved.

## Project Constraints

- Identifier strategy: applicable — preserve Spec/Task identities and exact Git ref identity; any new recovery record follows the existing Run identity boundary. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential or HTTP API change is proposed; use existing read boundaries only. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014 and ADR-0057 reserve Verification and Task status to the Daemon; ADR-0091 keeps QA terminal; ADR-0117 places each detector at the stage that can establish it; ADR-0148 preserves non-vacuous Verification probing. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0080 applies: preserve the distinction between environment-blocked QA rows and product failure or success.
  ADR-0020 applies at the inherited Agent result boundary: a parsed prompt result remains usable despite a later teardown exit, with the anomaly recorded before Verification.
  ADR-0038 applies: the same Agent Session has only the accepted single Verification Feedback repair; recovery cannot enlarge it.
  ADR-0096 applies: prove mechanical QA facts before spending an Agent turn and preserve the existing verdict semantics.
  ADR-0127 applies: process residue remains a readiness observation, not a synthetic Run or authority to settle work.
  ADR-0056 applies: authoring/recovery preserves distinct Task and Verification capacities; neither grants whole-queue concurrency.
  ADR-0093 applies: consistency uses explicit citations and does not infer behavior from their presence.
  ADR-0097 applies: carry forward QA only from declared, unchanged evidence; stale input requires supported revalidation.
  ADR-0104 applies: acceptance uses independent evidence with its actual origin; missing external evidence remains visible under the declared policy.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `.agents/skills/implement-task/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-techspec/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-tasks/references/task-template.md`, `skills/implement-task/SKILL.md`, `internal/speccheck/coherence.go`, `internal/baseline/assets/modules/spec-workflow.json`, `docs/agents/spec-routing.md`, `docs/agents/setup-context.json`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- Every promised behavior, API contract and success metric has a Task and evidence owner.
- A falsified premise or stale QA produces a supported recovery action that preserves prior evidence.
- Task dependencies and verification declarations are executable in the declared CI/toolchain context.

## Core Features

1. Extend coverage from PRD units to API contracts, success metrics and accepted-ADR obligations, with Task and evidence references. Semantic independent review checks the decision and rejected alternatives; citation presence alone is not obedience.
2. Record premise falsification and supersession as an amendment with preserved prior Results and evidence. Revalidate consumers after their prerequisites land.
3. Provide a supported stale-QA invalidation/recovery path with Daemon-owned state changes and retained earlier report/Result evidence. Never delete the prior QA report or hand-edit a completed status to make a stale gate disappear.
4. Represent temporal prerequisites and ordinal-generator dependencies explicitly. A requirement needing a future release or external observation remains a named prerequisite and never invents release authority.
5. Require property-shaped acceptance, appropriate use-case/persistence test seams, updates to the affected existing contract, narrow Spec commits and declared CI feasibility. Validate newly required test classes against the actual CI command.
6. Keep the existing Task Type, terminal QA and bounded corrective-Task contracts. Split or reauthor when the accepted corrective ceiling is reached.

## Non-Goals / Out of Scope

- Replacing the Daemon with an authoring Agent or starting the durable delivery owner, which belongs to 0127.
- Automatically approving a changed Spec, granting a release or erasing historical Results.
- Treating syntactic citation checks as a complete semantic reviewer.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
later Task Graph. At least one acceptance row replays the pre-existing fleet
observation identified by the source ownership index; a fixture authored by this
Spec alone is insufficient to validate the original premise.

## Decisions and open authority

The maintainer selected the remaining source intent for the implementation queue.
The technical candidate is reviewable; proposed paths, experimental limits and
any change to an accepted ADR must receive an operative decision before execution.
The existing through-merge authority requires independent review and required
checks and grants no release or waiver.

## Research basis

Local source inspection and the archived fleet observations establish the
remaining behavior. Secondbrain's agent-workflow and harness syntheses explain
why durable state, external acceptance and explicit stop conditions matter.
Exa-read official Codex noninteractive and Claude headless documentation supports
independent execution interfaces, not the correctness of this Roundfix design.
The queue plan records source URLs, their influence and remaining limitations.

The [source ownership index](references/_index.md) records the pre-adoption
path, type, primary owner and current owned copy. Secondary consumers link
that copy; lifecycle completion means routing, not verified implementation.

## Technical candidate

The [_techspec.md](_techspec.md) records the reviewable implementation map,
coverage and build order. It is a proposed candidate, not a completed authoring
gate or permission to dispatch. Exact governed grants and the named decisions
remain pending; no Task Graph or implementation result is claimed.
