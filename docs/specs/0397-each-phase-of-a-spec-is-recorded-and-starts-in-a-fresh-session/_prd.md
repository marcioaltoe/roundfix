---
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
status: active
created: 2026-11-16
surfaces: [backend, cli, docs]
type: feat
---

# Each phase of a Spec is recorded and starts in a fresh session

The usage limit of an Agent harness is spent in Supervisor sessions that nothing describes: no record says which Spec or phase a session worked on, and a session that authors, queues and follows a Spec re-reads its whole context on every call. After this Spec, Roundfix appends to one documented record which Spec and phase each Supervisor session works on and which Work Item each prompt of a Run served, so that a usage tool can attribute the limit by phase. The skills and the autonomous-work guide tell the Supervisor to end its session at each phase boundary with a Session Handoff, and Roundfix says so, refusing nothing, when a session goes on. This repository starts running its own Authoring Phases through the generated subagents, with no header in the skills it publishes, and a fix or a refactor is authored in the session unless the user asks for the declared author. Roundfix reads no harness transcript, prices nothing and reports no limit percentage.

## Sources

Adopted by this Spec's planning change (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-14-jarvis-long-orchestration-consumes-the-weekly-claude-limit.md`](references/2026-11-14-jarvis-long-orchestration-consumes-the-weekly-claude-limit.md), ai-memory message `01a1da5d-ff74-7020-84f3-e620d48aa249`, reported by jarvis on 2026-11-14. The sender asked for a reply with the decision and the operator replied on 2026-11-14. It declared no blocked work and asked for no completion notice.
2. [`references/2026-11-15-maintainer-answers.md`](references/2026-11-15-maintainer-answers.md), the maintainer's answers, recorded by the operator.
3. [`references/2026-11-15-request-read-by-the-idea-stage-judge.md`](references/2026-11-15-request-read-by-the-idea-stage-judge.md), the request the judge read.
4. [`references/2026-11-15-usage-by-context-and-session-identity-measurement.md`](references/2026-11-15-usage-by-context-and-session-identity-measurement.md) and [`references/2026-11-16-delivery-speed-and-plan-quality-baseline.md`](references/2026-11-16-delivery-speed-and-plan-quality-baseline.md), recorded by the author.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from. At the authoring base this project's ai-memory inbox holds only the adopted message, `docs/findings/` is empty, and the three open Backlog Entries share no component with this Spec: a Conexus report that waits for reproduction steps, the status queue table the maintainer ordered after this Spec, and the start cost of a session, which the Idea recorded.

## Prerequisites

None in `requires`. Authored against `2caa1d17` (v0.145.1); no other Spec is planned and the Delivery Queue is idle. Spec 0394, the declaration for OpenCode and Codex, is not authored yet and will extend the same installation code this Spec changes.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the Phase Record stores identifiers other components made (a Spec slug, a Run identifier, a harness session identifier) and its files are named by date. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; the record is a local file that holds no secret, and no command reads a harness transcript. Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0407 (new) records the main decision: "Roundfix appends to one record under Roundfix Home which Spec and phase each Supervisor session works on and which Work Item each prompt of a Run served." ADR-0408 (new) records the amendment: "Installation writes no Declared Skill Header into an authoring skill whose source the repository itself publishes". It changes where the header is written and when the agent hands a phase over, and the amended record stays the one that ships no name (ADR-0401: "Roundfix ships no model name"). ADR-0402 applies unchanged to this repository's declaration: "Its tool list is therefore closed by Roundfix". ADR-0405 applies, because its rule names a copy "of an authoring skill whose phase the repository declares", with one consequence that ADR-0408 states: a published skill carries no header to hide. ADR-0338 applies unchanged and is what the boundary rests on: "The mode claims the project's handoff before its inventory." ADR-0255 applies unchanged in what it removes and gains one more thing to remove, with the same manual door: "`roundfix gc` runs the same sweep without a budget". Not applicable: ADR-0290, ADR-0296 and ADR-0301, which cite ADR-0255; ADR-0353, ADR-0383 and ADR-0390, which cite ADR-0338; and ADR-0386, which ADR-0405 amended (reasons in the TechSpec's ADR Alternatives). Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the answers of 2026-11-14 to 2026-11-16 quoted in the authorization record, on the cycle BF terms, which keep the grants of the owned skills, of the Baseline source with its derived files and of the guides, and which put this repository's Project Config in this Spec's scope. No Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/_authorization.md`; bounded files: `.roundfixrc.yml`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/storage.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`, `internal/baseline/assets/modules/autonomous-work.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- A usage tool can read, from one place that does not depend on any repository's configuration, which Spec and phase a Supervisor session was working on at a given time and which Task, QA gate or review a prompt of a Run served.
- A Supervisor that follows the skills ends its session at each phase boundary with a Session Handoff, and the next session continues from the handoff alone.
- A session that goes on past a boundary is told so and loses nothing: no command refuses and no exit code changes.
- This repository's Authoring Phases the agent conducts run on declared models through the generated subagents, and the skills it publishes stay free of a model and an effort.
- A fix or a refactor is authored in the session, on the session's model, in every adopting repository, unless the user asks for the declared author.

## User Stories

1. As a person who supervises Spec work on a harness with a usage limit, I want each phase of each Spec recorded with the session that worked on it, so that my usage tool shows which phase spends the limit.
2. As the owner of a usage tool, I want the record in one documented place and format that does not change with a repository's configuration, so that I read it without the Roundfix binary and without each repository's settings.
3. As that owner, I want each prompt of a Run recorded with its Work Item and the session identifier its adapter reported, so that the tokens of a Run split by Task and by QA gate.
4. As a Supervisor, I want the skills to tell me to hand off and end the session at each phase boundary, so that the next phase starts from a small context.
5. As a person who opens the next session, I want one command to continue from the handoff, so that a boundary costs me one line.
6. As a person who authors a small fix in one session, I want to be told that the session crossed a boundary and nothing refused, so that the rule informs me and does not stop the work.
7. As the maintainer of this repository, I want its feature phases to run on the declared models through generated subagents, so that authoring stops running on whatever model the session holds.
8. As a maintainer of a repository that publishes its own authoring skills, I want installation to keep the header out of them, so that the published source never gains a model name.
9. As an adopter who declared authoring models, I want a fix authored in the session unless I ask, so that a small correction does not pay for the strongest model.

## Core Features

1. **Phase Record.** Roundfix keeps one append-only record in the directory `phase-record` under Roundfix Home, the same place on a machine whatever a repository configures. It is a set of day files named by the UTC date, one JSON object per line, each carrying the schema name `roundfix/phase-record/v1` and the time it was written. The directory and its files are readable by the user alone, and Roundfix never rewrites a line. Hazard: the record names Spec slugs, repository and worktree paths, session identifiers and times on the user's disk, and any process of that user can read it; it never holds prompt text, a token count, a price or an environment value, and a test proves that an entry holds only its declared fields.
2. **Supervisor Phase.** The record distinguishes seven phases of the Supervisor's work: `follow-up`, the four Authoring Phases `idea`, `prd`, `techspec` and `tasks`, `delivery` and `release`. `roundfix phase start <phase>` records that the session that runs it starts that phase, with the Spec named by `--spec`, the repository it ran in, the harness and the harness session read from the environment. Every phase but `follow-up` and `release` requires a Spec. Without a session identity in the environment the entry is written with none and the command says so. An unknown phase, a missing or malformed Spec and a directory outside a Git repository are refused with exit `2`, and a record that cannot be written exits `1`; nothing else fails. The command loads no Project Config, so a configuration problem never stops it.
3. **Phase Boundary and Boundary Notice.** A **Phase Boundary** follows each Authoring Phase and the delivery and release of a Spec. When `roundfix phase start` finds, among the last seven days of entries, that the same session last recorded a different phase or Spec, it prints the **Boundary Notice** on standard error, naming the earlier phase and Spec, and marks the new entry. Three cases are not a boundary: an earlier `follow-up`, `delivery` followed by `release` of the same Spec, and the same phase of the same Spec recorded again. The notice refuses nothing: standard output and the exit code are those of a command without it. Hazard: the notice is read by the agent that ran the command, so its text is fixed, names the handoff and the new session, and carries no instruction taken from the record.
4. **Prompt entries of a Run.** When a prompt of an Agent Session ends, the Daemon appends an entry with the Run, the Spec when the Run has one, the scope the Run already uses for Turn Usage (a Task, the QA gate, a review Batch or a session), the repository, the working directory the Agent Session ran in, the runtime, the session identifier the adapter reported and the times the prompt began and ended. A prompt that reported no session identifier is recorded without one. A prompt that never returns leaves no entry.
5. **Prompt entries of the pre-PR review.** `roundfix review` appends the same entry for each prompt of a pre-PR review, with no Run, the Spec its record names when it names exactly one, and the worktree it reviewed, whether a Delivery Queue or a person started it.
6. **A failed write changes nothing.** When the record cannot be written, the Daemon and `roundfix review` print one warning and the prompt, the Task, the Run and the review end as they would have. Hazard: the record is advisory, and a full disk must never fail a Run.
7. **Retention.** The Run Retention Sweep removes each day file whose date is older than its cutoff, in the same run that removes Runs, and its report names how many it removed. A dry run removes none.
8. **The skills record each phase.** The Roundfix skill's follow-up Skill Mode records `follow-up` when it starts, `delivery` when it starts to follow a Spec's delivery and `release` before a release. Each authoring skill records its own phase once the Spec's slug is known, in the session or in the Author Subagent that runs it. When the installed Roundfix does not know the command, the skill says so and continues.
9. **The session boundary rule.** The four authoring skills, the Roundfix skill's Skill Modes and the autonomous-work guide tell the Supervisor what to do at a Phase Boundary: write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed; tell the user what stays armed; start no further phase in the session. The next session starts with the follow-up Skill Mode, which claims the handoff. An Author Subagent ends no session: it returns its report and the session that started it applies the rule. With nobody present the session still ends at the boundary and says what waits; a scheduled wake-up is not a way across a boundary, because it returns to the same session. Adopters receive the guide's rule at their next `roundfix baseline update`.
10. **A published skill receives no header.** When the repository root holds the published source of an authoring skill whose phase is declared, `roundfix skills install --target project` and the skills stage of `roundfix baseline update` install that skill without the Declared Skill Header, say so in their output, and generate every declared subagent. `roundfix skills check`, the Doctor Command and the Update Check count the declaration as applied when that installed skill holds no model and no effort, report a header found there as a problem, and never say that a user-level copy of that skill hides a header. A repository that publishes no skill sees no change. Hazard: installation writes into a tree a build publishes, so the rule is read from the tree and needs no configuration key that an older binary would refuse.
11. **This repository declares its authoring models.** This repository's Project Config declares `fable` at effort `xhigh` for `idea`, `prd` and `techspec`, `fable` at effort `high` for `tasks`, and `sonnet` at effort `medium` for `research` with the two web tools and the two Exa tools as its listed tools. The five generated subagent files are committed, and the four authoring skills and their mirror stay byte-identical to each other. Hazard: a declared model name is written into committed files the harness reads; it is configuration and never a secret.
12. **A fix or a refactor is authored in the session.** In `write-techspec` and `write-tasks`, a skill the agent activated for a fix or a refactor runs the phase in the session and starts the phase's Author Subagent only when the user's request or one of the user's answers asks for it; the report says where the phase ran. A feature's phases are handed over as before. The Roundfix skill's Skill Modes say the same and ask no question about it.
13. **Documentation and glossary.** The user guide documents the record's place, its two entry kinds field by field, the phases, the notice, the retention and the new command, and says that the record is the contract an outside reader uses. The command reference, the Roundfix skill's references and the configuration guide follow the changed behavior. `CONTEXT.md` defines the four new terms and revises the terms this Spec changes, and ADR-0401 carries a notice naming ADR-0408.

## User Experience

A Supervisor, or an Author Subagent, records the phase it starts:

```text
$ roundfix phase start prd --spec 0412-example
Recorded prd of 0412-example for claude session 3f2a9c1e-…
```

The same session, later, starts the next phase without ending:

```text
$ roundfix phase start techspec --spec 0412-example
Recorded techspec of 0412-example for claude session 3f2a9c1e-…
roundfix: notice: this session already recorded prd of 0412-example; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.
```

The first line goes to standard output and the notice to standard error. At a boundary the Supervisor's last message names the Spec, the phase that ended, what stays armed and the one command that continues: the follow-up Skill Mode in a new session.

A usage tool reads the day files. A phase entry and a prompt entry look like this:

```json
{"schema":"roundfix/phase-record/v1","kind":"phase","time":"2026-11-20T13:05:11Z","phase":"prd","spec":"0412-example","repository":"/Users/me/dev/app","harness":"claude","session":"3f2a9c1e-…","boundary":false}
{"schema":"roundfix/phase-record/v1","kind":"prompt","time":"2026-11-20T15:41:02Z","started":"2026-11-20T15:22:40Z","ended":"2026-11-20T15:41:02Z","run_id":"run_20261120T152100Z_0123456789abcdef","spec":"0412-example","scope_kind":"task","scope_id":"task_02","repository":"/Users/me/dev/app","work_dir":"/Users/me/.roundfix/worktrees/app-1a2b3c4d/run_20261120T152100Z_0123456789abcdef","harness":"codex","session":"01a1cc3f-…"}
```

In this repository, and in any repository that publishes an authoring skill, installation says which header it withheld and why, and the three checks stay silent about it.

## Non-Goals / Out of Scope

- Automatic or proactive compaction: refused by the maintainer.
- Roundfix ending or opening a harness session, and a command a repository declares to open the next session with nobody present.
- A refusal, a failed check or a changed exit code when a session crosses a boundary.
- Reading harness transcripts, pricing tokens or computing a limit percentage in Roundfix, and any command that reports usage by phase.
- A command that lists or queries the record: the files are the interface.
- Changing the operator's usage script, which lives outside this repository.
- Recording from a read-only command: the Follow-up Command and its status view keep writing nothing.
- A session identity for Codex and OpenCode Supervisor sessions: not measured; such a session is recorded without one.
- The Declared Skill Header in a skill the repository publishes, and a configuration key that withholds it.
- The declaration for OpenCode and Codex: Spec 0394.
- Reducing what a session or a phase carries when it starts: the Backlog Entry the Idea recorded.
- Ending a session when its context passes a size.
- The orchestration session's model, which is the user's own configuration.

## Success Metrics

1. Success Metric: with a temporary home and a session identity in the environment, `roundfix phase start` for each of the seven phases appends exactly one line to the day file of the current UTC date, the line parses as an object with the schema name, the kind `phase` and only the declared fields, and the directory and the file are readable by their owner alone.
2. Success Metric: in the same home, a second phase recorded by the same session prints the Boundary Notice on standard error with the earlier phase and Spec, marks the entry, and leaves standard output and the exit code as they are without it; each of the three cases that are not a boundary prints no notice; a different session prints none.
3. Success Metric: without a session identity the entry is written with none and the command says so; an unknown phase, a missing Spec, a malformed Spec and a directory outside a Git repository each exit `2` and write nothing; an unwritable record exits `1`.
4. Success Metric: a fixture Run with a fake runner leaves one prompt entry for each prompt that ended, with its Run, scope, working directory, runtime, reported session identifier and both times, and a `roundfix review` with a fake runner leaves one for each review prompt with no Run; with an unwritable record the same Run and the same review end with their usual outcome and one warning each.
5. Success Metric: the Run Retention Sweep removes the day files older than its cutoff and keeps the others, names the count in its report, and removes none in a dry run.
6. Success Metric: in a temporary repository that holds the published source of one declared authoring skill and not of another, installation writes the header into the second only, generates every declared subagent, and a second run changes no byte; the three checks report the declaration as applied, and report a header placed by hand in the published skill as a problem. In a repository that publishes none, installation and the three checks print the bytes they print at the delivery base.
7. Success Metric: in this repository `roundfix skills check` prints no warning about the Authoring Model Declaration, the five generated files exist with their marker, no shipped `SKILL.md` holds a model or an effort, and each mirror equals its canonical copy.
8. Success Metric: the four authoring skills, the Roundfix skill's `modes` reference and the autonomous-work guide carry the boundary rule and the recording step; `write-techspec` and `write-tasks` carry the fix and refactor rule; `CONTEXT.md` defines **Phase Record**, **Supervisor Phase**, **Phase Boundary** and **Boundary Notice** and the changed **Declared Skill Header**, **Author Subagent** and **Session Handoff**; ADR-0401 carries an "Amended by ADR-0408" notice; the documentation, skill, Baseline and glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write, read through Exa on 2026-11-15.

- Claude Code documentation, *Manage costs effectively* (https://code.claude.com/docs/en/costs): "Claude Code sends your full conversation with every request", so "a one-line question in a session that has been open all day still draws usage for the whole conversation". This is why a session that crosses phases costs more on every call.
- The same page: "a session-scoped scheduled task fires on its interval even while the session is idle, sending your full context each time". This is why a scheduled wake-up is not a way across a boundary.
- The same page: "compacting a large context is itself a large request. When you want a fresh start instead of continuity, `/clear` costs nothing". This is why the rule ends the session with a handoff.
- Claude Code documentation, *Monitoring* (https://code.claude.com/docs/en/monitoring-usage/): `session.id` is "Claude Code's own session ID, which `/clear` replaces". This is why the record keys a phase entry by the harness session.
- Anthropic, *Effective context engineering for AI agents* (https://www.anthropic.com/research/effective-context-engineering-for-ai-agents): "overly aggressive compaction can result in the loss of subtle but critical context whose importance only becomes apparent later". This supports the maintainer's refusal of automatic compaction.
- Amp, *Amp, Rebuilt* (https://ampcode.com/news/neo): "we don't think it warrants the complexity introduced by many small, connected threads". This is the published case against the rule, and why it is an instruction with a notice that the record can later judge.

The two measurements this planning change commits under `references/` are this Spec's own evidence, not outside evidence; `## Validation` reports them.

## Unreachable Acceptance

- criterion: over the first full week after the release, at least 95% of the Claude cost of this repository's Supervisor work carries a Spec and a phase in the operator's usage report joined to the Phase Record.
  criterion: in that week the median context per call is under 150 thousand tokens in this repository's Supervisor sessions and under 200 thousand tokens in its authoring subagents, and at least 90% of the subagent cost of feature Specs runs on the generated subagents.
  criterion: over the first 20 Specs delivered after the release, the median time from a plan reaching the default branch to its delivery merge is at most 7.5 hours and the median duration of a Clean Implement Run is at most 71 minutes.
  criterion: over the same 20 Specs, at least 80% have a first Implement Run that ends Clean and the Implement Runs per delivered Spec are at most 1.25; in a week, under 5% of Task verdicts fail at settlement.
  criterion: the weekly limit points attributed to this repository on both accounts, divided by the Specs delivered in the week, fall against the week of 2026-11-08 to 2026-11-15.
  criterion: one full authoring cycle of a feature Spec in this repository, on the released binary, runs each phase the agent conducts on its declared model, which is also the pilot evidence Spec 0392 left unproven.
  reason: each needs the released binary, live sessions on the maintainer's subscriptions, the operator's usage script and weeks of delivered Specs, which no Verification or QA row may reach.
  satisfied-by: the operator's usage reports and Run listings after the release, recorded session evidence.

## Validation

The decisions were checked against these sources on 2026-11-15 and 2026-11-16.

- The usage measurement ([`references/2026-11-15-usage-by-context-and-session-identity-measurement.md`](references/2026-11-15-usage-by-context-and-session-identity-measurement.md)): 94.2% of a week's Claude cost on the Mac mini ran in sessions opened in a repository checkout and 0.5% in a Roundfix worktree; a command run by a Claude Code session reads the session's identity from the environment, and a subagent's command reads its parent's.
- The baseline of speed and plan quality ([`references/2026-11-16-delivery-speed-and-plan-quality-baseline.md`](references/2026-11-16-delivery-speed-and-plan-quality-baseline.md)), which gives the thresholds under `## Unreachable Acceptance`: each speed threshold is its baseline plus one standard error of a 20-Spec median, and the first-Run threshold is the level of the last 19 Specs.
- Characterization at `2caa1d17`. A prompt's usage is recorded once per prompt with its Run, scope, runtime and working directory, and the session identifier the adapter reported is read and then dropped for every prompt but a pre-PR review's. Every one of the 23 identifiers recorded for this repository's pre-PR reviews equals the session identifier of a Codex transcript on the same machine; the same was not measured for the Claude adapter. Nothing under Roundfix Home names a Spec's phase. The Follow-up Command and its status view write nothing.
- Installation, measured on 2026-11-16 in a throwaway clone of this repository with the declaration of Core Feature 11: it wrote a model and an effort into the four authoring skills, whose mirror then differed, generated the five subagent files and left the skills lock unchanged; with the four skills restored, `roundfix skills check` reported four headers as not the declared ones.
- The installed binary reads Project Config from the default branch and refuses a key it does not know, which is why Core Feature 10 adds none. The keys Core Feature 11 uses exist since the release that carries Spec 0392.
- The Secondbrain was read on 2026-11-15 (`wiki/index.md`, then `qmd query`): `wiki/concepts/compactacao-de-contexto-de-agentes.md` calls a restart with a handoff the way that depends least on the harness and says a poor handoff loses everything, which is why Core Feature 9 lists what the handoff carries; `wiki/concepts/claude-code-operacao-produtiva.md` gave the 150 thousand tokens of the second criterion.

## Idea-stage judgment

- item: the request in the session
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: no (0.26)
- unverified_assumption: no (0.18)
- market_question: no (0.07)
- verdict: not indicated
- stage: ran
- shape: changed — the record covers the Supervisor's sessions outside a Run and not only the Run; this repository adopts the generated subagents without the header; and a fix or a refactor is authored in the session unless the user asks for the Author Subagent.

The stage ran although the verdict was not indicated: the maintainer chose "Manter a Ideia" on 2026-11-15, after choosing "Spec com Ideia após o reset" on 2026-11-14, before the judge existed. The judged text is [`references/2026-11-15-request-read-by-the-idea-stage-judge.md`](references/2026-11-15-request-read-by-the-idea-stage-judge.md).

## Glossary

- adds: **Phase Record**
- adds: **Supervisor Phase**
- adds: **Phase Boundary**
- adds: **Boundary Notice**
- changes: **Declared Skill Header**
- changes: **Author Subagent**
- changes: **Session Handoff**

## Decisions

Autonomous route: no person present; decisions taken from the approved Idea, the maintainer's answers recorded in `references/2026-11-15-maintainer-answers.md`, the operator's cycle BF brief for Spec 0397 and the measurements above.

The maintainer's answers (AskUserQuestion), verbatim as the operator recorded them:

- 2026-11-14: "Jarvis 01a1da5d" — "Medidas já, sem Spec" and "Spec com Ideia após o reset"; "Orquestração mais barata" = "Opus medium"; automatic compaction refused: "passa para determinístico e pode perder informação importante".
- 2026-11-14: "adoção do modelo por fase da 0392 no Roundfix" — "Adotar na Spec 0397".
- 2026-11-14: "autorização" — "Mesmas do BE" (push, PR, merge, release após cada Spec, até 3 correções, Makefile/CI/go.mod/lint só se a correção exigir).
- 2026-11-15: "Ideia 0397" — "Manter a Ideia".
- 2026-11-15: "Modelo do subagente research em skills.authoring.claude do Roundfix" — "Sonnet, effort medium (Recommended)".
- 2026-11-16: "Cobertura do registro de fase" — "Run + sessões Supervisor (Recommended)": "Registra sessão, Spec e fase também fora do Run (autoria, entrega, release, triagem) e, dentro do Run, a Task ou o QA em curso."
- 2026-11-16: "Quem relata o consumo por fase" — "Roundfix só expõe (Recommended)": "Lugar e formato documentados e estáveis; o ai_usage.py cruza com os transcripts; o Roundfix não lê transcript de harness nem calcula preço."
- 2026-11-16: "Fronteira de sessão do Supervisor" — "Instrução e aviso (Recommended)": "Skills e guia mandam gravar handoff, dizer o que ficou armado e encerrar; o Roundfix avisa, sem bloquear, quando a mesma sessão registra duas fases; você abre a próxima com /roundfix follow-up."
- 2026-11-16: "Adoção dos modelos por fase no Roundfix" — "Só subagentes gerados (Recommended)": "Na 0397: autores por fase e pesquisador em Sonnet em .claude/agents/, sem cabeçalho nas skills deste repositório; uma skill digitada segue no modelo da sessão."
- 2026-11-16: "Correção pequena × declaração por fase" — "Exceção pela rota (Recommended)": "Na rota de correção ou refactor a skill escreve na sessão e só delega quando você pedir; feature sempre delega; muda o texto das skills para todos os adotantes."
- 2026-11-16: "Direção da Ideia" — "Original + respostas (Recommended)": "Registro de fase, fronteira de sessão e adoção dos modelos numa Spec; o custo de nascer de cada sessão vira Backlog Entry."
- 2026-11-16: "Acréscimo do Marcio ao critério de sucesso" — "Adicione ao critério de sucesso: com queda de performance nula ou quase nula e com aumento de qualidade de artefatos para implementações."
- 2026-11-16: "Aprovação da Ideia 0397" — "Aprovar e seguir (Recommended)".

The six answers dated 2026-11-16 were given in the night of 2026-11-15 to 2026-11-16, about 00:00 BRT. The scope the operator recorded from the maintainer's decisions of 2026-11-14 is: "o Run grava o id e a fase num lugar fixo para o ai_usage.py; o Supervisor fecha a sessão em cada fronteira de fase".

Taken from those answers, the Idea and the measurements, with no further interview:

- The Idea's Open Questions. A boundary for the notice is any change of phase or Spec in one session, with the three exceptions of Core Feature 3. A phase handed to an Author Subagent counts as a boundary, as the scope was decided. Roundfix learns the Supervisor's phase from a command the skills run, because no read-only command may write and an Idea stage runs no other Roundfix command. A harness that gives no session identity is recorded without one and receives no notice. The notice appears on the command that records the second phase and nowhere else: the status view stays a report of work, and it writes and reads no session. The record is kept as long as Runs are and removed by the same sweep. Installation knows a published skill from the tree, as ADR-0408 records. The user asks for the declared author on the fix route in the request or in an answer; the Skill Modes add no question. This repository's Research Subagent lists the two web tools and the two Exa tools, because the Secondbrain guide requires Exa and a sweep that needs a command stays with the author. The glossary takes the four names of `## Glossary`. A full authoring cycle here counts as the pilot evidence of Spec 0392, recorded by the operator.
- The record is a set of files, not a table of the Run Database with a command to read it: the request was a fixed place, and a reader that needs neither the binary nor a schema of the database keeps working across releases.
- One entry per prompt, written when the prompt ends, carries both times and the session identifier, which the adapter reports only with the result. A prompt that never returns is the accepted gap.
- The model of a call is left out of the record: each transcript already holds the model of each call, and a Run's selection is in its Turn Usage.
- `follow-up` and `release` may name no Spec: the first runs before a Spec exists, and a release can carry more than one.
- The success criterion's two added conditions are measured after the release, by the operator, over 20 delivered Specs; a week holds too few to decide. If a speed measure leaves its band, the boundary rule is the first thing to relax.
- `type: feat` and Release Ceiling `minor`: the Spec adds a command, a record and report lines. ADR-0407 and ADR-0408 record the decisions; ADR-0408 amends ADR-0401.

## Open Questions

- Does the boundary after a phase handed to an Author Subagent pay for itself? Default: it stays. The maintainer is asked again if the record's first two weeks show that the Supervisor session gains a median under 25 thousand tokens across a handed-over phase; at that growth, by the usage script's weights, a new session pays for itself only after about 68 further calls. The number is an inference from one week's prices.
- Does the Claude adapter report a session identifier equal to the one that names its transcript, as the Codex adapter does? Not measured. Default: the entry holds what the adapter reported, and the reader falls back to the working directory and the two times.
- Do Codex and OpenCode give a command the Supervisor session's identity? Not measured. Default: recorded without one; Spec 0394 can add it with their agent files.
- Will the usage script's owner read the record as documented? Default: yes; the first criterion of `## Unreachable Acceptance` is the check.
