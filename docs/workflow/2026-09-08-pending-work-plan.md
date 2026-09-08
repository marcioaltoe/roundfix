# Roundfix implementation portfolio — 2026-09-08

The complete triage assigns 23 ordinary sources to 11 Specs: 14 Backlog Entries
and nine Findings. The six Rollups remain as the license roots of 82 archived
members. Their dated dispositions distinguish shipped behavior, residual work
and superseded proposals; routing completion does not claim implementation.

Every ordinary source has exactly one primary owner and a Spec-local copy.
Secondary consumers link that copy. There are no remaining ordinary entries in
`docs/backlog/` or `docs/findings/`; the latter retains the six licensed Rollups.
The two newly discovered destination Inbox captures, from Fiscus and Pantheon,
also have owners. Secondbrain triage preserves original capture bodies and
changes their routing metadata only.

The portfolio contains PRDs, proposed technical candidates, owned source
indexes where applicable, and exact proposed authorization records. It has no
executable Task Graphs yet. No proposed grant, checker omission or unanswered
question authorizes tooling changes or live paid model calls.

## Confirmed policy and delivered foundation

- Codex is the default independent reviewer. An explicit `.roundfixrc.yml`
  review profile wins through existing `profiles.review`/`CategoryReview`
  resolution. This repository currently selects Codex / gpt-5.6-luna / max,
  with its declared Codex / gpt-5.6-sol / high fallback. Do not inject launcher
  flags that override that choice. Invalid configuration or unsupported review
  capability is an explicit refusal. [ADR-0151](../adr/0151-configured-review-profiles-select-the-independent-reviewer.md)
  records the confirmed decision.
- The approved delivery destination is automatic squash merge after independent
  review and required checks pass for the current candidate. Routine actions
  already granted do not need another prompt. Failed, skipped, stale, empty or
  absent evidence does not approve publication.
- Work branches use purpose prefixes such as `feat/`, `fix/` and `refactor/`;
  Run/Task branches retain their Roundfix namespace. The personal `ma/` rule is
  removed. This change already landed in PR177 and is not pending implementation
  in Spec 0125.
- PR176 and PR177 landed the initial research/capture/question rules, portfolio
  PRDs and branch-policy foundation. Their independent reviews and CI evidence
  concern those delivered changes, not these proposed product features.
- The model reference has already been refreshed with the observed provider
  IDs, prices and installed adapter versions. Availability of Fable 5.1 or
  Astra through a particular ACP runtime remains unproved; Spec 0123 owns the
  capability/access work. No paid probe is part of this triage.

## Current autonomy and the proposed path

Today `roundfix implement --spec <slug> --detach` can execute one ready Spec
Task Graph under the Daemon after its initiating session exits. Existing
process identity, locks, journal/events and Run Window are useful primitives.
They do not implement a durable all-Spec queue. `implement.auto_push` does not
create a PR or provide native independent pre-PR review. The legacy
`review_source.name: coderabbit` remains operational product configuration;
it is not the reviewer profile and has not been disabled by these documents.

The inspected Fluxus kickoff is a chat-owned coordinator. Reuse its inventory
and preparation ideas through the existing `implement-spec` entry point; do
not copy its second window script, bundled questions, no-review path or
conflicting implementation-loop instruction. No kickoff script was invoked.

The technical candidate in Spec 0127 proposes a small queue owner in the
existing Daemon with SQLite queue/item/action records. It runs one whole Spec
at a time, records external action intent and receipts, and reconciles observed
push/PR/merge outcomes before retrying after a restart. The intended sequence
is Implement and terminal QA, archive/commit, independent review, PR, checks,
squash merge, then revalidate the next Spec. The archive-first late-correction
trade-off remains a named decision; it does not reopen archived evidence or
implicitly authorize a new corrective Spec.

A newly captured static gap belongs to Spec 0127: production Implement warns
or renders `budget.max_run_duration`, but the inspected execution path does
not apply that deadline. The YAML value therefore cannot prove a bounded
unattended Run. Enforced cancellation/recovery must be tested before enabling
the durable queue. This observation is static evidence, not a timed live repro.

## Implementation order

All rows are authoring candidates. Dependencies identify contracts that must
exist before their consumers execute. Shared skills and canonical files require
serial integration and refreshed plans; a separate package is not evidence
that two whole-Spec deliveries can safely run concurrently.

| Order | Spec and result | Required predecessors / boundary |
| --- | --- | --- |
| 1 | [0119 — Spec-contained authority](../specs/0119-spec-contained-authorization/_techspec.md): approval by role, committed scope and command-source trust | Concrete initial grant; preserve valid legacy grants. |
| 2 | [0123 — Runtime readiness and models](../specs/0123-runtime-readiness-and-model-capabilities/_techspec.md): real capability/access evidence, profile provenance and pointer-owned mutex state | 0119 convention; access/spending policy; explicit ADR decision before changing fallback timing. |
| 3 | [0126 — Native independent review](../specs/0126-agent-review-before-pull-request/_techspec.md): configured reviewer, candidate receipts and operational CodeRabbit retirement | 0119, 0123; migration/adapter grant and finite correction limits. |
| 4 | [0122 — Verified content and settlement](../specs/0122-verified-content-and-terminal-settlement/_techspec.md): executable-source proof, consistent QA eligibility, authorized repair entry and independent Verification groups | 0119; retain Daemon ownership and original failure evidence. |
| 5 | [0121 — Complete Baseline decisions](../specs/0121-baseline-decisions-and-complete-regeneration/_techspec.md): typed edits, full/incremental Verification, owned outputs and evidence-based external lock reconciliation | 0119; exact generator/test scope and reviewed migration. |
| 6 | [0129 — Spec authoring and gate recovery](../specs/0129-spec-authoring-and-gate-recovery/_techspec.md): API/metric/ADR evidence, premise amendments, temporal dependencies and supported stale-QA recovery | 0119, 0121, 0122; consume 0126 semantic review. Characterize the safe recovery transition before executable Tasks. |
| 7 | [0125 — Repository identity and reconciliation](../specs/0125-repository-identity-and-run-branch-policy/_techspec.md): linked-worktree identity, legacy aliases and present-target/squash proof | 0119, 0122; purpose naming and deleted-target repair are already delivered. |
| 8 | [0124 — Measured Verification economics](../specs/0124-verification-capacity-and-measured-economics/_techspec.md): bounded comparative measurements, preserved failures and justified gate changes | 0123 mutex repair before analyzer gate; measure the operative gate after 0121/0122. |
| 9 | [0120 — Knowledge lifecycle](../specs/0120-knowledge-lifecycle-and-durable-capture/_techspec.md): capture durability, terminal dispositions, licenses, upstream boundary and stable Review Artifact retirement | 0119; 0126 receipts for the proposed retirement revision. ADR-0123 remains operative until proposed ADR-0152 is accepted. |
| 10 | [0128 — Bare stable release tags](../specs/0128-release-planning-with-bare-stable-tags/_techspec.md): strict parsing, exact ref identity and explicit ambiguity | 0119 convention; serialize the shared Roundfix skill after earlier edits. Read-only planning, no release/tag creation. |
| 11 | [0127 — Durable unattended delivery](../specs/0127-durable-unattended-spec-workflow/_techspec.md): enforced Run deadline, queue owner, receipts, recovery and through-merge delivery | 0119, 0122, 0125, 0126, 0129; runtime readiness and measured capacity from 0123/0124. Revalidate the whole approved portfolio before start. |

The existing per-Spec Implement path is the bootstrap executor. The new queue
cannot implement its own prerequisites before it exists. Before each dispatch,
author the approved Task Graph from that Spec's build order, validate the
operative source/grant and current prerequisite commits, and retain terminal
QA. No supervisor-authored feature code or tests are part of this plan.

## Source ownership

Each linked copy preserves its original observation and adds dated routing.
A `promoted` Backlog Entry or `done` Finding in the Spec records adoption, not
successful implementation. Pre-adoption paths and dates are in each owning
`references/_index.md`.

| Source | Type | Primary Spec |
| --- | --- | --- |
| [2026-08-26-the-baseline-does-not-name-a-home-for-authorization-records.md](../specs/0119-spec-contained-authorization/references/2026-08-26-the-baseline-does-not-name-a-home-for-authorization-records.md) | backlog | 0119 |
| [2026-09-08-verification-command-execution-has-an-explicit-trust-boundary.md](../specs/0119-spec-contained-authorization/references/2026-09-08-verification-command-execution-has-an-explicit-trust-boundary.md) | backlog | 0119 |
| [2026-08-26-the-upstream-only-knowledge-rule-has-no-detector.md](../specs/0120-knowledge-lifecycle-and-durable-capture/references/2026-08-26-the-upstream-only-knowledge-rule-has-no-detector.md) | backlog | 0120 |
| [2026-09-08-capture-inbox-contract-matches-the-fleet-door.md](../specs/0120-knowledge-lifecycle-and-durable-capture/references/2026-09-08-capture-inbox-contract-matches-the-fleet-door.md) | backlog | 0120 |
| [2026-09-08-completed-work-has-a-terminal-state-without-a-spec.md](../specs/0120-knowledge-lifecycle-and-durable-capture/references/2026-09-08-completed-work-has-a-terminal-state-without-a-spec.md) | backlog | 0120 |
| [2026-09-08-history-guide-defines-every-retired-family.md](../specs/0120-knowledge-lifecycle-and-durable-capture/references/2026-09-08-history-guide-defines-every-retired-family.md) | backlog | 0120 |
| [2026-09-08-secondbrain-capture-has-one-durability-owner.md](../specs/0120-knowledge-lifecycle-and-durable-capture/references/2026-09-08-secondbrain-capture-has-one-durability-owner.md) | backlog | 0120 |
| [2026-09-08-profile-omits-the-required-incremental-verification.md](../specs/0121-baseline-decisions-and-complete-regeneration/references/2026-09-08-profile-omits-the-required-incremental-verification.md) | finding | 0121 |
| [2026-09-08-skill-regeneration-declares-its-owned-outputs.md](../specs/0121-baseline-decisions-and-complete-regeneration/references/2026-09-08-skill-regeneration-declares-its-owned-outputs.md) | backlog | 0121 |
| [2026-09-08-verified-executable-source-is-dropped-before-settlement.md](../specs/0122-verified-content-and-terminal-settlement/references/2026-09-08-verified-executable-source-is-dropped-before-settlement.md) | finding | 0122 |
| [2026-09-08-acpx-runner-value-receivers-copy-its-mutex.md](../specs/0123-runtime-readiness-and-model-capabilities/references/2026-09-08-acpx-runner-value-receivers-copy-its-mutex.md) | finding | 0123 |
| [2026-09-08-support-claude-fable-5-1-and-gpt-6-astra.md](../specs/0123-runtime-readiness-and-model-capabilities/references/2026-09-08-support-claude-fable-5-1-and-gpt-6-astra.md) | backlog | 0123 |
| [2026-08-08-go-clean-testcache-clears-a-cache-the-gate-does-not-use.md](../specs/0124-verification-capacity-and-measured-economics/references/2026-08-08-go-clean-testcache-clears-a-cache-the-gate-does-not-use.md) | backlog | 0124 |
| [2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs.md](../specs/0124-verification-capacity-and-measured-economics/references/2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs.md) | backlog | 0124 |
| [2026-08-11-a-git-worktree-that-fails-only-under-load.md](../specs/0124-verification-capacity-and-measured-economics/references/2026-08-11-a-git-worktree-that-fails-only-under-load.md) | finding | 0124 |
| [2026-08-15-the-gate-runs-a-saturated-suite-inside-a-dense-run.md](../specs/0124-verification-capacity-and-measured-economics/references/2026-08-15-the-gate-runs-a-saturated-suite-inside-a-dense-run.md) | backlog | 0124 |
| [2026-09-08-force-stop-legacy-owner-has-an-unexplained-ci-failure.md](../specs/0124-verification-capacity-and-measured-economics/references/2026-09-08-force-stop-legacy-owner-has-an-unexplained-ci-failure.md) | finding | 0124 |
| [2026-08-12-three-consecutive-specs-measure-the-loop.md](../specs/0126-agent-review-before-pull-request/references/2026-08-12-three-consecutive-specs-measure-the-loop.md) | finding | 0126 |
| [2026-08-31-the-review-agent-rewrites-the-contract-it-was-asked-to-satisfy.md](../specs/0126-agent-review-before-pull-request/references/2026-08-31-the-review-agent-rewrites-the-contract-it-was-asked-to-satisfy.md) | backlog | 0126 |
| [2026-09-08-generated-review-rules-reactivate-a-disabled-provider.md](../specs/0126-agent-review-before-pull-request/references/2026-09-08-generated-review-rules-reactivate-a-disabled-provider.md) | finding | 0126 |
| [2026-09-08-implement-does-not-apply-the-declared-run-duration.md](../specs/0127-durable-unattended-spec-workflow/references/2026-09-08-implement-does-not-apply-the-declared-run-duration.md) | finding | 0127 |
| [2026-09-08-plan-releases-with-bare-stable-tags.md](../specs/0128-release-planning-with-bare-stable-tags/references/2026-09-08-plan-releases-with-bare-stable-tags.md) | backlog | 0128 |
| [2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md](../specs/0129-spec-authoring-and-gate-recovery/references/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md) | finding | 0129 |

Spec 0125 has no ordinary source adoption index: its residuals come from the
shared Run lifecycle Rollup, whose historical members retain their existing
license. Do not create a duplicate owned copy solely to populate a directory.

## Historical Rollup disposition

The six active license roots contain dated member-family dispositions and
primary/secondary owners. All 82 unique members remain licensed exactly once.

| License root | Residual owner(s) | Already delivered or explicitly excluded |
| --- | --- | --- |
| [Agent selection](../findings/2026-08-06-rollup-agent-selection-and-execution-environments.md) | 0123; 0124 capacity; 0121 lock reconciliation | Specs 0052/0056/0088/0089/0091/0103 provide existing adapters/profiles/capabilities/compiled fixtures. Catalog presence is not new-model access proof. |
| [Baseline and tooling](../findings/2026-08-06-rollup-baseline-and-derived-tooling.md) | 0121; 0119 authority; 0124 measurements | Specs 0054/0057/0061/0062/0067/0082/0084 provide the current Baseline foundation. Residuals concern specific typed edits/ownership/evidence gaps. |
| [QA and Verification](../findings/2026-08-06-rollup-qa-gates-and-verification-evidence.md) | 0122; 0124 economics; 0129 authoring; 0128 bare tags | 0083 truthful gate, 0095/0116 authoring checks, 0097 bootstrap/collision/deleted-target proof, 0105 characterization/equivalent QA evidence already shipped. Do not recreate them from older observations. |
| [Review and delivery](../findings/2026-08-06-rollup-review-and-delivery-convergence.md) | 0126; 0127 queue; 0120 retirement | Historical dirty-review/broad skip proposals conflict with current ADR-0139/0141; merge without independent review is superseded by the maintainer's decision. |
| [Run lifecycle](../findings/2026-08-06-rollup-run-lifecycle-and-branch-integrity.md) | 0125; 0124 capacity; 0127 durable owner | 0055/0059/0066/0068/0081/0103 already implement process/store/journal protections. Existing write mutex, single connection and busy timeout do not justify speculative HOME isolation. Personal branch prefix is already retired. |
| [Spec authoring](../findings/2026-08-06-rollup-spec-authoring-and-contract-enforcement.md) | 0129; 0119/0120/0122/0127 consume the relevant contracts | 0092/0096/0098/0118 already provide recovery/carry-forward foundations. Preserve prior Results, current Task Type and corrective ceilings. |

Optional historical proposals about fixture deadlines, external evidence
records, fragile literals, test mass and toolchain strategy receive a measured
accept/defer/reject disposition in Spec 0124. Adoption promises that each is
assessed against current evidence, not that every old suggestion must become
code. Unexplained CI/load failures remain unexplained until reproduced or
settled with evidence; a successful retry does not supply a root cause.

## Approval package and remaining decisions

Each Spec's authorization proposal lists exact governed files. None is a grant.
The current checker does not fully enforce the proposed undated filename's
approval state; Spec 0119 owns that gap, and a clean checker cannot bypass it.
The approved record must be independently merged into target ancestry before
its consuming tooling candidate. A separate commit inside the same PR would
be flattened by squash and does not preserve the required prior-authority
boundary. Recheck the final merged history; proposal-only PRs grant nothing. Existing confirmed reviewer, branch and through-merge choices remain
settled and need no repeated permission.

| Spec | Exact proposed scope | Governed paths |
| --- | --- | --- |
| 0119 | [_authorization.md](../specs/0119-spec-contained-authorization/_authorization.md) | 21 |
| 0120 | [_authorization.md](../specs/0120-knowledge-lifecycle-and-durable-capture/_authorization.md) | 12 |
| 0121 | [_authorization.md](../specs/0121-baseline-decisions-and-complete-regeneration/_authorization.md) | 20 |
| 0122 | [_authorization.md](../specs/0122-verified-content-and-terminal-settlement/_authorization.md) | 19 |
| 0123 | [_authorization.md](../specs/0123-runtime-readiness-and-model-capabilities/_authorization.md) | 2 |
| 0124 | [_authorization.md](../specs/0124-verification-capacity-and-measured-economics/_authorization.md) | 2 |
| 0125 | [_authorization.md](../specs/0125-repository-identity-and-run-branch-policy/_authorization.md) | 2 |
| 0126 | [_authorization.md](../specs/0126-agent-review-before-pull-request/_authorization.md) | 14 |
| 0127 | [_authorization.md](../specs/0127-durable-unattended-spec-workflow/_authorization.md) | 12 |
| 0128 | [_authorization.md](../specs/0128-release-planning-with-bare-stable-tags/_authorization.md) | 3 |
| 0129 | [_authorization.md](../specs/0129-spec-authoring-and-gate-recovery/_authorization.md) | 14 |

Pending decisions controlling execution:

1. **Consumption:** the current question asks whether to use only existing
   subscriptions or permit metered APIs with an explicit per-Run/day ceiling.
   No answer is recorded. Real model access/review exercises wait for sufficient
   authority; public catalog research and docs-only verification can proceed.
2. **Finite operation limits:** proposed serial whole-Spec delivery and at most
   two review correction cycles, with explicit queue duration/spending ceilings.
   Existing Run Window and Run Budget keep their separate meanings. The window
   controls starts; a Run already started follows its enforced budget. These
   candidate values/semantics are not inferred user answers.
3. **Bounded protected scope and trust:** approve the linked exact file lists,
   the Spec-contained grant convention and authored-command source boundary.
   Shared files must be revalidated after earlier changes; an expanded path or
   action needs its own authority. Sanctioned regeneration follows source
   approval, and generated guides use the public Baseline plan/apply path.
4. **Named architectural trade-offs:** archive-first final review and a separately
   authorized corrective Spec after a late blocker; stable-evidence retirement
   in proposed ADR-0152; and fallback readiness timing in Spec 0123. Existing
   accepted ADRs remain operative while these revisions are proposals.

No release, tag, deployment, destructive history cleanup, foreign-repository
mutation, scheduler installation or unlimited paid use is bundled into ordinary
queue delivery. A missing answer becomes a named blocker only for the action
that depends on it. It does not stop source triage, reviewable design or the
already authorized publication of these planning documents.

## Research and its effect

Local Secondbrain consultation used the index/qmd and the concepts
`agent-harnesses` and `agent-workflows-e-loop-engineering`, plus the destination
Inbox and prior Roundfix mirror. They support durable ownership, independent
checks, receipts and visible stops, and reveal shipped work that must not be
reimplemented. The Fiscus incremental-declaration and Pantheon provider-opt-out
captures supply concrete residuals assigned to 0121 and 0126.

Exa located and read these external sources:

- [Codex noninteractive execution](https://developers.openai.com/codex/noninteractive)
  and [CLI reference](https://developers.openai.com/codex/cli/reference) describe
  local execution/review surfaces. They inform the configured reviewer adapter;
  they do not prove complete structured review evidence in Roundfix.
- [Claude headless execution](https://code.claude.com/docs/en/headless) describes
  print/structured output and permission controls. Using it as an independent
  local-candidate reviewer is a proposed integration, not a demonstrated
  drop-in replacement or a promise of access through every ACP adapter.
- [CodeRabbit auto-review configuration](https://docs.coderabbit.ai/configuration/auto-review)
  and [review commands](https://docs.coderabbit.ai/reference/review-commands)
  confirm that disabling automatic review still permits manual requests.
  This supports the Pantheon finding: provider opt-out must not be undone by
  generated manual-review obligations. No provider request was triggered.

Local `codex exec --help`, Roundfix help/configuration and package inspection
were compared with those sources. No external query contained private code,
credentials or client records. The reusable research is captured in Secondbrain
at `inbox/secondbrain/2026-09-08-configured-reviewer-and-autonomous-queue-readiness.md`
for future ingestion. Source/code inspection is not live runtime or recovery
validation; model access, priced execution and unknown flake causes remain
explicit limitations.

## Validation scope

This plan records triage and technical candidates only. Its completion evidence
must include source/index uniqueness, preservation of original bodies, all
82 license targets, valid local links and the repository's documentation gate.
Stage-scoped checks cannot establish Task Graph readiness, approved grants,
implementation, terminal QA or a working durable queue. Current verification
results belong in the delivery record after the candidate tree is frozen.
