---
type: feat
status: promoted
created: 2026-11-15
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
reason: null
---

# The status mode shows the Spec queue as a table with forecast start and end times

## Opportunity

The maintainer asked in session on 2026-11-15 (08:52 BRT), after a `/roundfix --status` report: "É possível que no comando /roundfix --status seja informados além do que já temos uma tabela com a fila de specs programada, da ordem de implementacao mais antiga para a mais nova (incluindo possíveis mudanças, ou seja, a fila corrente), com o status e breve observação e para as que ainda não iniciaram, previsão de inicio e fim?"

The `--status` Skill Mode and `roundfix follow-up --status` would add one table of the current Spec queue, oldest scheduled first: each Spec with its stage, a short note, and, for a Spec that has not started, a forecast start and end.

## Value

The maintainer reads the whole delivery plan and its timing in one look instead of piecing it together from `deliver status`, the handoff and the session's memory. The hypothesis is that a forecast stated with its basis lets him plan his own time (authoring windows, Claude limit resets, releases) without asking the Supervisor.

## Shape

Non-binding:

- Rows come from what the command already reads: Delivery Queue items (position, stage, Tasks completed of total, blocker) and active Specs under `docs/specs/` not yet queued, in delivery order.
- Planned Specs that exist only in a handoff or an ai-memory message (for example 0397 on 2026-11-15) are rows only when the mode hands them in; `--status` reads them without accepting the handoff or consuming the message.
- The forecast chains from the running item: the duration comes from past Runs in this repository (for example the median per Spec or per Task, stated as the basis), and a wait the queue cannot know (the Quaestor holding the single verification slot, a release merge, an Idea stage waiting on the maintainer, a usage limit reset) is shown as a stated condition, such as "after 15/11 20:00 BRT", instead of a guessed time.
- Times print in the maintainer's zone; JSON keeps a machine-readable form with the basis of each forecast.

Example the maintainer was shown:

| # | Spec | Status | Note | Forecast start | Forecast end |
|---|---|---|---|---|---|
| 1 | 0396 | running (0/2 Tasks) | code Task with the Agent | 15/11 08:48 BRT (actual) | ~10:30 BRT |
| 2 | 0397 | planned | Idea after the limit reset | after 15/11 20:00 BRT | — (no plan) |

Routing decided by the maintainer on 2026-11-15: record as this Backlog Entry now; Idea-stage judgment and authoring after Spec 0397, one Fable authoring at a time.
