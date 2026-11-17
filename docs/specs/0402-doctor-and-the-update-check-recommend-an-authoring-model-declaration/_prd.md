---
spec: 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
status: active
created: 2026-11-17
surfaces: [backend, cli, docs]
type: fix
---

# Doctor and the Update Check recommend an Authoring Model Declaration

A repository can have the four authoring skills installed and declare no `skills.authoring.claude` in `.roundfixrc.yml`. Its Idea, PRD, TechSpec and Task Graph phases then run on whatever model the session uses, with no Declared Skill Header and no Author Subagent, and nothing says so: `roundfix doctor` prints its `authoring-models:` line and `roundfix update check` its `authoringModels` only when something is declared or a generated file remains (`internal/cli/authoring_reports.go`, `loadedAuthoringState`). Argus is in that state on 2026-11-17. After this Spec, both checks warn in that case and print the Recommended Authoring Declaration, the block to add to `.roundfixrc.yml`, with the models Roundfix itself declares. They never fail and never write.

## Sources

No Backlog Entry or Finding records this defect. The maintainer reported it on 2026-11-17 and chose the fix; the answers are in [`references/2026-11-17-maintainer-answers.md`](references/2026-11-17-maintainer-answers.md). The state of the adopters is in [`references/2026-11-17-installed-authoring-skills-without-a-declaration-measurement.md`](references/2026-11-17-installed-authoring-skills-without-a-declaration-measurement.md).

No open Backlog Entry and no unresolved Finding shares this context:

- The large-context entry changes how an authoring session starts, not which model it declares.
- The Delivery Queue, docs-folders and Conexus status entries change other components.

## Prerequisites

None in `requires`. Authored against `bcd716df` (`docs/plan-0402`, `v0.147.1`). Spec 0401 is planned and changes the phase-boundary rules in the authoring skills and the Roundfix skill's `modes.md`; this Spec changes the Roundfix skill's `setup.md` and one sentence of the `--update` mode in `modes.md`. Both raise the Roundfix skill's version, so whichever Spec delivers second records its version on the merge base, as the version record command does, and its Task rebases the `modes.md` sentence on the text it finds.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The new JSON state and fields are named in the TechSpec. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — both checks stay local and read files of the checkout; the Update Check's release lookup is unchanged. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0415 (new) records the decision and amends ADR-0401; ADR-0408 (published skills stay exempt) and ADR-0402 (the Research Subagent's tools) bind it; ADR-0405 and ADR-0413, which cite ADR-0401, are listed as not applying; the TechSpec relates each one. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-17, quoted in the authorization record. It covers the Roundfix skill's setup and modes references, their `SKILL.md` pair and the governed tests the record lists. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Goals

- A repository with an installed authoring skill it does not publish and no `skills.authoring.claude` gets `authoring-models: warn` from Doctor and `Authoring models:` with `authoringModels.state` `undeclared` from the Update Check.
- Both print the Recommended Authoring Declaration as a block ready for `.roundfixrc.yml`, and name `roundfix skills install --target project` as the step after it.
- Neither check fails, changes its exit code or the Update Check's `next`, or writes a file.
- The block's values have one source in the binary, and this repository's own declaration and the configuration guide cannot drift from it unnoticed.
- Every case reported today keeps its bytes: a declaration, a leftover generated file, a repository that publishes every installed authoring skill, no installed authoring skill, Project Config not read, and outside a Git repository.

## Core Features

1. The Recommended Authoring Declaration in `internal/config`: the model and reasoning effort of each phase and of `research`, and its rendering as YAML, with no `research.tools`.
2. The `undeclared` report in Doctor and the Update Check, text and JSON.
3. A repository contract test that holds this repository's `.roundfixrc.yml` and the configuration guide's block equal to the recommendation.
4. The guides, the Roundfix skill's setup and `--update` mode text, the glossary and ADR-0401's amendment line describe the report.

## User Experience

Today, in Argus:

```text
$ roundfix update check --format json | jq -c .authoringModels
null
```

With this Spec, `roundfix doctor` there prints, among its lines:

```text
authoring-models: warn (the authoring skills are installed without an Authoring Model Declaration: write-idea, write-prd, write-techspec, write-tasks; next: add the Recommended Authoring Declaration below to .roundfixrc.yml, then run roundfix skills install --target project)
  skills:
    authoring:
      claude:
        idea:
          model: fable
          reasoning_effort: xhigh
        ...
```

## Non-Goals / Out of Scope

- Writing `.roundfixrc.yml`, installing skills, or failing any check.
- A declaration that covers some phases and not others: a declared section stays the repository's choice and is reported as today.
- `roundfix skills check`, `roundfix skills install` and `roundfix baseline update`.
- A key that silences the warning, or a harness other than `claude`.
- Recommending `research.tools`.
- Changing this repository's `.roundfixrc.yml`, the Recommended Profile or the Model Catalog.

## Success Metrics

1. Success Metric: in a fixture repository with the four authoring skills installed and no declaration, Doctor prints the warn line and the block, and exits as it did; the Update Check's text and JSON carry the `undeclared` report with `installed`, `recommendation`, `problems: []`, no `next` member, and the report's `next` unchanged.
2. Success Metric: the same fixture with every installed authoring skill published, with no installed authoring skill, with a declaration, or with a leftover generated subagent produces the bytes it produced before this Spec.
3. Success Metric: the rendered block, loaded as Project Config, validates and yields exactly the recommended entries; it holds no `tools` key.
4. Success Metric: the repository contract test fails when this repository's `.roundfixrc.yml` changes one model or effort, or when the configuration guide's block changes one line.
5. Success Metric, outside evidence: the state the measurement recorded for Argus on 2026-11-17 (four installed authoring skills, no `skills:` key, no `.claude/agents`, `authoringModels` absent) reproduced in a fixture makes the built binary's Update Check report `undeclared` with those four skills, where the installed 0.147.1 printed no `authoringModels`.

## Acceptance evidence

The outside-evidence row rests on the measurement of the Argus checkout, a repository this Spec did not build, recorded in `references/` with its HEAD and the 0.147.1 output. The QA gate reproduces that state in a temporary fixture and never reads Argus.

## Unreachable Acceptance

None. The fixture is local, and the QA gate needs no network.

## Validation

The decisions were checked on 2026-11-17 against these sources:

- The measurement in `references/`.
- The code of `internal/cli/authoring_reports.go`, `authoring_models.go`, `doctor.go` and `update_check.go`, of `skills/authoring_inspect.go` and `skills/authoring.go`, and of `internal/config/authoring.go` and `recommendations.go`.
- ADR-0401, ADR-0402 and ADR-0408.
- A probe on a scratch worktree of `bcd716df` that made the checks report every installed authoring skill without a declaration and ran the whole `./internal/cli` package, to find the tests whose expectations this Spec changes.

## Glossary

- adds: **Recommended Authoring Declaration**
- changes: **Authoring Model Declaration**

## Decisions

- 2026-11-17: What happens when the authoring skills are installed and `skills.authoring.claude` is not declared? — "Avisar e sugerir (Recommended)": Doctor and the Update Check warn in that case and show the recommended block to paste into `.roundfixrc.yml`; they never fail and never write; a fix Spec after 0401.
- 2026-11-17: Which authorization applies to Spec 0402? — "Mesmas do BE (Recommended)": implement, commit, push, pull request, merge and release right after the Spec, up to three corrections. Makefile, CI, `go.mod` and lint only if a correction requires one.

Autonomous route: no person present during authoring; the remaining decisions were taken from the operator's brief for the TechSpec of Spec 0402 of 2026-11-17 and the measurement.

- The block's values live in one table in `internal/config`, and a repository contract test holds this repository's declaration and the configuration guide to it (ADR-0415).
- The block omits `research.tools` (ADR-0402, ADR-0415).
- A problem the inspection finds keeps today's report; the `undeclared` report applies only without one.
- `type: fix` with a `minor` Release Ceiling, because Doctor gains a warn state and the Update Check's JSON gains a state and two fields.

## Open Questions

None.
