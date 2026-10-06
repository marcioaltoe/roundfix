---
spec: 0240-a-lost-rollout-is-infrastructure
status: active
created: 2026-10-06
surfaces: [backend, cli, docs]
---

# A lost rollout is infrastructure

Between 2026-10-05 and 2026-10-06, Runs in Oraculum's Delivery Queue
(Roundfix 0.44 to 0.46, `codex-acp` 2.1.1) failed mid-turn with
`internal error -32603: no rollout found for thread id <uuid>`. Codex had
created each thread and never written its rollout file, so the adapter could
not resume it. Roundfix read the failure as the Task's. The Task settled
`failed`, a QA gate died before its report, and each attempt spent one of the
queue's retries. The same QA, re-run by hand on `claude / opus / high`, reached
a verdict ([the adopted Backlog Entry](references/2026-10-06-a-lost-codex-rollout-fails-the-task.md)).

A measurement during authoring showed the failure on the wire:
`{"code":-32603,"message":"Internal error","data":{"details":"no rollout found for thread id <uuid>"}}`,
answering `session/resume`, or `session/prompt` after Agent output. Today's
runner reports it as `agent/protocol error at session setup: Internal error`
or `... at session/prompt: Internal error`, and the phrase is lost
([investigation record](references/2026-10-06-lost-rollout-investigation.md)).

This is a bug fix: the Run, the QA gate and the Delivery Queue keep their
contracts, and a failure that belongs to the runtime stops being charged to the
Task. This minimal PRD exists for the downstream artifact contract; the design
lives in the [_techspec.md](_techspec.md) and ADR-0245.

## Prerequisites

None. Specs 0239, 0241 and 0242 are authored in the same cycle. If another of
them also raises the Roundfix Skill's version, the operator orders the queue so
that the later Spec raises it from the earlier one's value (ADR-0233
regenerates a raised version at merge).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the Run
  Event payload gains snake-case keys under the existing `daemon.task` kind,
  and the Delivery Queue gains one blocker token. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the runtime is still the locally
  installed acpx and adapter, no request, credential or forge read is added,
  and every test uses fake runners, a fake acpx and a temporary home. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0245 (this Spec) decides that a
  Lost Rollout is runtime infrastructure, recovered in a new Agent Session
  without repair or retry, and taken by the Fallback Chain before the First
  Handoff or before a QA report. ADR-0114: "Roundfix therefore treats opening
  an Agent Session as selection, not work", which ADR-0245 narrows for a lost
  first turn only; ADR-0050: "once Agent work begins, Roundfix fails the Work
  Item instead of switching models", kept for every other failure after Agent
  work. ADR-0057: "Genuine Agent execution, infrastructure, or
  unreadable-artifact failures may still be settled failed by the Daemon
  without a Verification command", which governs the third Lost Rollout.
  ADR-0010: "Stop Requests and infrastructure errors still halt the cycle",
  unchanged: a recovered loss is not an error the cycle sees. ADR-0020:
  "Without a parsed result, a nonzero exit stays a Batch failure exactly as
  before", kept; the runner only adds a flag and a detail. ADR-0242: "keeps
  only the error's `message`", which ADR-0245 widens by one field, the string
  `data.details`, read only for the lost-rollout match. ADR-0238 cites ADR-0050
  and decides the light tier's escalation, whose prompt pattern the recovery
  reuses and whose rule it leaves alone. ADR-0187 splits the Roundfix Skill by
  command and ADR-0189 ties an owned skill's version to its content, so the
  skill edit raises the version, and ADR-0233 regenerates it at merge.
  ADR-0237 decides how a Delivery Retry records a merge made outside the queue,
  and the new park keeps that rule. ADR-0184: "A TechSpec now declares
  numbered Surface Transcripts", answered in the TechSpec with the reason none
  applies. ADR-0240: "One QA partial policy, and rows a Run sandbox cannot
  reach" decides when this Spec's QA partial qualifies. The gate is bound by
  ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167, and
  ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check this Spec's
  consistency by citation and receipt. ADR-0178 authorizes each Task commit by
  its grant, ADR-0182 runs Settlement Checks before it, and ADR-0166 records
  undeclared paths; every Task declares its paths. ADR-0069 cites ADR-0050 but decides Baseline analysis; ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry; ADR-0160 and ADR-0170 cite ADR-0057 but decide the red repository gate and Task Carry-Forward of a completed Task; ADR-0192 cites ADR-0178 but decides how a conflict confined to declared derived paths is resolved; ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park; this Spec changes none of them, so none applies. ADR-0180 cites ADR-0069 but decides the built-in selections, and ADR-0181 and ADR-0217 cite ADR-0180 but decide where a configuration is compared with the recommendation and how a Cursor selection names its model; ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot; none of them changes here. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's canonical files and its
  `SKILL.md` mirror are Governed Paths. The maintainer authorized skill edits
  ("considere autorizado a ajustar todas as skills se necessário") and, on
  2026-10-06, the Governed Paths this Spec declares ("Concedo"). No other
  Governed Path changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0240-a-lost-rollout-is-infrastructure/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/profiles.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A lost rollout is recognized by its code and its phrase, and no other
  `-32603` changes meaning.
- A Lost Rollout gets no repair and spends no queue retry. The Task continues
  in a new Agent Session, on the next Fallback Selection when the loss came
  before the First Handoff or before the QA report, or on the same selection
  otherwise.
- Every recovery is recorded in the Run, and a QA fallback also in its report.
- An unrecovered Lost Rollout parks the queue item as `runtime-infrastructure`,
  and the retry from that park is not counted.
- What the record says about why the rollout was not written, and the upstream
  report, are kept with the Spec.

## User Stories

1. As an operator running an unattended Delivery Queue, I want a runtime that
   loses its session to cost me neither a Task nor a retry, so that a Spec does
   not park over an adapter defect.
2. As an operator, I want a QA gate whose Codex session was lost before the
   report to finish on my Claude fallback automatically and say so in its
   report, so that I do not re-run QA by hand.
3. As an operator reading a Run, I want each lost session, its step and the
   recovery chosen to appear in the Run Event Stream, so that I can tell a
   runtime defect from a Task failure.

## Core Features

1. **The runner names a Lost Rollout.** A JSON-RPC error whose code is `-32603`
   or `-32600` and whose `message` or string `data.details` contains
   `no rollout found for thread id` is marked as a Lost Rollout, with its step
   (now including `session/resume`) and the bounded detail (ADR-0245).
2. **The Daemon recovers it.** No Verification Feedback repair. The lost
   Agent Session is ended, and the Task continues in a new one. The next
   Fallback Selection takes the Task before the First Handoff, or for a QA Task
   while its seeded report is still pending. The same selection continues in a
   fresh session otherwise. Each Agent Session owner recovers at most two Lost
   Rollouts.
3. **The Run and the QA report record it.** A `daemon.task` Run Event with
   phase `rollout_lost`, a stderr notice, and, for a QA fallback, a Daemon-owned
   section in the seeded QA Report.
4. **The queue spends no retry.** A Run left unresolved by a third Lost Rollout
   parks as `runtime-infrastructure` in the `environment` Park Class, and its
   Delivery Retry is not counted against the retry limit.
5. **Docs and glossary.** `CONTEXT.md` gains **Lost Rollout** and **First
   Handoff**, and the fallback guides, the `deliver` guide and the Roundfix
   Skill's `profiles` and `deliver` references describe the recovery.

## Non-Goals / Out of Scope

- Fixing why Codex does not write the rollout; that is the adapter's and
  Codex's defect, and the upstream report text is recorded, not filed.
- Treating any other `-32603`, `-32600` or adapter error as infrastructure.
- Fallback after a First Handoff, or for review Runs and the pre-PR review.
- Runs without Agent Selection Profiles (legacy configuration), which keep
  today's Task failure.
- Setting a `CODEX_HOME` per Run or changing acpx arguments or TTLs.

## Success Metrics

1. Success Metric: stdout recorded from acpx 0.19.4 for a loss during
   `session/resume` and during `session/prompt` yields a Lost Rollout naming
   `session/resume` and `session/prompt`, with the phrase in its detail. An
   `Internal error` without the phrase, and the phrase under another code, are
   not Lost Rollouts.
2. Success Metric: with a fake runner, a loss before the First Handoff runs the
   Task's next prompt on fallback 1 in a new Agent Session, and a loss after it
   continues on the same selection in a new Agent Session. Both reach Task
   `completed` with no Verification Feedback prompt and one `rollout_lost`
   event each.
3. Success Metric: a QA Task whose first turn is lost while its report is
   pending runs on fallback 1. Its report carries the Daemon's fallback section
   with `retry_spent: false`.
4. Success Metric: a third Lost Rollout settles the Task `failed` with a reason
   starting `runtime infrastructure: lost rollout`. The Delivery Queue parks
   the item as `runtime-infrastructure`, and `deliver retry` returns it without
   raising its retry count or hitting the limit.
5. Success Metric: every existing runner, Daemon and delivery test passes
   unchanged.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Oraculum report of 2026-10-06, now in the Secondbrain at
  `inbox/roundfix/_triaged/2026-10-06-codex-acp-perde-o-rollout-e-a-task-falha-com-32603.md`:
  three thread ids, the message, the QA cost and the manual Claude re-run.
- The installed `@agentclientprotocol/codex-acp` 2.1.1 bundle
  (`dist/index.js`), read on 2026-10-06: Codex writes a rollout on the first
  user message; `resumeThread` recovers a missing rollout only from the live
  app-server; any error that is not a `RequestError` becomes `internalError`
  with the text in `data.details`.
- The installed acpx 0.19.4 bundle and a measurement with a fake adapter under
  a temporary `HOME`: a lost rollout during `session/resume` or
  `session/prompt` of a session with Agent messages exits `1` with the error
  line on stdout. Without Agent messages, acpx replaces the session silently.
- A census of this machine's `~/.acpx/sessions`: 934 of 993 stream files
  between 2026-09-22 and 2026-10-05 carry the error at `session/resume`, all
  replaced silently by acpx.
- The Codex tracker: openai/codex#16872 (turn completes, rollout never
  materializes, `failed to queue rollout items: channel closed`),
  openai/codex#42099 (zero-turn threads not persisted since 0.151.0) and
  openai/codex#28496 (clients must match `no rollout found` by text).

## Glossary

- adds: **Lost Rollout**
- adds: **First Handoff**
- changes: **Fallback Chain**
- changes: **Fallback Selection**
- changes: **Agent Work Started**

## Decisions

- Match the code and the phrase, reading `data.details` only for this match;
  see ADR-0245.
- Recover in a new Agent Session; fall back before the First Handoff or a QA
  report; at most two recoveries per owner; see ADR-0245.
- Park an unrecovered loss as `runtime-infrastructure` with an uncounted retry;
  see ADR-0245.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
