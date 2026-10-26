---
task: task_01
spec: 0310-the-python-profile-writes-its-starter-files
status: completed
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
   [`../0309-a-python-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md`](../../references/2026-10-22-a-baseline-profile-for-python-repositories.md)
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

- [ ] `python` resolves the three Starter Files in order, and no other profile resolves any.
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


## Result

Implemented the task_01 catalog slice. Modules may declare optional
`starterFiles`; loading validates the array shape, catalog-wide starter IDs
and paths, safe relative paths, and template existence and kind. The resolver
returns copies of verbatim template bytes in the supplied module order and
then declaration order, with SHA-256 content identities, and rejects unknown
modules. Python declares the Makefile, Project Config and workflow starters;
no repository Starter File writer was added.

Acceptance evidence:

- Python resolves the three paths in order and the other four profiles resolve
  none: `TestThePythonModuleDeclaresItsStarterFiles` passed. It checks IDs,
  module ownership, content identities, required template lines, absence of
  interpolation markers, returned-byte isolation, repeated module order and
  unknown-module refusal.
- Invalid declarations receive their own diagnostics:
  `TestAStarterFileDeclarationIsValidated` passed, covering unsafe, absolute,
  Windows and empty paths; duplicate paths; unknown and guide templates;
  missing and duplicate IDs; malformed arrays and entries; and compatibility
  when the optional field is absent.
- A second Managed Refresh reported `File changes: 0` and
  `Idempotence: verified`. The first refresh changed only
  `docs/agents/setup-context.json` (the catalog digest).

Focused checks and sanctioned regeneration, all exit 0:

1. `go test ./internal/baseline -run '^(TestThePythonModuleDeclaresItsStarterFiles|TestAStarterFileDeclarationIsValidated)$' -count=1`.
2. `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1` recorded Python module version 2.
3. `make baseline-digests` regenerated only `catalog.digest`,
   `catalog.normalized.json` and the four declared plan characterization
   goldens; its outputs were not hand-edited.
4. `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`, twice, each with a fresh temporary `HOME`.
   Existing dependency storage was reused via `GOPATH`, with `GOPROXY=off`
   and `GOSUMDB=off`; no dependency download was allowed. Both refreshes
   reported nested-carrier warnings for the embedded formatter and Source
   Baseline fixtures and left those carriers unchanged.
5. `go test ./internal/baseline -run '^(TestThePythonModuleDeclaresItsStarterFiles|TestAStarterFileDeclarationIsValidated|TestEmbeddedCatalog|TestCatalogCompatibility|TestCatalogDiagnosticCharacterization|TestBaselinePlanCharacterization|TestEveryBaselineModuleVersionIsRecorded)$' -count=1` passed after regeneration and refresh.
6. `git diff --check` passed.

Environment and scope: sandbox initialization rejected the symlinked writable
root even when the resolved working directory was supplied. Approved escalated
execution used the assigned checkout. The only pre-existing dirty path was
this Task file (its Daemon-owned status); that status was preserved. All
implementation and generated paths remain within this Task's declarations.

Not run: the declared Verification commands, full package suite, selected
repository Verification and incremental Verification remain Daemon-owned.
No existing test, other Task, Task Graph, skill, other module, profile, setup,
repository Makefile or repository Project Config was edited. No commit, push
or pull request was made.

## Carry-forward provenance

- Source Run: `run_20261026T111302Z_6e46ba39379a80ff`
- Source commit: `1861dbb464ae536691032a8eadaa181ad0cc58e2`
