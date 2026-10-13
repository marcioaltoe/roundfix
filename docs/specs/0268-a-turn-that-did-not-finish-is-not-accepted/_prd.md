---
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
status: active
created: 2026-10-13
surfaces: [backend, cli, docs]
---

# A turn that did not finish is not accepted

This is an engineering-framed minimal PRD for a bug fix. It enters through the
TechSpec (`docs/agents/spec-routing.md`), and its design lives in
[_techspec.md](_techspec.md) and ADR-0274.

On 2026-10-13 the QA gate of Spec 0265 lost four Runs to provider capacity.
`codex-acp` 2.2.2 streamed `threadStatus` `systemError`, then "Selected model
is at capacity. Please try a different model.", then answered `session/prompt`
with `stopReason` `end_turn` and `usage` `null`. Roundfix accepted the turn,
the seeded QA report stayed `verdict: pending`, the Run ended Unresolved and the
configured Claude fallback never ran
([adopted entry](references/2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md)).
The operator then ran the gate on `claude / opus / high`. That turn started its
load driver as a Claude Code background task and ended while the driver still
ran, and the report stayed `pending` again
([adopted entry](references/2026-10-13-a-claude-qa-turn-ends-while-its-work-runs-in-the-background.md)).

## Characterization

Measured on 2026-10-13 against `9963ffd3`, read-only:

- The recorded capacity stream
  (`~/.acpx/sessions/081b6e5c-bb31-4259-9941-598b885246eb.stream.ndjson`, 16
  lines) has `threadStatus` `active` on line 12, `systemError` on line 13, the
  capacity text as an `agent_message_chunk` on line 14, and
  `{"stopReason":"end_turn","usage":null}` on line 15. A normal turn ends with
  `threadStatus` `idle` before its result.
- Across the 1,380 recorded streams on this machine, 19 carry `systemError`,
  and every one is followed by a provider failure message: 9 capacity, 9
  `unexpected status 502 Bad Gateway` from a local gateway, and 1
  `stream disconnected before completion`. Each such turn ends `end_turn`; a
  later prompt in the same session returns to `active`.
- `internal/agent/acpx_runner.go` reads only the `stopReason` of the result:
  `RunPrompt` returns success for any non-empty stop reason, and
  `handleStdoutLine` drops `session_info_update` metadata. No file under
  `internal/` names `systemError` or `threadStatus`.
- The Claude QA stream of the same day shows a Bash tool call with
  `"run_in_background": true` whose update carries
  `_meta.claudeCode.toolResponse.backgroundTaskId` `bvlc7cgwx`, then a Monitor
  tool call whose response carries `taskId` `bekzicxjj`, then the message "I'll
  wait for the driver's events before touching the worktree." and `end_turn`.
  Seven recorded Claude streams carry `backgroundTaskId`.
- acpx 0.19.4's `--format json` output writes every ACP message as received
  (`JsonOutputFormatter.onAcpMessage`, which sanitizes only file reads), so
  the `session_info_update` lines and the tool call `_meta` reach Roundfix's
  stdout reader unchanged.
- `claude-agent-acp` 0.89.1 passes `process.env` to the Claude Agent SDK, and
  acpx 0.19.4 spawns the adapter with its own `process.env`. The adapter defers
  a turn's settlement for live background subagents only; its source says
  "shells-never-defer".
- No `internal/` file names `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS` or
  `backgroundTaskId`, and the QA prompt contract in
  `internal/agent/spec_prompt.go` says nothing about background work.

## Prerequisites

None for the code. Spec 0267, authored in the same cycle, also adds a sentence
to the QA prompt contract in `internal/agent/spec_prompt.go` and may raise the
Roundfix Skill's version. Whichever Spec merges second rebases its sentence and
its recorded version on the first (ADR-0233 regenerates a raised version at
merge); no Verification of this Spec reads the other Spec's work.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes. The
  `daemon.task` Run Event gains the phases `provider_system_error` and
  `turn_unfinished` with snake-case keys, and Agent Selection attempts gain the
  reason code `provider_system_error`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the runtime is still the locally
  installed acpx and adapter. One environment variable with a fixed value is
  added to the `claude` runtime's commands; no credential, request or forge
  read is added, and every test uses fake runners, a fake acpx and a temporary
  home. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — see the TechSpec's Project Constraints
  for the full list. ADR-0274 (this Spec) decides that a Provider System Error
  is runtime infrastructure and that background tasks are refused in Run
  sessions, with a QA turn that leaves one unfinished settled as an Unfinished
  Turn. ADR-0245: "When the loss happens before the Task's First Handoff, or in
  a QA Task while its seeded report still says `verdict: pending`, the next
  selection in the profile's Fallback Chain takes the Task", the rule this Spec
  reuses. ADR-0114: "A Fallback Selection may switch ACP Runtime automatically
  only while Agent work has not begun", narrowed by ADR-0274 exactly as
  ADR-0245 narrowed it. ADR-0050: "once Agent work begins, Roundfix fails the
  Work Item instead of switching models over potentially modified state", kept
  for every other failure. ADR-0020: "Without a parsed result, a nonzero exit
  stays a Batch failure exactly as before"; a parsed result is still read
  first, and a Provider System Error is decided from the stream beside it.
  ADR-0057: "Genuine Agent execution, infrastructure, or unreadable-artifact
  failures may still be settled failed by the Daemon without a Verification
  command", which governs an unrecovered Provider System Error and an
  Unfinished Turn. ADR-0189 and ADR-0233 govern the Roundfix Skill's version
  raise. ADR-0256's skills-Task rule applies because
  `internal/cli/deliver_workflow.go` is a covered source. ADR-0259 and ADR-0271
  bind the new tests. ADR-0244: "A Spec now carries a Glossary Declaration",
  answered by the Glossary section below. ADR-0184: "A TechSpec now declares
  numbered Surface Transcripts", answered in the TechSpec with the reason none
  applies. ADR-0242 cites ADR-0050 and decides how a blocked review names its
  protocol step; a review turn with a Provider System Error keeps that
  classification, and the review path's code does not change. ADR-0238 cites
  ADR-0050 and decides the light tier's escalation, which this Spec leaves
  alone. ADR-0069 cites ADR-0050 but decides Baseline analysis; ADR-0160 and
  ADR-0170 cite ADR-0057 but decide the red repository gate and Task
  Carry-Forward of a completed Task; ADR-0166 cites ADR-0057 and records
  undeclared paths, which every Task here avoids by declaring its paths;
  ADR-0250 cites ADR-0189 but decides Baseline module versions and the
  coverage record, neither of which changes. None of those five applies.
  ADR-0176 cites ADR-0166 but decides which text the citation check reads;
  ADR-0180 cites ADR-0069 and decides the built-in selections, and ADR-0181
  and ADR-0217 cite ADR-0180 but decide where a configuration is compared with
  the recommendation and how a Cursor selection names its model; this Spec
  changes none of them, so none applies. ADR-0257 cites ADR-0238 and keeps the
  Daemon the only full gate inside a Run, which these Tasks follow: their
  Agents run focused tests and the Daemon runs the declared Verification.
  ADR-0258 cites ADR-0257 and holds a Spec to the QA rerun classes, which this
  Spec answers with tests named in each QA row, no Surface Transcript to
  reproduce, and byte identity compared with the delivery base.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no lint, formatter, test-runner, Makefile or
  CI configuration changes. The Roundfix Skill's canonical files and its
  `SKILL.md` mirror are Governed Paths, granted by the cycle AG approval of
  2026-10-13. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0268-a-turn-that-did-not-finish-is-not-accepted/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/profiles.md`,
  `skills/roundfix/SKILL.md`.

## Goals

1. A turn whose stream reports a provider system error is never accepted as a
   finished turn, whatever its stop reason.
2. The configured Fallback Chain takes such a turn's Task under the same rule
   that already governs a Lost Rollout, and the QA report records it.
3. A Provider System Error that no fallback may take settles as runtime
   infrastructure, and its Delivery Queue retry is not counted.
4. A `claude` Agent Session in a Run cannot start a Claude Code background
   task, and the QA prompt says every command runs in the foreground.
5. A QA turn that ends with its report pending after starting a background task
   is reported as an Unfinished Turn that names each task.

## User Stories

1. As the operator of a Delivery Queue, I want a capacity reply to move my QA
   gate to the configured fallback, so that one provider's outage does not
   park my Spec with an empty report.
2. As the operator, I want an unrecovered provider failure parked as
   `runtime-infrastructure` with its message, so that I retry without spending
   the queue's retry budget and I know why it stopped.
3. As the operator, I want a QA turn that left work running to be named as
   unfinished, with the task it left, so that I do not read a `pending` report
   as the gate's answer.
4. As an Agent in a QA turn, I want the prompt to tell me that the turn's end is
   the gate's end, so that I run my commands in the foreground.

## Core Features

1. The ACPX Runner marks a Provider System Error: a turn whose last
   `threadStatus` before its `session/prompt` result is `systemError` fails
   with the bounded text the adapter streamed after the status. The sealed
   prompt boundary refuses it the same way.
2. The Agent Session owner recovers a Provider System Error before the First
   Handoff, or in a QA Task while its report is pending, through the next
   Fallback Selection; it records the `provider_system_error` Run Event and the
   QA report's `## Agent runtime fallback` lines. Otherwise the Task settles
   `failed` as runtime infrastructure, and the Delivery Queue parks
   `runtime-infrastructure`.
3. Every acpx command for the `claude` runtime carries
   `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1`; the runner records each Agent
   Background Task a turn starts; the QA prompt contract requires foreground
   commands; and a QA turn that ends with a pending report after starting one
   settles the QA Task `failed` as an Unfinished Turn naming each task.
4. The glossary, the configuration, usage and `deliver` guides, and the
   Roundfix Skill's `profiles` and `deliver` references describe the three
   terms, the recovery, the refusal and the park.

## Non-Goals / Out of Scope

- Retrying a Provider System Error on the same selection. A provider at
  capacity answers the same within seconds; the queue's uncounted retry is the
  later attempt.
- Fallback after a First Handoff, or for review scopes. A review turn with a
  Provider System Error fails as a review runtime failure and is not accepted.
- Proving that a background task is still alive, holding a turn open until it
  ends, or killing it. Closing the Agent Session already ends its processes.
- Background work that a shell command detaches by itself (`&`, `nohup`), and
  background facilities of runtimes other than `claude`.
- ACP v2's `error` stop reason; acpx 0.19.4 speaks ACP v1.
- Raising the Bash tool's timeout. A QA command longer than it is split into
  foreground steps.

## Success Metrics

1. Success Metric: replayed through the runner, the recorded capacity stream's
   turn is a Provider System Error carrying the capacity text, against a
   successful `end_turn` result at `9963ffd3`.
2. Success Metric: in a QA Run whose preferred selection returns a Provider
   System Error before the report, the fallback selection runs the gate and the
   report records the fallback with `- retry_spent: false`.
3. Success Metric: a Provider System Error with no eligible fallback parks the
   Delivery Queue item as `runtime-infrastructure: <scope> hit a provider
   system error: <detail>`, and its retry is not counted.
4. Success Metric: a QA turn that starts an Agent Background Task and ends with
   its report pending settles the QA Task `failed` with a reason naming the
   task's id, tool and description, and every acpx command of a `claude`
   selection carries `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1` while no `codex`
   command does.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write,
read through Exa and the Claude Code documentation during authoring:

- Agent Client Protocol, "v2 Prompt Lifecycle" RFD
  (https://agentclientprotocol.com/rfds/v2/prompt), section "Report failures
  after insertion": "In v1, `session/prompt` lasts until the end of the turn, so
  its error response can report a failure at any point: an expired credential,
  a provider outage, an exhausted quota." and "v1’s `session/prompt` is
  unchanged; it already reports these failures as an error response." A v1
  agent that ends a provider failure with `end_turn` reports it on a channel
  the protocol does not define, which is why Roundfix reads the stream's status
  beside the stop reason.
- Agent Client Protocol, "Prompt Turn"
  (https://agentclientprotocol.com/protocol/v1/prompt-turn), Stop Reasons:
  `end_turn` is "The language model finishes responding without requesting
  more tools". A turn whose model was never reached did not finish responding.
- Claude Code documentation, Environment variables
  (https://code.claude.com/docs/en/env-vars): `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS`
  "Set to `1` to disable all background task functionality, including the
  `run_in_background` parameter on Bash and subagent tools,
  auto-backgrounding, and the Ctrl+B shortcut". The refusal rests on this
  documented behavior; the Monitor tool is not named, which is why the turn
  record stays.
- Claude Code documentation, Tools reference
  (https://code.claude.com/docs/en/tools-reference): setting the variable
  "disables auto-backgrounding along with the rest of the background task
  functionality, so a command that reaches its timeout stops instead." This is
  why a long QA command is split into foreground steps.

## Research record

- Secondbrain: read `wiki/index.md`, then ran
  `qmd query "agent turn ends with provider capacity error accepted as finished fallback runtime" --all --files --min-score 0.3`.
  It returned the Roundfix mirror of `docs/user-guide/usage.md` (0.76), ADR-0114
  (0.63), the archived Spec 0026 and ADR-0050 (0.60) and ADR-0245 (0.50). They
  confirmed that the fallback boundary and its Lost Rollout exception are the
  rule to extend, not to replace. A second query,
  `qmd query "Claude background task run_in_background turn ends before job finishes" --all --files --min-score 0.3`,
  returned `raw/web/2026-08-08-youtube-MpZzWMdmQCE.md` (0.75), "Your coding
  agent doesn't always follow your rules", which supports refusing background
  tasks mechanically rather than by prompt wording alone.
- Exa: one search found the ACP "Prompt Turn" page and the "v2 Prompt
  Lifecycle" RFD quoted above; the RFD also names Codex's app-server split of a
  failed turn (`turn/completed` with status `failed`), which matches treating
  `systemError` as a failed turn. The Claude Code documentation was read for
  the background-task variable and the Agent SDK's `env` default
  (`process.env`).

- Advisory judge: `roundfix spec judge` at the TechSpec stage raised no
  judgment and suggested grouping the background-task entry with
  `docs/backlog/2026-10-13-qa-evidence-records-environment-values.md`
  (P 0.36). That entry stays apart: the maintainer's cycle AG decision gives it
  to Spec 0267, whose subject is what QA evidence may persist, not when a turn
  ends. At the Tasks stage it suggested the standard tier for every Task.

## Glossary

- adds: **Provider System Error**
- adds: **Agent Background Task**
- adds: **Unfinished Turn**
- changes: **Fallback Chain**
- changes: **Fallback Selection**
- changes: **Agent Work Started**

## Decisions

- A Provider System Error is recovered like a Lost Rollout but without a
  same-selection retry. See ADR-0274.
- Background tasks are refused in Run sessions and recorded when the refusal
  does not hold; a QA turn that leaves one with a pending report is an
  Unfinished Turn. See ADR-0274.

## Open Questions

None.
