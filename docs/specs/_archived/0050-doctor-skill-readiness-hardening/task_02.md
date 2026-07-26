---
task: task_02
spec: 0050-doctor-skill-readiness-hardening
status: completed
type: backend
complexity: medium
---

# Task 02: Centralize external skill hash compatibility

## Overview

Replace the duplicated lowercase-byte comparator with one pure internal hash
component that matches the installed skills CLI's `localeCompare` behavior.
Doctor-facing folder hashing and Baseline restoration must consume the same
implementation while preserving the existing digest shape.

## Requirements

1. MUST add one internal component that hashes path-plus-content pairs after
   stable American English Unicode collation.
2. MUST preserve slash-normalized paths, the absence of separators, lowercase
   SHA-256 output, and caller-owned input order.
3. MUST migrate `skills.SkillFolderHash` and Baseline external lock generation
   to the shared component.
4. MUST delete both copied lowercase-byte comparator implementations.
5. MUST add a pinned skills CLI 1.5.19 compatibility oracle covering
   punctuation, digits, case, Unicode, and nested paths.
6. MUST add a negative companion proving a path or content change alters the
   digest.
7. MUST preserve current real Repository Skill Set and Baseline restoration
   compatibility.

## Subtasks

- [x] Add the pure shared hash component.
- [x] Add the punctuation and Unicode compatibility oracle.
- [x] Migrate Repository Skill Set folder hashing.
- [x] Migrate Baseline external lock hashing.
- [x] Remove duplicated comparators and stale explanatory comments.
- [x] Exercise the real lock and restoration contracts.

## Acceptance Criteria

- [x] The pinned corpus produces the exact skills CLI 1.5.19 digest.
- [x] `_a` sorts before `-a`, and the Unicode cases match the pinned
      `en-US` oracle.
- [x] Changing one path or byte changes the digest.
- [x] Both production consumers call the shared implementation and contain no
      copied `strings.ToLower` path comparator.
- [x] The real repository remains ready and existing Baseline restoration
      hashes remain compatible.
- [x] `go.mod` and `go.sum` remain unchanged from Task 01.

## Context

- instruction: `docs/agents/go.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `skills/skills.go`
- interface: `skills/skills_test.go`
- interface: `internal/baseline/skills_restore.go`
- interface: `internal/baseline/skills_restore_test.go`
- interface: `internal/baseline/assets/lock-hash-compatibility-v1.json`

## Verification

- `rtk go test ./internal/skillhash ./skills ./internal/baseline -run 'Test(Sum|SkillFolderHash|CheckRepositoryMatchesRealRepository|SkillsRestore)' -count=1` — expected:
  the shared oracle, real Repository Skill Set, and Baseline restoration
  compatibility tests pass.
- `rtk go test -race ./internal/skillhash ./skills ./internal/baseline -run 'Test(Sum|SkillFolderHash|CheckRepositoryMatchesRealRepository|SkillsRestore)' -count=1` — expected:
  the shared collator and both consumers are race-free.

## References

- `_prd.md` → Goals 2–3; User Stories 2–3; Core Features 2–3; Success Metrics.
- `_techspec.md` → Shared external-skill hash; Integration Points; Testing
  Approach; Build Order 2.

## Result

Added one pure `internal/skillhash` component that copies caller-owned input,
slash-normalizes relative paths, applies stable American English collation,
and emits the existing separator-free lowercase SHA-256 digest. Repository
Skill Set folder hashing and Baseline external lock generation now adapt their
collected files to that component; both copied lowercase comparators and the
stale approximation comment were removed.

Verification:

- Pre-change `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test
  ./internal/skillhash -run '^TestSum' -count=1`: failed because `File` and
  `Sum` did not exist, establishing the missing shared component.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test
  ./internal/skillhash ./skills ./internal/baseline -run
  'Test(Sum|SkillFolderHash|CheckRepositoryMatchesRealRepository|SkillsRestore)'
  -count=1`: passed in all three packages.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache go test -race
  ./internal/skillhash ./skills ./internal/baseline -run
  'Test(Sum|SkillFolderHash|CheckRepositoryMatchesRealRepository|SkillsRestore)'
  -count=1`: passed in all three packages with race detection.
- The first sandboxed `make verify` attempt was blocked when Go tried to
  create the module-cache lock
  `/Users/marcio/go/pkg/mod/cache/download/golang.org/x/sync/@v/v0.22.0.lock`;
  an unfiltered package reproduction confirmed `operation not permitted`.
- `rtk env GOCACHE=/private/tmp/roundfix-task02-go-cache make verify`:
  passed outside the task sandbox with 2,399 Go tests, four focused skill
  contract tests, the Roundfix skill check, and the production build.
- `rtk git -c core.fsmonitor=false diff --check`: passed.

Acceptance evidence:

- `TestSumMatchesSkillsCLI1519` produced the pinned
  `2a46b6d704729eafc0148969028b9cc4030813059e1f7524def2f38b433011d4`
  digest for punctuation, digits, case, composed and decomposed Unicode, and
  nested paths while proving caller order stayed unchanged.
- `TestSumSortsUnderscoreBeforeHyphen` proved `_a` collates before `-a`;
  `TestSumChangesWhenPathOrContentChanges` proved either mutation changes the
  digest.
- `rtk rg -n "skillhash\.Sum" skills/skills.go
  internal/baseline/skills_restore.go` found the two production consumers.
  The companion search for the copied comparator and stale comment returned no
  matches.
- The focused suite passed `TestCheckRepositoryMatchesRealRepository`, the
  existing external compatibility fixture, and all selected
  `TestSkillsRestore` restoration and refusal contracts.
- Task 01 and final SHA-256 values for `go.mod` and `go.sum` are unchanged:
  `3764c738663bd617989809ef9ebdd9849bb07e4a20a4f043fd8d914d6d65d38a`
  and
  `277534b002103e8090d3a5c593c9bb3d642e34ef255198a01cc3a022ea8e1fd2`.

Follow-ups: none.
