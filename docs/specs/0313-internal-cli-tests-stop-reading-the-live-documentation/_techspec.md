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

A corrective Task, added after the first Run's QA gate (F-01), gives the
upgrade characterization test a fixture workspace: it was a fourteenth
reader, which inspected the live repository through `roundfix upgrade`'s
history notice. A second corrective Task, at the maintainer's request,
moves the suite guard's fingerprint to Git subprocesses and keeps every
detection class with a mode-only `os.Lstat` pass.

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

Measured on 2026-10-25 at the first Run's head `5299fcf6`, with task_01 and
task_02 completed, for the corrective Task:

6. `TestUpgradeStdoutAndExitCodesAreCharacterized`
   (`internal/cli/upgrade_characterization_test.go`) runs `roundfix upgrade`
   without `withCLIWorkspace`, so the command's working directory is the
   package directory inside this repository. After the upgrade,
   `printHistoryNotice` (`internal/cli/upgrade.go`) loads the Project Config
   and runs `historyInventory` on the repository root it finds. A run of
   that test alone with `-test.testlogfile` logs open and stat entries under
   `docs/history`, `docs/adr`, `docs/specs`, `docs/user-guide`,
   `docs/references` and `.roundfixrc.yml`. The test asserts only stdout and
   the exit code, so it passes with the trees unreadable; that is why the
   audit's unreadable-trees probe missed it. It is the only test of its file
   without a fixture workspace.
7. `internal/suiteguard` is not a cache input. `suiteguard.Main` digests
   the repository before and after `m.Run`, and the testing package logs
   file access only inside `m.Run`. A fixture test that ran a guarded
   package with `-test.testlogfile` found no entry under the fixture
   repository. With the line of item 6 added, a full run of the untagged
   test binary logged only `go.mod`, `internal/go.mod` and
   `internal/speccheck/testdata` inside the module. A line appended to
   `docs/user-guide/commands/gc.md` and `.agents/skills/roundfix/SKILL.md`
   then left the package `(cached)`.
8. `internal/suiteguard.fingerprint` walks the root with `filepath.WalkDir`,
   skips `.git` and paths a hand-written parser of the `.gitignore` files
   ignores, and records each remaining entry's `fs.FileMode` and the
   SHA-256 of its content or link target. A characterization table of 22
   leak cases, each run inside a guarded subprocess on a committed Git
   fixture, passes on it: tracked edits, deletions, an emptied file, a
   second edit of a dirty file, untracked creations, edits and deletions,
   empty-directory creation and removal, executable bits, symbolic-link
   retargeting and type change, non-executable permission changes of a file
   and a directory, a new directory holding only ignored files, a negated
   ignore rule, and three ignored writes that pass. The guard accepts a root
   that is not a Git work tree.

## System Architecture

- `internal/cli/repository_profile_docscontract_test.go` (new, task_01) and
  `internal/cli/live_documents_docscontract_test.go` (new, task_02): the
  moved tests.
- The seven source test files lose those functions; helpers that other tests
  use stay where they are.
- `internal/verifyselect/contracts_repository_test.go`: an `internal/cli`
  `docscontract` case and a pinned file count.
- `docs/agents/specific-repository.md`: the rule.
- `internal/cli/upgrade_characterization_test.go` (task_04): one
  `withCLIWorkspace(t)` line.
- `internal/suiteguard/suiteguard.go` (task_05): `fingerprint(root)` through
  Git and a mode pass; the ignore parser and `digestFile` go away.
  `internal/suiteguard/leak_cases_test.go` (new) holds the table;
  `suiteguard_test.go` and `residue_test.go` make their fixtures Git work
  trees.

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
3. API Contract: `suiteguard.Main(m, repoRoot)` keeps its signature, its
   messages and its exit codes on a Git work tree. The fingerprint is one
   state per watched path, a full `fs.FileMode` and a Git mode plus object
   ID, taken from:
   - `git ls-files -z --cached --others --exclude-per-directory=.gitignore`;
   - `git ls-files -z --others --directory --exclude-per-directory=.gitignore`
     for untracked directories;
   - `git update-index --add --remove --info-only -z --stdin` into a
     temporary `GIT_INDEX_FILE` outside the repository, read back with
     `git ls-files -z --stage`;
   - one `os.Lstat` per watched path and directory.

   Git runs without inherited `GIT_` variables and with
   `GIT_OPTIONAL_LOCKS=0`. A root outside a Git work tree fails with
   `suiteguard: fingerprint repository before tests:`.

### Surface Transcripts

None. No command surface changes.

### Invariants

1. No untagged `internal/cli` test opens the five documentation trees.
2. Each of the thirteen tests passes under `-tags docscontract` and fails on
   the inputs it failed on before.
3. Every `internal/cli` `docscontract` file is always selected.
4. The suite guard reports every leak class of Current behavior 8, stays
   silent on ignored writes, and leaves no repository path in a guarded
   package's test input log.

## Coverage Map

- Goal 1, Success Metrics 3 and 4 → Invariant 1, task_01, task_02, task_04,
  task_03.
- Goal 2, Core Feature 4, Success Metric 5 → Invariant 1, task_04, task_03.
- Goal 3, Core Feature 5, Success Metric 6 → API Contract 3, Invariant 4,
  task_05, task_03.
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
- task_04 proves its line with Go's test input log of the one test it
  changes. The QA gate runs the whole untagged test binary with
  `-test.testlogfile` (Success Metric 5), the measurement that found the
  fourteenth reader. A pass under unreadable trees does not prove that a
  test read nothing.
- task_05 writes the leak-case table first and runs it on the unchanged
  guard. The same table, the input-log test and the refusal test prove the
  Git fingerprint, and the QA gate sabotages each Git listing and the mode
  pass.
- The cache behavior and the unreadable-trees run are QA rows (task_03):
  each needs a full `internal/cli` run of three to five minutes, which a Task
  Verification should not spend twice.

## Build Order

1. Move the five repository-profile tests and classify the file (task_01).
2. Move the other eight, raise the count, write the guide rule (task_02)
   (depends on: 1, because both change
   `internal/verifyselect/contracts_repository_test.go`).
3. Give the upgrade characterization test a fixture workspace (task_04)
   (depends on: 2, so it runs alone after both moves; it shares no file with
   them).
4. Fingerprint the repository through Git in the suite guard (task_05)
   (depends on: 3, so it runs alone after task_04; it shares no file with
   the earlier Tasks).
5. The QA gate (task_03) (depends on: 1, 2, 3, 4).

No other Spec is required.

## Risks & Considerations

- A helper used only by a moved test becomes unused in the untagged build;
  Go accepts an unused function, and the Task moves such a helper with its
  test when no other test uses it.
- A pending plan that edits one of the thirteen bodies must follow the test
  to its new file; none of 0307, 0309 and 0310 does at authoring time.
- `make test` no longer runs the thirteen; `make verify-changed` and the Full
  Contract Run do.
- The guard now needs `git` and a Git work tree at the root. Checkouts,
  Run worktrees, CI and the Linux Lane (which runs `git init` and commits
  the candidate) all provide one; a `git archive` copy must be initialized
  and committed first, as the QA rows already do.
- A future test that runs a command in its package directory can read the
  live repository again. The repository rule of task_02 names the trees, and
  the QA test input log row is the way to find such a reader.

## Glossary

None.

## Decisions

- Move, not rewrite (ADR-0326).
- Pin the file count in the declarations test, so a third file is a
  deliberate change.
- Corrective (2026-10-25): fix the fourteenth reader (task_04). The guard's
  walk is outside Go's logging window (Current behavior 7), so the cache
  depends on task_04 alone.
- Corrective (2026-10-25, the maintainer's "As duas coisas"): the guard also
  fingerprints through Git (task_05). A mode-only `os.Lstat` pass keeps
  non-executable permission changes detectable; it runs outside `m.Run` and
  so is not a cache input. `--exclude-per-directory=.gitignore` rather than
  `--exclude-standard` keeps today's ignore scope: the user's global
  excludes file, which ignores `dist/` on the maintainer's machine, does not
  hide untracked paths from the guard. `--info-only` hashes
  without writing objects into the repository.
