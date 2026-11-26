---
spec: 0421-the-test-suite-meets-a-measured-time-target-at-every-gate
prd: _prd.md
created: 2026-11-25
---

# The test suite meets a measured time target at every gate — Technical Spec

## Executive Summary

This repository's own command decides what each gate runs, from two facts it can already read: the Run Marker with the Run's branch, and an Archive Record in the change. A Task settlement runs format, vet, build, the selected Repository Contract Tests and the comparison of the coverage record. The Complete Suite Run happens at the Delivery Queue's repository gate. Continuous integration runs only for a pull request, as the gate of its squash merge: five parallel jobs on `macos-26`, gathered by one job that keeps the required check's name, run the packages the change reaches and every Repository Contract Test once. Nothing runs after a merge. A release runs no test: before it builds anything it proves, from GitHub, that its commit is the squash merge of a pull request whose check passed on that same tree. A new repository tool, `cmd/gate-time`, prints each gate's time against a target it reads from one reference file outside every Governed Path. The product, the Daemon and the Project Config do not change.

The primary trade-off: a package test that a Task breaks and its declared Verification does not name is found at the queue's gate, where the item parks and no Agent Session repairs it. The PRD's first risk accepts this on the sixteen diagnostics, and the reference document counts each case. A second one comes with the maintainer's answer on continuous integration: no run follows a merge, so two pull requests that each pass can meet on main in a tree neither verified. The release's pull request closes that gap, because it runs every package and every contract and is refused unless its branch held main's tip.

The one Reserved Question of this design, what a release verifies, was answered on 2026-11-25: the Praetor applied the maintainer's decision on continuous integration and granted the Release Provenance designed here, with the bound of `release.yml` it needs. The PRD's `## Decisions` quotes that answer as the Praetor's application, and the authorization record is approved with it.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted. Build Order step 4 writes ADR-0453, the next free number, and the Gate Time Rule takes the next free number when the closing change writes it. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the product gains no endpoint, credential or network call. `cmd/gate-time windows` runs on the mini outside every gate, with the GitHub CLI session the operator already has, and its tests read recorded answers. The verification workflow keeps its permissions, and the publication job keeps its permissions and its Trusted Publishing. The release's first job reads two endpoints of GitHub with the workflow's own token and gains two read-only permissions for it, `checks: read` and `pull-requests: read`, which the answer on the release grants. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the obligations are the ones the PRD's row lists, with their receipts, and this design adds two. The Run Marker is read, never set or cleared (ADR-0283: "Every process a Run starts for an Agent or a Verification carries the Run Marker."). The maintainer changed where continuous integration runs (message `01a213bf-b859-7f60-91a1-201ca006c98a`, quoted in the PRD's `## Decisions`): it runs only for a pull request, so the run after a merge that one decision fixed ends (ADR-0253: "CI runs the Full Contract Run on every push to main and in the release workflow before any build or publication"), and each contract runs once for each event. Outside continuous integration the contracts rule stands (`docs/agents/specific-repository.md`: "pass before any pull request opens"). By the answer on the release, the Praetor's application of the same decision, ADR-0453 withdraws the release clause of ADR-0253 and the sentences of ADR-0452 on the release's verification job, through their notices. ADR-0452, ADR-0387, ADR-0259, ADR-0257, ADR-0182, ADR-0342, ADR-0253, ADR-0252, ADR-0326, ADR-0271, ADR-0281 and ADR-0311 are accounted for in `## ADR Alternatives`, with ADR-0258, ADR-0273, ADR-0282, ADR-0285, ADR-0286, ADR-0295, ADR-0298, ADR-0319, ADR-0384, ADR-0385, ADR-0404, ADR-0423 and ADR-0435, which cite one of them and decide nothing this design changes. Step 4 writes ADR-0453 through the lifecycle (`docs/agents/docs-layout.md`: "A retirement deletes the file."). Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`, `docs/agents/specific-repository.md`, `docs/agents/agent-instructions.md`.
- Tooling authority: applicable — the design proposes protected tooling mutation. On 2026-11-25 the maintainer authorized the nine paths below with their bounds, the regenerations, the changes of the decision records and three additions, adjusted by his answer on continuous integration; the record quotes the message. The Release Provenance changes the steps, the runner and the permissions of the job `verify` of `.github/workflows/release.yml`; that bound was granted the same day by the Praetor's application of the maintainer's decision on continuous integration, which the record holds as that application. The record is `approved`, granted on 2026-11-25. The Task Graph phase measured again over every file the Tasks declare and found the same nine. On 2026-11-25 the author measured `speccheck.GovernedPath` again over the 122 paths this design creates, edits or deletes, at `3f621ac6`, in a `go test -overlay` probe that wrote nothing to the repository: fifteen are governed, the nine below and six the design leaves unchanged (`internal/cli/baseline_release_gate_test.go`, `internal/gittest/gittest.go`, `internal/baseline/plan_test.go`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md`, and `dist/npm/roundfix/package.json`, which only a release changes). No path beyond the nine is needed. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0421-the-test-suite-meets-a-measured-time-target-at-every-gate/_authorization.md`; bounded files: `Makefile`, `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `.roundfixrc.yml`, `docs/agents/setup-context.json`, `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`, `internal/cli/baseline_plan_test.go`, `internal/spec/coverage_test.go`.

## System Architecture

Three existing parts are extended and one is new.

- The selector, `cmd/verify-select` and `internal/verifyselect`, gains a package plan, the shards of the CLI package and a build-tag filter for contracts. Its two fixed sets, `core` and `baseline`, and the Baseline pattern of the CLI package are removed.
- The Makefile keeps `make verify-changed` as the one command the Project Config names (`.roundfixrc.yml`: "verification: make verify-changed"). Its recipe runs what the plan prints. It gains `coverage-record-check`, `verify-part` and `gate-times`. The suite-time budget's targets and variable, `verify-changed-core`, `verify-changed-baseline` and the `docs-test-skipped` switch leave.
- `ci-verify.yml` runs for a pull request only: one matrix job of five parts and one gathering job. `release.yml` runs no test, and its job `verify` proves the Release Provenance.
- The gate-time tool, `cmd/gate-time` and `internal/gatetime`, is new. The release builds only `cmd/roundfix`, so it is not part of the product. It stays apart from the selector, which every gate compiles, because it reads the Run Database and calls the GitHub CLI.

### Where the Complete Suite Run lives

The Daemon and the Delivery Queue run the configured command through one function, `ExecVerifier.Verify` in `internal/daemon/daemon.go`, and the command can tell the gates apart by two facts.

1. Inside a Run. For a Task settlement, the Integrated Repository Verification and the precondition of the QA gate, the Daemon passes the Run's id, and the function sets `ROUNDFIX_RUN_ID` (CONTEXT.md: "The `ROUNDFIX_RUN_ID` variable naming the Run in each of its Agent Session and Verification processes."). The command runs in the Run Worktree or a Task Worktree, on a branch named `roundfix/run-<run id>` or `roundfix/run-<run id>-<task id>`.
2. The candidate of a Spec. The queue's stages run in a fixed order: Run, review, archive, gate, publish (`archiveCandidate`, then `gateCandidate`, in `internal/delivery/engine.go`). The gate (`commandDeliveryWorkflow.Gate` in `internal/cli/deliver_workflow.go`) runs the command in the item worktree, on the branch `roundfix/deliver-<slug>-<suffix>`, with no Run id, at the commit that archived the Spec. The change from the base therefore adds one file directly under `docs/history/specs/`.

The plan reads them in this order:

| Mode | When | Package tests | Fresh results |
| --- | --- | --- | --- |
| `settlement` | the Run Marker is set and the checked-out branch is that Run's | none | not applicable |
| `complete` | otherwise, the change adds an Archive Record | every package | yes, `-count=1` |
| `reach` | otherwise | the packages that hold a changed file and the packages that import them | as the caller sets |

The Run's branch is the second fact for `settlement` because the function leaves an inherited marker in place when it is given no Run id: a marker alone would let a queue owner started from a marked shell skip the Complete Suite Run in silence. A marker without the Run's branch falls to the wider modes. A review Batch, which runs in the user's checkout, is such a case and keeps its package tests.

The QA gate is not the home. Its precondition runs the same command on the same branch as a settlement, the product derives the QA Task's Verification (`spec.DerivedQAVerification`), and a failure of that precondition refuses the gate without an Agent turn.

Consequences:

- The Complete Suite Run runs once for each pass of the queue's gate. A retry of a parked item and a Base Update run it again.
- A failure parks the item with `gate-failed`, as a failure of that gate does today.
- A Run delivered outside the Delivery Queue gets the Complete Suite Run from the same command: the Supervisor runs `make verify-changed` on the branch after the archive commit. The repository's own guide says so.
- A stale local `main` widens the change list. It cannot turn a settlement into a complete run, and at the queue's gate the candidate's own Archive Record is always in the list.

### What a settlement runs

`make verify-changed` in the `settlement` mode runs `fmt-check`, `vet`, `build`, the Repository Contract Tests that Contract Relevance selects, and `coverage-record-check`. The checks of the repository against its own rules and records that are code tests today become Repository Contract Tests declared `//verify:always`, in files under the `docscontract` tag in their own packages:

- the six rule tests of `internal/testfixture` and the two rules step 1 adds;
- `TestRepositoryContractTestsDeclareTheirRelevance`;
- the two tests of `internal/cli` that compare this repository with its Baseline record.

That class is the one the decision describes (ADR-0252: "It is for contracts whose inputs are spread across the repository and that cost about a second."). The tag is `docscontract` because the coverage record is collected with it, so no test name leaves the record, and the two CLI tests compile into the tagged test binary a settlement already builds.

The comparison of the coverage record stays an untagged code test, so its sanctioned command and the Project Config value that names it do not change. `coverage-record-check` runs `TestCoverageEquivalence` and `TestCoverageRecordCountsNoPackageUnderDocs` by name in every mode. `skills-sync-check` and `skills-check` run in every mode when the change touches `skills/`, `.agents/skills/` or the `Makefile`, as today.

### What a pull request selects

In the `reach` mode the selector lists the module once with `go list -json ./...` and takes the closure of importers over imports, test imports and external test imports. A changed file that is not Go source selects the nearest package directory above it. `go.mod`, `go.sum`, the `Makefile` and any path outside a package select every package, as today, with one exception: a path under `docs/history/`, `docs/adr/`, `docs/specs/` or `docs/user-guide/` selects no package, because no code test may read those trees (ADR-0326: "A test that reads the repository's live documentation trees is a Repository Contract Test."). `.agents/skills/` selects the Baseline packages and their importers, as today. Two of the three version-only paths of a release, `CHANGELOG.md` and `dist/npm/roundfix/package.json`, lie outside every package, so the pull request of a release selects every package; a test of the selector pins it.

### The jobs

Measured on the macOS runner in pull request `45529968311`, with a cold build cache: format, vet and build 65 s; the packages outside the CLI package 431 s and the Baseline packages 119 s; the CLI package 564 s, plus 69 s for its Baseline tests; every contract 425 s; `make verify-docs` 232 s, of which 214 s repeat the regeneration contracts. By the author's estimate, steps 1 and 2 remove about 390 of the 1,109 seconds the test binaries outside the CLI package ran.

| Part | What it runs | Estimate |
| --- | --- | --- |
| `packages` | `fmt-check`, `vet`, `build`, the skills checks, then the code tests of every selected package except `internal/cli` | 375 s |
| `cli-1` | the code tests of `internal/cli`, shard 1 | 350 s |
| `cli-2` | the code tests of `internal/cli`, shard 2 | 350 s |
| `contracts-repo` | the contracts under the `repocontract` tag | 330 s |
| `contracts-docs` | the contracts under the `docscontract` tag, `coverage-record-check`, `spec-budget`, `spec-check` | 380 s |

Each estimate includes 30 s for checkout and setup, is sized for a change that selects every package, and is below the 480 s median with at least a fifth to spare. They rest on one run, and the delivering pull request is their first measurement.

- The workflow keeps one trigger, a pull request into main with the three event types it has today. The `push` trigger leaves, with every step that ran only for a push.
- One matrix job, `parts`, runs the five parts with `fail-fast: false` on `macos-26`, each from its own checkout with `GOFLAGS=-count=1`. Its timeout stays 60 minutes, and every part keeps the step that merges the current base branch.
- The gathering job keeps the id `verify` and the name `Verification gate`. It needs `parts`, runs with `if: ${{ always() }}` and fails unless the result of `parts` is `success`. Without that condition a failed part would skip the job, and GitHub reads a skipped required check as passing.
- A pull request runs the plan's packages and the Full Contract Run, which is the two contract parts together: every contract discovery finds, split by its build tag (ADR-0253: "prints one `go test` invocation per build tag for every contract that discovery finds"), each once. Continuous integration no longer calls `make verify-docs` by name: every check it lists runs in those two parts.
- No run happens on main, so no cache of the default branch exists for a pull request to restore, and the first run of each pull request compiles from nothing. That is the condition of the one measured run.
- The selector lists the top-level tests of `internal/cli` with `go/parser` and assigns each name to a shard by a hash of the name. Shard 1 runs `-run <pattern>`, and the last shard runs `-skip` with the union of the other patterns, so a test the listing misses still runs.
- Each part writes `go test -json` events to a file. `gate-time tests` prints the package lines, the output of each failed test and the twenty slowest tests, and writes that list to the job summary. The last step prints the Gate Time Report of the job and of its slowest test binary.

### The release

This subsection is the design the answer on the release granted on 2026-11-25. The job `verify` of `release.yml` keeps its id, and the publication job still needs it and keeps every line. The job stops running tests, loses the override and its comment, and runs on `ubuntu-latest`, because no test runs in it, with `contents: read`, `checks: read` and `pull-requests: read`. It keeps the checkout with the whole history and loses the Go setup. Its one step proves the **Release Provenance** of the release commit, which is the commit that last changed `dist/npm/roundfix/package.json`:

1. A pushed tag points at the release commit. A manual dispatch evaluates the release commit of the ref it runs on.
2. GitHub lists one pull request for that commit, merged into main, whose merge commit it is.
3. The commit's parent is an ancestor of the pull request's head, and the two have the same tree. The branch therefore held main's tip when it merged, the step that merges the base changed nothing, and the check ran on this tree.
4. The latest check run named `Verification gate` of that head, from GitHub Actions, completed with success.

Any other state fails the job, and so does an answer GitHub fails to give; nothing is built. Rule 3 needs no record of the run: Git proves it. Because the release commit changes a path outside every package, the check that passed ran every package and every contract on the macOS runner.

Measured on 2026-11-25 over the last fifteen release tags: fourteen meet the four rules, and `v0.148.1` fails rule 3, because its pull request merged behind main. The repository's branch protection does not require an up-to-date branch, so the runbook tells the operator to update the branch before merging and to run the manual dispatch before pushing the tag. A refused release commit is never published; the release is cut again as the next patch version.

### The numbers and the closing change

`docs/references/gate-time-targets.json` holds every number a gate reads, and `docs/references/gate-time-targets.md` is the reference document of Core Feature 1. Neither is a Governed Path. The `Makefile` and the workflows hold no number: they call `gate-time report` with a gate's id.

Each gate's row names its machine, its approved target and the window that judges it. Two gates have a window. `gate-time windows` reads the first 30 passing settlement checks after the window's start from the Run Database, and the first 20 passing runs of a pull request from GitHub, and prints each median and p95 against the target. The time of a run is the longest of its parts, each from its own start to its own end. The operator records the window's start, and later its result, in the JSON file. No run fails for its time: before a release the operator reads the windows, and a gate over its number holds the release.

The closing change the maintainer decided therefore edits documents only. It adds the Gate Time Rule as a new ADR, records each window's measured result and number in the two reference files, gives ADR-0453 its amendment notice and adds the glossary's link. The rule states that continuous integration runs only for a pull request, as his answer asks; ADR-0453 records it first, for the interval. The Gate Time Rule amends ADR-0453 and does not supersede it: ADR-0453 stays the record of the retirement that the amended ADRs cite, and deleting it would force the closing change to edit a test.

## Implementation Design

### Interfaces

```go
// internal/verifyselect
type Mode string // "settlement", "complete", "reach"

type Facts struct {
	RunMarker string   // ROUNDFIX_RUN_ID
	Branch    string   // the checked-out branch, empty when detached
	Changed   []string // ChangedPaths, as today
	Added     []string // the paths of Changed that the base does not hold
}

type Plan struct {
	Mode     Mode
	Reason   string   // the line printed to standard error
	Packages []string // "./internal/x"; empty at a settlement
	Fresh    bool
}

func PlanPackages(facts Facts, module []goListPackage, root string) (Plan, error)
func CLIShard(repoRoot string, shard, of int) (flag, pattern string, err error)
```

```go
// internal/gatetime
type Limits struct{ Seconds, MedianSeconds, P95Seconds int }

type Gate struct {
	ID, Name, Machine string
	Approved          Limits
	Measured, Number  *Limits
	Window            Window
}

func Load(path string) (Document, error)
func Number(measured, approved int) (number int, reached bool)
func Validate(doc Document, approved map[string]Limits) []error
func Report(gate Gate, seconds int, job string) string
func Windows(doc Document, source Source) ([]string, bool)
```

```text
1. PlanPackages returns "settlement" only when RunMarker is not empty and Branch is "roundfix/run-"+RunMarker or begins with that string followed by "-".
2. PlanPackages returns "complete" when the mode is not "settlement" and Added holds a path that matches docs/history/specs/*.md.
3. For any shard count, every top-level test of internal/cli is run by exactly one shard.
4. Number returns the smaller of approved and measured × 1.25 rounded up to the next ten seconds; reached is false when measured is above approved, and the gate then holds no number.
5. Validate fails on an approved value that differs from the one the contract holds, on a number above its approved value, and on a number that differs from Number(measured, approved).
6. gate-time opens the Run Database for reading only.
```

### Data Models

`docs/references/gate-time-targets.json`, schema version 1: `allowancePercent` (25), `rule` (null, then the ADR's number), and `gates`. A gate has `id`, `name`, `machine`, `approved`, `lab`, `window` (`kind`, `event`, `size`, `since`, `passing`, `closed`), `measured` and `number`. The ids are `task-settlement`, `complete-suite-run`, `before-pull-request`, `pull-request` and `test-binary`. `measured`, `number`, `since` and `closed` are null in the first delivery.

The Run Database is read, never written: the events of kind `daemon.verification` of this repository's Runs, pairing each `started: make verify-changed` with the next `command passed` of the same Task and attempt. This is the query the PRD's baseline used.

### API Contracts

1. API Contract: `verify-select -plan [-base <ref>]`. Standard output: one package argument per line, or nothing. Standard error: one line that begins `verify-select: packages:`. Exit 0; exit 1 when the module cannot be listed. When Git cannot list the changes it selects every package and says so.
2. API Contract: `verify-select -job packages|cli-1|cli-2 [-base <ref>]`. Standard output: the `go test` arguments of that part, or nothing when the plan selects none.
3. API Contract: `verify-select -contracts [-all] -tag docscontract|repocontract`. As today, restricted to one build tag. Without `-tag` it prints what it prints today.
4. API Contract: `make verify-changed`. Prerequisites `fmt-check`, `vet` and `build`, as today. It runs the plan's packages with `-parallel $(GO_TEST_PARALLEL)`, adding `-count=1` when the plan is fresh, then `verify-changed-contracts` and `coverage-record-check`. In the `settlement` and `complete` modes it ends with a Gate Time Report, and `make verify-docs` ends with the report of the gate `before-pull-request`. Neither fails for its time.
5. API Contract: `make verify-part PART=<part>`. It runs one row of the table of `### The jobs`; a contract part runs every contract of its tag. No invocation carries a longer `-timeout` than today: none for code tests, `20m` for contracts.
6. API Contract: `gate-time report -gate <id> (-seconds <n> | -since <unix seconds>) [-job <name>] [-targets <file>]`. One line on standard output. Exit 0 whatever the time; exit 2 for an unknown gate.
7. API Contract: `gate-time tests -events <file> [-top <n>] [-summary <file>]`. It renders `go test -json` events as text and lists the slowest tests. Exit 0; the recipe keeps the status of `go test`.
8. API Contract: `gate-time windows [-fixture <file>] [-latest]`, also `make gate-times`. One line for each gate that has a window. Exit 0; exit 1 when the document names a rule and a closed window is over its number.
9. API Contract: the required check `Verification gate` of `ci-verify.yml` succeeds only when every part succeeded, and the workflow has no trigger but `pull_request`.
10. API Contract: the step `Release provenance` of `release.yml`'s job `verify`. Input: the commit and the event of the run, and the workflow's token for `gh api`. On success it prints `provenance: <commit> is the squash merge of pull request #<n>, and Verification gate passed on its tree` and exits 0. Otherwise it prints one line `::error::provenance: <reason>` and exits 1. The reasons: the tag is not on the release commit; the commit is not the squash merge of one pull request into main; the pull request merged behind main; the check did not pass or is missing; GitHub could not be read.

### Surface Transcripts

The product's command surfaces do not change. These transcripts fix the lines of the PRD's `## User Experience` for the two repository tools. The QA gate runs the first from the Run Worktree, where the Run Marker and the Run's branch hold.

1. Surface Transcript: the package plan at a settlement.

   ```transcript
   $ go run -buildvcs=false ./cmd/verify-select -plan
   stdout:
   stderr:
   verify-select: packages: none at a settlement; the Complete Suite Run runs them once per Spec
   exit: 0
   ```

2. Surface Transcript: the Gate Time Report of a Task settlement.

   ```transcript
   $ go run -buildvcs=false ./cmd/gate-time report -gate task-settlement -seconds 31
   stdout:
   gate time: task settlement 31s (target: median 120s, p95 240s)
   stderr:
   exit: 0
   ```

3. Surface Transcript: the Gate Time Report of the Complete Suite Run.

   ```transcript
   $ go run -buildvcs=false ./cmd/gate-time report -gate complete-suite-run -seconds 352
   stdout:
   gate time: complete suite run 352s (target: 360s)
   stderr:
   exit: 0
   ```

4. Surface Transcript: the Gate Time Report of a verification job.

   ```transcript
   $ go run -buildvcs=false ./cmd/gate-time report -gate pull-request -job cli-2 -seconds 412
   stdout:
   gate time: pull request, job cli-2 412s (target: median 480s, p95 600s)
   stderr:
   exit: 0
   ```

5. Surface Transcript: the Gate Time Report once a rule and a number are recorded.

   ```transcript
   $ go run -buildvcs=false ./cmd/gate-time report -targets internal/gatetime/testdata/targets-with-rule.json -gate task-settlement -seconds 31
   stdout:
   gate time: task settlement 31s (number: median 120s, p95 240s; within the measured result)
   stderr:
   exit: 0
   ```

6. Surface Transcript: the windows, from recorded runs and settlements.

   ```transcript
   $ go run -buildvcs=false ./cmd/gate-time windows -fixture internal/gatetime/testdata/windows-fixture.json
   stdout:
   gate times: task settlement, 30 of 30 passing: median 104s of 120s, p95 212s of 240s, within
   gate times: pull request, 14 of 20 passing: window open
   stderr:
   exit: 0
   ```

A single run is compared with the p95 pair: `within the measured result`, `in the allowance` or `over the number`. The fixture of transcript 5 records the measured median of 104 s and p95 of 212 s that transcript 6 prints.

## Coverage Map

- Goal 1 → What a settlement runs, The jobs, The numbers and the closing change: each gate's composition, the Gate Time Targets files, `gate-time windows`
- Goal 2 → The jobs, The release: the one trigger, the matrix job, the gathering job, the shards, the Release Provenance
- Goal 3 → Where the Complete Suite Run lives, What a settlement runs: the package plan
- Goal 4 → The jobs, API Contracts: `gate-time report`, `gate-time tests`
- Goal 5 → The numbers and the closing change: the Gate Time Targets contract, the closing change
- Story 1 → the package plan, `make verify-changed`
- Story 2 → the matrix job, the package plan (`reach`), the shards
- Story 3 → the one trigger, the Full Contract Run on every pull request, the Release Provenance
- Story 4 → ADR-0453, the Gate Time Targets contract, the closing change
- Story 5 → `gate-time report`
- Story 6 → the Nested Go Tool rule, the One Build rule, the shared pinned history
- Core Feature 1 → `docs/references/gate-time-targets.json` and `.md`
- Core Feature 2 → the `settlement` mode, the rule tests as `//verify:always` contracts, `coverage-record-check`
- Core Feature 3 → the `complete` mode at the Delivery Queue's repository gate
- Core Feature 4 → the `reach` mode, the Full Contract Run in the two contract parts
- Core Feature 5 → the one trigger, the matrix job, the gathering job, the shards and their test
- Core Feature 6 → the step `Release provenance` of `release.yml`'s job `verify`
- Core Feature 7 → ADR-0453, its notices and its contract
- Core Feature 8 → `gate-time report`, `gate-time tests`, `gate-time windows`
- Core Feature 9 → step 1: the four nested test runs, the shared listing, the Nested Go Tool rule
- Core Feature 10 → `testfixture.ProductBinary`, the memoized `testfixture.FixtureBinary`, the One Build rule, the shared `gittest.PinnedHistory`
- Core Feature 11 → the prepared Baseline fixtures
- Core Feature 12 → the event-bound waits
- Core Feature 13 → the list of names in step 1 to step 4, the sanctioned re-record
- Core Feature 14 → the delivering pull request under the new jobs
- Core Feature 15 → the closing change
- Core Feature 16 → `gate-time windows` and its exit status, the guide's sentence, the release runbook
- Core Feature 17 → the guide's sentence, written by step 4
- Core Feature 18 → the Gate Time Targets files and their contract
- Core Feature 19 → step 4: `CONTEXT.md`, `docs/agents/specific-repository.md`, the Makefile's comments
- Success Metric 1 → the `settlement` mode
- Success Metric 2 → the `settlement` mode with the `//verify:always` contracts and `coverage-record-check`
- Success Metric 3 → the `reach` mode
- Success Metric 4 → the `complete` mode
- Success Metric 5 → the Nested Go Tool rule and the One Build rule, with their seeded cases
- Success Metric 6 → the contract `TestCLIShardsCoverEveryTestExactlyOnce`
- Success Metric 7 → the workflow tests of `internal/verifyselect`
- Success Metric 8 → ADR-0453 and its contract
- Success Metric 9 → the sanctioned re-record of the coverage record
- Success Metric 10 → the tests of `gate-time report` and of the Makefile, and the windows fixture
- Success Metric 11 → `gatetime.Validate` and the Gate Time Targets contract
- Success Metric 12 → the reference document and the two workflows
- Success Metric 13 → the cases of the release harness for the step `Release provenance`

Deviations from the PRD:

- Core Feature 13 — the coverage record also loses the two tests of the suite-time budget, which leave with ADR-0387, and gains the tests of the new code. The maintainer granted this wider bound on 2026-11-25, as addition (i).
- Core Feature 9 — `TestPartitionFollowsTheMakefileRecipes` and `TestPartitionCoversEveryTestExactlyOnce` change in step 3, with the sets they test, and not in step 1. Step 1 also repairs `TestDoctorAdapterNamesTheCursorLogin`, which compiles the test binary of another package and which the PRD does not name.
- Success Metric 4 — the QA gate proves it with the gate's command in three states of a scratch copy, not with a real item of the Delivery Queue, which needs Agents and a pull request.
- Success Metric 13 — the QA gate proves it with the harness and its recorded answers. The proof against GitHub is the manual dispatch the PRD's `## Unreachable Acceptance` declares.

## Absence Evidence

1. The gate-time tool, `cmd/gate-time` and `internal/gatetime` — command: `ls cmd internal | grep -i gate; grep -rn 'gate time' --include='*.go' --include=Makefile --include='*.yml' . | grep -v docs/`; result: no output; no gate prints its time and nothing reads a window.
2. The Gate Time Targets files — command: `ls docs/references | grep -i gate`; result: no output.
3. The package plan and its two facts — command: `grep -n 'ROUNDFIX_RUN_ID\|docs/history' Makefile internal/verifyselect/verifyselect.go internal/verifyselect/contracts.go cmd/verify-select/main.go`; result: no output; the selector reads neither the Run Marker nor an Archive Record, and its only use of imports is `coreImportsBaseline`.
4. The shards and the tag filter — command: `grep -n -i 'shard' Makefile .github/workflows/*.yml internal/verifyselect/*.go; grep -nE '"(plan|job|tag)"' internal/verifyselect/verifyselect.go`; result: no output; the only split of the CLI package is the Baseline pattern.
5. The matrix job and the gathering job — command: `grep -nE 'matrix|needs:' .github/workflows/ci-verify.yml`; result: no output; one job runs every step.
6. The Nested Go Tool rule and the One Build rule — command: `grep -rln -i 'nested' internal/testfixture; grep -rn 'ProductBinary' internal cmd`; result: one file, `internal/testfixture/real_time_budget_test.go`, where the word names a test directory; no rule reads a nested `go` call and no helper builds the product binary once.
7. The shared pinned history — command: `grep -nE 'sync\.|Once' internal/gittest/pinned_history.go`; result: no output; each call extracts again.
8. ADR-0453 — command: `ls docs/adr | grep -E '^045[3-9]'`; result: no output; the highest number is 0452.
9. The glossary terms — command: `grep -nE 'Gate Time|Complete Suite Run|Measured Dossier|Release Provenance' CONTEXT.md`; result: no output.
10. The Release Provenance — command: `grep -n -i -E 'provenance|check-runs|/pulls' .github/workflows/release.yml internal/verifyselect/*.go`; result: no output; the release reads nothing of a pull request.

## ADR Alternatives

- ADR-0452: no rejected alternative is taken. Every job that runs a test names `macos-26`, no larger runner is used (ADR-0452: "a larger runner was rejected without a measurement"), and the publication job keeps its runner (ADR-0452: "Running the whole release on macOS was rejected to preserve the proven publication job and its permissions."). ADR-0453 amends it by withdrawing the one job, the runner-type sentence and the 420 seconds. By the answer on the release it also withdraws the sentences on the release's verification job, with the release clause of ADR-0253. Rule Challenge of 2026-11-25: retire.
- ADR-0387: superseded by ADR-0453 and deleted. The design takes an alternative it rejected (ADR-0387: "A separate job for the Full Contract Run spends a second runner, checkout and build on every push"). Rule Challenge of 2026-11-25: retire.
- ADR-0453: written by step 4. It records the retirement, the three limits of the interval, what each gate runs in it and that continuous integration runs only for a pull request, and it rejects nothing this design takes.
- ADR-0259: the CLI package stays one package, run as two shards of one test binary (ADR-0259: "Splitting `internal/cli` into several packages would let `go test` overlap them"). `-parallel` stays 16 and no ceiling rises. The prepared Baseline fixtures copy the result of one sync and of one plan build that a test only reads (ADR-0259: "Fixtures copy what a test reads."); the `git init` of the suite's fixtures is not replaced by a template (ADR-0259: "Fixture templates copied instead of recreated were measured and left out"). ADR-0453 replaces its notice about the budget.
- ADR-0257: the Agent of a Run still runs neither selected command (ADR-0257: "Remove the runtime decisions. Rejected because of the adopter impact."). ADR-0453 gives it a notice: in this repository a settlement runs no package test, by the maintainer's answer "Só o leve" in the PRD's `## Decisions`, and continuous integration runs nothing on a push to main.
- ADR-0182: it records no rejected alternative. ADR-0453 gives its cost sentence a notice, by the same answer.
- ADR-0342: the per-Task command stays (ADR-0342: "Only the Integrated Repository Verification, no per-Task command."). The reason that alternative was rejected returns for package tests (ADR-0342: "a failure found after the last Task has no Agent Session to repair it"): the answer "Só o leve" accepts it, and the PRD's first risk counts each case.
- ADR-0253: amended by ADR-0453. The run on every push to main ends by the maintainer's answer of 2026-11-25, in the PRD's `## Decisions`, which needs no question. The design keeps the run and its source: every contract discovery finds runs on every pull request, with no list kept by hand (ADR-0253: "A hand-kept list of every contract for main would drift as `REPO_CONTRACT_TESTS` did.") and no switch that widens `make verify-docs`. It moves the net before the merge (ADR-0253: "A miss of the selection rule now surfaces on the push to main, after the merge, rather than before it"). The release clause leaves by the answer on the release, the Praetor's application of the maintainer's decision. Rule Challenge of 2026-11-25: change.
- ADR-0252: the directives and the classes are unchanged, and each new contract is declared beside its test (ADR-0252: "Keeping a central list of tests and paths in `cmd/verify-select` would drift from the tests it describes."). One consequence stops being true (ADR-0252: "so CI runs the selected contracts twice"), and ADR-0453 gives it a notice. Rule Challenge of 2026-11-25: change.
- ADR-0326: the moved tests stay in their packages (ADR-0326: "they exercise unexported `internal/cli` code"). The last shard's `-skip` names what the other shard runs in the same event, so it hides no test from a gate (ADR-0326: "it hides them from every gate that does not repeat the list"). Its sentence on where the Full Contract Run runs follows the notice of ADR-0253.
- ADR-0271: no real-time allowance is added (ADR-0271: "Raising each budget until the suite passes."). The four waiting tests end on an event.
- ADR-0281: the one build of a package writes its executable once, to a new path, and a test that changes it copies it first (ADR-0281: "Removing the old file before copying.").
- ADR-0311: the check wait and the job timeout it was sized on stay (ADR-0311: "the job timeout of this repository").
- ADR-0283: the plan detects a Run from the Run Marker, as decided. The branch only confirms that the checkout is that Run's. A temporary repository of a QA row, the case the ADR rejected detection by directory for, falls to the wider selection.
- ADR-0285: not applicable — the compile contract keeps its class, its tag and its place in the Full Contract Run, in `contracts-repo`.
- ADR-0435: not applicable — the design changes nothing it decides. The gate after a Base Update runs the same command, so the Complete Suite Run repeats there.
- ADR-0258: not applicable — this design rewords no Baseline clause.
- ADR-0273: not applicable — this design sends nothing from a Run.
- ADR-0282: not applicable — the Delivery Queue's item binary is not touched.
- ADR-0286: not applicable — it decides the Correction Review of a fixed finding.
- ADR-0295: not applicable — the selector's base stays `VERIFY_BASE`.
- ADR-0298: not applicable — it decides the attempts of a profile proof.
- ADR-0319: not applicable — it decides the order of the follow-up code.
- ADR-0384: not applicable — the Machine Verification Capacity is unchanged.
- ADR-0385: not applicable — the Verification Core Budget is unchanged, and no recipe sets `GOMAXPROCS`.
- ADR-0404: not applicable — it decides what the Idea-stage judge reads.
- ADR-0423: not applicable — it decides the Correction Sign at a merge.

## Integration Points

- The Daemon and the Delivery Queue are unchanged. They run the configured command and set the Run Marker as today; the selector reads the environment and the checked-out branch.
- The Run Database is read by `gate-time windows` through a read-only connection.
- The GitHub CLI: `gate-time windows` calls `gh run list` and `gh run view --json jobs` for the pull request runs of the workflow `CI | Verify` and takes each job's start and end. A part prints its report as a `::notice` annotation, as the workflow already does for the tested base.
- GitHub's REST API, from the release workflow: `gh api` reads the pull requests of a commit and the check runs of a head, with the workflow's token. `gh` is on every hosted runner.
- GitHub's branch protection is unchanged, because the required check keeps its name. It does not require a branch to be up to date, which is why the Release Provenance checks that for the release commit.
- The Delivery Queue reads the checks of a pull request's head and expects a run of each workflow that triggers for a pull request (`internal/delivery/head_check_set.go`). It reads nothing of a push to main, so the removed trigger changes nothing there.
- The Secondbrain export gains no line: ADR-0453 is marked `knowledge: false`, because it is about how this repository is tested.

## Testing Approach

Existing seams are used, and two are new.

- The selector: `PlanPackages` is a pure function over `Facts`, so its table tests need no Git and no environment. One Git-backed test with `gittest` proves `Added` on a commit that adds an Archive Record, and one table case pins that the version files of a release select every package. `TestCLIShardsCoverEveryTestExactlyOnce`, a `//verify:always` contract, computes each shard's tests from the patterns and the listing, and a seeded name in two shards fails it.
- The Makefile: the existing seam of `internal/verifyselect` runs the real Makefile with a stub `go`. New cases prove that the `settlement` mode starts no package test, that the `complete` mode runs `test -count=1 -parallel 16 ./...`, and that each gate prints its whole Gate Time Report, with one named test for each gate.
- The workflows: the existing tests that parse the two files are rewritten. `ci-verify.yml` triggers on `pull_request` and on nothing else; every job that runs a test names `macos-26`; the gathering job needs `parts`, always runs and fails on any result but success; no job of `release.yml` runs a test; the publication job needs `verify` and keeps its runner, its permissions and its steps; no step sets a suite-time budget.
- The Release Provenance: the harness of `internal/verifyselect` that runs a step of the release workflow with stub commands learns to run a step of the job `verify`. Each case builds a small repository with `gittest` and a stub `gh` that answers from recorded files: a squash merge of a branch that held main's tip passes; a merge behind main, a failed check, a missing check, a commit no pull request merged, a tag on another commit and an unreadable answer each fail with their line.
- The gate-time tool: unit tests cover `Number`, `Validate`, `Report` and `Windows`. The two new seams are `-fixture`, which gives `windows` recorded runs and settlements, and `-targets`, which gives `report` a document that holds a rule. No test calls GitHub or opens the Run Database.
- The live files: `internal/docscontract/gate_time_targets_test.go` applies `Validate` to `docs/references/gate-time-targets.json`, with the approved targets as constants, and proves the reference document's table holds the same numbers.
- The two new rules follow the six existing ones: a walk over the test sources, a frozen inventory with a reason for each entry, and a seeded violation.
- The repaired tests keep their assertions. Step 2 makes each shared directory unwritable, so a test that writes into it fails.
- The record of ADR-0453: a contract in `internal/docscontract` pins what it supersedes and amends and does not forbid a later `amended_by`, so the closing change passes it unchanged.

The QA gate runs the lab protocol of the reference document, the control cases of Success Metric 2 on a scratch copy and the transcripts. It proves Success Metric 4 with the same command in three states of a scratch copy: twice on a Run's branch with the marker, where no package test starts, and once on a branch that adds an Archive Record, where every package runs once.

## Build Order

1. Code tests stop running the Go tool over other packages and building per test.
   - `TestCoverageCollectionMatchesTheToolchainOnThisPlatform` and `TestJournalConsumerCorpusReplaysEveryConsumer` become Repository Contract Tests under the `docscontract` tag, `//verify:relevant` with the inputs their Task lists.
   - `TestDerivedQAVerificationPassesTheChecker` is merged into `TestDerivedQAVerificationIsAcceptedByTaskStageChecker`, the test it starts.
   - The two tests that collect the coverage record share one collection for each run of their package.
   - `testfixture.ProductBinary` builds `cmd/roundfix` once for a test binary, and `testfixture.FixtureBinary` builds each stub once. The tests that build inside a test take them.
   - The two rules, their sentences in the repository's own guide, and the re-record of the coverage record.
   - Governed Paths: `internal/spec/coverage_test.go`, `internal/cli/baseline_plan_test.go`, `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`.
2. Tests share what they only read, and wait on events (depends on: 1 — both edit the test fixture package and re-record the coverage record).
   - `gittest.PinnedHistory` extracts each set of paths once for a test binary, read-only. Its signature and its callers do not change.
   - The asset-sync tests that assert a later step take a prepared copy of the synced state, and `TestTransactionFailureMatrix` copies one prepared repository for each fault point.
   - The detach test and the three helper-residue tests end on an event the test controls.
   - Governed Paths: `docs/references/coverage-record.json`.
3. A Task settles on the light checks, the Complete Suite Run gets its home, and each gate reports its time (depends on: 1, 2 — the rule tests exist, and the lab numbers need both).
   - The package plan, `-job`, `-tag`, the shards and their test. The tests of the retired sets are renamed for what they prove; the Task Graph lists every name.
   - The rule tests and the two Baseline-record tests of the CLI package become `//verify:always` contracts.
   - The Makefile: `verify-changed` by plan, `coverage-record-check`, `verify-part`, `gate-times`, and the reports of `verify-changed` and `verify-docs`.
   - `cmd/gate-time`, `internal/gatetime`, the two Gate Time Targets files and their contract.
   - Governed Paths: `Makefile`, `docs/agents/setup-context.json` by its command, `docs/agents/specific-repository.md`, `.roundfixrc.yml` (the comment), `docs/references/coverage-record.json`.
4. Continuous integration runs only for a pull request and in parallel jobs, a release proves its provenance, and two time rules are retired (depends on: 3).
   - The two workflows and their tests: `ci-verify.yml` for a pull request only, and the Release Provenance in `release.yml`.
   - The suite-time budget leaves the Makefile, with its two tests.
   - ADR-0453, accepted, with the deletion of ADR-0387, the notices in ADR-0452, ADR-0253, ADR-0252, ADR-0259, ADR-0257 and ADR-0182, and its contract.
   - The glossary terms and the changed entry of the Full Contract Run, the guide's sentences on what each gate runs and on the Gate Time Rule, and the release runbook.
   - Governed Paths: `.github/workflows/ci-verify.yml`, `.github/workflows/release.yml`, `Makefile`, `docs/agents/setup-context.json` by its command, `docs/agents/specific-repository.md`, `docs/references/coverage-record.json`.
5. QA gate (depends on: 1, 2, 3, 4).

The Gate Time Rule is not a step. It is the closing change, after the windows.

## Risks & Considerations

- The estimates of the jobs rest on one run with a cold cache, and the shards are balanced by count, not by time. If the delivering pull request shows a part above 480 s, the shard count or the split of the contract parts changes before the merge, inside the same grant.
- A settlement may not fit 45 s. A change to the CLI package recompiles its tagged test binary; the design keeps that to one binary, and the QA gate measures it.
- The `settlement` mode rests on the Run's branch name. A test of the selector pins the prefix against `store.RunBranchPrefix`, so a rename in the product fails here first.
- The parts of one pull request merge the base separately. If main moves between two checkouts, the parts test two bases; each annotates its own, and the item is then behind its base and returns to the gate.
- Nothing verifies main after a merge, and a commit pushed to main without a pull request is verified by nothing. The next pull request that reaches a broken package finds it, and at the latest the release's pull request, which runs everything and must hold main's tip.
- The Release Provenance depends on GitHub's API when the release runs. An unreadable answer refuses the release, and running the job again repeats the check. A release commit that merged behind main costs a patch number.
- Pins move. A contract file added or removed changes the counts `TestRepositoryContractTestsDeclareTheirRelevance` pins. The tests that run the real Makefile need `test` to stay `$(GO) test -parallel $(GO_TEST_PARALLEL) $(PKGS)` and `GO := go`. A live-document contract forbids the retired lane's words in the Makefile and the workflows.
- `make verify` no longer runs the rule tests, because they become contracts. `make verify-changed` runs them in every mode, and the guide's rule already puts checks of the repository itself outside `make verify`.
- Step 2 may change a source of a command, in `internal/cli/detach.go`. The record of command surfaces is then re-recorded by its own command, which the maintainer granted as addition (ii).

## Glossary

- adds: **Release Provenance**
- changes: **Full Contract Run**
- not a term: **Nested Go Tool rule** — the name of a repository rule test, like the Hard Link rule
- not a term: **One Build rule** — the name of a repository rule test

## Decisions

- The Complete Suite Run lives at the Delivery Queue's repository gate, with no configuration key and no Daemon change.
- The checks of the repository against its own rules become `//verify:always` contracts. The comparison of the coverage record stays a code test that the Makefile names.
- The suite-time budget's targets and variable leave the Makefile.
- The Gate Time Rule amends ADR-0453.
- Five parts and two shards. The gathering job also names `macos-26`; it starts when the parts have ended, so one event never holds more than five macOS jobs. `GO_TEST_PARALLEL` stays 16.
- Continuous integration runs only for a pull request: the maintainer's answer of 2026-11-25, message `01a213bf-b859-7f60-91a1-201ca006c98a`, in the PRD's `## Decisions`. The targets, the window and the report of a push to main leave with it.
- The release commit is the commit that last changed the version file, and the Release Provenance is a step of the workflow, not a command of the product or a target of the Makefile. Both follow the answer on the release.
- The hold on a release is the exit status of `gate-time windows`, the guide's sentence and the runbook. The product's release plan does not change.
- Delegated Decision of 2026-11-25: Which packages does a change to a documentation tree select on a pull request?
- Delegated Decision of 2026-11-25: What does a review Batch run?
- Delegated Decision of 2026-11-25: Where does the Full Contract Run happen once continuous integration runs only for a pull request?

Sources read on 2026-11-25, and what each changed:

- GitHub Docs, "Troubleshooting required status checks", https://docs.github.com/en/pull-requests/how-tos/merge-and-close-pull-requests/troubleshooting-required-status-checks: a job that depends on a failed job is skipped and may not block merging. It is why the gathering job always runs.
- GitHub Docs, "Actions limits", https://docs.github.com/en/actions/reference/limits: five concurrent macOS jobs on the account's plan. It bounds the parts at five.
- The Go command, `go help testflag`: `-skip` uses the matching rules of `-run`. That one pattern splits the top-level tests in two is the author's inference, so a test proves it.
- Three published shard scripts, read as excerpts by the Research Subagent (Mattermost, gascity, Pulumi), list tests with `go test -list`, which compiles the package and runs its `TestMain`. The selector lists from the source with `go/parser` instead. They are examples, not primary sources.
- GitHub Docs, "About protected branches", read through Exa on 2026-11-25, https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches: a check that does not require an up-to-date branch may fail after the merge. With the repository's own setting, read through `gh api`, it is why rule 3 of the Release Provenance exists.
- GitHub Docs, the REST endpoints for commits and for check runs and the permissions each needs, read by the Research Subagent: they give the two calls of the Release Provenance and its two read-only permissions. Whether the calls answer without a token is not stated, so the job declares both.
- GitHub Docs, "Dependency caching reference": a pull request restores caches of its own ref and of the default branch only. It is why the estimates assume a cold cache.
- Secondbrain, `concepts/tags-protegidas-e-release-por-workflow-no-github.md` and `projects/quaestor-legacy/mirror/docs/references/ci-repetido-no-release-2026-11-16.md`, opened on 2026-11-25: a sibling repository moved its verification to the pull request and made main require an up-to-date branch. They support the recommended option and name its condition, which this repository checks in the release because a setting on GitHub is outside the Spec. The PRD's `references/sources.md`, section 9, holds every passage and measurement of this revision, and their digest went to the Secondbrain's owner as ai-memory message `01a213d5-8842-7002-bd41-ba95de2bec84`.
- Secondbrain, through the Indexus tool `pesquisar_acervos`: `sources/digest-roundfix-ci-rerun-merge-antigo-2026-10-06.md`, opened, is why every part keeps the step that merges the base. A second search, on light checks for each Task, returned no relevant page, and `wiki/index.md` did not open through the tool; both are limits of this consultation. The digest of these sources went to the Secondbrain's owner as ai-memory message `01a212c5-99e0-7bc2-877e-96845761f9a3`.

The two other options of the release question, a minimal publication set and the whole verification as first granted, were not taken. The answer gives the reasons: the first publishes a tag on any commit without proof of a pull request, and the second runs a verification after the merge.
