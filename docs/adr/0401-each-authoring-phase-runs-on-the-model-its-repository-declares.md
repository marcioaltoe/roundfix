---
status: accepted
amended_by: [ADR-0408, ADR-0415, ADR-0427]
created_at: 2026-11-14T11:00:00Z
updated_at: 2026-11-16T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Each authoring phase runs on the model its repository declares

Amended by ADR-0427 (2026-11-17): the session starts an Author Subagent with the printed model and effort and switches to the declared fallback when the declared model reaches its usage limit.

A repository declares in Project Config, per Agent harness, the model and
reasoning effort of each Authoring Phase and of research sweeps. Roundfix
ships no model name. Installation turns the declaration into a header in each
installed authoring skill, which the harness applies when the user types the
skill, and into one generated subagent per phase, to which the skill hands the
phase when the agent conducts it, because neither mechanism covers both ways a
phase starts.

## Context

Authoring an Idea, a PRD, a TechSpec or a Task Graph pays for the strongest
model, and the session that authors them also monitors Runs, opens pull
requests and follows releases, where that model is wasted. The switch was
manual. On 2026-11-13 indexus applied `model` and `effort` by hand to four
authoring skills and wrote a research subagent, and reported that the next
installation erases the edit: `roundfix skills install` and `roundfix baseline
update` replace each owned `SKILL.md` whole. The maintainer ranked the pains
on 2026-11-14: forgetting to switch the model, then research sweeps weighing
on the main session, then the erased header, and set the success criterion as
a full authoring cycle with every phase on its declared model and no manual
switch, the phases the agent conducts included.

Measured on Claude Code 2.1.325 on 2026-11-14, headless: a skill header's
`model` and `effort` apply when the user types the skill and are ignored when
the agent calls it, which is the open issue anthropics/claude-code#98898; the
`model` and `effort` of a subagent file under `.claude/agents/` apply when the
agent starts that subagent, in both directions from the session's model; a
model passed in the call that starts the subagent overrides the file; a
subagent starts another subagent, which runs on its own model; a mistyped
effort is ignored without a word, and a mistyped model fails the call. Claude
Code withholds its question tool from every subagent.

## Decision

**The declaration.** Project Config gains `skills.authoring.<harness>.<key>`,
where `<key>` is `idea`, `prd`, `techspec`, `tasks` or `research` and each
entry holds `model`, `reasoning_effort` or both. This decision supports the
harness `claude`; any other harness name is refused by name until a later
decision adds it. The section belongs to Project Config: in User Config it is
ignored with a warning, because what it produces is committed with the
repository. Nothing is declared by default, and a repository that declares
nothing keeps every byte Roundfix writes today.

**The Declared Skill Header.** `roundfix skills install --target project` and
the skills stage of `roundfix baseline update` write the declared `model` and
`effort` into the frontmatter of the installed `write-idea`, `write-prd`,
`write-techspec` and `write-tasks` under `.agents/skills/`, the one copy every
harness reads. A user-level target and a `--dir` target receive none, and the
shipped bundle never holds one.

**The Author Subagent.** The same installation writes, for each declared
phase, `.claude/agents/roundfix-<phase>-author.md` with the declared model and
effort, and for `research` the Research Subagent of ADR-0402. An authoring
skill the agent activated hands its phase to that subagent, passing no model
and no effort; a skill the user typed runs where it is, on its header. An
Author Subagent cannot ask the user, so the session gathers the decisions
first, and the subagent returns its open questions instead of guessing one.

**Amended by [ADR-0408](0408-a-repository-that-publishes-an-authoring-skill-receives-no-header-for-it-and-a-fix-is-authored-in-the-session.md) (2026-11-16).** Installation writes no Declared Skill Header into an authoring skill whose source the repository itself publishes under `skills/<name>/`, and still generates every declared subagent; this repository declares the section on those terms. An authoring skill the agent activated on the fix or refactor route runs its phase in the session and starts the Author Subagent only when the user asks. A feature's phases are still handed over, a typed skill still runs on its header where one is written, and Roundfix still ships no model name.

**Amended by [ADR-0415](0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md) (2026-11-17).** Doctor and the Update Check print a Recommended Authoring Declaration when the authoring skills are installed without one. Roundfix still applies no model name by itself.

**Generated files carry a marker.** The first line of a generated subagent's
body is `<!-- roundfix:generated skills.authoring.<harness>.<key> -->`.
Installation overwrites a file at a generated name only when it carries the
marker, deletes one whose entry is no longer declared only when it carries the
marker, and refuses, before it writes anything, when a file without the marker
sits at a generated name.

**The header is part of the installation.** An owned skill is still held to
the binary's copy by its version, and the lock entry still records the
installed tree, header included. `roundfix skills check`, Doctor and the
Update Check warn when the installed files do not match the declaration and
name `roundfix skills install --target project`; none of them fails for it.

**An older Roundfix.** A binary older than the release that carries this
decision cannot learn the section. It refuses every command that loads
configuration with `field skills not found in type config.configOverlay`, and
its installation writes the skills without the header. Once the Setup
Manifest records the newer release, a binary that knows ADR-0392 refuses
first, naming `roundfix upgrade`.

## Considered Options

- **Keep a header the repository wrote by hand.** The smallest change, and it
  adds no configuration. Rejected by the maintainer: it answers the pain he
  ranked last and leaves the phases the agent conducts on the session's model.
- **The header alone.** Rejected: measured, the header is ignored when the
  agent calls the skill, which is how most authoring starts.
- **A header with `context: fork`.** Measured to run an agent-called skill on
  the header's model. Rejected: a forked skill starts without the
  conversation and without the question tool, so a phase the user types would
  lose its interview.
- **A rule that tells the session to switch models.** Rejected: a rule is
  what the maintainer already had and forgot, and an agent cannot change the
  model of its own session.
- **Model names shipped as defaults.** Rejected: a name ages and assumes one
  vendor; the documentation shows a marked example.
- **One entry per skill.** Rejected: a phase is what the person thinks in,
  and `research` names no skill.
- **The section in User Config.** Rejected: the header and the generated
  files are committed, so one developer's choice would reach every clone.
- **A lock hash computed without the header.** Rejected: ADR-0359 keeps one
  rule, the installed tree, and what the external skills CLI does with the
  other hash was not verified for both lock versions.
- **Failing a check for an unapplied declaration.** Rejected: the installed
  skills stay usable, and `make verify` and the release plan read those
  checks.

## Consequences

- A repository that wrote a header by hand declares it instead. An
  installation that finds `model` or `effort` in an authoring skill with no
  entry for its phase says so once and names the key to declare.
- The external skills CLI's `skills update` replaces an installed skill whose
  hash differs from the catalog's, so it can erase a Declared Skill Header;
  the three checks then name the unapplied declaration.
- A user-level copy of an authoring skill hides the header in Claude Code;
  ADR-0405 reports it.
- This repository's `.agents/skills` are the source of the shipped bundle, so
  it does not declare the section until a later decision says how; a test
  keeps `model` and `effort` out of every shipped `SKILL.md`.
- A later decision adds the harnesses `opencode` and `codex` under the same
  keys, compiled into the agent files those harnesses read.
