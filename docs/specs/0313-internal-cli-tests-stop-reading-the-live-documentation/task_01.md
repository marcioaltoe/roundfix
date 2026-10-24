---
task: task_01
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
status: pending
type: test
complexity: medium
---

# Task 01: The repository-profile tests become documentation contracts

## Overview

Five `internal/cli` tests run Baseline and Doctor commands on this
repository's own profile and so read its guides and skills; any edit there
re-runs the whole package. This Task moves them, unchanged, into a new
`docscontract` file declared `//verify:always` (ADR-0326), and teaches the
repository contract declarations test that such a file is always selected.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-24,
   [`references/2026-10-24-internal-cli-tests-read-the-live-documentation-trees.md`](references/2026-10-24-internal-cli-tests-read-the-live-documentation-trees.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1 and 2 of
   `_techspec.md` for its first group.
2. MUST move `TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile`,
   `TestBaselineUpdateReachesCurrentWithARepositoryProfile`,
   `TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile` and
   `TestDoctorComparesARepositoryProfileWithItsSnapshot` from
   `internal/cli/baseline_update_repository_profile_test.go`, and
   `TestBaselineUpdatePreviewListsASetupSnapshotDifference` from
   `internal/cli/doctor_trailing_skills_test.go`, with their doc comments and
   bodies byte for byte, into the new file
   `internal/cli/repository_profile_docscontract_test.go`, which starts with
   a Suite, Invariant, Boundary IN and Boundary OUT header, a comment that
   these tests read the live documentation trees and run under the
   `docscontract` tag per ADR-0326, then `//go:build docscontract`, a blank
   line, `//verify:always`, a blank line and `package cli`. A helper that
   only moved tests use moves with them; every other test and helper keeps
   its bytes, and the two source files keep compiling without the tag.
3. MUST add to `TestRepositoryContractTestsDeclareTheirRelevance` in
   `internal/verifyselect/contracts_repository_test.go` a case before the
   `internal/suiteguard` case: a contract whose package is `internal/cli`
   and whose tag is `docscontract` wants `verifyselect.ContractAlways`, with
   a comment citing ADR-0326, and record its file; after the loop, fail
   unless exactly one such file was seen. Every other expectation keeps its
   bytes.
4. MUST NOT change a test body, the Makefile, the CI workflows, any
   production file or `internal/docscontract`.

## Subtasks

- [ ] Move the five tests into the new `docscontract` file.
- [ ] Classify the file in the repository contract declarations test.

## Acceptance Criteria

- [ ] Without the tag the five tests do not exist; with it they pass.
- [ ] The declarations test classifies the new file as always selected.

## Consumers

None. Test code moves inside one package and one repository test learns its class; no production contract changes.

## Context

- instruction: `internal/docscontract/corpus_test.go`
- instruction: `internal/verifyselect/contracts.go`
- instruction: `docs/adr/0326-a-code-test-never-reads-the-live-documentation-trees.md`
- instruction: `docs/references/2026-10-24-test-suite-audit.md`
- creates: `internal/cli/repository_profile_docscontract_test.go`
- interface: `internal/cli/baseline_update_repository_profile_test.go`
- interface: `internal/cli/doctor_trailing_skills_test.go`
- interface: `internal/verifyselect/contracts_repository_test.go`

## Verification

- `listed="$(go test -list '^(TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile|TestBaselineUpdateReachesCurrentWithARepositoryProfile|TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile|TestDoctorComparesARepositoryProfileWithItsSnapshot|TestBaselineUpdatePreviewListsASetupSnapshotDifference)$' ./internal/cli 2>&1)" || { printf '%s\n' "$listed"; exit 1; }; if printf '%s\n' "$listed" | grep -q '^Test'; then printf 'still untagged: %s\n' "$listed" >&2; exit 1; fi; out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile|TestBaselineUpdateReachesCurrentWithARepositoryProfile|TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile|TestDoctorComparesARepositoryProfileWithItsSnapshot|TestBaselineUpdatePreviewListsASetupSnapshotDifference)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile TestBaselineUpdateReachesCurrentWithARepositoryProfile TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile TestDoctorComparesARepositoryProfileWithItsSnapshot TestBaselineUpdatePreviewListsASetupSnapshotDifference; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -timeout 15m -run '^(TestRepositoryContractTestsDeclareTheirRelevance|TestNoContractNameCrossesABuildTagClass)$' ./internal/verifyselect && GOOS=windows go vet -tags 'repocontract docscontract' ./internal/cli` — expected: exit 0. Before this Task the untagged list prints the five tests, so the command fails at its first check; with the file's `//verify:always` removed the declarations test fails on the class. Afterwards the five run and pass only under the tag, the declarations test passes, and the package vets for Windows with both contract tags.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Core Feature 3; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior; API Contracts 1-2; Invariants 1-3; Testing Approach; Build Order 1
