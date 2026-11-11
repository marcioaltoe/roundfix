---
status: approved
granted: 2026-11-11
action: record the release of the Roundfix that last applied a Baseline change as the Setup Manifest's roundfixVersion, refuse a newer recorded release in configuration loading, roundfix init --scope project and the six Baseline commands, keep Doctor, the Update Check, the Run readers and roundfix stop working without Project Config, and document it in four command guides, the Roundfix skill's baseline and setup references, the glossary and ADR-0392
consuming: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/setup.md
  - internal/cli/baseline_human_test.go
  - internal/cli/baseline_plan_test.go
  - skills/roundfix/SKILL.md
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

# Approved authority for Spec 0382

On 2026-11-11 conexus reported, in ai-memory message
`01a1c84e-8dfb-7483-966c-061a3c293853`, that a developer ran Roundfix 0.96.0
against a repository configured for 0.137.7
(`references/2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md`).
The maintainer then answered the operator's cycle BC questions with
AskUserQuestion (via `/roundfix --follow-up`, 2026-11-11 ~14:45–14:50 BRT).
The answers this Spec rests on are, verbatim:

- Question "Relato do Conexus 01a1c84e (status disse 'nenhuma run ativa' com
  Run rodando; comando parado esperando autorização sem retorno ao agente;
  tela preta; Diego usava binário 0.96 em repo configurado para 0.137.7).
  Nenhum foi reproduzido. Qual destino?" → "Backlog + Spec de versão
  (Recommended)" — option text: "Spec F agora: binário mais velho que o
  Baseline do repositório recusa com mensagem clara. Os três sintomas viram
  Backlog Entry esperando passos de reprodução de você ou do Diego."
- Question "Spec F (binário mais velho que o Baseline do repositório recusa
  com mensagem clara e nomeia 'roundfix upgrade'): autorização e posição?" →
  "Mesmas do BB, após E (Recommended)" — option text: "Push, merge, release,
  teto minor; entregue depois de D e E, fechando o ciclo."

"Mesmas do BB" keeps the cycle BB rules of 2026-11-10 ("Mesmas do BA, teto
minor (Recommended)" and "Mesmas do BB, teto minor (Recommended)"), and
through them the cycle BA, AZ, AY, AX, AW, AV, AU, AT and AS answers, which
release to the Runs the owned skills and changes of skills that reach the
adopters, the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths
declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar
(módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
go.mod, lint, que pedem autorização específica).") and the 2026-09-30 grant of
every skill ("considere autorizado a ajustar todas as skills se necessário").

This Spec needs no Makefile, `.github/workflows`, `go.mod`, `go.sum`, lint or
formatter configuration change, and no Baseline module, template or profile
change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: feat`, so the Release Ceiling is `minor` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it after the merge, delivered after Specs 0380
and 0381.

The governed set was measured on 2026-11-11 on the branch `docs/plan-0382`,
based on `013af3b4`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed only this Spec's
  planning files, its adopted reference and ADR-0392 afterwards, and the probe
  file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every file the scratch prototype of tasks 01 and 02 changed
  or created, every planning file of this Spec, its adopted reference,
  ADR-0392 and `docs/agents/setup-context.json`.
- The six paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `baseline` and `setup` references, which task_03 changes, and
  the two governed CLI test files `internal/cli/baseline_human_test.go` and
  `internal/cli/baseline_plan_test.go`, whose three parity tests task_02
  gives the writer version (TechSpec → Declared breaks). The
  probe also named `.agents/skills/domain-modeling/SKILL.md`, which task_03
  only reads, and `docs/agents/setup-context.json`, which no Task changes;
  neither is granted.

## What is not governed

- Every other Go source and test file tasks 01 and 02 create or change under
  `internal/config/`, `internal/baseline/` and `internal/cli/`.
- The mirrors under `skills/roundfix/references/` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/**`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`
  and ADR-0392.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated
records. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
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
  module, template, profile, Setup Snapshot or plan golden, any guide under
  `docs/agents/`, `docs/agents/setup-context.json`,
  `docs/references/coverage-record.json`, any Archive Record, any skill other
  than `roundfix`, any Roundfix reference other than `baseline` and `setup`,
  or the `### QA settlement` section of any skill.
- No command, Task or QA row edits `roundfixVersion` in this repository's Setup
  Manifest or runs a Baseline command against this checkout.
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
