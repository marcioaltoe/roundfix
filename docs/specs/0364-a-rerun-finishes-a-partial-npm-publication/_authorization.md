---
status: approved
granted: 2026-11-06
action: give the release workflow's Publication Preflight a published state, add the tag-run step that compares each already published coordinate's decompressed tarball with this run's build, make the publish step skip each coordinate that compared equal, test the workflow's own step scripts with stub curl and npm, and document the recovery in the release runbook, the glossary and ADR-0146's amendment notice, as ADR-0374 decides
consuming: 0364-a-rerun-finishes-a-partial-npm-publication
paths:
  - .github/workflows/release.yml
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

# Approved authority for Spec 0364

On 2026-11-05 the operator recorded, from the v0.130.0 and v0.132.1 release
runs, that the npm publication twice stopped part way and that the rerun of
the failed job could never complete the version
(`references/2026-11-05-the-release-workflow-survives-a-transparency-log-409.md`).
The maintainer's answers this Spec rests on (AskUserQuestion, 2026-11-06
~19:48 BRT) are, verbatim:

- Question "O backlog 'release sobrevive a um 409 do transparency log' (npm
  parou no meio em v0.130.0 e v0.132.1; a correção fica em
  .github/workflows/release.yml) vira Spec agora, entregue depois da 0361?"
  → "Spec agora (Recommended)" — option text: "Fix entregue logo depois da
  0361: preflight aceita pacote já publicado com a mesma integridade e o
  rerun publica só o que falta. Libera editar release.yml nesta Spec."
- Question "Para a Spec do 409 (se sair agora), valem as mesmas regras do
  ciclo AZ: push, merge pela Delivery Queue e release patch logo após o
  merge, sem notas aos adotantes?" → "Mesmas do AZ (Recommended)" — option
  text: "Push, merge pela fila e release patch após o merge; sem notas de
  release nas caixas."

"Mesmas do AZ" keeps the cycle AZ rules of 2026-11-05, and through them the
cycle AY, AX, AW, AV, AU, AT and AS answers, the cycle AJ grant of
2026-10-19 "Mexer em Governed Paths declarados" (option text: "'Concedo'
para os caminhos que cada Spec declarar (módulos Baseline, registros,
skills, testes governados; fora: Makefile, CI, go.mod, lint, que pedem
autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário"). The first
answer's option text grants the CI file this Spec needs:
`.github/workflows/release.yml`, because the Publication Preflight and the
publish step live there. The design needs no other workflow, no Makefile,
Go module, lint or formatter change, and no separate script: the tests run
the workflow's own step scripts.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it, and no Task or QA row publishes to npm, pushes a tag or creates
a GitHub Release; the operator cuts the patch right after the merge.

The governed set was measured on 2026-11-06 on the branch `docs/plan-0364`,
based on `cb3ed035`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference,
  the removed Backlog Entry path and ADR-0374.
- The one path above is governed: `.github/workflows/release.yml`, which
  task_02 and task_03 change.

## Sanctioned regeneration

The repository-owned command writes the Behavior Surface Record. No record
is hand-edited.

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

## What is not governed

- `internal/verifyselect/release_workflow_test.go`, which task_02 creates
  and task_03 extends.
- `CONTEXT.md`, `docs/user-guide/release-runbook.md`,
  `docs/references/skill-coverage.json`,
  `docs/references/behavior-surfaces.json`, ADR-0146 and ADR-0374.
- Every planning file of this Spec and its adopted reference.
