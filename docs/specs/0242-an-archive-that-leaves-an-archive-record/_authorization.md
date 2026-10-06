---
status: approved
granted: 2026-10-06
action: make an archive leave an Archive Record and remove the Spec folder, fix every reader of archived Specs to work from the record or from Git, add Jev-advised archive planning with promotion upstream, and describe the record in the skills, guides, Baseline clauses and glossary
consuming: 0242-an-archive-that-leaves-an-archive-record
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/qa-gate/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - docs/references/coverage-record.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/derived_ownership_test.go
  - internal/baseline/derived_regeneration_repocontract_test.go
  - internal/spec/archive.go
  - internal/spec/archive_layout_characterization_test.go
  - internal/spec/archive_test.go
  - internal/spec/coverage_test.go
  - internal/speccheck/backlog.go
  - internal/speccheck/governed_repocontract_test.go
  - internal/suiteguardcontract/regeneration.go
  - skills/archive-spec/SKILL.md
  - skills/owned_skill_edit_repocontract_test.go
  - skills/qa-gate/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0242

On 2026-09-30 the maintainer asked for unattended work through every release
of the program. Of the skills the maintainer said "considere autorizado a
ajustar todas as skills se necessário". Of Baseline sources and guides the
maintainer said "Autorizar os dois".

On 2026-10-06 the maintainer decided the history clean-up:

- "Não quero nada no histórico que não seja relevante para o secondbrain".
- The clean-up "deve atuar também no que já existe no diretório". Spec 0243
  owns that part.
- Scope "Tudo, na ordem sugerida", of which this Spec is fourth, in v0.53.0.
- Existing history is migrated in batches after a `history-full` tag
  ("Aplicar em lotes via PR"). That is Spec 0243's work, not this one's.
- For the Governed Paths this Spec declares, the maintainer said "Concedo".

On 2026-10-06 the operator also decided, for this Spec only: "the briefing's
rule against editing the `### QA settlement` section is lifted for this Spec
only, because the maintainer's history decision ("Não quero nada no
histórico que não seja relevante para o secondbrain") and grant ("Concedo")
make the current text false". The three tables are rewritten identically in
the canonical archive-spec, qa-gate and Roundfix `SKILL.md` files and their
mirrors, which this record already bounds.

The first of the maintainer's decisions is the explicit approval ADR-0215 requires before an
archive may remove a Spec folder. It applies to the archives this Spec's
runtime performs from now on. It does not authorize deleting, moving or
rewriting anything already under `docs/history`.

The governed set was measured with `GovernedPath` on the authoring branch at
`40a7893d`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare. 24 of them are governed.

## Why each governed path is unavoidable

- `internal/spec/archive.go` — the cut replaces the move inside `Archive`.
- `internal/spec/archive_test.go` and
  `internal/spec/archive_layout_characterization_test.go` — the archive's
  own tests assert the record. Their corpus tests read the pinned commit
  instead of the folders.
- `internal/speccheck/backlog.go` — the promoted Backlog Entry check names
  archived Specs by folder only.
- `internal/suiteguardcontract/regeneration.go` — every operative sanctioned
  regeneration lives in archived authorization records today. The record
  must carry them, or the suite guard loses its authority when folders go.
- `internal/spec/coverage_test.go` and `docs/references/coverage-record.json`
  — the record lists three packages under `docs/history`, and it fails when
  their folders go.
- `internal/speccheck/governed_repocontract_test.go`,
  `internal/baseline/derived_regeneration_repocontract_test.go`,
  `internal/baseline/derived_ownership_test.go` and
  `skills/owned_skill_edit_repocontract_test.go` — the repository-contract
  tests that read the archived folders, measured by the ablation.
- The canonical archive-spec, qa-gate and Roundfix skill files and their
  `SKILL.md` mirrors. The skill-sync rule requires a Pull Request that
  changes CLI behavior to ship the skill update, and an owned skill's
  content changes only with its version.
- `internal/baseline/assets/modules/spec-workflow.json` and
  `context-workflow.json` — the two clauses that describe the archive.
  `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, the
  formatter golden `docs-layout.md` and
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json` are
  their sanctioned derived files.

## What is not governed

These are ordinary:

- the other Go sources and tests under `internal/spec`, `internal/cli`,
  `internal/worktree`, `internal/runcause`, `internal/speccheck`,
  `internal/specaudit`, `internal/suiteguardcontract`, `internal/gittest`,
  `internal/judge`, `internal/authorization` and `internal/baseline`;
- `internal/baseline/testdata/**`;
- the user guides, `CONTEXT.md` and
  `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`;
- the skill reference mirror `skills/roundfix/references/archive.md` and
  `skills/testdata/owned-skill-versions.json`.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors and the Baseline's
derived files. This declaration records the regeneration that follows the
approved edits and adds no source path.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

`make baseline-digests` regenerates the catalog snapshots, the plan goldens,
the formatter goldens and the profile digest.
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
renders this repository's guides and Setup Manifest, and a second refresh
reports no file change. No digest pin, golden or generated guide is
hand-edited.

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No deletion, move or rewrite of any file already under `docs/history`.
  Spec 0243 owns the migration and the `history-full` tag.
- No test, Verification command or QA row reaches a provider, reads a
  credential, or reads or writes the real `~/.roundfix`. Archive Advice is
  exercised only against a fake transport.
- No change to the `### QA settlement` section beyond the identical
  rewrite task_04 declares, and no change to `CHANGELOG.md`.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
