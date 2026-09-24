---
task: task_04
spec: 0164-knowledge-lifecycle-and-capture
status: pending
type: docs
complexity: medium
---

# Task 04: Every retired family has one documented home

## Overview

The layout guide names the backlog, findings and handoff history directories but not the ADR, review and Spec ones that `spec.ArchiveDir` resolves, and the resolver has no handoff family, so guide and resolver can disagree unnoticed.

This is an authorized tooling Task. It may change only `internal/spec/archive.go`, `internal/spec/archive_test.go`, `internal/baseline/assets/modules/context-workflow.json`, `internal/docscontract/publicdocs_test.go`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, the derived pins the sanctioned regeneration rewrites, and this Task file. Stop before any other governed mutation. The bounded set comes from [_authorization.md](_authorization.md).

## Requirements

1. MUST add `ArchiveKindHandoff` (`docs/history/handoffs`) and an exported `ArchiveKinds()` listing every kind to `internal/spec/archive.go`, without changing any existing `ArchiveDir` answer.
2. MUST extend `clause.context.docs-one-job-per-directory` in `internal/baseline/assets/modules/context-workflow.json`, keeping its identity and bumping versions by the catalog convention, so it names every directory `ArchiveKinds()` resolves.
3. MUST state in that clause that `rejected`, `deprecated` and `superseded` ADRs retire to `docs/history/adr/` while a `proposed` one stays, and that a finished orphan Review Artifact retires to `docs/history/reviews/` under ADR-0163 while a new Round never writes into history.
4. MUST state how each family records its disposition, with the literal phrase `dated disposition addendum` for Findings and Backlog Entries, lifecycle front matter for ADRs, the archive stamp for Specs, and a byte-identical move for Review Artifacts and handoffs.
5. MUST render `docs/agents/docs-layout.md` and `docs/agents/setup-context.json` through `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`, then run `make baseline-digests`, and MUST NOT hand-edit a pin; a second refresh changes nothing.
6. MUST add `TestDocsLayoutGuideNamesEveryHistoryFamily` to `internal/docscontract/publicdocs_test.go`, reading the repository's rendered guide and failing for any `ArchiveDir` of `ArchiveKinds()` it does not name.

## Subtasks

- [ ] Add the handoff family and the list of kinds.
- [ ] Extend the clause and regenerate the guide, manifest and pins.
- [ ] Add the guide-against-resolver test.

## Acceptance Criteria

- [ ] `ArchiveKinds()` includes the handoff family and every earlier kind resolves as before.
- [ ] The rendered guide names every history directory the resolver returns.
- [ ] The clause is sourced from the module and a second managed refresh is a no-op.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/archive.go`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `docs/agents/docs-layout.md`
- interface: `internal/docscontract/publicdocs_test.go`

## Verification

- `grep -q "dated disposition addendum" internal/baseline/assets/modules/context-workflow.json || exit 1; out="$(go test -count=1 -tags docscontract -v -run "^(TestArchiveKindsNameEveryRetiredFamily|TestDocsLayoutGuideNamesEveryHistoryFamily)$" ./internal/spec ./internal/docscontract 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestArchiveKindsNameEveryRetiredFamily TestDocsLayoutGuideNamesEveryHistoryFamily; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task the module lacks the phrase and neither named case exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — History families
