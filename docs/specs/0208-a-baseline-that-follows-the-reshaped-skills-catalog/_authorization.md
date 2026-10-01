---
status: approved
granted: 2026-09-30
action: give the Baseline setup snapshots the upstream names go, rust and typescript, retire go-cli, drop the removed review and triage skills, give the external-triage module the triage clauses, require typesafe-ai and crafting-effective-readmes in core, and bring this repository's skills to the state an adopter reaches after its next update
consuming: 0208-a-baseline-that-follows-the-reshaped-skills-catalog
paths:
  - internal/baseline/assets/setups/go-cli.json
  - internal/baseline/assets/setups/go-tui.json
  - internal/baseline/assets/setups/rust-cli.json
  - internal/baseline/assets/setups/typescript-bun.json
  - internal/baseline/assets/setups/go.json
  - internal/baseline/assets/setups/rust.json
  - internal/baseline/assets/setups/typescript.json
  - internal/baseline/assets/profiles/go-cli-tui.json
  - internal/baseline/assets/profiles/rust-cli.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/typescript.json
  - internal/baseline/assets/modules/external-triage.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/external-triage.md
  - internal/baseline/assets/source-baselines/index.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/external-triage.md
  - internal/baseline/derived_ownership_test.go
  - docs/agents/skill-dispatch.md
  - docs/agents/setup-context.json
  - skills/baseline_skill_contract_test.go
  - .agents/skills/review/SKILL.md
  - .agents/skills/typesafe-ai/LICENSE
  - .agents/skills/typesafe-ai/SKILL.md
  - .agents/skills/crafting-effective-readmes/README.md
  - .agents/skills/crafting-effective-readmes/SKILL.md
  - .agents/skills/crafting-effective-readmes/references/art-of-readme.md
  - .agents/skills/crafting-effective-readmes/references/make-a-readme.md
  - .agents/skills/crafting-effective-readmes/references/standard-readme-example-maximal.md
  - .agents/skills/crafting-effective-readmes/references/standard-readme-example-minimal.md
  - .agents/skills/crafting-effective-readmes/references/standard-readme-spec.md
  - .agents/skills/crafting-effective-readmes/section-checklist.md
  - .agents/skills/crafting-effective-readmes/style-guide.md
  - .agents/skills/crafting-effective-readmes/templates/internal.md
  - .agents/skills/crafting-effective-readmes/templates/oss.md
  - .agents/skills/crafting-effective-readmes/templates/personal.md
  - .agents/skills/crafting-effective-readmes/templates/xdg-config.md
  - .agents/skills/crafting-effective-readmes/using-references.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0208

On 2026-09-30 the maintainer expressly authorized editing the Baseline source
and its generated guides. Asked a structured question that named
`internal/baseline/assets/` and the guides under `docs/agents/`, the maintainer
answered "Autorizar os dois". The same day the maintainer said of the skills
"considere autorizado a ajustar todas as skills se necessário". On 2026-10-01
the maintainer asked for this adjustment before the next release ("esse ajuste
deve ser feito antes do próximo release"). The governed set was measured with
`GovernedPath` on `c3be3bc9`, through a `go test -overlay` probe that wrote
nothing to the repository.

## Why each governed path is unavoidable

- The seven setup snapshot paths are the four retired snapshots and the three
  that replace them. Each new one is seeded from its predecessor, then written
  by the asset sync (task_01).
- The three profiles name their setup (task_01).
  `standard-typescript-monorepo.json` also carries the golden digest pin,
  which `make baseline-digests` rewrites in every Task that edits a module.
- `core.json` drops `review` (task_01) and gains `typesafe-ai` and
  `crafting-effective-readmes` (task_03). `typescript.json` drops `triage`
  (task_01) and `crafting-effective-readmes` (task_03).
  `external-triage.json` gains its clauses (task_02).
- The two Standard TypeScript Monorepo goldens are regenerated, never
  hand-edited.
- The four Source Baseline files gain the eight clause rows in task_02. The
  corpus text, the identity, the force and the carrier of each row are written
  by hand. The regeneration fills the offsets, the digests and the identity
  record.
- `internal/baseline/derived_ownership_test.go` compares the sanctioned
  outputs with a frozen 2026-08-06 enumeration and names the setup paths that
  moved (task_01).
- `docs/agents/skill-dispatch.md` and `docs/agents/setup-context.json` are this
  repository's managed guide and Setup Manifest. The public managed refresh
  rewrites them (task_01, task_02, task_03).
- `skills/baseline_skill_contract_test.go` pins the digest of every upstream
  skill the lock declares. The constant takes the value the test reports once
  the lock set changes (task_03, task_04). No other line changes.
- The `.agents/skills/` files are the two skill trees that the public
  `baseline update` restores at `b3c45a4` (task_03), and the obsolete `review`
  tree the repository deletes after `baseline skills reconcile` (task_04). The
  restored trees are upstream content, written only by that command.

## What is not governed

Every test file except `derived_ownership_test.go`, every new test file,
`skills-lock.json`, `skills/recommended.txt`, `README.md`,
`docs/user-guide/usage.md`, the catalog snapshots, the plan goldens, and the
sanctioned parity files under `internal/baseline/testdata/` are ordinary.

## Sanctioned regeneration

The grant names the command, and the tree names its outputs.

```yaml
command: make baseline-digests
```

This repository's managed guides and Setup Manifest are rendered by
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`,
and a second refresh must report `File changes: 0`.

## Limits

- No upstream skill is edited, and `~/dev/skills` is read only.
- No Roundfix-owned skill is edited, and no `minimumVersion` in a setup
  snapshot changes.
- No `.agents/skills/` file outside the list above changes.
- No built-in profile identifier, rule identifier or existing clause changes,
  except that `rule.external-triage` moves from rule-level guidance to clauses.
- No Makefile target, lint or formatter configuration, CI workflow or
  dependency changes.
- No top-level test is removed. The one renamed test is named in task_01.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
