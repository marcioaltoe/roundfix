---
spec: 0126-agent-review-before-pull-request
status: active
created: 2026-09-08
surfaces: [cli, data, docs]
---

# Configurable review before a Pull Request

A repository chooses `codex`, `claude`, `coderabbit`, or explicit `none` for
pre-PR review. Codex remains the default and an explicit project choice takes
precedence. This confirmed policy replaces mandatory independent review and
mandatory removal of CodeRabbit. The implementation remains in authoring:
configuration migration, provider contracts, limits and remaining protected
mutations must be settled before Task dispatch.

## Project Constraints

- Identifier strategy: applicable — preserve the existing Run and Task identities and use immutable Git commit identities to identify the reviewed candidate. New work branches use purpose prefixes such as `feat/`, `fix/`, and `refactor/`; Roundfix-owned Run/Task branches retain their documented namespace. The maintainer explicitly removed the personal-prefix conflict on 2026-09-08. No new identifier format is approved here. Source: `docs/agents/domain.md`, `docs/agents/agent-instructions.md`.
- Authentication and HTTP: applicable — reuse the maintainer's existing authenticated runtimes and any explicitly selected provider service; do not introduce credentials, change authentication policy, or silently route subscription work through paid APIs. No backend HTTP guide exists for this CLI repository, so absence supplies no authorization. Native review access, tool permissions, and monetary or quota bounds still require confirmation. Source: `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — retain the active execution, evidence, and authoring contracts while proposing the configurable pre-PR review contract. Source: `docs/agents/domain.md`, `docs/agents/autonomous-work.md`, `docs/agents/spec-routing.md`.
  ADR-0014 applies: the Daemon runs Task Verification and settles the outcome.
  ADR-0019 applies to the historical Watch contract: Clean requires evidence that the Open Pull Request is merge-ready; new policy modes must not rewrite that recorded meaning.
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
  ADR-0151 retains the Codex default and explicit project-selection precedence for agent review.
  ADR-0153 applies: the Pre-PR Review Policy permits codex, claude, coderabbit or explicit none; enabled-review failures do not become configured omissions.
  ADR-0093 applies: mechanical citation accounting does not establish semantic review correctness.
  ADR-0097 applies to carried QA evidence consumed by delivery: retain declared unchanged evidence rather than inheriting a pass after its inputs move.
  ADR-0104 applies: acceptance uses independent evidence with its actual origin; missing external evidence remains visible under the declared policy.
  ADR-0154 applies to archive disposition: an explicitly user-authorized QA Archive Override preserves actual QA/Task evidence and does not satisfy independent delivery gates.
- Tooling authority: applicable — express maintainer authorization covers the canonical policy in [the narrow grant](references/2026-09-08-configurable-review-policy-authorization.md), bounded files: `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/autonomous-work.json`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md` and `docs/agents/setup-context.json`, with sanctioned digest regeneration. Other exact governed mutations remain proposed in [_authorization.md](_authorization.md); status proposed and a null grant authorize no mutation. Bounded proposed files: `.roundfixrc.yml`, `.coderabbit.yaml`, `internal/baseline/assets/modules/core.json`, `internal/baseline/assets/modules/autonomous-work.json`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/agents/openai.yaml`, `skills/roundfix/SKILL.md`, `skills/roundfix/agents/openai.yaml`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go`, `skills/baseline_skill_contract_test.go`. Sanctioned regeneration follows source approval. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A repository can select its pre-PR reviewer or intentionally omit review.
- Enabled review identifies the examined candidate and cannot clear a changed candidate.
- An unselected provider imposes no readiness, request or account dependency.
- Earlier Runs and review evidence retain their original meaning and remain readable.
- Corrections preserve the approved Spec and Verification within recorded limits.

## User Stories

1. As a maintainer, I want to select Codex, Claude or CodeRabbit before publication, so that the workflow uses my chosen reviewer.
2. As a maintainer, I want complete evidence for the current candidate, so that stale or failed review cannot clear another commit.
3. As a repository owner, I want explicit none to skip review, so that otherwise authorized delivery can proceed with QA and required checks.
4. As a Supervisor, I want enabled-review findings routed into bounded corrective work, so that corrections retain the approved contract.
5. As a maintainer inspecting delivery, I want to distinguish reviewed, intentionally disabled and failed review, so that an omission never looks like approval.

## Core Features

1. Resolve one Pre-PR Review Policy: `codex`, `claude`, `coderabbit`, or explicit `none`. Keep the Codex default and explicit project precedence. For agent providers, preserve applicable configured model/effort and selection provenance. None is not an Agent Runtime; CodeRabbit's service configuration is distinct from an ACP profile. Exact configuration schema and migration remain implementation design.
2. Enabled providers inspect the explicit current candidate before a PR exists. Codex and Claude use independent read-only reviewer sessions; CodeRabbit uses a supported local review surface. Evidence records repository/base/head, effective provider, coverage, findings and execution outcome. Unsupported capability, incomplete output, runtime failure, timeout or stale evidence blocks that selected mode; it never selects none.
3. Explicit none performs no reviewer/provider call, no reviewer readiness probe and no wait for review. It records a configured omission for the candidate, permits otherwise authorized publication and merge, and preserves QA, required checks and external branch protection. A provider-reported skip, missing response or failure is not this explicit policy.
4. Unselected or explicitly disabled providers receive no automatic or manual request from Roundfix or generated guidance. CodeRabbit remains available when selected; remove universal CodeRabbit requirements rather than requiring service removal. Historical PR-feedback commands and the local pre-PR interface remain distinguishable. Migration cannot silently opt a repository into a service or into none.
5. Enabled-review findings receive evidence-backed dispositions. Corrections remain runtime-owned through Roundfix, stay within the Spec and Verification, and consume approved limits. A changed candidate requires fresh enabled review. The proposed archive-first late-correction path requires sufficient corrective-Spec authority; it is still pending that policy decision.
6. Publication consumes either complete applicable enabled-review evidence or the explicit configured-omission record. Required repository/GitHub checks still cover the current head before merge. None does not waive those checks or approve unrelated actions.
7. Historical Runs, issues and receipts retain their original bytes and classifications. A historical Review Skipped or Clean Unverified result never becomes a configured-none record. Selecting a new policy does not retroactively rewrite old results.
8. Proposed order is implementation and terminal QA, archive and commit, configured review or explicit omission, publication, current-head checks and merge. Recording either review evidence or omission must not create an unexamined candidate change; an enabled review must be repeated if its candidate changes.

## User Experience

Show the selected policy, its source, current candidate, applicable limits and
outcome. With none, state that review was disabled by configuration and continue
through the other delivery gates. With an enabled provider, show its result or
the concrete failure. Enabling these choices does not change the current
project configuration or install, authenticate or invoke a provider.

## Non-Goals / Out of Scope

- A new implementation framework or the durable multi-Spec owner from Spec 0127.
- Changing GitHub branch protection, required checks, credentials or account plans.
- Editing upstream-managed review or PR-workflow skills.
- Reopening archived Specs or inheriting enabled-review validity across a candidate change.
- Automatically converting provider failure into none or buying API credits.
- Treating the provider choice as permission for unrelated external code disclosure; use the existing approved provider boundary.

## Success Metrics

| Acceptance observation | Evidence required |
| --- | --- |
| Every enabled provider works before a PR exists | Codex, Claude and CodeRabbit each inspect the candidate through their proved local interface and produce complete evidence. |
| Explicit none omits review | No reviewer process, readiness probe or provider call occurs; the durable record says configured omission and delivery continues through QA/checks. |
| Invalid or unavailable is not none | Unknown policy, unavailable selected capability, malformed output and provider-reported skip produce named failures without fallback to none. |
| Candidate changes invalidate enabled review | Changed head or base relationship requires new review. |
| Optional providers impose no hidden dependency | Each unselected provider can be absent without blocking review selection; generated instructions never invoke it. |
| Existing protections remain effective | A failed required check blocks merge even with none; external protection is not modified. |
| History and corrections remain truthful | Earlier records are preserved; corrections retain scope and fresh evidence. |
| External acceptance evidence | A selected reviewer exercises a real defect and correction from outside this Spec's fixtures; missing evidence is reported under the declared policy. |

## Decisions

- Confirmed on 2026-09-08: review may be Codex, Claude, CodeRabbit or none. The user's token "code" is interpreted as Codex. No current project choice is changed by declaring the supported modes.
- Codex remains the default; explicit project selection takes precedence. None permits no-review delivery while QA, required checks, authority and limits remain mandatory.
- This replaces the earlier mandatory-review and total-CodeRabbit-removal policy; the updated decision is recorded in ADR-0153.
- The Supervisor authors and dispatches; code/tests and Verification remain runtime/Daemon work.
- Archive-first late-correction behavior, schema details and remaining implementation grants are still proposed.

## Open Questions

- Correction cycles, wall-clock limits and API/subscription bounds for enabled providers remain pending; none incurs no reviewer call.
- Final provider-selection schema, migration of existing review profiles and legacy PR-feedback configuration, and unsupported-profile handling must be settled before implementation.
- Archive-first handling and authority for a late corrective Spec remain pending.
- Exact implementation paths beyond the narrow canonical policy grant remain proposed.

## Research and limitations

Secondbrain's existing reviewer/readiness captures preserve the prior mandatory
review policy; the latest maintainer instruction changes that policy. Exa
located and read the official [CodeRabbit CLI reference](https://docs.coderabbit.ai/cli/reference)
and [local review overview](https://docs.coderabbit.ai/overview/ide-cli-review).
They document local review without a PR and structured agent output, making
CodeRabbit a candidate pre-PR adapter. Published interfaces do not establish
Roundfix integration, installed-version behavior, account access or paid authority.
The TechSpec retains dated Codex/Claude interface observations and the new
CodeRabbit evidence. No provider was invoked by this authoring change.

## Source ownership and authoring state

The [source index](references/_index.md) records the primary owned evidence.
The provider-opt-out Finding still requires unselected providers to remain
unrequested; it no longer requires removing the selected CodeRabbit option.
The [_techspec.md](_techspec.md) and [_authorization.md](_authorization.md) make
remaining implementation reviewable. No Task Graph or implemented mode is
claimed by this policy correction.


## Authorized archive disposition

Consume the archive policy from Spec 0122 and ADR-0154. An applicable explicit
user authorization can archive the covered Spec with unmet QA, recording the
override and preserving original evidence. This does not reopen or complete
Tasks, declare QA passed, imply review approval or authorize publication/merge.
A durable workflow records the overridden archive and evaluates subsequent
actions against their own approval and gates; it does not retry the waived
archive prerequisite or ask again for the same applicable archive approval.
The absence of authority for a later action remains a separate visible blocker.
