---
spec: 0120-knowledge-lifecycle-and-durable-capture
prd: _prd.md
created: 2026-09-08
---

# Knowledge lifecycle, publication evidence and review retirement — Technical candidate

## Cleanup repair authorization — 2026-09-09

The maintainer approved the five compatibility repairs needed to deliver the
documentation cleanup. That bounded repair completed implementation and QA and
was archived as [Spec 0130](../../history/specs/0130-documentation-cleanup-compatibility/_prd.md).
Its [authorization](../../history/specs/0130-documentation-cleanup-compatibility/_authorization.md)
supersedes the former cleanup proposal in this folder. This does not approve or
complete the remaining lifecycle and source-adoption work in Spec 0120.

## Executive Summary

Extend existing lifecycle readers and canonical capture rules instead of adding another Inbox or scheduler. The trade-off is retaining explicit unknown lifecycle outcomes until durable evidence exists. A Review Artifact relocation is a path migration and must not decide whether its review is active.

This document makes the proposed work concrete for review. Authoring remains
open: the exact governed grant and the decisions named below are pending.
It is not an implementation-ready TechSpec, an executable Task Graph or
approval to mutate protected files. The remaining decisions are recorded in
[_prd.md](_prd.md) and [_authorization.md](_authorization.md); the dependencies
below define this Spec's place in the implementation order.

## Project Constraints

- Identifier strategy: applicable — retain existing Spec slugs, dated artifact basenames, Inbox destinations, and source provenance; a terminal disposition does not invent an implementation identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, authentication, or HTTP API change is proposed. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — preserve the one-way mirror and fleet Inbox boundary, the existing History Root, and accepted regeneration ownership under ADR-0149. Source: `docs/agents/docs-layout.md`.
  ADR-0083 applies to the source-adoption validation follow-up: adopted Findings and Backlog Entries move to one owning Spec, with no copy or stub left at the original path.
  ADR-0116 applies to cited-policy evidence: preserve checks that compare a claim with the cited record; source-origin absence is an additional structural check, not a substitute for semantic citation validation.
  ADR-0123 remains operative: retirement currently uses conservative local reachability and the Review Artifact resolver never writes into history. The proposed stable-evidence change does not yet supersede it.
  ADR-0152 is a proposed revision only, recorded for review; it creates no current obligation or grant.
- Tooling authority: applicable — express maintainer authorization on 2026-09-08 covers the terminal lifecycle instruction through [the narrow grant](references/2026-09-08-terminal-lifecycle-authorization.md); bounded files: `internal/baseline/assets/modules/context-workflow.json`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, with sanctioned digest regeneration. The broader [_authorization.md](_authorization.md) remains proposed, including all Go, test, companion-repository and unrelated skill changes. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Capture contract | `canonical context-workflow/secondbrain modules and companion Secondbrain Inbox contract` | Name one fleet door and distinguish local capture, local commit and observed remote publication. |
| Terminal disposition | `internal/spec/retirement.go and internal/speccheck/backlog.go` | Represent evidence-backed fulfilled or superseded intent without misclassifying it as declined. |
| History and licenses | `internal/spec/archive.go and internal/speccheck/citations.go` | Resolve all retired families and replacement license chains while preserving old records. |
| Upstream citation check | `internal/docscontract/publicdocs_test.go` | Reject evidentiary citations from glossary/guides into active or archived downstream cases. |
| Review Artifact retirement | `existing Review Artifact resolver and retirement readers in internal/spec` | Separate stable lifecycle proof from legacy path relocation and local object availability. |

The map extends current package owners. Paths that name a package are
implementation seams, not permission for arbitrary edits below that directory.
The exact governed files remain in the authorization proposal; ordinary source
changes must stay within this Spec's behavior. Revalidate shared files after
prerequisite Specs land rather than replacing their newer contracts.

## Implementation Design

### Capture and publication

Keep `inbox/<destination>/` in Secondbrain as the fleet door and the project as
the owner of its triage. A capture receipt distinguishes a created file, a local
commit and confirmed remote visibility. Use the existing publication owner;
record a bounded session-owned fallback only when that scheduler is absent.
The installed scheduler's last exit and mere presence cannot certify a new
capture's remote durability. Companion changes belong to Secondbrain's own
branch and authorization, never its mirrors or immutable raw sources.

### Dispositions and archive ownership

Extend the existing backlog/finding readers with a reasoned terminal disposition
for directly fulfilled and superseded work. Preserve old `done` absorption
meaning: a Spec owner is not implementation proof. Validate the evidence/ref
that supports terminal disposition and reject cycles, missing replacements,
path escape or claims resting only on a status label. History resolution covers
handoffs as a separate family and preserves basename/provenance on moves.

The six terminal Rollups were moved to `docs/history/findings/` during the
maintainer-directed cleanup on 2026-09-08. Their 82 members now have direct
active/archived Spec absorbers, and each archived Rollup has its own absorber.
Original observations and previous routing remain recorded. This concrete
migration uses the existing Spec-slug resolution contract; it does not require
support for archived Rollup basenames or assert pending implementation passed.

### Review retirement and upstream boundaries

ADR-0123 remains operative until the proposed replacement in ADR-0152 is
approved. The new retirement result records stable provider/recorded evidence,
including a squash outcome, or an explicit unknown reason. Fetching or pruning
a local Git object must not change the lifecycle verdict. Move a legacy
`docs/specs/_reviews` path independently of that verdict, retaining its active
or unknown state. Consume the local-candidate receipts from0126 when available;
missing historical receipts do not become successful evidence.

The citation detector reads the glossary and active agent guides, distinguishes
literal layout declarations from actual evidence links, and rejects active and
history Spec/Finding references as upstream authority. Do not exempt the entire
layout guide: it can contain a real forbidden citation. Keep this detector in
the existing docscontract boundary with positive and negative fixtures.

### Interfaces

The component responsibilities above define the boundary inputs and outcomes.
Preserve the existing public command, error, persistence and owner contracts
unless this design explicitly proposes their revision. Refusal occurs before
the dependent mutation and reports the missing fact; empty or absent evidence
cannot become a successful terminal result.

### Data Models

Retain source provenance, terminal disposition and absorber/license identity. A capture durability receipt identifies the observed durable destination; local file creation alone cannot attest remote storage.

### API Contracts

Extend existing capture, lifecycle and consistency boundaries. Any retirement or migration command must name its exact source/destination and refusal conditions; scheduler installation and external automation remain separately scoped.

## Coverage Map

- PRD Goal 1 → Capture contract.
- PRD Goal 2 → Terminal disposition.
- PRD Goal 3 → History and licenses, Review Artifact retirement.
- PRD Goal 4 → Upstream citation check.
- Core Feature 1 → Capture contract.
- Core Feature 2 → Capture contract.
- Core Feature 3 → Terminal disposition.
- Core Feature 4 → History and licenses.
- Core Feature 5 → History and licenses.
- Core Feature 6 → Upstream citation check.
- Core Feature 7 → Capture contract.
- Core Feature 8 → Review Artifact retirement.

The Testing Approach below describes the observations that must settle these
contracts. Task IDs and actual evidence are deliberately not invented during
proposal authoring; approved Task decomposition must assign every contract and
success metric before execution.

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

1. Capture states file/commit/remote-confirmed and scheduler-absent fallback; no mirror or raw writes.
2. Fulfilled, superseded, unresolved and declined intent with actual disposition evidence; merely changing status cannot close a defect.
3. Every one of the 82 original member relationships remains traceable after direct-Spec migration; all 88 resulting member/Rollup absorbers resolve. Standalone terminal closure requires declared reason/evidence, while an invalid existing pointer remains an error.
4. The same Review Artifact classifies identically before/after Git object availability changes; legacy relocation preserves liveness and exact observations.
5. Actual downstream evidence citations in active/history trees fail; directory-layout examples and accepted ADR links pass.

These are planned checks, not executed evidence. Exact commands, independent
groups where supported, required runtime access and failure expectations must
be authored after approval. Preserve the repository's declared Go/toolchain and
CI constraints; live paid calls require the separately recorded limit.

## Build Order

1. Capture/publication contract and companion changes (depends on: none).
2. Evidence-backed terminal dispositions (depends on: 1).
3. History-family and replacement-license resolution (depends on: 2).
4. Stable Review Artifact retirement and separate legacy relocation (depends on: 2, 3).
5. Upstream citation detector and canonical guidance (depends on: 1, 3, 4).
6. Cross-repository examples, migration negatives and terminal QA (depends on: 1, 2, 3, 4, 5).

This sequence is a proposed build order, not `_tasks.md`. Prerequisite merges,
exact governed authority and remaining decisions must be revalidated before it
becomes a Task Graph. Implementation and Verification remain runtime/Daemon
owned; the Supervisor authors and coordinates.

## Risks & Considerations

Retirement revises an accepted policy and is pending that decision. A source move must not erase evidence or leave a dangling archive license. Cross-repository publication must retain the actual owner and observed receipt rather than treating either checkout as the other.

## Decisions

- The maintainer selected complete source triage and the implementation portfolio; source ownership is now recorded. That intent is distinct from a concrete governed-file grant.
- Delivery through squash merge requires the configured pre-PR review policy outcome and passing required checks for the current candidate. Explicit none records intentional review omission; enabled-provider failure cannot select none. Releases, tags and paid consumption are not implied.
- Preserve configured reviewer selection; this Spec introduces no separate reviewer override.
- The proposed mechanisms and unresolved trade-offs above remain candidates. Existing accepted ADRs named in Project Constraints remain operative until any explicit revision is accepted.
- [Proposed ADR-0152](../../adr/0152-review-artifact-retirement-uses-stable-evidence.md) challenges the object-dependent retirement rule in ADR-0123; it does not supersede that rule yet.

## Cross-Spec dependencies

Required predecessor contracts: [0119](../0119-spec-contained-authorization/_techspec.md), [0126](../0126-agent-review-before-pull-request/_techspec.md).
Shared skills and canonical files require serial integration and revalidation
after predecessor changes. A predecessor reference is not an execution grant.

## Confirmed terminal lifecycle scope — 2026-09-08

The maintainer explicitly requires terminal Findings, Backlog Entries and
Rollups to leave the active family directories. The narrow canonical source
and generated-output authority is recorded in
[the terminal lifecycle grant](references/2026-09-08-terminal-lifecycle-authorization.md);
it does not grant the broader Go/test work proposed by this Spec.

Spec-owned adopted references remain with their consuming Spec. A standalone
terminal Finding without an actual absorber records `closure_reason` and
`closure_evidence`; it must not invent a Spec/Rollup just to enter history.
An invalid existing `absorbed_by` cannot be hidden by these fields. Mechanical
support for this explicit closure path remains required: the current checker
still demands an absorber for every archived Finding. The completed six-Rollup
migration has genuine Spec absorbers and already passes that existing check.


## Source adoption validation gap — 2026-09-08

The maintainer reported adopted Findings/Backlog Entries remaining in their
active source directories after entering Spec references. This is contrary to
[ADR-0083](../../adr/0083-adopted-sources-move-to-their-owning-spec.md) and
`write-prd` step 5: adoption is one move, never a copy or stub. The index's
`source` is an immutable provenance string, not a required existing file.
A background citation alone is not adoption and does not transfer ownership.

The current checkout contains no active Finding or Backlog Markdown files;
only the directory placeholders remain. A scan of both `docs/specs/*/references/_index.md`
and `docs/history/specs/*/references/_index.md` found 37 indexed Finding/Backlog
sources: all 37 destinations exist and all 37 original paths are absent.
No other occurrence was found within that indexed Roundfix corpus. This does
not establish the state of another project, branch or unindexed source.

### Reproduction and observed outcomes

A temporary Git fixture outside the repository exercised the public command:

```sh
rtk proxy /tmp/roundfix-terminal-lifecycle spec check 0090-a-gate-that-could-have-failed --strict --format json
```

The binary was built from the inspected source with Go 1.26.7. The fixture used
the citation fixture's guides and ADR-0083, a synthetic documentation-only PRD
with the four Project Constraints, and a valid five-column reference index.
The row named `docs/findings/2026-09-08-adoption-probe.md` or its Backlog
counterpart as `source`, owner `0090`, date `2026-09-08`, and the same basename
as the destination. Each case changed only that indexed source/destination and
its lifecycle metadata. Initial unrelated citation-fixture claims were removed
from this synthetic PRD before the controlled comparisons below.

| Case | Actual exit | Relevant diagnostic | Expected contract |
| --- | --- | --- | --- |
| Finding moved; original absent | 0 | none | accept |
| Finding destination and identical `done` original both present | 0 | none | reject duplicate original |
| Promoted Backlog moved; original absent | 0 | none | accept |
| Promoted Backlog destination; original remains `open` | 0 | none | reject duplicate original |
| Promoted Backlog destination and promoted original both present | 1 | SC-BACKLOG-UNMOVED | reject |
| Indexed destination missing | 1 | SC-REF-UNRESOLVED | reject |

The reproduction assertion fails because the two duplicate-origin cases return
zero. Positive controls and the existing missing-destination/promoted-Backlog
negative controls behave as expected. The reports contain 19 or 20 skips for
unrelated absent fixture surfaces; no Verification command or Task was executed.
These are detector observations, not terminal Spec QA or a fixed implementation.

### Cause and bounded follow-up

`internal/speccheck/citations.go:detectReferenceIndex` checks only the indexed
`path`. Its `indexedReference` representation discards `source`, so it cannot
prove original-path absence. `internal/speccheck/backlog.go:detectBacklogPromotion`
only rejects originals declaring `status: promoted`; an unchanged `open`
original escapes even when the index proves adoption. Current Finding lifecycle
validation likewise accepts the duplicated `done` source.

The canonical module already requires Backlog moves and terminal routing, but
the explicit shared adoption postcondition for Findings and Backlogs is stated
in the authorial skill/ADR rather than a dedicated module clause. Complete the
portable contract in `internal/baseline/assets/modules/context-workflow.json`
and its generated guidance, then validate source absence from the declared
index regardless of source status. Keep background citations and immutable
historical ownership distinct; never delete an original on basename equality
alone. Test retained originals, missing destinations, legitimate one-owner moves
and secondary consumers against the same contract. This proposed repair is not
implemented and does not widen an operative protected-file grant.

### Sources and influence

Secondbrain index and qmd located the mirrored ADR-0083 and the historical
`2026-07-25-spec-owned-reference-lifecycle.md` Finding. Both were read and checked
against the local ADR/skill; they establish the intended single-owner move.
The new evidence is the public detector's acceptance of duplicates, not a new
claim that source adoption was never designed. No matching measured detector
case was found in the destination Inbox.

Exa MCP read [Git's git-mv documentation](https://git-scm.com/docs/git-mv): a
successful move updates the index and still needs a commit. This supports the
move mechanics only; the local ADR owns placement policy and the controlled
CLI observations establish the validation gap. No proprietary source, private
record or credential was sent to Exa.
