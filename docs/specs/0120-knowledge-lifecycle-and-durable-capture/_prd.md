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
- Tooling authority: applicable — protected tooling mutation is proposed, not authorized. The proposal is `docs/specs/0120-knowledge-lifecycle-and-durable-capture/_authorization.md`; bounded files: `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/modules/secondbrain.json`, `internal/spec/archive.go`, `internal/spec/archive_test.go`, `internal/speccheck/backlog.go`, `internal/speccheck/backlog_test.go`, `internal/docscontract/publicdocs_test.go`, `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`, `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`. Source: `docs/agents/agent-instructions.md`.

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
The cross-Spec plan records those limitations. Sources are pending adoption
until this Spec's final scope is approved.
