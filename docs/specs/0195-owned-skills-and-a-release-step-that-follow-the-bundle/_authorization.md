---
status: approved
granted: 2026-09-30
action: make the owned-skill minimum follow the bundle, record owned skill versions, list and dispatch the Roundfix skill in every setup and profile, and add the release step that checks skills and guides
consuming: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - docs/agents/specific-repository.md
  - internal/baseline/assets/setups/go-cli.json
  - internal/baseline/assets/setups/rust-cli.json
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - docs/agents/skill-dispatch.md
  - docs/agents/agent-instructions.md
  - docs/agents/setup-context.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0195

On 2026-09-30 the maintainer approved the program "Três ondas", whose block
"skills em dia" covers the owned skills and the Baseline guides. The same day
the maintainer expressly authorized editing the Baseline source and its
generated guides, answering "Autorizar os dois" to a structured question that
named `internal/baseline/assets/` and the guides under `docs/agents/`, and
said of the owned skills "considere autorizado a ajustar todas as skills se
necessário". The maintainer also decided that day that the check of skills and
Baseline happens with each release: "Essa validação e ajuste de skills e
baseline deve acontecer junto com o release programado e ao final dos próximos
release". A later decision of the same day kept in this work the Roundfix
skill's membership in the `go-cli` and `rust-cli` setups with its dispatch
trigger, and left the membership of upstream skills to another Spec. The
governed set was measured with `GovernedPath` on `9e439dbb`, through a test
overlay that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  describe Doctor's `skills:` line and the managed refresh. task_01 changes
  both behaviors, so the two sentences that describe them change, and the
  skill's version rises.
- `docs/agents/specific-repository.md` is where a Spec author reads this
  repository's hard rules. task_02 adds the rule that an owned skill's content
  changes only with its version.
- `internal/baseline/assets/setups/go-cli.json` and `rust-cli.json` gain the
  Roundfix skill entry (task_03).
- `internal/baseline/assets/modules/autonomous-work.json` requires the skill
  and carries its trigger (task_03). `internal/baseline/assets/modules/core.json`
  carries the release clause (task_04).
- The two formatter goldens and the Standard TypeScript Monorepo digest pin
  are regenerated from those sources and are never hand-edited.
- `docs/agents/skill-dispatch.md`, `docs/agents/agent-instructions.md` and
  `docs/agents/setup-context.json` are this repository's managed guides and
  Setup Manifest. The public managed refresh rewrites them.

## What is not governed

`skills/skills.go`, `skills/skills_test.go`,
`internal/cli/baseline_update.go`, `internal/baseline/assets_sync.go`,
`docs/user-guide/commands.md`, `docs/user-guide/context-driven-development.md`,
`docs/user-guide/release-runbook.md`, the record
`skills/testdata/owned-skill-versions.json`, every new test file, the catalog
snapshots `internal/baseline/testdata/catalog.diagnostics.golden.json`,
`internal/baseline/testdata/catalog.digest` and
`internal/baseline/testdata/catalog.normalized.json`, the four plan goldens
under `internal/baseline/testdata/plan-characterization/`, and the parity
files `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`
and `internal/baseline/testdata/parity-corpus/v1/manifest.json` are ordinary.
The governed files `internal/cli/cli_test.go`,
`internal/docscontract/publicdocs_test.go`,
`skills/baseline_skill_contract_test.go`,
`internal/baseline/assets/setups/typescript-bun.json` and the two other
profiles are not touched.

## Sanctioned regeneration

The grant names the command, and the tree names its outputs. No digest pin,
golden, catalog snapshot or generated guide is hand-edited.

```yaml
command: make baseline-digests
```

The Roundfix skill's mirror is regenerated with `make skills-sync`. This
repository's managed guides and Setup Manifest are rendered by the public
Baseline update,
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
A second refresh must report `File changes: 0`.

## Limits

- No clause changes its identifier or enforcement level, and no clause is
  added or removed. One sentence is appended to one clause.
- No `minimumVersion` in a setup snapshot changes, and no entry other than
  the Roundfix skill is added to or removed from a setup.
- No owned skill other than the Roundfix skill is edited, and its
  `### QA settlement` section stays byte-identical.
- No Makefile target, lint or formatter configuration, CI workflow or
  dependency changes.
- No top-level test is renamed or removed, and no exported function signature
  changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
