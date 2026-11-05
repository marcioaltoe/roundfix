---
status: approved
granted: 2026-11-05
action: run the pre-PR review of Delivery Queue items beside the Queue Owner's other work up to a User Config Review Capacity, keep a queued item from starting while an earlier item is in review, record each item's Queue Turn in a Turn Record beside the owner's console log, print the turns and the review concurrency in deliver status, and document the behavior in the Roundfix skill, the deliver and configuration guides, CONTEXT.md and a new ADR
consuming: 0355-a-queue-item-shows-its-turn-and-the-pre-pr-review-gets-its-own-capacity
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
release_ceiling: minor
release_paths:
  - CHANGELOG.md
  - dist/npm/roundfix/package.json
  - internal/app/version.go
---

# Approved authority for Spec 0355

On 2026-11-05 at 08:07 BRT Fiscus reported, at the maintainer's request, a
Delivery Queue of about twenty-one Specs that had run for two days, with five
items in `reviewing` and no review running while one Implement Run held the
queue owner. At about 08:20 BRT the maintainer answered the operator's cycle
AZ batch 2 questions with AskUserQuestion. The answers this Spec rests on
are, verbatim:

- Fiscus points: question "Relato do Fiscus (fila de 21 Specs há 2 dias).
  Quais pontos viram Spec agora? Os que ficarem de fora viram Backlog Entry."
  — all four options selected; the one this Spec answers: "1. Vez na fila
  (Recommended)" — "O deliver status mostra 'aguardando capacidade' e a
  posição. O review pré-PR, que só lê, ganha uma concorrência própria,
  separada da do implement."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, com
  merge pela fila e conserto quando a entrega parar. Makefile, CI, go.mod e
  lint liberados, assim como mudanças de Baseline e skills que chegam aos
  adotantes. Release assim que cada Spec ficar pronta, com notas aos
  adotantes, aviso ao Fluxus e ao Fiscus."
- Release, the same day: "Quero que faça os releases hoje assim que
  estiverem prontos. E se existe no baseline regra para um release por dia
  por algum motivo deve ser removida." No such rule exists in the repository
  or the Baseline, so this Spec changes nothing for it.

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
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The release is one minor, cut by the operator right after this
Spec merges, so the Release Ceiling is `minor` and the release may change
only the three version-only files in `release_paths`: the changelog, the npm
package manifest and the application version. No Task performs the release.

The governed set was measured on 2026-11-05 on the branch `docs/plan-0355`,
based on `ad77a797`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0365.
- The three paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver` reference. `.agents/skills/domain-modeling/SKILL.md`
  is governed too; task_01 declares it as an `instruction:` it reads and
  changes nothing in it, so it is not listed.

## What is not governed

- Every Go file and test the Tasks declare under `internal/config`,
  `internal/delivery` and `internal/cli`.
- The mirror `skills/roundfix/references/deliver.md` and
  `skills/testdata/owned-skill-versions.json`.
- `CONTEXT.md`, `docs/user-guide/commands/deliver.md`,
  `docs/user-guide/configuration.md`, `docs/references/behavior-surfaces.json`
  and ADR-0365.
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

No mirror, version field or record is hand-edited. No top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or derived Baseline file, any guide under
  `docs/agents/`, `docs/references/skill-coverage.json`, any Archive Record,
  any skill other than `roundfix`, any Roundfix reference other than
  `deliver`, `internal/cli/cli_test.go`, or the `### QA settlement` section of
  any skill.
- No change to the Run Database schema, a Delivery Stage, the item line of
  `deliver status`, a blocker text, a Park Class or the exit codes of any
  `deliver` subcommand.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session, an ACP adapter or a real Queue Owner process;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
