---
status: approved
granted: 2026-11-04
action: make the Read-Only Run Commands (bare runs, runs list, runs show, attach and events) load without failing on a worktree.location refusal and print it as one warning, make runs list outside a Git repository list every repository, make the refusal name the detected Git root and its .git path, and state the contract in the Roundfix runs, events and review-runs references, the command and configuration guides and CONTEXT.md
consuming: 0347-read-only-run-commands-work-outside-a-repository
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/events.md
  - .agents/skills/roundfix/references/review-runs.md
  - .agents/skills/roundfix/references/runs.md
  - docs/references/coverage-record.json
  - internal/cli/cli_test.go
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

# Approved authority for Spec 0347

On 2026-11-04 the maintainer answered the operator's cycle AX questions with
AskUserQuestion, after reporting on 2026-11-03 that `roundfix runs list` and
`roundfix attach` refused outside a repository with
`worktree.location must not be inside the repository tree`. The answers this
Spec rests on are, verbatim:

- 0347: "Spec neste ciclo (Recommended)" — option text: "runs list/attach
  funcionam fora de repositório (leem o store global), e a recusa do
  worktree.location nomeia a raiz Git detectada e o caminho do .git, para o
  usuário achar o .git perdido."
- Rules (all four): "Mesmas do AW (Recommended)" — option text: "Sem
  supervisão: push e PR dos planos, entrega e merge pela fila, consertar e
  seguir em falha técnica; Makefile, .github/workflows, go.mod/go.sum e lint
  liberados quando o desenho pedir; mudança de Baseline que chega aos
  adotantes liberada."
- Release (all four): "Uma v0.131.0 no fim (Recommended)" — option text: "Uma
  minor (a mudança da Baseline e o runs list fora de repo são feat); notas aos
  8 adotantes, resposta ao Pantheon com o defeito corrigido e o contorno que
  podem largar."

The Rules answer keeps the cycle AW rules of 2026-11-03 ("Mesmas do AV
(Recommended)", option text: "Push e PR dos planos, entrega e merge pela
Delivery Queue, consertar e seguir em falha técnica (pular só se o conserto
contradisser o PRD). Makefile, .github/workflows, go.mod/go.sum e config de
lint/formatter liberados quando o desenho pedir."), and through them the cycle
AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor v0.131.0 after the four Cycle
AX Specs, so the Release Ceiling is `minor` and the release may change only
the three version-only files in `release_paths`: the changelog, the npm
package manifest and the application version. No Task performs the release.

The governed set was measured on 2026-11-04 on the branch `docs/plan-0347`,
based on `6a0ed904`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain internal` showed no
  change from it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference,
  ADR-0357 and the three release paths.
- The seven paths above are governed: the canonical Roundfix `runs`, `events`
  and `review-runs` references, the `SKILL.md` pair of the Roundfix skill whose
  version fields rise, `internal/cli/cli_test.go`, from which task_03 removes
  the test that pinned the old refusal, and the coverage record task_03
  re-records. The canonical `domain-modeling` skill is governed too, but
  task_04 only reads it as an instruction, so it is not listed.
  `dist/npm/roundfix/package.json` is governed too; it is a release path, not
  a Task path.

## What is not governed

- `internal/config/config.go`, `internal/config/worktree_location_test.go`.
- `internal/cli/cli.go`, `internal/cli/runs.go`, `internal/cli/runs_show.go`,
  `internal/cli/attach.go`, `internal/cli/events.go`,
  `internal/cli/runbrowser.go`, `internal/cli/read_only_run_commands_test.go`.
- The mirrors `skills/roundfix/references/runs.md`, `events.md` and
  `review-runs.md`, and `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/runs.md`, `attach.md`, `events.md`,
  `docs/user-guide/configuration.md`, `docs/references/behavior-surfaces.json`,
  `CONTEXT.md`.
- ADR-0357, every planning file of this Spec and its adopted reference.

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

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

No mirror, version field or record is hand-edited. task_03 removes the
top-level test `TestRunRunsListOutsideRepositoryRequiresAll`, so the coverage
record is re-recorded with its command.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module, profile, template or derived Baseline file, any guide under
  `docs/agents/`, `docs/references/skill-coverage.json`, any Archive Record,
  the `domain-modeling` skill, or the `### QA settlement` section of any skill.
- No change to how the Git root is detected (`findGitRoot`), to `runs causes`,
  or to any command other than the five Read-Only Run Commands beyond the
  longer refusal text.
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
