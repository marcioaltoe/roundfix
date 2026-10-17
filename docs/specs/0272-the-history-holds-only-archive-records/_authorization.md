---
status: approved
granted: 2026-10-16
action: make a retirement delete instead of reduce, give every citation of a removed file a Git Locator in the Archive Command, the Legacy Link Rebase, the History Sanitize Command and the history section of roundfix baseline update, hold the History Root to Archive Records in the Spec Consistency Check once the Baseline carries the rule, state that rule once in the context-workflow and secondbrain modules with recorded versions and regenerated derived files, describe it in the user guide, the Roundfix and archive-spec skills and the glossary, and remove this repository's retired history through its own Managed Refresh
consuming: 0272-the-history-holds-only-archive-records
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/roundfix/references/spec.md
  - docs/agents/docs-layout.md
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
  - docs/references/coverage-record.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/docscontract/publicdocs_test.go
  - internal/docscontract/testdata/corpus-golden.json
  - internal/spec/archive.go
  - internal/spec/archive_layout_characterization_test.go
  - internal/speccheck/backlog.go
  - internal/speccheck/backlog_test.go
  - internal/speccheck/coherence.go
  - internal/speccheck/constraints.go
  - skills/archive-spec/SKILL.md
  - skills/baseline_skill_contract_test.go
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0272

On 2026-10-16 the maintainer approved cycle AH through AskUserQuestion: scope
"A: history só com records, B: ergonomia da fila, C: guardas de autoria, D:
remetente do F6", order "B → A → C → D (Recomendado)", authorizations
"Governed Paths declarados" (the option text: "\"Concedo\" como em AE/AF/AG
para os Governed Paths que cada Spec declarar: módulos da Baseline, registros
de versão, make baseline-digests, skills e testes governados. Makefile,
workflows de CI, go.mod e configuração de lint ficam fora."), "Merge e
releases", "Merge da PR #545", "Avisos por push"; failures "Regras de AE–AG
(Recomendado)". The 2026-09-30 skills grant ("considere autorizado a ajustar
todas as skills se necessário") also stands. This Spec is cycle AH item A,
second in the queue, released as v0.84.0.

No live provider call is authorized by this record. Tests, Verification and QA
use temporary repositories, temporary homes and this repository's own tree. No
test, Verification command or QA row opens a network connection, reads the
real `~/.roundfix`, or reads or writes another repository.

The governed set was measured on the authoring branch `docs/plan-0272` with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository (`git status
--porcelain` was identical before and after). It ran against all 138 files the
Tasks declare and the planning files. The 24 paths above are the governed
files the Tasks change; `.agents/skills/domain-modeling/SKILL.md` and
`docs/agents/specific-repository.md` are governed too, but the Tasks only read
them as instructions and change neither.

## Why each governed path is unavoidable

- `internal/spec/archive.go` holds the Archive Command, which task_01 makes
  rewrite the outcome and each promoted copy with Git Locators.
- `internal/speccheck/constraints.go` and `internal/speccheck/coherence.go`
  wire `SC-HISTORY-RETAINED` beside the repository-level detectors and
  register its stage (task_03). `internal/speccheck/backlog.go` and
  `internal/speccheck/backlog_test.go` carry the Fix text that now says to
  delete a terminal Backlog Entry.
- `internal/docscontract/testdata/corpus-golden.json` and
  `internal/spec/archive_layout_characterization_test.go` pin the active
  corpus code set, which the new code joins at 0 (task_03); the
  characterization test also stats this repository's Findings family, which
  the cleanup removes (task_04).
- `.agents/skills/roundfix/references/archive.md`, `baseline.md`,
  `deliver.md`, `setup.md` and `spec.md` cover the Behavior Surfaces whose
  sources the Tasks change (`command archive`, `command history sanitize`,
  `command baseline update`, `command deliver *`, `command upgrade`,
  `command spec check`). `.agents/skills/archive-spec/SKILL.md` describes the
  Archive Relink and promotion. `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md` and `skills/archive-spec/SKILL.md` carry the
  skills' version fields, raised with their content.
- `internal/baseline/assets/modules/context-workflow.json` and
  `internal/baseline/assets/modules/secondbrain.json` carry the retirement
  clause and the conditional lesson capture (task_04).
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, the
  two formatter goldens, `docs/agents/docs-layout.md`,
  `docs/agents/secondbrain.md` and `docs/agents/setup-context.json` are
  derived from those modules and change only through `make baseline-digests`
  and the Managed Refresh. The set comes from the module edits of commits
  `e4dd6a68` and `1db414e2`; task_04 measures what the commands rewrite and
  changes nothing else.

## What is not governed

The Go sources and tests under `internal/spec` (task_04's
`internal/spec/history_entries.go` and its tests among them), `internal/cli`,
`internal/speccheck` and `internal/baseline` other than those listed above;
`internal/docscontract/corpus_test.go`; `internal/baseline/module-versions.json`
and the catalog and plan characterization goldens under
`internal/baseline/testdata/`; the mirrors `skills/roundfix/references/*.md`;
`skills/testdata/owned-skill-versions.json`;
`docs/references/behavior-surfaces.json`; `CONTEXT.md`; the user guide pages;
`docs/references/**` other than the coverage record; every file under `docs/history/`; and the ADR notes
written while planning.

## Sanctioned regeneration

The repository-owned commands resolve the Baseline's derived files, the skill
mirrors and the generated records. This declaration records the regeneration
that follows the approved edits and adds no source path.

```yaml
command: make baseline-digests
```

The Module Version Record step writes each changed module's version line and
the record; suiteguard accepts those writes only for this declared command:

```yaml
command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
outputs:
  - internal/baseline/module-versions.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/secondbrain.json
```

```yaml
command: make skills-sync
```

```yaml
command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
outputs:
  - skills/testdata/owned-skill-versions.json
```

```yaml
command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
outputs:
  - docs/references/behavior-surfaces.json
```

```yaml
command: go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
outputs:
  - docs/references/coverage-record.json
```

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
is the Managed Refresh: it renders this repository's guides and Setup Manifest
and applies the history section that removes the retired families, rewrites
their citations and the broken record links. A second refresh reports the
repository current. No digest pin, golden, generated guide or record is
hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml` or the CI
  workflows.
- No Baseline clause other than `clause.context.docs-one-job-per-directory`,
  `clause.context.findings-05-done`, `clause.context.findings-08-rollup`,
  `clause.context.findings-09-archive`,
  `clause.context.backlog-01-operational-contract` and
  `clause.secondbrain.capture-trigger` changes text;
  `clause.context.findings-10-live-work-health` and
  `clause.context.findings-11-rollup-closure` leave only through the
  `replaces` of `findings-09-archive`. No Source Baseline corpus, retention or
  snapshot changes.
- No change to the `### QA settlement` section of any skill, to another
  Spec's folder, or to the Archive Record schema, the Plan Digest algorithm or
  the Archive Records of delivered or cancelled Specs beyond link destinations
  to missing files.
- No test, Verification command or QA row reads or writes the real
  `~/.roundfix`, reaches the network, or reads another repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.

## Amendment of 2026-10-17

The first Run of this Spec's delivery stopped on 2026-10-17 at task_02's
settlement with `QA-AUTH-PATHS`: retiring the reduction tests changes the test
inventory in `docs/references/coverage-record.json`, which is a Governed Path,
and its regenerating command was not declared, so the suite guard also refused
the generator's write. This amendment adds that record to `paths` and its
command to the sanctioned regeneration, under the same cycle AH grant
("registros de versão" and "testes governados"); it adds no source path and
changes nothing else. It lands on the default branch before the delivery
consumes it.

## Amendment of 2026-10-17, second

The third Run of this Spec's delivery stopped on 2026-10-17 at task_04 after
two attempts. The cleanup removes `docs/history/findings` and
`docs/history/backlog`, and two Repository Contract Tests in Governed Paths
read them: `TestDocsLayoutGuideNamesEveryHistoryFamily` in
`internal/docscontract/publicdocs_test.go` requires the guide to name every
retired family, and `copyTrackedRepository` in
`skills/baseline_skill_contract_test.go` refuses a tracked file the working
tree has deleted. `make verify-contracts`, run on a copy of the attempt's tree
outside the repository, found only these two; with both changed as task_04
now requires, all 75 Repository Contract Tests and `go test ./...` passed.
This amendment adds the two files to `paths` under the same cycle AH grant
("testes governados"). The same attempt found that the history citation scan
stops on the tracked `CLAUDE.md` link; task_04 fixes it in
`internal/spec/history_entries.go`, which is not governed, in the same commit
as the cleanup (first planned as a corrective task_06, folded into task_04 on
2026-10-17 because it cannot settle before the cleanup). Nothing else
changes, and the amendment lands on the default branch before the delivery
consumes it.
