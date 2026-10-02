---
status: approved
granted: 2026-09-30
action: scope the backend bucket clause under a new identity that replaces the old one, render the declared workspace paths in the backend and frontend guides, and add the core clause that has the Agent ask the person to run a skill only a person can start, with their Source Baseline rows and retention disposition
consuming: 0216-baseline-wording-left-after-the-stack-wave
paths:
  - internal/baseline/assets/modules/backend.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - internal/baseline/assets/profiles/go-cli-typescript-monorepo.json
  - internal/baseline/assets/templates/index.json
  - internal/baseline/assets/templates/guides/backend.md
  - internal/baseline/assets/templates/guides/frontend.md
  - internal/baseline/assets/source-baselines/index.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/backend.md
  - internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md
  - internal/baseline/plan_test.go
  - docs/agents/skill-dispatch.md
  - docs/agents/setup-context.json
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0216

On 2026-09-30 the maintainer expressly authorized editing the Baseline source
and its generated guides, answering "Autorizar os dois" to a structured
question that named `internal/baseline/assets/` and the guides under
`docs/agents/`. On 2026-10-01 the maintainer authorized the whole next cycle,
Specs A to F, unattended ("Tudo, de A a F"); this Spec is item E. The governed
set was measured on 2026-10-02 with `GovernedPath` at `6364d3c9`, through a
test overlay that wrote nothing to the repository, against every file a
disposable-clone rehearsal of the three implementation Tasks changed. The
eighteen paths above are exactly the governed ones; the other changed paths
(new and existing tests outside `plan_test.go`, `internal/baseline/plan.go`,
catalog snapshots and plan goldens under `internal/baseline/testdata/`, and
`internal/cli/baseline_update_test.go`) are not governed.

## Why each governed path is unavoidable

- `modules/backend.json` carries the replaced bucket clause (task_01) and
  `modules/core.json` the person-only skill clause (task_03).
- The two built-in TypeScript profiles bind each declared workspace to its
  guide (task_02); the Standard TypeScript Monorepo profile's digest pin is
  rewritten by the sanctioned regeneration in every Task. The composed
  profile takes the identical workspace change so the fields its test
  compares with the Standard TypeScript Monorepo profile stay equal.
- `templates/guides/backend.md` and `frontend.md` gain the
  `workspace.location` token and `templates/index.json` declares it and raises
  both template versions (task_02).
- The Source Baseline index, manifest, identity record and two corpus files
  gain one row for each new clause of a module the Standard TypeScript
  Monorepo selects, because catalog validation refuses a required clause
  without its row and the regeneration maintains rows without creating them
  (task_01, task_03). The regeneration fills offsets, digests and the identity
  record.
- The three formatter goldens are rewritten by the sanctioned regeneration and
  never hand-edited.
- `internal/baseline/plan_test.go`: the structural-clause retention test
  expects the bucket clause `replaced` by its successor (task_01).
- `docs/agents/skill-dispatch.md` and `setup-context.json` are this
  repository's managed guide and Setup Manifest, rewritten by the public
  Managed Refresh.

## Sanctioned regeneration

```yaml
command: make baseline-digests
```

This repository's managed guides and Setup Manifest are rendered by
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`,
and a second refresh must report `File changes: 0`.

## Limits

- Exactly one clause identity is replaced,
  `clause.backend.prohibit-generic-layers`, and its replacement is declared;
  its Source Baseline row stays. No other clause changes its identifier,
  enforcement or text.
- No Skill Activation, trigger, setup snapshot or `auth.provider` behavior
  changes.
- No skill, vendored or owned, is edited.
- No adopter repository or the Secondbrain is written.
- No Makefile target, lint or formatter configuration, CI workflow or
  dependency changes.
- No top-level test is renamed or removed, and no exported function signature
  changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
