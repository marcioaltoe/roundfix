---
spec: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
status: active
created: 2026-11-05
surfaces: [cli, docs]
type: feat
---

# A queue item shows its turn, and the pre-PR review gets its own capacity

On 2026-11-05 Fiscus reported a Delivery Queue of about twenty-one Specs that had run for two days; at one moment five items stood in `reviewing` while no review process ran, because one Implement Run of another item held the Queue Owner, and `roundfix deliver status` said nothing about the wait. The stage an item shows is the step it will take next, not proof that the step is running: a Delivery Retry records `reviewing` before any review starts, and the Queue Owner advances one item at a time through every stage, so a retried item waits behind the whole delivery of the item in front of it. After this Spec the pre-PR review, which only reads, runs beside the owner's other work for up to a configured number of items at once, and `deliver status` prints one line per item that says whether it is advancing, reviewing or waiting, and for what and in which place in line.

## Sources

One source, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-fiscus-delivery-queue-takes-days.md`](references/2026-11-05-fiscus-delivery-queue-takes-days.md),
   reported by fiscus on 2026-11-05 at 08:07 BRT, copied verbatim. This Spec
   answers its point 1; points 2 and 4 belong to Spec 0356, point 3 to Spec
   0357 and point 5 to Spec 0358.

No open Backlog Entry or unresolved Finding exists on `main`. This project's
ai-memory inbox held three messages on 2026-11-05: the Fluxus skills request
(Spec 0354), this Fiscus report, and an informational Pantheon note on its
`oxfmt` pin, which shares no context with this Spec. None was consumed.

## Prerequisites

`requires: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root`. Spec
0352 is being delivered now and adds the test helper that marks a directory
as a Git repository the configuration loader accepts; this Spec's Project
Config test uses it, because after 0352 an empty `.git` directory no longer
marks a repository. 0352 also changes the configuration loader this Spec
extends with one key.

Authored against `ad77a797`, where Specs 0352 and 0353 are active. Specs
0354, 0356, 0357 and 0358 are authored in parallel; the planned delivery
order is 0354, 0356, 0355, 0357, 0358. Files this Spec may share with them:
the Delivery Engine and its retry (0356 warns on a retry whose new head moves
the review base; 0357 adds a check where the queue revalidates an item),
`deliver status` and the retry output in the `deliver` command, the
configuration loader (0352), the Roundfix skill's `SKILL.md` pair, version
record and `deliver` reference, the `deliver` and configuration user guides,
`CONTEXT.md` and the Behavior Surface Record. This Spec keeps each edit to
those files narrow and rebases on the earlier Specs before its plan merges.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the Turn
  Record is named after the existing Queue Owner id, and items keep their
  Spec slugs. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the owner and `deliver status` read the Run Database,
  the configuration and a local file. Source: `docs/agents/go.md`,
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0304 keeps the Run Database
  schema unchanged for a wait, ADR-0304: "A new `waiting` stage or an item column was rejected: a schema change makes
  older item binaries refuse the shared Run Database.", so the wait this Spec shows lives in a file the owner writes, not in a stage or
  column. The owner still waits for each Prerequisite Spec before it starts
  an item, and the new turn line names that wait. ADR-0279 keeps `deliver status`
  read-only and its `Owner:` line for a Queue Owner proven gone, ADR-0279: "Status stays read-only and never reclaims the owner record", and the turn lines appear only for an owner proven live. ADR-0199
  keeps the Queue Token Ceiling a limit on new starts that does not count the
  review, ADR-0199: "Tokens outside Runs, such as the pre-PR review, do not count.", so reviews beside other work spend no queue tokens. ADR-0225 keeps each queue step running the binary its item branch builds, ADR-0225: "Before each of those three children it runs the
  build in the item worktree", so each review builds into its own item worktree as today. ADR-0311 is the
  model for the new operator-only key, ADR-0311: "a User
  Config key, 60 minutes by default, that the operator changes without a
  release.", and its Check Wait is unchanged. ADR-0237, which cites the ceiling, carries no obligation here: this Spec
  does not change the retry's merge observation, ADR-0237: "A Delivery Retry records a merge made outside the queue". ADR-0365, written by
  this Spec, records the decision. Source: `docs/agents/domain.md`,
  `docs/agents/cli.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 2 answers of 2026-11-05 quoted in the authorization record;
  the governed files are the Roundfix skill files the documentation Task
  changes. No Makefile, workflow, Go module, lint, formatter or Baseline
  module change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Goals

- A Delivery Queue item that is waiting says so in `deliver status`: for
  what (the owner's turn, a review slot or its prerequisites) and its place
  in line, and an item whose step is running says that instead.
- The pre-PR review of an item never waits for the Implement Run or the
  publication of another item: up to the **Review Capacity** of items are
  reviewed at once, beside the one item the Queue Owner advances through
  every other stage.
- A queued item still starts from a default branch that holds every earlier
  item that merged, so the queue's ordering guarantee and its conflicts are
  no worse than today.
- The operator sets the Review Capacity in User Config and sees it on the
  `Limits:` line of `deliver start` and `deliver status`.
- The guides, the Roundfix skill, `CONTEXT.md` and ADR-0365 tell the operator
  and Agents what a stage means, how to read the `Turn:` lines and how to set
  the Review Capacity.

## User Stories

1. As the operator of a long queue, I want `deliver status` to say which item
   the owner is advancing and which ones wait, for what and in which place,
   so that I can tell a stalled queue from a busy one.
2. As the operator, I want an item I retried at `reviewing` to be reviewed
   while another item's Implement Run is still running, so that a review
   never waits hours behind unrelated work.
3. As the operator, I want several items that wait for review to be reviewed
   at once, up to a number I choose, so that a batch of retries clears in
   one review's time.
4. As a Spec author, I want a queued item to start only after the earlier
   items in its queue have merged or parked, as today, so that it starts from
   their merged code.
5. As an Agent reading `deliver status`, I want the turn lines in a fixed
   form and the review concurrency on the `Limits:` line, so that I can parse
   the queue's state without guessing.

## Core Features

1. The stage is the next step. `deliver status` keeps printing each item's
   stage, which names the step the item is at; the Roundfix skill, the
   `deliver` guide and `CONTEXT.md` say that a stage does not mean the step
   is running.
2. The **Queue Turn**. For a Queue Owner that `deliver status` proves live
   and that recorded its turns, status prints one `Turn:` line per item that
   is neither merged nor parked, in queue order, after the item lines and any
   `Warning:` lines and before any `Park:` line, in one of six fixed forms:
   advancing at its stage; reviewing, with the review slots held; waiting
   for a review slot, with its place among the waiting items and the slots
   held; waiting for its prerequisites, named; waiting for the owner's turn,
   with its place among the items that wait for it; or waiting for the owner
   to read it, for an item whose stage changed after the owner last recorded
   it.
3. The **Turn Record**. The Queue Owner writes the Queue Turn of every item
   it has not finished, its PID, start-time identity and Review Capacity, to
   a file beside its console log whenever a turn changes and every 30
   seconds. Status ignores a record whose PID or identity is not the queue's
   recorded owner, and a missing record prints no turn line. No stage, column
   or schema change is made.
4. The review gets its own capacity. The Queue Owner runs the pre-PR review
   of items in `reviewing` in position order, up to the Review Capacity at
   once, while it keeps advancing one item at a time through every other
   stage. It finds an item a Delivery Retry returned to `reviewing` within 30
   seconds even while it is inside another item's Implement Run. A review
   reads only its own item worktree in a read-only review session; the owner
   never prepares the worktree of an item under review.
5. Order is kept. A queued item starts its Implement Run only when no earlier
   item of the queue is between its Implement Run and its merge, and never
   ahead of an earlier queued item; an earlier item that is parked, merged or
   waiting only for its prerequisites does not hold it. Every review
   outcome parks or advances its item exactly as today, and a review that
   cannot run parks only its item as `delivery-error`.
6. The knob. `delivery.review_concurrency` is a User Config integer, `2` by
   default, at least `1`; a Project Config value is ignored with the warning
   other operator-only keys print. The owner reads it when it starts.
   `deliver start` and `deliver status` print it on the `Limits:` line as
   `review concurrency <n>` after `concurrency 1`.
7. Limits are unchanged. The Queue Token Ceiling, the queue deadline and the
   retry limit still apply to new starts and retries only; the machine-wide
   ceiling on Active Implement Runs and each Run's Verification Capacity do
   not count reviews, because a review is not an Implement Run.
8. The rules in writing. The Roundfix skill's `deliver` reference, the
   `deliver` and configuration user guides, `CONTEXT.md` and ADR-0365 state
   the stage meaning, the Queue Turn lines, the Turn Record, the Review
   Capacity and the start rule.

## User Experience

While the owner implements `0068-bravo`, two items are reviewed, one waits for a review slot, and `0083-hotel`, just retried, has not yet been seen by the owner:

```text
0066-alpha	reviewing	-	/worktrees/0066-alpha
0068-bravo	running	-	/worktrees/0068-bravo
0070-charlie	reviewing	-	/worktrees/0070-charlie
0074-delta	reviewing	-	/worktrees/0074-delta
0080-echo	queued	-	-
0081-golf	queued	-	-
0083-hotel	reviewing	-	/worktrees/0083-hotel
Turn: 0066-alpha reviewing; 2 of 2 review slot(s) held
Turn: 0068-bravo advancing at running
Turn: 0070-charlie reviewing; 2 of 2 review slot(s) held
Turn: 0074-delta waiting for a review slot (1 of 1 waiting); 2 of 2 review slot(s) held
Turn: 0080-echo waiting for prerequisites: 0079-foxtrot
Turn: 0081-golf waiting for the owner's turn (1 of 1 waiting)
Turn: 0083-hotel waiting for the owner to read it
Limits: deadline none, retries per item none, concurrency 1, review concurrency 2, tokens none
Usage: no Runs recorded
```

## Non-Goals / Out of Scope

- Running two Implement Runs of one queue at once, or starting a queued item
  while an earlier item is in review.
- Running the archive, the repository gate, publication, the check wait or
  the merge of two items at once.
- A new Delivery Stage, a Run Database column or any schema change.
- A `deliver start` flag or a Project Config key for the Review Capacity.
- Counting review tokens in the Queue Token Ceiling, or counting reviews
  against the ceiling on Active Implement Runs.
- Fiscus points 2 to 5: the merge-of-main review delta and the review segment
  base (Spec 0356), the numbered-file collision detector (Spec 0357) and large
  authored data files (Spec 0358).
- The Makefile, CI workflows, `go.mod`, `go.sum`, lint or formatter
  configuration, any Baseline module and any `### QA settlement` section.

## Success Metrics

1. Success Metric: with one item's Implement Run held open, an item a Delivery
   Retry returns to `reviewing` is reviewed before that Implement Run ends,
   and the owner never prepares its worktree while the review runs.
2. Success Metric: with a Review Capacity of 2 and three items in
   `reviewing`, at most two reviews run at once, the third starts when one
   ends, and the console log says once that it waits for a review slot.
3. Success Metric: a queued item behind an item in review does not create its
   item branch until that earlier item merges or parks; every existing
   Delivery Engine and `deliver` command test passes unchanged except for the
   `Limits:` line.
4. Success Metric: `deliver status` for a live owner with a Turn Record prints
   the six turn forms in queue order exactly as the Surface Transcript of the
   TechSpec, and prints no turn line for an owner proven gone, a record of
   another owner or no record.
5. Success Metric: `delivery.review_concurrency` defaults to `2`, accepts a
   User Config value of at least `1`, refuses `0` and negative values, is
   ignored with the operator-only warning in Project Config, and appears on
   the `Limits:` line of `deliver start` and `deliver status`.
6. Success Metric: the `deliver` reference of the Roundfix skill and its
   mirror, the `deliver` and configuration guides, `CONTEXT.md` (Review
   Capacity, Queue Turn, Turn Record and the revised Queue Owner) and
   ADR-0365 state the behavior.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Fiscus report, adopted verbatim in `references/` and written before this
  Spec existed: "Num momento, 5 itens estavam em `reviewing` sem nenhum
  processo de review rodando, enquanto um único `implement` (0068) ocupava o
  dono da fila. O `deliver status` mostra `reviewing`, mas não diz que o item
  espera a vez."
- GitHub's documentation on workflow concurrency
  (https://docs.github.com/en/actions/how-tos/write-workflows/choose-when-workflows-run/control-workflow-concurrency,
  read through Exa on 2026-11-05) separates waiting from running in the same
  way: "When a concurrent job or workflow is queued, if another job or
  workflow using the same concurrency group in the repository is in progress,
  the queued job or workflow will be `pending`."

## Unreachable Acceptance

- criterion: Fiscus's own queue shows its turn lines and reviews retried items while another item implements
  reason: a Run cannot reach the Fiscus repository or the MacBook that runs its queue; the QA gate proves the behavior on synthetic queues through the Delivery Engine and the `deliver` command
  satisfied-by: Fiscus's `deliver status` after the release that carries this Spec, reported back through ai-memory

## Validation

The decisions were checked on 2026-11-05 against these sources.

- Reading at `ad77a797`: the Queue Owner runs one loop that, item by item in
  position order, drives an item from its stage to `merged` or `parked`
  before it looks at the next; the pre-PR review runs inline in that loop. A
  Delivery Retry of a `review-findings` park records the stage `reviewing`
  and hands the item to the live owner, which reaches it only when its loop
  comes back to that position. `deliver status` prints the stored stage and
  blocker and a fixed `concurrency 1` on its `Limits:` line.
- A prototype on a scratch copy of `ad77a797` moved the review into a slot
  pool beside the loop, made a queued item wait for an earlier item in
  review, and continued a reviewed item without preparing its worktree
  again; all 363 tests of the Delivery Engine package and the 452 `deliver`
  tests of the command package passed unchanged.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd search "roundfix deliver queue review concurrency" --all --files`,
  which returned the Roundfix mirrors of the configuration and `deliver`
  guides, the Roundfix skill's `deliver` reference and `CONTEXT.md`, and the
  2026-09-24 Exa note on restructuring the queue. The Roundfix entity page
  records that a Run Database schema bump once blocked seven repositories
  until the next release, the cost ADR-0304 avoids and this Spec avoids too.
- Exa: the GitHub page quoted above.

## Glossary

- adds: **Review Capacity**
- adds: **Queue Turn**
- adds: **Turn Record**
- changes: **Queue Owner**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 2 brief for Spec 0355 of 2026-11-05 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-05 ~08:20 BRT), verbatim:

- Fiscus points: question "Relato do Fiscus (fila de 21 Specs há 2 dias).
  Quais pontos viram Spec agora? Os que ficarem de fora viram Backlog Entry."
  — all four options selected; the one this Spec answers: "1. Vez na fila
  (Recommended)" — "O deliver status mostra 'aguardando capacidade' e a
  posição. O review pré-PR, que só lê, ganha uma concorrência própria,
  separada da do implement."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, com
  merge pela fila e conserto quando a entrega parar. Makefile, CI, go.mod e
  lint liberados, assim como mudanças de Baseline e skills que chegam aos
  adotantes. Release assim que cada Spec ficar pronta, com notas aos
  adotantes, aviso ao Fluxus e ao Fiscus."
- Release, the same day: "Quero que faça os releases hoje assim que
  estiverem prontos. E se existe no baseline regra para um release por dia
  por algum motivo deve ser removida." No such rule exists in the repository
  or the Baseline, so this Spec changes nothing for it.

Taken from those answers, the reading and the prototype above, with no
further interview:

- The stage keeps its name and meaning (the next step); the wait is shown on
  a separate `Turn:` line, because renaming a stage or adding one changes
  the Run Database that older item binaries share. See ADR-0365.
- The Review Capacity is User Config only, like the ceiling on Active
  Implement Runs and the Check Wait: it spends the operator's provider
  sessions on the operator's machine, so a repository does not choose it. No
  `deliver start` flag, because recording it per queue needs a schema
  change. See ADR-0365.
- The default is `2`, not `1`: Fiscus measured five items waiting for
  review, and two reviews beside one Implement Run of Task Capacity 2 keep
  the owner within four Agent sessions, fewer than one machine already runs
  under the default ceiling of three Active Implement Runs. See ADR-0365.
- A queued item does not start while an earlier item is in review, so the
  review runs beside other work only when items reach `reviewing` out of
  order, as retries do. See ADR-0365.
- Four implementation Tasks (the documents with ADR-0365 first, then the
  knob and the `Limits:` line, the review slots, and the Turn Record and
  status lines) and a QA gate.

## Open Questions

None.
