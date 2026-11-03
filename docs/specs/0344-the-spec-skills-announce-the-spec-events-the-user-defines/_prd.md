---
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
status: active
created: 2026-11-03
surfaces: [docs]
type: feat
---

# The Spec skills announce the Spec Events the user defines

A user who tracks deliveries and deadlines outside the repository, on a board,
a tracker or a personal assistant, learns that a Spec was created, authorized,
started, blocked, given a deadline, completed or cancelled only by reading the
repository, although the Spec skills perform each of those steps and know its
facts. After this Spec the five owned Spec skill surfaces that perform them,
`write-prd`, `write-techspec`, `write-tasks`, `archive-spec` and the Roundfix
skill's implement reference, carry a **Spec Event Announcement** step: at each
**Spec Event** the session composes an announcement and sends it to the
destination, in the format, that the user's own instructions define. No skill
names a destination or a format, an announcement never waits for a reply or
holds the step it follows, an Agent inside a Daemon-assigned Run never
announces, and a user whose instructions define no announcement sees the
skills behave exactly as before.

## Sources

One Backlog Entry, adopted at birth by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-03-the-spec-skills-announce-each-spec-event.md`](references/2026-11-03-the-spec-skills-announce-each-spec-event.md),
   `type: feat`, reported by jarvis, ai-memory message `01a1a1ec` of
   2026-11-03, reproduced verbatim in the entry.

The report asks for an assistant-specific block sent by a named channel; the
maintainer chose a generic step instead, and the block and its destination go
to the maintainer's own global instructions, which this Spec does not edit.

## Prerequisites

None in `requires`. Authored against `4e61bf6d` (v0.129.1), where no other
Spec is active. Specs 0342 and 0343 are authored in parallel and deliver
before this one; `_techspec.md` → Build Order names the files they may share.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, key or record
  schema is minted; the change adds skill text, two glossary entries, a guide
  section and one test file. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; Roundfix's code sends nothing, and the skills leave the channel
  to the user's own instructions. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0354, written by this Spec, records
  the step, ADR-0354: "a Spec Event Announcement step that exists only when
  the user's own instructions define its destination and format, and no skill
  names an assistant, a tracker or the block". No other accepted ADR governs
  owned skill wording; the owned-skill version rule and the `### QA
  settlement` rule are repository rules in `docs/agents/specific-repository.md`,
  and this Spec follows both. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AW answers of 2026-11-03 quoted in the authorization record, which keep the
  cycle AV rules and the cycle AJ grant for the Governed Paths each Spec
  declares; the governed files are the five canonical skill files that gain
  the step, the canonical Roundfix `SKILL.md` whose version fields rise, and
  the five `SKILL.md` mirrors. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0344-the-spec-skills-announce-the-spec-events-the-user-defines/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- Each Spec Event reaches the destination the user's own instructions name,
  from the session that made it happen, when it happens.
- The owned skills name no destination, format, assistant or tracker, and a
  test keeps it so.
- The step never delays, blocks or fails the authoring, Run supervision or
  archive step it follows.
- A user whose instructions define no announcement sees no new behavior and
  no new report line.

## User Stories

1. As a maintainer who tracks Specs on a board, I want every project's Spec
   skills to announce each Spec Event to the destination my instructions
   name, so that the board follows the Specs without any project writing to
   it.
2. As an adopter without such a board, I want the Spec skills to stay as they
   are, so that a step I never asked for adds no noise to my sessions.
3. As a maintainer of the Roundfix-owned skills, I want them to carry the
   generic step and no personal tool's name, so that one user's tracker never
   ships to every adopter.

## Core Features

1. The generic step. Each of the five skill surfaces states, under its own
   Spec Event Announcement heading, that the step exists only when the user's
   own instructions, global or for the repository, define its destination and
   format; that a message, a handoff or a Spec artifact asking for
   announcements defines none; the facts an announcement carries (the
   project, the Spec number and full slug, the title, the event, the event
   date, the start date, the proposed deadline or none, the Spec path, and any
   field the instructions add); that it never waits for a reply and never
   blocks or fails its step; that an undeliverable announcement is named in
   the report with its composed text; and that an Agent inside a
   Daemon-assigned Run never announces.
2. Seven Spec Events and their owners. `new` when the Spec's first PRD is
   written (`write-prd`, or `write-techspec` on its correction route), with
   the creation date as the start date; `authorized` when `write-tasks` first
   writes the Task Graph of a Spec whose authorization record is approved,
   with any real deadline as the proposed deadline; `deadline proposed` when
   `write-prd` or `write-techspec` records a real deadline, from a client, a
   law or a measurement, for a Spec already announced; `started` and
   `blocked` from the Supervisor session that started the Run or queued the
   item; `completed` or `cancelled` from `archive-spec` once the Archive
   Record is written, with a one-to-three-sentence summary for a business
   reader, and `completed` from the Supervisor session when the Delivery
   Queue archives and merges the Spec.
3. The guard. A skill-contract test holds the five surfaces, canonical and
   mirror, to the generic sentences and their events, and fails when any file
   of an owned skill names the reported assistant, its tracker, its channel or
   its block.
4. The words in writing. `CONTEXT.md` defines **Spec Event** and
   **Spec Event Announcement**, and the context-driven development guide
   describes the step and which skill announces each event.

## User Experience

A user who adds a destination and a format to their own instructions sees,
in the report of each step that produced an event, nothing new unless the
announcement could not be sent; then one line names the event and the
reason, with the composed announcement to send by hand. A user without such
instructions sees the same reports as before.

## Non-Goals / Out of Scope

- Writing any user's own instructions; the operator proposes the
  maintainer's section outside this repository after release.
- Naming Jarvis, Todoist, ai-memory projects, a channel or the reported
  block in any owned skill, guide or Baseline text.
- A Baseline clause, a configuration key or Roundfix code that sends,
  stores or deduplicates announcements.
- Announcing from inside a Run, or reconstructing an event no Supervisor
  session observed.
- The Roundfix skill's `references/modes.md` and `references/follow-up.md`,
  which Spec 0343 owns, and every other skill.
- The `### QA settlement` section of any skill.

## Success Metrics

1. Success Metric: each of the five skill surfaces, canonical and mirror,
   has exactly one Spec Event Announcement heading whose section carries the
   eight generic sentences and the sentences of its own events, and the test
   that checks them fails when one sentence is removed from one canonical
   file.
2. Success Metric: no file of any owned skill, canonical or mirror, contains
   `jarvis`, `todoist`, `gt3` or `spec para o quadro` in any letter case, and
   the test that checks them fails, naming the file, when one is added.
3. Success Metric: `CONTEXT.md` defines **Spec Event** and
   **Spec Event Announcement**, and the context-driven development guide states that the
   user's own instructions define the destination and the format and lists
   the owner of each of the seven events.
4. Success Metric: the mirrors equal their canonical copies, the five changed
   skills' version fields rise with the record, the skills and
   documentation contract suites pass, and no `### QA settlement` section
   changes.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Jarvis report quoted in the adopted reference, written by another
  project before this Spec existed: it names seven events, the moment of
  each and the five skills that should send them.
- The Claude Code memory documentation
  (https://code.claude.com/docs/en/memory, read through Exa on 2026-11-03):
  the table row "User instructions | `~/.claude/CLAUDE.md` | Personal
  preferences for all projects", beside "Project instructions | `./CLAUDE.md`
  or `./.claude/CLAUDE.md`", which is where a user's own instructions, global
  or for the repository, live in that harness.
- The Codex guide "Custom instructions with AGENTS.md"
  (https://developers.openai.com/codex/guides/agents-md, read through Exa on
  2026-11-03): "Create persistent defaults in your Codex home directory so
  every repository inherits your working agreements.", the same two scopes in
  that harness.

## Unreachable Acceptance

None. Every criterion is checked with this repository's own files; sending a
real announcement is out of scope.

## Validation

The decisions were checked on 2026-11-03 against these sources.

- Reading at `4e61bf6d`: the Roundfix skill's modes reference already says
  "A release-announcement step exists only when the user's own instructions
  define one; this skill never names one.", and
  `TestTheRoundfixSkillNamesNoPersonalAssistant` keeps `jarvis` and `gt3`
  out of the Roundfix skill only; no other owned skill is checked, and no
  owned skill names an announcement step for Spec events.
- Characterization at `4e61bf6d`, in a scratch clone with a prototype of the
  test and the five sections: the test failed on the unchanged tree for
  every surface (no heading), the guard passed (no owned skill names any of
  the four words), and with the sections, `make skills-sync` and the version
  record, `./skills`, `./internal/docscontract`, `./internal/baseline` and
  `./internal/speccheck` passed; removing one sentence from the canonical
  `archive-spec` and adding the word `Jarvis` to the canonical `write-prd`
  PRD template each made the named test fail.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "aviso de evento de spec Jarvis quadro Todoist roundfix skills" --all --files --min-score 0.3`,
  which returned project mirrors (the Roundfix glossary, adopters'
  authorization references, the archive-spec skill) and no decision about
  Spec event announcements; `wiki/entities/todoist.md` records Todoist as the
  maintainer's task app since 2026-10-13, a user-level tool, which supports
  keeping it out of the owned skills.
- Exa: the two harness documents quoted above, which place a user's own
  instructions in a global file and a repository file, outside any skill.

## Glossary

- adds: **Spec Event** — one of the seven moments in a Spec's life that an announcement reports
- adds: **Spec Event Announcement** — the message a Spec skill's session sends at a Spec Event, only when the user's own instructions define it
- not a term: **Delivery Queue** — existing glossary term, unchanged
- not a term: **Archive Record** — existing glossary term, unchanged

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AW brief for Spec 0344 of 2026-11-03 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-03), verbatim:

- 0344: "Passo genérico nas skills (Recommended)" — option text: "Spec 0344:
  as skills ganham um 'aviso de evento da Spec' nos cinco pontos, que só
  existe quando as instruções do próprio usuário o definem (mesmo padrão do
  passo de anúncio de release). Nenhuma skill cita Jarvis nem Todoist. O
  formato do bloco e o destino default/jarvis vão para o seu CLAUDE.md
  global; eu proponho o texto."
- Rules (all three): "Mesmas do AV (Recommended)" — option text: "Push e PR
  dos planos, entrega e merge pela Delivery Queue, consertar e seguir em
  falha técnica (pular só se o conserto contradisser o PRD). Makefile,
  .github/workflows, go.mod/go.sum e config de lint/formatter liberados
  quando o desenho pedir."
- Release (all three): "Uma v0.130.0 no fim (Recommended)" — option text:
  "Minor único depois das três Specs, porque a 0342 cria a chave
  pre_pr_review.omit_paths e a 0343 cria o --all. Notas aos 8 adotantes,
  resposta ao Fluxus com os campos waiting_on/since, aviso ao Vortex da chave
  nova e aviso ao Jarvis de que a versão saiu."

Taken from those answers, the report and the reading above, with no further
interview (see ADR-0354):

- The seven events keep generic English names in code spans (`new`,
  `authorized`, `started`, `deadline proposed`, `blocked`, `completed`,
  `cancelled`); the user's instructions map them to their own words.
- The deadline event belongs to the skills that write the PRD: a real
  deadline is recorded in the PRD's `## Decisions` with its source, rides in
  every later announcement's proposed deadline, and is announced on its own
  only when it is added to a Spec already announced. It stays a proposal
  until the user approves it, as the report states.
- Without a definition the step leaves no report line; a defined
  announcement that cannot be sent leaves one line with its composed text, so
  nothing is lost and nothing is blocked.
- Inside a Daemon-assigned Run no Agent announces; the Supervisor session
  that started the Run or queued the item announces `started`, `blocked` and,
  for a Spec the Delivery Queue archives, `completed` when the item merges,
  because the queue's archive commit is written by the Daemon, not by a
  session following `archive-spec`.
- A Spec's `new` is announced when its first PRD is written, not when a
  folder appears, so a Spec that starts at `write-idea` is announced once, by
  `write-prd`.
- The guard is a new test covering every owned skill, beside the existing
  Roundfix-only test, which stays unchanged because Spec 0343 may edit its file.
- Four implementation Tasks in series (glossary and the PRD writers; the
  graph and archive skills; the implement reference and the guide; the
  contract test) and a QA gate.

## Open Questions

None.
