---
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
created: 2026-11-16
status: promoted
leverage: compounding
---

# The status view shows the Spec queue with a forecast start and end

## Overview

The status view of the Follow-up Command, and the `--status` Skill Mode that relays it, gain one table of the repository's current Spec queue, oldest scheduled first. Each row names a Spec, its status and a short note, with a forecast start and end when the Spec has not started, and its actual start and a forecast end when it is in progress. Each forecast is two times, a likely one and a limit, computed from how long this repository's own Specs took. A wait that Roundfix cannot date is shown as a named condition, never as a time.

It is for the person who supervises Spec work and plans their own time around it: when to answer a question, when an authoring window opens, when a release can go out. In this repository that is the maintainer. An adopter gets the same table from its own history.

It is valuable because the status view holds no time today, and the obvious estimate is wrong by a wide margin. A Clean Implement Run took a median of 47 minutes over the 30 days this repository's Run Database holds, and a plan took a median of 6.0 hours to go from the default branch to its delivery merge over the 58 Specs created since 2026-11-01 (baseline adopted by Spec 0397, `docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/references/2026-11-16-delivery-speed-and-plan-quality-baseline.md@f4580437cd6f3c354c2524f7f2eca9c924446d64`). Whoever wants the timing asks the Supervisor, who pieces it together from `roundfix deliver status`, `roundfix runs list`, the Session Handoff and its own context.

V1 is one Spec and covers this repository's status view only. It ships only if the forecast holds on history: the PRD replays the method over about 100 Specs of the last 30 days, and the limit must hold for at least 80% of them. The view stays read-only and as fast as it is today.

## Problem

On 2026-11-15, after a `/roundfix --status` report, the maintainer asked: "É possível que no comando /roundfix --status seja informados além do que já temos uma tabela com a fila de specs programada, da ordem de implementacao mais antiga para a mais nova (incluindo possíveis mudanças, ou seja, a fila corrente), com o status e breve observação e para as que ainda não iniciaram, previsão de inicio e fim?" The request is the adopted Backlog Entry [`references/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md`](references/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md).

**The status view holds no time.** It lists each pending Spec with its queue position, its stage and its Tasks completed of total. `roundfix runs list` shows when a Run started and how long it has run, and `roundfix deliver status` shows stages and blockers without a time. Roundfix records when each Run began and ended and nothing else about a delivery's time: a Delivery Queue item holds no time, and a new queue replaces a finished one together with its items.

**The queue the maintainer plans by is longer than the queue the view reads.** On 2026-11-15 his next two pieces of work were Spec 0397, planned only in a Session Handoff, and this feature, an open Backlog Entry with a Waiting Declaration. The status view showed neither. It lists a Spec only once its folder holds a PRD, so a Spec in its Idea stage is absent, and it reads no Backlog Entry.

**A duration guessed from a Run is several times too short.** The same baseline counts 1.58 Implement Runs per Spec, and 40% of the Specs had at least one Run that did not end Clean. The 6.0 hours start when the plan reaches the default branch, so the wait before the first Run, the review, the checks and the retries fill the distance from 47 minutes. The spread is as wide as the gap, and a single time hides it: a quarter of the plans merged within 3.7 hours, a quarter took more than 10.9, and one in ten more than 19.2.

**Some waits cannot be dated at all.** A Run Slot Wait names the Run that holds the slot and its repository, and nothing says when that Run ends. A parked item waits for the answer to its Pending Question. A Spec being authored waits for the maintainer to approve each phase, and a Backlog Entry waits on what its Waiting Declaration says. A time printed for any of them would be invented.

### Market Data

- The Kanban Guide (May 2025) defines the forecast this Idea uses. A service level expectation is "a forecast of how long it should take a work item to flow from started to finished", with "two parts: a period of elapsed time and a probability associated with that period (e.g., '85% of work items will be finished in eight days or less')", and it "should be based on historical cycle time".
- ProKanban says how to state it: "The key is to maintain communicating it as a range and level of confidence so as not to create confusion toward a guarantee." It also bounds the tool: "This is for single items only. If you have a batch of work, the SLE is not the appropriate tool." Chaining down a queue therefore needs its own check.
- ASOS checked 1,200 such forecasts against what happened. The 85th percentile was right in 361 of 400 (90%), the 70th in 336 of 400 (84%) and the 50th in 270 of 400 (68%). The amount of history, from 6 to 12 weeks, made no significant difference.
- On the size of the history: with 11 past items the next one falls inside the range already seen with a probability of 83.3%, and with 20, 90.5% (Herbig). A forecast "cannot produce variation the record does not contain": "tails are estimated from tails". Delays that come in runs leave a chained forecast "too narrow at the upper end, exactly where commitments are read" (MerceMay).
- Jenkins estimates a build from the mean of up to three recent builds. Its tracker records what happened when failed builds entered that mean: "Taking into account failing builds makes the estimate completely meaningless." The same report asks that with too few successful builds Jenkins "just give no estimate at all".
- Mergify's merge queue publishes `estimated_merge_at` for each pull request and returns "`null` when it cannot be estimated", "as behind a frozen batch", beside states such as `waiting_for_checks_slot` and `waiting_schedule`. It names the state where a time would be invented.
- A GitLab request of 2024 asks for a job's estimated end from "a statistic (p50, p95) for completed jobs". In the tools read, an estimate belongs to one job or one pull request already in a queue. None was found that carries a forecast through work that is planned and not started. This is an absence in the search, not proof.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | The Spec queue table | Critical | The status view shows one table of the current queue, oldest scheduled first: Delivery Queue items in queue order, then planned Specs not yet queued, then Specs being authored, then open Backlog Entries. Each row has the Spec, its status, a short note, a forecast start and a forecast end. The view's existing lines and JSON fields stay. |
| F2  | Status and note from recorded state | Critical | The status is the row's stage: the Delivery Queue stage with Tasks completed of total, the Authoring Phase in progress, or "not yet a Spec". The note is one short line taken from what Roundfix already records: a blocker with its Park Class and next command, the holder of a Run Slot Wait, a Pending Question, or a Backlog Entry's Waiting Declaration. |
| F3  | A forecast is a likely time and a limit | Critical | Every forecast shows two times, such as "~10:30, by 12:10": the median and the 85th percentile of how long this repository's delivered Specs took. A Spec in progress shows its actual start and a forecast end. |
| F4  | The end is the delivery merge | Critical | A Spec ends when its delivery is merged into the default branch, after its Run, review, checks and merge. The release stays on the view's next-release line. |
| F5  | A wait is a condition, and the queue chains from it | Critical | When a start depends on something Roundfix cannot date, the row shows the condition as its start, such as "after 0397 is merged", and its end as a duration after the start. The rows behind it chain the same way, and no time is invented. |
| F6  | Accepted by a backtest | Critical | Before the release, the method is replayed over about 100 Specs of this repository's history, each forecast made only from what was recorded before that Spec started. It is accepted only if the actual end falls at or before the limit for at least 80% of them. |
| F7  | Rows beyond the Delivery Queue | High | A Spec in its Idea, PRD, TechSpec or Tasks phase is a row, read from the Phase Record that Spec 0397 delivers and from the Spec folders. An open Backlog Entry is a row marked "not yet a Spec", with its declared wait and no forecast end. |
| F8  | The basis is stated | High | A footer under the table says what the forecasts were computed from: how many delivered Specs, over which period, on this machine, and which percentiles. With too little history a cell says so and shows no time. The JSON carries each row's likely time, limit, condition and basis in a form a program can read. |
| F9  | The Skill Mode relays the table | High | The `--status` Skill Mode shows the table as the command printed it and never computes or adjusts a forecast itself. |
| F10 | Read-only, at today's speed | High | The status view still writes nothing, accepts no handoff and consumes no message, and its time stays inside the band of the last KPI. |

## KPIs

The first three are read from the backtest the PRD runs, over the Specs of the last 30 days in this repository's Run Database, each forecast made from the history before that Spec started. The operator repeats them over the first 20 Specs delivered after the release.

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| The limit holds | The actual end falls at or before the forecast limit for >= 80% of the Specs | The backtest: for each delivered Spec, compare the limit forecast at its start with its delivery merge. |
| The likely time is a median | The actual end falls at or before the likely time for 40% to 60% of the same Specs | The same backtest. Outside the band, the likely time is biased and the method is corrected before release. |
| The limit is narrow enough to plan by | Median of the limit's duration divided by the actual duration <= 2.5 (proposal) | The same backtest. One distribution for every Spec gives about 2.7 on the plan-to-merge baseline, an interpolation between its third quartile of 10.9 hours and its 90th percentile of 19.2 against the median of 6.0. The PRD confirms the target or replaces it with the measured number. |
| Rows of the scheduled queue that the view shows | 100% of Delivery Queue items, Specs being authored and open Backlog Entries (today: queue items and Specs with a PRD) | A fixture that holds one row of each kind, and this repository's own status on the release day. |
| Forecast cells that show a time without a basis, or a time where a wait cannot be dated | 0 | A fixture with one row for each kind of wait and one repository with too little history; every cell is a pair of times under a stated basis, a named condition, or a statement that the history is too small. |
| Time of the status view | The release with this Spec no more than 0.10 s above the release before it, at the median of ten alternated runs each (0.83 s measured on 0.145.1) | `roundfix follow-up --status --format json` on the same checkout and Run Database with the two binaries, as in the [recorded baseline](references/2026-11-16-status-view-time-baseline.md). |

The last row is the maintainer's standing condition, "com queda de performance nula ou quase nula", applied to this command. A difference under 0.10 s cannot be told from chance with ten runs: the baseline puts the standard error of a ten-run median at 0.06 s.

## Feature Assessment

| Criteria        | Question                                            | Score  |
| --------------- | --------------------------------------------------- | ------ |
| Impact          | How much more valuable does this make the product?  | Strong: the delivery plan and its timing are read in one look, from recorded history instead of an agent's guess, where the view holds no time today and the median Clean Run is about an eighth of the median time from plan to merge. |
| Reach           | What % of users would this affect?                  | Strong: every adopter that delivers through the Delivery Queue receives the table with the binary. The forecast needs that repository's own delivered Specs on that machine, and the rows of Specs being authored need the skills of Spec 0397. |
| Frequency       | How often would users encounter this value?         | Must do: 30 Specs were delivered from 2026-11-08 to 2026-11-15, about four a day, so the queue and its timing change every few hours. |
| Differentiation | Does this set us apart or just match competitors?   | Maybe: cycle-time percentiles are established practice, and Jenkins and Mergify estimate one job or one pull request. A forecast carried through planned Specs with named conditions was not found in the tools read. |
| Defensibility   | Is this easy to copy or does it compound over time? | Strong: each delivered Spec adds to the history the forecast reads, and the backtest makes its accuracy a measured property of the repository. |
| Feasibility     | Can we actually build this?                         | Strong for the table and the conditions, which come from recorded state. Maybe for the forecast until the backtest passes: a queue item's times are not recorded, per-Task times are kept 7 days here, and the Run Database is one machine's. |

Verdict: proceed in the shape the maintainer's six answers gave, with the forecast gated by the backtest; three results would change it. A backtest that misses 80% with every basis the history supports would ship the table without clock forecasts and record what the history lacks. A limit wider than the third KPI allows would call for a basis that accounts for a Spec's size before release. A status view slower than the band would move the forecast behind an explicit request.

## Trade-off Insights

- **Recommended approach:** build the table from state Roundfix already records, forecast from this repository's own delivered Specs as two percentiles, name a condition wherever a wait cannot be dated, and accept the method only after it holds on history.
- **Key trade-offs:**
  - The merge as the end against the Run's end or the release. The merge is when the Spec is on the default branch; the median Clean Run is about an eighth of the median time from plan to merge, and a release point is the maintainer's decision.
  - Two times in a cell against one. The table is wider, and the reader sees both what is likely and what to plan by.
  - Conditions against earliest times. On a day with waits the table shows fewer clock times, and each one it shows is backed by history.
  - Rows from what the repository records against rows from the Session Handoff. The recorded rows can be checked; the maintainer's own example row came from a handoff and appears only once its Authoring Phase is recorded or its Backlog Entry exists.
  - A backtest against a log of forecasts. The view writes nothing, so a shown forecast is never stored; accuracy is proven by replay before release and again by the operator afterwards.
- **Risks identified:**
  - A forecast is read as a commitment. The pair of times, the stated basis and the conditions are the defence, and the PRD states that the limit is missed about one time in seven by design.
  - The history is small or stops describing the present. Run Retention keeps 7 to 30 days, the Runs of 2026-11-09 to 2026-11-11 ran on another machine, the orchestration model changed on 2026-11-14 and the share of first Runs that end Clean rose from 60% to 79% inside the baseline.
  - Specs differ in size. The baseline warns that "a Run's duration is not divided by its Tasks", so one distribution for every Spec can be too wide to plan by.
  - Delays come in runs. When another repository holds the machine, every item behind waits, and a chained limit is then too narrow at the end of the queue.
  - The Phase Record is written by a command the skills run. A phase the skills did not record, or a Spec authored on another machine, has no entry here.
  - The note prints repository text, a blocker or a Waiting Declaration, into a report an Agent reads. The PRD states it as data, bounded in length, and never as an instruction.
  - Reading 30 days of history on every status call can cost time. The last KPI bounds it.
- **Dissenting view:** the forecast rests on times Roundfix does not record for a queue item, so the simpler shape, the table now and the forecast after a measurement, would have been the safer first Spec. The scope as decided keeps the forecast in V1 and lets the backtest refuse it. A second dissent: a wait on another repository's delivery or on the maintainer's own answer turns every row behind it into a condition, so on such a day the table's value is the queue and the notes. A third: the wait that matters most here belongs to another repository, which only the machine-wide table could forecast.
- **Stretch goal (V2+):** one table for the whole machine, across the repositories that share the Roundfix Home, which turns a Run Slot Wait from a condition into a forecast. After it, the recent hit rate in the footer, recomputed from history on each call.

## Opportunity Scan

- **Original** — one table in this repository's status view, with a forecast from its own history. Medium effort, one Spec. Impact Strong, Reach Strong, Frequency Must do, Differentiation Maybe, Defensibility Strong, Feasibility Strong for the table and Maybe for the forecast.
- **More ambitious** — the whole machine's queue: one table across every repository on the machine, so the wait for a slot is forecast and not only named. Large effort. Impact Must do for someone who supervises several repositories on one machine; Feasibility Maybe, because the Follow-up Command reads one repository and each repository has its own Spec folders.
- **Simpler** — the table now and the forecast later: status, note, actual start and elapsed time, with the forecast columns in a second Spec after the accuracy is measured. Small effort. Impact Maybe: it leaves out what the request asked for.
- **Adjacent** — tell the maintainer when a Spec in progress outlives the repository's usual delivery time, pushed as a notification and not read on demand. Small effort. Impact Strong where nobody is watching; it needs the same forecast first.

**Chosen direction:** original — the maintainer's pick: "Tabela com previsão (Recommended)". The more ambitious alternative is the stretch goal, the simpler one is the fallback if the backtest fails, and the adjacent one is left for a later request.

## Integration with Existing Features

| Integration point | How |
| ----------------- | --- |
| The Follow-up Command's status view | Gains the table in its text and the same rows in its JSON, added beside the existing fields and never renaming one. The follow-up list is unchanged. |
| The `--status` Skill Mode | Relays the table. It still accepts no handoff and consumes no message. |
| The Delivery Queue | Gives the first rows, their order, stage, blocker, Park Class, Run Slot Wait and Pending Question. The table never reorders the queue. |
| The Run Database and Run Retention | Give when each Spec's Runs began and ended. The retention window is the window of the history. |
| The Phase Record of Spec 0397 | Gives the Authoring Phase a Spec is in before its PRD exists. |
| Backlog Entries and the Waiting Declaration | Give the last rows and their notes. |
| The Release Plan | Unchanged: the next-release line still says what a release brings, and the forecast ends at the merge. |
| `roundfix runs list` | Already prints a start and a duration in the machine's time zone; the table follows the same convention. |

## Prerequisite and Shared Files

Spec 0397 is planned on the default branch and is delivered before this one. F7 reads the Phase Record it adds, so the PRD names Spec 0397 under its Prerequisites, and the Task Graph requires it if a Verification of this Spec depends on the record.

Both Specs are expected to change these files, and the PRD and the TechSpec confirm the list against the Task files of Spec 0397:

- `CONTEXT.md`;
- the Roundfix skill: the version fields of `.agents/skills/roundfix/SKILL.md` and `.agents/skills/roundfix/references/modes.md`, with their mirrors under `skills/roundfix/`;
- `skills/testdata/owned-skill-versions.json`;
- `docs/references/behavior-surfaces.json`;
- the Skill Coverage Map, `docs/references/skill-coverage.json`, and the Coverage Record, `docs/references/coverage-record.json`.

The Task files of Spec 0397 name no source file, user-guide page or skill reference of the Follow-up Command.

## Out of Scope (V1)

- **The whole machine's queue in one table** — the maintainer chose this repository's table. Another repository's hold on the machine appears as a condition.
- **Rows that exist only in a Session Handoff or an ai-memory message** — not chosen. The status view reads what the repository and the Phase Record hold, and the mode still accepts no handoff.
- **A forecast for the release** — the end is the merge. When a release goes out is the maintainer's decision, and the next-release line already reports it.
- **A forecast end for a Backlog Entry** — nobody has decided to build it yet.
- **Storing forecasts, or showing the hit rate in the status** — the view writes nothing, and the maintainer chose the backtest alone. The hit rate is part of the stretch goal.
- **Sizing a Spec by its text** — a forecast comes from recorded durations only. No judge is asked and no text leaves the machine.
- **A notification when a Spec passes its limit** — the adjacent alternative, left for a later request.
- **Changing the queue's order, the Follow-up Priority, the follow-up list or `roundfix deliver status`** — the table reports the queue and changes nothing else.

## Idea-stage judgment

- item: the request in the session
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: unsure (0.37)
- unverified_assumption: unsure (0.39)
- market_question: no (0.07)
- verdict: undecided
- stage: ran
- shape: changed — a forecast became a likely time and a limit that end at the delivery merge, a wait that cannot be dated became a condition, the rows grew to Specs being authored and open Backlog Entries, and the forecast ships only after a backtest.

The verdict was undecided, and the maintainer chose "Rodar a Ideia" on 2026-11-16. The judged text is the adopted Backlog Entry from its title on: [`references/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md`](references/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md).

## Decisions

Every answer is the maintainer's, quoted as recorded in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md).

- Route of the request (2026-11-15): "Backlog agora, Spec depois (Recommended)". When to author it (2026-11-15, ~20:15 BRT): "Logo após o plano da 0397 (Recommended)".
- The Idea stage (2026-11-16, ~08:00 BRT): "Rodar a Ideia" — "write-idea em Fable xhigh com perguntas antes do PRD, como na 0397."
- Authorization (2026-11-16, ~08:00 BRT): "Mesmas do BE (Recommended)" — "Push, PR, merge, release minor depois da Spec, até 3 correções; Makefile/CI/go.mod/lint só se a correção exigir; com o Quaestor ocioso, entregue depois da 0397." The PRD quotes it in the authorization record.
- Size of V1 (2026-11-16, ~09:15 BRT): "Tabela com previsão (Recommended)" — "O pedido original, neste repositório: status, observação e previsão pelo histórico de Runs do próprio repositório."
- Rows beyond the queue (2026-11-16, ~09:15 BRT): "Autoria e backlog (Recommended)" — "Entram Specs em Ideia, PRD ou TechSpec (pelo registro de fase da 0397) e backlog aberto, marcados 'ainda não é Spec', com a espera declarada e sem previsão de fim."
- The forecast end (2026-11-16, ~09:15 BRT): "Merge no main (Recommended)" — "Fim é a entrega integrada (Run, revisão, checks e merge); a release segue na linha 'próxima release' do status."
- The form of a forecast (2026-11-16, ~09:15 BRT): "Provável e limite (Recommended)" — "Dois horários por célula, ex.: '~10:30, até 12:10' (mediana e percentil 85), com a base no rodapé."
- Waits (2026-11-16, ~09:15 BRT): "Condição e duração (Recommended)" — "Início 'depois de \<condição\>' e fim como duração após o início; as de trás encadeiam igual, sem horário inventado."
- Validation (2026-11-16, ~09:15 BRT): "Backtest antes (Recommended)" — "No PRD, o método é repetido sobre ~100 Specs dos últimos 30 dias e só é aceito se o fim real cair dentro do limite em pelo menos 80%."
- Standing addition to the success criterion, made for Spec 0397 on 2026-11-16 and carried to this Spec by the operator's brief: "com queda de performance nula ou quase nula". It is recorded in the answers adopted by Spec 0397 (`docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/references/2026-11-15-maintainer-answers.md@f4580437cd6f3c354c2524f7f2eca9c924446d64`) and is read here as the time of the status view.

The author's own, made in this stage:

- The question that opened the stage showed the judgment and did not offer the dismissal again, because the maintainer had already chosen "Rodar a Ideia".
- The Backlog Entry is adopted here: its status is `promoted`, it names this Spec, and it moved into `references/`. Its Waiting Declaration was removed, because the wait it declared ended when the maintainer ordered the authoring.
- No ADR is written at this stage. ADR-0409 and ADR-0410 are reserved for the PRD and the TechSpec.

## Open Questions

- Which recorded times give a delivered Spec its start and its end? Default: the start is the Spec's first Implement Run and the end is its delivery merge in the Git history. The PRD's backtest shows that both can be read for about 100 Specs, or says what is missing. If the history cannot support the forecast, the author returns to the maintainer before changing what the Delivery Queue records.
- What does a forecast depend on besides the history as a whole? Default: the backtest compares one distribution for every Spec with one that accounts for the number of Tasks, and keeps the plainest one that meets the first three KPIs.
- How do rows chain when one item's review overlaps the next item's Run? Default: the backtest measures how consecutive deliveries followed one another and the forecast uses that, with no assumption of a strict series.
- How much history is enough? Default: 11 delivered Specs inside the Run Retention window on this machine, the count at which the next duration falls inside the seen range with a probability above 80%. Below it the cell reads "not enough history" with the count.
- What does a Spec being authored show? Default: its Authoring Phase as the status, "after its plan is merged and queued" as the start and a duration as the end, under the answer about waits. This reads "ainda não é Spec" and "sem previsão de fim" in the answer about rows as said of Backlog Entries. If the maintainer meant both kinds, a Spec being authored shows no forecast end either.
- What does a Spec in progress show once it has passed its limit? Default: "past the limit" with how long ago, never a time that is already gone.
- In which order do the rows beyond the queue come, and how many Backlog Entries are listed? Default: Specs with a plan that are not queued, then Specs being authored, both by number, then every open Backlog Entry by date. Merged items of the current queue stay out, as today.
- What is a row for a Spec that has a folder and no Phase Record entry? Default: a row with the status "being authored" and no phase named.
- Which waits become named conditions? Default: a Run Slot Wait with its holder, a park with its Park Class, an unmerged prerequisite, a Pending Question, a Spec that is not queued, a Spec without Tasks and a Backlog Entry's Waiting Declaration.
- Which names does the glossary take? Default: the PRD's Glossary settles them. The candidates are a name for the pair of times of a row, one for what a forecast was computed from, and one for a wait shown in place of a time.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md`](references/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md), the Backlog Entry that holds the request, recorded by the operator.
2. [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md), the answers behind `## Decisions`.
3. [`references/2026-11-16-status-view-time-baseline.md`](references/2026-11-16-status-view-time-baseline.md), recorded by the author: the time of the status view in `## KPIs`.

Owned by Spec 0397 and linked at its path: the delivery speed baseline (`docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/references/2026-11-16-delivery-speed-and-plan-quality-baseline.md@f4580437cd6f3c354c2524f7f2eca9c924446d64`), the source of every Run and plan-to-merge number in this file. Spec 0397 can be archived before this plan reaches the default branch, and its Archive Relink rewrites only the links it finds there, so the PRD stage checks the two links into its folder after rebasing.

This project's ai-memory messages were listed on 2026-11-16 and none was pending. The two other open Backlog Entries, a false status reported without a reproduction and the cost of starting a session, belong to other contexts and stay where they are. The session listed no Research Subagent, so the author did the sweeps.

Secondbrain (read on 2026-11-16; `wiki/index.md`, then three `qmd query` runs, which returned nothing on delivery forecasting):

- `wiki/concepts/practical-process-automation-capitulo-14-11-process-visibility.md` — an operator "must have knowledge of the current process state at all times, and the ability to use this knowledge effectively in predicting future process states". It supports putting the forecast where the state is read.
- `wiki/sources/exa-roundfix-reestruturacao-fila-2026-09-24.md` — earlier research on the Spec queue; a search of it finds nothing on forecasting.

External sources (read through Exa on 2026-11-16):

- The Kanban Guide, May 2025, https://kanbanguides.org/english/ — the service level expectation; it gave the forecast its two parts and its source in history.
- ProKanban, "The Benefits of Defining Your Service Level Expectation", https://prokanban.org/blog/the-benefits-of-defining-your-service-level-expectation — a range and a level of confidence, for one item; it challenged chaining down a queue.
- Nick Brown, ASOS Tech Blog, "The Full Monte", https://medium.com/asos-techblog/the-full-monte-901d721b8532 — forecasts checked against outcomes; it supports the two percentiles and the 80% bar.
- Robert Herbig, SEP, "Improving Software Forecasting With Monte Carlo Simulation", https://sep.com/blog/software-forecasting-my-process/ — how many samples are enough; it gave the default of 11.
- MerceMay, "Forecasting Completion Dates from Historical Throughput", https://mercemay.com/articles/software-engineering/engineering-process/forecasting-completion-from-throughput/ — what a short history cannot show and what correlated delays do to a chain; both are risks above.
- Jenkins, issue JENKINS-21099, https://issues.jenkins.io/browse/JENKINS-21099, and its estimate from the last three successful builds, https://devops.stackexchange.com/questions/3768/how-many-builds-does-jenkins-need-to-run-before-it-has-estimated-times — no estimate without enough successes.
- Mergify, Merge Queue API, https://docs.mergify.com/api/merge-queue/ — an estimate that is `null` when it cannot be made, beside a named waiting state.
- GitLab, work item 504119, https://gitlab.com/gitlab-org/gitlab/-/work_items/504119 — a job's estimated end as a request of 2024.
