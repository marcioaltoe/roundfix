---
spec: 0126-agent-review-before-pull-request
prd: _prd.md
created: 2026-09-08
---

# Configured independent review before a pull request — Technical candidate

## Executive Summary

Add a native local-candidate review phase using the existing atomic review profile resolver and durable evidence store. Retire CodeRabbit as an operational dependency while retaining historical readers. The design spends a separate review session and revalidates its candidate instead of treating an absent or skipped provider check as approval.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The cross-Spec order and remaining
decisions are in [the portfolio plan](../../workflow/2026-09-08-pending-work-plan.md).

## Project Constraints

- Identifier strategy: applicable — preserve the existing Run and Task identities and use immutable Git commit identities to identify the reviewed candidate. New work branches use purpose prefixes such as `feat/`, `fix/`, and `refactor/`; Roundfix-owned Run/Task branches retain their documented namespace. The maintainer explicitly removed the personal-prefix conflict on 2026-09-08. No new identifier format is approved here. Source: `docs/agents/domain.md`, `docs/agents/agent-instructions.md`.
- Authentication and HTTP: applicable — reuse the maintainer's existing authenticated local runtimes; do not introduce credentials, change authentication policy, or silently route subscription work through paid APIs. No backend HTTP guide exists for this CLI repository, so absence supplies no authorization. Native review access, tool permissions, and monetary or quota bounds still require confirmation. Source: `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — retain the active execution, evidence, and authoring contracts while proposing a distinct pre-PR review contract. Source: `docs/agents/domain.md`, `docs/agents/autonomous-work.md`, `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon runs Task Verification and settles the outcome.
  ADR-0019 applies to the historical Watch contract: Clean requires evidence that the Open Pull Request is merge-ready; retirement must not rewrite that recorded meaning.
  ADR-0020 applies to retained acpx Batch execution: a valid parsed prompt result outranks a later teardown exit, which remains journaled; this does not establish success rules for a separate native reviewer.
  ADR-0038 applies: Daemon Verification Feedback permits one repair in the same Agent Session, distinct from the proposed independent-review correction limit.
  ADR-0056 applies to corrective Spec Runs: Task Capacity and Verification Capacity remain separate, and the exclusive temporary-failure retry retains its existing bound.
  ADR-0057 applies: only the Daemon writes Implement Task status; reviewer output and Agent handoffs cannot claim a terminal Task outcome.
  ADR-0091 applies: QA remains the terminal Task node rather than a separate invocation after graph completion.
  ADR-0096 applies: the Daemon-owned mechanical QA stage withholds the Agent turn on blocking facts and cannot loosen verdict semantics.
  ADR-0117 applies: authoring defects are checked at the stage that produces them; commit-dependent audits and user-surface evidence remain in the later gate.
  ADR-0127 applies to reviewer process lifecycle: readiness reports process residue without inventing Run records or settling work from the inventory.
  ADR-0139 applies: the Run Database and checkout guard preserve one Active Run per work target and reject competing mutation.
  ADR-0142 applies to preserved historical Review Source Evidence: expected-head classification and the distinction between Clean and Clean Unverified remain readable, while a new pre-PR contract requires an explicit decision.
  ADR-0080 applies to the inherited QA evidence cited in the adopted measurement: distinguish environment-blocked rows from product failure and successful execution.
  ADR-0151 applies: use Codex by default and preserve an explicit configured review profile, with effective selection provenance and explicit capability refusal.
  ADR-0093 applies: mechanical citation accounting does not establish semantic review correctness.
  ADR-0097 applies to carried QA evidence consumed by delivery: retain declared unchanged evidence rather than inheriting a pass after its inputs move.
  ADR-0104 applies: acceptance uses independent evidence with its actual origin; missing external evidence remains visible under the declared policy.
- Tooling authority: applicable — exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.roundfixrc.yml`, `.coderabbit.yaml`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/autonomous-work.json`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/agents/openai.yaml`, `skills/roundfix/SKILL.md`, `skills/roundfix/agents/openai.yaml`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`, `skills/baseline_skill_contract_test.go`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Configured reviewer | `internal/config/profiles.go, CLI profile preflight and daemon selection owner` | Resolve CategoryReview, exact effective tuple, source and explicit fallbacks. |
| Native review adapter | `internal/agent/review.go and internal/cli/review.go (new in existing packages)` | Review an explicit clean candidate in an independent session without changing repository content. |
| Candidate review receipt | `internal/store and internal/runevent; local-candidate evidence in reviewsource` | Persist repository/base/merge-base/head, selection, terminal verdict, findings and omissions outside the reviewed tree. |
| Provider retirement | `internal/config/config.go, CLI default provider factories and legacy command routes` | Remove operational CodeRabbit defaults/requests while preserving known legacy decoding/history. |
| Finding correction cycle | `existing Roundfix Task execution and receipt linkage` | Apply bounded runtime-owned corrections and require fresh review of the changed candidate. |
| Publication eligibility | `CLI delivery/pre-PR check and owned/canonical guidance` | Require current independent review and GitHub checks under the approved through-merge policy. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Selection and review boundary

Codex is the confirmed default, and explicit `.roundfixrc.yml` review selection
wins. Reuse `profiles.review` rather than introducing a parallel reviewer key.
The current project tuple is Codex/Luna/max with its declared Codex/Sol/high
fallback; the built-in profile remains a distinct default. Preserve User and
Project Config precedence, effective tuple and source. The automated launcher
omits selection flags unless the user supplied an actual invocation override.
Invalid config or absent required review capability is an explicit refusal;
only approved declared fallback reasons can activate another selection.

Create a fresh reviewer session over a clean candidate with repository reads
and denied writes. Codex's noninteractive review and Claude's structured
headless execution are integration surfaces, not interchangeable output
contracts. Characterize each allowed runtime's terminal/result shape and
permission behavior. A configured runtime without a proved adapter is named as
unsupported; it is never silently replaced by Codex. Do not reuse the current
sealed-prompt path as if it offered repository investigation: its empty working
directory, denied tools and 2 MiB input bound solve a different contract.

### Receipt and correction

A complete receipt fixes repository identity, base/merge-base/head, requested
profile and origin, effective runtime/model/effort, execution completion,
coverage/omissions, verdict and typed findings. Store it outside the reviewed
tree so recording review does not change its own candidate. Empty, truncated,
failed, cancelled, stale or missing results are not approval. Re-evaluate the
review range whenever the candidate changes and retain superseded receipts.

Disposition distinguishes implemented correction, evidence that disproves a
finding and an explicitly accepted real risk. The reviewer never implements
its own correction or weakens Verification. Corrections are separate runtime
Tasks/Runs inside their grant and consume a finite cycle budget. After the
accepted archive-first sequence, a late blocker parks publication and needs a
new corrective Spec with sufficient authority; do not reopen archived evidence
or inherit unlimited approval from the previous Spec.

### CodeRabbit migration and publication

`review_source.name` is the old PR-feedback provider, not an ACP runtime.
Decode the known historical CodeRabbit block as warned inert legacy data when
needed for compatibility; invalid new values still refuse with the correct
configuration path. Retired fetch/resolve/watch provider operations refuse
before network, Run creation or Agent work. Remove active factories, defaults,
requests and provider-specific canonical obligations, plus this project's
obsolete YAML/block in the approved change. Preserve original Runs/events and
artifact readers; never purge history to satisfy a string search.

The Pantheon evidence requires an explicit provider opt-out to remain effective.
Disabled automatic review is not authority to request that provider manually.
Do not mutate another repository or its User Config as part of this migration.
Publication requires the current candidate's native review and the repository's
required checks; a later push invalidates prior candidate evidence.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

A review receipt binds repository/base/merge-base/head, profile provenance, effective reviewer, completion, verdict, findings and omissions. Persist it outside the reviewed tree and preserve superseded receipts.

### API Contracts

The proposed native review command accepts an explicit clean local candidate and uses CategoryReview. It returns structured complete evidence or a named refusal/failure; it needs no fabricated PR number. Legacy provider routes refuse before network/Run creation.

## Coverage Map

- PRD Goal 1 → Configured reviewer, Native review adapter.
- PRD Goal 2 → Candidate review receipt.
- PRD Goal 3 → Provider retirement, Publication eligibility.
- PRD Goal 4 → Provider retirement, Candidate review receipt.
- PRD Goal 5 → Finding correction cycle.
- User Story 1 → Configured reviewer, Native review adapter.
- User Story 2 → Candidate review receipt.
- User Story 3 → Provider retirement, Publication eligibility.
- User Story 4 → Finding correction cycle.
- User Story 5 → Candidate review receipt, Provider retirement.
- Core Feature 1 → Configured reviewer and Native review adapter.
- Core Feature 2 → Candidate review receipt.
- Core Feature 3 → Provider retirement.
- Core Feature 4 → Finding correction cycle.
- Core Feature 5 → Finding correction cycle.
- Core Feature 6 → Publication eligibility.
- Core Feature 7 → Provider retirement and Candidate review receipt.
- Core Feature 8 → Publication eligibility.

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

1. No override uses the Codex default; explicit project Claude/OpenCode/Codex selection retains its actual tuple/provenance or reports unsupported capability. Invalid config never silently defaults.
2. Independent review reads the clean candidate, cannot change its files/index/refs, records a complete receipt and rejects empty/truncated/error/cancelled output.
3. Changed head/range invalidates review; corrections require fresh evidence and remain within approved Spec and Verification scope.
4. Known legacy config remains inspectable but cannot trigger CodeRabbit network requests; retired commands refuse before mutation and disabled-provider guidance never restores manual requests.
5. Historical records remain byte-identical and readable; pre-PR review operates without a synthetic PR number.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. CategoryReview selection and native review contract (depends on: none).
2. Read-only adapters and exact terminal/selection evidence (depends on: 1).
3. Durable candidate receipts and eligibility (depends on: 1, 2).
4. Operational CodeRabbit retirement and known-legacy migration (depends on: 1, 2, 3).
5. Bounded findings/correction and publication integration (depends on: 3, 4).
6. Owned skills/canonical guidance and sanctioned outputs (depends on: 1, 2, 3, 4, 5).
7. Real configured-runtime and public refusal journeys with terminal QA (depends on: 1, 2, 3, 4, 5, 6).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

The native adapter must prove read-only behavior and result completeness, not only process exit. The archive-first late-correction policy and monetary/cycle limits remain pending approval. Preserve opted-out providers and historical evidence during migration.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge is confirmed only with independent review and required checks approved for the current candidate; releases, tags and paid consumption are not implied.
- Reviewer policy is settled by [ADR-0151](../../adr/0151-configured-review-profiles-select-the-independent-reviewer.md): Codex by default, explicit review profile in `.roundfixrc.yml` taking precedence.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
