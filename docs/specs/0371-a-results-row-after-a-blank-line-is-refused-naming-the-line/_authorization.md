---
status: approved
granted: 2026-11-10
action: add the one reader of a QA Report's Results rows, refuse a table line no header leads with the report path, the line and the repair on every QA settlement path, in archive and in the mechanical stage, give the gate one repair turn, tell it to continue the seeded table in the Daemon prompt and the qa-gate skill, and document it in the qa-report, settle, implement and context-driven development guides, the Roundfix settle and implement references and CONTEXT.md, as ADR-0381 decides
consuming: 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
paths:
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/settle.md
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
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

# Approved authority for Spec 0371

On 2026-11-09 pantheon reported a QA Report whose `## Results` table held a
blank line between the rows the Daemon seeded and the rows its QA Agent wrote;
Roundfix read none of the rows after it and refused the report with
`rows_blocked_environment is 1; expected 0`
(`references/2026-11-10-a-qa-report-results-row-after-a-blank-line-is-read-by-no-reader.md`).
The answers this Spec rests on are, verbatim:

- The maintainer's decision for this defect, relayed by the jarvis session
  (ai-memory message `01a1c4c0`, 2026-11-10 09:18 UTC), answering the
  operator's question "1 Spec 0371 agora (recomendada) / 2 backlog no próximo
  ciclo (default sem resposta) / 3 só ajustar a skill qa-gate":

  > Resposta à sua msg 01a1c379 ("uma decisão para o Marcio de manhã, defeito do relatório de QA"). Marcio escolheu a opção 1 às 06:20 BRT de 10/11, por AskUserQuestion na sessão do Jarvis: abrir a Spec 0371. O leitor passa a ler a tabela de resultados inteira ou a recusar apontando a linha em branco exata, e a skill/prompt de QA passa a mandar continuar a tabela sem linha em branco. Autorizado escrever, entregar e soltar a v0.137.8, com as mesmas regras do ciclo BA. Mande os blocos SPEC PARA O QUADRO e o resumo final ao default/jarvis.

- Rules (AskUserQuestion via `/roundfix --follow-up`, 2026-11-10 ~15:35
  BRT), question "Com que autorizações eu toco sem supervisão as três Specs
  deste ciclo: 0372 (review do Conexus), 0371 (relatório de QA, que você já
  decidiu às 06:20 pelo Jarvis) e 0373 (regra de resposta)?", selected
  "Mesmas do BA, teto minor (Recommended)" — "Escrever, entregar pela fila,
  merge e release de cada Spec assim que pronta (patch nas duas correções; a
  regra pode sair como minor), teto de 2 Tasks corretivas, eventos ao Jarvis,
  sem notas às caixas, e Makefile/CI/go.mod/lint só onde a Spec exigir."

"Mesmas do BA" keeps the cycle BA rules, and through them the cycle AJ grant
of 2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo'
para os caminhos que cada Spec declarar (módulos Baseline, registros, skills,
testes governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere
autorizado a ajustar todas as skills se necessário"). The design needs no
Makefile, workflow, Go module, lint or formatter change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-10 on the branch `docs/plan-0371`,
based on `9fda0b65`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0381.
- The six paths above are governed: the `qa-gate` and Roundfix skills'
  `SKILL.md` pairs and the Roundfix skill's canonical implement and settle
  references, which task_03 changes.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated
records. No mirror, version field, digest or record is hand-edited.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

## What is not governed

- Every Go source and test file tasks 01 and 02 create or change under
  `internal/spec`, `internal/speccheck`, `internal/cli` and `internal/daemon`.
- The mirrors `skills/roundfix/references/settle.md` and
  `skills/roundfix/references/implement.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/qa-report.md`,
  `docs/user-guide/commands/settle.md`,
  `docs/user-guide/commands/implement.md`,
  `docs/user-guide/context-driven-development.md`, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json` and ADR-0381.
- Every planning file of this Spec and its adopted reference.
