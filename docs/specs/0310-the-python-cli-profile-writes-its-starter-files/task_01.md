---
task: task_01
spec: 0310-the-python-cli-profile-writes-its-starter-files
status: pending
type: backend
complexity: high
---

# Task 01: The python module declares its Starter Files and the catalog validates them

## Overview

The Baseline catalog has no way to carry a repository tooling file. This Task
adds Starter File declarations, a module field with a template kind of its
own, validates them when the catalog loads, resolves a profile's Starter Files
in order, and gives the `python` module its Makefile, Project Config and CI
workflow starters. Nothing writes them yet; task_02 does (`_prd.md` → Core
Features 1-2; ADR-0323).

## Requirements

1. MUST answer the Backlog Entry Spec 0309 adopted,
   [`../0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md`](../0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md)
   of 2026-10-22, its expected item 4.
2. MUST let a module declare `starterFiles` (`internal/baseline/catalog_load.go`
   allowed fields) and validate each declaration in
   `internal/baseline/catalog_validate.go` with the codes and rules of
   `_techspec.md` → API Contract 1 and Data Models 1-2, only when the field is
   present; a module without the field and every other template kind behave as
   before.
3. MUST create `internal/baseline/starter_files.go` with `StarterFile` and
   `(*Catalog).StarterFiles` of `_techspec.md` → Interfaces: the given
   modules in order, each module's declarations in their order, the template
   bytes and their `sha256:` content identity; an unknown module is an error.
4. MUST create the three templates under
   `internal/baseline/assets/templates/starters/python/` with the content of
   `_techspec.md` → API Contract 2 (English comments, no `{{`), add their three
   `template.starter.python.*` entries (version 1, kind `starter-file`, no
   tokens) to `internal/baseline/assets/templates/index.json` right before
   `template.guide.python`, and the three declarations to
   `internal/baseline/assets/modules/python.json` after
   `repositoryExtensions`, changing nothing else in the module by hand. Spec
   0309 creates that module; this Task's Context lists it under `creates:`
   only because it is absent from the planning tree.
5. MUST then run, in this order, and hand-edit none of their outputs:
   - `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
     which writes the next version into `python.json` and records it in
     `internal/baseline/module-versions.json`;
   - `make baseline-digests`, which rewrites the derived files this Task
     declares (measured on a prototype: `catalog.digest`,
     `catalog.normalized.json` and the four plan characterization goldens);
   - `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
     which rewrites `docs/agents/setup-context.json`, then the same command
     again, which must report `File changes: 0`; both with `HOME` set to a
     temporary directory.
   A derived file the commands leave unchanged is not edited; a file they
   change that this Task does not declare is a stop-and-report condition.
6. MUST create `internal/baseline/starter_files_test.go` with
   `TestThePythonModuleDeclaresItsStarterFiles` and
   `TestAStarterFileDeclarationIsValidated` as `_techspec.md` → Testing
   Approach describes them. Each top-level test and subtest calls
   `t.Parallel()` first.
7. MUST keep every existing test of `internal/baseline` passing without
   editing it.
8. MUST NOT write any Starter File into a repository, change any other module,
   profile or setup, any skill, any user guide page, `.roundfixrc.yml` or the
   Makefile of this repository.

## Subtasks

- [ ] Add the field, its validation and the resolver.
- [ ] Add the three templates, their index entries and the module's declarations.
- [ ] Record the module version, run `make baseline-digests` and the Managed Refresh twice, and add the tests.

## Acceptance Criteria

- [ ] `python-cli` resolves the three Starter Files in order, and no other profile resolves any.
- [ ] Each invalid declaration is refused with its own code.
- [ ] A second Managed Refresh of this repository reports no file change.

## Consumers

- `StarterFile`, `(*Catalog).StarterFiles` — consumers: `./internal/baseline`

## Context

- instruction: `docs/adr/0323-a-profile-starter-file-is-written-only-where-absent-and-only-on-a-confirmed-preview.md`
- instruction: `internal/config/config.go`
- interface: `internal/baseline/catalog_load.go`
- interface: `internal/baseline/catalog_validate.go`
- creates: `internal/baseline/starter_files.go`
- creates: `internal/baseline/assets/templates/starters/python/Makefile`
- creates: `internal/baseline/assets/templates/starters/python/roundfixrc.yml`
- creates: `internal/baseline/assets/templates/starters/python/verify.yml`
- interface: `internal/baseline/assets/templates/index.json`
- creates: `internal/baseline/assets/modules/python.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/starter_files_test.go`

## Verification

- `for name in TestThePythonModuleDeclaresItsStarterFiles TestAStarterFileDeclarationIsValidated; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/baseline && GOOS=windows go vet ./internal/baseline` — expected: exit 0. Before this Task neither test exists, so the command fails at its first check. Afterwards both exist and the whole package passes, including the catalog, the module record and the derived-file characterizations, and its tests compile for Windows.

## References

- `_prd.md` → Goals; Core Features 1-2; Success Metrics 1-2; User Story 1
- `_techspec.md` → Current behavior 1-3; Interfaces; Data Models 1-2; API Contracts 1-2; Testing Approach; Build Order 1
- ADR-0323
