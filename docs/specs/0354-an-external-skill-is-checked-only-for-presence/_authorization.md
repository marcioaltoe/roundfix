---
status: approved
granted: 2026-11-05
action: check a required external skill only for presence in Doctor, skills check and the managed refresh, stop comparing external skills with the Setup Snapshot, keep the owned-skill version rule, keep this repository's own lock and snapshot rule in its contract tests, and document the rule in the Roundfix skill, the Doctor and Baseline guides, the README, the glossary, ADR-0364 and the amendment notices of ADR-0221, ADR-0224, ADR-0241 and ADR-0264
consuming: 0354-an-external-skill-is-checked-only-for-presence
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/references/setup.md
  - docs/references/coverage-record.json
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

# Approved authority for Spec 0354

On 2026-11-05 at about 08:20 BRT the maintainer answered the operator's cycle
AZ batch 2 questions with AskUserQuestion, after Fluxus reported that
`roundfix doctor` and `roundfix baseline update` held its external skills,
installed from `marcioaltoe/skills`, to the Setup Snapshot
(`references/2026-11-05-only-roundfix-owned-skills-are-checked-by-digest-or-version.md`).
The answers this Spec rests on are, verbatim:

- Question "Pedido do Fluxus (Spec 0354): as skills do próprio Roundfix devem
  ter o mesmo digest do marcioaltoe/skills, ou versão igual ou maior. Contra
  qual cópia o doctor confere?" → "Cópia do binário, sem rede (Recommended)"
  — option text: "O catálogo marcioaltoe/skills copia as skills do Roundfix
  do main todo dia, então a cópia dentro do binário é a mesma ou mais nova. O
  doctor aceita digest igual ou versão igual ou maior que a do binário, sem
  acessar a rede."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, com
  merge pela fila e conserto quando a entrega parar. Makefile, CI, go.mod e
  lint liberados, assim como mudanças de Baseline e skills que chegam aos
  adotantes. Release assim que cada Spec ficar pronta, com notas aos
  adotantes, aviso ao Fluxus e ao Fiscus."
- The maintainer's standing instruction of the same day: "Quero que faça os
  releases hoje assim que estiverem prontos."

"Mesmas do AZ" keeps the cycle AZ rules of 2026-11-05, and through them the
cycle AY, AX, AW, AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. This design needs none of them, so none is listed in `paths`.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-05 on the branch `docs/plan-0354`,
based on `ad77a797`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its adopted reference and
  the five ADRs.
- The six paths above are governed: the Roundfix skill's `SKILL.md` pair, its
  canonical `baseline`, `runtime` and `setup` references, and the coverage
  record `docs/references/coverage-record.json`, which task_01 re-records
  because it removes and renames top-level tests.

## What is not governed

- `skills/repository.go`, `internal/cli/doctor.go`,
  `internal/cli/required_skills.go`, `internal/cli/baseline_update.go`,
  `internal/baseline/skills_trailing.go`, and every test file tasks 01 and 02
  create, change or delete under `skills`, `internal/cli` and
  `internal/baseline`.
- The mirrors `skills/roundfix/references/baseline.md`,
  `skills/roundfix/references/runtime.md` and
  `skills/roundfix/references/setup.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/doctor.md`, `docs/user-guide/commands/baseline.md`,
  `README.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
  ADR-0221, ADR-0224, ADR-0241, ADR-0264 and ADR-0364.
- Every planning file of this Spec and its adopted reference.

## Sanctioned regeneration

The repository-owned commands write the skill mirrors and the generated
records. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

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

No mirror, version field or record is hand-edited.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`,
  `skills-lock.json`, any Vendored Skill under `.agents/skills`, any Baseline
  module, profile, template, Setup Snapshot or derived Baseline file, any
  guide under `docs/agents/`, `docs/references/skill-coverage.json`, any
  Archive Record, any skill other than `roundfix`, any Roundfix reference
  other than `baseline`, `runtime` and `setup`, or the `### QA settlement`
  section of any skill.
- No change to the owned-skill version comparison, the bundle install, the
  lock rewrite of owned entries, `roundfix baseline skills restore`,
  `roundfix baseline skills reconcile` or the deletion of Retired Skills.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix` or the real home directory;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
