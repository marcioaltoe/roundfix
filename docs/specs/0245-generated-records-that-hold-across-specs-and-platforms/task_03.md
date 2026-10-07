---
task: task_03
spec: 0245-generated-records-that-hold-across-specs-and-platforms
status: completed
type: chore
complexity: medium
---

# Task 03: The derived declarations run both record steps, and the rule and glossary name them

## Overview

task_01 and task_02 give each record a command that computes it. This Task
makes a Pull Request conflict run those commands: the Baseline catalog
declaration runs the module record step before `make baseline-digests` and
covers each module's version line, and the Coverage Record gains its own
declaration. It writes the repository rule an author follows after a module
edit, adds the two glossary terms, and re-records the Coverage Record so it
holds the tests task_01 added.

## Requirements

1. MUST change the first `delivery.derived_paths` entry of `.roundfixrc.yml`
   to the declaration of `_techspec.md` → Data Models: the record path first
   in `paths`, the line-scoped module version lines with
   `match: '^  "version": [0-9]+,$'`, and the regeneration
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1 && make baseline-digests`.
   MUST add the Coverage Record declaration last. The Setup Manifest and
   owned-skill declarations keep their bytes and order.
2. MUST change `TestThisRepositoryDeclaresItsToolsAndDerivedPaths` in
   `internal/config/verification_tools_test.go` to expect exactly those four
   declarations, the one declared break there.
3. MUST add, after the owned-skill rule in `docs/agents/specific-repository.md`,
   this rule, with every other line kept:
   "- A Baseline module's content changes only together with its version.
   After editing a module under `internal/baseline/assets/modules/`, run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   before `make baseline-digests`; the record command keeps a version above
   every recorded one or writes the next free version into the module's own
   version line, records it in `internal/baseline/module-versions.json`, and
   requires that record to be declared in the Task. A Spec names the next
   version, never a number. `docs/references/coverage-record.json` is
   re-recorded only with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`,
   which writes the same bytes on any host."
4. MUST add **Module Version Record** and **Coverage Record** to `CONTEXT.md`
   through `domain-modeling`, after **Derived Path Declaration**, each with an
   `_Avoid_` line and ADR-0250 cited, and never citing a Spec.
5. MUST re-record `docs/references/coverage-record.json` with the command of
   Requirement 3 after task_01 and task_02, so it lists
   `TestEveryBaselineModuleVersionIsRecorded`.

## Subtasks

- [ ] Change the Baseline catalog declaration and add the Coverage Record declaration.
- [ ] Update the repository-config test.
- [ ] Write the repository rule and the two glossary terms.
- [ ] Re-record the Coverage Record.

## Acceptance Criteria

- [ ] The Project Config declares the record step before the digests, the
      module version lines and the Coverage Record, and the config test pins
      exactly those declarations.
- [ ] The repository rule tells an author to run the record step and name
      the next version, and `CONTEXT.md` defines both terms.

## Context

- interface: `.roundfixrc.yml`
- interface: `internal/config/verification_tools_test.go`
- interface: `docs/agents/specific-repository.md`
- interface: `CONTEXT.md`
- interface: `docs/references/coverage-record.json`
- instruction: `docs/user-guide/configuration.md`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`

## Verification

- `grep -qF -- "-record-module-versions -count=1 && make baseline-digests" .roundfixrc.yml || exit 1; grep -qF -- "-update-coverage-record -count=1" .roundfixrc.yml || exit 1; go test -count=1 ./internal/config -run '^TestThisRepositoryDeclaresItsToolsAndDerivedPaths$' || exit 1; tr -s '[:space:]' ' ' < docs/references/coverage-record.json | grep -qF -- '"TestEveryBaselineModuleVersionIsRecorded"'` — expected: exit 0; before this Task `.roundfixrc.yml` names neither record step, so the command fails.
- `for phrase in '**Module Version Record**:' '**Coverage Record**:'; do tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "$phrase" || { printf 'missing glossary term: %s\n' "$phrase" >&2; exit 1; }; done; for phrase in "-record-module-versions -count=1" "A Spec names the next version, never a number." "which writes the same bytes on any host"; do tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- "$phrase" || { printf 'missing rule phrase: %s\n' "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task neither term nor the rule exists, so the command fails.

## References

- `_prd.md` → Goal 2; Core Feature 3; Core Feature 4; Success Metric 3
- `_techspec.md` → Data Models; Invariant 4; Build Order 3
- ADR-0250; ADR-0233; ADR-0192

## Result

Implemented the Task 03 slice: the Project Config now declares the Baseline
module version record before the Baseline digests, scopes its module version
lines, keeps the Setup Manifest and owned-skill declarations in their existing
order, and declares the Coverage Record last. The repository-config test pins
the resulting four declarations. The repository rule and the two CONTEXT
glossary terms describe the record steps and cite ADR-0250.

The required Coverage Record generator re-recorded the existing platform
matrix and added the Task 01 baseline tests, including
`TestEveryBaselineModuleVersionIsRecorded`. The generated file was not hand
edited.

Acceptance evidence:

- Project Config and declaration shape: direct diff inspection confirms the
  record path is first, the module line matcher is
  `^  "version": [0-9]+,$`, the record command precedes `make baseline-digests`,
  and the Coverage Record is the fourth declaration. The focused config
  declaration and validation tests passed.
- Repository guidance and vocabulary: direct diff inspection confirms the
  rule follows the owned-skill rule, names the next-version behavior, and
  `CONTEXT.md` defines **Module Version Record** and **Coverage Record** after
  **Derived Path Declaration**, each with an `_Avoid_` line and ADR-0250.

Focused checks:

- `GOCACHE=/private/tmp/roundfix-0245-task03-gocache go test ./internal/config -run '^(TestDerivedLineDeclarationsAreReadAndValidated|TestVerificationToolsAreBareExecutableNames)$' -count=1` — exit 0.
- `GOCACHE=/private/tmp/roundfix-0245-task03-gocache go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1` — the generator wrote the 11 missing baseline test entries, including `TestEveryBaselineModuleVersionIsRecorded`; the package then exited 1 because suiteguard reported the intentional modification of `docs/references/coverage-record.json`.
- `gofmt -w internal/config/verification_tools_test.go` — completed.
- `git -c core.fsmonitor=false diff --check` — exit 0.
- `rg` inspection of `docs/references/coverage-record.json` confirmed the
  platform field, the Darwin/Linux/Windows platform-specific tests, and the
  newly generated `TestEveryBaselineModuleVersionIsRecorded` entry.

Daemon Verification remains pending; no command from this Task's authored
Verification section was run. Task status and checkboxes remain Daemon-owned.
No commit, push, or pull request was made.

## Carry-forward provenance

- Source Run: `run_20261007T115541Z_03741c3fcc656349`
- Source commit: `c78cfc85c9055790e29bde9a3321cc4c522bc317`
