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

This Spec is **approved for implementation**. On 2026-09-09 the maintainer
approved the twenty-one bounded governed paths recorded in
[_authorization.md](_authorization.md) and chose committed provenance as the
trust contract for executing authored Verification. The requested outcome, the
through-merge delivery boundary, the record schema, the protected paths, and
the command execution boundary are all settled; what remains is implementation
against that grant.

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
- Tooling authority: applicable — express maintainer authorization: "Aprovar os 21 caminhos", 2026-09-09, recorded in `docs/specs/0119-spec-contained-authorization/_authorization.md`; bounded files: `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/modules/context-workflow.json`, `internal/speccheck/constraints.go`, `internal/speccheck/constraints_characterization_test.go`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `.agents/skills/write-tasks/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`, `skills/write-tasks/SKILL.md`, `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `internal/speccheck/governed.go`, `internal/speccheck/governed_repocontract_test.go`. Sanctioned regeneration follows the approved source edits: `make skills-sync` rewrites the shipped bundle already enumerated above, and `make baseline-digests` rewrites the derived pins. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

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
   requires the configured review policy outcome and passing required checks
   for the candidate commit. Explicit none records intentional review omission;
   it does not grant a release or bypass.
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
  delivery may reach merge only with the configured review policy outcome and
  checks approved.
- Approved on 2026-09-09: one `_authorization.md` per Spec carrying typed
  `status`, `granted`, `action`, `consuming` and `paths` frontmatter with a
  preserved decision trail, over the twenty-one bounded governed paths listed
  in Project Constraints. Spec 0130 already delivered a grant in this shape, so
  the decision makes existing practice canonical rather than inventing a form.
- Approved on 2026-09-09: authored Verification executes on committed
  provenance. Commands run when the Spec artifacts carrying them are tracked in
  this repository and byte-identical to their committed bytes at the resolved
  revision. A Spec Root outside the Git tree, an untracked or modified artifact,
  or command text differing from the committed bytes requires an execution
  approval in that Spec's `_authorization.md` naming the approved revision.
  Read-only `spec check` stays available for every source.
- Resolved on 2026-09-09 from the standing operating rules rather than a new
  question: no paid API use is authorized by this Spec, and corrective work
  stays inside the standing cap of two corrective Tasks; a third stops the Run
  because the decomposition, not the code, is wrong.

### Declared intentional breaks

Everything outside this list must keep behaving exactly as it does today; an
observable change that is not named here is a regression, not a decision.

1. A Spec-contained record whose filename carries no date becomes a validated
   typed grant. A Spec citing such a record while it stays `proposed` or
   `granted: null` starts failing `roundfix spec check` where the date-keyed
   predicate let it pass silently.
2. The Governed Path set gains the owned shipped authoring templates. Changing
   `skills/write-prd/references/prd-template.md` or
   `skills/write-techspec/references/techspec-template.md` without a grant
   starts failing the changed-path audit where it previously passed.
3. `--run-verification`, Implement dispatch and Settle refuse authored commands
   from an untracked or modified Spec artifact, or from a Spec Root outside the
   repository's Git tree, **when no valid execution approval covers that
   source**, where they previously executed them unconditionally. A source
   carrying a valid approval still executes, so the break is the removal of
   unconditional trust, not the removal of the approval path.

Preserved without change: every legacy dated record under
`docs/workflow/authorizations/` that the suite guard still reads, every already
approved Spec-contained grant, archived Spec bytes, Daemon ownership of
Verification and of Task status, and the current public command, error and
persistence contracts.

### Regression locks

- A characterization corpus of today's constraint, governed-path and execution
  behavior is captured before the change and becomes the regression gate, not a
  test written afterwards to match the new code.
- Every new refusal is fail-closed on evidence and never on doubt: a path that
  legitimately passes today must still pass. A refusal fires only on a positive
  observation — an absent or malformed grant, a proven byte difference against
  the committed source — never on an unreadable or merely unrecognized input.

## Open Questions

None. The three questions that survived the 2026-09-08 clarification were
answered on 2026-09-09: the grant convention and its exact protected paths were
approved, the source-trust contract was set to committed provenance with
per-source approval bound to the approved revision rather than a reviewed
queue, and the spending and correction limits resolved to the standing rules
recorded in Decisions.

## Research basis

The inventory read the existing authorization-home Backlog, the Inbox report
about third-party Verification execution, and the current constraints and
changed-path readers. These inputs are adopted in the source index; their ownership does not claim
implementation has started. Secondbrain's earlier fleet workflow observations
and the Exa-read Codex/Claude execution documentation distinguish tool access,
input, and approval; this informs explicit approval state instead of inferred
authority. Those external interfaces do not validate the proposed Roundfix
schema. The [historical research record](https://github.com/marcioaltoe/roundfix/blob/6b8ea48725cbca13974eee0b400b3482202874f6/docs/workflow/2026-09-08-pending-work-plan.md) records the source URLs and their limitations; it is available at the cited Git revision, not as a current planning file.

The [source ownership index](references/_index.md) records the pre-adoption
path, type, primary owner and current owned copy. Secondary consumers link
that copy; lifecycle completion means routing, not verified implementation.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order. Its named decisions are settled and its governed grant is
approved, so it is the contract the Task Graph decomposes. It records no
implementation result; evidence belongs to the Tasks and the terminal QA gate.

## Squash history prerequisite

A new or widened governed grant must land independently in target ancestry
before its consuming squash delivery. A separate commit inside the same PR
would be flattened with the change and does not preserve prior approval.
Validate that boundary before dispatch and against the observed merged history.

## Maintainer-directed cleanup — 2026-09-08

The maintainer explicitly requested deletion of every file under
`docs/workflow/`, including the old authorization records, and the two dated
design documents. Their content remains in Git revision
`6b8ea48725cbca13974eee0b400b3482202874f6`; this Spec must not promise that their
old paths remain readable as current worktree files. Historical approval and
scope remain evidence to resolve from their operative revision, not new grants
recreated or widened by the cleanup. The Spec-contained approval for this Spec
was granted on 2026-09-09 over the exact bounded paths in Project Constraints;
it widens no historical grant, and adapting any reader outside those paths
requires its own bounded scope.
