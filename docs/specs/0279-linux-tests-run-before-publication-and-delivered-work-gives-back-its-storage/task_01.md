---
task: task_01
spec: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
status: completed
type: backend
complexity: high
---

# Task 01: The Delivery Queue gate runs the changed Go packages in a Linux Lane

## Overview

On 2026-10-17 Spec 0276's item binary copy failed only in GitHub Actions on
Linux with `text file busy` (pull request #568), because no Run on the
maintainer's macOS machine compiles a `//go:build linux` test (`_prd.md` →
Sources 1). This Task adds the Linux Lane: when Project Config opts in, the
Delivery Queue's repository gate runs `go test` over the candidate's changed Go
packages in a `golang:<go.mod version>` container that pulls nothing, reaches
no network and receives only what `go test` needs, and skips with a named
reason when the machine cannot run it. This repository opts in. It is
verifiable on its own through gate tests that use a fake `docker`.

## Requirements

1. MUST answer the operator observation of 2026-10-17 (Spec 0276, pull request
   #568, CI run 39053329570) cited in `_prd.md` → Sources 1 by implementing API
   Contracts 1 to 5 of `_techspec.md` (ADR-0288).
2. MUST add the boolean Project Config key `delivery.linux_lane`, default
   `false`, read like `delivery.item_binary`: Project Config replaces User
   Config, and the Delivery Queue reads it from the fetched default-branch
   commit. A non-boolean value is refused with the config package's existing
   error shape.
3. MUST make `Gate` run the lane after the repository Verification passes and
   only when the key is `true` (API Contract 1). With the key unset, `Gate`
   starts no `docker` or `go` process and its output is unchanged.
4. MUST select the changed packages exactly as API Contract 2 states, from the
   item worktree's Git data only, and MUST check the skip conditions in the
   order of API Contract 3, writing each named reason to the Queue Owner's
   console log and to `<artifact dir>/delivery/<slug>/linux-lane.log`.
5. MUST run the container exactly as API Contract 4 states. The hazards:
   - The tree comes from `git archive` of the head the engine passes and goes
     to the container's standard input; nothing is written into the item
     worktree.
   - The image is used only with `--pull never` after `docker image inspect`
     succeeds; the lane never pulls.
   - The container gets `--network none`, `--user 65534:65534`, the host module
     cache mounted read-only, and only the variables `HOME`, `GOCACHE`,
     `GOMODCACHE`, `GOFLAGS`, `GOPROXY`, `GOTOOLCHAIN` and `GOSUMDB`. No judge
     key, token or other variable of the Queue Owner reaches it.
   - The docker CLI runs with the Queue Owner's own environment from
     `commandDependenciesForContext(ctx).environ()`, not
     `deliveryCommandEnvironment`, and `docker` and `go` are resolved on that
     environment's `PATH`.
6. MUST follow API Contract 5 for the outcome. The 20-minute cap cancels: it
   runs `docker rm -f <name>` and fails the gate. A cancelled context also runs
   `docker rm -f <name>` and returns the context's error. Exit `125` skips; any
   other non-zero exit fails the gate with the log path.
7. MUST add these tests. Each calls `t.Parallel()` first, compiles for Windows,
   uses temporary repositories and artifact directories, and builds its fake
   `docker` and fake `go` with `testfixture.FixtureBinary`, placed on the
   `PATH` of the injected environment. A test MUST NOT write a literal
   executable, call `t.Setenv`, start a real container or wait on the
   20-minute cap.
   - `internal/config/linux_lane_test.go`:
     `TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt`.
   - `internal/cli/deliver_linux_lane_test.go`:
     `TestLinuxLaneRunsTheChangedPackagesInAContainer`, which asserts every
     argument of API Contract 4, the archive on standard input and the
     selected packages for a changed `.go` file, a changed `testdata` file and
     a changed `go.mod`; `TestLinuxLaneSkipsWithANamedReason`, one subtest per
     reason of API Contract 3, each asserting the exact console line and that
     no `docker run` started; `TestLinuxLaneFailureFailsTheGate`, for exit 1;
     `TestLinuxLaneCancellationRemovesTheContainer`, which cancels while the
     fake `docker run` blocks and asserts `docker rm -f <name>`; and
     `TestGateWithoutTheLaneStartsNoDocker`.

   No existing test is renamed or removed.
8. MUST set `delivery.linux_lane: true` in `.roundfixrc.yml` under the existing
   `delivery:` block, with a one-line comment naming ADR-0288, and change
   nothing else in that file.
9. MUST document the lane, each naming `Linux Lane`:
   - in `docs/user-guide/configuration.md`, the key and its default;
   - in `docs/user-guide/commands/deliver.md`, the gate step, the skip
     reasons, the pass, fail and cap lines and the log path;
   - in `.agents/skills/roundfix/references/deliver.md`.

   No text inside any `### QA settlement` section may change.
10. MUST add **Linux Lane** to `CONTEXT.md` through domain-modeling, with
    ADR-0288's link.
11. MUST run these commands, because this Task edits user guides and the
    Roundfix references. No record or mirror is hand-edited.
    - `make skills-sync`;
    - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
      to raise and record the Roundfix skill version;
    - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
      to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the config key and its test.
- [ ] Add the lane and call it from the gate.
- [ ] Add the lane tests with fake `docker` and `go` binaries.
- [ ] Opt this repository in.
- [ ] Document the lane in the guides, the reference and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A gate with the key on runs one bounded `docker run` over the changed packages, and a failing run fails the gate with the log path.
- [ ] Each machine condition of API Contract 3 skips with its named reason and starts no container.
- [ ] With the key off, the gate starts no `docker` or `go` process.

## Consumers

- The Project Config `delivery` schema and the Delivery Queue gate result — consumers: `./internal/config`, `./internal/cli`

## Context

- instruction: `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`
- instruction: `docs/adr/0225-a-queue-item-runs-the-roundfix-binary-its-branch-builds.md`
- instruction: `docs/adr/0281-a-test-replaces-an-executable-it-rebuilds-instead-of-writing-into-it.md`
- instruction: `internal/testfixture/fixture.go`
- interface: `internal/config/config.go`
- creates: `internal/config/linux_lane_test.go`
- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_linux_lane.go`
- creates: `internal/cli/deliver_linux_lane_test.go`
- interface: `.roundfixrc.yml`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLinuxLaneRunsTheChangedPackagesInAContainer|TestLinuxLaneSkipsWithANamedReason|TestLinuxLaneFailureFailsTheGate|TestLinuxLaneCancellationRemovesTheContainer|TestGateWithoutTheLaneStartsNoDocker)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLinuxLaneRunsTheChangedPackagesInAContainer TestLinuxLaneSkipsWithANamedReason TestLinuxLaneFailureFailsTheGate TestLinuxLaneCancellationRemovesTheContainer TestGateWithoutTheLaneStartsNoDocker; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt ' || { printf 'missing pass: TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt\n' >&2; exit 1; }; grep -Eq '^  linux_lane: true$' .roundfixrc.yml || { printf 'missing delivery.linux_lane opt-in\n' >&2; exit 1; }; for f in docs/user-guide/configuration.md docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Linux Lane' || { printf 'missing Linux Lane in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Linux Lane**:' || { printf 'missing glossary term Linux Lane\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/config || exit 1; go test -count=1 -run 'Deliver|Gate|Lane|ItemBinary' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the lane tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the lane and config tests, the opt-in, the documentation, the glossary, the skill version and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Core Feature 7; Success Metric 1; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 10; Invariants 1 and 2; Current behavior; Testing Approach; Build Order 1
- ADR-0288; ADR-0225; ADR-0281; ADR-0285; ADR-0273; ADR-0271; ADR-0280; ADR-0189; ADR-0256; ADR-0259


## Result

Implemented Task 01's Linux Lane slice for Daemon verification. Task status
and checklist settlement remain Daemon-owned; no declared Verification command,
repository Verification, commit, push or pull request operation was run.

### Implementation

- Added strict boolean `delivery.linux_lane` parsing, default off, with Project
  Config overriding User Config and explicit project `false` disabling it.
  Non-boolean values retain the scoped `parse config` diagnostic shape.
- After repository Verification succeeds, Gate reads the setting from an
  immutable commit resolved from the already-fetched delivery default ref.
  This also covers resumed gates and prevents an item's or owner's unmerged
  opt-in from deciding the setting.
- Selected sorted, unique nearest Go package directories from the candidate's
  Git tree and merge-base diff; root `go.mod` and `go.sum` changes select
  `./...`. The archive streams directly from the supplied candidate head to
  Docker, including when the checkout is dirty or has later commits.
- Added ordered machine skips and exact console notices, a command/output log,
  the locally inspected image, the required fixed script and container
  arguments, owner-environment Docker/Go resolution, and a read-only module
  cache mount. Only the seven prescribed variables are explicitly forwarded
  to the container.
- Nonzero test exits fail the gate with the lane log path; exit 125 skips.
  Cancellation and the 20-minute cap remove the named container through a
  fresh, bounded cleanup context. The cap is not waited on by tests.
- Opted this repository in with only the authorized line/comment in
  `.roundfixrc.yml`; documented Linux Lane in the guides and Roundfix
  reference; added the glossary term with ADR-0288. Synced skill mirrors and
  generated Roundfix version 0.1.96 and the Behavior Surface Record. QA
  settlement text is byte-identical to HEAD.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Enabled gate runs one bounded container for changed packages; failure names its log | `TestLinuxLaneRunsTheChangedPackagesInAContainer` checks every argument, the literal script, process environment/directory, archive contents, candidate-head provenance, package selection for Go/testdata/module/sum changes, sorting and deduplication, and worktree preservation. `TestLinuxLaneFailureFailsTheGate` checks exits 1, 126 and 127 and the exact reason/log path. `TestLinuxLaneCancellationRemovesTheContainer` cancels only after the compiled fake run signals it is blocking, then checks `docker rm -f <name>` and `context.Canceled`. |
| Every machine condition skips with its named reason and starts no container | `TestLinuxLaneSkipsWithANamedReason` covers all eight conditions, earlier-reason precedence, exact console notices and matching lane-log notices. The first seven start no `docker run`; the exit-125 condition necessarily makes one refused run attempt and starts no real container. All tests use compiled fake tools. |
| Disabled gate starts no Docker or Go process | `TestGateWithoutTheLaneStartsNoDocker` checks both default-branch disablement despite the owner's opposite loaded value and repository Verification failure before the lane. It asserts no tool starts, no lane log and no console line. The enabled test also proves that the fetched default enables the lane despite the owner's loaded false and a dirty checkout's false. |

### Focused checks run

- `go test ./internal/config -run LinuxLane -count=1`: exit 0.
- `go test ./internal/cli -run '^(TestLinuxLane|TestGateWithoutTheLane)' -count=1`:
  exit 0.
- `go test ./internal/cli ./internal/config -run 'LinuxLane|GateWithoutTheLane' -count=1`:
  exit 0.
- `go test -race ./internal/cli -run 'LinuxLane|GateWithoutTheLane' -count=1`:
  exit 0.
- `go test -race ./internal/cli ./internal/config -run 'LinuxLane|GateWithoutTheLane' -count=1`:
  exit 0 after the final implementation edit.
- `go test ./internal/config ./internal/cli -run '^(TestDeliveryItemBinary|TestDeliveryLinuxLane|TestAPrerequisiteIsMet)' -count=1`:
  exit 0; existing declaration and prerequisite behavior exercised.
- `GOOS=windows GOARCH=amd64 go test -c ./internal/cli -o /tmp/roundfix-0279-cli.test.exe`
  and the equivalent `./internal/config` compile: exit 0.
- Required authoring commands `make skills-sync`, the owned-skill version
  record command, and the docscontract skill-coverage record command: exit 0.
  Records and mirrors were generated, not hand-edited.
- `make baseline-digests`: exit 0, `changed: false`; no derived Baseline file
  changed.
- `git diff --check`: exit 0. Canonical/mirror byte comparisons and extraction
  of QA settlement sections against HEAD: identical. Changed-path audit:
  only this Task's declared paths plus this Task file.

### Limits and follow-ups

- The new lane tests were absent at HEAD (confirmed with `git cat-file`),
  providing the pre-change absence evidence.
- All runtime evidence uses disposable Git repositories, artifact directories,
  and `testfixture.FixtureBinary` tools. No real container, image pull,
  provider call, or external network request was made. Loopback is used only
  as an event signal for cancellation; it has no service dependency.
- Daemon-declared Verification and repository Verification remain unrun in
  this turn. The real 20-minute cap is implemented but deliberately not
  elapsed in a test; cancellation exercises container cleanup without a
  timing allowance.
- No work from Task 02 or Task 03 was included.
