---
task: task_01
spec: 0036-doctor-skill-readiness
status: completed
type: backend
complexity: high
---

# Task 01: Diagnose Repository Skill Set readiness

## Overview

Extend the Doctor Command with one blocking, deterministic Repository Skill Set
check. The slice compares installed Roundfix-owned skills with the running
binary, compares required external skills with `skills-lock.json`, reports
stable remediation, and proves the complete behavior without network access or
filesystem mutation.

## Requirements

1. MUST add structured repository-readiness types and a read-only checker at
   the existing `skills` ownership boundary.
2. MUST require all names from the binary's owned `Names()` and external
   `Recommended()` sets under `<git-root>/.agents/skills`.
3. MUST classify an owned skill as missing when its directory is absent and as
   outdated when required files are absent or changed, unexpected regular files
   exist, or symlinks replace versioned files.
4. MUST parse `skills-lock.json` locally, validate version 1 and every required
   external `computedHash`, and ignore unrelated installed directories and
   unrelated lock entries.
5. MUST reproduce the external skills CLI hash contract exactly: slash-normalized
   relative paths sorted lexically, SHA-256 over each path followed immediately
   by its file bytes, excluding `.git` and `node_modules` directories.
6. MUST add exactly one `skills: ok|failed` Doctor line after the Agent
   Selection Profile readiness line delivered by Spec 0041, with derived
   owned/external counts and lexically sorted missing/outdated names.
7. MUST make missing, outdated, malformed-lock, and unreadable-artifact results
   fail Doctor with `exitRunFailed`, while preserving every existing usage,
   stdout, stderr, and other readiness-check contract.
8. MUST print `roundfix skills install --target project` for owned failures and
   `bunx skills experimental_install && bunx skills update -p -y` for external
   failures, once each and in that order when both apply.
9. MUST perform no command execution, network access, downloads, installs,
   config writes, lock updates, or skill-tree mutations.
10. MUST wrap filesystem and decoding errors with the failed operation and path,
    and MUST reject unsafe skill names without reading outside the repository.

## Subtasks

- [x] Add repository readiness values and stable ownership classifications.
- [x] Add exact owned-tree comparison against embedded skill files.
- [x] Add strict required lock-entry validation and compatible external hashing.
- [x] Inject the checker into Doctor and render deterministic success/failure.
- [x] Add focused package tests for clean, missing, outdated, malformed, unsafe,
      symlink, mixed-ownership, and ordering cases.
- [x] Add exact-output Doctor tests for success, each ownership group, mixed
      remediation, checker errors, exit behavior, and no mutation.
- [x] Update Doctor command help for the appended Repository Skill Set check.

## Acceptance Criteria

- [x] A complete fixture matching the embedded owned bundle and required lock
      hashes reports `skills: ok` with the current derived 14/25/39 counts and
      exits zero when all other Doctor checks pass.
- [x] Removing one owned or external skill reports that name under `missing`,
      prints the ownership-specific command, and exits one.
- [x] Changing, adding, or removing a versioned file reports the affected skill
      under `outdated`; unrelated extra skills remain ignored.
- [x] A compatibility fixture produces the exact same digest as the installed
      skills CLI algorithm, including path normalization and exclusions.
- [x] Missing/malformed lock data, invalid hashes, unsafe names, unreadable
      artifacts, and symlinks fail deterministically without panic or traversal.
- [x] Mixed failures produce one `skills:` line, sorted names, and both
      remediation commands exactly once in owned-then-external order.
- [x] Doctor runs every existing check even when skill readiness fails and
      changes no repository, user-config, Run, or skill path.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/skill-dispatch.md`
- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- instruction: `docs/specs/_archived/0041-agent-selection-runtime-readiness/_techspec.md`
- interface: `skills/skills.go`
- interface: `skills/skills_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/cli.go`

## Verification

- `rtk go test ./skills -run 'Test(CheckRepository|SkillFolderHash)' -count=1`
  — expected: owned comparison, lock validation, hash compatibility, unsafe
  paths, stable ordering, and no-mutation cases pass.
- `rtk go test ./internal/cli -run 'TestRunDoctor' -count=1` — expected: exact
  `skills:` output, remediation, exit behavior, and existing Doctor lines pass.
- `rtk go test -race ./skills ./internal/cli -run 'Test(CheckRepository|RunDoctor)' -count=1`
  — expected: injected checks and filesystem reads are race-free.

## References

- `_prd.md` → Goals; User Stories 1–5; Core Features 1–5; Success Metrics.
- `_techspec.md` → Repository readiness contract; Owned-skill comparison;
  External lock and hash comparison; Doctor integration and output; Build
  Order 1–2.
- `docs/specs/_archived/0041-agent-selection-runtime-readiness/_techspec.md` → prerequisite
  profile-aware Doctor coordinator and output order.
- `CONTEXT.md` → Doctor Command; Repository Skill Set; Roundfix Skill.
- `docs/agents/skill-dispatch.md` → owned versus upstream-managed authority and
  synchronization boundaries.

## Result

Implemented one offline, read-only Repository Skill Set readiness check at the
`skills` ownership boundary and appended its deterministic `skills:` result
after Agent Selection Profile readiness in Doctor.

### Verification

- `rtk env GOCACHE=/tmp/roundfix-task01-go-cache go test ./skills -run 'Test(CheckRepository|SkillFolderHash)' -count=1`
  — passed (`ok roundfix/skills`).
- `rtk env GOCACHE=/tmp/roundfix-task01-go-cache go test ./internal/cli -run 'TestRunDoctor' -count=1`
  — passed (`ok roundfix/internal/cli`).
- `rtk env GOCACHE=/tmp/roundfix-task01-go-cache go test -race ./skills ./internal/cli -run 'Test(CheckRepository|RunDoctor)' -count=1`
  — passed for both packages with the race detector.
- `rtk env GOCACHE=/tmp/roundfix-task01-go-cache go test ./skills ./internal/cli -count=1`
  — passed for both complete package suites.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

The sandbox did not permit Go's default user cache, so focused checks used the
equivalent temporary `GOCACHE`; the Daemon remains responsible for running the
task's declared Verification commands verbatim.

### Acceptance evidence

- The complete embedded-owned plus required-external fixture returns derived
  `14` owned, `25` external, and Doctor renders `39 required` with exit zero.
- Owned and external removals are classified under sorted `missing` names;
  content changes, added files, and removed files are classified under
  `outdated`; unrelated installed skills and lock entries are ignored.
- The pinned compatibility fixture, nested slash-normalized paths, lexical
  ordering, and `.git`/`node_modules` exclusions pass through
  `SkillFolderHash`.
- Missing and malformed locks, wrong versions, absent entries, invalid hashes,
  unsafe names, unreadable artifact shapes, owned symlinks, and external
  symlinks have deterministic error or readiness outcomes with no traversal.
- Exact Doctor output tests prove success, owned-only, external-only, mixed,
  and checker-error lines; mixed remediation appears exactly once in
  owned-then-external order and all pre-existing checks still execute.
- A before/after fixture snapshot proves `CheckRepository` performs no
  mutation. Production code executes no commands, performs no network access,
  and exposes no write path.

### Follow-ups

None. Documentation and Roundfix Skill synchronization remain assigned to
their separate Tasks in this Spec.
