---
status: approved
granted: 2026-11-05
action: hand the terminal file to the Run Browser, the attach cockpit and the owning cockpit and make a terminal program refuse any other writer, print one closing line when the Run Browser ends without a selection and name the repository in it, give runs list a table at an interactive terminal with --format table|plain and --wide while keeping its stable line for every other reader, let attach and runs show take the short Run id, and document the behavior in the Roundfix skill, the runs and attach guides, the help text and a new ADR
consuming: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
paths:
  - .agents/skills/roundfix/SKILL.md
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

# Approved authority for Spec 0353

On 2026-11-05 at about 06:50 BRT the maintainer showed the operator his
MacBook terminal: `roundfix runs list` printed two lines of about three
hundred characters each, and bare `roundfix attach` returned to the prompt
twice with nothing; he wrote "A lista de runs está confusa, talves com
informação em excesso e o roundfix attach não mostra nada." At about 07:10
BRT he answered the operator's cycle AZ question with AskUserQuestion. The
answers this Spec rests on are, verbatim:

- Listing: "Tabela para pessoas (Recommended)" — option text: "No terminal
  interativo, runs list mostra cabeçalho, colunas alinhadas, id curto, nome
  do repositório, Spec, hora em BRT e duração. Branch e worktree só com
  --wide. Em pipe ou com --format plain, continua a linha estável de hoje
  para os agentes."
- Rules: "Mesmas do AX" — "Sem supervisão; merge pela fila; corrigir e
  continuar nas paradas; Makefile, CI, go.mod e lint liberados; mudanças de
  Baseline e skills que chegam aos adotantes liberadas."
- Release, later the same day: "Quero que faça os releases hoje assim que
  estiverem prontos. E se existe no baseline regra para um release por dia
  por algum motivo deve ser removida." No such rule exists in the repository
  or the Baseline; it was the operator's own habit, so this Spec changes
  nothing for it.

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
them. This design needs none of them (`golang.org/x/sys` is already a direct
requirement), so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The release is one minor, v0.133.0, cut by the operator right
after this Spec merges, on 2026-11-05, so the Release Ceiling is `minor` and
the release may change only the three version-only files in
`release_paths`: the changelog, the npm package manifest and the application
version. No Task performs the release.

The governed set was measured on 2026-11-05 on the branch `docs/plan-0353`,
based on `38cc5854`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  ADR-0363.
- The three paths above are governed: the Roundfix skill's `SKILL.md` pair and
  its canonical `runs` reference. `internal/cli/cli_test.go` is governed too;
  no Task declares or changes it.

## What is not governed

- Every Go file and test the Tasks declare under `internal/cli`,
  `internal/tui` and `internal/store`.
- The mirror `skills/roundfix/references/runs.md` and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/runs.md`, `docs/user-guide/commands/attach.md`,
  `docs/references/behavior-surfaces.json` and ADR-0363.
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
  `docs/agents/`, `docs/references/skill-coverage.json`, `CONTEXT.md`, any
  Archive Record, any skill other than `roundfix`, any Roundfix reference
  other than `runs`, `internal/cli/cli_test.go`, `internal/config`, or the
  `### QA settlement` section of any skill.
- No change to the Run Database schema, `newRunID`, or the stable line of
  `runs list`.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - opens an interactive terminal program;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
