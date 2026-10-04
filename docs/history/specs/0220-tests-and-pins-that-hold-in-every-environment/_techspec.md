---
spec: 0220-tests-and-pins-that-hold-in-every-environment
prd: _prd.md
created: 2026-10-03
---

# Tests and pins that hold in every environment — Technical Spec

## Executive Summary

Both defects are checks that read a proxy for the fact they assert. The death
tests read the answer to `kill(-pgid, 0)` as the life of a process group, and
on macOS a group whose members have all exited but are not yet reaped answers
`EPERM`. The skill contract tests read one hand-pinned digest as the state of
the upstream skills, and the supported restore moves that state. Each fix makes
the check read the record that changes with the fact: the process table's
group membership and the fixture's start time, and the lock and Setup Snapshot
that the restore itself writes. The primary trade-off is strictness: this
repository's `make verify` now fails when its skills trail the embedded Setup
Snapshot, so a later change that moves the snapshot must restore the skills in
the same change (ADR-0224).

## Project Constraints

- Identifier strategy: not applicable — no identifier is created or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local fixtures only; the restore's
  one source fetch runs in the Agent turn of Build Order 2 and never in
  Verification. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0213, ADR-0221 and ADR-0224 bind
  the two changes; ADR-0191, ADR-0125, ADR-0126 and ADR-0028 hold unchanged;
  ADR-0179 bounds the Governed Paths; ADR-0093, ADR-0117, ADR-0168, ADR-0176
  and ADR-0183 check consistency; the gate is bound by ADR-0080, ADR-0088,
  ADR-0091, ADR-0104, ADR-0155, ADR-0156 and ADR-0167; ADR-0189, ADR-0182,
  ADR-0184, ADR-0204, ADR-0206, ADR-0030, ADR-0096, ADR-0097, ADR-0098, ADR-0194,
  ADR-0195, ADR-0210 and ADR-0219 do not apply,
  as the PRD records. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md): the unattended program of 2026-09-30,
  extended 2026-10-01, the maintainer's "Continue" to cycle G/H/I on 2026-10-03,
  "Autorizar os dois" for the Baseline source and the guides, and "considere
  autorizado a ajustar todas as skills se necessário" for skills. No Makefile,
  lint, formatter, test-runner or CI file and no dependency changes. Bounded
  files: `.agents/skills/bubbletea/SKILL.md`,
  `.agents/skills/bubbletea/references/components.md`,
  `.agents/skills/bubbletea/references/emoji-width-fix.md`,
  `.agents/skills/bubbletea/references/golden-rules.md`,
  `.agents/skills/bubbletea/references/troubleshooting.md`,
  `.agents/skills/domain-modeling/CONTEXT-FORMAT.md`,
  `.agents/skills/domain-modeling/GLOSSARY-FORMAT.md`,
  `.agents/skills/domain-modeling/SKILL.md`,
  `.agents/skills/golang-concurrency/SKILL.md`,
  `.agents/skills/golang-concurrency/references/channels-and-select.md`,
  `.agents/skills/golang-concurrency/references/pipelines.md`,
  `.agents/skills/golang-concurrency/references/sync-primitives.md`,
  `.agents/skills/golang-context/SKILL.md`,
  `.agents/skills/golang-context/references/cancellation.md`,
  `.agents/skills/golang-context/references/http-services.md`,
  `.agents/skills/golang-error-handling/SKILL.md`,
  `.agents/skills/golang-error-handling/references/error-creation.md`,
  `.agents/skills/golang-error-handling/references/error-handling.md`,
  `.agents/skills/golang-error-handling/references/error-wrapping.md`,
  `.agents/skills/golang-lint/SKILL.md`,
  `.agents/skills/golang-lint/references/linter-reference.md`,
  `.agents/skills/golang-testing/SKILL.md`,
  `.agents/skills/golang-testing/references/benchmarks.md`,
  `.agents/skills/golang-testing/references/coverage.md`,
  `.agents/skills/golang-testing/references/examples.md`,
  `.agents/skills/golang-testing/references/integration-testing.md`,
  `.agents/skills/golang-testing/references/mocking.md`,
  `.agents/skills/no-workarounds/SKILL.md`,
  `.agents/skills/systematic-debugging/CREATION-LOG.md`,
  `.agents/skills/systematic-debugging/SKILL.md`,
  `.agents/skills/systematic-debugging/condition-based-waiting-example.ts`,
  `.agents/skills/systematic-debugging/condition-based-waiting.md`,
  `.agents/skills/systematic-debugging/defense-in-depth.md`,
  `.agents/skills/systematic-debugging/find-polluter.sh`,
  `.agents/skills/systematic-debugging/root-cause-tracing.md`,
  `.agents/skills/systematic-debugging/test-pressure-1.md`,
  `.agents/skills/systematic-debugging/test-pressure-2.md`,
  `.agents/skills/systematic-debugging/test-pressure-3.md`,
  `.agents/skills/testing-boss/SKILL.md`,
  `.agents/skills/testing-boss/references/ai-writes-tests.md`,
  `.agents/skills/tui-design/SKILL.md`,
  `.agents/skills/tui-design/references/app-patterns.md`,
  `.agents/skills/tui-design/references/visual-catalog.md`,
  `skills/baseline_skill_contract_test.go`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

No production file changes and no package is added.

| Component | Where | Change |
| --- | --- | --- |
| Fixture group reading | new `internal/cli/detach_fixture_group_darwin_test.go`, `detach_fixture_group_linux_test.go`, `detach_fixture_group_other_test.go` | Lists the live members of a process group from the process table |
| Death-test probe and cleanup | `internal/cli/implement_detach_teardown_test.go` | Ends on no live member and a changed or absent start time; no skip |
| Zombie-group characterization | new `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` in the teardown file | Holds a zombie group leader and proves both the platform answer and the reading |
| Restored skills | the eleven trees under `.agents/skills/` and `skills-lock.json` | Written by the supported restore |
| Lock comparison | `skills/baseline_skill_contract_test.go` | Each lock entry against its installed tree; constant removed |
| Snapshot comparison | new `internal/baseline/repository_skills_snapshot_test.go` | `TrailingSetupSkills` over this repository is empty |

### Why an `EPERM` arrives

XNU's `kill()` sends a negative PID to `killpg1`. `killpg1` finds the group,
iterates it with a filter that drops `SZOMB` members, counts members the sender
may signal, and returns `EPERM` when it counted none. A process group exists
until its last member is reaped. So between the moment the last fixture
process exits and the moment launchd reaps it, the probe answers `EPERM`, and
`kill(pid, 0)` on the zombie leader answers success. Linux answers success for
the same group until the reap, so it only waits longer. The probe meets the
window more often when the machine is loaded and the reaper is slow, which is
the Daemon's Verification during a delivery. A live member that the sender may
not signal, for example under a sandbox, also answers `EPERM`, so `EPERM`
alone cannot mean ended.

## Implementation Design

### Interfaces

```go
// internal/cli/detach_fixture_group_<os>_test.go
// detachFixtureGroupLiveMembers lists the PIDs whose process group is pgid and
// that have not exited. An exited, unreaped process is not listed. An error is
// a failed reading, never an empty group.
func detachFixtureGroupLiveMembers(pgid int) ([]int, error)
```

Darwin reads `kern.proc.pgrp` through `unix.SysctlKinfoProcSlice` and drops
entries whose `Proc.P_stat` is `SZOMB` (5 in `<sys/proc.h>`, named as a test
constant). Linux scans `/proc/<pid>/stat`, reads the state and process-group
fields after the last `)`, drops state `Z`, and skips a PID that vanished
during the scan. Other Unix systems keep the signal probe and report a live
member only on success.

```go
// internal/cli/implement_detach_teardown_test.go
// detachFixtureEnded reports whether the recorded fixture has ended: its group
// has no live member, or its PID now carries a start identity other than the
// one recorded while it was alive. The fixture leads its group, so a live
// fixture is always a live member.
func detachFixtureEnded(t *testing.T, pid int, identity string) (bool, string)
```

The probe records `store.OwnerProcessIdentity(ctx, pid)` once it has seen the
fixture alive, before it kills the inner test binary.

### Data Models

None. No schema, record or file format changes; `skills-lock.json` keeps its
shape and only the restore writes its values.

### API Contracts

None. No command, flag, output or exit code changes; the Spec changes tests and
vendored skill content.

### Surface Transcripts

None. No command surface changes.

## Coverage Map

- Goal 1 → Fixture group reading; Death-test probe and cleanup
- Goal 2 → Zombie-group characterization; Fixture group reading
- Goal 3 → Restored skills; Snapshot comparison
- Goal 4 → Lock comparison; Snapshot comparison
- Core Feature 1 → Fixture group reading; Death-test probe and cleanup;
  Zombie-group characterization
- Core Feature 2 → Restored skills; Lock comparison; Snapshot comparison
- Success Metric 1 → Zombie-group characterization; Death-test probe and
  cleanup
- Success Metric 2 → Lock comparison; Snapshot comparison
- Success Metric 3 → Restored skills; Snapshot comparison

## Integration Points

- The process table, read through `golang.org/x/sys/unix` on Darwin, which
  `go.mod` already requires, and through procfs on Linux.
- The upstream skills repository at `b3c45a4`, reached once by the restore in
  the Agent turn. A local checkout passed with `--source-dir` is an equivalent
  source when it is at that commit.

## Testing Approach

1. `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` starts `cat` as
   the leader of its own process group, with a pipe as its standard input, so
   the child lives until the test closes that pipe. It asserts that the reading lists the child while it lives,
   then closes that input so the child exits, and does not wait for it. Once
   the reading lists no live member while
   `kill(pid, 0)` still succeeds, it asserts the platform answer to
   `kill(-pgid, 0)` (`EPERM` on Darwin, success on Linux), asserts that
   `detachFixtureEnded` reports ended, then reaps the child and asserts
   `ESRCH`. It fails on the tree before this Spec, because the test does not
   exist, and it would fail against a reading that treats a zombie as live.
2. `TestImplementDetachChildEndsWhenItsTestBinaryDies` and
   `TestDetachSurvivorEndsWhenItsTestBinaryDies` keep their structure and their
   `SIGKILL` assertions, use `detachFixtureEnded`, and never call `t.Skip`.
3. A new lock-mismatch test copies one upstream-managed tree and its lock entry
   into a temporary repository, changes one byte, and asserts that the lock
   comparison names that skill.
4. `TestThisRepositoryHoldsItsRequiredSkillsAtTheSetupSnapshot` reads the
   Baseline Profile from `docs/agents/setup-context.json` and the skill names
   from `skills-lock.json`, and asserts that `TrailingSetupSkills` returns no
   name. Its failure message names each trailing skill and the restore command.

## Build Order

1. The fixture group reading, the zombie-group characterization and the death
   tests without a skip. Files: `internal/cli/implement_detach_teardown_test.go`
   and the three new `internal/cli/detach_fixture_group_*_test.go` files.
2. The restore and the record-reading skill tests, in one Task, because the
   restore fails the old pin and the new snapshot test fails without the
   restore. The restore is run with `roundfix baseline update --repo . --yes`
   built from this tree, and changes these 43 files under `.agents/skills/`:
   - `bubbletea`: `SKILL.md`, `references/components.md`,
     `references/emoji-width-fix.md`, `references/golden-rules.md`,
     `references/troubleshooting.md`;
   - `domain-modeling`: `SKILL.md`, the new `GLOSSARY-FORMAT.md`, and the
     removed `CONTEXT-FORMAT.md`;
   - `golang-concurrency`: `SKILL.md`, `references/channels-and-select.md`,
     `references/pipelines.md`, `references/sync-primitives.md`;
   - `golang-context`: `SKILL.md`, `references/cancellation.md`,
     `references/http-services.md`;
   - `golang-error-handling`: `SKILL.md`, `references/error-creation.md`,
     `references/error-handling.md`, `references/error-wrapping.md`;
   - `golang-lint`: `SKILL.md`, `references/linter-reference.md`;
   - `golang-testing`: `SKILL.md`, the new `references/benchmarks.md`,
     `references/coverage.md` and `references/examples.md`, and
     `references/integration-testing.md`, `references/mocking.md`;
   - `no-workarounds`: `SKILL.md`;
   - `systematic-debugging`: `SKILL.md`, `CREATION-LOG.md`,
     `condition-based-waiting-example.ts`, `condition-based-waiting.md`,
     `defense-in-depth.md`, `find-polluter.sh`, `root-cause-tracing.md`,
     `test-pressure-1.md`, `test-pressure-2.md`, `test-pressure-3.md`;
   - `testing-boss`: `SKILL.md`, `references/ai-writes-tests.md`;
   - `tui-design`: `SKILL.md`, `references/app-patterns.md`,
     `references/visual-catalog.md`.

   It also changes `skills-lock.json`,
   `skills/baseline_skill_contract_test.go` and the new
   `internal/baseline/repository_skills_snapshot_test.go`. Independent of
   step 1.
3. The QA gate. (depends on: 1, 2)

## Risks & Considerations

- **A restore that writes more than measured.** The authoring clone measured
  exactly the 43 files above. Upstream `b3c45a4` is fixed, so a second restore
  writes the same bytes; a path outside the grant fails the Task's commit, and
  the record would need a widened grant first.
- **No network in the Agent turn.** The restore then needs a local source at
  `b3c45a4` through `--source-dir`. The Task names both sources; neither
  appears in Verification.
- **A reading the sandbox hides.** A sandbox that denies process information
  could return an empty group. The Darwin reading fails on a `sysctl` error
  rather than treating it as empty, and the characterization test fails when
  the reading does not list a member it knows is alive, so a hidden group
  stops the run instead of passing it.
- **PID reuse.** The kernel does not hand out a PID that is still a process
  group ID, so a group with the recorded ID and a leader of another start
  identity is a new group, and the old one has ended.

## Decisions

- Read group membership and start time; do not interpret `EPERM`. Treating
  every `EPERM` as ended was rejected because a live member that may not be
  signalled answers the same.
- Platform-split test files over `ps`, because procps may be absent from a
  minimal Linux image and `ps` adds a process per poll.
- Restore and replace the pin in one Task, and hold this repository to its
  snapshot and lock. See ADR-0224.
- Leave the release-plan Backlog Entry open, as the PRD's Non-Goals state.
