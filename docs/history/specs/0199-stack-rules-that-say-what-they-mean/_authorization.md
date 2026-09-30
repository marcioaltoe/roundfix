---
status: approved
granted: 2026-09-30
action: make the TypeScript, Bun, backend and frontend stack rules name what they govern, run tests through the package script, state once that a rule governs over a skill default, and state REST as the suggested HTTP mode
consuming: 0199-stack-rules-that-say-what-they-mean
paths:
  - internal/baseline/assets/modules/bun.json
  - internal/baseline/assets/modules/typescript.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/backend.json
  - internal/baseline/assets/modules/frontend.json
  - internal/baseline/assets/skill-activations.json
  - internal/baseline/assets/contract-v1.json
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/guides/typescript-bun.md
  - internal/baseline/assets/templates/guides/bun.md
  - internal/baseline/assets/templates/guides/backend.md
  - internal/baseline/assets/templates/guides/frontend.md
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md
  - docs/agents/agent-instructions.md
  - docs/agents/skill-dispatch.md
  - docs/agents/setup-context.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0199

On 2026-09-30 the maintainer approved the program "Três ondas" and, for the
stack rules, decided "Defeitos agora, regras novas depois": release v0.22.0
takes the safe corrections to existing stack rules, and new clauses wait for a
later wave. The same day the maintainer answered "REST como padrão" for the
suggested HTTP mode and "Autorizar os dois" to a structured question that named
`internal/baseline/assets/` and the guides under `docs/agents/`. The maintainer
also said "considere autorizado a ajustar todas as skills se necessário" and
authorized `.roundfixrc.yml`; this Spec edits no skill and no Project Config and
relies on neither statement.

The governed set was measured with `GovernedPath` on `9e439dbb`, through a test
overlay that wrote nothing to the repository. Twenty-one of the thirty-three changed
paths are governed. All lie in the Baseline source or among this repository's
managed guides.

## Why each governed path is unavoidable

- The five module files are the only source of the clauses and of the guide
  versions. `bun.json` and `typescript.json` carry the four reworded stack
  clauses (task_01). `core.json` carries the dependency clause and the clause
  that requires the skills (task_02). `backend.json` and `frontend.json` change
  only the versions of their module and guide, because their guides gain a
  scope sentence (task_03).
- `skill-activations.json` is the only source of the testing trigger (task_02).
- The four guide templates are where a scope sentence can stand without
  rewording a clause the Source Baseline retains by its bytes, and
  `templates/index.json` records each template's version (task_01, task_03).
- `contract-v1.json` and the profile hold two of the three stale statements of
  the suggested HTTP mode (task_03). The profile also holds the formatter
  digest pin, which every Task's regeneration rewrites.
- The five formatter goldens are regenerated from the modules and templates
  and are never hand-edited.
- `docs/agents/agent-instructions.md`, `skill-dispatch.md` and
  `setup-context.json` are this repository's managed guides and Setup
  Manifest. The public Managed Refresh rewrites them.

## What is not governed

The new test files `internal/baseline/stack_rule_wording_test.go`,
`internal/baseline/stack_core_wording_test.go`,
`internal/baseline/stack_scope_and_http_default_test.go` and
`internal/cli/baseline_http_default_statement_test.go`, the public guide
`docs/user-guide/context-driven-development.md`, the catalog snapshots
`internal/baseline/testdata/catalog.diagnostics.golden.json`,
`internal/baseline/testdata/catalog.digest` and
`internal/baseline/testdata/catalog.normalized.json`, and the four plan goldens
under `internal/baseline/testdata/plan-characterization/` are ordinary. The
governed files `internal/baseline/plan_test.go`,
`internal/cli/baseline_documentation_contract_test.go`,
`internal/cli/baseline_human_test.go` and `internal/cli/cli_test.go` are not
touched: the new tests call their helpers from new files.

## Sanctioned regeneration

The grant names the command, and the tree names its outputs. No digest pin,
golden or generated guide is hand-edited.

```yaml
command: make baseline-digests
```

This repository's managed guides and Setup Manifest are rendered by the public
Baseline update,
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
A second refresh must report `File changes: 0`.

## Limits

- No clause is added or removed, and no clause changes its identifier or its
  enforcement level.
- No backend or frontend clause changes a byte.
- No skill, no setup, no decision declaration, no Source Baseline corpus and
  no parity corpus changes.
- No Makefile target, lint, formatter or test-runner configuration, CI
  workflow or dependency changes.
- No top-level test is renamed or removed, and no exported function signature
  changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
