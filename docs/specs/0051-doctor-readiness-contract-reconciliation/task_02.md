---
task: task_02
spec: 0051-doctor-readiness-contract-reconciliation
status: pending
type: backend
complexity: medium
---

# Task 02: Make external skill hash ordering total

## Overview

Give the shared external-skill digest one total path order even when American
English collation treats distinct Unicode strings as equal. The slice remains
pure and independently verifiable through permutation tests plus the existing
Baseline consumer contract.

## Requirements

1. MUST retain slash normalization and American English collation as the
   primary path order.
2. MUST break a zero collation comparison between distinct normalized paths
   with ordinary Go string ordering.
3. MUST produce the same digest for every permutation of the same file set,
   including precomposed and decomposed Unicode path pairs.
4. MUST retain `_a` before `-a`, the current path-plus-content digest shape,
   and caller-owned slice order.
5. MUST keep `internal/skillhash` as the only comparator implementation used by
   Repository Skill Set and Baseline consumers.
6. MUST NOT change a lock file, Baseline asset, parity artifact, recommendation
   list, or upstream-managed skill.

## Subtasks

- [ ] Add a failing permutation case for collation-equal distinct paths.
- [ ] Define the normalized raw-path tie-breaker after collator equality.
- [ ] Preserve punctuation, case, path-depth, and content-change regressions.
- [ ] Exercise the existing Baseline restoration consumer against the shared
      implementation.

## Acceptance Criteria

- [ ] Every tested permutation of the tie corpus yields one pinned digest.
- [ ] The punctuation and Unicode compatibility corpus retains its primary
      American English collation order.
- [ ] Changing one path or one content byte changes the digest.
- [ ] `Sum` does not mutate the caller's file order.
- [ ] Baseline restoration tests pass without adding a comparator or changing
      a generated artifact.

## Context

- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/skillhash/hash.go`
- interface: `internal/skillhash/hash_test.go`
- interface: `internal/baseline/skills_restore.go`
- interface: `internal/baseline/skills_restore_test.go`

## Verification

- `rtk go test ./internal/skillhash -run 'TestSum'` — expected: total-order,
  compatibility, non-mutation, and negative digest cases pass.
- `rtk go test ./internal/baseline -run 'TestSkillsRestore'` — expected:
  Baseline restoration remains compatible with the shared hash authority.
- `rtk go test -race ./internal/skillhash ./internal/baseline -run 'Test(Sum|SkillsRestore)'` — expected: affected hash consumers pass under the race detector.

## References

- `_prd.md` → Core Features 2; User Story 2; Success Metrics.
- `_techspec.md` → Deterministic hash order; Testing Approach; Build Order 2.
