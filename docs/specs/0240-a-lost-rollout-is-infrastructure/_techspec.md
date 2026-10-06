---
spec: 0240-a-lost-rollout-is-infrastructure
prd: _prd.md
created: 2026-10-06
---

# A lost rollout is infrastructure — Technical Spec

## Executive Summary

Codex writes a thread's rollout lazily, and `codex-acp` 2.1.1 can resume a
thread without one only from the app-server that created it. When a later
`acpx prompt` process reconnects with `session/resume`, a thread whose rollout
was never written fails with `-32603 Internal error` and the phrase
`no rollout found for thread id` in `data.details`. acpx 0.19.4 hides that
case while the session has no Agent message. Once the session has done work,
it exits `1`. Today the ACPX Runner reports the failure as a generic protocol
error and drops the phrase. The Daemon then settles the Task `failed`, or a QA
gate dies before its report, and the Delivery Queue spends a retry.

The fix has three parts. The runner marks the failure as a Lost Rollout. The
Daemon's Agent Session owner ends the lost session and continues the Task in a
new one, with no repair. The Fallback Chain takes the Task before its First
Handoff, or in a QA Task while the seeded report is pending. The Delivery
Queue parks an unrecovered loss as `runtime-infrastructure` and does not count
its retry. The trade-off accepted: the detector rests on wording the adapter
does not publish as a contract, and a fallback may inherit a tree a lost first
turn changed (ADR-0245).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the Run
  Event payload gains snake-case keys under `daemon.task`, and the queue gains
  the blocker token `runtime-infrastructure`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the runtime remains the local
  acpx and adapter, no request or credential is added, and tests use fake
  runners and a fake acpx. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0245 (this Spec) governs the
  matcher, the recovery, the fallback boundary and the uncounted retry.
  ADR-0114: "Roundfix therefore treats opening an Agent Session as selection,
  not work", narrowed by ADR-0245 only for a lost turn before the First
  Handoff or a QA report. ADR-0050: "once Agent work begins, Roundfix fails the
  Work Item instead of switching models", kept for every other failure.
  ADR-0057: "Genuine Agent execution, infrastructure, or unreadable-artifact
  failures may still be settled failed by the Daemon without a Verification
  command", which settles the third Lost Rollout. ADR-0010: "Stop Requests and
  infrastructure errors still halt the cycle", unchanged, because a recovered
  loss returns no error and an exhausted one settles the Task. ADR-0020:
  "Without a parsed result, a nonzero exit stays a Batch failure exactly as
  before", kept. ADR-0242: "keeps only the error's `message`", widened by
  ADR-0245 to the string `data.details` for the lost-rollout match only.
  ADR-0238's light-tier escalation keeps its rule; the recovery reuses its
  prompt pattern. ADR-0187, ADR-0189 and ADR-0233 govern the skill edit and its
  version. ADR-0237's merge-outside-the-queue retry is unchanged. ADR-0184: "A
  TechSpec now declares numbered Surface Transcripts", answered below with the
  reason none applies. ADR-0240 decides the QA partial policy. The gate is
  bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155, ADR-0156 and
  ADR-0167; ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183 check
  consistency; ADR-0166, ADR-0178 and ADR-0182 bind each Task commit. ADR-0069 cites ADR-0050 but decides Baseline analysis; ADR-0096 and ADR-0097 cite ADR-0080 but decide the gate's machine stage and row carry; ADR-0160 and ADR-0170 cite ADR-0057 but decide the red repository gate and Task Carry-Forward of a completed Task; ADR-0192 cites ADR-0178 but decides how a conflict confined to declared derived paths is resolved; ADR-0229 cites ADR-0167 but decides how an operator archive resumes a park; this Spec changes none of them, so none applies. ADR-0180 cites ADR-0069 but decides the built-in selections, and ADR-0181 and ADR-0217 cite ADR-0180 but decide where a configuration is compared with the recommendation and how a Cursor selection names its model; ADR-0194, ADR-0195 and ADR-0210 cite ADR-0097 but decide what a QA row records, when it is observed again and its evidence snapshot; none of them changes here. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — task_01 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths. Express maintainer
  authorization: "considere autorizado a ajustar todas as skills se
  necessário", and the maintainer's answer "Concedo" of 2026-10-06 for the
  Governed Paths this Spec declares. Bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/profiles.md`,
  `skills/roundfix/SKILL.md`. No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0240-a-lost-rollout-is-infrastructure/_authorization.md`.

## System Architecture

No new package, command or flag. Four existing components change:

- **ACPX Runner** (`internal/agent`): the prompt-stream trace tracks
  `session/resume` and reads `data.details` of an error line. A new file holds
  the matcher and `DescribeLostRollout`.
- **Agent Session owner** (`internal/daemon/agent_session_owner.go`):
  `Run` recognizes a Lost Rollout before its selection-failure branch, ends
  the session, and either activates the next Fallback Selection or re-activates
  the same one under a fresh session name. Then it repeats the lost turn.
- **QA step** (`internal/daemon/task_engine.go`): gives the QA owner a
  report-pending predicate and appends the Daemon-owned fallback section to the
  seeded report.
- **Delivery Queue** (`internal/delivery`, `internal/store/delivery.go`,
  `internal/cli/deliver_workflow.go`): reads the Run's exhausted `rollout_lost`
  event, parks `runtime-infrastructure`, and retries it without counting.

## Implementation Design

### Interfaces

```go
// internal/agent/lost_rollout.go
const LostRolloutPhrase = "no rollout found for thread id"

type LostRollout struct {
	Runtime    string // from SelectionFailureError.Runtime when present
	Step       string // the ProtocolFailure step
	Detail     string // bounded text that carried the phrase
	PromptSent bool
}

func DescribeLostRollout(err error) (LostRollout, bool)

// ProtocolFailure gains: LostRollout bool; Detail string
```

```go
// internal/daemon: the owner gains handedOff bool, rolloutRecoveries int and
// reportPending func() bool (nil outside a QA scope).
const maxLostRolloutRecoveries = 2

// internal/delivery
const BlockerRuntimeInfrastructure = "runtime-infrastructure"
// RunResult gains: RuntimeInfrastructure string // "" or the park detail
```

### Invariants

1. A JSON-RPC error line is a Lost Rollout when its `code` is `-32603` or
   `-32600` and its `message`, or its `data.details` when that is a string,
   contains `LostRolloutPhrase`. No other field of `data` is read, and no other
   error line changes meaning.
2. `session/resume` is a tracked protocol request, so a failure answering it
   names `session/resume`, not `session setup`.
3. `ProtocolFailure.Detail` is set only for a Lost Rollout. It holds the text
   that carried the phrase, bounded like the message (one line, at most
   `ProtocolFailureMessageLimit` bytes on a rune boundary). The error types,
   `Reason`, `Err`, the exit mapping and the selection/batch classification are
   unchanged.
4. `DescribeLostRollout` finds a `ProtocolFailure` with `LostRollout` in the
   chain of a `SelectionFailureError` or `BatchFailureError` and reports false
   for every other error, including context cancellation and `StopError`.
5. In `agentSessionOwner.Run`, a Lost Rollout is handled before the
   selection-failure branch, unless the context is done. The owner ends the
   lost Agent Session and counts one recovery. It never returns the lost turn's
   error to the Task engine while a recovery remains, so no Verification
   Feedback repair starts for it.
6. The owner sets `handedOff` when a prompt returns without error. A recovery
   activates the next Fallback Selection when one exists and either
   `handedOff` is false (a Task scope) or `reportPending()` reports true (the QA
   scope). Otherwise it re-activates the same selection.
7. A re-activated or fallback session has a name no earlier session of the
   owner used: the candidate's name plus `-rollout-NN`, where `NN` is the
   recovery count. The lost session id is never resumed.
8. The prompt sent to the new session is the owner's first Task prompt, then
   the notice `The previous Agent Session was lost by the ACP Runtime; the
   working tree holds its changes. Keep them and continue.`, then the lost
   prompt when it differs from the first Task prompt.
9. Each recovery publishes a `daemon.task` Run Event with phase
   `rollout_lost` and payload `scope_kind`, `scope_id`, `selection`, `step`,
   `detail`, `recovery` (`fallback`, `new_session` or `exhausted`),
   `next_selection` when `fallback`, and `retry_spent: false`. It writes one
   stderr line starting `roundfix:` that names the scope, the step and the
   recovery. A `fallback` recovery also publishes the existing
   `agent_selection_fallback` notification with `reason_code` `rollout_lost`.
10. A third Lost Rollout in one owner publishes `recovery: exhausted` and
    returns an error whose text starts `runtime infrastructure: lost rollout`.
    For a Task, `agentFailureReason` keeps that text as the failure reason. For
    the QA step, the QA Task settles `failed` with that reason instead of
    halting the cycle. The Run ends Unresolved.
11. Before a QA fallback session starts, the Daemon appends to the seeded
    report (still `verdict: pending`) a section `## Agent runtime fallback`
    with the lines `- lost_rollout: <selection> at <step>: <detail>`,
    `- fallback: <selection>` and `- retry_spent: false`. It never changes the
    frontmatter or the Results table.
12. The delivery workflow marks an Unresolved Run as runtime infrastructure
    when its journal holds a `daemon.task` event with phase `rollout_lost` and
    recovery `exhausted`. The queue then parks
    `runtime-infrastructure: <scope_id> lost its rollout at <step>` before any
    `qa-environment-partial` or `run-unresolved` decision.
13. `ClassifyPark` puts `runtime-infrastructure` in the `environment` class
    with the next action `run roundfix deliver retry <slug>; a retry from
    runtime-infrastructure is not counted`. A Delivery Retry from that blocker
    skips the retry-limit refusal and leaves `retry_count` unchanged. Its
    re-entry stage and carry-forward are those of `run-unresolved`.

### Data Models

The Run Event payload of Invariant 9 is new under an existing kind. The queue
item keeps its schema: the blocker is a string, and `retry_count` simply does
not move. No migration.

### API Contracts

1. API Contract: lost rollout from a prompt process. `ACPXRunner.RunPrompt`
   with the stdout recorded for a loss during `session/resume` returns a
   `SelectionFailureError`, and during `session/prompt` a `BatchFailureError`,
   as today. `DescribeLostRollout` reports `session/resume` or
   `session/prompt`, `PromptSent` false or true, and a detail containing the
   phrase. The error text gains `; lost rollout: <detail>` after the JSON-RPC
   code.
2. API Contract: recovery. Per Invariants 5 to 10, the owner's `Run` returns
   the result of the continued turn, or the exhausted error after the third
   loss.
3. API Contract: QA record. Per Invariant 11.
4. API Contract: queue park and retry. Per Invariants 12 and 13; `deliver
   status` prints `Park: <slug> environment: run roundfix deliver retry <slug>;
   a retry from runtime-infrastructure is not counted`.

### Surface Transcripts

None. The changed surfaces are the Run Event Stream, a stderr notice and a
Delivery Queue park reached only when a live adapter loses a session. They are
proved at the runner seam with stdout recorded from acpx 0.19.4, at the owner
and QA seams with fake runners, and at the queue seam with a fake executor and
journal. QA re-measures acpx with a fake adapter.

## Coverage Map

- Goal 1 → Invariants 1 to 4; API Contract 1.
- Goal 2 → Invariants 5 to 8; API Contract 2.
- Goal 3 → Invariants 9 and 11; API Contracts 2 and 3.
- Goal 4 → Invariants 10, 12 and 13; API Contract 4.
- Goal 5 → Integration Points; Build Order 1.
- User Story 1 → API Contracts 2 and 4.
- User Story 2 → Invariants 6 and 11; API Contract 3.
- User Story 3 → Invariant 9.
- Core Feature 1 → API Contract 1.
- Core Feature 2 → API Contract 2.
- Core Feature 3 → Invariants 9 and 11.
- Core Feature 4 → API Contract 4.
- Core Feature 5 → Build Order 1; Vocabulary Contract.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 2.
- Success Metric 4 → Testing Approach 2 and 3.
- Success Metric 5 → Testing Approach 4.

## Integration Points

- `@agentclientprotocol/codex-acp` 2.1.1 maps a missing thread to
  `{"code":-32603,"message":"Internal error","data":{"details":"no rollout found for thread id <uuid>"}}`.
  It also recovers an unwritten rollout only from the app-server that created
  the thread.
- acpx 0.19.4 reconnects every `prompt` with `initialize` and `session/resume`
  (or `session/load`) when the queue owner is gone. After a failed resume it
  creates a fresh session itself only when the session record holds no Agent
  message (`shouldFallbackToNewSession`). Otherwise it prints the error line on
  stdout and exits `1`.
- Roundfix sets `CODEX_PATH` and no `CODEX_HOME`, so every Run shares
  `~/.codex`. The [investigation record](references/2026-10-06-lost-rollout-investigation.md)
  holds the findings, the census and the upstream report text, which is not
  filed.

Research record. The Secondbrain was read first (`wiki/index.md`, then a
`qmd` query on the defect). It returned only the triaged Oraculum report and
the Roundfix mirror; no prior decision covers a lost rollout. Exa found
openai/codex#16872, #42099, #28496, #22064 and #32728, which show the same
error from other clients and the lazy, unflushed rollout write. The installed
adapter and acpx bundles, read locally, decided the matcher (`data.details`,
`-32603`) and the boundary (acpx hides the empty case).

## Testing Approach

1. Runner seam: a new `internal/agent/lost_rollout_test.go` drives `RunPrompt`
   through `runFakeACPXPrompt` with the stdout in the investigation record's
   appendix: a loss during `session/resume` and during `session/prompt`. It
   also covers an `Internal error` with other details, the phrase under code
   `-32602`, the phrase in `message` with code `-32600`, and a detail longer
   than the limit.
2. Owner and QA seams: a new `internal/daemon/lost_rollout_recovery_test.go`
   with the package's fake runner. It covers a loss on the first turn with a
   fallback configured, a loss on the Verification Feedback turn after a
   handoff, a loss with no fallback configured, a QA loss while the report is
   pending, and three losses in one owner.
3. Queue seam: new `internal/delivery/runtime_infrastructure_test.go` and
   `internal/cli/deliver_runtime_infrastructure_test.go`. They cover the park
   from a fake executor result, the Park Class and next action, an uncounted
   retry at the retry limit, and the journal read of an exhausted event.
4. Existing runner, Daemon and delivery tests stay unchanged and green.

## Build Order

1. The glossary, the fallback and `deliver` guides and the Roundfix Skill's
   `profiles` and `deliver` references describe the Lost Rollout, the First
   Handoff, the recovery and the park; the skill version is raised and
   recorded.
2. The ACPX Runner names a Lost Rollout (API Contract 1).
3. The Daemon recovers it, records it in the Run and the QA report, and the
   queue parks an unrecovered one without counting its retry (API Contracts 2
   to 4) (depends on: 1, 2).
4. The final QA gate (depends on: 3).

## Risks & Considerations

- If the adapter rewords the phrase, the matcher stops matching and the
  failure returns to today's Task failure, which is visible. QA re-reads the
  installed adapter.
- A fallback after a lost first turn inherits whatever the turn wrote. The
  prompt says so, and Verification still decides the Task.
- A runtime that loses every session costs two extra sessions per Task before
  the park. The bound keeps that finite.
- The census shows acpx silently replacing empty sessions daily. That stays
  acpx's behavior and is out of scope.

## Vocabulary Contract

- emits: `internal/delivery/engine.go`
  pattern: `runtime-infrastructure`
  documented-in: `docs/user-guide/commands/deliver.md`
- emits: `internal/daemon/agent_session_owner.go`
  pattern: `rollout_lost`
  documented-in: `docs/user-guide/usage.md`

Glossary terms adopted in `CONTEXT.md` by task_01: **Lost Rollout** and
**First Handoff**; **Fallback Chain**, **Fallback Selection** and **Agent Work
Started** gain the Lost Rollout exception.

## Glossary

- not a term: **ACPX Runner** — the existing runner component in `internal/agent`, named as code, not a domain term

## Decisions

- Match code and phrase; read `data.details` for this match only; see
  ADR-0245.
- Recover in a fresh Agent Session, never by resuming the lost id; fall back
  only before the First Handoff or a pending QA report; see ADR-0245.
- Bound recoveries at two per owner and park the third as
  `runtime-infrastructure` with an uncounted retry; see ADR-0245.
