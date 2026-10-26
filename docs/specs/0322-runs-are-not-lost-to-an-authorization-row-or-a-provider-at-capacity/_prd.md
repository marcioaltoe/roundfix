---
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
status: active
created: 2026-10-26
surfaces: [backend, docs]
---

# Runs are not lost to an authorization row or a provider at capacity

Minimal PRD for a bug fix, minted by `write-techspec`. Two stops parked four
of this repository's Delivery Queue items whose work was sound: two needed a
QA Archive Override and two an operator's retry. A QA row blocked because the Spec's own
authorization forbids its lookup qualifies only when it names the network in
one exact sentence, so the honest rows of Specs 0317 and 0309, which named
what their authorizations actually forbade, parked `qa-environment-partial`.
A provider answering "Selected model is at capacity" ends a QA turn that had
already written a verdict, or a Task's repair turn, with no fallback although
the profile names one, so Specs 0300 and 0302 parked `runtime-infrastructure`.
After this Spec the row qualifies in the form
`blocked (environment: <subject> prohibited by Spec authorization)` with a
bounded subject, and a Provider System Error in any Task or QA turn moves the
work to the next Fallback Selection, reopening a QA verdict the failed turn
wrote; every other qualification rule and the exhausted-chain park keep their
bytes.

## Sources

Two Backlog Entries, written and adopted at birth by this Spec's planning
change (see `references/_index.md`):

1. [`references/2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md`](references/2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md),
   `type: fix`: Spec 0317's row QA-08 and Spec 0309's row R07, both archived
   with the QA Archive Override, and Spec 0258's rows 12–15 as the control.
2. [`references/2026-10-26-a-provider-at-capacity-parks-a-turn-the-fallback-chain-could-finish.md`](references/2026-10-26-a-provider-at-capacity-parks-a-turn-the-fallback-chain-could-finish.md),
   `type: fix`: Runs `run_20261023T022528Z_d6dffdb8e615fb09` (Spec 0300) and
   `run_20261023T070317Z_62af7ae2274ed43c` (Spec 0302), recovery `exhausted`
   with a configured fallback.

Both answer one question, a Run lost to a stop that is not the work's fault,
and each changes one shared policy that every caller reads. No open Backlog
Entry or unresolved Finding shares this context: at the authoring base
`docs/backlog/` and `docs/findings/` hold no entry.

## Prerequisites

None in `requires`. Authored against `d08f8557`, where Specs 0319 and 0320 are
active and deliver first. This Spec reads neither's code. It raises the
`qa-gate` and `roundfix` owned-skill versions that Spec 0319 also raises, and
regenerates `docs/agents/setup-context.json`, which Spec 0320 also
regenerates; both are written by the record and refresh commands at delivery
time, so the delivery base decides the bytes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; a QA Results status form and one line of the QA report's
  `## Agent runtime fallback` section are added. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories,
  fake runtimes and fake runners. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0335: "A stop that is not the
  work's fault does not cost the Run", written by this planning change, which
  amends ADR-0276 (the Network-Prohibited Row) and ADR-0274 (the Provider
  System Error window), both amended in place with a dated line. ADR-0295
  cites ADR-0276 only for its Delivery Base, the half of ADR-0276 this Spec
  leaves unchanged. The shared QA partial policy, Lost Rollout recovery, the
  selection-start fallback and the Light Tier's skip, yield and escalation
  keep their recorded decisions. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AQ answers of 2026-10-26 quoted in the authorization record, which carry the
  cycle AJ grant for the Governed Paths each Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/profiles.md`, `.agents/skills/roundfix/references/settle.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`, `internal/baseline/assets/modules/autonomous-work.json`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/archive-spec/SKILL.md`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A QA partial whose only unmet rows are exempt qualifies when an
  outside-evidence row names what the Spec's authorization forbids in the
  form `blocked (environment: <subject> prohibited by Spec authorization)`,
  in settlement, `roundfix qa-report accept`, `roundfix archive`, Daemon
  settlement and the Delivery Queue alike.
- Nothing else about a qualifying partial loosens: provenance, counts,
  skipped, finding and declared rows, the refusal messages and every report
  that qualifies or fails today keep their outcome.
- A Provider System Error in a Task's first turn, in its repair turn after the
  First Handoff, or anywhere in a QA turn moves the work to the next Fallback
  Selection; only an exhausted chain parks `runtime-infrastructure`.
- A QA fallback session never inherits a verdict the failed turn wrote; the
  report records the verdict it discarded.
- The Baseline clauses, the QA, archive and Roundfix skills, the guides,
  ADR-0335 and `CONTEXT.md` say both rules.

## Core Features

1. The subject form. `internal/spec` reads a Results row as a
   Network-Prohibited Row when its status is
   `blocked (environment: <subject> prohibited by Spec authorization)`, the
   subject is 1 to 80 Unicode code points after trimming and holds no
   parenthesis, and its provenance names the outside-evidence row. A row
   counts toward one exempt kind at most.
2. The fall-through. `internal/daemon` lets the next Fallback Selection take
   a Task or QA turn after any Provider System Error, resends a failed repair
   turn's prompt after the notice, and for QA rewrites a non-pending
   frontmatter verdict to `pending` and adds `- discarded_verdict: <verdict>`
   to the report's `## Agent runtime fallback` section.
3. The Baseline clauses. The sentences of
   `clause.spec.project-constraints-06-outside-evidence` and
   `clause.autonomous.loop-01-qa-once` name the subject form; the derived
   files and this repository's managed guides are regenerated.
4. Documentation. The `qa-gate`, `archive-spec` and `roundfix` skills and
   their references, the settle, archive, qa-report, usage, configuration and
   context-driven guides, and the glossary entries **Network-Prohibited Row**,
   **Provider System Error**, **First Handoff** and **Fallback Chain**.

## User Experience

The QA gate of a Spec whose authorization forbids reading another repository
records `| R07 | blocked (environment: another repository prohibited by Spec
authorization) | outside-evidence row |` and names the authorization text in
its report. Settlement accepts the partial, the Delivery Queue does not park
`qa-environment-partial`, and the archive needs no override. When the QA turn
has written `verdict: partial` and the provider answers "Selected model is at
capacity", the operator reads `qa qa hit a provider system error: Selected
model is at capacity. Please try a different model.; recovery: fallback` on
the progress stream; the Claude fallback opens the report at `verdict:
pending` with `- discarded_verdict: partial` under `## Agent runtime
fallback`, finishes the gate, and the item moves on without a park.

## Non-Goals / Out of Scope

- Reading `_authorization.md` to confirm the subject (rejected by ADR-0276
  and ADR-0335), or accepting any status that mentions the authorization.
- Changing the Network-Denied Row, the pre-PR Pull Request row, Unreachable
  Acceptance coverage, the `### QA settlement` table or the refusal messages.
- Lost Rollout recovery: it already continues in a fresh session when no
  fallback applies and never parks while a selection remains.
- Review Batches, which keep returning a Provider System Error unchanged.
- Light Tier skip, yield and escalation rules, the
  subscription rule for provider routing, and the selection-start fallback.
- Retrying the same selection after a Provider System Error (ADR-0274).
- The `roundfix events` projection that skips `daemon.task` recovery events
  without a `task` field; the selection category already names the cause.

## Success Metrics

1. Success Metric: Spec 0309's and Spec 0317's archived Results tables, with
   their frontmatter counts, qualify; Spec 0258's still qualify.
2. Success Metric: a subject of 80 multibyte code points qualifies, one of 81
   does not, an empty or blank subject, a subject with a parenthesis, a
   changed suffix case and a row without outside-evidence provenance do not,
   and a status in both the subject and network-denied forms counts once.
3. Success Metric: `roundfix settle`, `roundfix qa-report accept`,
   `roundfix archive` and Daemon settlement agree on the two archived shapes
   and on the parenthesis shape.
4. Success Metric: a Task whose repair turn after the First Handoff meets a
   Provider System Error completes on the fallback, whose prompt carries the
   notice and the failed repair prompt; without a fallback it still settles
   `runtime infrastructure: task task_01 hit a provider system error: …`.
5. Success Metric: a QA turn that wrote `partial`, `pass` or `fail` and then
   meets a Provider System Error completes on the fallback, which receives the
   report at `verdict: pending` with `- discarded_verdict: <verdict>`; a
   pending report gains no such line; a single-selection profile still
   settles the QA Task as runtime infrastructure.
6. Success Metric: the Light Tier, Lost Rollout, review and Delivery Queue
   park tests pass unchanged.
7. Success Metric: the Baseline clauses, this repository's
   `docs/agents/spec-routing.md` and `docs/agents/autonomous-work.md`, the
   skills and the guides carry the subject form; a second Managed Refresh
   changes nothing, the owned-skill versions are recorded, and the Behavior
   Surface record is re-recorded last.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Archive Records
  `docs/history/specs/0317-the-baseline-sends-handoffs-and-messages-through-ai-memory.md`
  and `docs/history/specs/0309-a-python-baseline-profile-for-uv-pytest-and-ruff.md`,
  written by the queue and the maintainer before this Spec existed: each
  carries `disposition: qa-override` with a `qa_override_reason` naming the
  row "outside the accepted wording". The Results rows they cite are those of
  Success Metric 1.
- The Test Anything Protocol 14 specification
  (https://testanything.org/tap-version-14-specification.html), read on
  2026-10-26: its grammar keeps the directive fixed and the reason free,
  `Directive := " # " ("todo" | "skip") (" " Reason)?` with
  `Reason := [^\n]+`. The subject form follows it: a fixed prefix and suffix,
  a free but bounded subject.
- The LangSmith LLM Gateway model-fallback documentation
  (https://docs.langchain.com/langsmith/llm-gateway-fallbacks), read on
  2026-10-26: the gateway "Calls each fallback in order until one succeeds,
  returns a status that does not trigger another fallback, or the chain is
  exhausted", with `503` among the trigger codes. The fall-through and its
  exhausted-chain park follow the same contract.

## Unreachable Acceptance

- criterion: after release, a Delivery Queue item whose QA or repair turn
  meets a real "Selected model is at capacity" reply completes on its
  fallback selection without a `runtime-infrastructure` park.
  reason: it needs a real provider at capacity, which no Verification or QA
  row may reach or provoke.
  satisfied-by: the operator's queue log and the Run Event Journal of the next
  Run that meets the reply (event `provider_system_error` with recovery
  `fallback` in a QA scope or after a First Handoff), recorded session
  evidence.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Characterization at `d08f8557`: the Results tables of the three archived
  reports, read with `git show` at `078869bf`, `06a9b7ce` and `ef1a4bab`, give
  Spec 0258 four exact network rows and a Pull Request row (qualifies), Spec
  0309 one row in another subject among four environment rows and Spec 0317
  one among three (both refused by `QAReportEligibility`).
- Run Event Journal, read with `sqlite3 -readonly` and `roundfix events`:
  in both QA Runs the recovery was `exhausted` with a `claude/opus/high`
  fallback configured, and the turn's tool calls had replaced
  `verdict: pending` with `verdict: partial` before the capacity reply
  (Run `…d6dffdb8…` cursors 2487–2496, Run `…62af7ae2…` cursors 4124–4225).
  Task 01 of Run `…62af7ae2…` met the same reply before its First Handoff and
  fell back as designed (cursor 566).
- Root causes, read at `d08f8557`: `internal/spec/qa.go`
  `qaReportResultsCounts` compares the status with `QANetworkProhibitedStatus`
  by `==`; `internal/daemon/agent_session_owner.go`
  `recoverProviderSystemError` requires `!owner.handedOff` for a Task and
  `reportPending()` for QA.
- A prototype of Tasks 01, 02 and 03 in scratch clones passed the full
  `./internal/spec`, `./internal/daemon` and `./internal/baseline` packages,
  the CLI partial-shape, provider and Baseline tests and the documentation
  contracts, and every sabotage in `_techspec.md` → Testing Approach failed its
  named case.
- Secondbrain, read only: `qmd query "roundfix QA outside-evidence row
  prohibited by Spec authorization qualifying partial" --all --files
  --min-score 0.3` returned the Roundfix mirror's `docs/agents/spec-routing.md`,
  the archived Specs 0235 and 0300 and the `qa-gate` and `archive-spec`
  skills; `qmd query "model at capacity provider system error fallback chain
  agent selection" --all --files --min-score 0.3` returned the Fallback
  Chain, Light Tier and Lost Rollout ADRs and
  `docs/references/model-selection.md`. None
  records a decision to keep the exact network sentence or the First Handoff
  boundary for a Provider System Error; they confirmed the ADR set above.
- Exa: the TAP 14 grammar and the LangSmith fallback contract quoted under
  Acceptance evidence. Other gateway documentation read the same day
  (llmgateway.io routing) ends a stream on an error after output reached the
  client instead of falling back; Roundfix differs because the Daemon owns the
  tree and the report, and ADR-0335 records the reopened verdict that makes a
  mid-turn fallback safe.

## Glossary

- changes: **Network-Prohibited Row**
- changes: **Provider System Error**
- changes: **First Handoff**
- changes: **Fallback Chain**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AQ batch 2 brief for Spec 0322 of 2026-10-26, the maintainer's cycle AQ answers quoted in `_authorization.md`, and the characterization above.

Taken from those sources, with no further interview:

- Accept the subject form in the parser and teach it in the skills, instead
  of only tightening the skill wording, because a row that names the network
  when the authorization forbade another repository would be untrue.
- Bound the subject at 80 code points without parentheses, and count each row
  toward one exempt kind.
- Reword the two Baseline clauses so the portable rule matches the parser.
- Remove the First Handoff and pending-verdict window for a Provider System
  Error in Task and QA scopes, and reopen a QA verdict the failed turn wrote.
- One new ADR, ADR-0335, with dated amendments to ADR-0274 and ADR-0276.

## Open Questions

None.
