---
spec: 0126-agent-review-before-pull-request
status: active
created: 2026-09-08
surfaces: [cli, data, docs]
---

# Independent review before a Pull Request

Roundfix currently depends on CodeRabbit for review feedback and evidence, even
when a repository wants to review locally before opening a Pull Request. The
maintainer requests removal of that operational dependency and a local,
independent Codex or Claude Code review that proves which candidate it examined.
Reviewer selection is confirmed: Codex by default, with the explicit configured
review profile taking precedence. This PRD remains in authoring because limits,
migration details and exact protected-file authority are pending. It does not
authorize execution.

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

## Goals

- A repository can obtain an independent local review before creating a Pull Request.
- Review evidence identifies the examined candidate and becomes unusable when that candidate changes.
- Roundfix operation and owned skills require no CodeRabbit service, request, configuration, or approval signal.
- Historical Runs, Review Issues, and evidence remain readable after the operational integration is retired.
- Corrections preserve the approved Spec and Verification instead of rewriting the contract to close a finding.

## User Stories

1. As a maintainer, I want an independent agent to review my candidate before publication, so that review can catch defects before a Pull Request opens.
2. As a maintainer, I want evidence identifying the exact candidate and reviewer, so that a stale or skipped review cannot clear a later commit for merge.
3. As a repository owner, I want CodeRabbit removed from the operational workflow, so that its availability and account limits cannot block local delivery.
4. As a Supervisor, I want review findings routed into bounded corrective work, so that fixes remain subject to the Spec, the Daemon's Verification, and fresh review.
5. As a maintainer investigating an earlier Run, I want its original review evidence preserved, so that migration does not erase the delivery record.

## Core Features

1. Independent local review examines an explicit repository candidate and the applicable Spec and standards in a fresh reviewer session. The reviewer reports findings and coverage without mutating code, tests, instructions, Verification, or Git state. Codex is the default; an explicit effective project review profile takes precedence. Supporting a runtime does not authorize additional billing.
2. Review evidence records the candidate's base and head, the effective reviewer selection, the result, omissions, and execution outcome. A changed head or base relationship, incomplete coverage, malformed response, runtime failure, timeout, or interruption cannot produce a passing review.
3. CodeRabbit retirement removes active provider requests, defaults, command routing, configuration requirements, and mandatory owned-skill instructions. Explicit legacy usage receives a migration diagnostic before side effects. Existing `none` support must not be assumed; legacy command retirement and inherited configuration migration require a confirmed public contract.
4. Findings receive evidence-backed dispositions. A disproved finding is distinct from accepting a real risk. Unresolved blocking findings prevent publication; a failed check or review cannot be automatically accepted or waived. Optional suggestions can be retained as backlog intent under the approved policy.
5. Corrective implementation remains owned by the selected runtime through a Roundfix Run. It cannot weaken the Spec or Verification, invent infrastructure to suppress a finding, or silently expand protected-file authority. A blocking finding after archive requires a new corrective Spec with its own approval; the completed or archived Spec remains unchanged. Any candidate change requires fresh review. Correction, time, and spend limits remain pending and must be enforced once approved.
6. The publication boundary checks that independent review and required checks cover the current candidate. Until that proof and the approved limits exist, conditional authority to deliver through merge does not enable unattended merge.
7. Earlier Runs, Review Issues, comments, and evidence retain their original provenance and remain inspectable. Keeping historical CodeRabbit text does not retain a live provider integration or reinterpret an old outcome as independent local review.
8. Proposed delivery order is implementation and terminal QA, archive and commit, independent review of the resulting candidate, then publication and current-head checks before merge. Archive invalidates any earlier review. A review receipt must not add an unreviewed commit; no artifact-only exception is approved by this proposal.

## User Experience

The maintainer chooses an explicit candidate and can see which reviewer will
examine it, which limits apply, whether the review completed, and why publication
is blocked. Findings name the affected location and evidence. The result
distinguishes completed review from unavailable, incomplete, stale, or failed
execution. Legacy users receive a bounded migration action, not an invented
successful disabled mode.

## Non-Goals / Out of Scope

- A new implementation or orchestration framework.
- Uninstalling an organization-wide GitHub App or altering unrelated repositories.
- Editing upstream-managed `review` or `github-pr-workflow` skills.
- Reopening completed or archived Specs to apply late review corrections, or treating an archived candidate as covered by a review of its prior head.
- Buying API credits, changing credentials, or treating unknown quota as unlimited.
- Allowing a reviewer to implement its own corrections or waive failed gates.
- Implementing the durable multi-Spec delivery owner; that is Spec 0127.
- Renaming historical branches or changing the already permitted Roundfix Run namespace.

## Success Metrics

| Acceptance observation | Evidence required |
| --- | --- |
| A candidate receives review before a PR exists | A live selected native runtime examines the candidate and produces inspectable, complete evidence. |
| A changed candidate cannot reuse the prior verdict | Public behavior rejects stale review after the head or reviewed base relationship changes. |
| Archive and late findings preserve evidence | A review of the pre-archive head cannot clear the archive commit; a blocking post-archive finding leaves the original Spec unchanged and blocks publication pending an authorized corrective Spec. |
| The reviewer cannot alter its contract | The public review flow leaves the examined tree unchanged and refuses attempted writes or reports an execution failure. |
| CodeRabbit is operationally absent | Public legacy paths refuse before provider calls; the active configuration and owned workflow contain no requirement to invoke the retired service. |
| History remains readable | Earlier Run and review artifacts retain their original content and render after migration. |
| External acceptance evidence | A reviewer exercises a real defect and correction from a repository this Spec did not build, with origin recorded. Missing external evidence is reported under the repository's declared policy rather than invented. |

## Decisions

- The maintainer requested complete removal of the operational CodeRabbit dependency and review before opening a Pull Request.
- The maintainer confirmed conditional delivery through merge only after Specs and limits are approved, with independent review and required checks passing for the current commit.
- The Supervisor authors and dispatches; implementation and tests remain runtime work through a Run.
- Proposed archive-first review preserves the current archive boundary and exact-head requirement. Its trade-off is a new corrective Spec and gate for late defects, potentially waiting for a new grant; the maintainer has not accepted that policy choice yet.
- This artifact is a concrete proposal in authoring. Sources are background inputs pending an implementation commitment; none has been adopted, moved, or marked promoted by this PRD.

## Open Questions

- Reviewer policy is settled: Codex by default, with the explicitly configured review profile taking precedence. The current project tuple and declared fallbacks are retained.
- Correction cycles, wall-clock limits, and paid API/subscription bounds — maintainer decision pending. Two corrective review cycles is a proposal, not an approved allowance; the existing rule to re-examine more than two QA corrective Tasks still applies.
- Legacy command/configuration retirement and the pre-PR readiness contract — architecture decision pending before TechSpec authoring can settle implementation.
- Archive-first final review and the authority for any later corrective Spec — proposed ordering pending confirmation; neither an archive exception nor approval inheritance is assumed.
- Exact protected paths and sanctioned generated outputs — maintainer decision pending in the proposed authorization record.
- The naming bootstrap is resolved by the maintainer's purpose/Run-namespace decision. Supervisor feature implementation remains prohibited; implementation still belongs to a Roundfix Run.

## Research and limitations

Secondbrain consultation read the index and queried independent review and
durable delivery. The prior request at
`/Users/marcio/dev/secondbrain/inbox/roundfix/_triaged/2026-08-25-review-source-unica-e-a-frota-desligando-o-coderabbit.md`
identified the same missing pre-PR path; current local inspection confirmed the
limitation without treating its earlier fleet counts as current measurements.

Exa located and read [OpenAI code review](https://learn.chatgpt.com/docs/code-review)
and [Claude headless execution](https://code.claude.com/docs/en/headless).
They establish native review and headless structured-output primitives; they
do not establish Roundfix integration, runtime access, or a paid allowance.
Local CLI help confirmed available flags, and a parser-only probe confirmed the
Codex target/custom-prompt conflict. No live reviewer was invoked for this PRD.

## Confirmed reviewer decision — 2026-09-08

The maintainer selected Codex as the default independent reviewer and requires
an explicit reviewer in `.roundfixrc.yml` to take precedence. Reuse the existing
`profiles.review` resolution instead of introducing another reviewer key or
forcing invocation flags that override the project. Current Project Config
selects Codex / gpt-5.6-luna / max with the declared Codex / gpt-5.6-sol / high
fallback; preserve that actual tuple. Built-in, User Config and Project Config
provenance remain visible. Invalid configuration or unavailable required review
capability is a named refusal, never a silent substitution with the default.

`review_source.name: coderabbit` is the legacy external PR-feedback provider,
not an Agent Selection Profile. The new native review must consume the review
profile; changing this planning record does not yet implement that adapter or
remove CodeRabbit. Review sessions are independent of implementation sessions,
and a changed candidate invalidates their evidence. The granted default-policy
decision does not approve otherwise proposed governed mutations or paid calls.

The [source ownership index](references/_index.md) records the pre-adoption
path, type, primary owner and current owned copy. Secondary consumers link
that copy; lifecycle completion means routing, not verified implementation.

The provider-opt-out evidence is owned in [the adopted Finding](references/2026-09-08-generated-review-rules-reactivate-a-disabled-provider.md). Its acceptance must prove that retired/disabled CodeRabbit never receives an automatic or manual request from generated guidance or the new delivery path.

## Technical candidate

The [_techspec.md](_techspec.md) records the reviewable implementation map,
coverage and build order. It is a proposed candidate, not a completed authoring
gate or permission to dispatch. Exact governed grants and the named decisions
remain pending; no Task Graph or implementation result is claimed.
