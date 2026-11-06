---
status: approved
granted: 2026-11-06
action: make the Delivery Queue's checking stage wait for every check run, commit status and GitHub Actions workflow run of the pull request head and for each Expected Workflow of the head's workflow files, park checks-failed naming each failed check and its link, name the pending checks in the console log, the checks-timeout blocker and the Turn line of deliver status, and document the rule in the Roundfix skill, the deliver and configuration guides, CONTEXT.md, a new ADR and notices in ADR-0311 and ADR-0365
consuming: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
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

# Approved authority for Spec 0359

On 2026-11-05 Fluxus reported, in two ai-memory messages (11:05 and 12:35
BRT), that the Delivery Queue merged its pull requests #466 and #469 with
`CI | Backend` about to fail and #477 before its CI ran. On 2026-11-06 the
maintainer answered the operator's cycle AZ batch 3 questions with
AskUserQuestion. The answers this Spec rests on are, verbatim:

- Question "O Fluxus relatou que a Delivery Queue fez merge dos PRs #466 e
  #469 com o CI vermelho e do #477 antes de o CI terminar. O ruleset do main
  não exige checks, então o GitHub dava o merge como CLEAN, e a fila confia
  só nisso. O que a fila deve fazer?" → "Esperar todos e parar
  (Recommended)" — option text: "Spec agora: na etapa de checks, a fila
  espera todos os workflows que o head disparou terminarem, até o prazo de
  sempre. Se algum falhar, mesmo sem ser obrigatório, o item para com o nome
  do check e o link."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, merge
  pela fila, conserto e retry quando a entrega parar. Release assim que cada
  Spec ficar pronta, com notas aos adotantes e resposta ao Fluxus e ao
  Oraculum."
- The standing instruction of 2026-11-05: "Quero que faça os releases hoje
  assim que estiverem prontos. E se existe no baseline regra para um release
  por dia por algum motivo deve ser removida." No such rule exists in the
  repository or the Baseline, so this Spec changes nothing for it.

The Rules answer keeps the cycle AZ rules of 2026-11-05 ("Mesmas do AX":
sem supervisão; merge pela fila; corrigir e continuar nas paradas; Makefile,
CI, go.mod e lint liberados; mudanças de Baseline e skills que chegam aos
adotantes liberadas), and through them the cycle AX, AW, AV, AU, AT and AS
answers, the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths
declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar
(módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
go.mod, lint, que pedem autorização específica).") and the 2026-09-30 grant
of every skill ("considere autorizado a ajustar todas as skills se
necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them: `gopkg.in/yaml.v3` is already a direct
dependency in `go.mod`. None is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is `type: fix`, so the release is one patch, cut by the
operator right after this Spec merges; the Release Ceiling is `patch` and the
release may change only the three version-only files in `release_paths`: the
changelog, the npm package manifest and the application version. No Task
performs the release.

The governed set was measured on 2026-11-06 on the branch `docs/plan-0359`,
based on `6ceb7ff5`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its references and
  ADR-0369.
- The three paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver` reference. `.agents/skills/domain-modeling/SKILL.md`
  is governed too; task_01 declares it as an `instruction:` it reads and
  changes nothing in it, so it is not listed.

## What is not governed

- Every Go file and test the Tasks declare under `internal/delivery` and
  `internal/cli`.
- The mirror `skills/roundfix/references/deliver.md` and
  `skills/testdata/owned-skill-versions.json`.
- `CONTEXT.md`, `docs/user-guide/commands/deliver.md`,
  `docs/user-guide/configuration.md`, `docs/references/behavior-surfaces.json`,
  ADR-0369, ADR-0311 and ADR-0365.
- Every planning file of this Spec and its references.

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

No mirror, version field or record is hand-edited. No top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or derived Baseline file, any guide under
  `docs/agents/`, `docs/references/skill-coverage.json`, any Archive Record,
  any skill other than `roundfix`, any Roundfix reference other than
  `deliver`, `internal/delivery/park_class.go`, `internal/cli/cli_test.go`,
  or the `### QA settlement` section of any skill.
- No change to the Run Database schema, a Delivery Stage, a blocker name, a
  Park Class, a configuration key or the exit codes of any `deliver`
  subcommand.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session, an ACP adapter or a real Queue Owner process;
  - reaches the network or GitHub;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
