---
task: task_02
spec: 0054-tooling-task-and-verification-hygiene
status: pending
type: infra
complexity: medium
---

# Task 02: Ship the sanctioned regeneration target and a portable cache

## Overview

Expose the update modes as one sanctioned command, make the repository gate
behave identically inside and outside the ACP sandbox by defaulting the Go
build cache to a repository-local path, and stop the bare compile check from
leaving an unignored binary in the tree. This Task owns every authorized
change to the build tooling files.

## Requirements

1. MUST add a regeneration target that runs every update mode from Task 01 in
   dependency order, is idempotent, and is registered as a phony target.
2. MUST default the Go build cache to a repository-local, ignored path when
   the environment does not set one, and MUST leave an explicitly exported
   cache untouched so developer and CI overrides keep winning.
3. MUST ignore the repository-local cache directory and the path a bare
   `go build ./cmd/roundfix` writes, so neither can be reported as untracked
   or swept into a commit.
4. MUST NOT change what the gate verifies beyond cache determinism.
5. MUST change only `Makefile` and `.gitignore`, plus this Task file — the
   exact protected paths this Spec authorizes for build tooling.

## Subtasks

- [ ] Add the regeneration target and register it as phony.
- [ ] Default and export the build cache only when unset.
- [ ] Add the ignore entries for the cache directory and the bare-build
      binary, each with a comment naming its source.

## Acceptance Criteria

- [ ] The regeneration target runs every update mode and, on an unchanged
      repository, leaves the working tree clean; a second consecutive run
      also changes nothing.
- [ ] With no cache exported, the gate runs against the repository-local
      cache; with one exported, that value is used unchanged.
- [ ] `go build ./cmd/roundfix` leaves no path reported as untracked.
- [ ] The changed-path set for this Task is exactly `Makefile`,
      `.gitignore`, and this Task file.

## Context

- interface: `Makefile`
- interface: `.gitignore`

## Verification

- `grep -q '^baseline-digests:' Makefile` — expected: the sanctioned target exists.
- `grep -q 'GOCACHE ?=' Makefile` — expected: the cache defaults only when unset.
- `grep -q '^/roundfix$' .gitignore` — expected: the bare-build binary is ignored.
- `go build -buildvcs=false ./cmd/roundfix && git status --porcelain --untracked-files=all | grep -v '^ M docs/specs/' ; test $? -eq 1` — expected: the bare build leaves nothing untracked.
- `git status --porcelain -- Makefile .gitignore` — expected: both files are the only build-tooling paths this Task touched.

## References

`_prd.md` → User Story 3, Core Features 1, 5, 6, Project Constraints:
Tooling authority; `_techspec.md` → Build Order 2 and 3, Interfaces: Make
targets.
