---
spec: 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
status: active
created: 2026-11-18
surfaces: [backend, cli, docs]
type: fix
---

# Adopting the recommendations works under load and gives the researcher the web

Two defects that adopters met while adopting the recommendations of v0.148.0. First, `profiles configure` and `profiles check --apply` prove every tuple they write under the built-in 30-second setup deadline and ignore `defaults.profile_proof_setup_timeout`, although the configuration guide and the Roundfix skill say `profiles configure` honors it; on a loaded machine Fiscus and Fluxus could not adopt a Recommended Profile and edited `.roundfixrc.yml` by hand. Second, the Recommended Authoring Declaration declares no research tools, so the `roundfix-researcher` subagent it produces holds only `Read, Grep, Glob` while its own text says it reads the web, and nothing that Doctor or the Update Check prints names the `research.tools` key. After this Spec both commands prove under the configured deadline, the recommendation gives the Research Subagent the built-in tools `WebSearch` and `WebFetch`, and a repository whose research entry declares no tools is told which ones to add.

## Sources

This Spec has no Backlog Entry; the operator minted it from the session of 2026-11-18. The maintainer's answers and both field reports are in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md): the first defect was reported by fiscus and fluxus in their sessions of 2026-11-18, the second was reported by conexus, message 01a1f18f. The code facts were measured by the operator and measured again by the author in [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md).

The two defects share no component: one is the proof path of the two commands that write Agent Selection Profiles, the other is the Recommended Authoring Declaration and the two checks that print it. The maintainer chose one Spec for both on 2026-11-18 (`## Decisions`): both were found by adopters taking the recommendations of the same release, they ship in the same patch release, and the second alone would be a one-Task Spec paying a whole Run, QA gate, review and release. No open Backlog Entry or unresolved Finding concerns the proof deadline or the research tools.

## Prerequisites

None. Authored against `a6a4a15b` (v0.148.0 with the plans of 0407 and 0408). Specs 0407 and 0408 deliver before this one and also change `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `skills/testdata/owned-skill-versions.json` and the Roundfix skill (0408: its `deliver` reference and version); no Verification of this Spec depends on their work, and each Task of this Spec regenerates the shared records on the tree it runs on.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the change reads one configuration key and one declaration entry. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; the proof runner and the checks read local files, and the web tools are names written into a generated subagent file, never called by Roundfix. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0441 (new) records the research tools decision and amends ADR-0415; ADR-0298 is the decision the proof fix restores; ADR-0402 keeps the generator's fixed list, ADR-0402: "The generated file's `tools` holds `Read`, `Grep` and `Glob`, then each name of `skills.authoring.<harness>.research.tools` in the declared order."; ADR-0181, ADR-0344, ADR-0401 and ADR-0430 are listed only to record why they do not apply; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE" of 2026-11-18 for Spec 0411, quoted in the authorization record. It covers the Roundfix skill's `SKILL.md` in both copies and its canonical `profiles.md`, `setup.md` and `modes.md` references, and grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/profiles.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/modes.md`.

## Goals

- `profiles configure` and `profiles check --apply` prove every tuple under the `defaults.profile_proof_setup_timeout` that User Config and Project Config declare, Project Config winning, as `profiles validate` does, so a slow machine can adopt a Recommended Profile by raising that key.
- An adopter that pastes the Recommended Authoring Declaration gets a Research Subagent that can search and read the web.
- An adopter that pasted the earlier block, or declared a research entry without tools, is told by Doctor and the Update Check which tools to add and how to apply them; nothing is rewritten for it.
- Every guide and skill reference that names the deadline's consumers or the research tools says what the code does.

## Core Features

1. `profiles configure` and `profiles check --apply` read the configured proof setup deadline from User Config and Project Config before they prove, even when another key of those files would stop the configuration from loading, and refuse an invalid value before any proof. A file that does not parse as YAML refuses the command with a message that names the file and the line of the YAML error.
2. The Recommended Authoring Declaration lists `tools: [WebSearch, WebFetch]` under `research`, and no MCP tool; the generated subagent's fixed tools stay `Read, Grep, Glob`, and no shell or write tool can be added.
3. Doctor (`authoring-models: warn`) and the Update Check (a text line and `authoringModels.recommendedResearchTools`) name the recommended research tools when the declared research entry has no `tools` key; an empty list, `tools: []`, is the repository's choice and is not reported.
4. The configuration guide, the `profiles`, `doctor` and `update` guides, the Roundfix skill's `profiles`, `setup` and `modes` references, the glossary entry for the Recommended Authoring Declaration, ADR-0415 and the new ADR-0441 state both rules.

## Non-Goals / Out of Scope

- The fallback candidate probe of `internal/agent/fallback.go`, measured and out of scope: no command reaches it, because commit `d7fdd8cf` (#32) removed the last call of `probeRuntimeSelection`, so its built-in deadline affects no user. Removing it, with the unused Agent probe of the health checker, is the Backlog Entry [`2026-11-19-the-fallback-probe-and-the-health-checker-agent-probe-are-dead-code.md`](../../backlog/2026-11-19-the-fallback-probe-and-the-health-checker-agent-probe-are-dead-code.md).
- The Baseline semantic analysis proof of `internal/baselineacp`: the `baseline` command reads no Roundfix configuration, and an analysis that times out preserves every Source Baseline Entry (ADR-0344). It is the Backlog Entry [`2026-11-19-the-baseline-semantic-analysis-proves-under-a-fixed-30-second-deadline.md`](../../backlog/2026-11-19-the-baseline-semantic-analysis-proves-under-a-fixed-30-second-deadline.md).
- Reading any other key of the configuration in `profiles configure`, such as access policy or runtime commands.
- Adding web tools to the generated subagent's fixed list, adding MCP tool names to the recommendation, or rewriting the generated subagent's text.
- Comparing a declared research tool list with the recommendation when the list exists, and reporting model or effort differences of the authoring declaration.
- A snapshot date for the Recommended Authoring Declaration; the profile snapshot `2026-11-17` stays.

## Success Metrics

1. Success Metric: with a fake runner, every proof of `profiles configure` and of `profiles check --apply` carries the deadline that User Config and Project Config declare, Project Config winning, and the built-in 30 s when neither declares one, for a write to either scope.
2. Success Metric: `profiles configure` refuses a configuration whose deadline is invalid with the same message `profiles show` prints, exit 2, before any proof and without writing, and refuses a configuration file that does not parse as YAML with a message naming the file and the line of the error; it still rewrites a profile when an untouched category of the same file is invalid.
3. Success Metric: the rendered Recommended Authoring Declaration ends its research entry with `tools: [WebSearch, WebFetch]`, the subagent generated from it lists `Read, Grep, Glob, WebSearch, WebFetch` and no shell or write tool, and the configuration guide's block equals the rendering.
4. Success Metric: Doctor and the Update Check name the recommended tools for a research entry without a `tools` key, in text and JSON, and stay silent for `tools: []`, for the web tools and for a list of MCP tools; neither changes an exit code or the Update Check's `next`.
5. Success Metric: this repository's own declaration lists every recommended research tool, held by the repository contract test.
6. Success Metric, outside evidence: the field reports of fiscus and fluxus, which this Spec did not design, show the adoption failing under load while `profiles validate` passed, and the published Claude Code documentation states that a subagent whose `tools` field is set gets only the listed tools.

## Acceptance evidence

The outside-evidence row rests on two sources this Spec did not write. The field reports are recorded in [`references/2026-11-18-maintainer-answers.md`](references/2026-11-18-maintainer-answers.md): Fiscus measured each Codex proof at about 45 s under a load of 13 to 23 and proved by hand with `profiles validate`, and Fluxus raised the deadline to `3m` and still saw `profiles check --apply` time out. The published Claude Code tools reference, <https://code.claude.com/docs/en/tools>, says of a subagent definition: "`tools` only: the subagent gets only the listed tools."; the planning change commits that passage in [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md) section 5, so the QA gate checks the row without the network.

## Unreachable Acceptance

- criterion: on a loaded machine, an adopter adopts a Recommended Profile with `profiles check --apply` after raising `defaults.profile_proof_setup_timeout`
  reason: the QA gate may not start an Agent session or an ACP adapter, so it cannot time a real proof; the fake runner of Success Metric 1 records the deadline each proof receives
  satisfied-by: the next adoption by fiscus or fluxus after the release that carries this Spec

## Validation

The decisions were checked on 2026-11-18 against the code of `internal/cli/profiles_configure.go`, `profiles_check_apply.go`, `profiles_validate.go`, `internal/agent/fallback.go`, `internal/baselineacp/analyzer.go`, `internal/config/authoring_recommendation.go`, `internal/cli/authoring_reports.go` and `skills/authoring.go`, ADR-0298, ADR-0344, ADR-0402 and ADR-0415, and two prototypes on a scratch worktree, all recorded in `references/2026-11-18-characterization.md`. The Secondbrain query `qmd query "roundfix profile proof setup timeout configure"` returned only mirrors of this repository's own ADR-0298, ADR-0181 and guides, and `qmd query "roundfix-researcher subagent web tools"` returned this repository's ADR-0402 and Vortex's reference `projects/vortex/mirror/docs/references/2026-11-14-autorizacao-pins-de-modelo-researcher-e-lock-v2.md`, which created a read-only researcher by hand; neither changed the decisions, and ADR-0181's rule that a deliberate choice is not argued with shaped the silence for `tools: []`. Exa returned the Claude Code subagent guide and tools reference quoted above, which confirm that the listed tools are the boundary.

## Glossary

- changes: **Recommended Authoring Declaration**

## Decisions

- 2026-11-18: The defect is confirmed (`profiles configure` and `profiles check --apply` ignore the deadline; the fallback probe uses the constant); how is it handled? — "Spec fix 0411 após a 0408 (Recommended)": TechSpec on Opus, no Idea stage, "Mesmas do BE", patch release, queued right after 0408 and before 0409.
- 2026-11-18, about 23:35 BRT: Where does the Conexus report on the researcher's tools go? — "Juntar à Spec 0411 (Recommended)": one Spec for both defects, the same patch release, right after 0408.
- 2026-11-18: Is the workaround sent to Conexus? — "Mandar o contorno (Recommended)"; the operator sends it, not this Spec.
- 2026-11-19, about 00:05 BRT: Which release does the Spec take, now that it adds an optional JSON member and a report line? — "Patch, como autorizado (Recommended)".
- 2026-11-19, about 00:05 BRT: Do Doctor and the Update Check warn adopters whose research entry has no `tools`? — "Sim, avisar com o comando (Recommended)": nothing is rewritten by itself, and `tools: []` silences it on purpose.
- 2026-11-19, about 00:05 BRT: `profiles configure` now refuses when the other scope's file has a YAML syntax error; is that accepted? — "Aceitar, com mensagem clara (Recommended)": the message names the file and the line.
- 2026-11-19, about 00:05 BRT: Which findings of the planning become Backlog Entries? — "(c) bug do spec check, (a) remover código morto, (b) prazo da análise do baseline": all three, written by this planning change as [`2026-11-19-spec-check-aborts-on-a-verification-operand-longer-than-255-bytes.md`](../../backlog/2026-11-19-spec-check-aborts-on-a-verification-operand-longer-than-255-bytes.md), [`2026-11-19-the-fallback-probe-and-the-health-checker-agent-probe-are-dead-code.md`](../../backlog/2026-11-19-the-fallback-probe-and-the-health-checker-agent-probe-are-dead-code.md) and [`2026-11-19-the-baseline-semantic-analysis-proves-under-a-fixed-30-second-deadline.md`](../../backlog/2026-11-19-the-baseline-semantic-analysis-proves-under-a-fixed-30-second-deadline.md).

Autonomous route: no person present; decisions taken from the operator's brief for Spec 0411 of 2026-11-18, the coordinator's scope addition of about 23:35 BRT, the maintainer's answers above and the characterization.

- The two defects stay in one Spec, by the maintainer's answer, although they share no component.
- The proof fix reads only `defaults.profile_proof_setup_timeout`, because `profiles configure` must keep working on a file that another invalid key would stop from loading (characterization, section 3).
- The unreachable fallback probe and the Baseline semantic analysis keep their built-in deadline; each is a Backlog Entry, by the maintainer's answer of 2026-11-19.
- The recommendation adds only the built-in `WebSearch` and `WebFetch`; MCP names differ per machine (ADR-0441).
- The generator's fixed list stays `Read, Grep, Glob`, so a repository that declares no web tool keeps the subagent off the network (ADR-0402 holds).
- A research entry without a `tools` key is reported, never rewritten; `tools: []` and any declared list are the repository's choice (ADR-0441).
- `type: fix` with a `patch` Release Ceiling, by the maintainer's answers of 2026-11-18 and 2026-11-19; the new JSON member and output line are optional additions that report what the fix recommends.

## Open Questions

None.
