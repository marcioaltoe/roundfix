---
status: accepted
created_at: 2026-10-13T00:00:00Z
updated_at: 2026-10-13T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A turn that did not finish is not accepted

On 2026-10-13 the QA gate of Spec 0265 lost four Runs to provider capacity, one
of them after 45 minutes of work. In each, `codex-acp` streamed `threadStatus`
`systemError`, then the text "Selected model is at capacity. Please try a
different model.", and then answered `session/prompt` with `stopReason`
`end_turn`. Roundfix took `end_turn` as a finished turn, so the report stayed
`pending` and the configured Claude fallback never ran. The operator then ran
the gate on Claude Opus by hand. That turn
started its load driver as a Claude Code background task, said it would wait
for the driver's events, and ended its turn while the driver still ran. The
report stayed `pending` again, and the Run said only that.

Roundfix now treats both endings as unfinished. A turn whose last `threadStatus`
before its `session/prompt` result is `systemError` is a Provider System Error.
The ACPX Runner returns it as a failure that carries the bounded text the
adapter sent after the status, whatever the stop reason. The Daemon recovers it
the way ADR-0245 recovers a Lost Rollout: the next Fallback Selection takes the
Task before its First Handoff, or a QA Task while its seeded report is still
`verdict: pending`, and the QA fallback is recorded in that report. A Provider
System Error gets no retry on the same selection, because a provider at
capacity gives the same answer seconds later. When no fallback may take the
Task, the Task settles `failed` as runtime infrastructure, and the Delivery
Queue parks the Run as `runtime-infrastructure`; a retry from that park is not
counted.

Background tasks are refused in Run sessions. Every acpx command Roundfix runs
for the `claude` runtime carries `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1`,
which Claude Code documents as disabling "the `run_in_background` parameter on
Bash and subagent tools, auto-backgrounding, and the Ctrl+B shortcut". The QA
prompt also says to run every command in the foreground. Monitor is not named
by that documentation, and an adapter may ignore the variable. So the runner
also records each Agent Background Task a turn starts, from the task id that
`claude-agent-acp` puts in a tool call's `_meta`. A QA turn that ends with its
report pending after starting one is an Unfinished Turn. The QA Task settles
`failed` with a reason that names each task by id, tool and description, and
the Run records the same names.

Matching the capacity text was rejected. The same status carried a 502 from a
gateway and a dropped stream in the recorded sessions, and all three are
provider failures. Retrying the same selection was rejected for the reason above.
Holding the turn open until the job ends was rejected, because
`claude-agent-acp` holds a turn only for background subagents, never for
shells, and Roundfix cannot see whether a task is still alive. Prompt wording
alone was rejected: an Agent does not always follow a prompt rule, and the
refusal and the record still hold when it does not. A fallback after a First Handoff stays rejected, as ADR-0245
decided.

## Consequences

A capacity reply no longer settles a QA gate as `pending` or spends a counted
queue retry, and a QA gate may now move from Codex to the Claude subscription
mid-turn, recorded in its report. The detector rests on `codex-acp`'s
`threadStatus` metadata, which is not a published contract. If the field goes
away, the turn is accepted as it is today, so the failure stays visible in the
Run. A QA command that needs more than the Bash tool's longest timeout must now
be split into foreground steps. An Unfinished Turn names a task the turn started
and the turn did not show finished, not one that Roundfix proved is still
running.
