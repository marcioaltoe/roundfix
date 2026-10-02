---
task: task_02
spec: 0216-baseline-wording-left-after-the-stack-wave
status: pending
type: backend
complexity: medium
---

# Task 02: The backend and frontend guides name their declared workspace

## Overview

The two built-in TypeScript profiles declare `packages/backend` and `packages/frontend`, but the backend and frontend guides say only "the workspace". This Task binds each declared workspace to the guide that governs it and renders the bound path in the guide's scope sentence through a new `workspace.location` token. A guide with no bound workspace, which includes every repository-owned profile, renders the token empty and keeps its bytes. It answers the Backlog Entry "Baseline wording left after the stack wave" of 2026-09-30 (item: the declared workspace paths).

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST add `"guide": "guide.backend"` to the `packages/backend` workspace entry and `"guide": "guide.frontend"` to the `packages/frontend` entry of both `internal/baseline/assets/profiles/standard-typescript-monorepo.json` and `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`, editing in place, so the composed profile's compared fields stay equal.
2. MUST insert `{{workspace.location}}` after "workspace" in the scope sentence of `templates/guides/backend.md` and `templates/guides/frontend.md`, keeping every line break where it is, and MUST list the token and raise the version of both templates in `templates/index.json` as the TechSpec's "Data Models" gives.
3. MUST render the token in `internal/baseline/plan.go`: an empty default for `guide.backend` and `guide.frontend`, and, for a guide a built-in profile binds, " at " followed by the backticked path, or the sorted paths joined with ", " and a final " and ", from a `profileWorkspaceLocations` function that skips an entry without a guide or with a path that is not safe and relative.
4. MUST update `TestTheBackendAndFrontendGuidesSayWhatTheyGovern` to expect "at `packages/backend`" and "at `packages/frontend`" in the two scope sentences.
5. MUST create `internal/baseline/workspace_location_test.go` with the three tests the TechSpec's Testing Approach 2 names.
6. MUST run `make baseline-digests` twice (the second reports `"changed":false`), then the Managed Refresh twice; the second MUST report `File changes: 0`. MUST NOT hand-edit a snapshot, golden or digest.
7. MUST NOT change any clause, decision, Skill Activation, setup snapshot or skill, and MUST NOT change an exported function signature.

## Subtasks

- [ ] Bind the workspaces in both profiles.
- [ ] Add the token to the templates and the template index.
- [ ] Render the token and update and create the tests.
- [ ] Regenerate and refresh twice.

## Acceptance Criteria

- [ ] A Standard TypeScript Monorepo plan's backend and frontend guides name `packages/backend` and `packages/frontend` in their scope sentences.
- [ ] A guide with no bound workspace renders the generic sentence; two paths bound to one guide render sorted and joined; an unsafe path is never rendered.
- [ ] Every workspace of every built-in profile binds a guide its modules render.
- [ ] A second regeneration and a second Managed Refresh change nothing.

## Context

- instruction: `docs/adr/0222-a-baseline-guide-says-only-what-holds-for-the-repository-that-reads-it.md`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/backend.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/frontend.md`
- interface: `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/templates/guides/backend.md`
- interface: `internal/baseline/assets/templates/guides/frontend.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/stack_scope_and_http_default_test.go`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- creates: `internal/baseline/workspace_location_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheBackendAndFrontendGuidesNameTheirDeclaredWorkspace|TestAGuideWithoutADeclaredWorkspaceKeepsItsGenericScope|TestEveryBuiltInWorkspaceBindsAGuideItsProfileRenders|TestTheBackendAndFrontendGuidesSayWhatTheyGovern|TestTheComposedProfileTakesTheComposedSetup|TestTheComposedProfilePlanConverges|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheBackendAndFrontendGuidesNameTheirDeclaredWorkspace TestAGuideWithoutADeclaredWorkspaceKeepsItsGenericScope TestEveryBuiltInWorkspaceBindsAGuideItsProfileRenders TestTheBackendAndFrontendGuidesSayWhatTheyGovern TestTheComposedProfilePlanConverges; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the three new tests do not exist and the existing scope test expects no path, so the command fails.

## References

- `_prd.md` → Goal 1; Core Feature 2; Success Metric 1; Success Metric 4; Declared breaks
- `_techspec.md` → Exact texts; Interfaces; Data Models; Version changes; Existing tests that change; API Contract 1; API Contract 2; Testing Approach 2; Build Order 2
- ADR-0059, ADR-0061, ADR-0067, ADR-0190, ADR-0204, ADR-0222
