---
status: approved
granted: 2026-11-27
action: carry the maintainer's language rule into the Baseline as a new core rule with a required Reader Language decision and an optional Code Language Exception, reword the canonical-terms domain clause, add the review and QA instructions, record the decision in a new ADR, and refresh this repository's guides with the Reader Language English
consuming: 0417-the-baseline-carries-the-language-rule
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/setup-context-driven/SKILL.md
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
  - docs/references/coverage-record.json
  - internal/baseline/assets/decisions.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/profiles/go-cli-tui.json
  - internal/baseline/assets/profiles/go-cli-typescript-monorepo.json
  - internal/baseline/assets/profiles/python.json
  - internal/baseline/assets/profiles/rust-cli.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/templates/guides/agent-instructions.md
  - internal/baseline/assets/templates/index.json
  - internal/baseline/plan_test.go
  - internal/cli/baseline_human_test.go
  - internal/cli/baseline_plan_test.go
  - internal/cli/baseline_release_gate_test.go
  - internal/docscontract/publicdocs_test.go
  - skills/qa-gate/SKILL.md
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

# Authority for Spec 0417

This record is `approved`, granted on 2026-11-27 13:40 BRT by the maintainer. The record must be on the default branch before the Run starts (`docs/agents/spec-routing.md`: "A new or widened grant must land in the delivery target's ancestry before the consuming squash delivery").

## What the maintainer has decided

Message `01a21de1-c3f0-76d0-bc5e-784e4cbffed0`, delivered by the Praetor, the session that coordinates the software projects, from `default/praetor` to this project's ai-memory inbox; the maintainer answered in the Praetor's session. Verbatim:

```text
DECISÃO DE MARCIO
pedido: 01a21dd4-1880-7c40-b7fa-7fa6badaad42
pergunta: Roundfix, Spec 0417 (regra de código em inglês no Baseline): autoriza mudar os 28 arquivos protegidos do Baseline, rodar as 7 regenerações listadas e ir até onde? Nada em Makefile, CI, go.mod ou lint.
resposta: "Tudo, release até minor (Recomendado)"
alcance: Spec 0417 do Roundfix no mini, exatamente como no pedido: os 28 caminhos listados em 01a21dd4-1880; as regenerações (a) a (g) listadas ali; implement, commit, push, pull_request, merge e release com teto minor. Nada em Makefile, CI, go.mod, lint ou .roundfixrc.yml, e nada na medição da 0421.
quando: 27/11/2026 13:40 BRT

DECISÃO DE MARCIO
pedido: 01a21d98-452e (confirmação da decisão do Praetor 01a21d98-beb5)
pergunta: Confirma a exceção declarada: o repositório cujo dono decidiu manter o código existente em outro idioma (como o Quaestor legado, decisão de 24/11, 01a20f35) declara isso com a fonte, e a regra de renomear não vale ali?
resposta: "Confirmo (Recomendado)"
alcance: a Task 2 da 0417 (language.code_exception) fica no plano.
quando: 27/11/2026 13:40 BRT
```

The 28 paths of the request `01a21dd4-1880` are the `paths` above, in the same order. Outside this grant: `Makefile`, `.github/workflows/`, `go.mod`, lint configuration, `.roundfixrc.yml`, and everything that measures the Gate Time Targets of Spec 0421.

## Sanctioned regeneration

The seven commands (a) to (g) of the request, with the outputs the TechSpec's API Contract 7 declares. No derived file is edited by hand.

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/context-workflow.json
```

```yaml
command: make baseline-digests
outputs:
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/
  - internal/baseline/testdata/catalog.digest
  - internal/baseline/testdata/catalog.normalized.json
  - internal/baseline/testdata/catalog.diagnostics.golden.json
  - internal/baseline/testdata/plan-characterization/
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text --decision language.reader=English
outputs:
  - docs/agents/agent-instructions.md
  - docs/agents/domain.md
  - docs/agents/setup-context.json
```

This repository answers `language.reader=English` and declares no `language.code_exception`.

```yaml
command: make skills-sync
outputs:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/review.md
  - .agents/skills/setup-context-driven/SKILL.md
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/setup-context-driven/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/setup-context-driven/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

The coverage record may gain only the tests the Task Graph names and lose only the tests it names as renamed or merged, with their assertions.
