---
status: approved
granted: 2026-11-04
action: run the Task commit's Secret Scan as the last Settlement Check of every Verification attempt so a first-attempt Secret Match returns as Verification Feedback, list each secret refusal as its own item of the Causes report, name the failed and skipped Tasks of a run-unresolved park with the secret repair as its next action, and document the behavior in the Roundfix skill, the implement, runs and deliver guides, the glossary and a new ADR
consuming: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/runs.md
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

# Approved authority for Spec 0351

On 2026-11-04 at about 22:05 BRT the maintainer answered the operator's cycle
AY questions with AskUserQuestion, after Fiscus reported (ai-memory message
`01a1a775`, 2026-11-04 13:47 BRT) that a Secret Match of a Task commit failed
its Task with no repair turn and was visible neither in `roundfix runs
causes` nor in the Delivery Queue's park. The answers this Spec rests on are,
verbatim:

- Repair: "Scanner entra nos checks (Recommended)" — option text: "O scanner
  roda junto com os settlement checks da tentativa 1. A recusa
  (arquivo:linha:padrão, sem o valor) volta como Verification Feedback no
  único reparo, e a tentativa 2 roda o scanner de novo. Mantém a regra de um
  reparo só."
- Tracked identical value: "Não aceitar (Recommended)" — option text: "O
  scanner fica como está, olhando só as linhas adicionadas. O feedback diz
  como trocar por um valor que nenhum padrão casa. Um valor já rastreado
  também pode ser um segredo vazado, então copiá-lo não deveria passar."
- Rules: "Mesmas do AX (Recommended)" — option text: "Sem supervisão; merge
  pela fila; corrigir e continuar nas paradas; Makefile, CI, go.mod e lint
  liberados; mudanças de Baseline e skills que chegam aos adotantes
  liberadas."
- Release: "v0.132.0 ao fim (Recommended)" — option text: "Release minor
  depois do merge, com a nota aos 8 adotantes e a resposta ao Fiscus. É o
  único release de hoje, dentro do teto de 1 por dia."

The Rules answer keeps the cycle AX rules of 2026-11-04 (Rules "Mesmas do AW
(Recommended)", option text: "Sem supervisão: push e PR dos planos, entrega e
merge pela fila, consertar e seguir em falha técnica; Makefile,
.github/workflows, go.mod/go.sum e lint liberados quando o desenho pedir;
mudança de Baseline que chega aos adotantes liberada."), and through them the
cycle AW, AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19 "Mexer
em Governed Paths declarados" (option text: "'Concedo' para os caminhos que
cada Spec declarar (módulos Baseline, registros, skills, testes governados;
fora: Makefile, CI, go.mod, lint, que pedem autorização específica).") and
the 2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor, v0.132.0, after this Spec
merges, so the Release Ceiling is `minor` and the release may change only the
three version-only files in `release_paths`: the changelog, the npm package
manifest and the application version. No Task performs the release; the
operator cuts it after 00:00 BRT on 2026-11-05.

The governed set was measured on 2026-11-04 on the branch `docs/plan-0351`,
based on `de579ea2`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec and its adopted reference.
- The five paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `deliver`, `implement` and `runs` references. The canonical
  `domain-modeling` skill is governed too, but task_04 only reads it as an
  instruction, so it is not listed.

## What is not governed

- Every Go file and test the Tasks declare under `internal/daemon`,
  `internal/runevent`, `internal/runcause`, `internal/delivery` and
  `internal/cli`.
- The mirrors `skills/roundfix/references/deliver.md`,
  `skills/roundfix/references/implement.md` and
  `skills/roundfix/references/runs.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/implement.md`, `docs/user-guide/commands/runs.md`,
  `docs/user-guide/commands/deliver.md`,
  `docs/references/behavior-surfaces.json`, `CONTEXT.md` and ADR-0361.
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
  `deliver`, `implement` and `runs`, or the `### QA settlement` section of
  any skill.
- No change to `internal/secretscan`, `internal/runcause/signatures.json`,
  `internal/runcause/testdata/`, `internal/store` or the QA Report scan.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
