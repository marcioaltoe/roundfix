---
status: approved
granted: 2026-11-04
action: make the bundle install of roundfix baseline update and roundfix skills install remove the ref of each skills-lock.json entry it rewrites for an installed Roundfix-owned skill and report it, make roundfix doctor warn with DR-SKILL-LOCK-STALE naming each installed owned skill whose entry still pins another copy, and teach the rule in the doctor, baseline and skills command guides, the Roundfix skill's baseline and setup references and CONTEXT.md
consuming: 0349-the-baseline-update-leaves-owned-skill-lock-entries-unpinned
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/setup.md
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

# Approved authority for Spec 0349

On 2026-11-04 the maintainer answered the operator's cycle AX questions with
AskUserQuestion, after marcioaltoe-skills reported that `roundfix baseline
update` keeps a stale `ref` in the `skills-lock.json` entries of the
Roundfix-owned skills and that `skills update` then reinstalls the older
copy. The answers this Spec rests on are, verbatim:

- 0349: "Lock coerente + doctor avisa (Recommended)" — option text: "Spec
  nova: o baseline update grava um ref coerente com o conteúdo que instala
  (ou remove o ref das entradas que o Roundfix gerencia, a Spec escolhe depois
  de medir o skills CLI) e o `roundfix doctor` avisa quando ref e
  computedHash de uma entrada não batem, nomeando o rebaixamento que o skills
  update faria."
- Placement: "No ciclo AX, mesma v0.131.0 (Recommended)" — option text:
  "Ficam depois da 0346 na fila (as três mexem em Baseline/skills); uma
  v0.131.0 no fim com as seis; mesmas regras do AX (merge pela fila,
  Makefile/CI liberados). Nota ao catálogo e aos adotantes ao final."
- Rules (all cycle AX Specs): "Mesmas do AW (Recommended)" — option text:
  "Sem supervisão: push e PR dos planos, entrega e merge pela fila, consertar
  e seguir em falha técnica; Makefile, .github/workflows, go.mod/go.sum e
  lint liberados quando o desenho pedir; mudança de Baseline que chega aos
  adotantes liberada."
- Release (all cycle AX Specs): "Uma v0.131.0 no fim (Recommended)" — option
  text: "Uma minor (a mudança da Baseline e o runs list fora de repo são
  feat); notas aos 8 adotantes, resposta ao Pantheon com o defeito corrigido e
  o contorno que podem largar."

Those answers keep the cycle AW rules of 2026-11-03 (Rules "Mesmas do AV
(Recommended)"), and through them the cycle AJ grant of 2026-10-19 "Mexer em
Governed Paths declarados" (option text: "'Concedo' para os caminhos que cada
Spec declarar (módulos Baseline, registros, skills, testes governados; fora:
Makefile, CI, go.mod, lint, que pedem autorização específica).") and the
2026-09-30 grant of every skill ("considere autorizado a ajustar todas as
skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor, v0.131.0, after the six
cycle AX Specs, so the Release Ceiling is `minor` and the release may change
only the three version-only files in `release_paths`: the changelog, the npm
package manifest and the application version. No Task performs the release.

The governed set was measured on 2026-11-04 on the branch `docs/plan-0349`,
based on `6a0ed904`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain internal` showed no
  change from it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference,
  the two ADR files and the three release paths.
- The four paths above are governed: the canonical Roundfix `baseline` and
  `setup` references and the Roundfix `SKILL.md` pair whose version fields
  rise. The canonical `domain-modeling` skill is governed too, but task_03
  only reads it as an instruction, so it is not listed.
  `dist/npm/roundfix/package.json` is governed too; it is a release path, not
  a Task path.

## What is not governed

- `internal/baseline/owned_skills_lock.go`,
  `internal/baseline/owned_skills_lock_test.go`,
  `internal/baseline/owned_skills_lock_unpin_test.go`.
- `internal/cli/baseline_update.go`, `internal/cli/cli.go`,
  `internal/cli/doctor.go`, `internal/cli/owned_skills_lock_test.go`,
  `internal/cli/baseline_update_test.go`,
  `internal/cli/doctor_trailing_skills_test.go`,
  `internal/cli/doctor_skills_lock_test.go`.
- `skills/owned_lock_guidance_test.go`, the mirrors
  `skills/roundfix/references/baseline.md` and
  `skills/roundfix/references/setup.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/baseline.md`,
  `docs/user-guide/commands/skills.md`,
  `docs/references/behavior-surfaces.json`, `CONTEXT.md`.
- The ADR files and every planning file of this Spec and its adopted
  reference.

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
  formatter or test-runner configuration, `.roundfixrc.yml`, this
  repository's root `skills-lock.json` (Spec 0350's), any Baseline module,
  profile, Setup Snapshot, template or derived Baseline file, any guide under
  `docs/agents/`, `docs/references/skill-coverage.json`, any Archive Record,
  the `domain-modeling` skill, any skill other than the Roundfix skill, or the
  `### QA settlement` section of any skill.
- No change to an external lock entry, to `source`, `sourceType` or
  `skillPath` of an owned entry, or to the readiness contract of
  `roundfix skills check`.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session;
  - runs the skills CLI (`skills`, `npx skills`, `bunx skills`);
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
