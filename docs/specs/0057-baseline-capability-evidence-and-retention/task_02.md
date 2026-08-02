---
task: task_02
spec: 0057-baseline-capability-evidence-and-retention
status: completed
type: backend
complexity: medium
---

# Task 02: Resolve symlinked executables without running them

## Overview

Executable capability discovery inspects each PATH candidate with a link-level
stat and then requires a regular file, so a symlink never qualifies. Every tool
installed through Homebrew or Docker Desktop — the norm on maintainer machines —
reports missing while being present and working. This Task resolves a bounded
link chain and judges the target, still without executing anything, and reports
the candidate it inspected instead of an empty result.

## Requirements

1. MUST resolve a symlinked PATH candidate through a bounded chain of links to
   a regular executable, and treat that as discovered evidence.
2. MUST NOT execute any candidate or target; discovery stays offline and
   side-effect free per ADR-0087.
3. MUST bound the chain so a cycle terminates, and report a cycle, a broken
   link, and a non-executable target as three distinct reasons.
4. MUST report the inspected candidate whenever one existed, and report the
   absence of any candidate distinctly from a candidate that failed.
5. MUST keep a directly-present regular executable resolving exactly as it does
   today.
6. MUST leave every other capability evidence kind unchanged.

## Subtasks

- [ ] Add bounded link resolution to the candidate probe.
- [ ] Distinguish cycle, broken link, and non-executable target.
- [ ] Report the inspected candidate and the absence case.
- [ ] Confirm no candidate is executed.

## Acceptance Criteria

- [ ] A one-hop symlink to a regular executable resolves as discovered.
- [ ] A multi-hop chain within the bound resolves as discovered.
- [ ] A link cycle terminates and reports the cycle reason.
- [ ] A broken link reports the broken-link reason, not absence.
- [ ] A link to a non-executable target reports the non-executable reason.
- [ ] A directly-present regular executable still resolves as discovered.
- [ ] No candidate is executed, proven by a probe target that would record its
      own invocation.
- [ ] The characterization corpus from task 01 shows only the intended
      capability status changes.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestExecutableCandidateResolution -count=1` —
  expected: exit 0; all seven probe cases assert their own reason.
- `go test ./internal/baseline -run TestExecutableCandidateNeverExecutes -count=1`
  — expected: exit 0; a target that would record its invocation records nothing.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 3; Core Features 4.
- `_techspec.md` → Implementation Design: Interfaces; Build Order 2.
- ADR-0087.

## Result

### Implementation

- Replaced link-level executable rejection with a PATH-ordered resolver that
  follows at most 64 symlink hops, detects visited-path cycles, resolves
  relative and absolute targets, and judges only the final regular file's
  executable bits.
- Preserved the first failed existing PATH candidate when no later candidate
  resolves, while returning `not-found` with no candidate when PATH contains no
  candidate.
- Projected successful and failed probes into capability evidence: the
  inspected candidate remains `sourcePath`; failures carry `broken-link`,
  `link-cycle`, or `not-executable`; true absence carries `not-found` without a
  path.
- Added filesystem-backed coverage for the seven resolver cases, evidence
  projection, and a symlinked script whose invocation would create a marker.
  The implementation uses only `Lstat` and `Readlink`; it has no process launch
  path.

### Focused checks

- Pre-change signal:
  `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/baseline -run '^TestExecutableCandidateResolution$/one-hop_symlink$'`
  failed to compile because `resolveExecutableCandidate` did not exist.
- Resolver cases:
  `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/baseline -run '^TestExecutableCandidateResolution$/(direct_regular_executable|one-hop_symlink|multi-hop_symlink|link_cycle|broken_link|non-executable_target|absent_candidate)$'`
  reported 8 passing tests.
- Evidence and no-execution checks:
  `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/baseline -run '^(TestExecutableCandidateNeverExecutes|TestExecutableEvidenceDistinguishesFailureFromAbsence|TestCapabilityAuditNoExecution)$'`
  reported 5 passing tests.
- Fresh combined implementation check:
  `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/baseline -run '^(TestExecutableCandidateResolution|TestExecutableCandidateNeverExecutes|TestExecutableEvidenceDistinguishesFailureFromAbsence|TestCapabilityAuditNoExecution)$'`
  reported 13 passing tests.
- Existing evidence/ranking checks:
  `GOCACHE=/private/tmp/roundfix-task02-gocache rtk go test ./internal/baseline -run '^(TestProfileAlignmentCapabilityEvidenceRanking|TestProfileAlignmentDiscoversDeclaredRepositoryFormatter)$'`
  reported 5 passing tests.
- `rtk git diff --check` reported no whitespace errors.
- `rtk git -c core.fsmonitor=false status --porcelain` and
  `rtk git diff --name-only` listed only this Task file,
  `internal/baseline/profile_alignment.go`, and
  `internal/baseline/profile_alignment_test.go`.
- A bounded source search found no `exec.Command`, `exec.CommandContext`,
  `os.StartProcess`, or `syscall.Exec` reference in the changed implementation
  or tests.

### Acceptance criterion evidence

1. `TestExecutableCandidateResolution/one-hop_symlink` resolves the inspected
   link to its executable target with one hop and no rejection reason.
2. `TestExecutableCandidateResolution/multi-hop_symlink` resolves two relative
   links to the executable target with two hops.
3. `TestExecutableCandidateResolution/link_cycle` terminates with
   `link-cycle`, the inspected candidate, and no resolved target.
4. `TestExecutableCandidateResolution/broken_link` returns `broken-link`, and
   `TestExecutableEvidenceDistinguishesFailureFromAbsence/failed_candidate`
   retains that candidate in invalid evidence.
5. `TestExecutableCandidateResolution/non-executable_target` returns
   `not-executable` for a link whose regular target has no executable bit.
6. `TestExecutableCandidateResolution/direct_regular_executable` preserves the
   direct candidate as both inspected and resolved with zero hops.
7. `TestExecutableCandidateNeverExecutes` discovers a symlinked script through
   the real evidence collector and confirms its invocation marker is absent;
   `TestCapabilityAuditNoExecution` preserves the broader command-free audit.
8. No characterization fixture or golden file changed. The declared
   `TestBaselinePlanCharacterization` command was not run because the Daemon
   owns Task Verification; its corpus comparison remains pending that gate.
9. The final status and diff-name postflight listed only paths under
   `internal/baseline/` and this Task file.

### Daemon-owned Verification

No command from `## Verification` was run in this Agent turn. The Daemon owns
those commands and the terminal Task status.
