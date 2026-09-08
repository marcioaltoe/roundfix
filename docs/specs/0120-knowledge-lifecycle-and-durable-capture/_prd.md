---
spec: 0120-knowledge-lifecycle-and-durable-capture
status: active
created: 2026-09-08
surfaces: [backend, cli, docs]
---

# Knowledge has one capture door and a truthful terminal state

Fleet sessions use Secondbrain's Inbox while the canonical layout still
declares a local Inbox. Completed evidence and intent outside a new Spec have
no truthful terminal route, history guidance is incomplete, and capture-time
commit instructions conflict with an installed autosync owner. The outcome
is a consistent lifecycle that preserves provenance and distinguishes local
capture from confirmed remote durability.

This Spec is **in authoring**. Its sources were triaged on 2026-09-08; its
protected mutations and cross-repository contract changes await approval.
There is no executable Task Graph yet.

## Project Constraints

- Identifier strategy: applicable — retain existing Spec slugs, dated artifact basenames, Inbox destinations, and source provenance; a terminal disposition does not invent an implementation identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, authentication, or HTTP API change is proposed. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — preserve the one-way mirror and fleet Inbox boundary, the existing History Root, and accepted regeneration ownership under ADR-0149. Source: `docs/agents/docs-layout.md`.
  ADR-0083 applies to the source-adoption validation follow-up: adopted Findings and Backlog Entries move to one owning Spec, with no copy or stub left at the original path.
  ADR-0116 applies to cited-policy evidence: preserve checks that compare a claim with the cited record; source-origin absence is an additional structural check, not a substitute for semantic citation validation.
  ADR-0123 remains operative: retirement currently uses conservative local reachability and the Review Artifact resolver never writes into history. The proposed stable-evidence change does not yet supersede it.
  ADR-0152 is a proposed revision only, recorded for review; it creates no current obligation or grant.
- Tooling authority: applicable — express maintainer authorization on 2026-09-08 covers the terminal lifecycle instruction through [the narrow grant](references/2026-09-08-terminal-lifecycle-authorization.md); bounded files: `internal/baseline/assets/modules/context-workflow.json`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, with sanctioned digest regeneration. The broader [_authorization.md](_authorization.md) remains proposed, including all Go, test, companion-repository and unrelated skill changes. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- Capture guidance names the actual fleet door and the owner responsible for
  durability, without duplicate sync jobs or false claims of remote storage.
- A resolved or superseded observation can leave the active queue with
  durable evidence even when no new implementation Spec is warranted.
- Every retired documentation family has one documented home and valid
  provenance; active Rollups do not disappear while licensing their members.
- Upstream instructions cannot depend on downstream case evidence, including
  archived cases, without the consistency checker detecting the violation.

## Core Features

1. Canonical layout, Secondbrain instructions, and the destination Inbox
   contract agree on the fleet capture door. Existing historical references
   remain discoverable; no manual mirror synchronization is introduced.
2. Capture identifies whether the local scheduler or the capturing session
   owns publication. Local capture, local commit, and remote confirmation are
   distinguishable evidence. A configured scheduler is not a successful push.
3. Findings and Backlog Entries can record evidence-backed closure or
   supersession without inventing a new Spec or a one-member Rollup. An
   unresolved observation remains active or explicitly deferred with reason.
4. History guidance covers retired Specs, findings, intent, handoffs, decision
   records, and review artifacts using the existing resolver. A moved artifact
   preserves its original observation and obtains a dated disposition addendum.
5. Archive-license resolution preserves existing chains. Closing a Rollup
   requires a valid durable replacement for every member license; it cannot
   silently weaken the current archive check.
6. A deterministic detector rejects upstream citations to downstream evidence
   in both active and history trees. Literal directory declarations defining
   the layout are distinguished from evidentiary citations, with positive
   and negative cases.
7. Relevant research capture keeps URLs, dates, summary, project relevance,
   limitations, and influence on the decision. It feeds later ingestion and
   never claims the pending digest was already integrated into the wiki.

8. Review Artifact retirement must use stable recorded or provider evidence, with explicit unknown outcomes when evidence is unavailable, rather than local object-store availability. Preserve valid squash receipts. Relocation from legacy docs/specs/_reviews is independent of the active/retired decision. This proposed change must explicitly revise ADR-0123 before its conflicting behavior changes.

## Non-Goals / Out of Scope

- Editing read-only mirrors or immutable research sources.
- Replacing the current Secondbrain sync machinery or creating another cron.
- Declaring every historical issue fixed, deleting evidence, or changing an
  archived Spec's original requirements.
- The authorization schema itself, owned by Spec 0119.

## Success Metrics

- The same capture/disposition example receives one unambiguous destination
  from canonical guidance and the destination contract.
- Tests accept a valid historical provenance chain and reject a broken one;
  no active Rollup member loses its license during retirement.
- Upstream-citation negative controls fail for active and archived cases,
  while layout declarations remain valid.

## Decisions

- Confirmed: archive completed/deprecated documentation, triage the Roundfix
  Inbox, preserve useful research, and put authorization inside Specs.
- Observed: the local Secondbrain autosync job is installed with a 7,200-second
  interval, 33 runs and last exit code zero; its script commits, rebases and
  pushes. This does not establish remote durability for this session's files.
- Proposed: honor that existing owner rather than requiring each capture to
  create a competing commit; define an explicit session-owned fallback when
  autosync is absent. The fallback and cross-repository changes need approval.

## Open Questions

- Approve the exact protected mutations and the limited companion contract
  updates in Secondbrain.
- Confirm scheduler-owned capture publication and the no-scheduler fallback.
- Confirm the terminal disposition vocabulary and evidence shape before
  changing lifecycle readers. No new status is used to bypass today's rules.

## Research basis

The 2026-09-08 inventory and triage retain the six relevant fleet observations:
the obsolete local Inbox, incomplete history guidance, capture-time commits,
outside-Spec terminal states, authorization placement, and the missing
upstream-only detector. Current code distinguishes already delivered history
mechanics from remaining contract gaps. Secondbrain's own agent and Inbox
contracts and the inspected autosync script establish ownership locally;
Exa-read execution documentation supports distinguishing authority from
supplied content but does not decide this repository's lifecycle vocabulary.
The [historical research record](https://github.com/marcioaltoe/roundfix/blob/6b8ea48725cbca13974eee0b400b3482202874f6/docs/workflow/2026-09-08-pending-work-plan.md) preserves those limitations at its cited Git revision. Ordinary source ownership is recorded in the reference index. The exact
implementation grant remains pending independently of adoption.

The [source ownership index](references/_index.md) records the pre-adoption
path, type, primary owner and current owned copy. Secondary consumers link
that copy; lifecycle completion means routing, not verified implementation.

## Technical candidate

The [_techspec.md](_techspec.md) records the reviewable implementation map,
coverage and build order. It is a proposed candidate, not a completed authoring
gate or permission to dispatch. Exact governed grants and the named decisions
remain pending; no Task Graph or implementation result is claimed.

## Confirmed terminal lifecycle decision — 2026-09-08

The maintainer requires all terminal Findings, Backlog Entries and Rollups to
leave the active family directories for their `docs/history/` family, without
waiting for consuming Spec implementation, QA or publication. Spec-owned
adopted references remain part of their self-contained Spec. A real standalone
terminal disposition without an absorber records non-empty `closure_reason`
and `closure_evidence`; a fabricated owner is prohibited, and invalid existing
absorption pointers remain errors.

The six terminal Rollups and their 82 member links have now been migrated to
history with real direct Spec absorbers. This resolves their current placement
using existing code. Supporting the standalone closure fields in mechanical
validation remains implementation work. The
[narrow canonical grant](references/2026-09-08-terminal-lifecycle-authorization.md)
covers the requested instruction and generated guide/manifest only; the broader
proposed authorization remains unchanged.
