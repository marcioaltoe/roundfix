---
spec: 0405-authoring-model-fallback
created: 2026-11-17
status: proposed
leverage: quick-win
---

# An Authoring Phase falls back to its declared fallback model when the declared model's limit runs out

> Authored on Opus 5.5, not on Fable, the model the repository declares for the Idea phase: Fable's quota for this account ran out on 2026-11-17 and resets on 2026-11-22 at 20:00 BRT. The same limit is what this Spec is about.

## Overview

A repository declares in Project Config the model and reasoning effort of each Authoring Phase under `skills.authoring.claude.<phase>` (ADR-0401). Each declared phase gets an Author Subagent. When the declared model's usage limit runs out partway through a phase, the Author Subagent ends on the error. The skill then tells the Supervisor to report the error and to author in the session only if the user says so. The phase stops until the maintainer answers.

This Idea lets the declaration name, for each phase, a fallback model and effort. When the Author Subagent ends on the usage limit of the declared model's family, the Supervisor notifies the maintainer, starts the phase's Author Subagent again on the fallback model with the brief and the partial output, and finishes the phase. The Spec's remaining phases run on the fallback. The next Spec starts again on the declared model. The Phase Record (Spec 0404) records the model the phase actually ran on and why.

The maintainer's goal (2026-11-17): "Se a cota do Fable acabar enquanto estiver trabalhando em alguma skill, automaticamente deve ir para o Opus High ou X-High para não perder qualidade." V1 is the minimal version: the declared fallback, the automatic restart with a notice, and the record.

## Problem

On 2026-11-17 the TechSpec Author Subagent of Spec 0404 was declared `fable` through `skills.authoring.claude.techspec` (Spec 0392, ADR-0401). Partway through the phase it ended with "You've hit your Fable limit · resets Nov 22 at 8pm" (HTTP 429). The Supervisor stopped, asked the maintainer, and restarted the work on an Opus agent by hand. This was the authoring skills' current rule working as written: "When it cannot start or ends on a model error, report the error and author here only if the user says so." In an unattended cycle, the stop lasts until the maintainer reads the message, and the restart is improvised: which model, which effort, and what to pass on from the partial work all go unrecorded.

Claude Code does not cover this case (docs read 2026-11-17, CLI 2.1.326):

- **Frontmatter:** neither skill nor subagent frontmatter has a fallback field. Both take only `model` and `effort`.
- **Native fallback chain:** `fallbackModel` / `--fallback-model` applies to subagents since v2.1.247. It switches only when the model is overloaded or unavailable, and "authentication, billing, rate-limit … errors … never trigger a switch" (model-config, Fallback model chains).
- **A limited subagent ends:** it fails with the error, and its partial output reaches the parent (sub-agents, API errors in subagents).
- **Resume keeps the model:** a resumed or messaged subagent keeps the model it was started with, so continuing the failed agent goes nowhere.
- **A new call can change it:** a new Agent call can pass a model alias and, since v2.1.292, an effort, both overriding the subagent file (sub-agents, Choose a model / Choose an effort level).
- **Family limits are escapable:** the Fable, Opus and Sonnet limits apply only to their own family, so switching family keeps work going. The session and weekly limits are shared across all models, and no fallback helps there (errors, You've hit your session limit).

The maintainer also changes accounts mid-week (on 2026-11-17 he planned to move to an account whose limits reset on Wednesday or Thursday). A fallback must therefore never become permanent: the next Spec has to try the declared model again.

## Core Features

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | Declared fallback | Critical | Each entry of `skills.authoring.claude.<phase>` may name a fallback model and reasoning effort. Roundfix ships no model name, as in ADR-0401. A fallback must be one model alias the harness accepts when starting a subagent (`sonnet`, `opus`, `haiku`, `fable`). |
| F2  | Automatic restart with notice | Critical | When the Author Subagent ends on the usage limit of the declared model's family, the Supervisor tells the maintainer which phase switched, from what to what and why. It then starts the same Author Subagent on the fallback model and effort, with the original brief, the partial output and the files already written, and the phase continues. |
| F3  | Trigger limited to the model's limit | Critical | Only a model-family usage limit (Fable, Opus or Sonnet limit, HTTP 429) triggers the switch. A session or weekly limit, which is shared by every model, stops and reports as today. Any other model error keeps today's rule. |
| F4  | Fallback holds for the Spec | High | Once a phase of a Spec fell back because a model's limit ran out, that Spec's later Authoring Phases whose declared model is that model start on their fallbacks. The next Spec starts on the declared models. |
| F5  | Phase Record names the model actually used | High | The Phase Record of Spec 0404 records the phase's actual model and effort, that it was a fallback, and the error that caused it. F4 reads this record to know a Spec already fell back. |
| F6  | This repository's declaration | Medium | Idea, PRD and TechSpec fall back to `opus`/`xhigh`. Tasks stays declared `opus`/`high` and falls back to `sonnet`/`high`. The repository declares them only after a released Roundfix knows the key. |

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| Phases stopped waiting for the maintainer after a model-family limit | 0 | Phase Records with a fallback reason and no maintainer intervention, across the cycles after release |
| Time from limit error to the phase continuing | < 2 minutes | Timestamps of the error notice and of the fallback start in the Phase Record |
| Fallback phases recorded with model, effort and reason | 100% | Every Phase Record whose model differs from the declaration carries a reason |
| Next Spec back on the declared model | 100% | First phase of the next Spec records the declared model when its limit is available |

## Feature Assessment

| Criteria        | Question | Score |
| --------------- | -------- | ----- |
| Impact          | How much more valuable does this make the product? | Strong: removes an unattended stop in the most expensive phases |
| Reach           | What % of users would this affect? | Strong: every adopter that declares an authoring model on Claude Code |
| Frequency       | How often would users encounter this value? | Maybe: weekly near the end of a quota window, never otherwise |
| Differentiation | Does this set us apart or just match competitors? | Maybe: Claude Code's own chain leaves rate limits out, and Roundfix fills that gap |
| Defensibility   | Is this easy to copy or does it compound over time? | Maybe: grows in value with Spec 0404's per-model statistics |
| Feasibility     | Can we actually build this? | Must do: the per-call model and effort override already exists and the declaration is already parsed |

Leverage: Quick Win.

## Trade-off Insights

- **Recommended approach:** a declared fallback per phase, plus the restart rule in the four authoring skills' "Declared model" section, plus the record in the Phase Record. No extra subagent file is generated, because a per-call model and effort override the subagent file.
- **Key trade-offs:** the per-call override takes model aliases only, so a fallback cannot be a full model ID (for example a pinned `claude-opus-5-5`). A generated `-fallback` subagent file would allow full IDs, but it doubles the generated files.
- **Risks identified:** limit errors are recognized by their text, which Claude Code can change. A fallback in the same family as the limited model does not escape the limit, and Roundfix knows no family by name; it can only compare the limit's name with the declared and fallback model names. The Agent tool asks for effort only when a skill or the user explicitly requests it, so the skills must state the rule plainly. The restarted agent starts fresh: anything the first agent decided but did not write down is lost. A fallback that the user starts by typing the skill (header route) is not covered, because the main session cannot switch its own model.
- **Dissenting view:** the persistence rule. Keeping the fallback for the rest of the Spec gives the Spec consistent authorship and avoids a second limit stop. But it keeps Opus even after the maintainer changes to an account with Fable available partway through the Spec, which trying the declared model at each phase would avoid. The maintainer chose consistency within a Spec (2026-11-17).
- **Stretch goal (V2+):** a fallback chain of more than one model; fallback for the Research Subagent; `roundfix doctor` warning about a fallback in the same family as its declared model.

## Opportunity Scan

- **More ambitious:** Roundfix watches account quotas (`/usage`, the status line's `rate_limits`) and switches models before the limit hits. This is medium effort, and Roundfix sees no quota API today. Impact Strong, Feasibility Maybe.
- **Simpler:** a fixed rule in the Baseline, "on a Fable limit, restart on Opus xhigh", with no configuration. This is small effort, but it puts a model name in Roundfix, which ADR-0401 rejects. Impact Strong, Defensibility Pass.
- **Adjacent:** the session and weekly limits: wait for the reset and continue (Claude Code's `autoContinueAtUsageLimit`). This is a different problem, already handled by the harness.

**Chosen direction:** original. The declared fallback keeps the decision of which model to use with the repository, as ADR-0401 requires. Its only enabler, the per-call model and effort override, is already in the installed harness (maintainer, 2026-11-17: "fallbacks na declaração").

## Out of Scope (V1)

- **Session and weekly limits:** they are shared across every model, so switching models does not help; they stop and report as today.
- **Overload and unavailable-model errors:** Claude Code's `fallbackModel` already covers them.
- **The header route,** where the user types the skill and it runs in the main session: the session cannot switch its own model, so the user's `/model` stays the remedy.
- **Runs' Agent Selection Profiles:** they already have preferred and fallback selections (ADR-0050). This Spec is about Authoring Phases.
- **Harnesses other than Claude Code:** the declaration supports only `claude` today.

## Idea-stage judgment

- item: the request in the session (scratch/spec-0405/request.md)
- judged-by: judge (typesafe/jev-1.13-20260917)
- shape_open: yes (0.90)
- unverified_assumption: no (0.24)
- market_question: no (0.09)
- verdict: indicated
- stage: ran
- shape: changed — the open "frontmatter or Baseline rule" became a declared fallback per phase with a restart rule in the authoring skills, because Claude Code's own fallback chain excludes rate limits and no frontmatter declares a fallback

## Decisions

All from the maintainer on 2026-11-17, about 17:50 BRT:

- **Shape:** fallbacks in the declaration ("fallbacks na declaração (Recommended)"), with the restart rule in the skills and no generated fallback subagent file.
- **This repository's fallbacks:** Idea, PRD and TechSpec fall back to Opus xhigh ("Opus xhigh / Tasks Opus high (Recommended)"). For Tasks, on approving the Idea (about 18:10 BRT): "Vamos manter opus high e sonnet high só para fallback ter algo." Tasks stays declared Opus high and falls back to Sonnet high.
- **No family warning:** Roundfix ships no model name, so it cannot tell a model's family. Instead the declaration refuses a fallback that names its declared model, and the switch is refused when the fallback's name contains the name of the limit that ran out.
- **Switch:** switch automatically and notify ("Troca e avisa (Recommended)").
- **Trigger:** only the model's own limit ("Só limite do modelo (Recommended)"). A session or weekly limit stops and reports.
- **Persistence:** "Fica no reserva na spec e volta ao modelo declarado na próxima spec".
- **Record:** the Phase Record records the model actually used and why ("Sim, no Phase Record (Recommended)"). This Spec therefore requires Spec 0404.
- **Delivery:** "Mesmas do BE (Recommended)": implement, commit, push, pull request, merge and release; up to 5 corrections; Makefile, CI, go.mod and lint only if a correction requires it.
- **Decision record:** the persistence rule and the restart rule will be recorded as ADRs from ADR-0427 onward in the PRD; they amend ADR-0401 and ADR-0408.

## Open Questions

- **Which later phases stay on the fallback.** The maintainer said the fallback holds for the rest of the Spec. Read literally, a Fable limit in the TechSpec would also move Tasks, declared Opus, to its Sonnet fallback, which he declared "só para fallback ter algo". Default: only later phases whose declared model is the model whose limit ran out start on their fallback; the others start on their declared model.
- **Rollout order.** The binary must learn the new key before this repository's config uses it (ADR-0408). Default: the declaration follows a release that knows the key.
