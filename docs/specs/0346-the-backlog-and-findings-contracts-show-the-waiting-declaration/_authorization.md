---
status: approved
granted: 2026-11-04
action: add the two commented Waiting Declaration keys and one sentence to the Baseline's copyable Backlog and Findings Operational Contracts in the context-workflow module, record the module's next version, regenerate the derived Baseline files and this repository's managed docs layout guide and Setup Manifest, and keep the user guide's Findings template, the glossary, the Skill Coverage Map and the Behavior Surface Record in step
consuming: 0346-the-backlog-and-findings-contracts-show-the-waiting-declaration
paths:
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
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

# Approved authority for Spec 0346

On 2026-11-04 the maintainer answered the operator's cycle AX questions with
AskUserQuestion, after the follow-up inventory of that day listed the open
Backlog Entry this Spec adopts. The answers this Spec rests on are, verbatim:

- 0346: "Sim, Spec própria (Recommended)" — option text: "Spec pequena: as
  duas chaves comentadas nos dois blocos de frontmatter + uma frase; versão
  do módulo e arquivos derivados seguem as regras da Baseline; chega aos
  adotantes no baseline update."
- Rules (all four Specs of the cycle): "Mesmas do AW (Recommended)" — option
  text: "Sem supervisão: push e PR dos planos, entrega e merge pela fila,
  consertar e seguir em falha técnica; Makefile, .github/workflows,
  go.mod/go.sum e lint liberados quando o desenho pedir; mudança de Baseline
  que chega aos adotantes liberada."
- Release (all four Specs of the cycle): "Uma v0.131.0 no fim (Recommended)"
  — option text: "Uma minor (a mudança da Baseline e o runs list fora de repo
  são feat); notas aos 8 adotantes, resposta ao Pantheon com o defeito
  corrigido e o contorno que podem largar."

Those answers keep the cycle AW rules of 2026-11-03 ("Mesmas do AV
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
`release`. The Release answer asks for one minor v0.131.0 after the four
cycle AX Specs, so the Release Ceiling is `minor` and the release may change
only the three version-only files in `release_paths`: the changelog, the npm
package manifest and the application version. No Task performs the release;
the operator cuts it.

The governed set was measured on 2026-11-04 on the branch `docs/plan-0346`,
based on `6a0ed904`, with `speccheck.GovernedPath`.

- The probe was a `go test -overlay` probe in `internal/speccheck` that wrote
  nothing to the repository: `git status --porcelain internal` showed no
  change from it, and the probe file never existed in the tree.
- It covered every file task_01 declares under `creates:` or `interface:`,
  every planning file of this Spec, its adopted reference, its index and the
  three release paths.
- The five paths above are governed: this repository's managed docs layout
  guide and Setup Manifest, the Standard TypeScript docs layout golden, the
  `context-workflow` module and the Standard TypeScript profile.
  `dist/npm/roundfix/package.json` is governed too; it is a release path, not
  a Task path.

## What is not governed

- `internal/baseline/module-versions.json`,
  `internal/baseline/waiting_declaration_clause_test.go`.
- `internal/baseline/testdata/catalog.diagnostics.golden.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json` and the four plan
  characterization goldens under
  `internal/baseline/testdata/plan-characterization/`.
- `docs/user-guide/context-driven-development.md`,
  `docs/references/behavior-surfaces.json`,
  `docs/references/skill-coverage.json`, `CONTEXT.md`.
- Every planning file of this Spec, its adopted reference and its index.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files, this
repository's managed guides and the generated records. This declaration
records the regeneration that follows the approved edits and adds no source
path.

The Module Version Record step writes the module's version line and the
record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/context-workflow.json
```

```yaml
command: make baseline-digests
```

```yaml
command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
outputs:
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
```

```yaml
command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
outputs:
  - docs/references/behavior-surfaces.json
```

The Managed Refresh runs with `HOME` at a temporary directory that is removed
afterwards, and a second refresh reports no file change. No digest pin,
golden, generated guide or record is hand-edited. No top-level test is
removed or renamed, so the coverage record is not re-recorded.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, any Baseline
  module other than `context-workflow`, any clause of that module other than
  `clause.context.backlog-01-operational-contract` and
  `clause.context.findings-01-frontmatter`, the Source Baseline corpus,
  manifest, accounting or index, `internal/baseline/assets/retention/**`, the
  Setup Snapshots, any skill, any Archive Record, any ADR, or the
  `### QA settlement` section of any skill.
- No change to how `roundfix follow-up` reads or validates a Waiting
  Declaration, and no new Spec Check finding.
- No Task performs a release, creates a tag or edits `CHANGELOG.md`, a
  version file or a package manifest.
- No test, Verification command or QA row:
  - reads or writes the real `~/.roundfix`;
  - starts an Agent session or an ACP adapter;
  - reaches the network;
  - sends a request to a provider, the judge or a reviewer;
  - reads another repository;
  - records an environment value;
  - signals a process it did not start, or ends a process by name.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
