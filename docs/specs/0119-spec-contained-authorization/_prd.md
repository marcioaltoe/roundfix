---
spec: 0119-spec-contained-authorization
status: active
created: 2026-09-08
surfaces: [backend, cli, docs]
---

# A Spec carries its authority

The maintainer wants every pending Spec to carry the authorizations needed to
implement and deliver it. Today grants may live outside the Spec, historical
discovery assumes one workflow directory, and executing authored Verification
does not have a clearly stated source-trust boundary. A maintainer should be
able to inspect one Spec and distinguish approved actions from proposed work.

This Spec is **in authoring**. The requested outcome and through-merge delivery
boundary are confirmed. Its proposed schema, protected paths, and command
execution boundary await the maintainer's decision; there is no executable
Task Graph and no grant to change protected tooling yet.

## Project Constraints

- Identifier strategy: applicable — preserve Spec slugs, Task IDs, and existing diagnostic identities; authorization belongs to one consuming Spec rather than a new global identity service. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no authentication system or HTTP endpoint is introduced. Existing credentials remain in their current runtime boundary. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0014 keeps Verification with the Daemon. Source: `docs/agents/spec-routing.md`.
  ADR-0130 keeps the audit on governed paths and preserves the historical governed set; this Spec changes grant discovery without exempting those paths.
  ADR-0149 makes the grant name the regeneration command and the ownership tree name its outputs; preserve that division.
  ADR-0057 keeps the Daemon the exclusive writer of Implement Task status; authorizations do not let Agents settle themselves.
  ADR-0096 requires mechanical facts, including bounded-path authority, before the QA Agent turn; extend that check rather than bypassing it.
  ADR-0117 places a defect check at the stage that can produce it; validate authored grants at authoring and actual changed files after commits exist.
  ADR-0020 is not applicable to the implementation scope: parsed prompt result versus acpx teardown exit classification remains unchanged.
  ADR-0038 is not applicable to the implementation scope: the one Verification repair allowance remains unchanged and is not widened by a grant.
  ADR-0056 is not applicable to the implementation scope: Task Capacity, Verification Capacity, and the temporary-failure retry remain unchanged.
  ADR-0127 is not applicable to the implementation scope: reporting process residue as a readiness fact remains unchanged.
- Tooling authority: applicable — protected tooling mutation is proposed, not authorized. The reviewable proposal is `docs/specs/0119-spec-contained-authorization/_authorization.md`; bounded files: `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/modules/context-workflow.json`, `internal/speccheck/constraints.go`, `internal/speccheck/constraints_characterization_test.go`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `.agents/skills/write-tasks/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`, `skills/write-tasks/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`. Source: `docs/agents/agent-instructions.md`.

## Goals

- A Spec remains self-contained when it moves into history: its approval,
  scope, limits, and consuming work travel with it.
- A reader and the execution machinery reach the same answer about which
  actions are approved. A proposed, absent, contradictory, or withdrawn
  authorization grants nothing.
- Running authored commands uses the approved source and execution scope;
  importing a third-party Spec does not silently authorize its shell commands.
- Existing granted work remains readable and valid during migration.

## User Stories

1. As the maintainer, I want to approve exact actions and protected paths
   inside the Spec, so that implementation needs no repeated permission for
   work already covered by that approval.
2. As the Supervisor, I want to identify missing authority before dispatch,
   so that unattended work stops at a concrete boundary and preserves evidence.
3. As a reviewer, I want to compare actual changed files and publication
   actions with the consuming grant, so that generic approval cannot conceal
   a tooling or deployment expansion.

## Core Features

1. Every new Spec records the approval state, maintainer decision and date,
   permitted actions, exact protected paths, consuming Spec, sanctioned
   regeneration, and relevant limits. Its PRD and TechSpec both point to the
   same operative record in Project Constraints.
2. The proposed record is separate from the commit that consumes it. An
   executor cannot widen its own grant to make its change pass. Amendments
   require a new recorded maintainer decision before dependent work resumes.
3. The canonical authoring rules and templates require this placement. All
   preserved legacy grants continue to resolve without rewriting their
   historical text or granting new actions.
4. Grant validation depends on the artifact's role, not the presence of a
   date in its filename. Missing approval and malformed or escaping paths
   are actionable diagnostics, never successful validation.
5. Approval for implementation, commit, push, PR creation, merge, and release
   remains distinguishable. The session's approved through-merge policy
   requires passing independent review and required checks for the candidate
   commit; it does not grant a release or bypass.
6. Execution of authored Verification states which Spec source and commands
   were approved and what effects they may have. Read-only checking remains
   available without command execution. A changed or untrusted source needs
   the appropriate decision before shell execution; no new sandbox or
   credential policy is implied.

## User Experience

The maintainer reviews the complete grant in the Spec and answers one pending
decision at a time. A refusal names the missing action, protected path,
changed source, or unresolved limit. Once a bounded action is approved, the
Supervisor proceeds without asking for that same approval again.

## Non-Goals / Out of Scope

- Rewriting archived Specs or migrating every historical authorization.
- A signature service, global permission database, or new authentication layer.
- Authorizing arbitrary shell from imported documentation, broad tooling
  classes, production deployments, paid model calls, or automatic risk waivers.
- Changing Task status ownership or replacing the Daemon's Verification.

## Success Metrics

- Positive and negative checks distinguish a valid consuming grant from
  proposals, wrong-Spec grants, missing fields, path traversal, and scope
  expansion without a new approval.
- An archived new Spec retains an operative approval reference; an existing
  legacy grant resolves unchanged.
- A third-party Spec can be inspected without executing its commands.

## Decisions

- Confirmed on 2026-09-08: authorizations must be inside Specs and this must
  become canonical guidance.
- Confirmed on 2026-09-08: after approval of Specs and limits, autonomous
  delivery may reach merge only with independent review and checks approved.
- Proposed: one `_authorization.md` per Spec, with explicit proposed/approved
  state and a preserved decision trail. This is not an approved schema yet.

## Open Questions

- Approve the proposed grant convention and exact protected paths.
- Set the source-trust contract for executing imported Verification and
  identify whether bounded approval is per Spec revision or a reviewed queue.
- Set model-spending and correction limits for consuming Specs. Until then,
  no paid probe or new execution authority is implied.

## Research basis

The inventory read the existing authorization-home Backlog, the Inbox report
about third-party Verification execution, and the current constraints and
changed-path readers. These are inputs awaiting adoption, not claims that
implementation has started. Secondbrain's earlier fleet workflow observations
and the Exa-read Codex/Claude execution documentation distinguish tool access,
input, and approval; this informs explicit approval state instead of inferred
authority. Those external interfaces do not validate the proposed Roundfix
schema. The cross-Spec plan records the source URLs and their limitations.
