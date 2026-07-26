---
task: task_02
spec: 0050-doctor-skill-readiness-hardening
status: pending
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

- [ ] Add the pure shared hash component.
- [ ] Add the punctuation and Unicode compatibility oracle.
- [ ] Migrate Repository Skill Set folder hashing.
- [ ] Migrate Baseline external lock hashing.
- [ ] Remove duplicated comparators and stale explanatory comments.
- [ ] Exercise the real lock and restoration contracts.

## Acceptance Criteria

- [ ] The pinned corpus produces the exact skills CLI 1.5.19 digest.
- [ ] `_a` sorts before `-a`, and the Unicode cases match the pinned
      `en-US` oracle.
- [ ] Changing one path or byte changes the digest.
- [ ] Both production consumers call the shared implementation and contain no
      copied `strings.ToLower` path comparator.
- [ ] The real repository remains ready and existing Baseline restoration
      hashes remain compatible.
- [ ] `go.mod` and `go.sum` remain unchanged from Task 01.

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

