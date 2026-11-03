---
spec: 0343-the-follow-up-names-what-waits-and-where-each-item-went
status: active
created: 2026-11-03
surfaces: [backend, cli, docs]
type: feat
---

# The follow-up names what waits on someone else and where each item went

A `/roundfix --follow-up` run lists the open Backlog Entries and unresolved
Findings, and its agent decides which become Specs, but nothing records an
item that waits on something outside the session: on 2026-11-02 Fluxus got
three Findings and a Backlog Entry that each waited on a vendor, a supplier or
the maintainer, the agent authored none of them, the report said why only in
free text, and the `repository-inbox` source listed the folder marker
`docs/_inbox/.gitkeep` as a note. After this Spec a Backlog Entry or a Finding
can carry a **Waiting Declaration**, two frontmatter fields `waiting_on` and
`since`, which makes it a **Waiting Item** that `roundfix follow-up` holds back
by default, counts on its source line and lists with `--all`; the follow-up
mode writes that declaration when it leaves an item waiting and its report
names where every listed Finding and Backlog Entry went; and the inbox source
ignores hidden files and folders.

## Sources

One source, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-02-the-follow-up-names-what-waits-and-where-each-item-went.md`](references/2026-11-02-the-follow-up-names-what-waits-and-where-each-item-went.md),
   `type: feat`, reported by fluxus, ai-memory message `01a19bf3` of
   2026-11-02, with its three requests and the evidence line of the Fluxus
   run.

No open Backlog Entry or unresolved Finding exists on `main`, so no other
source shares this context. Specs 0342 and 0344 are authored in parallel for
other reports.

## Prerequisites

None in `requires`. Authored against `4e61bf6d` (v0.129.1), where no other
Spec is active. Specs 0342 and 0344 are authored in parallel; 0342 delivers
before this Spec and 0344 after it. `_techspec.md` → Build Order names the
files they may share.

## Project Constraints

- Identifier strategy: not applicable — no identifier or slug is minted; two
  optional frontmatter keys are read, and the follow-up JSON gains optional
  keys inside the unchanged `roundfix/follow-up/v1` schema. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the judge's transport, keys and recipients stay as they
  are, and fewer items reach it. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0319 owns the Follow-up Priority,
  ADR-0319: "The order is computed in code, never asked."; this Spec adds a
  class outside its four and asks the judge about fewer items. ADR-0338 makes
  the follow-up mode claim the handoff and leave no handed-in message pending,
  ADR-0338: "No message the mode handed in stays pending when the turn ends";
  this Spec extends the same accounting to every listed Finding and Backlog
  Entry. ADR-0330 keeps the command away
  from ai-memory, ADR-0330: "Roundfix's own code never talks to ai-memory";
  nothing here changes that. ADR-0353, written by this Spec, records the
  waiting class, its default exclusion and where the declaration is
  validated. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AW answers of 2026-11-03 quoted in the authorization record, which keep the
  cycle AV rules and the cycle AJ grant for the Governed Paths each Spec
  declares; the governed files are the Roundfix skill's `SKILL.md` pair and
  its canonical `follow-up`, `deliver` and `modes` references. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0343-the-follow-up-names-what-waits-and-where-each-item-went/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/follow-up.md`, `.agents/skills/roundfix/references/modes.md`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint, formatter or Baseline module change.

## Goals

- A folder marker or any hidden file under `docs/_inbox/` is not pending work
  for the follow-up or the delivery plan.
- A Finding or Backlog Entry that waits on someone else says so in its own
  frontmatter, and the follow-up shows it as waiting, with the reason and the
  date, instead of leaving the decision in an agent's free text.
- `roundfix follow-up` leaves Waiting Items out of its list by default without
  hiding them: each source says how many it held back, and `--all` lists them.
- A malformed declaration is reported, never silently accepted or silently
  dropped.
- The follow-up mode's report names the destination of every listed Finding
  and Backlog Entry, so a missing destination is visible as a gap.
- A repository that declares no wait gets the same report bytes as before.

## Core Features

1. Hidden inbox entries. The `repository-inbox` source of the follow-up and
   the inbox rows of `roundfix deliver plan` skip every file and folder under
   `docs/_inbox/` whose name starts with a dot, such as `.gitkeep`, and
   everything inside a hidden folder.
2. Waiting Declaration. An open Backlog Entry or a pending or partial Finding
   may declare `waiting_on` (who or what it waits on) and `since` (the
   `YYYY-MM-DD` date the wait began). With both valid, the item is a Waiting
   Item in a class of its own, `waiting`.
3. Held back by default. `roundfix follow-up` lists no Waiting Item; the
   backlog and findings sources keep counting every item they read and say
   how many of them wait and were held back. No Waiting Item is sent to the
   judge.
4. `--all`. `roundfix follow-up --all` also lists the Waiting Items, after
   every ranked item, oldest `since` first, each with its `waiting_on` and
   `since`; they are still not sent to the judge. The status view refuses
   `--all` and is otherwise unchanged.
5. Validation. A declaration with an empty or non-text `waiting_on`, a missing
   `since`, a `since` that is not a calendar date, or a `since` without
   `waiting_on` makes no Waiting Item: the item stays in the list in its
   ordinary class with the reason, and no other source fails.
6. The follow-up mode. `/roundfix --follow-up --all` runs the command with
   `--all`; an item the plan leaves waiting gets its declaration written;
   and the report names, for each listed `finding` and `backlog` item, the
   Spec it became, the Spec it was grouped into, or what it waits on and
   since when, and how many Waiting Items the command held back.
7. The rules in writing. The follow-up and deliver guides, the Roundfix
   skill's follow-up, deliver and modes references, `CONTEXT.md` and ADR-0353
   state them.

## Non-Goals / Out of Scope

- Adding `waiting_on` and `since` to the Baseline's Backlog and Findings
  Operational Contracts: ADR-0353 defers it, and this Spec's planning change
  records it as an open Backlog Entry.
- A `roundfix spec check` finding for a malformed declaration, or any Spec
  Check change.
- A waiting count or any other change in the `--status` view, which reads no
  Backlog Entry or Finding.
- A Waiting Declaration on an ai-memory message, a handed-in handoff entry, a
  note or an active Spec.
- Expiring a wait by date, reminding anyone, or reading `since` against the
  clock.
- Changing the four classes, the four questions, the judge's thresholds, its
  keys or its monthly ceiling.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, lint or formatter
  configuration, any Baseline module and any `### QA settlement` section.

## Success Metrics

1. Success Metric: with `docs/_inbox/.gitkeep`, a file inside
   `docs/_inbox/.obsidian/` and a visible note, the follow-up's
   `repository-inbox` source counts 1 and lists only the visible note, and
   `roundfix deliver plan` prints one inbox row, the visible note's.
2. Success Metric: a pending Finding declaring `waiting_on` and a valid
   `since` is not listed by `roundfix follow-up`; the `findings` source counts
   it, its JSON carries `waiting: 1` and its text line ends with
   `1 waiting, held back without --all`; with a judge key set the judge
   receives requests for the other items only.
3. Success Metric: `roundfix follow-up --all` lists that Finding after every
   ranked item with `priority_name` `waiting`, `basis` `waiting`, its
   `waiting_on` and `since`, oldest `since` first, and the judge still
   receives no request for it.
4. Success Metric: each malformed declaration (blank `waiting_on`, a list as
   `waiting_on`, no `since`, an impossible date as `since`, `since` alone)
   lists the item in its ordinary class with the documented
   `waiting_error`, and every other source keeps its state and count.
5. Success Metric: for a repository with no declaration, the text and JSON of
   `roundfix follow-up` equal, byte for byte, those of the delivery base's
   binary on the same repository; only the help text adds `--all`.
6. Success Metric: `roundfix follow-up --status --all` exits 2 with the
   refusal on standard error and nothing on standard output, and
   `roundfix follow-up --status` prints the same bytes as before.
7. Success Metric: the four Surface Transcripts of `_techspec.md` reproduce
   through the built binary.
8. Success Metric: the follow-up and deliver guides, both copies of the
   follow-up, deliver and modes references and `CONTEXT.md` describe the
   hidden-entry rule, the Waiting Declaration, `--all` and the report's
   destinations, and each mirror equals its canonical copy.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Fluxus report quoted in the adopted reference, written before this Spec
  existed: on roundfix 0.129.0 `roundfix follow-up --format json` listed a
  `.gitkeep` note and three Findings and a Backlog Entry that waited on others,
  with no recorded reason.
- The Getting Things Done guide "Getting Your Inbox to Zero"
  (https://gettingthingsdone.com/wp-content/uploads/2014/10/GettingYourInboxToZero.pdf,
  read through Exa on 2026-11-03): "If delegating, track the waiting for
  reminder (if you need to track this getting completed) in your @Waiting For
  folder in email or Waiting For list in your list manager." and "If
  deferring, track the action reminder in @Action folder in email, Calendar,
  or a Next Actions list in your list manager.", which keep what waits on
  others on a list apart from the next actions.
- The page "Getting Things Done: Organizing “Waiting For” & Own Action
  Reminders" (https://www.wellthatmakessense.com/organizing-the-chaos/, read
  through Exa on 2026-11-03), quoting the book: "When the next action on
  something is up to someone else, you don’t need an action reminder. Just a
  trigger about what you’re waiting for and from whom." and "It’s important
  for this category in particular to include the date that each item is
  requested for each entry", the two facts `waiting_on` and `since` record.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files and
temporary repositories.

## Validation

The decisions were checked on 2026-11-03 against these sources.

- Reading at `4e61bf6d`: `deliveryPlanIntent` walks `docs/_inbox/` with
  `filepath.WalkDir` and records every regular file, and it is the only
  reader of that folder, shared by `roundfix deliver plan` and the follow-up's
  `repository-inbox` source; `readFollowUpMarkdown` decodes the frontmatter
  into a struct that ignores unknown keys; `followUpOrder` places every item
  in one of four classes; the judge receives every listed item.
- Characterization at `4e61bf6d` with that commit's `bin/roundfix`, in a
  temporary repository holding `docs/_inbox/.gitkeep`,
  `docs/_inbox/.obsidian/app.json`, a pending Finding with `waiting_on` and
  `since`, and an open Backlog Entry: `roundfix follow-up` listed four items,
  `other note - docs/_inbox/.gitkeep: .gitkeep` and
  `other note - docs/_inbox/.obsidian/app.json: app.json` among them, the
  Finding as `health` with no trace of its wait, and `repository-inbox:
  checked 2`; `roundfix deliver plan` printed both inbox rows.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "follow-up waiting on external dependency backlog finding blocked status field" --all --files --min-score 0.3`,
  which returned only Roundfix project mirrors (the deliver and follow-up
  guides, the modes reference, ADR-0319, ADR-0338 and one older decision
  record on outside evidence), and
  `qmd query "GTD lista aguardando terceiros waiting for revisão semanal" --all --files --min-score 0.3`,
  whose digest `wiki/sources/digest-gtd-esclarecer-triagem-agentes-2026-10-14.md`
  records the maintainer's own practice of a separate "Aguardando" list,
  reviewed weekly with the other lists, and whose entity page
  `wiki/entities/todoist.md` records a next-actions filter that excludes the
  `aguardando` label. Neither decides a Roundfix field.
- Exa: the two Getting Things Done sources quoted above.

## Glossary

- adds: **Waiting Declaration** — the `waiting_on` and `since` pair of a Backlog Entry or a Finding
- adds: **Waiting Item** — a Follow-up Item whose Waiting Declaration is valid
- changes: **Follow-up Item** — a hidden file or folder in `docs/_inbox/` is not one
- changes: **Follow-up Priority** — the fifth class `waiting`, outside the four
- changes: **Follow-up Command** — Waiting Items are held back unless `--all`

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AW brief for Spec 0343 of 2026-11-03 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-03), verbatim:

- 0343: the maintainer typed (Other): "Os tres e por padrão sem waiting quando
  usado só --follow-up, mas preciso de uma flag ou parametro --all (Algo
  assim) para incluir tudo, incluindo os waiting". The recommended option it
  answered was "Os três, waiting_on + since — Campos no frontmatter:
  waiting_on: <quem/o quê> e since: <AAAA-MM-DD>. O item sai como classe
  'waiting', com motivo e data. Eu respondo ao Fluxus com o nome dos campos."
- Rules (all three): "Mesmas do AV (Recommended)" — option text: "Push e PR
  dos planos, entrega e merge pela Delivery Queue, consertar e seguir em falha
  técnica (pular só se o conserto contradisser o PRD). Makefile,
  .github/workflows, go.mod/go.sum e config de lint/formatter liberados quando
  o desenho pedir."
- Release (all three): "Uma v0.130.0 no fim (Recommended)" — option text:
  "Minor único depois das três Specs, porque a 0342 cria a chave
  pre_pr_review.omit_paths e a 0343 cria o --all. Notas aos 8 adotantes,
  resposta ao Fluxus com os campos waiting_on/since, aviso ao Vortex da chave
  nova e aviso ao Jarvis de que a versão saiu."

Taken from those answers, ADR-0319, ADR-0338 and the reading above, with no
further interview:

- All three requests are in scope, with the field names Fluxus proposed,
  `waiting_on` and `since`.
- The flag is `--all`, as the maintainer named it; `roundfix follow-up` has
  no other flag of that name, and `--all` of `roundfix deliver lane` is
  another command. The Skill Mode accepts `all` or `--all` after
  `follow-up`.
- The declaration is validated by the Follow-up Command, its only reader:
  a malformed one lists the item in its ordinary class with the reason.
  `roundfix spec check` and the Baseline's contracts do not change
  (ADR-0353); a Backlog Entry records the Baseline contract change for a
  later cycle.
- A Waiting Item is never sent to the judge, with or without `--all`.
- `--status` shows no waiting count and refuses `--all`.
- Hidden means a name that starts with a dot; the rule lives in the one
  inbox walk, so the follow-up and `roundfix deliver plan` change together.
- Three implementation Tasks in series (hidden inbox entries; the waiting
  class and `--all`; the follow-up mode) and a QA gate.

## Open Questions

None.
