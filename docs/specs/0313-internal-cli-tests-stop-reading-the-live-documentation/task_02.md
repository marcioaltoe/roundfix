---
task: task_02
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
status: pending
type: test
complexity: medium
---

# Task 02: The other live-document tests become documentation contracts

## Overview

Eight more `internal/cli` tests read this repository's guides, skills and
records: the skill set it holds, the examples and HTTP statement its guides
carry, the events help, and the release plan checks run on it. This Task
moves them, unchanged, into a second `docscontract` file declared
`//verify:always`, raises the pinned file count to two, and writes the rule
of ADR-0326 into the repository guide.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-24,
   [`references/2026-10-24-internal-cli-tests-read-the-live-documentation-trees.md`](references/2026-10-24-internal-cli-tests-read-the-live-documentation-trees.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1 and 2 of
   `_techspec.md` for its second group.
2. MUST move `TestARepositoryMissingARequiredExternalSkillIsReported` and
   `TestThisRepositoryHoldsEveryRequiredExternalSkill` from
   `this_repository_skill_set_test.go`, `TestBaselineExamplesParse` and
   `TestProjectConstraintPlanningExamplesParse` from
   `baseline_documentation_contract_test.go`,
   `TestEventsHelpDocumentsAgentSelectionFilter` from `cli_test.go`,
   `TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault` from
   `baseline_http_default_statement_test.go`, and
   `TestReleasePlanChecksOnAnAdoptedRepository` and
   `TestReleasePlanChecksWriteNothing` from `releaseplan_checks_test.go`, all
   under `internal/cli/`, with their doc comments and bodies byte for byte,
   into the new file `internal/cli/live_documents_docscontract_test.go`,
   with the same header shape as task_01's file. A helper that only moved
   tests use moves with them; every other test and helper keeps its bytes.
3. MUST change the expected number of `internal/cli` documentation contract
   files in `TestRepositoryContractTestsDeclareTheirRelevance` from one to
   two, and nothing else in that file.
4. MUST add one bullet to `## Repository rules` in
   `docs/agents/specific-repository.md`, after the HARD RULE on repository
   contracts: "A test in a code package that reads this repository's
   `docs/history`, `docs/adr`, `docs/specs`, `docs/user-guide` or
   `.agents/skills` lives in a `docscontract` file declared
   `//verify:always`, so a documentation edit never invalidates the
   package's cached test result (ADR-0326)." The rest of the file keeps its
   bytes.
5. MUST NOT change a test body, the Makefile, the CI workflows, any
   production file or `internal/docscontract`.

## Subtasks

- [ ] Move the eight tests into the second `docscontract` file.
- [ ] Raise the pinned file count and write the guide rule.

## Acceptance Criteria

- [ ] Without the tag none of the thirteen tests exists; with it all thirteen pass.
- [ ] The declarations test pins two always-selected `internal/cli` files, and the guide states the rule.

## Consumers

None. Test code moves inside one package and one repository test changes a count; no production contract changes.

## Context

- instruction: `docs/adr/0326-a-code-test-never-reads-the-live-documentation-trees.md`
- instruction: `docs/references/2026-10-24-test-suite-audit.md`
- creates: `internal/cli/live_documents_docscontract_test.go`
- interface: `internal/cli/this_repository_skill_set_test.go`
- interface: `internal/cli/baseline_documentation_contract_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/baseline_http_default_statement_test.go`
- interface: `internal/cli/releaseplan_checks_test.go`
- interface: `internal/verifyselect/contracts_repository_test.go`
- interface: `docs/agents/specific-repository.md`

## Verification

- `listed="$(go test -list '^(TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile|TestBaselineUpdateReachesCurrentWithARepositoryProfile|TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile|TestDoctorComparesARepositoryProfileWithItsSnapshot|TestBaselineUpdatePreviewListsASetupSnapshotDifference|TestARepositoryMissingARequiredExternalSkillIsReported|TestThisRepositoryHoldsEveryRequiredExternalSkill|TestBaselineExamplesParse|TestProjectConstraintPlanningExamplesParse|TestEventsHelpDocumentsAgentSelectionFilter|TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault|TestReleasePlanChecksOnAnAdoptedRepository|TestReleasePlanChecksWriteNothing)$' ./internal/cli 2>&1)" || { printf '%s\n' "$listed"; exit 1; }; if printf '%s\n' "$listed" | grep -q '^Test'; then printf 'still untagged: %s\n' "$listed" >&2; exit 1; fi; out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestARepositoryMissingARequiredExternalSkillIsReported|TestThisRepositoryHoldsEveryRequiredExternalSkill|TestBaselineExamplesParse|TestProjectConstraintPlanningExamplesParse|TestEventsHelpDocumentsAgentSelectionFilter|TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault|TestReleasePlanChecksOnAnAdoptedRepository|TestReleasePlanChecksWriteNothing)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestARepositoryMissingARequiredExternalSkillIsReported TestThisRepositoryHoldsEveryRequiredExternalSkill TestBaselineExamplesParse TestProjectConstraintPlanningExamplesParse TestEventsHelpDocumentsAgentSelectionFilter TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault TestReleasePlanChecksOnAnAdoptedRepository TestReleasePlanChecksWriteNothing; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -timeout 15m -run '^(TestRepositoryContractTestsDeclareTheirRelevance|TestNoContractNameCrossesABuildTagClass)$' ./internal/verifyselect && tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- 'never invalidates the package' && GOOS=windows go vet -tags 'repocontract docscontract' ./internal/cli` — expected: exit 0. Before this Task the untagged list prints the eight tests, so the command fails at its first check; with one of the eight left untagged the list check fails, and with the count left at one the declarations test fails. Afterwards none of the thirteen exists without the tag, the eight pass under it, the declarations test passes, and the guide states the rule.

## References

- `_prd.md` → Sources; Goals; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior; API Contracts 1-2; Invariants 1-3; Testing Approach; Build Order 2
