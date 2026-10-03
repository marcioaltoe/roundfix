---
task: task_04
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
status: completed
type: backend
complexity: high
---

# Task 04: A Verification command the shell cannot parse is malformed, and uncommitted sources are named

## Overview

On 2026-10-01 an amended Verification command of Spec 0205 held escaped
backticks inside its code span, so the Task Graph carried a truncated command
that the shell rejected. The authoring probe called it honest because it
failed before the work, and the complete task_03 failed twice. The same probe
runs the commands of uncommitted Task files without saying so. This Task makes
the one Verification prober ask the shell to parse each command first and
report a rejected command as malformed, adds a Spec Consistency finding for a
truncated span, and names each uncommitted source the probe ran.

## Requirements

1. MUST answer the Backlog Entry of 2026-09-30, "Authored Verification runs
   without a provenance check", and its addendum of 2026-10-01, as
   `_techspec.md` → The parse check and the probe report states; the
   provenance decision is ADR-0226's: name, never refuse.
2. MUST make `ProbeCommands` run `sh -n -c <command>` before `Verify` and,
   when the shell rejects it, record an unknown verdict whose cause wraps
   `ErrVerificationMalformed` and the parser's message, without running the
   command (invariant 7); when `sh` cannot be started it MUST run the command
   as today. The Daemon's pre-work refusal MUST therefore refuse a Task whose
   command is malformed.
3. MUST make `roundfix spec check --run-verification` print the verdict
   `malformed` with the parser's message for such a command, in text and in
   JSON, and exit `1` (API Contract 6, Surface Transcript 3).
4. MUST make the same command list, after `Verification tree: HEAD`, one
   `Uncommitted Verification source: <path> (<untracked|modified>)` line per
   probed `_tasks.md` or Task file that `git status` lists, and add the
   `uncommitted` array to the JSON verification report, empty when every
   source is committed; nothing is refused for it.
5. MUST add `CodeVerifyTruncated` and `TruncatedVerification` to
   `internal/speccheck/verification.go` and call it from
   `internal/speccheck/citations.go` for every Task that is not completed:
   one `SC-VERIFY-TRUNCATED` error per Verification bullet whose first code
   span ends in a backslash or whose text after that span holds an odd number
   of backticks (API Contract 1).
6. MUST NOT change `internal/speccheck/constraints.go` or
   `internal/speccheck/coherence.go`, and MUST keep the existing probe, spec
   check and Verification detector tests passing unedited.
7. MUST add the tests named in Verification to the new files
   `internal/daemon/verification_probe_malformed_test.go`,
   `internal/speccheck/verification_truncated_test.go` and
   `internal/cli/spec_check_provenance_test.go`; the malformed-command tests
   MUST prove the command never ran, for example by a marker file it would
   create staying absent.
8. MUST describe `SC-VERIFY-TRUNCATED`, the `malformed` verdict and the
   uncommitted-source line in `docs/user-guide/commands/spec.md` and the
   Roundfix Skill's `spec` reference; MUST raise the Roundfix Skill's version
   by one patch level above the version on this Task's starting tree in both
   front-matter fields, run `make skills-sync`, and re-record the version.

## Subtasks

- [ ] Parse each command before the prober runs it.
- [ ] Report `malformed` and the uncommitted sources in `spec check`.
- [ ] Add the truncated-span finding.
- [ ] Add the tests.
- [ ] Describe it in the guide and the Skill, and raise its version.

## Acceptance Criteria

- [ ] The prober reports the truncated command of Spec 0205's delivery as
      unknown with `ErrVerificationMalformed`, and a marker the command would
      create stays absent; a parsable failing command stays honest.
- [ ] The Daemon refuses a Task whose command is malformed before any Agent
      Session starts.
- [ ] `roundfix spec check --run-verification` prints `malformed` and exits
      `1` for it, lists each untracked or modified source, and lists none for
      a committed Spec; the JSON report carries the same facts.
- [ ] A pending Task whose Verification bullet ends its span in a backslash,
      or leaves a stray backtick after it, yields `SC-VERIFY-TRUNCATED`; a
      completed Task and a clean bullet do not.
- [ ] The guide and the Skill describe the finding, the verdict and the line,
      the mirrors equal their canonical files, and the raised version is
      recorded.

## Context

- creates: `internal/daemon/verification_probe_malformed_test.go`
- creates: `internal/speccheck/verification_truncated_test.go`
- creates: `internal/cli/spec_check_provenance_test.go`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/speccheck/verification.go`
- interface: `internal/speccheck/citations.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/daemon/task_engine.go`
- instruction: `internal/spec/task.go`
- instruction: `docs/specs/0219-a-delivery-that-survives-archive-requeue-and-review/references/2026-09-30-authored-verification-runs-without-a-provenance-check.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProbeReportsACommandTheShellCannotParse|TestProbeRunsNoMalformedCommand|TestPreWorkProbeRefusesAMalformedCommand)$" ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProbeReportsACommandTheShellCannotParse TestProbeRunsNoMalformedCommand TestPreWorkProbeRefusesAMalformedCommand; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestVerificationSpanEndingInBackslashIsTruncated|TestVerificationLineWithAStrayBacktickIsTruncated|TestCleanVerificationBulletIsNotTruncated)$" ./internal/speccheck 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestVerificationSpanEndingInBackslashIsTruncated TestVerificationLineWithAStrayBacktickIsTruncated TestCleanVerificationBulletIsNotTruncated; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestSpecCheckReportsAMalformedCommand|TestSpecCheckNamesAnUncommittedVerificationSource|TestSpecCheckNamesNoSourceWhenCommitted)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSpecCheckReportsAMalformedCommand TestSpecCheckNamesAnUncommittedVerificationSource TestSpecCheckNamesNoSourceWhenCommitted; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `for file in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "SC-VERIFY-TRUNCATED" || { printf 'missing finding in %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "malformed" || { printf 'missing verdict in %s\n' "$file" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "Uncommitted Verification source" || { printf 'missing source line in %s\n' "$file" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && out="$(go test -count=1 -v -run "^TestEveryOwnedSkillVersionIsRecorded$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass\n' >&2; exit 1; }` — expected: exit 0; before this Task neither the guide nor the reference names the finding or the line, so the command fails.

## References

- `_prd.md` → Goals; User Story 6; Core Features 7-9; Success Metric 6; Acceptance evidence
- `_techspec.md` → The parse check and the probe report; Data Models; Interfaces; API Contract 1; API Contract 6; Surface Transcript 3; Testing Approach 4; Build Order 4
- ADR-0226; ADR-0148; ADR-0182


## Result

Implemented the Task 04 slice for the 2026-09-30 Backlog Entry and its
2026-10-01 malformed-command addendum, following ADR-0226's decision to name
uncommitted sources without refusing them. Task status and declared
Verification remain Daemon-owned.

Acceptance evidence from focused implementation checks:

1. The shared prober runs `sh -n -c` before the verifier. A parser rejection
   produces an unknown verdict with `VerificationUnknownError` wrapping
   `ErrVerificationMalformed` and the parser message. The new
   `TestProbeReportsACommandTheShellCannotParse` exercises the truncated
   `grep` command and a parsable failing command; `TestProbeRunsNoMalformedCommand`
   proves a marker stays absent even when the malformed script places the
   marker command on an earlier line. Parser-startup fallback and cancellation
   also have focused tests.
2. `TestPreWorkProbeRefusesAMalformedCommand` exercises `TaskCycle` and
   observes the refusal reason, zero Agent calls, zero verifier calls, and
   zero commits. The existing Daemon refusal consumes the shared unknown
   verdict; no second parser or Task-engine policy was added.
3. `TestSpecCheckReportsAMalformedCommand` proves text and JSON report
   `malformed`, include the shell parser message, exit 1, and leave the marker
   absent. `TestSpecCheckNamesAnUncommittedVerificationSource` proves the
   modified graph and untracked Task are listed immediately after the HEAD
   line, JSON carries matching path/state objects, unrelated files are
   excluded, and the command still executes. `TestSpecCheckNamesNoSourceWhenCommitted`
   proves no source line and an empty JSON array for committed artifacts.
4. The authored-content detector reports one `SC-VERIFY-TRUNCATED` error per
   affected bullet with its actual path and line. The new truncated-span tests
   cover both malformed shapes, duplicate bullets, pending/in_progress/failed
   Tasks, completed-Task exclusion, clean spans, balanced trailing spans, and
   prose-only bullets. The Task-stage detector calls it after the completed
   Task exclusion. `constraints.go` and `coherence.go` are unchanged.
5. The command guide and canonical Skill reference describe the finding,
   malformed verdict, parser fallback, source line, and JSON fields. Both
   Skill front-matter versions rose from the starting 0.1.21 to 0.1.22.
   `make skills-sync` regenerated the mirrors; a byte comparison confirmed
   both mirror pairs match, and the version recorder added 0.1.22.

Checks performed:

- Initial focused regression compile failed on the absent
  `ErrVerificationMalformed`, establishing the missing implementation signal.
- `env GOCACHE="$PWD/.gocache" go test ./internal/daemon ./internal/speccheck ./internal/cli -run 'TestProbe|TestPreWorkProbe|TestVerificationSpan|TestVerificationLine|TestCleanVerification|TestSpecCheck|TestVerificationProbe' -count=1`
  exited 0 in all three packages after the final code edits. The first test
  run exposed a parser-message capitalization assumption and a fixture
  line-number offset; both test mistakes were corrected before this pass.
- `env GOCACHE="$PWD/.gocache" go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0. The initial invocation using the host cache was denied by the
  sandbox; the workspace cache resolved that access issue.
- `make baseline-digests` exited 0 and reported no derived changes.
- `rtk make verify-incremental` initially exited 2 because existing tests
  could not read the process table or bind a local HTTP test-server port in
  the sandbox. The permission-enabled rerun exited 0, including formatting,
  vet, the full existing Go suite, Skill checks, and build. Logs were captured
  at `/tmp/task04-incremental.log` and
  `/tmp/task04-incremental-unsandboxed.log`.
- Direct byte/content inspection confirmed the two Skill mirrors and all
  three documented surface terms. `git diff --check` reported no errors.

The starting worktree's only changed path was this Task's Daemon-written
status. Newly changed paths belong to this slice or its sanctioned Skill
regeneration. Existing tests, other Task files, and `_tasks.md` were not
edited. No commits, pushes, or pull requests were made. The commands in
`## Verification` were not run; the Daemon retains that gate and settlement.
No follow-up outside this slice was implemented.

## Carry-forward provenance

- Source Run: `run_20261003T213431Z_8d253be688342055`
- Source commit: `28eefd0482804df8ab59af3487162a023e13b718`
