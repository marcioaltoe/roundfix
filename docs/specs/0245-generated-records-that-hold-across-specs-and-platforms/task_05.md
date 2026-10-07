---
task: task_05
spec: 0245-generated-records-that-hold-across-specs-and-platforms
status: completed
type: chore
complexity: medium
---

# Task 05: The derived declarations cover every file a module edit regenerates

## Overview

Corrective Task for QA finding F1 (row Q04, Blocks-Completion) in
`qa/qa-report-2026-10-07.md`. Two branches that each edit a different clause
of `context-workflow.json`, and each run the Module Version Record step,
`make baseline-digests` and the Setup Manifest refresh, conflict at merge in
three places no Derived Path Declaration covers:

- the `formatter.goldenDigest` line of
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`;
- the formatter golden fixture
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md`;
- `docs/agents/domain.md`, which the Managed Refresh writes.

ADR-0192 refuses a conflict outside declared paths, so the item parks instead
of reaching the next free module version. This Task declares those outputs
and adds a test that proves every file a module edit regenerates is
declared, so the next uncovered output fails in the repository suite rather
than at a merge.

## Requirements

1. MUST extend `delivery.derived_paths` in `.roundfixrc.yml` so these
   outputs are declared. Use whole-path declarations only for files that
   are entirely generated: the formatter golden fixture tree, and every
   `docs/agents/` file the Setup Manifest marks as managed. Use a
   line-scoped declaration for the profile's `goldenDigest` line, because
   the profile also holds hand-written decisions. First read how
   `internal/config` matches `paths` and `lines` patterns, and use only the
   forms it supports. Keep the declaration order: the module record and
   digests run before the Setup Manifest refresh.
2. MUST add a repository test under the `docscontract` build tag. It copies
   the repository to a temporary Git checkout and edits one guidance clause
   of a Baseline module. It then runs the declared regeneration commands in
   declaration order and asserts that every path Git reports as changed
   matches a declaration. A line-scoped path counts as declared only when
   every changed line matches its `match` pattern. The test MUST fail when
   any declaration this Task adds is removed.
3. MUST update `TestThisRepositoryDeclaresItsToolsAndDerivedPaths` in
   `internal/config` so it pins the new declarations.
4. MUST re-record the Coverage Record through its declared command after
   adding the test, and MUST NOT edit the record by hand.
5. MUST NOT change a module's content, a profile's decisions or any
   production Go file outside `internal/config` tests and the new test.

## Subtasks

- [ ] Map every path the module-edit regeneration writes, in a scratch copy.
- [ ] Declare those paths in `.roundfixrc.yml`.
- [ ] Add the regeneration-coverage test and prove it fails when one new declaration is removed.
- [ ] Pin the declarations in the config test and re-record the Coverage Record.

## Acceptance Criteria

- [ ] The new test passes, and fails when any declaration this Task adds is removed.
- [ ] The config test pins the new declarations.
- [ ] `make verify` and `make verify-docs` exit 0.

## Verification

- `go test -count=1 -tags docscontract ./internal/config/... ./internal/baseline/... -run 'RegenerationIsDeclared' -v 2>&1 | grep -q -- '--- PASS: Test[A-Za-z]*RegenerationIsDeclared' || exit 1`
- `grep -qF 'goldenDigest' .roundfixrc.yml || exit 1; grep -qF 'formatter-fixtures' .roundfixrc.yml || exit 1; go test -count=1 ./internal/config -run '^TestThisRepositoryDeclaresItsToolsAndDerivedPaths$' || exit 1`

## References

- `qa/qa-report-2026-10-07.md` → F1; `qa/evidence/2026-10-07/Q04.txt`
- `_prd.md` → Requirement 4
- ADR-0250; ADR-0233; ADR-0192

## Context

- interface: `.roundfixrc.yml`
- interface: `internal/config/verification_tools_test.go`

## Result

Implementation evidence:

- Extended `delivery.derived_paths` in declaration order with the profile
  `goldenDigest` line, the complete formatter golden fixture tree, and the
  managed `docs/agents/` tree alongside the existing Setup Manifest output.
- Added `TestRegenerationIsDeclared` under the `docscontract` build tag. It
  copies the repository into a temporary Git checkout, edits a context module
  guidance clause, runs every declared regeneration command, checks changed
  paths and changed lines, and runs removal ablations for each declaration
  added by this Task.
- Updated `TestThisRepositoryDeclaresItsToolsAndDerivedPaths` with the exact
  declarations and re-recorded `docs/references/coverage-record.json` through
  its declared Coverage Record command.

Focused checks:

- `GOCACHE=/tmp/roundfix-task05-gocache go test -tags docscontract ./internal/config -run '^TestRegenerationIsDeclared$' -count=1` — PASS; the regeneration test and all three declaration-removal ablations passed.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 ./internal/config -run '^TestThisRepositoryDeclaresItsToolsAndDerivedPaths$'` — PASS.
- `git diff --check` — PASS.

Acceptance evidence:

- New regeneration coverage and declaration-removal behavior are exercised by
  the focused passing test; the Daemon still owns the authored Verification.
- The config pin test passed with the new declarations.
- The Coverage Record contains `TestRegenerationIsDeclared`, generated by the
  declared recorder command.
- `make verify` and `make verify-docs` remain for Daemon Verification.

Feedback repair evidence:

- Settlement reported `SC-SPEC-PATH-PINNED` because the scratch authorization
  helper named the active Spec directory. The helper now creates its
  regeneration authorization under a generic temporary legacy authorization
  path, so the test has no dependency on an archivable Spec directory.
- `GOCACHE=/tmp/roundfix-task05-feedback-gocache go test -tags docscontract ./internal/config -run '^TestRegenerationIsDeclared$' -count=1 -failfast` — PASS.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/config/regeneration_declared_test.go`

## Carry-forward provenance

- Source Run: `run_20261007T130416Z_67a7acccb80313ba`
- Source commit: `bb8ae254de8d6d6c4fe78be93c47b3e0b803c1d4`
