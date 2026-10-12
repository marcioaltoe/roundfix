---
status: approved
granted: 2026-10-11
action: record the Better Auth base path as an optional repository decision rendered by the backend guide and the Setup Manifest, let baseline update record declared decisions, report a lock-matching skill that differs from its Setup Snapshot as a difference that is kept instead of trailing, state the vendored-skill divergence order and the two missing overriding rules in the Baseline and the roundfix skill, and align the setup-context-driven skill, the user guide, the command references and the glossary
consuming: 0257-the-baseline-respects-the-repository
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/setup.md
  - .agents/skills/setup-context-driven/SKILL.md
  - docs/agents/setup-context.json
  - docs/agents/skill-dispatch.md
  - docs/references/coverage-record.json
  - internal/baseline/assets/decisions.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/rust.json
  - internal/baseline/assets/profiles/go-cli-typescript-monorepo.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/templates/guides/backend.md
  - internal/baseline/assets/templates/index.json
  - skills/roundfix/SKILL.md
  - skills/setup-context-driven/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0257

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources, guides and
`.roundfixrc.yml` the maintainer said "Autorizar os dois".

On 2026-10-11 the maintainer settled cycle AE. For this Spec, Spec D, the
Better Auth base path becomes a recorded repository decision in the Setup
Manifest ("Decisão no Setup Manifest"), with `/api/auth` as the default that
keeps current adopters unchanged; `DR-SKILL-TRAILS-SNAPSHOT` stops calling a
possibly newer skill trailing and stops suggesting its downgrade; and the
Baseline states the order for handling a diverging vendored skill. The
maintainer's answer to "Grants dos caminhos governados" was "Concedo" for the
Governed Paths each Spec declares. The Makefile and `.github/workflows` are
not covered for this Spec.

No live provider call is authorized by this record. Tests, Verification and QA
use the embedded catalog, temporary repositories, fake restore dependencies
and this repository's own files. No adopter repository is read or written by a
Task.

The governed set was measured with `GovernedPath` on the authoring branch at
`7a8b7f70`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository. It ran against every file the Tasks
declare, and 17 of them are governed. The derived files among them were
measured by an authoring prototype in a scratch clone, which applied task_01
and task_04 in order and ran the sanctioned commands after each.

## Why each governed path is unavoidable

- `internal/baseline/assets/decisions.json` gains the `auth.base_path`
  decision (task_01).
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` and
  `go-cli-typescript-monorepo.json` select it (task_01). task_04's
  `make baseline-digests` also rewrites the Standard TypeScript profile's
  golden digest.
- `internal/baseline/assets/templates/guides/backend.md` and
  `templates/index.json` carry the new token and the template version
  (task_01).
- `internal/baseline/assets/modules/core.json` and `rust.json` hold the three
  reworded clauses and their rule versions; the Module Version Record step
  writes their version lines (task_04).
- The Standard TypeScript golden `docs/agents/skill-dispatch.md` is rewritten
  by `make baseline-digests` from the core edit (task_04).
- `docs/agents/setup-context.json` is rewritten by this repository's Managed
  Refresh after task_01 and task_04, and `docs/agents/skill-dispatch.md` after
  task_04.
- `.agents/skills/setup-context-driven/SKILL.md` and its mirror describe the
  base path (task_01). `.agents/skills/roundfix/references/baseline.md`,
  `references/setup.md`, `SKILL.md` and the `SKILL.md` mirror describe the
  declared decision, the Setup Snapshot Difference and the divergence order
  (task_02 to task_04). An owned skill's content changes only with its
  version, and the record command raises both version fields.
- `docs/references/coverage-record.json` is re-recorded by its command after
  task_03 renames the tests whose meaning inverts.

## What is not governed

These paths are ordinary: `CONTEXT.md`, the two new ADRs and the addendum to
ADR-0221, the user guide pages, `docs/references/behavior-surfaces.json`,
`internal/baseline/module-versions.json`, `internal/baseline/testdata/**`,
`skills/roundfix/references/*.md`, `skills/testdata/owned-skill-versions.json`,
and every Go source and test file the Tasks declare.

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
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/rust.json
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

`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
renders this repository's guides and Setup Manifest, and a second refresh
reports no file change. No digest pin, golden, generated guide or record is
hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml` or the CI
  workflows.
- No change to the Source Baseline corpus, manifest, accounting or index,
  `internal/baseline/assets/retention/**`, the Setup Snapshots, the asset
  sync, or any module other than `core` and `rust`.
- No change to the `### QA settlement` section of any skill, to archived
  Specs, existing Archive Records or `CHANGELOG.md`.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, or reads or writes an adopter repository.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
