---
spec: 0420-roundfix-stops-verifying-on-linux
prd: _prd.md
created: 2026-11-24
---

# Roundfix stops verifying on Linux — Technical Spec

## Executive Summary

The design removes code and adds almost none. The Delivery Queue's repository gate ends when the repository Verification passes, the lane's two source files, its command, its doctor check and its 36 tests leave, and the two keys move to the registry of deprecated keys that the configuration loader already has. The job `Verification gate` and the release's verification run on the runner `macos-26`, and the Repository Contract Test that compiles every test for Windows compiles them for Linux as well. No table, column, stage, command, flag, key or exit code is added. The measurements behind the design are in [`references/2026-11-24-techspec-measurements.md`](references/2026-11-24-techspec-measurements.md).

The primary trade-off is the one the maintainer's decision accepts: after this Spec no gate runs a test on Linux or on Windows, so a defect that only one of them produces is found by nobody. The design keeps the one proof that costs seconds and runs nothing there: the code and the tests of both platforms compile before a pull request merges.

Two smaller trade-offs follow from where the verification runs. The release verifies in a job of its own that the publishing job needs, so the verification holds no publishing permission and a failed publication is re-run without it; the cost is that a wrong tag is refused after the verification and not before. The suite-time budget keeps its number, measured on Linux, until the first push to main measures the suite on the new runner, because a Task reaches no runner and the number changes only with a measurement.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. The decision record takes the number 0452, the next free one. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added, and the product stops calling the local `docker` command. The release's new job holds `contents: read` and nothing else; the publishing job keeps its permissions and its Trusted Publishing. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the PRD's row stands. ADR-0452 is the record a Task writes; it supersedes ADR-0288, ADR-0294, ADR-0297, ADR-0303 and ADR-0316 and amends in part ADR-0252, ADR-0285, ADR-0295, ADR-0384, ADR-0385 and ADR-0423. ADR-0031, ADR-0236, ADR-0250, ADR-0253, ADR-0256, ADR-0304, ADR-0349, ADR-0365, ADR-0387, ADR-0398 and ADR-0450 are the ones the design obeys or uses. ADR-0056, ADR-0259, ADR-0279, ADR-0311, ADR-0369, ADR-0395 and ADR-0406 are neighbours listed only to record why they do not apply. ADR-0034, ADR-0148, ADR-0271, ADR-0373, ADR-0412 and ADR-0435 are listed for the same purpose, each because it cites one of those neighbours or ADR-0349; this stage added them to the PRD's row. `## ADR Alternatives` relates each one. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: the maintainer's two decisions of 2026-11-24, "Tirar de vez" and "Os 8 arquivos (Recomendado)", quoted with their provenance in the authorization record, which grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. The governed set was measured again with `speccheck.GovernedPath` over every file the Build Order names (`references/2026-11-24-techspec-measurements.md`, section 1): the twelve granted paths are the only governed ones the design changes, and each change stays inside the bound the record gives its path. The design changes no `go.mod`, lint or formatter file, no Baseline asset, no other workflow and no skill reference beyond the four granted. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0420-roundfix-stops-verifying-on-linux/_authorization.md`; bounded files: `.roundfixrc.yml`, `Makefile`, `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/setup.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/runtime.md`, `skills/roundfix/SKILL.md`, `docs/references/coverage-record.json`, `docs/agents/specific-repository.md`.

## System Architecture

No package, file of production code or seam is added. Every change is inside a module that exists.

### The repository gate

`(*commandDeliveryWorkflow).Gate` in `internal/cli/deliver_workflow.go` trims the Shared Build Caches, takes one unit of the Machine Verification Capacity, runs the repository Verification with its Verification Core Budget, and then asks the lane. The last part leaves: the function returns a passed result as soon as the Verification passes. `internal/cli/deliver_linux_lane.go` leaves whole. Nothing outside the lane uses a declaration of that file except the doctor check, which leaves with it.

`internal/delivery/engine.go` loses the lane's field of the gate result, the warning it wrote on the item, the paragraph it added to the Pull Request body and the two constants that named them. `internal/delivery/park_class.go` keeps one case for a blocker that older versions stored.

### The command and the help

`internal/cli/deliver_lane.go` leaves whole, with its case in the dispatch of `runDeliverCommand`. The word `lane` then reaches the branch that already refuses any unknown `deliver` command. The usage text of `internal/cli/deliver.go` loses three lines and the root help of `internal/cli/cli.go` loses one word.

### The two retired keys

`internal/config/config.go` loses the two fields, their validation and their merge, and `deprecatedConfigKeys` gains two entries. The loader removes a registered key from the parsed document before it validates or decodes anything, whatever the kind of its value, and prints one line per key and load. Project Config, User Config and a configuration read at a commit all pass through that one function; the read at a commit discards the warnings, as it does today for every deprecated key.

### Doctor

`internal/cli/doctor.go` loses the check, its dependency and its four codes, and `internal/cli/health.go` loses the check's name. The checks before and after it keep their order.

### The two workflows and the compile contract

`.github/workflows/ci-verify.yml` changes one line, the runner of its only job. `.github/workflows/release.yml` gains a job `verify` on the same runner, which takes the two verification steps from the job `release`; that job gains `needs: verify` and keeps every other step on the runner it has. `internal/verifyselect` gains the contracts that read the four workflow files of this repository, and its compile contract gains Linux.

### The decision record, the glossary, the guide and the skill

One step writes ADR-0452, deletes the five records it supersedes, adds the amendment notices to six others, removes the four glossary entries of the lane and moves the citations of two entries. Another removes the lane from three pages of the user guide and four references of the Roundfix skill, lists the two keys among the deprecated ones, and corrects two sentences of the release runbook.

## Implementation Design

### Interfaces

The gate result, in `internal/delivery/engine.go`:

```go
type GateResult struct {
	Passed  bool
	Reason  string
	Blocker string
}
```

The blocker older versions stored, in `internal/delivery/park_class.go`:

```go
// retiredLaneContainerBlocker is a blocker that only versions before the
// Linux Lane was retired recorded. Nothing produces it.
const retiredLaneContainerBlocker = "linux-lane-container"
```

The two entries appended to `deprecatedConfigKeys`, in this order, in `internal/config/config.go`:

```go
{
	path:    []string{"delivery", "linux_lane"},
	name:    "delivery.linux_lane",
	removal: "the Linux Lane was retired, so remove it",
},
{
	path:    []string{"delivery", "linux_lane_image"},
	name:    "delivery.linux_lane_image",
	removal: "the Linux Lane was retired, so remove it",
},
```

The compile contract, in `internal/verifyselect/platform_test_compile_repocontract_test.go`:

```go
// compiledPlatforms are the release platforms no gate runs a test on.
var compiledPlatforms = []string{"linux", "windows"}

// compileEveryTest runs go vet over every package of dir for goos.
func compileEveryTest(ctx context.Context, dir, goos string) ([]byte, error)

func TestEveryTestCompilesForLinuxAndWindows(t *testing.T)
```

The runner the workflow contracts hold, in `internal/verifyselect/verifying_runner_test.go`:

```go
// verifyingRunner is the label of every job that runs a test, a vet or a contract.
const verifyingRunner = "macos-26"
```

The jobs of `.github/workflows/release.yml`:

```yaml
jobs:
  verify:
    runs-on: macos-26
    permissions:
      contents: read
    steps:
      # checkout with full history, Go setup, then:
      - name: Verify gate
        run: make verify
      - name: Verify every contract
        run: make verify-contracts
  release:
    needs: verify
    runs-on: ubuntu-latest
    # permissions, env and every other step as today
```

### Data Models

No schema change. A lane blocker is text in the item's `blocker` and a lane skip is text in the item's `warning`; both columns keep their meaning, and no stored row is rewritten. `docs/references/coverage-record.json`, `docs/references/behavior-surfaces.json` and `docs/references/skill-coverage.json` keep their schemas.

### API Contracts

1. API Contract: the repository gate. When the repository Verification passes, `Gate` returns a passed result. It no longer reads the Project Config at the default branch head, resolves no `docker`, writes no `linux-lane.log` and prints no line that names the lane. The engine's gate step then sets the stage `publishing` and leaves the item's warning as it found it, and the publication adds no paragraph to the Pull Request body. The unit of the Machine Verification Capacity is taken before the Verification and released after it under the step name `repository gate`, and the Verification Core Budget is computed once and passed to the Verification, as today. A failed Verification parks the item as today.
2. API Contract: `roundfix deliver lane`. The word is no command: with any arguments that hold no help flag, the command prints the two lines of Surface Transcript 1 on standard error and exits `2`. With a help flag it prints the usage of `deliver` and exits `0`, as any other word does today. The usage is the one of Surface Transcript 2, and the root help names the `deliver` commands in the line `  roundfix deliver <plan|start|status|resume|retry|stop> [<slug> ...]`.
3. API Contract: the two retired keys. A configuration that names `delivery.linux_lane` or `delivery.linux_lane_image` loads, in Project Config and in User Config, whatever the value is: a boolean, a text, a mapping, a list or nothing. Each named key is removed before validation, and the load prints once for it `config: <key> is deprecated and ignored; the Linux Lane was retired, so remove it` on standard error. A file that names both prints `delivery.linux_lane` first, and a key named in both files prints one line. A `delivery` mapping that held only these keys loads as an empty one. A configuration read at a commit removes the keys the same way and prints nothing. Any other unknown key under `delivery` fails the strict decode with today's message.
4. API Contract: doctor. `roundfix doctor` runs no check named `linux-lane` and prints none of `DR-LANE-IMAGE-UNKNOWN`, `DR-LANE-DOCKER-MISSING`, `DR-LANE-DOCKER-UNREACHABLE` and `DR-LANE-IMAGE-MISSING`, whatever the configuration says. The `verification-cores` line is followed by the `codex` line. Every other check keeps its name, its place and its text.
5. API Contract: what the lane recorded.
   - A parked item whose blocker begins `linux-lane-container:` is classed `environment`, as today, and its next step is `run roundfix deliver retry <slug>; the Linux Lane was retired, so the gate no longer runs it`. Status prints it in the `Park:` line, and the Pending Question of the follow-up carries the same text. The retry counts against the item's retry limit as today.
   - A parked item whose blocker begins `gate-failed: Linux Lane change failure:` keeps the `finding` class and the next step of every `gate-failed` item.
   - An item's warning that holds a `linux-lane-skipped:` segment is printed by status as stored. No code writes, reads or removes that segment. When the engine rewrites the item's low-disk warning, every other segment keeps its text and its order and the low-disk segment comes last.
   - The lane logs on disk are neither read nor deleted.
6. API Contract: the workflow `CI | Verify`. The job `verify` names `runs-on: macos-26`. No other byte of the file changes: the workflow's name, its triggers, its permissions and its concurrency, the job's name `Verification gate` and its timeout, and each step with its name, identifier, condition, environment and command, in today's order.
7. API Contract: the workflow `release`.
   - The job `verify` runs on `macos-26` with `contents: read`, for a pushed tag and for a manual dispatch. Its steps are the checkout with the full history and no persisted credentials, the Go setup with the options the release uses today, `Verify gate` and `Verify every contract`, the last two with today's names, commands and comments and with no condition. The checkout's comment that explains the full history moves here with the steps that need it.
   - The job `release` gains `needs: verify` and loses those two steps. Its runner, its permissions, its environment and every other step keep their bytes and their order: the checkout, the Go setup, the Node setup, `Guard Trusted Publishing runtime`, `Validate tag`, `Publication preflight`, `Cross-compile and stage`, `Compare published packages`, `Publish to npm` and `GitHub Release`. The comment of its checkout says that no step of this job reads the history any more and that the option stays as it was.
   - The triggers, the input and the workflow's name do not change, and neither job carries a condition that lets `release` run when `verify` did not succeed.
8. API Contract: the compile contract. `TestEveryTestCompilesForLinuxAndWindows` replaces `TestEveryTestCompilesForWindows`, in a file that keeps the `repocontract` tag and the Contract Relevance of the file it replaces. For each of `linux` and `windows` it runs `go vet -tags "repocontract docscontract" ./...` at the repository root with `GOOS` set and `GOARCH` inherited, and fails naming the platform and printing the compiler's output. A subtest runs the same function over a temporary module with no dependency, whose one test file builds only for Linux and does not compile: the Linux run fails with an output that names that file, and the Windows run passes. The Makefile's `REPO_CONTRACT_TESTS` names the new test in the place of the old one. The README says, after the sentence that names the platform packages, that the test suite runs on macOS only and that the Linux and Windows packages are compiled and not tested.
9. API Contract: the three statements about where the gate runs.
   - The Makefile's comment above `verify-changed` and the sentence of `docs/agents/specific-repository.md` no longer call continuous integration the Linux gate. Each says that the selective gate and continuous integration run on macOS and that a failure only Linux or Windows produces is a known limit, because those platforms are only compiled.
   - `SUITE_BUDGET_SECONDS` keeps the value `420`. Its comment keeps the measurement it states, says that the runs measured were on the Linux runner the gate used until this change, that the gate now runs on the standard macOS runner, where the suite has not been measured, and that the first push to main there is the measurement the number waits for.
10. API Contract: this repository's Project Config. `.roundfixrc.yml` loses the line `linux_lane: false` and the six comment lines above it. The `delivery` mapping then begins with the comment of `require_defect_class`.
11. API Contract: the decision record and what it changes.
    - `docs/adr/0452-roundfix-verifies-on-macos-and-the-delivery-queue-has-no-linux-lane.md` is accepted, carries a true Knowledge Mark and has the title `Roundfix verifies on macOS, and the Delivery Queue has no Linux Lane`. After the title it names what it supersedes and what it amends, each with the date of its `created_at`, in the form the records of this repository use: `Supersedes ADR-0288, ADR-0294, ADR-0297, ADR-0303 and ADR-0316 (<date>); all are retired, and Git keeps their text.`
    - It states as rules API Contracts 1 to 8 and the first item of API Contract 9, and that the suite-time budget keeps its number until the macOS runner measures the suite. It restates the rule of the tests' tools with its words: a test that needs an External Test Tool skips naming it, a tool outside the Base Test Tools is resolved through `testfixture.RequireTool`, and the rule in `internal/testfixture` refuses any other way.
    - It records the options given up: a light Linux job that only builds and vets, no test in continuous integration, the whole release on the macOS runner, three release jobs that validate the tag first, the moving label `macos-latest`, a larger runner, and a class of `unclassified` for the stored blocker.
    - The five superseded files are deleted in the same change.
    - Each of ADR-0252, ADR-0285, ADR-0295, ADR-0384, ADR-0385 and ADR-0423 gains `amended_by: [ADR-0452]` in its header, the date of the change in `updated_at`, and one paragraph that begins `**Amended by [ADR-0452](<file>) (<date>).**` after the text `## ADR Alternatives` quotes from it. No other byte of those six files changes, so every sentence this Spec quotes from them stays where it is.
    - `CONTEXT.md` loses the entries Linux Lane, Lane Environment Skip, Lane Container Failure and Lane Change Failure. Base Test Tools and External Test Tool keep their definitions and cite ADR-0452 in the place of the deleted record.
    - `.secondbrain-export` loses the lines of ADR-0288, ADR-0303 and ADR-0316 and gains the line of ADR-0452 in its place by number, after the last line of a decision record.
12. API Contract: the guide and the skill.
    - `docs/user-guide/commands/deliver.md` loses the sections `Linux Lane` and `Run the Linux Lane on demand`, the command's line of the usage, `linux-lane-container` in the Park Class table, and the lane in each sentence about the gate's capacity and core budget. `.agents/skills/roundfix/references/deliver.md` loses the same.
    - `docs/user-guide/commands/doctor.md` and `.agents/skills/roundfix/references/setup.md` lose the check and its four codes. The `implement` and `runtime` references lose their one mention each.
    - `docs/user-guide/configuration.md` loses the two keys from its reference of keys and lists them among the deprecated keys, each with its warning line. The `setup` reference lists them the same way.
    - `docs/user-guide/release-runbook.md` says, in the rehearsal of the Publication Preflight and in the order of the release workflow, that the workflow verifies on a macOS runner first and validates the version after it.
    - The mirrors under `skills/roundfix/` are written by `make skills-sync` and the version fields by the record command.
13. API Contract: the records.
    - `docs/references/coverage-record.json` is written by its own command in the step that deletes the tests. It loses the 36 names of `references/2026-11-24-techspec-measurements.md`, section 6, and no other, and gains the names of the tests that step adds.
    - `docs/references/behavior-surfaces.json` is written by its own command in each step that changes a help text, a configuration key or a page of the user guide.
    - `docs/references/skill-coverage.json` loses the entry `command deliver lane` in the step that removes the command. In the step that changes the release runbook, the entries `guide docs/user-guide/release-runbook.md` and `command release plan` gain a Coverage Review dated on the day of the change.
    - The number of `docscontract` files that `TestRepositoryContractTestsDeclareTheirRelevance` pins rises by one in each step that adds such a file.

### Surface Transcripts

1. Surface Transcript: the removed command, in any directory.

   ```transcript
   $ roundfix deliver lane
   stdout:
   stderr:
   roundfix: unknown deliver command "lane"
   Run 'roundfix deliver --help' for usage.
   exit: 2
   ```

2. Surface Transcript: the usage of `deliver`, in any directory.

   ```transcript
   $ roundfix deliver --help
   stdout:
   Usage:
     roundfix deliver plan [--json] [<slug>...]
     roundfix deliver start <slug>...
     roundfix deliver status
     roundfix deliver resume
     roundfix deliver retry <slug>
     roundfix deliver stop

   Records and advances an ordered queue of Specs. The detached owner survives the
   calling terminal. A blocked item is parked and the owner continues with the
   next item.

   Commands:
     plan    Report which Specs are approved to run
     start   Validate and record a new queue, then start its detached owner
     status  Print every queued Spec's stage, blocker and item worktree
     resume  Start a detached owner for the persisted queue
     retry   Return one parked Spec to the queue and reach its owner
     stop    Prove and terminate the persisted queue owner

   Flags:
     --max-duration <duration>  Set a positive queue duration
     --max-retries <n>          Set the per-item retry limit to at least 1
     --max-tokens <n>           Set the queue token ceiling to at least 1
   stderr:
   exit: 0
   ```

3. Surface Transcript: a configuration that names both keys. The command runs in a new Git repository whose `.roundfixrc.yml` holds a `delivery` mapping with `linux_lane: true` and `linux_lane_image: golang:1.26`, with a home directory that holds no Roundfix file.

   ```transcript
   $ roundfix runs list
   stdout:
   No Runs found.
   stderr:
   config: delivery.linux_lane is deprecated and ignored; the Linux Lane was retired, so remove it
   config: delivery.linux_lane_image is deprecated and ignored; the Linux Lane was retired, so remove it
   exit: 0
   ```

The root help and the lines of doctor and of status are not transcripts. A transcript reads text in angle brackets as a wildcard, so the root help's line would also match the line that still names `lane`, and the outputs of doctor and of status depend on the machine or on a prepared Run Database. Named tests assert those bytes (`## Testing Approach`).

### Invariants

```text
1. The Delivery Queue, its commands and doctor resolve and start no docker, and no test starts a container.
2. The Run Database's schema and every stored item keep their bytes.
3. Every line doctor prints for a check other than linux-lane keeps its bytes and its place.
4. Every line status prints keeps its bytes, except the next step of a stored linux-lane-container blocker.
5. A key that is not registered as deprecated still fails the strict decode.
6. The job Verification gate keeps its name, its timeout and its nine steps in order.
7. Every job that runs make verify, make verify-changed, make verify-docs, make verify-contracts, go test or go vet names the runner macos-26.
8. The job release runs no step before the job verify succeeded.
9. The steps the job release keeps have the bytes they have today.
10. The coverage record loses 36 names, each of the lane, and no other.
11. A test finds ADR-0452 by its title, never by its number.
12. An amended decision record keeps every byte of its body and gains one notice; its header gains amended_by and the date in updated_at.
13. No text inside a QA settlement section of a skill changes.
```

## Coverage Map

- Goal 1 → the repository gate, the command and doctor (API Contracts 1, 2 and 4, invariant 1)
- Goal 2 → the two retired keys (API Contract 3, Surface Transcript 3)
- Goal 3 → the two workflows and the compile contract (API Contracts 6 and 7, invariants 6 to 8)
- Goal 4 → the two workflows and the compile contract (API Contract 8), with the README's sentence
- Goal 5 → the decision record, the glossary, the guide and the skill (API Contracts 11 and 12)
- Story 1 → `Gate`, `GateResult` (API Contract 1)
- Story 2 → `deprecatedConfigKeys` (API Contract 3, Surface Transcript 3)
- Story 3 → `jobs.verify` of both workflows (API Contracts 6 and 7)
- Story 4 → `ClassifyPark`, `retiredLaneContainerBlocker` (API Contract 5)
- Story 5 → the README's sentence (API Contract 8)
- Story 6 → the usage texts, doctor, the guide and the skill (API Contracts 2, 4 and 12)
- Core Feature 1 → `Gate`, `GateResult`, the engine's gate and publication steps (API Contract 1)
- Core Feature 2 → `runDeliverCommand`, the usage texts (API Contract 2, Surface Transcripts 1 and 2)
- Core Feature 3 → `deprecatedConfigKeys` (API Contract 3, invariant 5)
- Core Feature 4 → doctor and the health check names (API Contract 4, invariant 3)
- Core Feature 5 → `ClassifyPark` and the engine's warning (API Contract 5, invariants 2 and 4)
- Core Feature 6 → `.github/workflows/ci-verify.yml` (API Contract 6, invariant 6)
- Core Feature 7 → `.github/workflows/release.yml` (API Contract 7, invariants 8 and 9)
- Core Feature 8 → `TestEveryTestCompilesForLinuxAndWindows`, the README (API Contract 8)
- Core Feature 9 → the Makefile's budget comment (API Contract 9)
- Core Feature 10 → `.roundfixrc.yml` (API Contract 10)
- Core Feature 11 → ADR-0452, the six notices, the glossary, the export list, the Makefile's comment and the guide's sentence (API Contracts 9 and 11, invariant 12)
- Core Feature 12 → the guide and the skill (API Contract 12), the records (API Contract 13)
- Core Feature 13 → `references/2026-11-24-techspec-measurements.md`, sections 4, 6 and 7, the tests of `## Testing Approach`, the coverage record (API Contract 13, invariant 10)
- Success Metric 1 → `TestTheRepositoryGateEndsAtTheRepositoryVerification`, `TestAPassedGateLeavesNoLaneWarningAndNoPullRequestParagraph`
- Success Metric 2 → `TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand`, `TestDeliverHelpReproducesSurfaceTranscriptWithoutTheLane`, `TestTopLevelUsageNamesDeliverRetry`
- Success Metric 3 → `TestTheRetiredLaneKeysLoadWithOneWarningEach`, `TestRunsListReproducesSurfaceTranscriptRetiredLaneKeys`
- Success Metric 4 → `TestDoctorRunsNoLinuxLaneCheck`, with the doctor tests of the measurements file, section 7, which pin every other line
- Success Metric 5 → `TestStatusPrintsWhatTheLaneRecorded`, `TestADeliveryRetryReturnsAnItemTheLaneParked`, `TestAStoredLaneBlockerKeepsItsParkClass`
- Success Metric 6 → `TestEveryVerifyingJobRunsOnTheStandardMacOSRunner`, `TestTheReleaseVerifiesBeforeItBuildsOrPublishes`, `TestTheVerificationGateKeepsItsNamesAndItsSteps`, `TestMainAndReleaseRunEveryContract`
- Success Metric 7 → `TestEveryTestCompilesForLinuxAndWindows`
- Success Metric 8 → `TestTheLinuxVerificationDecisionSupersedesTheLaneRecords`, `TestNoLiveDocumentDescribesTheLinuxLane`, `TestTheGuideAndTheSkillListTheRetiredLaneKeys`, `TestCoverageEquivalence`, with the existing documentation, skills and glossary contracts
- Success Metric 9 → `verifyingRunner` in `TestEveryVerifyingJobRunsOnTheStandardMacOSRunner`, read by the QA gate against the labels the PRD's `## Acceptance evidence` quotes

Deviations from the PRD:

- Core Feature 9 asks for a number and a comment that state a measurement taken on the macOS runner. A Task reaches no runner, and the delivering pull request runs the incremental Verification, which is not the suite the budget times. The number stays, the comment says what it was measured on and that the macOS measurement is pending, and the operator writes the measured number after the first push to main, as the second criterion of the PRD's `## Unreachable Acceptance` says.
- Core Feature 11 asks that every citation of a deleted record in a live document become a citation of ADR-0452 or a Git Locator. The glossary's links do. ADR-0423 names ADR-0288 in plain text inside a sentence about its own context; that sentence keeps its bytes, and the amendment notice after it names ADR-0452 as the successor (invariant 12).
- Core Feature 12 names three pages of the user guide. The design also changes two sentences of `docs/user-guide/release-runbook.md`, because the release's verification now precedes its tag validation and that page states their order.
- Core Feature 13 asks for named tests that state on the unchanged tree what the lane does. The unchanged tree already holds them: the 36 tests the coverage record lists. This stage's probes add what they do not cover (`references/2026-11-24-techspec-measurements.md`, section 4). No new test pins the lane, because the step that removes the lane would have to delete it in the same change; each test of `## Testing Approach` fails on the unchanged tree instead.

## Absence Evidence

1. ADR-0452 — command: `ls docs/adr | rg '^045'`; result: two files, numbered 0450 and 0451.
2. The two keys in the registry of deprecated keys — command: `rg -n 'linux_lane' internal/config/config.go`; result: six lines, the two fields of the overlay at 470 and 471 and the validation at 1873 to 1877; none between 768 and 794, where `deprecatedConfigKeys` is declared.
3. A macOS runner in a workflow or in the Makefile — command: `rg -n -i 'macos' .github/workflows Makefile`; result: no match.
4. A job that another job needs — command: `rg -n 'needs:' .github/workflows`; result: no match.
5. Linux in the compile contract — command: `rg -n 'GOOS' internal/verifyselect`; result: one match, `GOOS=windows` in `windows_test_compile_repocontract_test.go`.
6. A contract that reads a workflow's runner — command: `rg -n 'runs-on' internal --glob '*_test.go'`; result: one match, the Baseline's starter workflow in `internal/baseline/starter_files_test.go`, which is not a workflow of this repository.
7. A next step that says the lane was retired — command: `rg -n 'was retired' internal/delivery internal/cli --glob '!*_test.go'`; result: no match.
8. A statement in the README about what is tested — command: `rg -n -i 'tested|test suite|compiled' README.md`; result: one match, about the minimum tested version of a tool, at line 39.
9. `internal/verifyselect/verifying_runner_test.go`, `internal/verifyselect/platform_test_compile_repocontract_test.go` and the two documentation contracts — command: `ls internal/verifyselect internal/docscontract | rg 'runner|platform|lane|linux'`; result: no match.
10. A Coverage Review of the release runbook for this change — command: `rg -n '"review"' docs/references/skill-coverage.json | rg 'release'`; result: two matches, both dated 2026-11-06 for Spec 0364.

## ADR Alternatives

- ADR-0452: the record a Task of this Spec writes in step 4 of the Build Order, absent until then; it rejects nothing yet. API Contract 11 fixes what it states and the options it records as given up.
- ADR-0288: superseded by ADR-0452 and deleted, by the maintainer's decision of 2026-11-24 that the PRD's `## Rule challenges` names. The design takes none of the options it rejected, which were other ways of running the lane; it removes the rule itself (ADR-0288: "The Delivery Queue's repository gate runs a Linux Lane when Project Config opts in.").
- ADR-0294: superseded by ADR-0452 and deleted, by the same decision. The lane's image, its command and its two kinds of failure leave with it (ADR-0294: "The lane runs on demand."). One rule of it is not about the lane and moves to ADR-0452 with its words (ADR-0294: "A test that needs an External Test Tool skips naming it."). No option it rejected is taken.
- ADR-0297: superseded by ADR-0452 and deleted, by the same decision (ADR-0297: "The object store is mounted read-only."). Nothing is mounted any more, so none of the mounts it rejected is taken.
- ADR-0303: superseded by ADR-0452 and deleted, by the same decision (ADR-0303: "Doctor checks the lane it will need."). What a skip already recorded stays readable, and the field and the column it rejected for that record are not added.
- ADR-0316: superseded by ADR-0452 and deleted, by the same decision (ADR-0316: "Its Park Class follows the cause."). A stored container blocker keeps the class it gave, and no option it rejected is taken.
- ADR-0252: amended in part, by the same decision. The sentence that names the platform of continuous integration (ADR-0252: "A failure that only another platform produces is a known limit") gains a notice: the gate and continuous integration now run on macOS. Its selection rule stands, and neither option it rejected is taken: the contracts are selected by their Contract Relevance, with no central list.
- ADR-0285: amended in part, by the same decision: its contract compiles the tests for Linux as well (ADR-0285: "A repository contract compiles every test for Windows."). The options it rejected stay rejected: the gap is not accepted as a platform-only limit, no static scan replaces the compiler, and the contract keeps its default Contract Relevance. The Makefile edit it had no authority for is granted to this Spec for the contract's name.
- ADR-0295: amended in part, by the same decision: only the example that names the lane changes (ADR-0295: "such as the Delivery Queue and the Linux Lane, keep their reads"). The rule and the options it rejected are untouched.
- ADR-0384: amended in part, by the same decision: the gate's unit covers the repository Verification alone (ADR-0384: "takes one unit for its repository Verification and its Linux Lane together"). The capacity, its default and the options it rejected stand; no table and no second unit is added.
- ADR-0385: amended in part, by the same decision: the lane's container leaves the receivers of the budget (ADR-0385: "Verification and Linux Lane container receive it"). The number, its two names and the options it rejected stand.
- ADR-0423: amended in part, by the same decision: its context names the lane's key as the delivery setting the Queue Owner already read at the default branch head (ADR-0423: "so that an Item Branch cannot change it"), and the notice says the key is retired. Its rule stands (ADR-0423: "never from the item's tree or from the owner's checkout"), and no option it rejected is taken.
- ADR-0253: obeys it (ADR-0253: "CI runs the Full Contract Run on every push to main and in the release workflow before any build or publication."). The run stays in both places, and in the release it is in the job the publishing job needs. The hand-kept list and the widened documentation step it rejected are not taken.
- ADR-0387: obeys it (ADR-0387: "changing the number takes a new measurement"). The number stays until the macOS runner measures the suite; the design does not raise it alone, which the record rejected. It also rejected a separate job for the Full Contract Run on a push: the push keeps its one job, and the release's new job holds the repository Verification and the Full Contract Run together.
- ADR-0236: obeys it (ADR-0236: "merges the default branch tip it fetched at job start"). The step that merges the base and prints the `tested-base` notice keeps its bytes; only the job's runner changes. The options it rejected are not taken.
- ADR-0250: obeys it (ADR-0250: "The coverage record lists tests for every release platform."). Linux stays a release platform, the tests that build only for Linux stay in the record, and the record is written by its own command. Comparing only the host's subset, which it rejected, is not taken.
- ADR-0031: obeys it (ADR-0031: "cross-compiles the Go binary for darwin/linux/windows on both architectures"). The platforms, the packages and the single cross-compiling job do not change, and the other channel it deferred is not added.
- ADR-0365: obeys it (ADR-0365: "a schema change makes older item binaries refuse the shared Run Database"). The stored blocker and the stored warning are text the Run Database already holds, and the column it rejected is not added.
- ADR-0304: obeys it (ADR-0304: "No stage, column or schema is added."). The stage and the column it rejected are not added.
- ADR-0398: obeys it (ADR-0398: "A parked item always has a way back through a Delivery Retry"). An item the lane parked returns through a Delivery Retry, which pushes nothing, as it required.
- ADR-0450: obeys it (ADR-0450: "Success requires a proven live owner unless every item is already merged or parked."). The retry of an item the lane parked ends as every retry does, and the test of it runs with a proven live owner. The two options it gave up are not taken.
- ADR-0256: obeys it (ADR-0256: "is a Lagging Surface unless a covering skill file changed in the same range"). Every Behavior Surface the removal moves has a covering reference that changes in the same range, and the two surfaces of the release runbook carry a Coverage Review. The record is written by its command, never judged by a model, which it rejected.
- ADR-0349: uses it as it is (ADR-0349: "the plan accepts without asking while the fresh plan stays at or below it"). The authorization record carries a `patch` ceiling. It rejected answering a manual classification from the ceiling, and the design leaves the classification of a `refactor` to the operator. No Task performs the release.
- ADR-0056: not applicable — no capacity setting is read or changed; the gate's unit is the one ADR-0384 gives it.
- ADR-0259: not applicable — its rule and its table of Sequential Tests do not change; the tests this Spec adds follow it as every test does.
- ADR-0279: not applicable — the owner lines of status keep their bytes.
- ADR-0311: not applicable — the Check Wait, its key, its default and its bounds are not changed.
- ADR-0369: not applicable — the Head Check Set is derived as before; the workflow and the job keep their names.
- ADR-0395: not applicable — when an item stops and what the owner logs are not changed.
- ADR-0406: not applicable — the authoring probe's checkout is used as it is.
- ADR-0034: not applicable — the Worktree Bootstrap is not changed.
- ADR-0148: not applicable — the Verification probe is used as it is.
- ADR-0271: not applicable — no real-time budget is added or raised.
- ADR-0373: not applicable — how a reader proves a Queue Owner gone is not changed.
- ADR-0412: not applicable — the release plan's reading of a commit is not changed.
- ADR-0435: not applicable — a Pull Request behind its base is handled as before.

## Integration Points

- GitHub Actions, through two workflow files. The runner `macos-26` is a standard GitHub-hosted runner, free for a public repository (`_prd.md`, `## Acceptance evidence`). Its image lists four of the eight External Test Tools, `bash`, `jq`, `node` and `rustfmt`; the other four, `cmp`, `grep`, `gzip` and `sh`, are part of macOS, as on the maintainer's machine. The workflow adds no step for them: a test whose tool is absent skips naming it, as the rule requires. The image lists no Docker, and nothing asks for it.
- The Go cache of the job `Verification gate`. The setup action keys its cache by the runner's operating system and the step that aligns `GOCACHE` reads `go env GOCACHE`, so neither holds a Linux path. The first run on the new runner starts with an empty cache.
- npm Trusted Publishing, unchanged: the publishing job keeps its runner, its `id-token: write` and the workflow's file name, which is what each package's trusted publisher is bound to.
- Docker, removed: after this Spec the Delivery Queue and doctor neither resolve nor start it.
- The Run Database, read as before. No statement changes.

No Agent, ACP adapter, ai-memory call or network call is added, and no test or Verification reaches GitHub.

## Testing Approach

Every test attaches to a seam that exists. No test starts a container, reaches the network or reads the real home directory.

- The gate, in `internal/cli`: `TestTheRepositoryGateEndsAtTheRepositoryVerification` drives `Gate` through the fixture the surviving gate tests use, over a temporary repository whose fetched default branch says `linux_lane: true` and names an image, with a fake runner and a recording stand-in for `docker` on the PATH the command environment carries. It asserts a passed result, that the stand-in was never started, that no `linux-lane.log` exists and that the console holds no line of the lane. `TestAPassedGateLeavesNoLaneWarningAndNoPullRequestParagraph`, in `internal/delivery`, drives the engine with a gate that passes and asserts that the item reaches `publishing` with an empty warning and that the Pull Request body is the publication's own.
- The command and the help, in `internal/cli`: `TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand` and `TestDeliverHelpReproducesSurfaceTranscriptWithoutTheLane` assert Surface Transcripts 1 and 2 whole. `TestTopLevelUsageNamesDeliverRetry` keeps its name and pins the root help's new line.
- The keys: `TestTheRetiredLaneKeysLoadWithOneWarningEach`, in `internal/config`, loads Project Config and User Config that name the keys with a boolean, a text and a mapping, alone and together, and asserts the loaded configuration, the exact warnings and their count. It also asserts that an unknown key beside them still fails with today's message and that a configuration read at a commit prints nothing. `TestRunsListReproducesSurfaceTranscriptRetiredLaneKeys`, in `internal/cli`, asserts Surface Transcript 3 whole.
- Doctor, in `internal/cli`: `TestDoctorRunsNoLinuxLaneCheck` runs doctor with `linux_lane: true` in the checkout's configuration and the recording stand-in, and asserts that no line begins `linux-lane:`, that no `DR-LANE-` code is printed and that the stand-in was never started. The doctor tests that pin lines keep their names and their other lines.
- What the lane recorded: `TestAStoredLaneBlockerKeepsItsParkClass`, in `internal/delivery`, asserts API Contract 5 for the blocker and that a stored skip segment survives a rewrite of the low-disk warning. `TestStatusPrintsWhatTheLaneRecorded` and `TestADeliveryRetryReturnsAnItemTheLaneParked`, in `internal/cli`, run over a temporary Run Database with one parked item and one merged item; the retry uses the owner fake of `internal/cli/deliver_retry_owner_proof_test.go`.
- The workflows, in `internal/verifyselect`, untagged like the tests that already read them: `TestEveryVerifyingJobRunsOnTheStandardMacOSRunner` reads every file of this repository's `.github/workflows` and asserts invariant 7 in both directions. `TestTheReleaseVerifiesBeforeItBuildsOrPublishes` asserts API Contract 7 and invariant 8. `TestTheVerificationGateKeepsItsNamesAndItsSteps` asserts invariant 6. `TestMainAndReleaseRunEveryContract` keeps its name and reads the release's two steps in the job `verify`. None of them reads a Baseline asset.
- The compile contract: `TestEveryTestCompilesForLinuxAndWindows`, with its subtest over a temporary module.
- The documents, in `internal/docscontract`: `TestTheGuideAndTheSkillListTheRetiredLaneKeys` reads the configuration page and the `setup` reference. `TestTheLinuxVerificationDecisionSupersedesTheLaneRecords` finds ADR-0452 by its title and asserts its status, its mark, the five absent files, the six headers and notices, the glossary and the export list. `TestNoLiveDocumentDescribesTheLinuxLane` reads the user guide, both copies of the Roundfix skill, the agent guides, the README, the Makefile, the workflows and the Project Config for the five strings of Success Metric 8, and allows `Linux Lane` only in a line that holds a retired key's warning.

Each of these tests fails on the unchanged tree. Two are partly controls that also hold there, the status of a stored skip and the retry of a parked item, so no Verification rests on them alone. A Verification never depends on a test that can skip on the host: none of these resolves an External Test Tool.

Existing tests that change are named in `references/2026-11-24-techspec-measurements.md`: section 6 lists the 36 that leave with their eight files, and section 7 the ones that keep their names and lose the lane from a fixture, a stub or a pinned line.

## Build Order

1. The product loses the lane: API Contracts 1 to 5 and 10, with the entries of API Contract 13 for the coverage record, the Behavior Surface Record and the map's entry `command deliver lane`. Files: the two deleted lane files, the eight deleted test files and the changed test files of the measurements file, sections 6 and 7, `internal/cli/deliver.go`, `internal/cli/cli.go`, `internal/cli/deliver_workflow.go`, `internal/cli/doctor.go`, `internal/cli/health.go`, `internal/delivery/engine.go`, `internal/delivery/park_class.go`, `internal/config/config.go`, `.roundfixrc.yml`, the three records and seven new test files. The coverage record and the Behavior Surface Record are written by their commands, after the last edit.
2. The verification runs on macOS: API Contracts 6 to 9 (depends on: nothing; it shares no file with step 1). Files: `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `Makefile`, `README.md`, `docs/agents/specific-repository.md`, `internal/verifyselect/full_contract_run_test.go`, the deleted `internal/verifyselect/windows_test_compile_repocontract_test.go` and two new test files.
3. The guide and the skill: API Contract 12 and the Coverage Review of API Contract 13 (depends on: 1 and 2, whose results it describes; it shares the Behavior Surface Record and the map with step 1). Edit the four guide pages and the four canonical references, then run `make skills-sync`, the skill record command and the Behavior Surface record command of the authorization record. Files: `docs/user-guide/commands/deliver.md`, `docs/user-guide/commands/doctor.md`, `docs/user-guide/configuration.md`, `docs/user-guide/release-runbook.md`, the four references and their four mirrors, both `SKILL.md` files, `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json`, `docs/references/skill-coverage.json`, `internal/verifyselect/contracts_repository_test.go` and one new test file.
4. The decision record and the glossary: API Contract 11 (depends on: 1, 2 and 3, because its sweep reads what they leave and it shares `internal/verifyselect/contracts_repository_test.go` with step 3). Files: the new file of ADR-0452, the five deleted records, the six amended records, `CONTEXT.md`, `.secondbrain-export`, `internal/verifyselect/contracts_repository_test.go` and one new test file. The header entries and the export lines are written by hand, as the guides allow (`docs/agents/secondbrain.md`: "the only line added by hand is the path of an ADR, in order"; `docs/agents/docs-layout.md`: "In the change that adds the later ADR, its header gains").
5. The QA gate (depends on: 1, 2, 3 and 4).

Files shared with other Specs: the Roundfix skill's version fields, the owned-skill version record, the coverage record, the Behavior Surface Record and the pinned number of `docscontract` files. Each step regenerates them on the tree it runs on and writes the next version and the next count, never a number. No Verification of this Spec depends on another Spec, and the Task Graph needs no `requires`.

The governed paths fall in three steps: two in step 1, four in step 2 and six in step 3 (`references/2026-11-24-techspec-measurements.md`, section 1). Step 4 changes none.

## Risks & Considerations

- The first run on the new runner is the delivering pull request's. The runner has 3 processors and 7 GB of memory, and this suite is bound by the processes it starts. A suite that cannot pass there parks the item before anything merges, and the question returns to the maintainer; a larger runner is billed and is not an option of this Spec.
- The budget's number may be wrong for the new runner. An overrun fails only the last step of a push to main, after every other step ran. The operator reads the first push and corrects the number with the measured time before the release; that correction changes the Makefile again and is the session's to place under the grant or under a new decision.
- The release's new shape is first exercised by the manual dispatch, which runs the job `verify` and the job `release` up to the Publication Preflight. The steps a tag adds keep their bytes and their runner.
- A wrong tag is now refused after the verification. The run publishes nothing either way.
- Four test files that build only for Linux no longer run anywhere. They stay compiled and stay in the coverage record, as the decision accepts.
- The Queue Owner that delivers this Spec is the installed version. It reads the lane's key at the default branch head, which says `false` today and names no key after the merge; both mean no lane. A queue of another repository whose owner is older is unaffected, because it never named the key.
- A repository that names a key after the upgrade gets one warning per key on every command that loads its configuration, until the key is removed.
- ADR-0452's number can be taken while the Pull Request waits. The operator then renames the file, its links in the six notices and the two glossary entries, and its export line; invariant 11 keeps the tests independent of the number.
- The amendment notices follow sentences this Spec quotes. A Task that rewrote one of those sentences would break a Claim Receipt of this plan, which is why invariant 12 forbids it.

## Glossary

None.

## Decisions

- The PRD's fifth open question is settled by measurement: every Behavior Surface the removal moves has a covering reference among the four the grant holds, so the removal needs no excuse (`references/2026-11-24-techspec-measurements.md`, section 3). The release runbook's two surfaces are covered only by the `release` reference, which the grant does not hold and which does not describe the workflow's order, so the PRD gains a `## Skills` section with those two entries and step 3 records their Coverage Review, as Spec 0364 did for the same page.
- Delegated Decision of 2026-11-24: How is the suite-time budget measured on the macOS runner before the first push to main is judged?
- Delegated Decision of 2026-11-24: Which label names the runner, the moving `macos-latest` or a versioned one?
- Delegated Decision of 2026-11-24: How does the release keep its verification before the publication: a job of its own that the publishing job waits for, or one job on the macOS runner?
- Delegated Decision of 2026-11-24: Which Park Class and next step does status show for a stored `linux-lane-container` blocker once no code produces it?
- The product step and the documentation step are two steps. A rehearsal settled it: with the lane removed from the help, the record written and the map's entry removed, every documentation contract passes while the guide and the skill still describe the lane (`references/2026-11-24-techspec-measurements.md`, section 5).
- The decision record, the glossary and the five deletions are one step. The same rehearsal settled it: deleting the five records breaks the glossary's links at once, and the Claim Receipts of this plan stay provable because the check reads a deleted record from Git.
- The compile contract is renamed and keeps one name in the Makefile's list. Its proof on a temporary module is a subtest, so the list gains no entry, which the authorization record's bound on the Makefile would not hold.
- The contract inherits `GOARCH`, as the one it replaces does. No file of the repository is constrained by architecture.
- ADR-0452 carries a true Knowledge Mark: it states that the product has no Linux Lane and that two keys are retired, which someone outside this repository's development can use (`docs/agents/secondbrain.md`: "that someone outside this repository's development can use").
- A stale skip warning is not removed when an item passes the gate again after the upgrade. A source settles it: the PRD's Core Feature 5 says that no stored record is rewritten.
- The publishing job's checkout keeps its options. No step left in that job reads the history, but a change to the job that publishes is first exercised by a real tag, so only its comment changes.
- The PRD's `## Glossary` declares the two entries this Spec changes; this stage names no other term.
- The Secondbrain pages and the published sources consulted on 2026-11-24, with what each changed, are recorded in `references/2026-11-24-techspec-measurements.md`, section 8.
