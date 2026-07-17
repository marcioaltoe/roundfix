---
task: task_01
spec: 0036-doctor-skill-readiness
status: pending
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
6. MUST add exactly one `skills: ok|failed` Doctor line after the existing
   readiness lines, with derived owned/external counts and lexically sorted
   missing/outdated names.
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
11. MUST fix `sync-setups` so a valid current digest for `source.type: repo`
    takes precedence over same-path content in the external setup checkout,
    while non-repository skill refresh behavior remains unchanged.

## Subtasks

- [ ] Add repository readiness values and stable ownership classifications.
- [ ] Add exact owned-tree comparison against embedded skill files.
- [ ] Add strict required lock-entry validation and compatible external hashing.
- [ ] Inject the checker into Doctor and render deterministic success/failure.
- [ ] Add focused package tests for clean, missing, outdated, malformed, unsafe,
      symlink, mixed-ownership, and ordering cases.
- [ ] Add exact-output Doctor tests for success, each ownership group, mixed
      remediation, checker errors, exit behavior, and no mutation.
- [ ] Add setup synchronization regression tests for conflicting repo-owned
      external content and continuing external-source refresh.
- [ ] Update Doctor command help for the appended Repository Skill Set check.

## Acceptance Criteria

- [ ] A complete fixture matching the embedded owned bundle and required lock
      hashes reports `skills: ok` with derived 14/24/38 counts and exits zero
      when all other Doctor checks pass.
- [ ] Removing one owned or external skill reports that name under `missing`,
      prints the ownership-specific command, and exits one.
- [ ] Changing, adding, or removing a versioned file reports the affected skill
      under `outdated`; unrelated extra skills remain ignored.
- [ ] A compatibility fixture produces the exact same digest as the installed
      skills CLI algorithm, including path normalization and exclusions.
- [ ] Missing/malformed lock data, invalid hashes, unsafe names, unreadable
      artifacts, and symlinks fail deterministically without panic or traversal.
- [ ] Mixed failures produce one `skills:` line, sorted names, and both
      remediation commands exactly once in owned-then-external order.
- [ ] Doctor runs every existing check even when skill readiness fails and
      changes no repository, user-config, Run, or skill path.
- [ ] `sync-setups` preserves a Roundfix-owned digest when an external checkout
      contains conflicting content at the same path, and the canonical and
      embedded setup-context-driven suites pass together.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/skill-governance.md`
- interface: `skills/skills.go`
- interface: `skills/skills_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/cli.go`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_sync_setups.py`

## Verification

- `rtk go test ./skills -run 'Test(CheckRepository|SkillFolderHash)' -count=1`
  — expected: owned comparison, lock validation, hash compatibility, unsafe
  paths, stable ordering, and no-mutation cases pass.
- `rtk go test ./internal/cli -run 'TestRunDoctor' -count=1` — expected: exact
  `skills:` output, remediation, exit behavior, and existing Doctor lines pass.
- `rtk go test -race ./skills ./internal/cli -run 'Test(CheckRepository|RunDoctor)' -count=1`
  — expected: injected checks and filesystem reads are race-free.
- `rtk make setup-context-check` — expected: repo-owned digest precedence,
  external refresh, and all setup-context-driven flows pass.

## References

- `_prd.md` → Goals; User Stories 1–5; Core Features 1–6; Success Metrics.
- `_techspec.md` → Repository readiness contract; Owned-skill comparison;
  External lock and hash comparison; Doctor integration and output; Build Order
  1–3.
- `CONTEXT.md` → Doctor Command; Repository Skill Set; Roundfix Skill.
- `docs/agents/skill-governance.md` → owned versus external authority and
  synchronization boundaries.
- `docs/findings/2026-07-17-sync-setups-repo-owned-digest-drift.md` → reproduced
  source-precedence failure and required regression boundary.

## Result

Pending.
