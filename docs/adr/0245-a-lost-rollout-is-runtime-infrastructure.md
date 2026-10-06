---
status: accepted
created_at: 2026-10-06T00:00:00Z
updated_at: 2026-10-06T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A lost rollout is runtime infrastructure

Between 2026-10-05 and 2026-10-06 a repository's Delivery Queue lost three
Runs to `internal error -32603: no rollout found for thread id <uuid>`. Codex
had created each thread and never written its rollout file, so when the next
`acpx prompt` reconnected with `session/resume`, `codex-acp` 2.1.1 could not
find the thread and answered with JSON-RPC `-32603`, message `Internal error`,
and the text in `data.details`. acpx 0.19.4 replaces a session that fails to
resume with a fresh one only while the session holds no Agent message, so a
session that has done work fails the prompt. Roundfix read the failure as the
Agent's: the Task settled `failed`, a QA gate died before its report, and each
attempt spent one of the queue's retries. Nothing in the failure was the
Agent's or the Spec's.

Roundfix now treats that failure as runtime infrastructure. The ACPX Runner
marks a JSON-RPC error as a Lost Rollout only when its code is `-32603` or
`-32600` and its `message` or its string `data.details` contains `no rollout
found for thread id`. Every other `-32603` keeps today's meaning. The runner
reads `data.details` for this match alone and keeps it, bounded like the
message, as the failure's detail.

A Lost Rollout gets no Verification Feedback repair. The Daemon ends the lost
Agent Session and continues the Task in a new one. When the loss happens before
the Task's First Handoff, or in a QA Task while its seeded report still says
`verdict: pending`, the next selection in the profile's Fallback Chain takes
the Task. Otherwise, or when the chain has no next selection, the same
selection continues in a fresh Agent Session. The new session gets the Task
prompt, a notice that the working tree holds the lost session's changes, and
the turn that was lost. Each Agent Session owner recovers at most two Lost
Rollouts; a third settles the Task `failed` as runtime infrastructure. The Run
records every recovery, a QA fallback is recorded in the seeded QA Report, and
the Delivery Queue parks an unrecovered Run as `runtime-infrastructure`. A
retry from that park is not counted against the queue's retry limit.

This narrows ADR-0114 for one case. A Fallback Selection still never takes
over work that a turn finished and handed back. It may take a Task whose first
turn was lost, even though that turn may have changed the tree, because the
runtime discarded the turn before it could be handed back. The maintainer
decided on 2026-10-06 that a QA gate may fall back to the Claude subscription
this way, automatically and recorded in its report.

Matching any `-32603` was rejected because acpx and the adapters use that code
for model errors and failed turns too. Retrying the lost session's id was
rejected because Codex cannot find the thread, so resuming it fails again.
Letting acpx replace the session silently was rejected because the Agent would
lose the conversation without a record. Falling back after a First Handoff was
rejected because the fallback would inherit an attempt the Daemon had already
accepted as a handoff.

## Consequences

The failure reported on 2026-10-06 no longer fails a Task or spends a queue
retry, and an operator can see each recovery in the Run Event Stream. The
detector rests on wording that `codex-acp` and Codex do not publish as a
contract. If the phrase changes, the failure falls back to today's Task
failure, which is visible in the Run, rather than becoming a silent retry. The
Fallback Chain can now take a QA gate after its first turn began, while the
seeded report is still pending.
