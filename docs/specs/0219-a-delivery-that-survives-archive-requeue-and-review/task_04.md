---
task: task_04
spec: 0219-a-delivery-that-survives-archive-requeue-and-review
status: pending
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
