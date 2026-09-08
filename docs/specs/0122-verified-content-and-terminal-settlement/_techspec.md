---
spec: 0122-verified-content-and-terminal-settlement
prd: _prd.md
created: 2026-09-08
---

# Verified content, repair entry and truthful settlement — Technical candidate

## Executive Summary

Make the Daemon settle only the content and verdict it actually proved. Reuse existing Task execution and Git integration while moving omission checks before completion and sharing QA eligibility across every consumer. The trade-off is a typed recoverable refusal where the old path emitted a warning and reported success.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The cross-Spec order and remaining
decisions are in [the portfolio plan](../../workflow/2026-09-08-pending-work-plan.md).

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

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Repair entry policy | `internal/spec/spec.go and task.go; internal/daemon/task_engine.go` | Freeze an approved repair contract before starting an Agent on a known-red prerequisite. |
| Independent Verification groups | `Task parsing, authoring consistency and daemon Verification loop` | Aggregate independent failures without executing checks whose setup dependencies failed. |
| Delivered-path proof | `FilterStageablePaths, commitTask and shared daemon callers` | Preserve tracked executable source and block omission of required output before completion. |
| QA eligibility | `internal/spec eligibility reader; archive, daemon, Implement and derived QA command` | Apply one declared-acceptance policy while preserving the original verdict and unproven actions. |
| Repository postcondition | `Implement candidate verification and worktree integration` | Verify the exact committed candidate before target movement, cleanup or push. |
| Settlement guidance | `owned Task/QA/archive/Roundfix skills and canonical spec module` | Keep the same entry, verification and recovery contracts visible at each surface. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Entry and feedback

Replace command-string equality as the authority switch with a typed committed
entry policy. A repair Task names its Spec/Task, approved source, bounded repair
scope, selected repository command identity and prior failure evidence. Execute
and journal the initial precondition. Normal Tasks remain blocked by a red
gate. The repair allowance permits entry for that known failure only; unknown
execution failures, wrong/stale/proposed grants and an Agent-edited policy
refuse. The selected gate and focused effect assertion must pass at settlement.
Keep the existing repair-turn ceiling and do not introduce skip_precondition.

Independent Verification groups require explicit dependencies. Preserve ordered
legacy command lists; never infer arbitrary shell independence. A setup group
must succeed before its dependent checks run. Independent sibling groups each
report their failure before the single repair feedback is assembled, and a
blocked dependency is distinct from a passing or failed check. The Agent cannot
weaken group edges or command identity while reloading its Task. Authoring
validation rejects cycles and missing dependency names.

### Stage, commit and eligibility

Use index/tracked-source evidence, not executable mode alone, to classify source.
At minimum preserve the captured tracked executable and its mode. A new
executable needs explicit declared source evidence and path bounds; an arbitrary
untracked binary, escaping symlink or generated cache remains excluded. A
required path omitted from the candidate prevents completed/Clean settlement
and retains the surface with a typed reason. Check shared commit callers, QA
paths and recovery rather than changing one filter call.

Put QA eligibility in the Spec package. Settlement, Implement Clean, archive and
the derived read-only QA command consume the same result: original verdict,
eligibility and unproven actions. Preserve the accepted declared-only partial
policy without turning partial into pass. Validate actual declared acceptance
relationships, not only equal counts; fail/precondition failures and environment
rows lacking the accepted equivalent evidence remain ineligible.

### Candidate postcondition

After the Task Graph and its authored QA settle, execute the selected complete
repository Verification once against the exact committed Run candidate before
integrating or cleaning it. Integration must establish the destination equals
that verified candidate through the existing fast-forward contract. A failure,
cancellation, changed candidate or unrelated target dirt retains recovery and
prevents Clean/push. This is the concrete interpretation of the PRD's integrated
candidate: prove the commit before moving it and prove the destination after,
rather than running the gate on unrelated user edits after destructive cleanup.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Record the exact authorized repair precondition, independent group dependencies, verified path/content evidence and final committed repository postcondition. Preserve earlier failures and Daemon-owned terminal state.

### API Contracts

Keep Implement/Daemon settlement public outcomes truthful. The proposed repair-entry and independent-group declarations require explicit validation; no free-form skip flag or manual terminal-state edit substitutes for them.

## Coverage Map

- PRD Goal 1 → Delivered-path proof, Repository postcondition.
- PRD Goal 2 → Delivered-path proof.
- PRD Goal 3 → Repository postcondition.
- PRD Goal 4 → QA eligibility.
- Core Feature 1 → Delivered-path proof.
- Core Feature 2 → Delivered-path proof.
- Core Feature 3 → Repository postcondition.
- Core Feature 4 → QA eligibility.
- Core Feature 5 → Settlement guidance.
- Core Feature 6 → Repair entry policy.
- Core Feature 7 → Independent Verification groups.

The Testing Approach below describes the observations that must settle these
contracts. Task IDs and actual evidence are deliberately not invented during
proposal authoring; approved Task decomposition must assign every contract and
success metric before execution.

## Integration Points

Local repository evidence and the adopted sources define the concrete seams.
The [owned source index](references/_index.md) records each primary source.
The [portfolio plan](../../workflow/2026-09-08-pending-work-plan.md) records
secondary consumers and prerequisite Specs. External research was read through
Exa and compared with local Secondbrain history; the PRD and portfolio plan
retain links and describe its effect. Published interfaces support feasibility,
not a claim that the proposed runtime or behavior already exists.

## Testing Approach

Use focused tests at the named package seams for deterministic rules, real
Git/store/process boundaries for integration behavior, and the authored public
QA Task for user-visible acceptance. Do not infer a terminal pass from source
inspection or a focused fixture. Required observations:

1. Normal red precondition starts no Agent; approved exact repair can enter and must make the same gate pass; stale/wrong/Agent-edited authority and shell rephrasing cannot bypass it.
2. Two independent failures reach the single feedback turn; a failed setup blocks dependent commands; cyclic/undeclared groups refuse before execution.
3. Tracked executable content/mode, intentional deletion/rename, required omitted source, untracked binary, symlink escape and target dirt have distinct truthful outcomes.
4. One pass/declared-partial/fail/precondition/missing/newest-report matrix agrees across every QA consumer while retaining unproven actions.
5. Postcondition failure, cancellation and candidate change prevent integration cleanup/push; successful integration equals the verified commit.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Typed authorized red-precondition entry (depends on: none).
2. Explicit independent Verification groups and complete failure feedback (depends on: 1).
3. Delivered-path proof before completed settlement (depends on: 1).
4. Shared QA eligibility and derived public consumer (depends on: 1).
5. Committed-candidate repository postcondition (depends on: 2, 3, 4).
6. Canonical and owned-skill settlement guidance (depends on: 1, 2, 3, 4, 5).
7. Real Git/public CLI recovery and terminal QA (depends on: 1, 2, 3, 4, 5, 6).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Repair entry is authority, not a waiver: its policy must be frozen before the Agent turn and checked again at settlement. Aggregating arbitrary dependent shell commands is unsafe, so legacy sequencing stays intact. Preserve test names where practical to avoid unrelated coverage-record changes.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge is confirmed only with independent review and required checks approved for the current candidate; releases, tags and paid consumption are not implied.
- Reviewer selection follows the [confirmed portfolio policy](../../workflow/2026-09-08-pending-work-plan.md); this Spec does not introduce a separate override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
