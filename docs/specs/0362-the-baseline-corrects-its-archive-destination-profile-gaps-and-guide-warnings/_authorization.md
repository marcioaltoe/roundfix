---
status: approved
granted: 2026-11-06
action: give the Baseline archive clause its current destination and regenerate the derived Baseline files and this repository's docs layout guide and Setup Manifest, refuse a repository Baseline Profile that lacks a decision a template it renders requires in baseline profile validate, baseline plan and baseline update, stop the nested-carrier warning for a guide the repository wrote under docs/agents, and document both rules in the baseline command guide, the Roundfix skill's baseline reference and the setup skill, as ADR-0372 decides
consuming: 0362-the-baseline-corrects-its-archive-destination-profile-gaps-and-guide-warnings
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/setup-context-driven/SKILL.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - skills/roundfix/SKILL.md
  - skills/setup-context-driven/SKILL.md
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

# Approved authority for Spec 0362

On 2026-11-06 the maintainer answered the operator's cycle AZ batch 4
questions with AskUserQuestion, after onioncry reported that the docs layout
guide still names the retired Secondbrain inbox
(`references/2026-11-06-the-archive-clause-names-the-retired-secondbrain-inbox.md`)
and sql-library reported the opaque render failure of a repository profile
and the nested-carrier warning on its own guide
(`references/2026-11-06-a-repository-profile-without-a-template-decision-fails-opaquely.md`).
The answers this Spec rests on are, verbatim:

- Question "Dois defeitos da Baseline relatados hoje: (1) OnionCry: a regra
  de arquivamento do docs-layout ainda manda copiar arquivos para o
  "Secondbrain inbox" aposentado; (2) sql-library: com profile próprio sem
  verification.incremental, o baseline update falha com um erro de render
  opaco, e o nested-carrier-conflict acusa um guia só por citar o nome do
  marcador entre crases. O que fazer?" → "Uma Spec de correção
  (Recommended)" — option text: "Uma Spec fix: a regra passa a mandar uma
  mensagem ai-memory ao default/marcioaltoe-secondbrain. O profile validate e
  o baseline update nomeiam a decisão que falta, a linha a acrescentar e o
  comando --decision. O detector deixa de casar o nome do marcador citado em
  prosa."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, merge
  pela fila, conserto e retry quando parar. Release assim que cada Spec ficar
  pronta; resposta ao OnionCry, ao sql-library e ao Jarvis."
- The maintainer's standing instruction of 2026-11-05: "Quero que faça os
  releases hoje assim que estiverem prontos. E se existe no baseline regra
  para um release por dia por algum motivo deve ser removida."

"Mesmas do AZ" keeps the cycle AZ rules of 2026-11-05, and through them the
cycle AY, AX, AW, AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário"). The Baseline source, the guides
under `docs/agents/` and the skills are granted to this Spec as its design
needs them; it needs no Makefile, workflow, Go module, lint or formatter
change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-06 on the branch `docs/plan-0362`,
based on `9a285d77`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its two adopted
  references and ADR-0372.
- The ten paths above are governed: the `spec-workflow` module, the Standard
  TypeScript profile and its docs layout golden, this repository's docs
  layout guide and Setup Manifest that task_01 changes, and the Roundfix
  skill's `SKILL.md` pair and canonical baseline reference and the setup
  skill's `SKILL.md` pair that task_04 changes.

## Sanctioned regeneration

The repository-owned commands write the module version record, the derived
Baseline files, the setup manifest, the skill mirrors and the generated
records. No mirror, version field, digest or record is hand-edited.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/spec-workflow.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/docs-layout.md
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
  - skills/roundfix/SKILL.md
  - .agents/skills/setup-context-driven/SKILL.md
  - skills/setup-context-driven/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

## What is not governed

- Every Go source and test file tasks 01 to 03 create or change under
  `internal/baseline` and `internal/cli`, the module version record, the
  three catalog files and the four plan characterization goldens.
- The mirror `skills/roundfix/references/baseline.md` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/baseline.md`,
  `docs/references/behavior-surfaces.json` and ADR-0372.
- Every planning file of this Spec and its adopted references.
