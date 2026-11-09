---
status: approved
granted: 2026-11-09
action: make the backend Baseline module require and dispatch the seven official inngest-* skills in place of inngest, make inngest a Retired Skill, move the upstream Setup Snapshots and this repository's vendored catalog skills to catalog commit e45466d2 through the Catalog Skill Bump with the derived files, the asset-sync parity row and the Setup Manifest following through their commands, and document the change in ADR-0206, CONTEXT.md, the baseline command guide and the Roundfix skill's baseline reference
consuming: 0370-the-inngest-capability-takes-the-seven-official-inngest-skills
paths:
  - .agents/skills/context7-cli/references/setup.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/modules/backend.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
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

# Approved authority for Spec 0370

On 2026-11-09 the skills catalog reported that it replaced `inngest` with
seven official `inngest-*` skills
(`references/2026-11-09-marcioaltoe-skills-inngest-replaced-by-seven-official-skills.md`,
reported by marcioaltoe-skills). The maintainer's answers this Spec rests on
(AskUserQuestion via `/roundfix --follow-up`, 2026-11-09 ~11:20 BRT) are,
verbatim:

- Question "O catálogo marcioaltoe/skills trocou `inngest` por 7 skills
  oficiais `inngest-*`. A Baseline do Roundfix ainda instala e despacha
  `inngest` (módulo backend, setups typescript e go-cli-typescript-bun,
  fixtures). Mudança de Baseline exige `baseline update` nos adotantes. Como
  seguir?" → "Spec F agora (Recommended)" — option text: "Spec fix na
  Baseline: a capacidade opcional Inngest instala as 7 skills e despacha
  pelos gatilhos sugeridos (funções/steps, eventos, flow control, setup, CLI,
  middleware); fixtures atualizadas. Release logo após o merge; os adotantes
  pegam pelo Update Check."
- Question "Para as novas Specs (D testes de disco cheio, E e F se
  aprovadas), valem as mesmas regras do ciclo BA e o corte da meia-noite?" →
  "Mesmas do BA + corte (Recommended)" — option text: "Sem supervisão, merge
  pela fila, teto de 2 corretivas, release patch assim que cada Spec ficar
  pronta, Makefile/CI/go.mod/lint liberados, Jarvis avisado. Nenhum item
  começa depois das 21:00 BRT; a fila para antes da meia-noite e o resto
  fica para amanhã."

"Mesmas do BA" keeps the cycle BA and AZ rules, and through them the cycle
AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option text:
"'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grant of every skill
("considere autorizado a ajustar todas as skills se necessário"). The first
answer's option text grants the Baseline source, its snapshots, fixtures and
goldens. The design needs no Makefile, workflow, Go module, lint or
formatter change, although the answers would allow them.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Spec is a `type: fix`, so the Release Ceiling is `patch` and
the release may change only the three version-only files in `release_paths`.
No Task performs it; the operator cuts the patch right after the merge.

The governed set was measured on 2026-11-09 on the branch `docs/plan-0370`,
based on `cfc08683`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` showed no change from
  it, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `instruction:`, measured by running the planned commands on a scratch copy,
  every planning file of this Spec, its adopted reference and ADR-0380.
- The thirteen paths above are governed: the Baseline module, profile,
  snapshots and formatter golden task_01 changes through its commands, the
  Setup Manifest, the vendored `context7-cli` file the bump moves, and the
  canonical and mirrored Roundfix skill files task_02 changes.

## Sanctioned regeneration

The repository-owned commands write every derived file; nothing they write is
hand-edited. task_01 runs them in this order. The bump fetches commit
`e45466d2b9e15480eb286a4c19b513494a867bce` read-only from the operator's
clone `/Users/marcio/.roundfix-operator/scratch/spec-0370/skills`, with
`GIT_ALLOW_PROTOCOL=file`.

The Module Version Record step writes the changed module's version line;
suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/backend.json
```

```yaml
command: make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0370/skills SKILLS_REF=e45466d2b9e15480eb286a4c19b513494a867bce
outputs:
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
  - skills-lock.json
  - .agents/skills/context7-cli/references/setup.md
  - docs/agents/setup-context.json
```

```yaml
command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
outputs:
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
```

```yaml
command: roundfix baseline skills restore --profile go-cli-tui --skill <each catalog entry> --source-dir <fetched> --repo . --confirm-plan <digest> --format text
outputs:
  - skills-lock.json
  - .agents/skills/context7-cli/references/setup.md
```

The asset-sync parity row follows through Spec 0208's jq transform with
Spec 0317's composed-snapshot filter, quoted in `_techspec.md` → API
Contract 3, and the regeneration fills its digests and the corpus manifest:

```yaml
command: make baseline-digests
```

The managed refresh renders the Setup Manifest from a binary built after
the digests; the second run reports no change:

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
```

task_02 mirrors and records the owned skill and re-records the Behavior
Surface Record last:

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1
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

No digest pin, golden, generated guide, snapshot, vendored skill file, lock
entry, mirror or record is hand-edited; the parity row is changed only by the
quoted `jq` transform. No top-level test is removed or renamed, so the
coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module other than `backend`, any profile's source fields, template, Source
  Baseline, retention file or coverage entry, any frozen parity fixture other
  than `v1/fixtures/asset-sync.json` and `v1/manifest.json`, any rendered
  guide under `docs/agents/` other than `setup-context.json`, any owned skill
  other than `roundfix`, any vendored folder by hand, any ADR other than
  ADR-0206, or the `### QA settlement` section of any skill.
- No write to the operator's clone
  `/Users/marcio/.roundfix-operator/scratch/spec-0370/skills`, and no read or
  write of `/Users/marcio/dev/skills`.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads the operator's clone, except task_01's bump and the QA gate's
    second bump row, both read-only and with `GIT_ALLOW_PROTOCOL=file`;
  - reads `~/dev/skills` or another repository;
  - records an environment value.
- No Task releases, tags or deploys; the release is the operator's, under the
  Release Ceiling above. Verification stays Daemon-owned, and Task status
  stays Daemon-written.
