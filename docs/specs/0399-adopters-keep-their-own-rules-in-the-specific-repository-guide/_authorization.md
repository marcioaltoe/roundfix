---
status: approved
granted: 2026-11-16
action: state in the Baseline guides, the root pointer and the starter of docs/agents/specific-repository.md that a repository's own rules and variations go only in docs/agents/specific-repository.md, with an example; name each Misplaced Repository Rule and its destination in roundfix baseline update and roundfix doctor without moving it; refresh this repository's managed guides; and describe both in the Roundfix and setup-context-driven skills, the user guide and the glossary
consuming: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/setup-context-driven/SKILL.md
  - docs/agents/agent-instructions.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - docs/agents/skill-dispatch.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/repository-extension.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/templates/extensions/specific-repository.md
  - internal/baseline/assets/templates/guides/agent-instructions.md
  - internal/baseline/assets/templates/guides/skill-dispatch.md
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/root/repository-extension.md
  - skills/roundfix/SKILL.md
  - skills/setup-context-driven/SKILL.md
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

# Approved authority for Spec 0399

The Spec adopts one Backlog Entry, recorded by the operator from the maintainer's report of 2026-11-16. The maintainer's answers are recorded whole in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md). The answers this record rests on are, verbatim (AskUserQuestion, 2026-11-16 ~18:25 BRT):

- Alcance: "Texto + detecção (Recommended)" — regra e exemplo em agent-instructions, skill-dispatch e specific-repository; baseline update/doctor, ao achar edição local num arquivo do Baseline ou regra no AGENTS.md raiz, diz qual trecho mover para specific-repository.md.
- Inventário: "Ler os checkouts (Recommended)" — só leitura dos checkouts dos adotantes em ~/dev; sem commit nem mensagem a outro projeto; inventário vira referência da Spec e casos de teste.
- Autorização: "Mesmas do BE (Recommended)" — push, PR, merge, release logo após a Spec, até 3 correções; Makefile/CI/go.mod/lint só se a correção exigir.

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it. These include the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas as skills se necessário") and the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`. The Spec is a `type: fix`. Its Release Ceiling is `minor`, not `patch`, because the design adds a public output: the `misplacedRepositoryRules` field of `roundfix/baseline-update-result/v1` and the `baseline-rules` line of `roundfix doctor`. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge, up to three corrections included.

## How the governed set was measured

The author measured the governed set on 2026-11-16 on the branch `docs/plan-0399`, based on `bcf0e675`, with `speccheck.GovernedPath`.

- The measurement used a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository. `git status --porcelain` showed only this Spec's planning files, the adopted Backlog Entry's move and ADR-0411 afterwards. The probe file never existed in the tree.
- It covered 155 paths:
  - every file under the derived Baseline directories that `make baseline-digests` scans;
  - every source, test, guide, skill, record and document file the TechSpec names or expects the Tasks to create;
  - this Spec's planning files and ADR-0411;
  - three controls.
- The 22 paths above are the governed ones this Spec's design changes:
  - the three Baseline modules and the five template files;
  - the Standard TypeScript profile's golden pin and the three formatter goldens it renders from the changed guides;
  - this repository's managed `agent-instructions.md`, `docs-layout.md`, `skill-dispatch.md` and Setup Manifest;
  - the Roundfix skill's `SKILL.md` pair and its canonical `baseline` and `setup` references;
  - the `setup-context-driven` skill's `SKILL.md` pair.
- These governed paths are not granted, and no Task may change them:
  - the Source Baseline corpus and manifests under `internal/baseline/assets/source-baselines/`;
  - every Setup Snapshot under `internal/baseline/assets/setups/`;
  - the other profiles;
  - the formatter golden `AGENTS.md`, which carries no repository-extension pointer;
  - `docs/agents/specific-repository.md`;
  - `internal/cli/cli_test.go`;
  - the three controls: the `Makefile`, `.github/workflows/ci.yml` and `go.mod`.

`write-tasks` measures the exact derived set of the Baseline edit on a scratch copy before the Run. If `make baseline-digests` or the starter change rewrites a governed path outside `paths`, the operator amends this record before the Run, outside the change it authorizes.

## What is not governed

- Every Go source, test and fixture file under `internal/baseline` and `internal/cli` that the Tasks declare. This includes `internal/baseline/module-versions.json`, the catalog snapshots and plan goldens under `internal/baseline/testdata/`, `vendored_skill_clauses_test.go` and `plan_test.go`.
- This repository's `AGENTS.md`.
- The mirrors under `skills/roundfix/references/`, and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/baseline.md`, `docs/user-guide/commands/doctor.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json` and ADR-0411.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the derived records, the skill mirrors and this repository's managed guides. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/repository-extension.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - AGENTS.md
  - docs/agents/agent-instructions.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - docs/agents/skill-dispatch.md
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/setup-context-driven/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/setup-context-driven/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

No mirror, version field, digest, golden or record is hand-edited.

## Limits

- No change to any of these:
  - the Makefile, the CI workflows, `go.mod` or `go.sum`;
  - the lint, formatter or test-runner configuration;
  - `.roundfixrc.yml` or `skills-lock.json`;
  - any Vendored Skill under `.agents/skills`;
  - any Source Baseline, Setup Snapshot, or profile other than `standard-typescript-monorepo.json`;
  - `docs/agents/specific-repository.md`;
  - `docs/references/coverage-record.json`;
  - any Archive Record;
  - the `### QA settlement` section of any skill.
- No change to any file of another repository. The QA gate may read the adopted checkouts the inventory names, and may run `roundfix baseline update` there without confirmation, which writes nothing. It confirms each `git status --porcelain` unchanged afterwards. Never with `--yes`, `--confirm-plan` or `--promote`.
- No test or Task Verification:
  - reads another repository, except a temporary one the test itself created;
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
