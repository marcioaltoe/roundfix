---
spec: 0392-each-authoring-phase-runs-on-its-declared-model
created: 2026-11-14
status: proposed
leverage: strategic-bet
---

# Each authoring phase runs on its declared model

## Overview

A repository that adopts Roundfix declares, once, in its Project Config, the model and reasoning effort for each authoring phase (Idea, PRD, TechSpec, Tasks) and for research sweeps. Roundfix then makes each phase run on that declaration without anyone switching the session's model by hand: through the skill header when the user starts the phase by typing the skill, and through a generated subagent when the agent conducts the phase.

It is for the person who supervises Spec authoring through an Agent harness, in any adopting repository. The first version serves Claude Code completely; a second Spec compiles the same declaration into agent files for OpenCode and Codex.

It is valuable because authoring a Spec is where a stronger model pays off and orchestration is where it is wasted, and today the right model depends on someone remembering to switch at every phase boundary.

V1 is ambitious on purpose: the maintainer chose the shape that covers the agent-conducted path and reaches the three harnesses, in two Specs, over a smaller change that only preserved a hand-written header.

## Problem

The maintainer ranked three pains (2026-11-14): forgetting to switch the model, then research sweeps weighing on the main session, then the hand-written header that installation erases.

**Forgetting to switch.** A Supervisor session monitors Runs, checks Specs, opens pull requests and follows releases for hours, and in the same session authors an Idea, a PRD or a TechSpec. The orchestration wants an economical model and the authoring wants the strongest one. The switch is manual (`/model`, `/effort`), so either the authoring runs on the economical model or the orchestration runs on the expensive one. The maintainer's own workaround is a rule in his personal instructions ("Model by phase"); a rule still depends on being remembered, by the person or by the agent.

**Sweeps on the main session.** Codebase, Secondbrain and web sweeps fill the main session's context and spend the strong model's tokens on reading. `write-techspec` already says "Delegate the exploration to a subagent when one is available", but no adopting repository has one unless someone writes it by hand, and a hand-written one is easy to leave with a shell.

**The header is erased.** Indexus applied `model` and `effort` to four authoring skills and a research subagent by hand (2026-11-13), and Oraculum recorded a Backlog Entry to do the same. `roundfix skills install` and `roundfix baseline update` replace each owned `SKILL.md` whole, so the next Roundfix release erases the edit in every repository that made it.

**The header alone does not reach the main path.** Measured on Claude Code 2.1.325 (2026-11-14, see [the measurement](references/2026-11-14-skill-model-override-measurement.md)): a skill's `model` is applied when the user types `/skill-name` and ignored when the agent calls the skill. This is the open issue anthropics/claude-code#98898. In an agent-conducted cycle, which is how most Roundfix authoring happens, a header changes nothing today.

### Market Data

- Model by phase is common in coding tools and is nearly always set by the user: aider (architect, editor and weak models), Cline (separate Plan and Act models), Kilo Code (a model pinned per agent), Cursor (a `model` key per subagent file). Claude Code ships one built-in pairing, `opusplan`.
- The comparisons read for this study show no model-by-phase setting in spec-driven frameworks (Spec Kit, OpenSpec, BMAD). This is absence from comparisons, not a reading of each framework's configuration reference.
- Measured evidence exists for code editing, not for Spec authoring: aider's 2024 benchmark has a strong architect with a cheaper editor at 85.0% against 79.7% for the strong model alone; Anthropic reports a lead agent with cheaper subagents beating a single agent by 90.2% on its research evaluation, at about 15 times the tokens of a chat. No published measurement of a strong model against a cheap one for PRD or TechSpec authoring was found.
- Claude Code documents `model` and `effort` in skill frontmatter and `model`, `effort`, `tools` and `disallowedTools` in subagent files. OpenCode and Codex carry a model only in agent files, never per skill.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | Declared model per authoring phase | Critical | Project Config declares the model and reasoning effort for each authoring phase and for research sweeps. Roundfix ships no model name: a repository that declares nothing behaves exactly as today. |
| F2  | The agent-conducted phase runs on the declaration | Critical | When the agent conducts authoring, each phase is delegated to a generated author subagent that runs on the declared model and effort, because the skill header is ignored on that path. The decisions the phase needs are gathered from the user before the delegation. |
| F3  | The typed skill runs on the declaration | High | Installation writes the declared model and effort into the installed authoring skills, so a phase the user starts by typing the skill runs on them. Installation and Baseline updates keep the header, and no check reports it as a difference from the shipped skill. |
| F4  | Generated research subagent | High | A generated subagent runs on the declared research model with a tool list that holds no shell and no write tool. The authoring skills send codebase, Secondbrain and web sweeps to it and keep the synthesis and the decisions. |
| F5  | The same declaration for OpenCode and Codex | Medium | The declaration is compiled into the agent files those harnesses read, so the author and research subagents exist there too. Second Spec. |
| F6  | The user-level copy notice offers removal | Medium | The notice about a user-level copy of a Roundfix skill offers removing the copy beside refreshing it, and says that such a copy hides the repository's declared header. Reported by Pantheon on 2026-11-13. |

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Manual model or effort switches in a full authoring cycle | 0 | Count `/model` and `/effort` commands in the session record of one Idea-to-Tasks cycle in a pilot repository. |
| Authoring turns on the declared model | 100% on both paths | Read `message.model` of each assistant record in the session record, for a phase started by typing the skill and for a phase the agent conducts. |
| Research sweep tokens on the research model | > 80% | Usage by model in the session record of an authoring cycle, sweeps against the rest. |
| Files re-edited by hand after a Roundfix release | 0 | `git status` after `roundfix baseline update` in three adopting repositories that declared the setting. |
| Shell or write tools available to the research subagent | 0 | A test over the generated file's tool list, and one sweep in the pilot that tries a write. |

## Feature Assessment

| Criteria        | Question                                            | Score  |
| --------------- | --------------------------------------------------- | ------ |
| Impact          | How much more valuable does this make the product?  | Strong: removes a manual step that costs money or quality at every phase boundary. |
| Reach           | What % of users would this affect?                  | Strong: every adopter who authors through Claude Code in the first Spec, and through OpenCode or Codex in the second. |
| Frequency       | How often would users encounter this value?         | Strong: every authoring phase of every Spec. |
| Differentiation | Does this set us apart or just match competitors?   | Strong: no spec-driven framework in the comparisons read declares the model per phase in project configuration. |
| Defensibility   | Is this easy to copy or does it compound over time? | Maybe: easy to copy, and it rests on harness behavior that broke twice in four months. |
| Feasibility     | Can we actually build this?                         | Maybe: installation ignores Project Config and overwrites blindly, a new Project Config section is refused by an older binary, and Roundfix generates no agent file today. |

## Trade-off Insights

- **Recommended approach:** declare once in Project Config and enforce through the mechanism that is reliable on each path: the skill header for a typed skill, a generated subagent for an agent-conducted phase.
- **Key trade-offs:** generating from a declaration against preserving a hand-written header (generation needs a configuration surface and file management, preservation solves only the pain ranked last); a mechanism against a rule (a rule works on any harness and is forgotten, which is the problem being solved); shipped model names against an explicit declaration (names age and assume one vendor, so none is shipped); one Spec against two (two keep each within one context).
- **Risks identified:** the harness may change how it applies a skill's or a subagent's model again; a model identifier ages and differs per harness; an older Roundfix refuses a Project Config with a section it does not know; an author subagent cannot interview the user, so a phase that needs answers must collect them first; a research subagent reads untrusted text, so its tool list must exclude shell and write in fact, not in prose; this repository's own skills are the shipped bundle, so a declaration here must never reach the canonical files.
- **Dissenting view:** the operator recommended the hybrid in one Spec, without the OpenCode and Codex files, as the smallest shape that meets the success criterion; no adopter asked for authoring on those harnesses, and the maintainer had first scoped the audience to Claude Code. A second dissent: preserving a hand-written header is the cheapest change and adds no configuration surface.
- **Stretch goal (V2+):** report tokens by model for each authoring phase, so a repository sees what its declaration costs.

## Opportunity Scan

- **Original** — a `skills` section, the header at installation and a generated research subagent. Meets the success criterion only when the user types the skill. Medium effort.
- **Simpler** — installation preserves the header a repository wrote by hand; documentation shows the example. Small effort; Feasibility Must do, Impact Maybe, Differentiation Pass. Solves the pain ranked last.
- **Hybrid** — the original plus delegation of each agent-conducted phase to a subagent on the declared model. Medium effort; meets the criterion on both paths.
- **More ambitious** — the hybrid plus a generated author subagent and the same declaration compiled for OpenCode and Codex. Large effort, two Specs; Reach rises to Strong.

**Chosen direction:** more ambitious — the maintainer's pick on 2026-11-14, after the measurement showed the header alone fails the agent-conducted path.

## Sub-Features

- **0392, Claude Code:** F1, F2, F3, F4 and F6.
- **0394, OpenCode and Codex:** F5, after 0392 is delivered.

## Integration with Existing Features

| Integration point | How |
| ----------------- | --- |
| Project Config | Gains the declaration; User Config keeps only what belongs to a machine. |
| `roundfix skills install --target project` and `roundfix baseline update` | Apply the declaration when they install or refresh the authoring skills; user-level targets receive no header. |
| `roundfix skills check`, the Doctor Command and the Update Check | Read a declared header as part of the installation, and carry the new user-level copy notice. |
| The Roundfix skill's Skill Modes and the authoring skills | Say when a phase is delegated and to which generated subagent. |
| Agent Selection Profiles | Unchanged: they choose the Agent Session of a Run; this idea concerns the Supervisor's authoring. |

## Out of Scope (V1)

- **Model names shipped by Roundfix** — names age and assume one vendor; the documentation shows an example and the repository declares its own.
- **Changing how a Run selects its Agent** — Agent Selection Profiles already do that for Daemon-assigned work.
- **An author subagent that interviews the user** — a subagent cannot ask; the phase collects its decisions before delegating.
- **Repairing the harness** — the Claude Code defect is reported upstream; this idea works with the harness as measured.
- **Adopting the declaration in the Roundfix repository itself** — its skills are the shipped bundle; the maintainer decides after the release.
- **Cost reporting by phase** — kept as the stretch goal.

## Decisions

- Pain order: forgetting to switch the model, then sweeps on the main session, then the erased header (maintainer, 2026-11-14).
- Audience: anyone who authors through Claude Code in an adopting repository; nothing assumes the maintainer's machine (maintainer, 2026-11-14).
- Success: a full authoring cycle with every phase on its declared model and effort and no manual switch, the agent-conducted path included (maintainer, 2026-11-14).
- Shape: the more ambitious direction, in two Specs (maintainer, 2026-11-14).
- The user-level copy notice change (Pantheon) belongs to this Spec, not to one of its own (maintainer, 2026-11-14).
- Roundfix ships no model name (follows from the audience decision).

## Open Questions

- Does Claude Code apply a generated subagent's `model` and `effort` when the agent starts it? Default: yes, as documented; the PRD measures it with the same session-record reading before F2 is committed.
- Does an interactive session behave like the headless measurement for a typed skill and for an agent-called skill? Default: yes; measure once in the pilot.
- How does a repository on an older Roundfix behave once its Project Config gains the declaration? Default: the declaration requires the release that introduces it, and the older binary's refusal is named in the upgrade notes.
- Is the declaration made per phase or per skill? Default: per phase (idea, PRD, TechSpec, tasks, research), since a phase is what the user thinks in.
- Does Codex ignore frontmatter keys it does not know? Indexus cites the parser's tests; the official documentation read here does not say. Default: verify before F3 writes a header into a copy Codex also reads.
- Does the Roundfix repository adopt the declaration? Default: no, until the maintainer decides after the release.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-13-indexus-model-and-effort-per-skill-and-a-research-subagent.md`](references/2026-11-13-indexus-model-and-effort-per-skill-and-a-research-subagent.md), ai-memory message `01a1d68d-7d5c-7a22-8a37-85643843fe88`, copied verbatim; reported by indexus on 2026-11-13. It declared no blocked work and expects no reply.
2. [`references/2026-11-13-pantheon-roundfix-skills-left-the-user-level.md`](references/2026-11-13-pantheon-roundfix-skills-left-the-user-level.md), ai-memory message `01a1d6a1-2ce5-78d3-9998-e41bfdbd65e8`, copied verbatim; reported by pantheon on 2026-11-13. Its second request (the notice text) is F6; its first concerned the operator's own notes and is done. It declared no blocked work and asked for no notice.
3. [`references/2026-11-14-skill-model-override-measurement.md`](references/2026-11-14-skill-model-override-measurement.md), recorded by the operator.
4. [`references/2026-11-14-maintainer-answers.md`](references/2026-11-14-maintainer-answers.md), the maintainer's answers behind every line of `## Decisions`.

Secondbrain (read on 2026-11-14; `wiki/index.md`, then `qmd query`):

- `projects/oraculum/mirror/docs/backlog/2026-11-13-modelo-por-fase-nas-skills-de-autoria.md` — Oraculum plans the same hand-applied header and research subagent, and notes that a Roundfix installation erases them. It confirmed that the pain repeats in each repository.
- `wiki/concepts/custos-e-limites-de-agentes-de-codigo.md` — a subagent can inherit its parent's model and effort and multiply cost unnoticed; a strong model is justified for reusable assets and wasted on throughput. It supported declaring the research model instead of inheriting it.

External sources (read through Exa on 2026-11-14):

- Claude Code skills, https://code.claude.com/docs/en/skills — `model` and `effort` in skill frontmatter, applied for the current turn.
- Claude Code subagents, https://code.claude.com/docs/en/sub-agents — `model`, `effort`, `tools` and `disallowedTools` in a subagent file.
- Claude Code model configuration, https://code.claude.com/docs/en/model-config — aliases move over time; `opusplan`.
- anthropics/claude-code#98898, https://github.com/anthropics/claude-code/issues/98898 — the override is ignored when the agent calls the skill; open, confirmed by the measurement above. anthropics/claude-code#81618, closed 2026-10-23.
- OpenCode agents, https://opencode.ai/docs/agents/ — a model per agent file.
- Codex configuration reference, https://learn.chatgpt.com/docs/config-file/config-reference — model settings live in configuration and agent files, not in a skill.
- aider, https://aider.chat/2024/09/26/architect.html — the architect and editor benchmark quoted in Market Data.
- Anthropic, https://www.anthropic.com/research/multi-agent-research-system — the lead agent and subagents result quoted in Market Data.
- Cline, https://docs.cline.bot/core-workflows/plan-and-act; Cursor, https://cursor.com/docs/subagents; Kilo Code, https://kilo.ai/articles/roo-to-kilo-migration-guide — model by phase set by the user.
- Simon Willison, https://simonwillison.net/2025/jun/16/the-lethal-trifecta/ — why a subagent that reads untrusted text must not hold shell or write tools.
