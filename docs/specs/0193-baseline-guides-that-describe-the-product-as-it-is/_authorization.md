---
status: approved
granted: 2026-09-30
action: make the Baseline clauses state what the shipped product does, remove the duplicated backend clause, keep repository-only citations out of shipped guidance, and correct the repository guide's archive wording
consuming: 0193-baseline-guides-that-describe-the-product-as-it-is
paths:
  - internal/baseline/assets/modules/backend.json
  - internal/baseline/assets/modules/autonomous-work.json
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
  - internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md
  - internal/baseline/assets/profiles/standard-typescript-monorepo.json
  - docs/agents/autonomous-work.md
  - docs/agents/agent-instructions.md
  - docs/agents/spec-routing.md
  - docs/agents/docs-layout.md
  - docs/agents/setup-context.json
  - docs/agents/specific-repository.md
  - internal/speccheck/backlog.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0193

On 2026-09-30 the maintainer approved the program "Três ondas", whose block
"skills em dia" includes the Baseline guides. The same day the maintainer
expressly authorized editing the Baseline source and its generated guides,
answering "Autorizar os dois" to a structured question that named
`internal/baseline/assets/` and the guides under `docs/agents/`. The maintainer
also said "considere autorizado a ajustar todas as skills se necessário"; this
Spec edits no skill and does not rely on that statement. The governed set was
measured with `GovernedPath` on `9e439dbb`, through a test overlay that wrote
nothing to the repository. One measured path lies outside the Baseline source:
`internal/speccheck/backlog.go`, the checker's Backlog reader. The program
block authorizes making the guides true, and the Backlog clause is true only
when the checker accepts the status it names.

## Why each governed path is unavoidable

- The five module files are the only source of the clauses. `backend.json`
  loses its duplicate entry (task_01). `autonomous-work.json` carries the loop
  and hook clauses (task_02). `core.json` and `spec-workflow.json` carry the
  authorization, verification and outside-evidence clauses (task_03).
  `context-workflow.json` carries the lifecycle clauses (task_04).
- The five formatter goldens and the Standard TypeScript Monorepo digest pin
  are regenerated from those modules and are never hand-edited.
- `docs/agents/autonomous-work.md`, `agent-instructions.md`, `spec-routing.md`,
  `docs-layout.md` and `setup-context.json` are this repository's managed
  guides and Setup Manifest. The public managed refresh rewrites them.
- `docs/agents/specific-repository.md` is repository-owned. task_04 corrects
  the one sentence that still names `_archived`.
- `internal/speccheck/backlog.go` is the Spec checker's Backlog reader, and it
  holds the list of terminal statuses. task_04 adds `deferred` to that list,
  so the checker accepts the status the Backlog clause now names. A clause
  that named a status the checker reports as unknown would be the drift this
  Spec removes. The edit is one case label; no detector identifier, message
  or severity changes.

## What is not governed

`internal/spec/retirement.go`, the new test
files under `internal/baseline/`, `internal/delivery/`, `internal/speccheck/`
and `internal/spec/`, the catalog snapshots
`internal/baseline/testdata/catalog.diagnostics.golden.json`,
`internal/baseline/testdata/catalog.digest` and
`internal/baseline/testdata/catalog.normalized.json`, and the four plan
goldens under `internal/baseline/testdata/plan-characterization/` are
ordinary. The governed files `internal/speccheck/backlog_test.go`,
`internal/cli/cli_test.go`, `internal/docscontract/publicdocs_test.go` and
`docs/references/coverage-record.json` are not touched.

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

- No clause changes its identifier or enforcement level, and no clause other
  than the duplicate backend entry is removed.
- No skill, no setup snapshot, no Source Baseline corpus, no parity corpus,
  no Makefile target, no CI workflow and no dependency changes.
- No top-level test is renamed or removed, and no exported function signature
  changes.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
