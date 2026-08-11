---
task: task_02
spec: 0051-doctor-readiness-contract-reconciliation
status: completed
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

- [x] Add a failing permutation case for collation-equal distinct paths.
- [x] Define the normalized raw-path tie-breaker after collator equality.
- [x] Preserve punctuation, case, path-depth, and content-change regressions.
- [x] Exercise the existing Baseline restoration consumer against the shared
      implementation.

## Acceptance Criteria

- [x] Every tested permutation of the tie corpus yields one pinned digest.
- [x] The punctuation and Unicode compatibility corpus retains its primary
      American English collation order.
- [x] Changing one path or one content byte changes the digest.
- [x] `Sum` does not mutate the caller's file order.
- [x] Baseline restoration tests pass without adding a comparator or changing
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

## Result

`Sum` now applies ordinary Go string ordering to distinct normalized paths
only when American English collation reports equality. The existing collation
order, slash normalization, path-plus-content digest shape, and caller-owned
slice order remain unchanged.

The new regression test reproduced the defect before the implementation
change: the precomposed-first permutation produced
`dc5f4008143541df5bfab2e71ac7f1a0cbd8f9963ea00747b33f630894108f91`
instead of the pinned total-order digest
`2692b0010be42e60169e655aa45186b0d836be24c81a1c3aa9652d2614fc4451`.
After the change, both permutations produce the pinned digest.

Acceptance evidence:

- `TestSumIsPermutationIndependentForCollationEqualPaths` passes for both
  precomposed-first and decomposed-first input.
- `TestSumPreservesSkillsCLI1519PrimaryCollation` and
  `TestSumSortsUnderscoreBeforeHyphen` preserve the punctuation, case,
  path-depth, and Unicode primary-collation corpus.
- `TestSumChangesWhenPathOrContentChanges` passes for both negative digest
  cases.
- `TestSumPreservesSkillsCLI1519PrimaryCollation` confirms that `Sum` leaves
  the caller-owned slice unchanged.
- `TestSkillsRestore*` passes through the existing `skillhash.Sum` consumer;
  no Baseline comparator, asset, lock file, or generated artifact changed.

Verification:

- `rtk go test ./internal/skillhash -run 'TestSum'` — passed, 8 tests.
- `rtk go test ./internal/baseline -run 'TestSkillsRestore'` — passed, 14
  tests.
- `rtk go test -race ./internal/skillhash ./internal/baseline -run
  'Test(Sum|SkillsRestore)'` — passed, 22 tests in 2 packages.

Follow-ups: none.
