---
status: approved
granted: 2026-11-14
action: reword the grouping clause of the Baseline spec-workflow module so a Spec is sized by one context and by Tasks that earn their verification cycle, with its derived files and this repository's managed refresh; carry the same criterion into the write-idea, write-prd, write-techspec and write-tasks skills, the Roundfix skill's spec reference and the spec command guide; and record the glossary change and ADR-0400's amendment notice in ADR-0208
consuming: 0391-the-grouping-bound-becomes-a-criterion
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/modules/spec-workflow.json
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

# Approved authority for Spec 0391

The Spec adopts the maintainer's observation and answers of 2026-11-14,
given through AskUserQuestion in `/roundfix --follow-up` at about
06:45–06:58 BRT (cycle BF); the maintainer, 2026-11-14. The record of them is
copied verbatim under `references/`.

- The observation, added as free text to his answer about the scope of Spec
  0392: "E também uma observação sobre a regra de 4 tasks que talvez tenha que
  virar uma nova redação: O objetivo da spec com limite de 4 tasks era porque
  tinhamos specs gigantes no inicio com varios contextos na mesma spec e
  tambem specs com tasks pequenas que consumiam runs com grandes tempos de
  verificação (rotina normal) e faziam a spec demorar muito mais do que o
  suficiente."
- The rule question: "Sobre a sua observação: hoje a regra diz "agrupe fontes
  numa Spec só enquanto ela couber em quatro Tasks de implementação mais o QA"
  (Baseline `spec-routing`, skills write-prd e write-techspec, termo Grouping
  Bound do CONTEXT.md, ADR-0208) e não diz o motivo. Você lembrou os dois
  motivos: evitar Spec gigante com vários contextos e evitar Task pequena que
  gasta um ciclo inteiro de verificação. Como a regra deve ficar?" Answer:
  "Trocar o número pelo critério" — option text: "Sai o limite fixo. A Spec
  passa a ser medida por um contexto só e por Tasks que justifiquem cada uma o
  seu ciclo de verificação; o 4 vira referência, não teto. Muda o ADR-0208 e o
  termo do glossário. Vira Spec neste ciclo."
- The authorization question: "Se a nova redação da regra das 4 Tasks virar
  Spec neste ciclo (primeira pergunta), qual autorização vale para ela? Ela
  muda um módulo do Baseline e duas skills de autoria, que os adotantes
  recebem no próximo `baseline update`." Answer: "Mesmas do BE (Recommended)"
  — option text: "Push, PR e merge; release patch depois dela; até 3
  correções; Makefile, CI, go.mod e lint só quando uma correção exigir,
  descrito no relatório."

The cycle BE terms that answer keeps are the answers of 2026-11-13, which
release to the Runs the Baseline source (`internal/baseline/**`) with every
derived file, the guides under `docs/agents/` and the owned skills, the cycle
AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are granted only for a correction that requires one;
this Spec's design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`: it corrects the wording of an existing
rule to the purpose the maintainer states and adds no command, flag or binary
behavior, so the Release Ceiling is `patch` ("release patch depois dela"), and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-14 on the branch `docs/plan-0391`,
based on `251ed81e`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: the probe file lived only in the overlay, and
  `git status --porcelain internal/speccheck` was empty afterwards.
- It covered every file the Tasks declare under `interface:` or
  `instruction:`, every file the prototype of tasks 01 and 02 changed (the
  derived Baseline files included), every planning file of this Spec, its
  adopted reference, ADR-0400 and the Source Baseline corpus of the grouping
  clause.
- The sixteen paths above are governed: the spec-workflow module, the Standard
  TypeScript profile and its formatter golden of the spec routing guide, this
  repository's spec routing guide and setup manifest (task_01), and the
  canonical `SKILL.md` of `write-idea`, `write-prd`, `write-techspec`,
  `write-tasks` and `roundfix`, the canonical Roundfix `spec` reference and the
  five `SKILL.md` mirrors under `skills/` (task_02). The probe also reported
  the `domain-modeling` skill, which task_02 only reads, and
  `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`,
  which no Task may change; neither is granted.

## What is not governed

- The Go test files tasks 01 and 02 change:
  `internal/baseline/grouped_sources_clauses_test.go` and
  `skills/acceptance_rules_test.go`.
- `skills/roundfix/references/spec.md`, `skills/testdata/owned-skill-versions.json`,
  `internal/baseline/module-versions.json` and the catalog files and plan
  goldens under `internal/baseline/testdata/`.
- `CONTEXT.md`, `docs/user-guide/commands/spec.md`,
  `docs/references/behavior-surfaces.json`, ADR-0400 and the amendment notice
  in ADR-0208.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the module version record, the derived
Baseline files, this repository's managed guide and setup manifest, the skill
mirrors, the owned-skill version record and the Behavior Surface Record. This
declaration records the regeneration that follows the approved edits and adds
no source path.

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
  - docs/agents/setup-context.json
  - docs/agents/spec-routing.md
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

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline
  module other than `spec-workflow`, any Source Baseline, Baseline template or
  Setup Snapshot, any guide under `docs/agents/` other than through the
  managed refresh, any page under `docs/user-guide/` other than
  `commands/spec.md`, `docs/references/coverage-record.json`, any Archive
  Record, `CHANGELOG.md`, any file under `internal/cli/`, `internal/spec/`,
  `internal/speccheck/` or `internal/judge/`, any skill other than
  `write-idea`, `write-prd`, `write-techspec`, `write-tasks` and `roundfix`,
  any Roundfix reference other than `spec`, or the `### QA settlement` section
  of any skill.
- No clause is added, removed or renamed, and no code, check or judge question
  counts a Spec's Tasks.
- No test, Verification command or QA row:
  - reaches the network, calls ai-memory or runs the real `gh`;
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task sends an ai-memory message or a Spec Event Announcement; the
  operator does.
