---
task: task_02
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
status: completed
type: backend
complexity: medium
---

# Task 02: Doctor and the Delivery Queue warn below the Free-Space Floor

## Overview

One location inspection, `inspectFreeSpace`, reads the free space where the
worktree location, each declared Shared Build Cache and the temporary
directory live. Doctor's `disk-space:` line uses it and so also names the
temporary directory, and the Delivery Queue uses it before an item's
Implement Run and before its repository gate to record a `low-disk` warning
segment that `roundfix deliver status` prints, without ever parking or
delaying the item. The `doctor` and `deliver` documents and `CONTEXT.md`
describe the Free-Space Floor.

## Requirements

1. MUST add `inspectFreeSpace(ctx, loaded)` to
   `internal/cli/doctor_disk_space.go`, returning in order the worktree
   location, each declared Shared Build Cache and the temporary directory
   (`os.TempDir()`, label `temporary directory`), each read at its nearest
   existing directory through the existing `availableBytes` command
   dependency, and make `defaultDoctorDiskSpaceResults` render its results
   with today's codes, texts, threshold, `next: roundfix gc` and exit
   behavior, as `_techspec.md` → API Contract 4 says.
2. MUST change `TestDoctorWarnsWhenALocationIsLowOnDisk` to expect the third
   read and its `DR-DISK-LOW` or `DR-DISK-UNREAD` finding, and update every
   other test that pins the `disk-space:` line on the default branch when
   the Task runs.
3. MUST add to `internal/delivery/engine.go` the `FreeSpaceChecker`
   interface, `EngineDependencies.FreeSpace`, `WarningLowDisk = "low-disk"`
   and the check API Contract 5 describes: called at the start of
   `runCandidate`, before `RunSpec`, and of `gateCandidate`, before `Gate`;
   one `low-disk: <locations joined by ", ">` segment that replaces an
   earlier one and stays before a `linux-lane-skipped:` segment; the item
   stored when its warning changes; the log line
   `roundfix: warning: Delivery Queue item <slug>: low-disk: <locations> before its Implement Run; the queue continues`
   (or `before its repository gate`) word for word; the segment removed when
   the result is empty; nothing changed with a nil checker. It MUST NOT park,
   fail or delay the item, and `LaneSkipBody` and the pull request body keep
   their bytes.
4. MUST add `internal/cli/deliver_free_space.go` with the command
   implementation of `FreeSpaceChecker`, which turns each `inspectFreeSpace`
   location below 10,240 MiB into `<label> <path> has <N> MiB free, below 10240 MiB`
   and skips an unreadable one, and wire it into
   `commandDeliveryEngineDependencies` (`internal/cli/deliver_workflow.go`).
5. MUST add `TestCommandFreeSpaceNamesEachLocationBelowTheFloor` in
   `internal/cli/deliver_free_space_test.go` and
   `TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate` in
   `internal/delivery/low_disk_warning_test.go`, with the rows `_techspec.md`
   → Testing Approach names, following each package's rule on
   `t.Parallel()`, with the injected `availableBytes` and the engine's fake
   workflow and real queue store; no test reads the real disk's free space.
6. MUST write, word for word as `_techspec.md` → API Contract 8 gives them,
   the `disk-space:` bullet and `DR-DISK-LOW` row changes in
   `docs/user-guide/commands/doctor.md` and the canonical
   `.agents/skills/roundfix/references/setup.md`, and the paragraph starting
   "Before an item's Implement Run and before its repository gate" in
   `docs/user-guide/commands/deliver.md` and the canonical
   `.agents/skills/roundfix/references/deliver.md`; and add through
   `domain-modeling` the `CONTEXT.md` entry **Free-Space Floor** with its text
   and `_Avoid_` line, beside **Build Cache Bound**, with `ADR-0370` written
   as the glossary's Markdown link, adding no other capitalized phrase.
7. MUST run `make skills-sync` to copy the canonical references to
   `skills/roundfix/references/`, then raise the Roundfix skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last guide and skill edit; never
   hand-edit a mirror, a version field or a record. Of the Governed Paths the
   Task may change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/deliver.md`,
   `.agents/skills/roundfix/references/setup.md` and
   `skills/roundfix/SKILL.md`.
8. MUST keep edits in `internal/delivery/engine.go` and
   `internal/cli/deliver_workflow.go` local to `runCandidate`,
   `gateCandidate`, `EngineDependencies`, `NewEngine`, a new helper and
   `commandDeliveryEngineDependencies`: Spec 0359 edits `checkCandidate` in
   the same file.

## Subtasks

- [ ] Extract `inspectFreeSpace` with the temporary directory.
- [ ] Update the Doctor test for the third location.
- [ ] Add the queue seam, the `low-disk` segment and its log line.
- [ ] Add the command checker and wire it.
- [ ] Add the CLI and engine tests.
- [ ] Document the floor, add the glossary entry, sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Doctor names the temporary directory below 10,240 MiB and nothing at or above it, with an unchanged exit code.
- [ ] The queue records one `low-disk` segment before the Run and the gate when a location is below the floor, keeps it before a lane-skip segment, removes it when space returns, and the item still advances.
- [ ] Without a checker nothing changes.
- [ ] The `doctor` and `deliver` documents and `CONTEXT.md` describe the Free-Space Floor.

## Context

- instruction: `docs/adr/0370-the-makefile-bounds-its-build-cache-and-the-queue-warns-of-low-disk-and-names-a-full-one.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/delivery/lane_skip_warning_test.go`
- interface: `internal/cli/doctor_disk_space.go`
- interface: `internal/cli/doctor_disk_space_test.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_free_space.go`
- creates: `internal/cli/deliver_free_space_test.go`
- creates: `internal/delivery/low_disk_warning_test.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- Doctor's `disk-space:` line — consumers: `./internal/cli`
- `delivery.EngineDependencies` and the item warning — consumers: `./internal/delivery`, `./internal/cli`
- The `doctor` and `deliver` documents, the Roundfix skill version and the glossary — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestDoctorWarnsWhenALocationIsLowOnDisk|TestCommandFreeSpaceNamesEachLocationBelowTheFloor)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorWarnsWhenALocationIsLowOnDisk TestCommandFreeSpaceNamesEachLocationBelowTheFloor; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; grep -q 'temporary directory' internal/cli/doctor_disk_space.go || { echo 'no temporary directory location' >&2; exit 1; }; go test -count=1 -timeout 30m ./internal/cli && GOOS=windows go build -buildvcs=false ./internal/cli` — expected: exit 0. Before this Task the command checker test does not exist, so its pass line is missing; afterwards Doctor and the checker read the three locations, and the CLI suite passes and builds for Windows.

- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate ' || { echo 'missing pass: TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/delivery && GOOS=windows go build -buildvcs=false ./internal/delivery` — expected: exit 0. Before this Task the test does not exist; afterwards the queue warns before the Run and the gate without parking, and the delivery suite passes.

- `for f in docs/user-guide/commands/doctor.md .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'or the temporary directory has less than the Free-Space Floor of 10 GiB (10,240 MiB) free' || { echo "missing floor phrase: $f" >&2; exit 1; }; done; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -q -- 'it records a .low-disk: <label> <path> has <N> MiB free, below 10240 MiB. segment' || { echo "missing low-disk warning in $f" >&2; exit 1; }; done; grep -F -- '**Free-Space Floor**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Free-Space Floor\*\*(:| —( |$))' || { echo 'missing term: Free-Space Floor' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task no document names the temporary directory or the `low-disk` warning, so the first check fails; afterwards the documents, the mirrors, the glossary entry, the recorded skill version and the Behavior Surface fingerprint pass.

## References

- `_prd.md` → Goals 3 and 4; Core Features 2, 3 and 5; Success Metric 3; Success Metric 4; Success Metric 6; Glossary
- `_techspec.md` → Current behavior 3 and 4; Interfaces; API Contracts 4 and 5; API Contract 8; Invariants 4 and 5; Declared breaks; Testing Approach; Build Order 2
- ADR-0370; ADR-0334


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-context, golang-error-handling, golang-structs-interfaces, golang-safety, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, systematic-debugging, no-workarounds.

- Doctor: extracted `inspectFreeSpace`, retaining resolution errors and reading
  the worktree location, each declared Shared Build Cache and `os.TempDir()`
  in order through the injected capacity dependency at existing directories.
  Doctor retains its finding text, threshold, next action and exit behavior,
  with the additional `temporary directory` finding. Updated all pinned
  disk-space output expectations found in the CLI tests. Before the production
  change, the focused below-threshold test failed because it observed two
  reads rather than the required three; afterward all eight Doctor table rows
  passed, including unreadable, boundary and failed-exit cases.
- Queue warnings: added `FreeSpaceChecker`, its optional dependency and command
  implementation, and checks at the start of `runCandidate` and
  `gateCandidate`. The helper replaces/removes the single low-disk segment,
  stores changed warnings before the boundary runs, preserves the final lane
  reason and writes the exact warning line for each low observation.
  The six new engine rows exercise low space at both boundaries, recovery at
  the gate, lane ordering with a semicolon-containing reason, healthy removal,
  unchanged low space and a nil checker. Every row reaches `merged` through
  the real queue store; assertions inspect persistence before the Run, gate
  and publication and verify the exact PR body. `LaneSkipBody` and publication
  code are unchanged. Existing lane-skip publication tests also passed.
- Command checker: wired through `commandDeliveryEngineDependencies`; its
  three table rows verify the ordered location messages below 10,240 MiB,
  no warning at the floor, skipped unreadable observations and resolution to
  existing ancestors. Every capacity read in these tests is injected.
- Documentation: applied API Contract 8's Doctor bullet and finding-row text,
  Delivery Queue paragraph and Free-Space Floor glossary definition beside
  Build Cache Bound, with the accepted ADR link and required Avoid line.
  Focused glossary and command-documentation contracts passed. A read-only
  normalized-text check confirmed the Delivery paragraph word for word in
  the guide, canonical reference and mirror, and byte comparisons confirmed
  the canonical/mirrored skill files match.

Focused checks (all exit 0 after implementation):

- `go test ./internal/delivery ./internal/cli -run 'TestTheQueueWarns|TestDoctorWarnsWhen|TestCommandFreeSpace|TestAGateLaneSkip|TestAQuietGate' -count=1 -timeout 5m`
- `go test ./internal/cli -run 'Test(RunDoctor|Doctor|CommandFreeSpace|DeliverStatus|DeliverCommand)' -count=1 -timeout 5m`
- `go test -tags docscontract ./internal/docscontract -run 'TestGlossary|TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem|TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem' -count=1 -timeout 5m`
- `git diff --check`

Sanctioned regeneration (all exit 0, in this order after the final guide and
reference edits):

1. `make skills-sync`
2. `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   raised both Roundfix skill copies from `0.1.195` to `0.1.196` and recorded
   the version through its owner.
3. `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   regenerated only the changed Doctor and Deliver guide fingerprints.

The declared Verification commands, full package suites and Windows builds
remain Daemon-owned and were not run. Task status, the Task Graph and other
Task files were not edited; no commit, push or pull request was made.


### Verification feedback, attempt 1

Inspected the Daemon artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261106T150809Z_42a50d3dad874640/verification/batch-001-attempt-1.log`.
The configured gate remains unresolved: its reported failure is the
`internal/verifyselect` suite guard detecting a detached Git filesystem
monitor after the test assertions passed. The supplied transcript also
records the delivery package passing. No diagnostic log body is reproduced
here.

Investigation found a pre-existing launch route outside this Task's slice:
`TestPartitionFollowsTheMakefileRecipes` calls `makeRecipeCLITests`, which
runs the copied Makefile in the repository with an inherited environment.
The Makefile's immediately evaluated `BUILD_COMMIT` calls unqualified Git
status even for a dry run. A read-only query of the specific Git setting
confirmed that user configuration enables fsmonitor. The owning test helper,
Makefile expression and suite guard have no changes from this Task.

A disposable Python/subprocess probe copied the Makefile, replaced recursive
make as the existing test helper does, and supplied a recording Git stub,
disposable HOME/TMPDIR/cache paths and a dry run. It exited 0 and asserted
exactly the revision and status calls during Makefile parsing, without real
Git, a build, a real cache or a daemon. This establishes the unisolated call
route; the original artifact does not attribute the surviving helper to an
individual test.

Follow-up for the Supervisor: repair the verification fixture's inherited Git
configuration at its subprocess boundary, using the existing isolated Git
environment conventions in `internal/gittest`, and prove the owning suite
leaves no monitor. This is outside task_02; no Makefile, verification-test,
suite-guard or global Git configuration change was made, and no process was
signalled. The failure was not suppressed or classified as a passing gate.

Focused recheck:
`go test ./internal/cli ./internal/delivery -run 'TestDoctorWarnsWhen|TestCommandFreeSpace|TestTheQueueWarns' -count=1 -timeout 5m`
exited 0 for both packages. `git -c core.fsmonitor=false diff --check` also
exited 0. Only this Result evidence changed during the feedback turn; the
implementation, documentation and generated records retain their prior
bytes. The full configured Verification and declared Verification were not
rerun; Task status remains Daemon-owned.
