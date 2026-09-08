# Roundfix pending work plan — 2026-09-08

The maintainer requested a complete inventory and implementation plan, archival
of finished or deprecated documentation, refreshed model-selection references,
Roundfix Inbox triage, support planning for Claude Fable 5.1 and GPT 6 Astra,
removal of the operational CodeRabbit dependency, independent review before PR
creation, and canonical Spec-local authorization records. Commits and pushes
are authorized.

## Decisions confirmed in this session

- The future unattended workflow may implement, review in an independent
  session, open a pull request, and merge automatically after the Specs and
  their limits are approved. A merge requires passing review and required
  checks for the current commit. Failures outside the approved limits stay
  blocked; unattended execution does not authorize bypassing a gate.
- Relevant research must be retained in Secondbrain for later consultation.
- New Specs must contain their authorizations, and this placement must become
  a canonical Baseline contract.

## Decisions being collected

- Independent reviewer policy: one reviewer selected by profile, both Codex
  and Claude for each review, or Codex only.
- Paid API and subscription limits for implementation, review, and any live
  model-capability probes.
- Final architecture and exact protected paths for the changes that replace
  CodeRabbit and align the autonomous orchestration skills.

## Initial inventory

| Surface | Observed state | Treatment |
| --- | --- | --- |
| `docs/specs/` | No active Spec PRD; 103 Spec folders already under `docs/history/specs/` | Do not invent completed work or archive a Spec without its Task and QA proof. |
| `docs/backlog/` | Six entries | Check current implementation, remove duplicates through explicit routing, and promote remaining intent into Specs. |
| `docs/findings/` | Six Rollups and four other Findings | Preserve original evidence; record current closure or residual work in dated addenda and preserve archive licenses. |
| Secondbrain `inbox/roundfix/` | Twelve pending entries | Triage oldest first into exactly one Finding, Backlog Entry, or recorded discard per entry. |
| `docs/handoffs/` | Two historical handoffs | Both moved to `docs/history/handoffs/` with byte-identical contents under the maintainer's archival authorization. |

## Architectural facts established

CodeRabbit is not a required GitHub check for this repository. The inspected
main-branch protection requires `Verification gate` and `Validate PR title`.
The product nevertheless accepts only `review_source.name: coderabbit`; a
configuration-only switch to `none` is invalid. Removal needs an intentional
product contract, replacement tests, active documentation and skill updates.

The 2026-08-31 decision retained `implement-spec` and rejected adopting
`kickoff`. The inspected Fluxus skill contains useful queue and authorization
requirements, but also repeats a window mechanism now owned by Roundfix,
groups user questions, and assumes a disabled review mode that Roundfix does
not yet implement. Its autonomous delivery must be evaluated against the
Daemon's ownership of Task execution and Verification, not copied unchanged.

The current `Review Source` domain refers to external feedback on an open PR.
A review generated before a PR exists needs its own contract; it must not fake
a PR identifier to satisfy the old interfaces. Reviewer output must identify
the base and head commits, and a changed head invalidates the review.

The maintainer removed the mandatory personal prefix on 2026-09-08. New work
branches use purpose types (`feat/`, `fix/`, `refactor/`, and the other accepted
commit types); Roundfix Run/Task branches retain their existing `roundfix/run-`
namespace and derived names. This resolves the previously recorded naming
bootstrap conflict without a runtime code change or an execution exception.
The Supervisor still delegates feature code and tests through Roundfix. The
approved canonical/global instruction subset is recorded inside Spec 0125;
repository identity and reconciliation changes remain proposed.

## Progress and remaining decisions

- Inventory and triage completed: twelve original Inbox Entries processed,
  ten destination artifacts created, two evidence-backed discards, and no
  pending entry remaining in the Roundfix Inbox at that checkpoint.
- Two historical handoffs and the false-green Finding were archived. The
  remaining six Rollups and three original Findings retain their residuals;
  all 81 historical member licenses remain attached to active Rollups.
- A dedicated Backlog Entry records Fable 5.1 and Astra support. The refreshed
  model reference records 22 present and two absent IDs in the OpenRouter
  snapshot, API pricing overrides, exact version slugs, and four installed
  runtime/adapter versions. OpenCode advertises 422 IDs but neither new model;
  Codex/Claude ACP selection remains unmeasured. No paid probe was performed.
- Exa and local research for independent review and durable delivery is
  recorded in the design and in Secondbrain research captures.
- PRDs and proposed Spec-contained authorization records make the portfolio
  reviewable. TechSpecs, source adoption, and Task Graphs await the decisions
  that control their architecture, execution limits, and protected scope.

The CodeRabbit product integration is still active. Neither the planning
documents nor a passing PRD-stage checker disables it. No implementation Run,
paid model call, new PR, or merge has been started by this planning work.

This file records the cross-Spec plan. Individual Task state will live only in
the Task files; this plan does not certify implementation or QA completion.

### Validation checkpoint — 2026-09-08

`make verify-docs` passed after correcting three ADR attributions and recording
the relevant ADR scope in the PRDs. The full strict consistency check reports
nine Specs with zero findings. Each still reports fourteen omitted checks
because TechSpecs, Task Graphs, and adopted source indexes have not been
authored; no Task Verification ran. These are consistent planning artifacts,
not executable or approved Specs. The current checker does not enforce the
proposed authorization filename's approval state; Spec 0119 explicitly owns
that gap.

Forty-two local Markdown links were checked. Sixteen original observation
bodies were preserved, the two handoffs remained byte-identical, all twelve
Inbox source bodies survived triage, and all 81 Rollup licenses resolve.
The archive/triage commit `5e2c2cb7` was pushed to the existing `ma/` branch.

`go vet ./...` failed with the 31 existing copied-lock diagnostics now captured
in a Finding; no analyzer suppression or runtime change was made. A full
implementation test run is not claimed for this documentation-only work.

## Proposed Spec portfolio

The inventory found no active Spec numbered above 0118. The identifiers below
are reserved for this planning session; a reservation does not authorize
implementation. Splitting authorization from knowledge retention keeps two
independently testable contracts out of one oversized Spec.

| Spec | Outcome | Prerequisites |
| --- | --- | --- |
| 0119 — Spec-contained authorization | Each Spec carries its approved actions, exact protected paths, limits, and evidence of approval; the checker and executor distinguish approval from a proposal and preserve legacy grants. Verification execution observes that same trust boundary. | Maintainer approval of the concrete grant format and scope. |
| 0120 — Knowledge lifecycle and durable capture | One fleet Inbox, truthful terminal dispositions, complete history locations, clear autosync ownership, and an enforced prohibition on downstream evidence in upstream instructions. | 0119 for the authorization convention; current Secondbrain ownership verified. |
| 0121 — Complete Baseline decisions and regeneration | HTTP changes preserve all typed values; HTTP defaults have one owner; Greenfield refresh cannot deadlock on stale carriers; skill regeneration declares its outputs. | Exact protected file grant. |
| 0122 — Verified content and terminal settlement | Required executable files survive commit; Clean proves the resulting repository content; declared-only partial QA follows the already accepted policy consistently. | Exact protected file grant if an existing governed test changes. |
| 0123 — Runtime readiness and new models | Fable 5.1 and Astra are selectable only under advertised runtime capabilities and actual access evidence; runtime state is not copied with its mutex; fallbacks are proved when used. | Approved access/budget policy; current catalog research. |
| 0124 — Verification capacity and economics | Fresh controlled measurements explain gate cost and resource contention; unknown flakes retain evidence; gate/cache changes follow measurements and preserve negative controls. | 0123 mutex fix before enabling vet; tooling grant before gate changes. |
| 0125 — Repository identity and branch policy | Linked worktrees share repository identity; work branches express purpose while Run/Task branches retain the Roundfix namespace; legacy and squash reconciliation preserve unproved work. | Branch policy confirmed; repository-identity migration and reconciliation decisions pending. |
| 0126 — Independent review before a pull request | A local independent review certifies a specific candidate commit; bounded corrections trigger a fresh review; operative CodeRabbit dependencies are removed across product and owned skills. | Reviewer policy and concrete removal/migration grant. |
| 0127 — Durable unattended Spec delivery | A durable Supervisor drives approved Specs through implementation, independent review, QA, archive, PR, checks, and merge, stopping at exhausted limits or missing authority. | 0119, 0122, 0125, 0126; use 0123 readiness and the capacity evidence from 0124. |

Start with the authorization contract and the CodeRabbit replacement work;
finish the integrity/readiness prerequisites before enabling unattended
multi-Spec delivery. Independent implementation waves are permitted only when
their declared files and repository resources do not overlap. No release or
tag is implied by this queue.

### Source coverage

- All six original Backlog Entries have a residual owner: cache/gate economics
  in 0124, authorization and upstream-only rules in 0119/0120, and review
  contract protection in 0126.
- Ten of the twelve Inbox Entries become six Backlog Entries and four
  Findings; the other two record that their exact deleted-target reconcile
  incidents were already implemented by 0097. The present-target squash edge
  remains in 0125.
- The false-green Finding closes against 0083 with a fresh negative-control
  test. The three other ordinary Findings retain their dated measurements
  and route residuals into 0122–0127.
- The six Rollups retain the licenses of their 81 historical members. Their
  dated routing addenda distinguish shipped subfamilies from remaining work;
  a `done` absorption marker is never used as evidence of implementation.
- Fable 5.1 and Astra receive a dedicated Backlog Entry and are absorbed by
  0123 when its Spec is authored.

## Proposed authorization convention

Each new Spec should own `_authorization.md`, linked from the Project
Constraints in its PRD and TechSpec. It records scope, approval state, the
maintainer's decision and date, exact protected repository-relative paths,
sanctioned regeneration, permitted execution/publication actions, and limits.
The approval record lands in a separate commit before the authorized tooling
change. A proposed record is inspectable documentation and grants no action.

The checker must validate this file by its role, not by a date embedded in its
filename. Existing dated records under `docs/workflow/authorizations/` remain
valid historical evidence; no archived grant is rewritten merely to migrate
placement. Archive moves the Spec and its authorization together and repairs
active references through the existing resolver.

The confirmed through-merge policy belongs in each consuming Spec. Model
access, monetary limits, maximum correction attempts, and protected paths
still need concrete decisions. A permission denial, missing answer, expired
grant, changed candidate commit, failed reviewer, or failed check cannot be
converted to approval by timeout or by an empty response.

## Research and its effect

The local Secondbrain query on independent review and autonomous delivery
returned the fleet's earlier CodeRabbit-removal request and Roundfix's earlier
Spec portfolio. Reading those sources showed that QA-only delivery was an
explicit fleet gap and that several historical proposals had already shipped.
That evidence drives deduplication and a separate pre-PR review contract.

- Secondbrain: `inbox/roundfix/_triaged/2026-08-25-review-source-unica-e-a-frota-desligando-o-coderabbit.md`.
- Secondbrain: `projects/roundfix/mirror/docs/design/2026-08-12-the-spec-set-the-evidence-asks-for.md`.
- [Codex CLI reference](https://developers.openai.com/codex/cli/reference.md),
  read through Exa: a noninteractive review can target a local branch or
  commit. Its mutually exclusive target/prompt arguments constrain the adapter.
- [Codex code review](https://developers.openai.com/codex/code-review),
  read through Exa: the review flow produces prioritized findings without
  editing the worktree. A runtime contract still needs to prove completion
  and attach the reviewed commit.
- [Claude headless mode](https://code.claude.com/docs/en/headless.md), read
  through Exa: print mode and structured JSON output support a local review
  adapter. Using it for a pre-PR diff is a design inference, not a claim that
  Claude publishes a drop-in CodeRabbit replacement.

These sources establish interfaces and prior experience. They do not prove
that Roundfix already implements the proposed review or queue, that every
installed adapter exposes every new model, or that unknown CI flakes share
one root cause.

## Branch-policy decision — 2026-09-08

The mandatory personal prefix is removed. `refactor/` is used rather than
`refact/` to match the repository's commit vocabulary. The compatible
`branch.prefix` decision now defaults to the placeholder pattern `<type>/`;
this repository's stored value is changed through the public Baseline workflow.
Legacy manifests can still be read, but their personal values do not override
the new canonical purpose rule. Existing remote branches and original evidence
are retained. The current work continues on `feat/purpose-based-branch-policy`.
The owned `implement-spec` skill and its shipped copy also replace their former
personal-prefix guardrail; version 0.0.3 delegates naming to this policy and
preserves the Roundfix Run/Task namespace.

The prior claim that a naming exception was needed to bootstrap a Run is no
longer current. Reviewer choice, resource limits, the broader tooling grants,
and unfinished TechSpecs/Task Graphs remain pending; this naming decision does
not answer those separate questions.

The retired Run-prefix Finding moves to `docs/history/findings/` with its
original observation preserved and a dated policy-disposition addendum. The
active Run lifecycle Rollup licenses it as a fourteenth member, in addition to
the 81 historical members preserved at the earlier portfolio checkpoint.
The Secondbrain triage pointer follows its history destination.
