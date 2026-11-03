---
status: approved
granted: 2026-11-03
action: make the follow-up's repository-inbox source and the delivery plan's inbox rows skip hidden files and folders, read a Waiting Declaration (waiting_on and since) on open Backlog Entries and unresolved Findings into a waiting class that roundfix follow-up holds back by default, counts on its source lines, never sends to the judge and lists with --all, report a malformed declaration on its item, and teach the Roundfix skill's follow-up mode to accept --all, write the declaration on an item it leaves waiting and name every listed Finding's and Backlog Entry's destination in its report
consuming: 0343-the-follow-up-names-what-waits-and-where-each-item-went
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/follow-up.md
  - .agents/skills/roundfix/references/modes.md
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

# Approved authority for Spec 0343

On 2026-11-03 the maintainer answered the operator's cycle AW questions with
AskUserQuestion, in the `/roundfix --follow-up` mode, after Fluxus reported
that items waiting on someone else leave the follow-up plan without a record
and that `docs/_inbox/.gitkeep` is listed as a note. The answers this Spec
rests on are, verbatim:

- 0343: the maintainer typed (Other): "Os tres e por padrão sem waiting quando
  usado só --follow-up, mas preciso de uma flag ou parametro --all (Algo
  assim) para incluir tudo, incluindo os waiting". Context: the recommended
  option was "Os três, waiting_on + since — Campos no frontmatter:
  waiting_on: <quem/o quê> e since: <AAAA-MM-DD>. O item sai como classe
  'waiting', com motivo e data. Eu respondo ao Fluxus com o nome dos campos."
- Rules (all three): "Mesmas do AV (Recommended)" — option text: "Push e PR
  dos planos, entrega e merge pela Delivery Queue, consertar e seguir em falha
  técnica (pular só se o conserto contradisser o PRD). Makefile,
  .github/workflows, go.mod/go.sum e config de lint/formatter liberados quando
  o desenho pedir."
- Release (all three): "Uma v0.130.0 no fim (Recommended)" — option text:
  "Minor único depois das três Specs, porque a 0342 cria a chave
  pre_pr_review.omit_paths e a 0343 cria o --all. Notas aos 8 adotantes,
  resposta ao Fluxus com os campos waiting_on/since, aviso ao Vortex da chave
  nova e aviso ao Jarvis de que a versão saiu."

Those answers keep the cycle AV rules of 2026-11-02 (Rules "Mesmas do AU
(Recommended)", option text: "Sem supervisão: push, PR do plano, entrega e
merge pela fila; Makefile/CI/go.mod/lint liberados se o desenho pedir; falha
técnica = consertar e seguir; findings do review corrigidos ou dispensados com
evidência; mensagem do Pantheon consumida quando a Spec estiver no main."), and
through them the cycle AU answer of 2026-11-01 "Mesmas do AT (Recommended)",
the cycle AS answers of 2026-10-31 ("Consertar e seguir (Recommended)" and the
Tooling answer "Makefile, .github/workflows, go.mod / go.sum, Config de
lint/formatter"), the cycle AJ grant of 2026-10-19 "Mexer em Governed Paths
declarados" (option text: "'Concedo' para os caminhos que cada Spec declarar
(módulos Baseline, registros, skills, testes governados; fora: Makefile, CI,
go.mod, lint, que pedem autorização específica).") and the 2026-09-30 grant of
every skill ("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor, v0.130.0, after Specs 0342,
0343 and 0344 merge, so the Release Ceiling is `minor` and the release may
change only the three version-only files in `release_paths`: the changelog,
the npm package manifest and the application version. No Task performs the
release; the operator cuts it.

The governed set was measured on 2026-11-03 on the branch `docs/plan-0343`,
based on `4e61bf6d`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec and its adopted reference.
- The five paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver`, `follow-up` and `modes` references. The canonical
  `domain-modeling` skill is governed too, but task_01 and task_02 only read
  it as an instruction, so it is not listed.

## What is not governed

- `internal/cli/deliver_plan.go`, `internal/cli/deliver_plan_test.go`,
  `internal/cli/follow_up.go`, `internal/cli/follow_up_inventory.go`,
  `internal/cli/cli.go`, `internal/cli/follow_up_test.go`,
  `internal/cli/follow_up_status_test.go`,
  `internal/cli/follow_up_waiting_test.go`.
- The mirrors `skills/roundfix/references/deliver.md`,
  `skills/roundfix/references/follow-up.md` and
  `skills/roundfix/references/modes.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/follow-up.md`, `docs/user-guide/commands/deliver.md`,
  `docs/references/behavior-surfaces.json`, `CONTEXT.md`, ADR-0353 and the
  open Backlog Entry this planning change adds.
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
  any skill other than `roundfix`, or the `### QA settlement` section of any
  skill.
- No change to `internal/speccheck`, `internal/judge`, the follow-up's
  `--status` view beyond refusing `--all`, or the `roundfix/follow-up/v1` and
  `roundfix/follow-up-messages/v1` schema names.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
