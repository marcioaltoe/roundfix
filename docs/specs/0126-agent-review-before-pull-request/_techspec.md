---
spec: 0126-agent-review-before-pull-request
prd: _prd.md
created: 2026-09-08
---

# Configurable pre-PR review — Technical candidate

## Executive Summary

Add explicit provider selection and a truthful none path around the existing
candidate-review boundaries. Codex, Claude and CodeRabbit are enabled choices;
none performs no review and retains QA/checks. This revises the earlier provider
retirement proposal. The policy is confirmed; schema, adapters, limits and
remaining exact implementation grants are still in authoring.

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

## System Architecture

| Component | Existing package or proposed seam | Responsibility |
| --- | --- | --- |
| Review policy resolver | `internal/config` and CLI preflight | Resolve explicit provider/none, default and provenance; preserve valid agent-profile configuration. |
| Provider adapters | `internal/agent`, `internal/reviewsource`, proposed local review CLI | Inspect a candidate through Codex, Claude or CodeRabbit's supported local surface. |
| Review/omission record | `internal/store` and `internal/runevent` | Persist candidate identity and either enabled-review evidence or explicit configured omission. |
| Legacy configuration migration | Existing config and PR-feedback command routes | Preserve intentional service selection and history without universal provider requirements. |
| Corrective work | Existing Roundfix Task/Run execution | Apply authorized fixes and obtain fresh enabled review. |
| Publication eligibility | CLI delivery boundary and canonical guidance | Require the selected policy's outcome and all other applicable gates. |

Package seams are not unrestricted mutation grants. Revalidate exact files
after predecessor Specs land.

## Implementation Design

### Policy resolution

Separate provider selection from the Agent Selection Profile. Accepted policy
choices are codex, claude, coderabbit and explicit none. Keep the confirmed
Codex built-in default, User Config overrides and stronger explicit Project Config overrides. Reuse applicable profiles for an
agent provider's model/effort, with visible provenance; do not treat coderabbit
or none as ACP runtimes. A provider/profile mismatch is a named configuration
problem, never permission to silently use a different provider.

The final YAML field and legacy precedence mapping remain to be authored.
Existing `profiles.review` and `review_source.name` represent different facts;
inspect both during migration and preserve explicit intent. Merely reading a
legacy provider block cannot enable it, and absence cannot select none. A
configured runtime outside the three named providers needs an explicit supported
mapping or a migration refusal, not an invented fourth enabled option.

### Provider interfaces

For Codex and Claude, create an independent reviewer session with repository
reads and denied writes. Preserve configured selection and characterize actual
result shape, cancellation and inherited permissions. Retain the dated native
interface observations below as inputs, not current proof.

For CodeRabbit, use its documented local CLI review surface before the PR.
Its existing GitHub PR-feedback service remains a distinct interface. The CLI
reference documents agent JSON output and local scope selection. Characterize
the installed version, final result shape, exact candidate coverage, file
mutation behavior, authentication and cancellation before treating its output
as an accepted receipt. Scope defaults cannot omit the candidate's relevant
files; an empty or partial review is not complete evidence. Provider access
uses its existing approved boundary; selection alone introduces no new
credential or code-disclosure authorization.

Only the selected enabled provider's readiness is required. Unknown selection,
missing capability, failed/skipped/incomplete output or unavailable credentials
blocks that mode. No such failure activates none. Declared agent-model fallbacks
retain their approved scope and do not substitute another provider implicitly.

### Explicit none and durable outcomes

None bypasses reviewer process creation, provider requests, reviewer readiness
and review waiting. Record repository/candidate identity, effective policy and
configuration provenance as an intentional configured omission outside the
reviewed tree. It is not a passing review, an empty receipt or a reused historical
Review Skipped result. The publication boundary accepts this distinct outcome
when other authority, QA and required checks are satisfied.

Enabled review records repository, base/merge-base/head, provider, applicable
model/effort and provenance, execution completion, coverage/omissions, verdict
and findings. Invalid, absent, stale or incomplete records cannot authorize the
selected enabled path. Revalidate current candidate and policy at publication
and recovery; neither old enabled evidence nor an old none record overrides a
changed configuration. This storage design is proposed implementation.

### Corrections, publication and compatibility

Corrections remain runtime-owned, bounded by the consuming Spec and limits,
and require new enabled review of the changed candidate. The archive-first
late-correction trade-off remains pending. With none, the absence of reviewer
findings does not erase QA failures or previously recorded unresolved work.

CodeRabbit remains optional and selectable. Remove unconditional defaults and
instructions that force that service across all policies. Unselected or
explicitly disabled providers receive no generated manual or automatic request.
Do not uninstall the GitHub App, mutate another repository or silently rewrite
User Config. Existing historical Runs/events and outcomes retain their meaning.
Legacy command migration must be explicit; a post-PR service result is not
fabricated pre-PR evidence.

Both enabled review and none lead to publication only under the approved action
scope. Required checks still cover the head being merged, and repository/GitHub
protection remains binding. No administrator bypass is added.

### Data Models

Persist a policy outcome distinguishing completed enabled review from explicit
configured omission. Both identify repository, candidate and policy provenance;
only completed review carries reviewer verdict/coverage/findings. Preserve prior
attempts and historical outcome types. Final schema remains a candidate.

### API Contracts

A proposed local review operation resolves the candidate and policy, produces
structured evidence or configured omission, or returns a named failure. None
needs no reviewer/session/PR identity. Final command flags and YAML schema must
be authored and validated before Task dispatch; no runnable syntax is invented
by this proposal.

## Coverage Map

- PRD Goals 1 and 3; User Stories 1 and 3; Core Features 1–4 → policy resolver, provider adapters and legacy migration.
- PRD Goals 2 and 4; User Stories 2 and 5; Core Features 2, 3 and 7 → review/omission records.
- PRD Goal 5; User Story 4; Core Feature 5 → corrective work.
- Core Features 6 and 8 → publication eligibility and policy-aware delivery order.

## Integration Points

The [owned source index](references/_index.md) preserves the adopted evidence.
Spec 0119 supplies authority; Spec 0123 supplies runtime capability evidence.
Spec 0127 consumes the policy result, including none. ADR-0153 records the
maintainer correction while ADR-0151 retains default/profile precedence.

## Testing Approach

1. Omitted selection uses the confirmed default; each explicit codex/claude/coderabbit/none choice takes precedence with provenance. Invalid values and incompatible profiles refuse.
2. Each enabled provider reviews a real local candidate before a PR exists, covers its intended files, leaves source/index/refs unchanged and emits a complete validated result.
3. None creates no reviewer process, performs no reviewer readiness probe or provider request, records configured omission and permits otherwise authorized publication.
4. Failed, missing, skipped, cancelled, truncated or stale enabled review never selects none; candidate or policy changes invalidate prior eligibility.
5. Unselected providers may be absent without blocking the selected path. Generated instructions do not call opted-out providers; intentionally selected CodeRabbit remains functional.
6. Required-check failure blocks merge with every mode, including none. History remains readable and unchanged; corrections preserve scope and current evidence.
7. Real external acceptance evidence and terminal QA exercise each approved mode and distinguish unavailable environment from success. Provider calls await their recorded access and limits.

These are planned observations, not executed verification. Author exact commands
and controls after resolving the remaining design and grants.

## Build Order

1. Policy/schema resolution, compatibility plan and source provenance (depends on: none).
2. Explicit-none outcome and no-call boundary (depends on: 1).
3. Independent agent adapters and local CodeRabbit adapter (depends on: 1).
4. Candidate-bound records and policy-aware eligibility (depends on: 2, 3).
5. Corrective work and publication integration (depends on: 4).
6. Owned skill alignment, canonical guidance and sanctioned outputs (depends on: 1, 2, 3, 4, 5).
7. Real mode-specific journeys and terminal QA (depends on: 1, 2, 3, 4, 5, 6).

The Task Graph does not exist yet. This order does not grant implementation or
claim completion of the narrow canonical instruction change.

## Risks & Considerations

None intentionally omits independent review and must remain observable. Optional
providers have different coverage, result and account contracts. Defaults,
legacy configuration and agent profiles need unambiguous migration. Time,
consumption and correction limits plus archive-first late-correction policy
remain pending. No missing limit is interpreted as unlimited.

## Decisions

- ADR-0153 confirms codex, claude, coderabbit or none and replaces mandatory review/CodeRabbit removal. Codex remains the default; explicit project selection wins.
- QA, required checks, action authority and limits remain mandatory for every policy. Failure of an enabled provider never selects none.
- The dated narrow grant authorizes the canonical policy and its generated guidance only; broader adapters/configuration/test work remains proposed.

## Cross-Spec dependencies

Required predecessor contracts: [0119](../0119-spec-contained-authorization/_techspec.md), [0123](../0123-runtime-readiness-and-model-capabilities/_techspec.md).
Shared files integrate serially; revalidate their current contracts after each
predecessor. No dependency reference is an implementation grant.

## Dated native-interface evidence — 2026-09-08

The removed design recorded read-only inspection of `codex-cli 0.153.4` and
Claude Code `2.1.263`; these are dated observations, not promises about future
installed versions. Codex native target arguments and a custom review prompt
were mutually exclusive: combining `--base main` with a positional prompt
returned parser exit 2 before model work. Characterize how the selected review
mode receives the adopted Spec and standards without inventing an instructions
flag or assuming that help output proves the actual structured result.

The Claude investigation distinguished `dontAsk`, the allowed-tool list and
the effective tool surface. It proposed Read/Glob/Grep with a supplied diff,
with hooks and inherited MCP permissions included in the read-only proof.
Unrestricted Bash was not part of that proposal. The read documentation said
`--bare` did not use subscription OAuth and `--max-budget-usd` controlled API
spending, not subscription quota. Revalidate those constraints against the
installed adapter before any authorized live exercise; no such call has run.

Sources read through Exa: [Codex CLI reference](https://developers.openai.com/codex/cli/reference),
[Claude headless execution](https://code.claude.com/docs/en/headless) and
[Claude CLI reference](https://code.claude.com/docs/en/cli-reference).
The observations narrow the required adapter characterization; they do not
prove access, successful review or sufficient permission isolation.


## CodeRabbit local-interface evidence — 2026-09-08

Exa found the official [local review overview](https://docs.coderabbit.ai/overview/ide-cli-review)
and read the [CLI reference](https://docs.coderabbit.ai/cli/reference). The returned
reference documents local tracked-change review, committed/uncommitted scope
choices and one JSON object per stdout line with agent output. It is a viable
surface to characterize for the selected coderabbit mode; the fetched excerpt
does not prove its entire terminal schema or Roundfix integration. No service
was invoked, installed or authenticated here.

Secondbrain's earlier reviewer/readiness captures were consulted after an
index/qmd search. They preserve the prior mandatory-review/removal decision;
this maintainer correction supersedes that policy rather than rewriting those
historical captures. Published provider documentation establishes interfaces;
only the maintainer authorizes none and makes CodeRabbit optional.


## Canonical application and verification limit — 2026-09-08

The narrow policy grant landed as commit `888cb62` before the source edits.
Core and autonomous-work modules were changed and the public Baseline update
applied reviewed plan `sha256:c1767cd21b6f1824f9f3b73dfb733f02422b580eb2d3391a0e7f0ca06f3acf36`
at exit 0, updating only agent-instructions, autonomous-work and the Setup
Manifest. This proves those generated postimages, not a working provider mode.

`make baseline-digests` with Go 1.26.7 failed at
`TestReadoptionCompatibilityMaintainedFixture`: suiteguard rejected the generated
catalog digest/normalized writes because its authorization discovery still
requires the deleted workflow directory. Partial outputs were restored to
their pre-attempt bytes. The known seven-file compatibility proposal is still
pending. Consequently the shipped formatter fixture retains the old loop text
and strict Spec Check reports `SC-LOOP-ORDER-DIVERGENT`. This is unresolved
regeneration, not an exemption or a passing repository gate. No fixture, test,
runtime/config parser or provider adapter was manually changed.


The final precedence clarification used public Baseline plan
`sha256:13697f951a96ed88161731dc0a833e1c12370de31375852547fc56782c6ed643`:
apply exited 0 and a fresh preview reported current with no file changes.
Independent Standards and Spec rechecks found no remaining reviewed policy
inconsistency. All eleven pending Specs still report the known
`SC-LOOP-ORDER-DIVERGENT` regeneration failure and no other current detector
finding. Their absent Task Graphs remain explicit skips, not implementation
readiness. The two preserved nested-carrier warnings remain unchanged.


## Authorized archive disposition

Consume the archive policy from Spec 0122 and ADR-0154. An applicable explicit
user authorization can archive the covered Spec with unmet QA, recording the
override and preserving original evidence. This does not reopen or complete
Tasks, declare QA passed, imply review approval or authorize publication/merge.
A durable workflow records the overridden archive and evaluates subsequent
actions against their own approval and gates; it does not retry the waived
archive prerequisite or ask again for the same applicable archive approval.
The absence of authority for a later action remains a separate visible blocker.
