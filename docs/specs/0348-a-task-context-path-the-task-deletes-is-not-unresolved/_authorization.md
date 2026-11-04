---
status: approved
granted: 2026-11-04
action: make the Spec Consistency Check resolve a Task Context path that the checked tree no longer holds when HEAD or the Delivery Base tracks it, make the SC-REF-UNRESOLVED Task Context fix suggest removing the entry or declaring the deletion, and teach the write-tasks skill and Task template, the Roundfix spec reference, the spec command guide and the glossary that a removed file is listed only under deletes
consuming: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/spec.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
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

# Approved authority for Spec 0348

On 2026-11-04 the maintainer answered the operator's cycle AX questions with
AskUserQuestion, after Fiscus reported that a Task deleting a file its own
Task Context listed failed its settlement on `SC-REF-UNRESOLVED`. The
answers this Spec rests on are, verbatim:

- 0348: "Sim, as três sugestões (Recommended)" — option text: "Spec 0348:
  SC-REF-UNRESOLVED aceita caminho de Context que a própria Task apaga no
  diff dela ou que existe no tree base; o Fix: sugere remover a entrada em
  vez de criar o arquivo; write-tasks orienta não listar em Context um
  arquivo que a Task apaga."
- Rules (all four): "Mesmas do AW (Recommended)" — option text: "Sem
  supervisão: push e PR dos planos, entrega e merge pela fila, consertar e
  seguir em falha técnica; Makefile, .github/workflows, go.mod/go.sum e lint
  liberados quando o desenho pedir; mudança de Baseline que chega aos
  adotantes liberada."
- Release (all four): "Uma v0.131.0 no fim (Recommended)" — option text:
  "Uma minor (a mudança da Baseline e o runs list fora de repo são feat);
  notas aos 8 adotantes, resposta ao Pantheon com o defeito corrigido e o
  contorno que podem largar."

The Rules answer keeps the cycle AW rules of 2026-11-03 ("Mesmas do AV
(Recommended)", option text: "Push e PR dos planos, entrega e merge pela
Delivery Queue, consertar e seguir em falha técnica (pular só se o conserto
contradisser o PRD). Makefile, .github/workflows, go.mod/go.sum e config de
lint/formatter liberados quando o desenho pedir."), and through them the
cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor v0.131.0 after the four
cycle AX Specs, so the Release Ceiling is `minor` and the release may change
only the three version-only files in `release_paths`: the changelog, the npm
package manifest and the application version. No Task performs the release.

The governed set was measured on 2026-11-04 on the branch `docs/plan-0348`,
based on `6a0ed904`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  the three release paths.
- The six paths above are governed: the canonical `write-tasks` skill and
  Task template, the canonical Roundfix `spec` reference, and the `SKILL.md`
  pairs of the two skills whose version fields rise. The canonical
  `domain-modeling` skill is governed too, but task_02 only reads it as an
  instruction, so it is not listed. `dist/npm/roundfix/package.json` is
  governed too; it is a release path, not a Task path.

## What is not governed

- `internal/speccheck/citations.go`, `internal/speccheck/task_context_base.go`,
  `internal/speccheck/task_context_base_test.go`,
  `internal/speccheck/archive_record_test.go`,
  `internal/daemon/task_deletes_settlement_test.go`,
  `skills/task_context_deletion_test.go`.
- The mirrors `skills/roundfix/references/spec.md` and
  `skills/write-tasks/references/task-template.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/spec.md`, `docs/references/behavior-surfaces.json`,
  `CONTEXT.md`.
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
  - .agents/skills/write-tasks/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-tasks/SKILL.md
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
  the `domain-modeling` skill, or the `### QA settlement` section of any skill.
- No change to `SC-TOOLING-UNDECLARED`, the mechanical authorization audit,
  the scope audit, the `creates:` and `deletes:` labels or the `deletes:`
  settlement refusal, or to the other `SC-REF-UNRESOLVED` findings.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - reaches the network or fetches from a remote;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
