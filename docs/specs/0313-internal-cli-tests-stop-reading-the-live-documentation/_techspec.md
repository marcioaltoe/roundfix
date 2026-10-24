---
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
prd: _prd.md
created: 2026-10-24
---

# internal/cli tests stop reading the live documentation — Technical Spec

## Executive Summary

Go records every file a test process opens inside the module and reuses the
package's result only while those files are unchanged. Thirteen
`internal/cli` tests open this repository's guides, skills and records, so
any edit there re-runs the whole package. This Spec moves those thirteen
test functions, with their doc comments and unchanged bodies, into two new
`internal/cli` files under `//go:build docscontract` with `//verify:always`
(ADR-0326). The selective gate's discovery already finds `docscontract`
tests in any package and runs them with `-tags docscontract`; the repository
contract declarations test learns that `internal/cli` documentation contract
files are always selected.

The trade-off accepted: the thirteen leave `make test` and the push-to-main
"Verify" step, and run in `make verify-changed` and `make verify-contracts`,
compiling `internal/cli` once more with the tag.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0326, as `_prd.md` → Project
  Constraints states. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — protected paths limited to the
  repository guide line and the two governed test files the tests leave. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0313-internal-cli-tests-stop-reading-the-live-documentation/_authorization.md`;
  bounded files: `docs/agents/specific-repository.md`,
  `internal/cli/cli_test.go`,
  `internal/cli/baseline_documentation_contract_test.go`.

## Current behavior

Read and characterized on 2026-10-24 at `56dd1a76`.

1. The thirteen tests and their files:
   `baseline_update_repository_profile_test.go`
   (`TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile`,
   `TestBaselineUpdateReachesCurrentWithARepositoryProfile`,
   `TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile`,
   `TestDoctorComparesARepositoryProfileWithItsSnapshot`);
   `doctor_trailing_skills_test.go`
   (`TestBaselineUpdatePreviewListsASetupSnapshotDifference`);
   `this_repository_skill_set_test.go`
   (`TestARepositoryMissingARequiredExternalSkillIsReported`,
   `TestThisRepositoryHoldsEveryRequiredExternalSkill`);
   `baseline_documentation_contract_test.go` (`TestBaselineExamplesParse`,
   `TestProjectConstraintPlanningExamplesParse`); `cli_test.go`
   (`TestEventsHelpDocumentsAgentSelectionFilter`);
   `baseline_http_default_statement_test.go`
   (`TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault`);
   `releaseplan_checks_test.go` (`TestReleasePlanChecksOnAnAdoptedRepository`,
   `TestReleasePlanChecksWriteNothing`). Together about 150 s of test time.
2. With `docs/history`, `docs/adr`, `docs/specs`, `docs/user-guide` and
   `.agents/skills` removed, and again unreadable, exactly these thirteen
   fail.
3. In a worktree whose result was cached, a line appended to a command
   guide re-ran `internal/cli` (162.6 s, then 237.0 s); with the thirteen
   skipped it stayed `(cached)`. A Task file's `status:` change left both
   forms cached.
4. `internal/verifyselect` `DiscoverContracts` finds tests under
   `docscontract` in every package; `TestRepositoryContractTestsDeclareTheirRelevance`
   expects every contract outside a short list to be `package`-relevant,
   and nine `internal/docscontract` files to be `always`.
5. The Parallel Test Package rule excludes `docscontract` files, and
   `TestEveryTestCompilesForWindows` vets every package with both contract
   tags for Windows.

## System Architecture

- `internal/cli/repository_profile_docscontract_test.go` (new, task_01) and
  `internal/cli/live_documents_docscontract_test.go` (new, task_02): the
  moved tests.
- The seven source test files lose those functions; helpers that other tests
  use stay where they are.
- `internal/verifyselect/contracts_repository_test.go`: an `internal/cli`
  `docscontract` case and a pinned file count.
- `docs/agents/specific-repository.md`: the rule.

## Implementation Design

### Interfaces

```go
// internal/cli/<group>_docscontract_test.go
// Suite / Invariant / Boundary header, a comment naming ADR-0326, then:
//go:build docscontract

//verify:always

package cli
```

### Data Models

None. No persisted format changes.

### API Contracts

1. API Contract: each moved test keeps its name, doc comment and body byte
   for byte; only its file and imports change.
2. API Contract: `TestRepositoryContractTestsDeclareTheirRelevance` expects
   `ContractAlways` for every contract in `internal/cli` whose tag is
   `docscontract`, and exactly the number of such files the Task reaches (1
   after task_01, 2 after task_02).

### Surface Transcripts

None. No command surface changes.

### Invariants

1. No untagged `internal/cli` test opens the five documentation trees.
2. Each of the thirteen tests passes under `-tags docscontract` and fails on
   the inputs it failed on before.
3. Every `internal/cli` `docscontract` file is always selected.

## Coverage Map

- Goal 1, Success Metrics 3 and 4 → Invariant 1, task_01, task_02, task_03.
- Goal 2, Core Features 1 and 2, Success Metric 1 → API Contract 1,
  Invariant 2, task_01, task_02.
- Goal 3, Core Feature 3, Success Metric 2 → API Contract 2, Invariant 3,
  task_01, task_02.

Deviations from the PRD: None.

## Absence Evidence

1. An `internal/cli` documentation contract — command:
   `grep -ln "go:build docscontract" internal/cli/*_test.go`; result: no
   match.

## ADR Alternatives

- ADR-0326: this Spec implements it. Its rejected options: fixture copies for
  every test (they would stop checking what ships), a `-skip` list in the
  Makefile (it hides the tests from gates that do not repeat it), and moving
  them to `internal/docscontract` (they exercise unexported `internal/cli`
  code).

## Integration Points

None. No external system, network or provider is reached.

## Testing Approach

- The moved tests are their own proof: each passes under the tag; the
  untagged `go test -list` of their names prints none of them.
- `TestRepositoryContractTestsDeclareTheirRelevance` proves the class and the
  count.
- The cache behavior and the unreadable-trees run are QA rows (task_03):
  each needs a full `internal/cli` run of three to five minutes, which a Task
  Verification should not spend twice.

## Build Order

1. Move the five repository-profile tests and classify the file (task_01).
2. Move the other eight, raise the count, write the guide rule (task_02)
   (depends on: 1, because both change
   `internal/verifyselect/contracts_repository_test.go`).
3. The QA gate (task_03) (depends on: 1, 2).

No other Spec is required.

## Risks & Considerations

- A helper used only by a moved test becomes unused in the untagged build;
  Go accepts an unused function, and the Task moves such a helper with its
  test when no other test uses it.
- A pending plan that edits one of the thirteen bodies must follow the test
  to its new file; none of 0307, 0309 and 0310 does at authoring time.
- `make test` no longer runs the thirteen; `make verify-changed` and the Full
  Contract Run do.

## Glossary

None.

## Decisions

- Move, not rewrite (ADR-0326).
- Pin the file count in the declarations test, so a third file is a
  deliberate change.
