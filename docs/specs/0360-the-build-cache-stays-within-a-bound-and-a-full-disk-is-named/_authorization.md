---
status: approved
granted: 2026-11-06
action: hold this repository's Makefile Go build cache to a Build Cache Bound through a new cmd/gocache-bound command that removes only idle entries, align the Project Config declaration of that cache and regenerate the setup manifest, inspect the temporary directory in the Doctor disk-space check, warn in the Delivery Queue below the Free-Space Floor before an item's Implement Run and repository gate, park a full filesystem during the item binary build or the gate as disk-full with an uncounted retry, and document the rule in the deliver and doctor guides, the README, the repository rules guide, the Roundfix deliver and setup references, the glossary, ADR-0370 and the amendment notices of ADR-0327 and ADR-0334
consuming: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/setup.md
  - .roundfixrc.yml
  - Makefile
  - docs/agents/setup-context.json
  - docs/agents/specific-repository.md
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

# Approved authority for Spec 0360

On 2026-11-06 the maintainer answered the operator's cycle AZ batch 3
questions with AskUserQuestion, after oraculum reported that this
repository's Go build cache filled the Mac mini's internal disk
(`references/2026-11-05-the-roundfix-build-cache-fills-the-mac-mini-internal-disk.md`)
and the same cache stopped the Delivery Queue that evening
(`references/2026-11-06-the-repository-build-cache-stays-within-a-size-bound.md`).
The answers this Spec rests on are, verbatim:

- Question "O cache de build das Runs do Roundfix
  ($TMPDIR/roundfix-go-build-501) chegou a 26 GB e encheu o disco interno do
  mini duas vezes ontem. Às 15:58 BRT derrubou o Postgres do Oraculum, e às
  19:07 BRT parou a fila. Acabei de limpá-lo. Como evitar que volte?" →
  "Spec: limite + aviso (Recommended)" — option text: "O Makefile limpa ou
  poda o cache acima de um tamanho antes de compilar. O doctor e a fila
  avisam quando o disco do cache fica com pouco espaço, e um park por disco
  cheio diz isso em vez de 'delivery-error'."
- Rules: "Mesmas do AZ (Recommended)" — option text: "Sem supervisão, merge
  pela fila, conserto e retry quando a entrega parar. Release assim que cada
  Spec ficar pronta, com notas aos adotantes e resposta ao Fluxus e ao
  Oraculum."
- The maintainer's standing instruction of 2026-11-05: "Quero que faça os
  releases hoje assim que estiverem prontos. E se existe no baseline regra
  para um release por dia por algum motivo deve ser removida."

"Mesmas do AZ" keeps the cycle AZ rules of 2026-11-05, and through them the
cycle AY, AX, AW, AV, AU, AT and AS answers, the cycle AJ grant of 2026-10-19
"Mexer em Governed Paths declarados" (option text: "'Concedo' para os
caminhos que cada Spec declarar (módulos Baseline, registros, skills, testes
governados; fora: Makefile, CI, go.mod, lint, que pedem autorização
específica).") and the 2026-09-30 grant of every skill ("considere autorizado
a ajustar todas as skills se necessário").

The Makefile, `.github/workflows`, `go.mod`, `go.sum` and the lint and
formatter configuration are grantable to this Spec when its design needs
them. The design needs the Makefile, for the Build Cache Bound the
maintainer's answer names ("O Makefile limpa ou poda o cache acima de um
tamanho antes de compilar"), so `Makefile` is listed in `paths`, with
`docs/agents/setup-context.json`, whose recorded Makefile digest the change
makes stale. It needs no workflow, module, lint or formatter change.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`:
the changelog, the npm package manifest and the application version. No Task
performs it; the operator cuts it right after the merge.

The governed set was measured on 2026-11-06 on the branch `docs/plan-0360`,
based on `6ceb7ff5`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, every planning file of this Spec, its two adopted
  references and the three ADRs.
- The eight paths above are governed: the Makefile, `.roundfixrc.yml`, the
  setup manifest and the repository rules guide
  `docs/agents/specific-repository.md` that task_01 changes, and the Roundfix
  skill's `SKILL.md` pair and canonical `deliver` and `setup` references
  that tasks 02 and 03 change.

## What is not governed

- `cmd/gocache-bound/main.go`, `internal/sharedcache/bound.go`, and every
  Go source and test file tasks 01 to 03 create or change under
  `internal/sharedcache`, `internal/verifyselect`, `internal/cli`,
  `internal/delivery`, `internal/store` and `internal/docscontract`.
- The mirrors `skills/roundfix/references/deliver.md` and
  `skills/roundfix/references/setup.md`, and
  `skills/testdata/owned-skill-versions.json`.
- `docs/user-guide/commands/deliver.md`, `docs/user-guide/commands/doctor.md`,
  `README.md`, `CONTEXT.md`, `docs/references/behavior-surfaces.json`,
  ADR-0327, ADR-0334 and ADR-0370.
- Every planning file of this Spec and its adopted references.

## Sanctioned regeneration

The repository-owned commands write the setup manifest, the skill mirrors and
the generated records. This declaration records the regeneration that follows
the approved edits and adds no source path.

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
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

No mirror, version field, digest or record is hand-edited.

## Limits

- No change to the CI workflows, `go.mod`, `go.sum`, the lint, formatter or
  test-runner configuration, `skills-lock.json`, any Vendored Skill under
  `.agents/skills`, any Baseline module, profile, template, Setup Snapshot or
  derived Baseline file other than `docs/agents/setup-context.json`, any
  guide under `docs/agents/` other than `specific-repository.md`,
  `docs/references/skill-coverage.json`, `docs/references/coverage-record.json`,
  any Archive Record, any skill other than `roundfix`, any Roundfix
  reference other than `deliver` and `setup`, or the `### QA settlement`
  section of any skill.
- No change to the Makefile beyond the bound's variables, target and
  prerequisites; no recipe removes a cache entry younger than two hours or
  empties a cache.
- No change to the Shared Build Cache trim rules, the Run start and gate
  trims, the GC Command's report, the Linux Lane, or a `run-unresolved` park.
- No test, Verification command or QA row:
  - reads, trims or fills a real build cache, the real `~/.roundfix`, the
    real home directory or the real disk;
  - runs `go clean -cache`;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value or writes a real-looking secret literal;
  - signals a process it did not start, or ends a process by name.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
