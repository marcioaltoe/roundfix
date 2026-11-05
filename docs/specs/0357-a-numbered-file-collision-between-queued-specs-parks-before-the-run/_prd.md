---
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
status: active
created: 2026-11-05
surfaces: [cli, docs]
type: feat
---

# A numbered file collision between queued Specs parks before the Run

When several Specs of one Delivery Queue each add a database migration, every
one of them takes the next migration number of the main it was planned or
started from, so two or more of them claim the same number. Nothing in the
queue notices until a later merge: the second item has already spent its
Implement Run, its pre-PR review, its repository gate and its checks, and then
parks on a pull request conflict that the operator resolves by renumbering by
hand. On 2026-11-05 fiscus reported five queued Specs that all started at
migration `0052`. After this Spec a repository can declare which of its files
are numbered in sequence, and Delivery Revalidation, before the item's first
Run, compares the numbered files the item adds with those already on the
default branch and those added by every earlier item of the queue that has not
merged. A collision parks the item at once in the conflict Park Class, naming
both files, the earlier item that owns the other one and the next free number,
so the operator renumbers one plan before any Run is spent on it.

## Sources

One source, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-fiscus-delivery-queue-report.md`](references/2026-11-05-fiscus-delivery-queue-report.md),
   reported by fiscus on 2026-11-05 at 08:07 BRT (ai-memory message
   `01a1ab64-9074-7520-ad17-6f970973d186`), copied verbatim. This Spec answers
   its point 3: "Migrations drizzle de Specs paralelas colidem sempre" and the
   suggestion "um detector na revalidação que acuse colisão de numeração com
   itens à frente na fila antes do Run". Points 1, 2, 4 and 5 belong to Specs
   0355, 0356 and 0358.

No open Backlog Entry or unresolved Finding exists on `main`. This project's
ai-memory inbox held three pending messages on 2026-11-05: this report, the
fluxus skills request that Spec 0354 answers, and a pantheon notice that asks
for nothing. None shares this context, and none was consumed.

## Prerequisites

None in `requires`. Authored against `ad77a797`, where Specs 0352 and 0353
are active and are being delivered. The planned delivery order is
0354 → 0356 → 0355 → 0357 → 0358, so 0355 and 0356 merge before this Spec runs.
This Spec may share these files with them, keeps its edits narrow and rebases
before its plan merges; its Verification depends on none of them:

- the Delivery Queue engine and its Park Class table (Spec 0355 owns queue
  scheduling and status);
- the `deliver` user guide and the Roundfix skill's `deliver` reference
  (Spec 0355 owns the reference);
- the configuration loader (Spec 0352 owns Git root detection there);
- the Roundfix skill's `SKILL.md` pair and its version record, `CONTEXT.md`
  and the Behavior Surface Record, which every Spec of the cycle re-records.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or stored;
  the detector reads file names and Spec slugs that already exist. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the detector reads the repository through Git and the
  Run Database. Source: `docs/agents/go.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0312 skips Delivery Revalidation
  for an Archived Item Resume, ADR-0312: "runs neither Delivery Revalidation
  nor the Implement Command"; the new check lives inside Delivery
  Revalidation, so it does not run there either. The merge-time recovery of
  pull request conflicts and derived paths, and the prerequisite wait, are
  not changed and carry no obligation here. ADR-0367, written by this Spec,
  records the decision. Source: `docs/agents/domain.md`,
  `docs/agents/autonomous-work.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ batch 2 answers of 2026-11-05 quoted in the authorization record;
  the governed files are the Roundfix skill files the docs Task changes. No
  Makefile, workflow, Go module, lint, formatter or Baseline module change,
  and `.roundfixrc.yml` of this repository is not changed. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Goals

- A repository declares, in Project Config, the paths whose files are
  numbered in sequence; without a declaration nothing changes.
- Before an item's first Run, the queue finds every numbered file the item
  adds whose number is already taken in the same directory on the default
  branch or by an earlier unmerged item of the queue.
- A collision parks the item before any Run, with one line that names each
  colliding pair, the earlier item or the default branch, and the next free
  number, and a Delivery Retry refuses until the plan is renumbered.
- The user guides, the Roundfix skill, `CONTEXT.md` and ADR-0367 state the
  declaration, the rule and the operator's answer.

## User Stories

1. As a repository maintainer whose migrations are numbered, I want to
   declare the migration paths once, so that the queue knows which files take
   a number.
2. As the queue operator, I want an item whose plan takes a migration number
   an earlier queued item already took to park before its Run, so that no Run,
   review or check is spent on a branch that will conflict.
3. As the queue operator, I want the park to name both files, the earlier item
   and the next free number, so that I can renumber the plan without
   investigating.
4. As the queue operator, I want an item whose branch is older than a main
   that took its number since to park the same way, so that a continued item
   branch is held to the same rule.
5. As a maintainer of a repository without numbered files, I want nothing to
   change, so that my queue behaves exactly as before.

## Core Features

1. The declaration. Project Config, with User Config beneath it, gains a list
   of path patterns under the delivery settings, with the same path grammar as
   a Derived Path Declaration. It defaults to empty, which turns the check off;
   a project list replaces the user list, and an explicit empty list disables
   it. An unsafe path or an invalid pattern is refused with a configuration
   error that names the key. The queue reads the list from the default-branch
   commit, as it reads the derived paths.
2. The numbered file. A file a declared pattern matches is a Numbered File when
   its own name starts with one or more digits; the digits are its number,
   read by value, and its directory groups it. A matching file whose name does
   not start with a digit, such as a journal, is not one.
3. What an item adds. An item adds the Numbered Files its branch adds since
   its merge base with the default branch, and the Numbered Files that a Task
   of its Spec not yet completed declares it creates and that its merge base
   does not hold. An item without a branch adds what its plan on the default
   branch declares.
4. The collision. At Delivery Revalidation, a Numbered File the item adds
   collides with a Numbered File of the same directory and number that the
   default branch holds, or that an earlier item of the same queue adds while
   it is not merged, whatever their names. Later items, merged items and
   other queues never collide with it.
5. The park. A collision parks the item before its Run as
   `numbered-file-collision`, in the conflict Park Class, with every colliding
   pair, the earlier item's slug or `default branch` for each, and the next
   free number of each directory: one above the highest number the default
   branch and the earlier unmerged items hold there, written with the item's
   own number of digits. Status and the Pending Question print the next action:
   renumber the plan and any file the branch holds, or deliver the earlier
   item first and merge the default branch, then retry.
6. The retry. A Delivery Retry of an item that has not run repeats the check
   and refuses, leaving the item unchanged, while a collision remains.
7. The rules in writing. The configuration and `deliver` user guides, the
   Roundfix skill's `deliver` reference, `CONTEXT.md` and ADR-0367 state the
   declaration, what an item adds, the collision, the park, its next action and
   what the check does not cover.

## User Experience

A repository declares its migrations:

```yaml
delivery:
  numbered_files: [drizzle/*.sql, drizzle/meta/*_snapshot.json]
```

Item `0070-notas` is queued behind `0068-indices`, which parked after its Run
with `drizzle/0052_indices.sql` on its branch. When `0070-notas` starts, its
plan declares `drizzle/0052_notas.sql`, so `roundfix deliver status` shows:

```text
0070-notas	parked	numbered-file-collision: drizzle/0052_notas.sql collides with drizzle/0052_indices.sql (0068-indices); next free number: drizzle/0053	/worktrees/0070-notas
Park: 0070-notas conflict: in /worktrees/0070-notas, renumber the item's files the blocker names from its next free number, in the Task declarations on the item branch and in any file the branch holds, or deliver the earlier item it names first and merge the default branch into the item branch; then run roundfix deliver retry 0070-notas
```

No Run started for `0070-notas`. After the operator commits the renumbered
plan on the item branch, `roundfix deliver retry 0070-notas` repeats the check
and returns the item to the queue.

## Non-Goals / Out of Scope

- Renumbering by the queue. The maintainer chose the detector; renaming files,
  rewriting a migration journal or snapshot, or running a generator stays the
  operator's or the Agent's work, so a renumbering step is left as a follow-up.
- Waiting for the earlier item, as a prerequisite does: a plan that names a
  taken number must change whether or not the earlier item merges first.
- A number the Run will generate without its plan declaring it. Before the
  Run the queue can only see declared and committed files; such a number is
  found once a branch holds it, for every later item, and otherwise at merge
  as today.
- A Delivery Retry of an item whose Run already started; the check runs only
  where Delivery Revalidation runs, so such an item still meets the collision
  at merge.
- The read-only Delivery Plan, pull request conflict recovery, derived path
  regeneration and the Archived Item Resume, which keep their behavior.
- Folders numbered in sequence, as drizzle-kit v1 writes them; the declaration
  covers files.
- A default list of patterns, and any change to this repository's own Project
  Config.
- Queue scheduling and status lines other than the new park (Spec 0355), the
  pre-PR review (Specs 0356 and 0358), the Makefile, the CI workflows,
  `go.mod`, `go.sum`, lint or formatter configuration, any Baseline module and
  any `### QA settlement` section.

## Success Metrics

1. Success Metric: a configuration that declares `numbered_files` patterns
   reads them from Project Config over User Config, an explicit empty list
   clears them, and an unsafe path or an invalid pattern is refused naming
   `delivery.numbered_files`; without the key every existing Delivery
   Revalidation test passes unchanged.
2. Success Metric: in temporary repositories, an item whose plan or branch
   adds a Numbered File whose directory and number an earlier unmerged item's
   branch or plan, or the default branch, already holds is reported with the
   exact pair, owner and next free number; a later item, a merged item, a
   completed Task, a name without a leading digit and a renumbered plan report
   nothing.
3. Success Metric: a queue whose item collides parks it as
   `numbered-file-collision: <detail>` before any Run, in the conflict Park
   Class with the exact next action, and a retry of that item refuses while the
   collision remains and proceeds once it is gone.
4. Success Metric: `roundfix deliver status` prints the collision park and its
   next action for a parked item.
5. Success Metric: the configuration and `deliver` user guides, both copies of
   the Roundfix `deliver` reference, `CONTEXT.md` and ADR-0367 state the key,
   the rule, the park and its next action, and the three Park Class tables list
   `numbered-file-collision` in the conflict row.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The fiscus report, quoted in the adopted reference and written before this
  Spec existed: "Na fila, cinco Specs (0068, 0070, 0074, 0079 e outra)
  começaram todas na 0052. Cada merge obriga as demais a renumerar à mão".
- The drizzle-orm discussion "Best way to deal with migration merge
  conflicts?" (https://github.com/drizzle-team/drizzle-orm/discussions/1104,
  read through Exa on 2026-11-05): "both your branch (`HEAD`) and `dev`
  created migrations numbered 0007", with `both added:
  db/migrations/meta/0007_snapshot.json` in the conflict it shows.
- The drizzle-orm discussion "I built a GitHub Action that catches migration
  collisions in PRs before they hit main"
  (https://github.com/drizzle-team/drizzle-orm/discussions/5522, read through
  Exa on 2026-11-05): "two developers branch off main, both run `drizzle-kit
  generate`, and both create a migration with the same sequence number. The
  second PR to merge silently breaks the migration history".

## Validation

The decisions were checked on 2026-11-05 against these sources.

- Reading at `ad77a797`: Delivery Revalidation runs in the item worktree right
  after it is created or continued, before the first Run, and on a Delivery
  Retry of an item that has not run; it parks only on strict findings, and
  records production-premise overlaps as a warning. No code reads a number in
  a file name, and nothing compares the files of one queue item with another
  item's branch or plan. A pull request conflict is found only when the owner
  merges the default branch for a conflicting Pull Request, after the Run,
  review, gate and checks; a conflict outside the derived paths parks as
  `pull-request-conflict`.
- Reading at `ad77a797`: the item branch, once recorded, is reused by a
  retry, so a park after the worktree exists is answered on that branch; a
  prerequisite park happens before the worktree exists and is released by the
  owner.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd search "drizzle migration numbering collision" --all --files` (one
  fluxus techspec mirror unrelated to numbering) and
  `qmd search "roundfix delivery queue revalidation park" --all --files`,
  which returned the Roundfix mirrors of ADR-0312, the `deliver` guide, the
  archived Spec 0180 that introduced Delivery Revalidation, the `deliver`
  reference and `CONTEXT.md`. None records a numbered file check; Spec 0180
  confirms revalidation was meant as the place to re-check a later Spec against
  the main it starts from.
- Exa: the two drizzle-orm discussions quoted above, and the drizzle v1 notes
  (https://notes.drizzle.run/articles/migrations-v1/), which move migrations
  into timestamped folders; that is why folders are out of scope and why the
  key is generic, not drizzle-specific.

## Glossary

- adds: **Numbered File**
- adds: **Numbered File Collision**
- changes: **Delivery Revalidation**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ batch 2 brief for Spec 0357 of 2026-11-05 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-05 ~08:20 BRT), verbatim:

- Fiscus points: question "Relato do Fiscus (fila de 21 Specs há 2 dias).
  Quais pontos viram Spec agora? Os que ficarem de fora viram Backlog Entry."
  → all four options selected, of which this Spec answers "3+5. Migrations e
  catálogos" — "Um detector de número de migration repetido entre Specs da
  fila, rodado antes da Run. Uma categoria nova para arquivos grandes de dados
  escritos pelo projeto, revisados à parte em vez de omitidos. São desenhos
  maiores." (the second sentence is Spec 0358).
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, com
  merge pela fila e conserto quando a entrega parar. Makefile, CI, go.mod e
  lint liberados, assim como mudanças de Baseline e skills que chegam aos
  adotantes. Release assim que cada Spec ficar pronta, com notas aos
  adotantes, aviso ao Fluxus e ao Fiscus."
- Release, the same day: "Quero que faça os releases hoje assim que estiverem
  prontos. E se existe no baseline regra para um release por dia por algum
  motivo deve ser removida." No such rule exists in the repository or the
  Baseline, so this Spec changes nothing for it.

Taken from those answers, the reading above and the sources, with no further
interview:

- The detector is generic: any directory whose files start with a sequence
  number, declared by path patterns. See ADR-0367.
- The declaration is off by default, because file names that start with digits
  are common outside sequences (dated Backlog Entries and Findings, `01-intro`
  pages), and a default would park unrelated items. See ADR-0367.
- The check sits in Delivery Revalidation and parks; it neither waits nor
  renumbers. See ADR-0367.
- Precedence follows queue order: the later item renumbers.
- Three implementation Tasks (the declaration, the rule and park in the
  engine, the Git reading at revalidation) plus a documentation Task with
  ADR-0367, and a QA gate.

## Open Questions

None. Whether fiscus plans declare their migration files under `creates:` is
not verified from here (another repository); when a plan does not, the check
still holds every later item to the file once an earlier branch holds it.
