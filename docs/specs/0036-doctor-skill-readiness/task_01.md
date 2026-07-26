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
- [ ] Update Doctor command help for the appended Repository Skill Set check.

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
- [ ] Doctor help continues to name both Agent Selection Profiles and the
      appended Repository Skill Set readiness check.

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
- `rtk go test ./internal/cli -run 'Test(RunDoctor|ProfilesDocumentationContractMatchesPublicGuidance)' -count=1`
  — expected: exact `skills:` output, remediation, exit behavior, existing
  Doctor lines, and public help terminology pass.
- `rtk go test -race ./skills ./internal/cli -run 'Test(CheckRepository|RunDoctor)' -count=1`
  — expected: injected checks and filesystem reads are race-free.
- `rtk go test ./skills -count=1` — expected: the complete skills suite,
  including the real repository baseline contract, passes.

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

## Rework Trigger

Fresh real-repository evidence invalidated the completed compatibility claim:
after `bunx skills update -p -y` normalized 25 external lock entries, Doctor
reported 21 matching skills as outdated.

The installed `skills` CLI versions 1.5.19 and 1.5.20 sort relative paths with
JavaScript `localeCompare`, while `SkillFolderHash` currently uses Go byte
ordering. The pinned fixture encodes the Go ordering (`SKILL.md` before
`references/guide.md`) and therefore proves the implementation against itself
instead of against the CLI. Repair the production ordering and fixture, add a
regression that distinguishes these orderings, and prove the complete skills
suite against the refreshed real repository.

## Rework Trigger 2

The first Task 03 full gate exposed
`TestProfilesDocumentationContractMatchesPublicGuidance`: Doctor help no longer
names `Agent Selection Profiles`, even though Spec 0041 remains the prerequisite
readiness authority. Restore that existing public term while keeping the new
Repository Skill Set wording additive, then prove the focused documentation
contract before Task 03 runs.

## Result

Repaired the external skill hash compatibility contract. `SkillFolderHash` and
the Baseline external-lock adapter now reproduce the installed skills CLI
`localeCompare` ordering for repository paths, while Baseline setup snapshots
retain their separate byte-ordered content-digest contract.

The compatibility fixture now pins the CLI-derived `501e315e...` digest and
explicitly fails if its file order no longer distinguishes JavaScript locale
ordering from Go byte ordering. A read-only real-repository test proves all 25
refreshed external skill trees against `skills-lock.json`; catalog compatibility
fixtures carry the resulting lock-adapter artifact digest.

### Verification

- `rtk env GOCACHE=/tmp/roundfix-task01-final-skills-cache go test ./skills -run 'Test(CheckRepository|SkillFolderHash)' -count=1`
  — passed (`ok roundfix/skills`).
- `rtk env GOCACHE=/tmp/roundfix-task01-final-cli-cache go test ./internal/cli -run 'TestRunDoctor' -count=1`
  — passed (`ok roundfix/internal/cli`).
- `rtk env GOCACHE=/tmp/roundfix-task01-final-race-cache go test -race ./skills ./internal/cli -run 'Test(CheckRepository|RunDoctor)' -count=1`
  — passed for both packages with the race detector.
- `rtk env GOCACHE=/tmp/roundfix-task01-final-full-cache go test ./skills -count=1`
  — passed (`ok roundfix/skills`), including the real repository baseline.
- `rtk env GOCACHE=/tmp/roundfix-task01-rework-baseline-cache go test ./internal/baseline -run 'Test(EmbeddedCatalog|AssetsSyncCompatibilityMatchesMaintainedPythonContract|SkillsRestore|CatalogCompatibility)' -count=1`
  — passed (`ok roundfix/internal/baseline`).
- `rtk env GOCACHE=/tmp/roundfix-task01-rework-verify-cache make verify`
  — blocked after `2393` passing tests by the pre-existing, out-of-scope
  `TestProfilesDocumentationContractMatchesPublicGuidance` failure:
  Doctor help is missing `Agent Selection Profiles`.

The sandbox does not permit Go's default user cache, so verification used
task-local temporary `GOCACHE` directories. The Daemon remains responsible for
running the declared Verification commands verbatim.

### Acceptance evidence

- AC 1: the complete fixture and exact Doctor tests prove derived `14` owned,
  `25` external, and `39` total counts with exit zero.
- AC 2: focused repository and Doctor tests prove owned and external removals
  are sorted under `missing`, select the ownership-specific command, and exit
  one.
- AC 3: focused tests prove changed, added, and removed versioned artifacts are
  `outdated`, while unrelated installed and lock entries remain ignored.
- AC 4: the pinned fixture hashes `references/guide.md` before `SKILL.md` to
  `501e315e486dc59cbaa999085edd4312d35bbc690947dfb59af2e86722466aa9`;
  the regression rejects Go byte ordering, and the real repository test matches
  every refreshed external lock hash.
- AC 5: focused tests prove malformed locks, invalid hashes, unsafe names,
  unreadable artifacts, and symlinks fail without panic or traversal.
- AC 6: exact Doctor tests prove one sorted `skills:` line and one copy of each
  remediation command in owned-then-external order.
- AC 7: race, no-mutation, and exact Doctor tests prove all existing checks
  still run and the checker performs no repository, config, Run, lock, or skill
  mutation.

### Follow-ups

Restore the pre-existing `Agent Selection Profiles` wording required by
`TestProfilesDocumentationContractMatchesPublicGuidance` in its owning CLI
documentation slice; no `internal/cli` file changed in this rework.
