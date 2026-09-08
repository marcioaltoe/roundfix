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
approval to mutate protected files. The remaining decisions are recorded in
[_prd.md](_prd.md) and [_authorization.md](_authorization.md); the dependencies
below define this Spec's place in the implementation order.

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
  ADR-0154 applies: an explicit user-authorized QA Archive Override changes archive eligibility only, preserving original QA/Task evidence and separate delivery gates.
  ADR-0138 is applicable: preserve one commit per verified Task and opt-in push only at Clean; the proposed postcondition supplies additional completion evidence without granting new delivery actions.
- Tooling authority: applicable — express maintainer authorization covers the QA Archive Override policy through [the narrow grant](references/2026-09-08-authorized-qa-archive-override.md); bounded files: `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/autonomous-work.json`, `docs/agents/docs-layout.md`, `docs/agents/skill-dispatch.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, with sanctioned digest regeneration. All Go, tests, skills and remaining implementation mutations stay proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `internal/spec/archive.go`, `internal/spec/archive_test.go`, `internal/cli/archive.go`, `internal/cli/archive_test.go`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/qa-gate/SKILL.md`, `skills/archive-spec/SKILL.md`, `internal/baseline/assets/modules/spec-workflow.json`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `skills/write-tasks/SKILL.md`, `skills/write-tasks/references/task-template.md`, `.agents/skills/implement-task/SKILL.md`, `skills/implement-task/SKILL.md`, `internal/speccheck/coherence.go`, `docs/agents/spec-routing.md`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Repair entry policy | `internal/spec/spec.go and task.go; internal/daemon/task_engine.go` | Freeze an approved repair contract before starting an Agent on a known-red prerequisite. |
| Independent Verification groups | `Task parsing, authoring consistency and daemon Verification loop` | Aggregate independent failures without executing checks whose setup dependencies failed. |
| Delivered-path proof | `FilterStageablePaths, commitTask and shared daemon callers` | Preserve tracked executable source and block omission of required output before completion. |
| QA eligibility | `internal/spec eligibility reader; archive, daemon, Implement and derived QA command` | Share evidence classification, then apply action-specific eligibility; preserve verdict/unproven actions and keep authorized archive override separate from settlement success. |
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

Put QA evidence classification in the Spec package. Settlement, Implement
Clean, archive and the derived read-only QA command share original verdict and
unproven actions; each applies its own eligibility. Archive additionally accepts
a valid QA Archive Override, which cannot satisfy settlement or Clean. Preserve the accepted declared-only partial
policy without turning partial into pass. Validate actual declared acceptance
relationships, not only equal counts; fail/precondition failures and environment
rows lacking the accepted equivalent evidence remain ineligible.

### Authorized QA archive override

Extend the Archive Command and Spec archive boundary with an explicit request
that consumes user authority for the identified Spec/revision and unmet QA
prerequisite. Define final CLI syntax and approval record schema in authoring;
the current command accepts only a slug and the existing ArchiveRequest has
no override input. Merely adding a flag or reading a pre-existing true marker
must not fabricate the required authority.

All non-QA Tasks remain completed and the source reference/index, destination
and execution-ownership checks remain active. A covered terminal QA Task may
be incomplete or failed for archive purposes only; retain its state and report
instead of setting completed/pass to satisfy the old all-Tasks check. An
otherwise invalid graph is not silently repaired by this exception. A known
stale-QA diagnostic must retain its original evidence and be distinguished from
unrelated graph corruption before the overridden archive proceeds.

Stamp `qa_override: true` and a durable approval/evidence reference. Record
actual QA verdict/state or absence, scope/revision, date/source and supplied
reason. Exact additional field names remain to be authored; the boolean marker
alone is insufficient provenance. Report archive-with-override as such. A
normal successful/qualifying partial archive never receives a synthetic override,
and the authored declined-QA contract remains distinct.

An already recorded applicable user approval can be consumed without asking
again. Wrong-Spec, out-of-scope, changed-source or proposed records do not cover
the action. This is archive authority only: preserve Run outcomes, QA verdicts,
Task statuses, review-policy results, external checks and separate publication
permissions. No archive request, including an override, starts a reviewer or
pushes by itself. User approval of this capability is not an override of any
particular queued Spec.

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
- PRD Goal 5 → Authorized QA archive override and Settlement guidance.
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

- Core Features 8 and 9 → Authorized QA archive override, QA eligibility and public archive outcome.

## Integration Points

Local repository evidence and the adopted sources define the concrete seams.
The [owned source index](references/_index.md) records each primary source.
The prerequisite Specs are listed below. Secondary consumers reference the
primary owner's adopted source instead of duplicating it. External research was read through
Exa and compared with local Secondbrain history. The
[historical research record](https://github.com/marcioaltoe/roundfix/blob/6b8ea48725cbca13974eee0b400b3482202874f6/docs/workflow/2026-09-08-pending-work-plan.md)
retains the consulted sources, their influence and limitations after the plan
was removed from the current tree. Published interfaces support feasibility,
not a claim that the proposed runtime or behavior already exists.

## Testing Approach

Use focused tests at the named package seams for deterministic rules, real
Git/store/process boundaries for integration behavior, and the authored public
QA Task for user-visible acceptance. Do not infer a terminal pass from source
inspection or a focused fixture. Required observations:

1. Normal red precondition starts no Agent; approved exact repair can enter and must make the same gate pass; stale/wrong/Agent-edited authority and shell rephrasing cannot bypass it.
2. Two independent failures reach the single feedback turn; a failed setup blocks dependent commands; cyclic/undeclared groups refuse before execution.
3. Tracked executable content/mode, intentional deletion/rename, required omitted source, untracked binary, symlink escape and target dirt have distinct truthful outcomes.
4. One shared evidence-classification matrix covers pass/declared-partial/fail/precondition/missing/newest-report across QA consumers. Verify action-specific archive override never becomes settlement success or discards unproven actions.
5. Postcondition failure, cancellation and candidate change prevent integration cleanup/push; successful integration equals the verified commit.

Archive override controls additionally cover missing, failed, incomplete,
stale and malformed QA; valid explicit/current authority versus absent, wrong
or changed-source authority; unchanged QA Task/report bytes; incomplete non-QA
Task refusal; invalid references and occupied destination refusal; idempotent
inspection of an overridden archive; and no change to Run Clean or independent
delivery eligibility. Confirm a prior covered approval is consumed without
reconfirmation, while mere capability authorization supplies no per-Spec grant.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Typed authorized red-precondition entry (depends on: none).
2. Explicit independent Verification groups and complete failure feedback (depends on: 1).
3. Delivered-path proof before completed settlement (depends on: 1).
4. Shared QA evidence classification, action-specific eligibility and derived public consumer (depends on: 1).
5. Authorized archive-only QA override and durable provenance (depends on: 4).
6. Committed-candidate repository postcondition (depends on: 2, 3, 4).
7. Canonical and owned-skill settlement guidance (depends on: 1, 2, 3, 4, 5, 6).
8. Real Git/public CLI recovery and terminal QA (depends on: 1, 2, 3, 4, 5, 6, 7).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Repair entry is authority, not a waiver: its policy must be frozen before the Agent turn and checked again at settlement. Aggregating arbitrary dependent shell commands is unsafe, so legacy sequencing stays intact. Preserve test names where practical to avoid unrelated coverage-record changes.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge requires the configured pre-PR review policy outcome and passing required checks for the current candidate. Explicit none records intentional review omission; enabled-provider failure cannot select none. Releases, tags and paid consumption are not implied.
- Preserve configured reviewer selection; this Spec introduces no separate reviewer override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.

## Cross-Spec dependencies

Required predecessor contracts: [0119](../0119-spec-contained-authorization/_techspec.md).
Shared skills and canonical files require serial integration and revalidation
after predecessor changes. A predecessor reference is not an execution grant.


## QA override policy evidence — 2026-09-08

The maintainer explicitly requested the capability to archive with QA override
when requested or authorized. The existing archive-spec skill already permits
an explicit archive-anyway exception and stamps `qa_override: true`, but its
all-Tasks-completed prerequisite does not distinguish the terminal QA Task.
The old canonical clause allowed only passing QA. `internal/spec/archive.go`
requires all Tasks completed and has no override field in ArchiveRequest;
`internal/cli/archive.go` accepts only the slug.

Read-only CLI inspection confirmed the current help exposes no override. A
parser-only invocation with the nonexistent slug `__qa_override_parser_probe__`
and `--qa-override` returned exit 2, `unexpected argument "--qa-override"`,
before configuration/Spec loading or mutation. This is a characterization of
unsupported syntax, not a proposed final flag or an executed archive. No Spec
was moved and no QA result was modified.

Secondbrain index/qmd found and the session read the Fluxus 2026-08-25 repeated
override report and Conexus 2026-08-19 corroboration under the triaged Roundfix
Inbox. Those observations distinguish a legitimately declined gate from an
exception for missing/failed QA; they informed separate eligibility outcomes
rather than silently marking all no-report cases passed. Their historical
counts are not used as a current fleet inventory.

Exa MCP read [Git's move documentation](https://git-scm.com/docs/git-mv), which
supports the rename/index mechanics only. The maintainer supplies override
authority, the local skill supplies the existing marker, and local source/help
inspection establishes the runtime gap. No external source is claimed to
validate this repository's authorization policy.

The source-policy grant landed separately as `904b69c`. The public Baseline
applied plan `sha256:dff71221c7822544ebc6b0750181c92aad174e13492be8cc7d4942fdc17978c3`
at exit 0, updating docs-layout, autonomous-work, skill-dispatch and the Setup
Manifest from the three source modules. Runtime/test/skill-source integration
remains proposed. The known sanctioned-regeneration authorization-discovery
failure remains a repository delivery blocker; generated local guides are not
proof of complete fixture regeneration or a passing repository gate.


Final trigger/source clarification applied public Baseline plan
`sha256:fd866d5297c49907a3bbdf1b123f4e96706c8fbf41819500494529ee6878e328`
at exit 0. The archive trigger now explicitly requires non-QA Task completion,
so an overridden incomplete QA gate does not prevent skill activation. Both
changed clauses and the trigger match their generated guides; a fresh preview
reported current with no file changes. The two nested-carrier warnings remain.

`make baseline-digests` was attempted with Go 1.26.7 and failed at
`TestReadoptionCompatibilityMaintainedFixture` because suiteguard still cannot
consume the relocated authorization records. Its two partial catalog outputs
were restored. Strict inspection of eleven pending Specs reports the known
`SC-LOOP-ORDER-DIVERGENT` fixture inconsistency and no additional findings;
missing Task Graphs remain skips. These results do not establish runtime
implementation, complete regeneration or terminal QA success.
