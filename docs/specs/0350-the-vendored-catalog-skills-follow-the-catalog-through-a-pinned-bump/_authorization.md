---
status: approved
granted: 2026-11-04
action: add the Makefile target skills-third-party-bump that pins this repository's vendored catalog skills, their skills-lock.json entries and the five upstream Setup Snapshots to one full commit of the marcioaltoe/skills main through the asset sync and the confirmed restore, retire skills-update, run the bump once at c8b2e85 from the operator's clone, make the setup tests read the catalog revision from the lock, add a documentation contract that holds every vendored folder to its snapshot tree, and describe the bump in CONTEXT.md and the repository rules
consuming: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
paths:
  - .agents/skills/systematic-debugging/CREATION-LOG.md
  - .agents/skills/systematic-debugging/test-academic.md
  - .agents/skills/systematic-debugging/test-pressure-1.md
  - .agents/skills/systematic-debugging/test-pressure-2.md
  - .agents/skills/systematic-debugging/test-pressure-3.md
  - Makefile
  - docs/agents/setup-context.json
  - docs/agents/specific-repository.md
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
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

# Approved authority for Spec 0350

On 2026-11-04 the maintainer answered the operator's Cycle AX questions with
AskUserQuestion, in the `/roundfix --follow-up` mode. The Cycle AX addendum
(about 07:15 BRT) answered the two reports of marcioaltoe-skills, the skills
catalog repository. The answers this Spec rests on are, verbatim:

- 0350: "Bump fixado e repetível (Recommended)" — option text: "Um alvo de
  Makefile move as 33 entradas para o commit atual da main do catálogo (ref =
  SHA, não 'main'), reinstala e atualiza o hash; mantém o lock reprodutível e
  os digests da Baseline determinísticos. Rodo o alvo agora (c8b2e85) e a cada
  ciclo."
- Placement: "No ciclo AX, mesma v0.131.0 (Recommended)" — option text:
  "Ficam depois da 0346 na fila (as três mexem em Baseline/skills); uma
  v0.131.0 no fim com as seis; mesmas regras do AX (merge pela fila,
  Makefile/CI liberados). Nota ao catálogo e aos adotantes ao final."
- Rules (all Cycle AX Specs): "Mesmas do AW (Recommended)" — option text:
  "Sem supervisão: push e PR dos planos, entrega e merge pela fila, consertar
  e seguir em falha técnica; Makefile, .github/workflows, go.mod/go.sum e lint
  liberados quando o desenho pedir; mudança de Baseline que chega aos
  adotantes liberada."
- Release (all Cycle AX Specs): "Uma v0.131.0 no fim (Recommended)" — option
  text: "Uma minor (a mudança da Baseline e o runs list fora de repo são
  feat); notas aos 8 adotantes, resposta ao Pantheon com o defeito corrigido e
  o contorno que podem largar."

Those answers keep the Cycle AW rules of 2026-11-03 ("Mesmas do AV
(Recommended)", option text: "Push e PR dos planos, entrega e merge pela
Delivery Queue, consertar e seguir em falha técnica (pular só se o conserto
contradisser o PRD). Makefile, .github/workflows, go.mod/go.sum e config de
lint/formatter liberados quando o desenho pedir."), and through them the
cycle AJ grant of 2026-10-19 "Mexer em Governed Paths declarados" (option
text: "'Concedo' para os caminhos que cada Spec declarar (módulos Baseline,
registros, skills, testes governados; fora: Makefile, CI, go.mod, lint, que
pedem autorização específica).") and the 2026-09-30 grants of Baseline
sources, guides and `.roundfixrc.yml` ("Autorizar os dois") and of every
skill ("considere autorizado a ajustar todas as skills se necessário").

The design needs the Makefile: the decision itself is a Makefile target.
It is listed in `paths`, with the Setup Manifest the Makefile edit makes
stale. `.github/workflows`, `go.mod`, `go.sum` and the lint and formatter
configuration are not needed and not listed.

This record grants `implement`, `commit`, `push`, `pull_request`, `merge` and
`release`. The Release answer asks for one minor, v0.131.0, after the Cycle
AX Specs, so the Release Ceiling is `minor` and the release may change only
the three version-only files in `release_paths`. No Task performs the release.

The governed set was measured on 2026-11-04 on the branch `docs/plan-0350`,
based on `6a0ed904`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain` was identical before
  and after, and the probe file never existed in the tree.
- It covered every file the Tasks declare under `creates:`, `interface:` or
  `deletes:`, every planning file of this Spec, its adopted reference and
  ADR-0360.
- The thirteen files in `paths` are governed. `skills-lock.json`, the derived
  files under `internal/baseline/testdata/`, every test file, `CONTEXT.md`,
  `docs/references/behavior-surfaces.json`, ADR-0360 and the planning files
  are not.

## Why each governed path is unavoidable

- `Makefile`: the maintainer decided on a Makefile target; task_01 adds
  `skills-third-party-bump` and retires `skills-update`.
- `docs/agents/setup-context.json`: the Setup Manifest records the Makefile's
  declaration digest and the catalog digest; the managed refresh rewrites
  them after task_01's Makefile edit and after task_02's bump.
- The five `internal/baseline/assets/setups/*.json`: the asset sync is their
  only writer (ADR-0191), and ADR-0224 holds this repository to the `go`
  snapshot, so the bump must move them.
- The five `.agents/skills/systematic-debugging/` files: the catalog removed
  them; the confirmed restore deletes them (ADR-0264).
- `docs/agents/specific-repository.md`: task_03 adds the repository rule that
  the bump is the only writer of the vendored skills.

## Sanctioned regeneration

The repository-owned commands write every derived file; nothing they write is
hand-edited. task_02 runs the bump, which runs the commands below in this
order, from a repository fetched read-only from the operator's clone
`/Users/marcio/.roundfix-operator/scratch/spec-0350/skills` at
`c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84`, with `GIT_ALLOW_PROTOCOL=file`:

```yaml
command: make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0350/skills SKILLS_REF=c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84
outputs:
  - internal/baseline/assets/setups/go-cli-typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/python.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
  - skills-lock.json
  - .agents/skills/systematic-debugging/CREATION-LOG.md
  - .agents/skills/systematic-debugging/test-academic.md
  - .agents/skills/systematic-debugging/test-pressure-1.md
  - .agents/skills/systematic-debugging/test-pressure-2.md
  - .agents/skills/systematic-debugging/test-pressure-3.md
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
  - .agents/skills/systematic-debugging/CREATION-LOG.md
  - .agents/skills/systematic-debugging/test-academic.md
  - .agents/skills/systematic-debugging/test-pressure-1.md
  - .agents/skills/systematic-debugging/test-pressure-2.md
  - .agents/skills/systematic-debugging/test-pressure-3.md
```

```yaml
command: make baseline-digests
```

The managed refresh renders the Setup Manifest; task_01 runs it after its
Makefile edit and the bump runs it twice, the second reporting no change:

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/setup-context.json
```

## Limits

- No change to `.github/workflows`, `go.mod`, `go.sum`, the lint, formatter
  or test-runner configuration, `.roundfixrc.yml`, any production `.go` file,
  any Baseline module, profile, template, coverage entry, Source Baseline or
  retention file, any owned skill under `.agents/skills/` or `skills/`, any
  vendored folder by hand, any rendered guide under `docs/agents/` other than
  `setup-context.json` and `specific-repository.md`, any user guide page, or
  the `### QA settlement` section of any skill.
- No write to the operator's clone
  `/Users/marcio/.roundfix-operator/scratch/spec-0350/skills`, and no read or
  write of `/Users/marcio/dev/skills`.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads the operator's clone, except task_02's bump and the QA's second
    bump row, both read-only and with `GIT_ALLOW_PROTOCOL=file`;
  - reads `~/dev/skills` or another repository;
  - records an environment value.
- No release, tag or deployment by a Task. Verification stays Daemon-owned,
  and Task status stays Daemon-written.
