---
spec: 0398-the-status-view-shows-the-spec-queue-with-forecasts
status: active
created: 2026-11-16
surfaces: [backend, cli, docs]
type: feat
---

# The status view shows the Spec queue with a forecast start and end

The status view of the Follow-up Command holds no time, and it lists a Spec only once its folder holds a PRD. Whoever plans around the delivery asks the Supervisor, who pieces the timing together from several commands and its own context. After this Spec, the status view shows one table of the repository's Spec queue, oldest scheduled first: Delivery Queue items, planned Specs, Specs being authored and open Backlog Entries, each with its status and a short note. A Spec carries a forecast start and end as a likely time and a limit, computed from this repository's own delivered Specs, and a wait Roundfix cannot date is shown as a named condition with a duration. The view still writes nothing, and the `--status` Skill Mode relays the table as the command printed it.

## Sources

Adopted by this Spec's planning change (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md`](references/2026-11-15-the-status-mode-shows-the-spec-queue-with-forecasts.md), the Backlog Entry that holds the maintainer's request, recorded by the operator.
2. [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md), the maintainer's answers, recorded by the operator.
3. [`references/2026-11-16-forecast-backtest.md`](references/2026-11-16-forecast-backtest.md), recorded by the author: the replay of the method over 109 delivered Specs.
4. [`references/2026-11-16-status-view-time-baseline.md`](references/2026-11-16-status-view-time-baseline.md), recorded by the author: the time of the status view before this Spec.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from. At the authoring base this project's ai-memory inbox is empty and `docs/findings/` is empty. The two other open Backlog Entries share no component with this Spec: a Conexus report of a false status that waits for reproduction steps, and the start cost of a session, which waits for a week of numbers.

## Prerequisites

Spec 0397 is delivered before this one, and the Tasks stage names it in `requires`: this Spec reads the Phase Record that Spec 0397 adds and extends the package that owns it. Authored against `f4580437`, while Spec 0397 was in delivery.

Both Specs change these files, each inside its own sentences and fields:

- `CONTEXT.md`;
- the Roundfix skill: the version fields of `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, and the `modes` reference in `.agents/skills/roundfix/references/modes.md` with its mirror;
- `skills/testdata/owned-skill-versions.json`;
- `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json`;
- `internal/cli/cli.go`, the package directory `internal/phaserecord`, and ADR-0407, which gains this Spec's amendment notice.

No version number is written by hand: the record commands raise it on the tree they run on.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; a row carries a Spec slug, a Backlog Entry's path and times that other components recorded. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; the status view asks no judge, and its one existing call to the forge is unchanged. Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0409 (new) records the forecast and where a time may be shown (ADR-0409: "a named condition with a duration everywhere else"). ADR-0410 (new) records the rows (ADR-0410: "it lists the Specs being authored, read from the Phase Record and the Spec folders, and the open Backlog Entries") and amends two records. The status view read no Backlog Entry (ADR-0353: "reads no Backlog Entry or Finding"), and it now reads the open ones. Roundfix reported nothing from the Phase Record (ADR-0407: "Roundfix reports nothing from it"), and the status view now names an Authoring Phase from it; it reports no usage, so the rest of that paragraph stands. The status view still takes nothing from ai-memory (ADR-0338: "claim no handoff and consume no message"). Run Retention sets the window of the history unchanged (ADR-0255: "is removed whole in one write transaction per Run"). Not applicable: ADR-0383, ADR-0390 and ADR-0404, which cite ADR-0353; and ADR-0290, ADR-0296 and ADR-0301, which cite ADR-0255 (reasons in the TechSpec's ADR Alternatives). Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-16 quoted in the authorization record, which keeps the grant of the owned skills. No Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0398-the-status-view-shows-the-spec-queue-with-forecasts/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/follow-up.md`, `.agents/skills/roundfix/references/modes.md`, `skills/roundfix/SKILL.md`.

## Goals

- The person who supervises Spec work reads the whole queue in one table: what is being delivered, what is planned, what is being authored and what is still a Backlog Entry, in the order it will be worked.
- A Spec in progress, and a planned Spec behind it, show when they are likely to start and to merge and the time to plan by.
- No time is invented: where Roundfix cannot date a wait, the row names the wait and says how long the Spec takes once it starts.
- Every forecast says what it was computed from, and the method is the one that held on this repository's history.
- The status view stays read-only and takes no more time than it does today, within what a measurement can tell.

## User Stories

1. As a person who supervises Spec work, I want one table of the Spec queue in the status view, so that I stop piecing the plan together from several commands.
2. As that person, I want a Spec in progress to show when it will likely merge and a limit, so that I know when to look again.
3. As that person, I want a planned Spec behind a Spec in progress to show when it will likely start and merge, so that I plan my own time around it.
4. As that person, I want a wait that nobody can date shown as what it waits on, so that I am not given a time that means nothing.
5. As that person, I want a Spec being authored and an open Backlog Entry in the same table, so that the queue I plan by is the queue I see.
6. As a reader of a forecast, I want its basis stated, so that I know how much to trust it.
7. As an adopter with little history, I want to be told that the history is too small, so that I do not plan by three deliveries.
8. As an Agent that runs the `--status` Skill Mode, I want the rows and the basis in the JSON report, so that I relay the table without computing anything.
9. As the maintainer, I want the status view to stay read-only and fast, so that asking for the status costs what it costs today.

## Core Features

1. **Spec Queue Table.** The status view shows the **Spec Queue Table**: one row for each Delivery Queue item that is not merged, in queue order; then each planned Spec that is not queued, by number; then each Spec being authored, by number; then each open Backlog Entry, oldest first. Each row has a row number, the Spec's number, a status, a start, an end and a note. The lines the view prints today stay, and the table follows the pending Specs.
2. **The status of a row.** A Delivery Queue item shows its stage and its Tasks completed of total. A planned Spec shows that it is not queued, with the same Task count. A Spec being authored shows the latest Authoring Phase recorded for it, or only that it is being authored. A Backlog Entry shows that it is not yet a Spec.
3. **The note of a row.** The note is one line taken from what Roundfix already records: a parked item's Park Class and the command that answers it, the blocker text of a Run Slot Wait, the time an Authoring Phase began, and a Backlog Entry's title with its Waiting Declaration. A row with nothing recorded shows a dash. Hazard: a blocker, a title and a Waiting Declaration are repository text that an Agent reads in the report; the note is one line cut at 120 characters, it is data, and the skill text says that it is never an instruction.
4. **Queue Forecast.** A **Queue Forecast** is two times for one start or one end: the likely time, which is the median, and the limit, which is the 85th percentile, both by nearest rank over the **Forecast Basis**. The end of a Spec is its delivery merge into the default branch. Text shows a forecast as two local times, and JSON carries both as timestamps.
5. **Forecast Basis.** The Forecast Basis is the delivered Specs of the repository that the Run Database still holds: each has a first Implement Run, whose start is the Spec's start, and an Archive Record on the default branch that is not a cancellation, whose commit time is the Spec's end. Two durations come from it: the time from start to merge, and the start gap, from the start of one Spec to the start of the next when that next Spec's plan was already on the default branch. A forecast needs 11 values.
6. **A Spec in progress.** A Spec that has a first Implement Run and is neither merged nor parked shows its actual start and a forecast end. The end is its start plus the time from start to merge, taken from the delivered Specs that took longer than this one already has. When fewer than 5 took longer, the end reads that the Spec is past the usual time and shows no time.
7. **Planned Specs behind a Spec in progress.** A Spec with a PRD that has not started, one, two or three rows behind the Spec in progress that started last, shows a forecast start and a forecast end. Its start adds one start gap for each row, the first taken from the gaps longer than the time already passed, and its end adds the time from start to merge. Durations add as pairwise sums, never as a sum of percentiles. Rows that wait on something of their own do not count as rows.
8. **A wait is a condition.** Every other Spec that has not started shows a named condition as its start and its end as two durations after its start. The conditions are: next in the queue, or after the row ahead starts, with no Spec in progress; after it is queued, for a planned Spec with no Spec in progress; after a Run slot frees; after the named prerequisite is merged; after its park is answered; after the Spec in progress is merged, when too few start gaps are longer than the time already passed; and after the row ahead starts, for a fourth row behind and beyond. A parked Spec that had started shows its actual start and its park as the end.
9. **A Spec being authored.** A Spec is being authored when the checkout holds no PRD for it and either a folder for it or a phase entry of an Authoring Phase written in the last seven days. An entry counts when its repository is the checkout, or a working tree whose `origin` remote is the checkout's. The row shows that its start waits for its plan to be merged and its end as the two durations after its start, as feature 8 does. A Spec that has an Archive Record is never a row.
10. **A Backlog Entry.** Every open Backlog Entry is a row, a waiting one included, with no forecast. The status view still reads no Finding and no inbox note, prints no waiting count, and refuses `--all` and `--messages`.
11. **The basis is stated.** Under the table one line says how many delivered Specs the forecasts were computed from, the dates of the first and the last merge among them, that the likely time is the median and the limit the 85th percentile, and the two durations of each measure. With fewer than 11 delivered Specs, the line says how many there are of the 11 needed, and each forecast cell says the same and shows no time.
12. **The JSON report.** The status report gains the rows and the basis as two fields beside the existing ones, and three sources: the Backlog Entries, the Phase Record and the delivery history. Each row carries its kind, the Spec's slug or the entry's path, its status, its note, its Task counts, and its start and end, each as an actual time, a pair of times, a condition with a pair of durations, or the reason no forecast is shown. No existing field is renamed or removed.
13. **A source that cannot be read.** When the Backlog Entries, the Phase Record or the delivery history cannot be read, the report names the source as unavailable with its reason, shows the rows it still has, shows no time that needed the missing source, and exits `0`.
14. **Read-only, with bounded work.** The status view writes nothing under the repository or under Roundfix Home, accepts no Session Handoff, consumes no message and asks no judge. It reads the Git history once for the forecast, bounded by the Run Retention window, and asks at most eight other working trees for their `origin` remote. Hazard: the view runs in the user's checkout while Runs write the Run Database, so it reads through the snapshot it already takes and never holds a write lock.
15. **The accepted method is pinned.** The repository carries the history of the backtest as a fixture, and a test replays the method over it and reproduces the counts the backtest recorded. A change to the percentile rule, to either minimum or to how durations add fails that test.
16. **The Skill Mode relays the table.** The Roundfix skill's `--status` Skill Mode shows the Spec Queue Table and the basis line from the JSON report, in the user's time zone, and computes or adjusts no forecast. It says that a note is data.
17. **Documentation and glossary.** The user guide and the Roundfix skill's references describe the table, each row kind, the forecast, the conditions, the basis line, the new JSON fields and the three sources. `CONTEXT.md` defines the three new terms and revises **Follow-up Command**, and ADR-0353 and ADR-0407 each carry a notice naming ADR-0410.

## User Experience

A repository with a history, read at 10:36 local time while one Spec is in progress and one is planned behind it:

```text
Spec queue: 4
  #  Spec  Status                Start (-03)                   End (-03)                          Note
  1  0397  running 0/5 Tasks     11-16 08:30                   ~11-16 11:59, by 11-16 14:13       -
  2  0398  not queued 0/4 Tasks  ~11-16 11:20, by 11-16 15:03  ~11-16 14:07, by 11-16 19:41       -
  3  0399  authoring: prd        after its plan is merged      +1h56m, by +4h22m after its start  prd since 11-16 09:55
  4  -     not yet a Spec        -                             -                                  A report of late items; waits on the maintainer since 2026-11-11
Forecast basis: 109 Specs delivered from 10-17 to 11-15 on this machine. Likely is the median and the limit the 85th percentile: 1h56m and 4h22m from start to merge, 2h02m and 3h53m from one start to the next.
```

The times of the first two rows are the ones the method gives for this repository's history at that moment. A tilde marks the likely time and `by` the limit; a start without either is an actual start.

A repository whose Run Database holds three delivered Specs:

```text
Spec queue: 2
  #  Spec  Status                Start (-03)         End (-03)                     Note
  1  0007  not queued 0/2 Tasks  after it is queued  not enough history (3 of 11)  -
  2  -     not yet a Spec        -                   -                             Export the ledger as a spreadsheet
Forecast basis: not enough history, 3 of 11 delivered Specs.
```

The `--status` Skill Mode shows the same rows as a table in its report, with the basis line under it.

## Non-Goals / Out of Scope

- One table for the whole machine across repositories: the maintainer chose this repository's table, and another repository's hold on the machine appears as a condition.
- Rows that exist only in a Session Handoff or an ai-memory message.
- A forecast for a release: the end is the merge, and the next-release line is unchanged.
- A forecast for a Backlog Entry.
- Plan to merge as a basis, and a clock time for a planned Spec with no Spec in progress ahead: both were refused by the backtest.
- A forecast that depends on the number of Tasks or on a Spec's text, and any judge call.
- Storing forecasts, a log of forecasts against outcomes, and the hit rate in the status.
- Recording when a Delivery Queue item starts or reaches a stage.
- A notification when a Spec passes its limit.
- Merged items of the current queue as rows.
- The follow-up list, `roundfix deliver status` and `roundfix runs list`: unchanged.
- Findings and inbox notes in the status view.
- A setting for the percentiles, the minimums or the seven days.

## Success Metrics

1. Success Metric: in a temporary repository and home that hold one Delivery Queue item, one planned Spec, one Spec folder without a PRD, one open Backlog Entry and one waiting Backlog Entry, the status view prints the Spec Queue Table with those five rows in that order, each with its status and note, after the pending Specs and before the next release, and every line the view printed at the delivery base is still printed.
2. Success Metric: the replay test reads the fixture of 109 delivered Specs and reproduces the recorded counts: the limit holds for 80 of 98 ends forecast at the start; for 55 of 63 starts and 53 of 63 ends one row behind; for 36 of 39 starts and 36 of 39 ends two rows behind; and for 19 of 19 starts and 17 of 19 ends three rows behind. The test fails when the percentile rule, the minimum of 11 or the pairwise sum is changed.
3. Success Metric: with a fixture history of at least 11 delivered Specs, a Spec in progress shows its actual start and an end whose two times are its start plus the median and the 85th percentile of the times from start to merge longer than the time already passed; with fewer than 5 longer, the end reads that the Spec is past the usual time.
4. Success Metric: with the same history, planned Specs one, two and three rows behind a Spec in progress each show a start and an end as pairs of times that match the pairwise sums; the fourth shows a condition and two durations; a row that waits on a prerequisite is skipped by the count.
5. Success Metric: each condition of Core Feature 8 appears for its case with the two durations after the start, and no cell of those rows holds a clock time.
6. Success Metric: with fewer than 11 delivered Specs every forecast cell and the basis line say how many there are of 11, and no cell holds a time or a duration.
7. Success Metric: a Spec is a row as being authored for a folder without a PRD, for a phase entry recorded in the checkout and for one recorded in a clone with the same `origin` remote; it is not a row for an entry of another remote, for an entry older than seven days, or when it has an Archive Record; with a PRD it is a planned row.
8. Success Metric: an open Backlog Entry whose Waiting Declaration holds 300 characters and a line that reads as an instruction is printed as one note line of at most 120 characters; an entry that is not open is not a row; a Finding is never a row.
9. Success Metric: the JSON report carries the rows, the basis and the three sources, and every field it carried at the delivery base, with the same value, for the fixture of metric 1.
10. Success Metric: with an unreadable Phase Record, with a Git history that cannot be read, and with an unreadable Backlog Entry, the report exits `0`, names the source as unavailable and prints the rows that remain.
11. Success Metric: across a run of the status view no file under the repository or the home changes, no judge request is made, the Git history is read once for the forecast whatever the number of delivered Specs, and `--all` and `--messages` are refused as before.
12. Success Metric: the user guide and the Roundfix skill's references describe the table, the row kinds, the forecast, the conditions, the basis line and the JSON fields; `CONTEXT.md` defines **Spec Queue Table**, **Queue Forecast** and **Forecast Basis** and the revised **Follow-up Command**; ADR-0353 and ADR-0407 carry an "Amended by ADR-0410" notice; the documentation, skill and glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write, read through Exa on 2026-11-16.

- The Kanban Guide, May 2025 (https://kanbanguides.org/english/): a service level expectation is "a forecast of how long it should take a work item to flow from started to finished", and it "should be based on historical cycle time". This is why a forecast is a duration from a Spec's start, read from the repository's own deliveries.
- The same guide: "The SLE itself has two parts: a period of elapsed time and a probability associated with that period". This is why each forecast names its percentile.
- ProKanban, *The Benefits of Defining Your Service Level Expectation* (https://prokanban.org/blog/the-benefits-of-defining-your-service-level-expectation): "The key is to maintain communicating it as a range and level of confidence so as not to create confusion toward a guarantee." This is why a cell holds two times and the basis line is printed.
- Nick Brown, ASOS Tech Blog, *The Full Monte* (https://medium.com/asos-techblog/the-full-monte-901d721b8532): of 1,200 forecasts checked against outcomes, the 85th percentile was right in 361 of 400 and the 50th in 270 of 400. This is the published measure the 80% bar is read against.
- Jenkins, issue JENKINS-21099 (https://issues.jenkins.io/browse/JENKINS-21099): with too few successful builds Jenkins should "just give no estimate at all". This is why a small history shows no time.
- Mergify, *Merge Queue* API (https://docs.mergify.com/api/merge-queue/): `estimated_merge_at` is "`null` when it cannot be estimated". This is why a wait that cannot be dated is a named condition.

The backtest and the time baseline this planning change commits under `references/` are this Spec's own evidence, not outside evidence; `## Validation` reports them.

## Unreachable Acceptance

- criterion: over the first 20 Specs delivered after the release, the limit of the end forecast at each Spec's start holds for at least 16 and the likely time for 8 to 12.
  criterion: over the same 20 Specs, the limit is at most 2.5 times the actual time from start to merge at the median.
  criterion: on the maintainer's checkout and Run Database, the status view of the released binary takes no more than 0.10 s above the release before it, at the median of ten alternated runs each.
  criterion: in a live Delivery Queue with a Spec in progress and a planned Spec behind it, the released status view shows a start and an end for the Spec behind.
  reason: each needs the released binary, the maintainer's Run Database with real deliveries and days of delivered Specs, which no Verification or QA row may reach.
  satisfied-by: the operator's repeat of the backtest count on the Run listing and the Git history after 20 deliveries, and the timing and the live report recorded in the session.

## Validation

The decisions were checked against these sources on 2026-11-16.

- The backtest ([`references/2026-11-16-forecast-backtest.md`](references/2026-11-16-forecast-backtest.md)), a walk-forward replay over 109 Specs delivered from 2026-10-17 to 2026-11-15. The limit held for 80 of 98 ends forecast at the start (81.6%), the likely time for 46 (46.9%), and the limit was 1.91 times the actual at the median. Re-made one to four hours after the start it held for 80.0% to 88.0%. For a Spec one, two and three rows behind a Spec in progress it held for 84.1%, 92.3% and 89.5% of the ends. Each kind meets the three measures the Idea set: at least 80%, 40% to 60%, and at most 2.5.
- The same replay refused three bases. Plan to merge held for 74 of 96 (77.1%). The start and the end of a planned Spec with no Spec in progress ahead held for 59 and 60 of 78. Each spans a wait that ends when a person acts.
- The size of a Spec did not earn a basis of its own: by the number of Tasks the limit held for 81 of 98, against 80, with a wider limit.
- The time of the status view ([`references/2026-11-16-status-view-time-baseline.md`](references/2026-11-16-status-view-time-baseline.md)): a median of 0.83 s over five runs of Roundfix 0.145.1, with a standard error of 0.06 s for a ten-run median. Reading the Git history for 31 days of this repository took 0.04 to 0.05 s in three runs.
- Characterization at `f4580437`. The status view lists Delivery Queue items and Specs with a PRD, with a position, a stage and a Task count, and no time. A Delivery Queue item stores no time, and a new queue replaces a finished one with its items. A Run stores when it was created and completed, and the listing of a repository's Runs reads the terminal ones too. The status view reads no Backlog Entry, refuses `--all` and `--messages`, and its sources are six. A Run Slot Wait writes its holders into the item's blocker text.
- The Secondbrain was read on 2026-11-16 (`wiki/index.md`, then three `qmd query` runs): it holds nothing on delivery forecasting. `wiki/concepts/practical-process-automation-capitulo-14-11-process-visibility.md` supports reading the forecast where the state is read: an operator "must have knowledge of the current process state at all times, and the ability to use this knowledge effectively in predicting future process states".

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

## Glossary

- adds: **Spec Queue Table**
- adds: **Queue Forecast**
- adds: **Forecast Basis**
- changes: **Follow-up Command**

## Decisions

Autonomous route: no person present; decisions taken from the approved Idea, the maintainer's answers recorded in `references/2026-11-16-maintainer-answers.md`, the operator's cycle BF brief for Spec 0398 and the backtest above.

The maintainer's answers (AskUserQuestion), verbatim as the operator recorded them:

- 2026-11-15: "Roteamento" — "Backlog agora, Spec depois (Recommended)"; "quando autorar" — "Logo após o plano da 0397 (Recommended)".
- 2026-11-16: "Etapa de Ideia" — "Rodar a Ideia": "write-idea em Fable xhigh com perguntas antes do PRD, como na 0397."
- 2026-11-16: "Autorização" — "Mesmas do BE (Recommended)": "Push, PR, merge, release minor depois da Spec, até 3 correções; Makefile/CI/go.mod/lint só se a correção exigir; com o Quaestor ocioso, entregue depois da 0397."
- 2026-11-16: "Tamanho da V1" — "Tabela com previsão (Recommended)": "O pedido original, neste repositório: status, observação e previsão pelo histórico de Runs do próprio repositório."
- 2026-11-16: "Linhas além da fila" — "Autoria e backlog (Recommended)": "Entram Specs em Ideia, PRD ou TechSpec (pelo registro de fase da 0397) e backlog aberto, marcados 'ainda não é Spec', com a espera declarada e sem previsão de fim."
- 2026-11-16: "Fim previsto" — "Merge no main (Recommended)": "Fim é a entrega integrada (Run, revisão, checks e merge); a release segue na linha 'próxima release' do status."
- 2026-11-16: "Formato" — "Provável e limite (Recommended)": "Dois horários por célula, ex.: '~10:30, até 12:10' (mediana e percentil 85), com a base no rodapé."
- 2026-11-16: "Esperas" — "Condição e duração (Recommended)": "Início 'depois de \<condição\>' e fim como duração após o início; as de trás encadeiam igual, sem horário inventado."
- 2026-11-16: "Validação" — "Backtest antes (Recommended)": "No PRD, o método é repetido sobre ~100 Specs dos últimos 30 dias e só é aceito se o fim real cair dentro do limite em pelo menos 80%."
- 2026-11-16: "Spec em autoria mostra previsão?" — "Mostra previsão (Recommended)": "Uma Spec em autoria recebe o provável e o limite pelo histórico de plano-até-merge; só o backlog fica sem fim."
- 2026-11-16: "Aprovação da Ideia 0398" — "Aprovar e seguir (Recommended)".
- 2026-11-16: "Base da previsão de Spec em autoria" — "Aceitar a troca (Recommended)": "Spec em autoria mostra 'depois do merge do plano' + duração provável e limite da entrega; sua regra de 80% é respeitada."

The maintainer's standing condition of 2026-11-16, made for Spec 0397 and carried to this Spec by the operator's brief, is "com queda de performance nula ou quase nula". It is read here as the time of the status view.

Taken from those answers, the Idea and the backtest, with no further interview:

- A Spec being authored shows a forecast, as the maintainer decided, and its basis is the time from start to merge, not plan to merge. His first answer names the plan-to-merge history, and his answer on validation accepts a method only at 80%: plan to merge held for 77.1%. Asked on 2026-11-16, about 11:40 BRT, with that number, he accepted the swap. The row keeps the form he chose for a wait: "after its plan is merged" and the two durations of the delivery.
- The Idea's Open Questions. A Spec's start is its first Implement Run and its end is the commit that added its Archive Record; both were read for all 109 Specs. The forecast depends on nothing but the history: a basis by Task count held no better. Rows chain by the start gap, which the history gives whatever the queue overlapped. Eleven delivered Specs are the minimum, and five for a forecast re-made after time has passed. A Spec in progress past the usual time shows no time. The rows come in the order of Core Feature 1, every open Backlog Entry is listed, and merged items stay out. A Spec with a folder and no phase entry reads as being authored. The conditions are those of Core Feature 8. The glossary takes the three names of `## Glossary`.
- A clock time is shown only behind a Spec in progress, up to three rows: the replay held there and reached no fourth row.
- A phase entry counts for the repository by its `origin` remote, because this repository's own authors work in other clones, and only for seven days, the window Spec 0397 reads for its notice.
- The table is added to the existing report and the pending Specs lines stay, so a reader of the current text and JSON loses nothing.
- The basis line and the conditions are fixed English text; the Skill Mode relays them in the user's language as it relays the rest of the report.
- `type: feat` and Release Ceiling `minor`: the Spec adds report lines and JSON fields. ADR-0409 and ADR-0410 record the decisions; ADR-0410 amends ADR-0353 and ADR-0407.

## Open Questions

- Does the forecast hold in a repository that delivers differently from this one? Not measured. Default: the same method and minimums; each repository's basis line states its own history.
- Is a start gap the right chain when a queue holds several items and reviews overlap? Default: yes; the replay held for 55 of 63 starts, and the operator repeats the count after the release.
- Does the Git read stay inside the time band in a repository with a long history? Default: yes, because it is bounded by the Run Retention window; the third criterion of `## Unreachable Acceptance` is the check.
