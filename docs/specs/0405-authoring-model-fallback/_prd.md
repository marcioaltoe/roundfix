---
spec: 0405-authoring-model-fallback
status: active
created: 2026-11-17
surfaces: [backend, cli, docs]
type: feat
---

# An Authoring Phase falls back to its declared fallback when the declared model's usage limit runs out

When a phase's Author Subagent hits the usage limit of its declared model partway through, the phase stops today until the maintainer answers. On 2026-11-17 the TechSpec Author Subagent of Spec 0404 ended on "You've hit your Fable limit · resets Nov 22 at 8pm", and the Supervisor had to ask before restarting the work on Opus by hand. After this Spec, a repository declares a fallback model and effort for each Authoring Phase. When the declared model's own limit ends the Author Subagent, the session tells the maintainer and starts the same Author Subagent again on the fallback with the brief and the work already done. The Spec's later phases on that model stay on their fallbacks, the next Spec starts on the declared models, and the Phase Record says which model ran and why.

## Sources

1. [`references/2026-11-17-maintainer-answers.md`](references/2026-11-17-maintainer-answers.md), the request, the Idea-stage judgment and the maintainer's answers, recorded by the author.
2. [`references/2026-11-17-harness-and-test-measurements.md`](references/2026-11-17-harness-and-test-measurements.md), recorded by the author: what the Claude Code documentation and the installed CLI say about fallbacks and usage limits, the Governed Paths the design changes, and the existing tests a prototype of its surfaces moves.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from. This project's ai-memory inbox was listed on 2026-11-17 and held no pending message, and `docs/findings/` holds no Finding. None of the open Backlog Entries shares this Spec's context. The nearest one is about the start cost of a session, which is a matter of tokens, not of limits.

## Prerequisites

Specs 0401, 0402, 0403 and 0404 are delivered before this one. The Tasks stage names them in `requires`. This Spec was authored against `a8236a87`, the delivery of Spec 0400; none of the four was delivered then, and Spec 0404's plan was not yet on the default branch.

- Spec 0404 adds to a phase entry the model and effort that ran an Authoring Phase, and the flags the session uses to report them. This Spec's record of a fallback stands beside that entry, and its persistence reads the repository identity that Spec 0404 adds.
- Spec 0401 changes the phase command and the phase section of the four authoring skills.
- Specs 0402 and 0403 change the Roundfix skill, its references, the glossary, the surface record and the configuration guide.

Each of these files is also changed by this Spec, inside its own sentences. No version number is written by hand: the record commands raise it on the tree they run on.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the fallback record carries the phase, the Spec slug and the repository identity other components define. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added, and no judge or provider is asked. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0427 (new) records the declared fallback and its trigger (ADR-0427: "Only a model's own usage limit switches"). It amends ADR-0401, under which the session starts the Author Subagent "passing no model and no effort" (ADR-0401: "passing no model and no effort"), and keeps its rule that the repository names every model (ADR-0401: "Roundfix ships no model name"). ADR-0428 (new) records how long a fallback holds (ADR-0428: "the next Spec starts on the declared models"). The fallback is a new kind of entry in the Phase Record, which keeps its other rules (ADR-0407: "It never holds prompt text, a token count, a price or an environment value"), and the phase command still records a phase as it does after Spec 0401 (ADR-0414: "records each phase as before and reads no earlier entry"). A new Project Config key cannot reach this repository in the change that teaches the binary that key (ADR-0408: "refuses a key it does not know"), and an older binary refuses a repository a newer one wrote (ADR-0392: "every command that loads the repository's configuration"). The fallbacks of a Run are unchanged (ADR-0050: "once Agent work begins, Roundfix fails the Work Item instead of switching models over potentially modified state"). Each changed surface ships with its skill text (ADR-0256: "The Roundfix-owned skills must describe the behavior Roundfix ships"). The Release Ceiling in the authorization record answers the release plan (ADR-0349: "the Spec's authorization record carries the highest impact its release grant answers"). The TechSpec added ADR-0429 (ADR-0429: "A refusal exits `3`"). Not applicable: ADR-0405, which cites ADR-0401 for a user-level skill copy; ADR-0410, which cites ADR-0407 for the status view, which keeps reading phase entries only; ADR-0412, which cites ADR-0349 for the release impact of planning references; and ADR-0415, which cites ADR-0401 for the recommended declaration, which gains no fallback in this Spec. Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-17, quoted in the authorization record, which keeps the grant of the owned skills. It grants no Makefile, workflow, Go module, lint or formatter change except for a correction that requires one. This repository's Project Config declares its fallbacks in the operator's step after the release, and no Task makes it. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0405-authoring-model-fallback/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- A phase whose declared model's usage limit runs out finishes on its declared fallback with no question to the maintainer.
- The maintainer is told every time a phase switches, from which model to which and why.
- The fallback holds for the rest of the Spec on the model that ran out, and never past the Spec.
- The Phase Record says which model ran each phase and why it was not the declared one.
- A repository that declares no fallback authors exactly as it does today.

## User Stories

1. As the maintainer, I want to declare a fallback model and effort for each Authoring Phase, so that the phase keeps its quality when the declared model's limit runs out.
2. As a Supervisor whose Author Subagent ended on its model's usage limit, I want one command to tell me whether to fall back and to what, so that I restart it without asking.
3. As the maintainer, I want to be told which phase switched, from which model to which and why, so that I know which model wrote which artifact.
4. As a Supervisor starting the next phase of a Spec that already fell back, I want to know to start it on its fallback, so that I do not lose a start to a limit that has not reset.
5. As the maintainer, I want the next Spec to start on the declared models, so that a new account or a reset limit is used again.
6. As the maintainer, I want a session or weekly limit, or any other error, to stop the phase as today, so that a fallback is never used when it cannot help.
7. As an adopter that declares no fallback, I want my authoring unchanged, so that an update costs me nothing.

## Core Features

1. **A phase declares a fallback.** An entry `idea`, `prd`, `techspec` or `tasks` of the **Authoring Model Declaration** may hold a fallback with a model and, optionally, a reasoning effort. The research entry takes none. A fallback that names the entry's own declared model is refused, because it escapes no limit of that model. A refusal names the key and the accepted values. The configuration guide says the model must be a name the harness accepts when a subagent is started, an alias in Claude Code, and that a Roundfix older than the release refuses a Project Config that holds the key.
2. **Roundfix decides the fallback from the limit message.** One command takes the phase, the Spec and the error the Author Subagent ended on. It applies an **Authoring Fallback** only when the error says a named model's usage limit ran out, the name is that of the phase's declared model, a fallback is declared and its model is not the limited one. It then records the fallback and prints the model and effort to start the Author Subagent on. Otherwise it records nothing and prints why the phase stops: a session or weekly limit is shared by every model, a spend limit is not a model's limit, the error is not a limit, the limit is not the declared model's, no fallback is declared, or the fallback is the limited model. Hazard: the error is text a harness wrote; only the line that names the limit is kept, cut to 200 characters on one line, and it is printed and recorded as data.
3. **The session switches and tells.** The four authoring skills tell the session, when the Author Subagent ends on an error, to run that command. When it prints a fallback, the session tells the user which phase switched, from which model to which and the limit that ran out, and, without waiting, starts the same Author Subagent with the printed model and effort, the original brief, the output the first run returned and the files it wrote. When the command refuses, the session reports the error and authors in the session only if the user says so, as today. When the installed Roundfix does not know the command, the skill behaves as today and says so.
4. **The fallback holds for the Spec.** A second command prints the model and effort a phase's Author Subagent starts on. It names the fallback when an earlier phase of the same Spec, in the same repository, fell back because the limit of this phase's declared model ran out and this phase declares a fallback. Otherwise it says to start on the declared model, passing no model and no effort. The authoring skills tell the session to run it before starting an Author Subagent. The next Spec starts on its declared models.
5. **The Phase Record holds the fallback.** The **Phase Record** gains a fallback entry: when it was recorded, the phase, the Spec, the repository, the session, the declared model and effort, the fallback model and effort, the name of the limit and the limit line. Readers of phase entries skip it, and the Run Retention Sweep removes it with its day file. The restarted Author Subagent records its phase entry with the model it runs on, as Spec 0404 makes it.
6. **One step, and only for a started subagent.** A fallback has no fallback: when the fallback's own limit runs out, the command refuses and the phase stops. A phase the user starts by typing the skill runs in the session, which cannot change its own model, and keeps the harness's remedy.
7. **Documentation and glossary.** The user guide describes the key and the two commands, and the Roundfix skill's phase reference describes them for a Supervisor. `CONTEXT.md` adds **Authoring Fallback** and revises **Authoring Model Declaration**, **Author Subagent** and **Phase Record**. ADR-0401 carries a notice naming ADR-0427, and ADR-0407 one naming ADR-0428.

## User Experience

The lines below show the content of each surface. The TechSpec fixes the command and flag names and the exact wording.

A Project Config declares the fallbacks:

```yaml
skills:
  authoring:
    claude:
      techspec:
        model: fable
        reasoning_effort: xhigh
        fallback:
          model: opus
          reasoning_effort: xhigh
```

The TechSpec Author Subagent ends on its model's limit, and the session asks Roundfix:

```text
$ roundfix phase fallback techspec --spec 0412-example --error "You've hit your Fable limit · resets Nov 22 at 8pm"
Start roundfix-techspec-author with model opus and effort xhigh, the fallback of fable: "You've hit your Fable limit · resets Nov 22 at 8pm"
```

The session tells the user, in its own words: the TechSpec of Spec 0412 switched from fable xhigh to opus xhigh because the Fable limit ran out, and it goes on.

The weekly limit ends it instead:

```text
$ roundfix phase fallback techspec --spec 0412-example --error "You've hit your weekly limit · resets Mon 12:00am"
No fallback for techspec of 0412-example: a weekly limit is shared by every model. Report the error.
```

In a Spec whose Idea fell back on the Fable limit, the next phase declared on Fable starts on its fallback, and a phase declared on Opus starts as declared:

```text
$ roundfix phase model prd --spec 0412-example
Start roundfix-prd-author with model opus and effort xhigh, the fallback of fable: idea of 0412-example fell back when the Fable limit ran out.
$ roundfix phase model tasks --spec 0412-example
Start roundfix-tasks-author with no model and no effort: the declared model opus.
```

## Non-Goals / Out of Scope

- Session, weekly and spend limits: they are shared across models or are not a model's limit, and they stop the phase as today.
- Overload and an unavailable model: the harness's own fallback chain covers them.
- A phase the user starts by typing the skill, which runs in the session.
- The fallbacks of a Run's Agent Selection Profiles.
- A chain of more than one fallback, and a fallback for the Research Subagent.
- Harnesses other than Claude Code: the declaration supports only `claude`.
- The recommended declaration of the Doctor Command and the Update Check: it gains no fallback in this Spec.
- This repository's Project Config: the operator declares its fallbacks after the release.

## Success Metrics

1. Success Metric: after the release, every Authoring Phase whose Author Subagent ends on its declared model's own usage limit, with a fallback declared, finishes on the fallback without a question to the maintainer.
2. Success Metric: every fallback leaves one fallback entry in the Phase Record with the limit line, and one phase entry of the restarted Author Subagent.
3. Success Metric: the first phase of the next Spec starts on its declared model.
4. Success Metric: a repository that declares no fallback sees the same commands, outputs and skill behavior as before.

## Idea-stage judgment

- item: the request in the session (scratch/spec-0405/request.md)
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: yes (0.90)
- unverified_assumption: no (0.24)
- market_question: no (0.09)
- verdict: indicated
- stage: ran
- shape: changed — the open "frontmatter or Baseline rule" became a declared fallback per phase with a restart rule in the authoring skills, because Claude Code's own fallback chain excludes rate limits and no frontmatter declares a fallback

## Glossary

- adds: **Authoring Fallback**
- changes: **Authoring Model Declaration**
- changes: **Author Subagent**
- changes: **Phase Record**

## Decisions

- 2026-11-17: Where the fallback is declared — "fallbacks na declaração (Recommended)": a fallback per phase in the Authoring Model Declaration and a restart rule in the authoring skills, with no generated fallback subagent. See ADR-0427.
- 2026-11-17: This repository's fallbacks — "Opus xhigh / Tasks Opus high (Recommended)": the Idea, the PRD and the TechSpec fall back to Opus xhigh.
- 2026-11-17: The Tasks' fallback — "Vamos manter opus high e sonnet high só para fallback ter algo.": the Tasks stay declared Opus high and fall back to Sonnet high.
- 2026-11-17: Ask or act — "Troca e avisa (Recommended)": switch and tell the maintainer. See ADR-0427.
- 2026-11-17: What triggers the switch — "Só limite do modelo (Recommended)": only the declared model's own usage limit; a session or weekly limit stops and reports. See ADR-0427.
- 2026-11-17: How long the fallback holds — "Fica no reserva na spec e volta ao modelo declarado na próxima spec". See ADR-0428.
- 2026-11-17: Where the switch is recorded — "Sim, no Phase Record (Recommended)". See ADR-0428.
- 2026-11-17: The Idea draft — "Aprovar e seguir (Recommended)".
- 2026-11-17: Delivery authorization — "Mesmas do BE (Recommended)".
- Autonomous route: no person present; decisions taken from the maintainer's answers above and from the harness measurements in the references. The author settled three points from them: only later phases on the model that ran out stay on their fallback, so the Tasks, declared Opus, are not moved to Sonnet by a Fable limit; Roundfix compares names instead of knowing model families, refusing a fallback that names its declared model; and the recommended declaration of the Doctor Command gains no fallback in this Spec.
- 2026-11-17: Authoring model of this Spec — every phase was authored on Opus 5.5, `claude-opus-5-5`, not on the declared Fable, whose quota for the account ran out on 2026-11-17 and resets on 2026-11-22 at 20:00 BRT.

## Open Questions

- Which later phases stay on the fallback: default, only those whose declared model is the model that ran out (ADR-0428). The maintainer's literal answer would move every later phase.
- Whether the recommended declaration should carry fallbacks: default, no, until the maintainer asks.
