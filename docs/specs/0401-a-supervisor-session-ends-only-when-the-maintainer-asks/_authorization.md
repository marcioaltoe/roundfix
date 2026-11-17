---
status: approved
granted: 2026-11-16
action: record every Supervisor Phase without marking a boundary or printing the Boundary Notice, and state that no phase of a Spec ends the Supervisor session, which ends only when the maintainer asks or the work stops, in the autonomous-work Baseline clause and this repository's managed guide, the four authoring skills, the Roundfix skill's modes and phase reference, the glossary, ADR-0407 and the user guide
consuming: 0401-a-supervisor-session-ends-only-when-the-maintainer-asks
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/modes.md
  - .agents/skills/roundfix/references/phase.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
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
release_ceiling: patch
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0401

The Spec adopts one Backlog Entry, recorded by the operator from the maintainer's message and answer of 2026-11-16. The maintainer's answers are recorded whole in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md). The answers this record rests on are, verbatim (AskUserQuestion, 2026-11-16 and 2026-11-17):

- Mensagem do Marcio logo depois: "Esse negocio de criar sessões a cada etapa está atrasando demais."
- Fronteira de fase → "Só por Spec (Recommended)" (~21:00 BRT): sessão encerra só depois da entrega+release da Spec; fases de autoria em agentes; Spec fix ajusta skills e o aviso do roundfix phase.
- 16/11 ~21:45 BRT: Spec 0401 (fronteira só por Spec) autorização → "Mesmas do BE (Recommended)".
- 17/11 ~10:00 BRT: Marcio notou que ao fim de implementação de Spec a sessão para pedindo reinício (prejudica autonomous work) → fronteira "Só quando você pedir (Recommended)": nenhuma parada automática (nem fase nem Spec, presente ou não); handoff só ao encerrar de fato; Phase Record segue medindo; aviso do roundfix phase sai; ajustar ~/.claude/CLAUDE.md. Emendar plano 0401 antes de rodar; skills só depois da 0401.

The answer of 2026-11-17 changes the decision and not the authorization: the operator's brief for the amendment records it as unchanged, "Mesmas do BE" of 2026-11-16. Adjusting the maintainer's own instruction file outside the repository is not part of this Spec.

"Mesmas do BE" keeps the cycle BE rules of 2026-11-13 and, through them, the grants of the cycles before it. These include the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas as skills se necessário") and the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados", whose option text reads: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização específica)."

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter configuration are granted only for a correction that requires one. This Spec's design needs none of them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release`. The Spec is a `type: fix`, and its Release Ceiling is `patch`: the design adds and removes no command, flag, record key or exit code; the `boundary` key stays in the schema, always `false`. The Boundary Notice it removes was an advisory line on standard error that never changed standard output or the exit code, and `docs/agents/cli.md` makes the placement of standard error, not each diagnostic, public. The release may change only the three version-only files in `release_paths`: the changelog, the npm package manifest and the application version. No Task performs it; the operator cuts it after the merge, up to three corrections included.

## How the governed set was measured

The author first measured the governed set on 2026-11-17 on the branch `docs/plan-0401`, based on `379803c8`, and measured it again for the amended plan on a scratch worktree of `bcd716df`, with `speccheck.GovernedPath`.

- The measurement used a `go test -overlay` probe in `internal/speccheck` that wrote nothing to the repository. The probe file never existed in the tree.
- The second measurement covered 50 paths: every source, test, guide, skill, record and document file the TechSpec names or expects the Tasks to change, every derived file the regeneration commands rewrote on that worktree with the texts reworded and the writer changed, this Spec's PRD, ADR-0413, ADR-0414, ADR-0407, the revised Backlog Entry, and three controls.
- The 17 paths above are the governed ones this Spec's design changes:
  - the `autonomous-work` Baseline module, its Standard TypeScript formatter golden and the profile's golden pin;
  - this repository's managed `docs/agents/autonomous-work.md` and its Setup Manifest;
  - the Roundfix skill's `SKILL.md` pair and its canonical `modes` and `phase` references;
  - the `SKILL.md` pairs of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`.
- These governed paths are not granted, and no Task may change them: the three controls, the `Makefile`, `.github/workflows/ci.yml` and `go.mod`.

`write-tasks` measured the exact derived set on a scratch worktree before the Run. If a sanctioned regeneration rewrites a governed path outside `paths` on the tree a Task runs on, the operator amends this record before the Run, outside the change it authorizes.

## What is not governed

- `internal/phaserecord/record.go`, `internal/phaserecord/record_test.go`, `internal/phaserecord/read_test.go`, `internal/cli/phase.go`, `internal/cli/phase_test.go` and `internal/cli/follow_up_queue_test.go`.
- `internal/baseline/module-versions.json`, the catalog snapshots and plan goldens under `internal/baseline/testdata/`, and `internal/baseline/phase_boundary_clause_test.go`.
- The mirrors `skills/roundfix/references/modes.md` and `skills/roundfix/references/phase.md`, `skills/testdata/owned-skill-versions.json`, `skills/phase_boundary_test.go` and `skills/phase_reference_test.go`.
- `docs/user-guide/commands/phase.md`, `docs/user-guide/usage.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, ADR-0407, ADR-0413 and ADR-0414.
- Every planning file of this Spec and its references.

## Sanctioned regeneration

Repository-owned commands write the derived records, the skill mirrors and this repository's managed guide. This declaration records the regeneration that follows the approved edits and adds no source path.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/autonomous-work.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
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
  - `.roundfixrc.yml`, `skills-lock.json` or the generated subagents under `.claude/agents/`;
  - any Vendored Skill under `.agents/skills`;
  - any Baseline module other than `autonomous-work.json`, any template, Source Baseline or Setup Snapshot, or any profile other than `standard-typescript-monorepo.json`;
  - `AGENTS.md`, `CLAUDE.md` or `docs/agents/specific-repository.md`;
  - `docs/references/coverage-record.json`;
  - any Archive Record;
  - the `### QA settlement` section of any skill.
- No change to any file of another repository, and no read of one.
- No test or Task Verification:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - signals a process it did not start, or ends a process by name;
  - waits on elapsed time instead of the event it asserts.
- The QA gate replays the measured phase entries in a temporary Roundfix Home and a temporary Git repository it creates and removes, never in the operator's.
- No Task sends an ai-memory message, a Session Handoff or a Spec Event Announcement; the operator does.
