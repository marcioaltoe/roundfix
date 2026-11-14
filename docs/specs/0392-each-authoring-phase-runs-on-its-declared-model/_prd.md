---
spec: 0392-each-authoring-phase-runs-on-its-declared-model
status: active
created: 2026-11-14
surfaces: [backend, cli, docs]
type: feat
---

# Each authoring phase runs on its declared model

Authoring an Idea, a PRD, a TechSpec or a Task Graph pays for the strongest
model, and the session that authors them also monitors Runs and follows
releases, where that model is wasted; the switch is manual, so it is
forgotten. After this Spec a repository declares once, in Project Config, the
model and reasoning effort of each authoring phase and of research sweeps for
Claude Code, and Roundfix makes each phase run on that declaration with no
manual switch. A phase the user starts by typing the skill runs on a header
that installation writes into the installed skill and keeps across releases.
A phase the agent conducts is handed to a generated subagent on the declared
model, because Claude Code ignores a skill's header when the agent calls the
skill. Research sweeps go to a generated subagent that cannot run a command or
write a file. Roundfix ships no model name, and a repository that declares
nothing keeps today's behavior. The notice about a user-level copy of a
Roundfix skill gains the command that removes the copy, and names a copy that
hides the declared header whatever its version.

## Sources

Two ai-memory messages, adopted by this Spec's planning change with the
records behind its decisions (see `references/_index.md`):

1. [`references/2026-11-13-indexus-model-and-effort-per-skill-and-a-research-subagent.md`](references/2026-11-13-indexus-model-and-effort-per-skill-and-a-research-subagent.md),
   ai-memory message `01a1d68d-7d5c-7a22-8a37-85643843fe88`, copied verbatim;
   reported by indexus on 2026-11-13. It declared no blocked work and expects
   no reply.
2. [`references/2026-11-13-pantheon-roundfix-skills-left-the-user-level.md`](references/2026-11-13-pantheon-roundfix-skills-left-the-user-level.md),
   ai-memory message `01a1d6a1-2ce5-78d3-9998-e41bfdbd65e8`, copied verbatim;
   reported by pantheon on 2026-11-13. Its second request, the notice text, is
   Core Feature 6; its first concerned the operator's own notes and is done.
   It declared no blocked work and asked for no notice.
3. [`references/2026-11-14-maintainer-answers.md`](references/2026-11-14-maintainer-answers.md),
   the maintainer's answers of 2026-11-14, recorded by the operator.
4. [`references/2026-11-14-skill-model-override-measurement.md`](references/2026-11-14-skill-model-override-measurement.md),
   recorded by the operator, and
   [`references/2026-11-14-subagent-model-and-effort-measurement.md`](references/2026-11-14-subagent-model-and-effort-measurement.md),
   recorded by the author.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from.
At the authoring base this project's ai-memory inbox holds no pending message,
`docs/findings/` is empty, and the one open Backlog Entry, a Conexus report of
a false status and a silent authorization wait, shares no context with this
Spec.

## Prerequisites

None in `requires`. Authored against `251ed81e` (v0.143.3) and rebased on
`f17949de`, which adds the plan of Spec 0391 and changes no code. Specs 0391
and 0393 are authored beside this one against the same base and deliver in
the order 0391, 0392, 0393. Files this Spec shares with them, each edited in
its own part:

- The four authoring skills, `write-idea`, `write-prd`, `write-techspec` and
  `write-tasks`, canonical and mirrored: 0391 rewrites their grouping
  sentences and the Task sizing guidance, and 0393 may change when the Idea
  stage runs. This Spec adds one section under its own heading to each and
  changes no existing sentence.
- The version fields of those four skills and of the Roundfix skill, and
  `skills/testdata/owned-skill-versions.json`: each Spec raises the next
  version with the record command on the tree its Task runs on and never
  writes a number.
- `CONTEXT.md`: 0391 changes **Grouping Bound**; this Spec adds its own terms
  and changes **Stale User Skill Copy**.
- The Roundfix skill's references: 0391 edits `spec`; this Spec edits
  `setup`, `baseline` and `modes`, and 0393 owns the judge's reference text.
- `docs/references/behavior-surfaces.json` and
  `docs/references/skill-coverage.json`: each Spec re-records them after its
  own last guide edit.
- The Baseline module `spec-workflow`, its version record and the derived
  Baseline files belong to 0391 and 0393; this Spec changes no Baseline
  module.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the
  generated subagents carry five fixed names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added, and a declared model name is written into files the
  harness reads and is never a secret. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0401 (new) records the main
  decision: "A repository declares in Project Config, per Agent harness, the
  model and reasoning effort of each Authoring Phase and of research sweeps."
  ADR-0402 (new) bounds the research subagent: "Its tool list is therefore
  closed by Roundfix". ADR-0405 (new) amends ADR-0386: "They also name,
  whatever its version, a copy in Claude Code's user-level skill directory of
  an authoring skill whose phase the repository declares". ADR-0386 is
  amended in its comparison and its commands and keeps its limit: "None of
  them writes to those directories." ADR-0392 applies unchanged and covers an
  older binary once the Setup Manifest records the newer release: "A Roundfix
  binary older than that release refuses, before it reads Project Config,
  every command that loads the repository's configuration". ADR-0359 applies
  unchanged: installations still "write, in each existing lock entry that
  names an installed owned skill" the installed tree's hash, header included.
  ADR-0364 and ADR-0189 apply unchanged, so a header never makes an installed
  skill fail: "An installed owned copy passes when the version it declares is
  equal to or higher than the version the running binary carries", and "The
  minimum version of an owned skill is now the version the binary carries for
  it". ADR-0393 applies unchanged: the hash follows the lock's version, and
  "Roundfix never changes that version". ADR-0341 applies unchanged: "A
  restore keeps the Claude Skill Links.", through which Claude Code reaches
  the installed skill and its header. ADR-0371 applies unchanged: the Update
  Check gains lines and "It is read-only and always" exits zero for a report.
  Not applicable: ADR-0221, ADR-0224, ADR-0233, ADR-0241, ADR-0250, ADR-0256
  and ADR-0264, which cite ADR-0189, ADR-0359 or ADR-0364 (reasons in the
  TechSpec's ADR Alternatives). Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  answers of 2026-11-14 quoted in the authorization record, on the cycle BF
  terms, which keep the grant of the owned skills; the Governed Paths are the
  `SKILL.md` pairs of the Roundfix skill and of the four authoring skills and
  the Roundfix skill's canonical `baseline`, `modes` and `setup` references.
  No Makefile, workflow, Go module, lint or formatter change, no Baseline
  module change, and no change to this repository's Project Config. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0392-each-authoring-phase-runs-on-its-declared-model/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- A repository declares once, in Project Config, the model and reasoning
  effort of each authoring phase and of research sweeps for Claude Code.
  Roundfix ships no model name, and a repository that declares nothing sees
  no change.
- A phase the agent conducts runs on its declared model and effort, with the
  decisions it needs gathered from the user first.
- A phase the user starts by typing the skill runs on its declared model and
  effort, and installations and Baseline updates keep it so.
- Research sweeps run on the declared research model, in a subagent that
  holds no shell and no file-writing tool.
- The read-only checks say when the installed files do not match the
  declaration, and name each reported user-level copy with the command that
  removes it.

## User Stories

1. As a person who supervises Spec authoring through Claude Code, I want to
   declare each phase's model and effort once in the repository, so that
   nobody switches the session's model at a phase boundary.
2. As that person asking the agent to author a Spec, I want each phase to run
   on its declared model although the agent calls the skill, so that the
   authoring gets the strong model and the orchestration keeps the economical
   one.
3. As a user who types an authoring skill, I want that turn to run on the
   declared model and the next Roundfix release not to erase the setting, so
   that nothing is edited by hand again.
4. As an adopter whose authoring reads web pages and knowledge-base notes, I
   want sweeps to run in a subagent that cannot run a command or write a
   file, so that an instruction injected into what it reads has nothing to
   run.
5. As a developer whose home holds a copy of a Roundfix skill, I want the
   notice to give me the command that removes it, so that the repository's
   copy and its header are what the harness loads.
6. As an adopter who changed the declaration or updated Roundfix, I want
   `roundfix skills check`, Doctor and the Update Check to say when the
   installed files do not match, so that one command puts them right.
7. As a teammate still on an older Roundfix, I want the documentation to say
   what my binary prints on a repository that declares the section, so that I
   know to upgrade.

## Core Features

1. **Authoring Model Declaration.** Project Config gains a `skills` section
   whose `authoring` part declares, per Agent harness, an entry for each
   **Authoring Phase**, `idea`, `prd`, `techspec` and `tasks`, and one for
   `research`. An entry holds a model, a reasoning effort or both. This Spec
   accepts the harness `claude` and refuses any other harness name, any other
   key and an empty entry by name. A model is one token the harness
   understands, never validated against a list of names; a reasoning effort
   for `claude` is one of `low`, `medium`, `high`, `xhigh` and `max`, because
   Claude Code ignores a mistyped effort without a word. The section belongs
   to Project Config: in User Config it is ignored with a warning. Hazard: a
   declared model and effort are written into committed files the harness
   reads; they are configuration and never a secret.
2. **Installation applies the declaration.** `roundfix skills install
   --target project` and the skills stage of `roundfix baseline update` read
   the declaration and, after writing the bundle, write the **Declared Skill
   Header**, the declared model and effort, into the frontmatter of each
   installed authoring skill whose phase is declared, and write one generated
   subagent file for each declared entry under `.claude/agents/`. Every
   generated file carries a marker line. Hazards: installation overwrites and
   deletes files, so it overwrites a file at a generated name only when it
   carries the marker, deletes a marked file whose entry is no longer
   declared, and refuses before writing anything when a file without the
   marker sits at a generated name; it also overwrites each installed
   `SKILL.md` whole, as today, so when an installed authoring skill held a
   model or effort that no entry covers, it says so once and names the key to
   declare. A user-level target and a `--dir` target receive no header and no
   generated file. A second installation with the same declaration changes
   nothing. The lock entry of an installed skill records the installed tree,
   header included.
3. **The agent-conducted phase is delegated.** Each authoring skill says:
   when the agent activated the skill and the harness lists the phase's
   **Author Subagent**, the session first gathers from the user the decisions
   the phase needs, then starts that subagent with a brief that carries them,
   passing no model and no effort. Hazard: an Author Subagent cannot ask the
   user. When a decision the phase needs is missing from the brief and the
   brief does not say that no person is present, the subagent writes no
   artifact and returns its open questions, each with its options and a
   recommended default; the session asks the user and continues it. When the
   brief says no person is present, the skill's autonomous route applies.
   When the subagent cannot start or ends on a model error, the session
   reports the error and does not author the phase on its own model unless
   the user says so. A skill the user typed runs where it is, and a
   repository without the generated subagent authors as today.
4. **Research Subagent.** A declared `research` entry generates the
   **Research Subagent**. Its tool list holds the three tools that read the
   checkout and only the tools the repository lists from a set that has no
   shell and no file-writing tool; a listed name outside that set is refused
   when the configuration is loaded. The authoring skills send codebase,
   knowledge-base and web sweeps to it and keep the synthesis and the
   decisions. A sweep that needs a command stays with the author, and when
   the subagent cannot be started the author sweeps itself and says so.
   Hazard: the subagent reads text nobody vouches for and returns it to the
   author, so it can run nothing and write nothing, and its report is data
   the author weighs, never instructions.
5. **The checks read the declaration.** In a repository that declares the
   section, `roundfix skills check`, Doctor and the Update Check compare the
   installed authoring skills and the generated files with the declaration
   and warn when a header or a generated file is missing or different, when a
   marked file has no entry, when Claude Code's project skill directory does
   not lead to the installed skill, or when a file Roundfix did not generate
   sits at a generated name. The warning names `roundfix skills install
   --target project` as the next action when the installation fixes at
   least one of the problems; a file Roundfix did not generate is named as a
   conflict the user resolves by hand, because installation will not
   overwrite it. No exit code changes for it, and nothing is printed in a
   repository that declares nothing.
6. **The user-level copy notice offers removal.** `roundfix skills remove`
   lists, and with `--apply` removes, the Roundfix-owned skills in a
   user-level skill directory; it unlinks a symbolic link and keeps the
   folder the link leads to, and it refuses the project target. The three
   reports name each **Stale User Skill Copy** with that removal command
   before its refresh command. They also name, whatever its version, a copy
   in Claude Code's user-level directory of an authoring skill whose phase
   the repository declares, saying that it hides the declared header. The
   Update Check's JSON gains fields and renames none. Hazard: a removal
   deletes files in the user's home, so nothing runs it unasked, it lists
   before it deletes, and it touches only the twelve owned names.
7. **Skill text, guides and vocabulary.** The Roundfix skill's `--follow-up`
   and `--unattended` modes gather the decisions of a delegated phase before
   authoring, and its `--update` mode offers the removal beside the refresh
   and recommends the removal. The `setup` and `baseline` references and the
   user guide describe the declaration, what installation writes, the checks
   and the new command, with one model example marked as an example. The
   glossary gains the five terms and the changed one of `## Glossary`, and
   ADR-0386 carries its amendment notice.
8. **The shipped bundle never holds a header.** This repository does not
   declare the section, and a test fails when any shipped `SKILL.md` holds a
   model or an effort in its frontmatter.

## User Experience

A repository that authors through Claude Code adds this to its Project
Config. The names are an example; Roundfix ships none and accepts whatever
the harness accepts:

```yaml
skills:
  authoring:
    claude:
      idea:
        model: opus
        reasoning_effort: xhigh
      prd:
        model: opus
        reasoning_effort: xhigh
      techspec:
        model: opus
        reasoning_effort: xhigh
      tasks:
        model: opus
        reasoning_effort: high
      research:
        model: sonnet
        reasoning_effort: medium
        tools: [WebSearch, WebFetch]
```

`roundfix skills install --target project` then names each header it applied
and each subagent file it wrote, and the repository commits them with the
skills. From then on, typing `/write-prd` runs that turn on the declared
model, and asking the agent for a PRD makes the session ask its questions
first and hand the phase to `roundfix-prd-author`, which sends its reading to
`roundfix-researcher`. A session on an economical model stays on it.

After a later edit of the section, or an update that erased a header, Doctor
prints one `authoring-models` line that warns and names the installation
command; when everything matches, the line says so. A developer whose home
still holds `write-prd` sees it named in the user-level notice with the
command that removes it.

A developer on a Roundfix older than the release that introduces the section
sees every command that loads configuration refuse the repository, naming the
unknown field `skills` and no fix; once a Baseline change applied by the newer
release has recorded its version, the refusal names both versions and
`roundfix upgrade`.

## Non-Goals / Out of Scope

- Agent files for OpenCode and Codex compiled from the same declaration: Spec
  0394, authored after this one is delivered. The section is keyed by harness
  so that 0394 adds `opencode` and `codex` beside `claude` and renames
  nothing.
- Model names shipped as defaults, and any list of model names Roundfix
  validates against.
- Changing how a Run selects its Agent: Agent Selection Profiles are
  untouched.
- An Author Subagent that interviews the user, and repairing Claude Code's
  skill-header defect.
- Declaring the section in the Roundfix repository itself; the maintainer
  decides after the release, and Core Feature 8 guards the bundle until then.
- Reporting tokens or cost by phase.
- Keeping a header a repository wrote by hand without declaring it.
- A lock hash computed without the header, and any change to what the
  external skills CLI does on `skills update`.
- Removing a user-level copy during an installation, an upgrade or a check,
  and naming a user-level copy that is not older in a repository that
  declares nothing.
- Creating Claude Code's project skill links for the authoring skills; the
  check names a missing one.
- A friendlier refusal in binaries already released: they are out of reach.
- Any Baseline module, guide under `docs/agents/`, or Vendored Skill.

## Success Metrics

1. Success Metric: in a temporary repository whose Project Config holds the
   example above, `roundfix skills install --target project` leaves each of
   the four installed authoring skills equal to the shipped file plus its
   declared model and effort lines, writes the five generated files with
   their marker, and a second run changes no byte; `roundfix baseline update`
   does the same through its skills stage. Without the section, the
   installed tree equals the shipped bundle byte for byte and nothing is
   created under `.claude/agents/`.
2. Success Metric: for every accepted declaration the generated Research
   Subagent's tool list is the three reading tools followed by the declared
   names, and a declaration that lists a shell, a file-writing tool or the
   tool that starts subagents is refused when the configuration is loaded,
   naming it.
3. Success Metric: installation refuses, before any write, a file without the
   marker at a generated name; it deletes a marked file whose entry was
   removed and leaves every other file under `.claude/agents/` untouched.
4. Success Metric: with a declaration whose files are missing or different,
   `roundfix skills check`, Doctor and the Update Check each print the
   warning with the installation command and keep their exit codes; with
   matching files Doctor prints its `authoring-models` line as `ok`; without
   the section every one of the three prints the bytes it prints at the
   delivery base.
5. Success Metric: with a temporary home, the three reports name the removal
   command before the refresh command; a copy at the binary's version of a
   declared authoring skill in Claude Code's user-level directory is named;
   `roundfix skills remove` lists without `--apply`, removes with it, unlinks
   a link while its target survives, and refuses the project target.
6. Success Metric: no shipped `SKILL.md` holds a model or an effort in its
   frontmatter, and the mirror of each owned skill equals its canonical copy.
7. Success Metric: the four authoring skills and the Roundfix skill's `modes`
   reference carry the delegation text; `CONTEXT.md` defines **Authoring
   Model Declaration**, **Authoring Phase**, **Declared Skill Header**,
   **Author Subagent** and **Research Subagent** and the changed **Stale User
   Skill Copy**; ADR-0386 carries an "Amended by ADR-0405" notice; the
   documentation, skill and glossary contracts pass.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write,
read through Exa on 2026-11-14.

- Claude Code, *Create custom subagents*, frontmatter reference
  (https://code.claude.com/docs/en/sub-agents): `model` is the "Model to use:
  `sonnet`, `opus`, `haiku`, `fable`, a full model ID such as
  `claude-opus-5-5`, or `inherit`", and `effort` is the "Effort level when
  this subagent is active. Overrides the session effort level". This is why
  a generated subagent can carry a phase's declaration.
- The same page, *Available tools*: "The first filter removes these tools,
  even when listed in the `tools` field", and the list that follows holds
  `AskUserQuestion`. This is why the session gathers the decisions before it
  delegates, and why an Author Subagent returns questions.
- The same page, *Choose a model*: the first source of a subagent's model is
  "The per-invocation `model` parameter, which Claude can pass when it starts
  the subagent." This is why the skill text forbids passing one.
- Claude Code, *Extend Claude with skills*, frontmatter reference
  (https://code.claude.com/docs/en/skills): a skill's `model` is the "Model
  to use when this skill is active. The override applies for the rest of the
  current turn and isn't saved to settings.", and on precedence, "Enterprise
  over personal, and personal over project." The first is why the header
  serves a typed skill; the second is why a user-level copy hides it.
- anthropics/claude-code issue 98898
  (https://github.com/anthropics/claude-code/issues/98898), open on
  2026-11-14: "A skill's `model:` and `effort:` frontmatter is applied when
  the user types `/skill-name`, but silently ignored when Claude invokes the
  same skill via the Skill tool." This is why the header alone does not
  reach a phase the agent conducts.
- OWASP, *LLM06:2025 Excessive Agency*
  (https://genai.owasp.org/llmrisk/llm062025-excessive-agency/): "Limit the
  extensions that LLM agents are allowed to call to only the minimum
  necessary." and "Avoid the use of open-ended extensions where possible
  (e.g., run a shell command, fetch a URL, etc.)". This is why the Research
  Subagent's list is closed and a tool that fetches a URL is the
  repository's explicit choice.

The measurements this planning change commits under `references/` are this
Spec's own evidence, not outside evidence; `## Validation` reports them.

## Unreachable Acceptance

- criterion: in a pilot adopting repository on the released binary, one full
  authoring cycle the agent conducts (Idea, PRD, TechSpec, Tasks) records
  every authoring turn on its declared model and effort and no `/model` or
  `/effort` command, read from `message.model` and `effort` in the session
  and subagent transcripts.
  criterion: in the same pilot, a phase the user starts by typing its skill
  in an interactive session runs that turn on the declared model and effort.
  criterion: in the same pilot, more than 80% of the tokens of the research
  sweeps are used by the declared research model, read from the usage by
  model of the session.
  criterion: in the same pilot, a sweep that asks the Research Subagent to
  write a file writes none.
  criterion: after `roundfix baseline update` with the released binary in
  three adopting repositories that declare the section, `git status` shows no
  file that someone must edit by hand again.
  reason: each needs a released binary, a live Claude Code session on a
  subscription and adopting repositories, which no Verification or QA row may
  reach.
  satisfied-by: the operator's pilot evidence after the release, recorded
  session evidence.

## Validation

The decisions were checked on 2026-11-14 against these sources.

- Characterization at `251ed81e`. Installation writes every shipped file
  over the installed one and compares nothing, so a hand-written header is
  erased; `roundfix baseline update` runs the same installation.
  `roundfix skills install` and `roundfix skills check` read no Project
  Config, only the Git root; the history stage of `roundfix baseline update`
  already loads it. An installed owned skill is held to the binary by its
  declared version only. After the installation the lock entry of each owned
  skill receives the hash of the installed tree, and Doctor compares that
  hash with the installed tree. Roundfix writes nothing under `.claude/`
  except one optional link for the Roundfix skill, and deletes an installed
  skill only through the Retired Skill proof. The user-level notice names a
  copy only when it is older than the binary and offers only the refresh.
  The built binary, given a Project Config with a `skills` section, refused
  every command that loads configuration with `field skills not found in
  type config.configOverlay` and installed the bundle without a header.
- The author's measurement on Claude Code 2.1.325, headless, twenty-two runs
  ([`references/2026-11-14-subagent-model-and-effort-measurement.md`](references/2026-11-14-subagent-model-and-effort-measurement.md)).
  The first Open Question of the Idea is settled: the `model` and the
  `effort` of a subagent file apply when the agent starts the subagent, in
  both directions from the session's model. A `tools` list was a boundary
  the subagent could not cross, with a control that crossed it. A subagent
  started another, which ran on its own model. A model passed in the call
  overrode the file. A header with both keys applied to a typed skill reached
  through a link and not to the same skill called by the agent. A mistyped
  effort was ignored without a word; a mistyped model failed the call. A
  header with `context: fork` ran an agent-called skill on the header's
  model; it is not used, because a forked skill has neither the conversation
  nor the question tool.
- Codex and OpenCode read the same installed copy. Codex's skill parser
  reads `name`, `description` and `metadata` and drops unknown keys, and
  OpenCode adds a skill's body "without frontmatter" to the conversation, so
  a header written for Claude Code changes nothing there. The fifth Open
  Question of the Idea is settled by reading both; neither was run.
- Secondbrain, read only: `wiki/index.md`, then four `qmd query` searches.
  `raw/web/akita/2026-04-25-akita-llm-benchmarks-vale-a-pena-misturar-2-modelos.md`
  challenges the premise: in three rounds a strong planner with a cheap
  executor lost to the strong model alone, and delegation paid only for
  numerous independent tasks. It shaped two choices: a phase is handed whole
  to one subagent and never split between a planner and an executor, and
  only sweeps, which are independent reading, go to the cheaper model.
  `wiki/concepts/claude-code-operacao-produtiva.md` says the design must
  exist before the spawn, with limits and a stop point, and that
  coordination cost does not disappear; it shaped the brief and the rule
  that an Author Subagent stops and returns its questions.
  `raw/web/2026-08-18-aihero-ai-coding-user-vs-project-skills-126e1198.md`
  describes a user-level skill as personal and global and a project skill as
  shared with every clone, which is the reason the declaration lives in the
  repository and a user-level copy is reported. The search for prior notes
  on a subagent that reads untrusted text returned nothing relevant. The
  Idea already read the Oraculum Backlog Entry and the page on costs and
  limits of coding agents.
- Exa: the six published sources of `## Acceptance evidence`, the Claude
  Code passages on nested subagents ("By default, a subagent can spawn
  subagents of its own, up to three layers below the main conversation.")
  and on a skill's `context: fork`, the Codex parser source and the OpenCode
  skills page. The Idea's `## Sources` lists the sources read before.

## Glossary

- adds: **Authoring Model Declaration**
- adds: **Authoring Phase**
- adds: **Declared Skill Header**
- adds: **Author Subagent**
- adds: **Research Subagent**
- changes: **Stale User Skill Copy**

## Decisions

Autonomous route: no person present; decisions taken from the approved Idea, the maintainer's answers of 2026-11-14 recorded in `references/2026-11-14-maintainer-answers.md`, the operator's cycle BF brief for Spec 0392 and the measurements above.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-14), verbatim:

- 2026-11-14: "O Indexus (mensagem 01a1d68d, pedido seu de 13/11) pede que `roundfix skills install --target project` e o Baseline escrevam `model` e `effort` no cabeçalho das skills e gerem um subagente de pesquisa em `.claude/agents/`, a partir de uma seção `skills` do `.roundfixrc.yml`. É feature (PRD + TechSpec + Tasks, release minor v0.144.0). O que faço com o pedido?" — "Spec completa": "Uma Spec com a seção `skills` no Project Config, o cabeçalho model/effort na instalação, o subagente de pesquisa só leitura e as skills de autoria passando a usá-lo. Cabe em 4 Tasks + QA."
- 2026-11-14: "O Pantheon (mensagem 01a1d6a1) relata a sua decisão de 13/11: skills do Roundfix só instaladas no projeto, nunca no nível do usuário (já removidas no mini e no MacBook; conferi no mini). Ele pede que o aviso de cópia de usuário antiga, em `skills check`, `doctor` e `update check`, pare de indicar só `refresh: roundfix skills install --target …`. O que faço com o texto do aviso?" — "Spec pequena (Recommended)": "Correção própria: o aviso passa a oferecer apagar a cópia de usuário, com o comando, ao lado do refresh. 1 Task + QA, release patch. A memória antiga é aposentada de qualquer forma."
- 2026-11-14: "Você respondeu "Spec pequena" para o ajuste do aviso de cópia de usuário (Pantheon). Pelo motivo que você lembrou, ela seria justamente uma Spec de 1 Task pagando um ciclo inteiro de Run, QA, review e release. O ajuste é do mesmo contexto da Spec do Indexus (instalação de skills e cópias de usuário), e ninguém está bloqueado por ele. Junto ou mantenho separado?" — "Juntar à Spec do Indexus (Recommended)": "O ajuste do aviso entra numa Task da Spec do Indexus que já mexe em `skills check`. Um ciclo de entrega e uma release (minor v0.144.0); a Spec adota as duas mensagens e continua em 4 Tasks + QA."
- 2026-11-14: "As Specs que saírem das duas respostas acima precisam de autorização de push, merge e release. A entrega local só entra na fila com o Quaestor ocioso (regra do mini, 1 vaga de verificação). Qual autorização vale para elas?" — "Mesmas do BE (Recommended)": "Push, PR e merge; release após cada Spec (patch para correção, minor para feature); até 3 correções por Spec; Makefile, CI, go.mod e lint só quando uma correção exigir, descrito no relatório."
- 2026-11-14: "Etapa de Ideia da 0392, pergunta 1 (o problema). O pedido do Indexus junta três dores diferentes: trocar de modelo à mão, o cabeçalho posto à mão que a instalação apaga, e as varreduras de pesquisa pesando na sessão principal. Qual delas mais pesa hoje para você?" — "1 > 3 > 2".
- 2026-11-14: "Etapa de Ideia da 0392, pergunta 2 (quem usa). A pesquisa confirmou que só o Claude Code troca modelo e esforço pelo cabeçalho da skill; Codex e OpenCode ignoram essas chaves e só aceitam modelo em arquivos de agente próprios. Hoje você contorna com a regra "Model by phase" do seu CLAUDE.md e com o `/model` à mão. Para quem a primeira versão precisa funcionar?" — "Claude Code, qualquer adotante (Recommended)": "Quem autora Specs pelo Claude Code em qualquer repositório que adota o Roundfix. Nada supõe a sua máquina: sem Fable, Exa ou ai-memory embutidos; Codex e OpenCode seguem como hoje."
- 2026-11-14: "Etapa de Ideia da 0392, pergunta 3 (sucesso). Você ordenou as dores: esquecer de trocar o modelo, depois as varreduras na sessão principal, depois o cabeçalho apagado. Daqui a um mês, qual resultado faria você dizer que a feature valeu?" — "Ciclo sem troca manual (Recommended)": "Num ciclo inteiro de autoria (Ideia, PRD, TechSpec, Tasks) cada fase roda no modelo e esforço declarados sem nenhum `/model` ou `/effort` digitado, inclusive quando é o agente que chama a skill. Medido pelo modelo que cada turno registra."
- 2026-11-14: "Etapa de Ideia da 0392, pergunta 4 (a forma). Medi que o Claude Code 2.1.325 ignora o `model` do cabeçalho da skill quando é o agente que a invoca (só aplica quando você digita `/skill`), e o seu critério de sucesso inclui esse caminho. Qual forma a ideia toma?" — "Mais ambiciosa": "A híbrida, mais um subagente autor gerado com o modelo forte e a mesma declaração compilada para arquivos de agente do OpenCode e do Codex. Duas Specs."
- 2026-11-14: "Rascunho da Ideia 0392 (cada fase de autoria roda no modelo declarado): duas Specs, a 0392 para o Claude Code (declaração no Project Config, subagente autor e subagente de pesquisa gerados, cabeçalho na instalação, aviso de cópia de usuário) e a 0394 para OpenCode e Codex. […] Aprova?" — "Aprovado (Recommended)".

Taken from those answers, the Idea and the measurements, with no further
interview:

- The Idea's Open Questions. A subagent file's model and effort apply when
  the agent starts it: measured, yes. The declaration is made per phase, as
  the Idea's default says, and per harness, so that Spec 0394 extends it.
  Codex and OpenCode ignore the header's keys: settled by reading both. An
  older Roundfix: the section requires the release that introduces it, the
  documentation says what an older binary prints, and ADR-0392 covers the
  case once the newer release has applied a Baseline change. An interactive
  session and the Roundfix repository's own adoption stay open, below.
- The section keeps the name `skills` the maintainer's answer gives it, and
  its effort key is `reasoning_effort`, the word Project Config already uses
  in `profiles`; installation writes it as the harness's own `effort`.
- The reasoning effort for `claude` is validated against the five documented
  levels because a mistyped one is ignored in silence; a model is not
  validated against names because none is shipped and a wrong one fails the
  call aloud.
- A phase is handed whole to one Author Subagent, which reads the installed
  skill from its folder, so the delegation does not depend on a harness
  loading the skill by name and Spec 0394 can reuse it.
- The Research Subagent's fixed tools read the checkout only. Web tools and
  MCP tools are the repository's to list, because the maintainer decided
  that no search provider and no memory server is built in, and because a
  tool that fetches a URL can also send out what was read.
- A hand-written header is not preserved; the maintainer chose generation.
  Installation names the key to declare when it finds one.
- An unapplied declaration warns and never fails a check, as a stale lock
  entry does today.
- The removal is a Roundfix command, not a printed shell line, so that it is
  the same on every system, touches only the owned names and lists first.
- A user-level copy that is not older is named only where it hides a
  declared header: in Claude Code's user-level directory, for a declared
  phase's skill.
- The `--update` Skill Mode recommends the removal over the refresh: the
  maintainer decided on 2026-11-13 that the Roundfix skills live only in the
  project, and a refreshed copy goes stale again at the next release.
- `type: feat` and Release Ceiling `minor`: the Spec adds a configuration
  section, a command and report fields. ADR-0401, ADR-0402 and ADR-0405
  record the decisions; ADR-0405 amends ADR-0386.

## Open Questions

- Does an interactive Claude Code session behave like the headless
  measurement, for a typed skill and for a subagent the agent starts?
  Default: yes; the second criterion of `## Unreachable Acceptance` measures
  it once in the pilot.
- Does the Roundfix repository declare the section? Default: no, until the
  maintainer decides after the release. Its installed skills are the source
  of the shipped bundle, so adopting it needs a way to keep the header out of
  that source first.
- How does the external skills CLI's `skills update` treat an authoring
  skill whose installed copy carries a header, for each lock version? Not
  measured. Default: it may replace the copy, and Core Feature 5 reports the
  missing header with the command that restores it.
