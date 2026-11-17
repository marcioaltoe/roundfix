---
status: approved
granted: 2026-11-17
action: let each Authoring Phase declare a fallback model and effort, add the commands that decide and record an Authoring Fallback and print the model a phase starts on, tell the four authoring skills and the Roundfix skill to switch and tell, and describe all of it in the guides, the glossary and the amended ADR-0401 and ADR-0407
consuming: 0405-authoring-model-fallback
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/phase.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
  - release
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0405

The maintainer's answers are recorded whole in [`references/2026-11-17-maintainer-answers.md`](references/2026-11-17-maintainer-answers.md). The answer this record rests on is, verbatim (AskUserQuestion, 2026-11-17 about 17:50 BRT): authorization for Spec 0405 → "Mesmas do BE (Recommended)".

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it:

- the grant of 2026-09-30 for every skill: "considere autorizado a ajustar todas as skills se necessário";
- the cycle AJ grant of 2026-10-19, "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

Under the same rules, the Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one, with up to five corrections. This Spec's design needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`. The Spec is a `type: feat`. Its Release Ceiling is `minor` because the design adds a Project Config key, `fallback` under an entry of `skills.authoring.claude`, and two command paths, `roundfix phase model` and `roundfix phase fallback`. The release may change only the three version-only files in `release_paths`. No Task performs it; the operator cuts it after the merge, corrections included.

## How the governed set was measured

The author measured the governed set on 2026-11-17 in a clone on the branch `docs/plan-0405`, based on `a8236a87`, with `speccheck.GovernedPath`. [`references/2026-11-17-harness-and-test-measurements.md`](references/2026-11-17-harness-and-test-measurements.md) records the measurement.

- The measurement used a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository.
- It covered every source, test, guide, skill, record and configuration file the TechSpec names or expects the Tasks to create, ADR-0401, ADR-0407 and the new ADRs, and four controls.
- The twelve paths above are the governed ones this Spec's design changes.
- These governed paths are not granted, and no Task may change them:
  - `internal/cli/cli_test.go`; new tests go in new files;
  - `internal/docscontract/publicdocs_test.go`;
  - `.roundfixrc.yml`;
  - the controls: the `Makefile`, `.github/workflows/ci-verify.yml` and `go.mod`.

## What is not governed

- Every Go source and test file under `internal/config`, `internal/phaserecord` and `internal/cli` that the Tasks declare, `internal/cli/cli.go` included.
- The mirrors under `skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json` and the tests under `skills/`.
- `docs/user-guide/commands/phase.md`, `docs/user-guide/configuration.md`, `docs/user-guide/usage.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, `CONTEXT.md`, ADR-0401, ADR-0407, ADR-0427, ADR-0428 and ADR-0429.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the skill mirrors and the derived records. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/write-techspec/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - skills/write-tasks/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest or record is hand-edited.

## Operator step after the release

This repository's `.roundfixrc.yml` gains its fallbacks only after a released Roundfix that knows the key is installed on every machine that runs it. Until then, the binary of the Delivery Queue reads Project Config from the default branch and refuses an unknown key (ADR-0408). The operator then adds, under `skills.authoring.claude`:

- `fallback: {model: opus, reasoning_effort: xhigh}` to `idea`, `prd` and `techspec`;
- `fallback: {model: sonnet, reasoning_effort: high}` to `tasks`, which stays declared `opus`/`high`.

That change is the operator's, under the maintainer's answers of 2026-11-17, and no Task makes it.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`, except for a correction that requires one;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any file under `.claude/agents/`;
  - any Baseline module, template, profile, Source Baseline or Setup Snapshot, and any file under `docs/agents/`;
  - any other skill, Vendored Skill or Archive Record;
  - the `### QA settlement` section of any skill.
- No change to any file of another repository, and no read of one.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`, `codex` or `claude`;
  - reads or writes the real `~/.roundfix`, `~/.codex`, `~/.claude` or the real home directory;
  - starts an Agent session, an Author Subagent or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - routes an OpenAI or Anthropic model through OpenRouter;
  - signals a process it did not start, or ends a process by name;
  - reads the clock to decide a value, or waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
