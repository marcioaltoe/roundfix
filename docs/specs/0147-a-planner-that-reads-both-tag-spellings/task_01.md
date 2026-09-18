---
task: task_01
spec: 0147-a-planner-that-reads-both-tag-spellings
status: pending
type: backend
complexity: medium
---

# Task 01: Read a stable version with or without the prefix

## Overview

The parser requires a leading `v` and calls everything else malformed. This
slice makes the prefix optional, keeps every rejection intact, and carries the
spelling it read so nothing downstream has to guess.

## Requirements

1. MUST accept `MAJOR.MINOR.PATCH` and `vMAJOR.MINOR.PATCH` as stable versions.
2. MUST keep rejecting malformed input, prerelease identifiers and build
   metadata, with the messages they produce today.
3. MUST consume only a leading `v` immediately followed by the numeric core, so
   a tag like `version-1.2.3` stays malformed.
4. MUST record, on the parsed value, whether the tag carried the prefix.
5. MUST leave every existing parser test passing with its assertions unchanged.
6. MUST name the recorded field `Prefixed`, so the Task's own check and later
   readers agree.

## Subtasks

- [ ] Make the prefix optional and record it.
- [ ] Keep the numeric-core and canonical-identifier rules as they are.
- [ ] Cover both spellings, the prerelease and build-metadata rejections, and a
      near-miss like `version-1.2.3`.

## Acceptance Criteria

- [ ] `1.2.3` and `v1.2.3` both parse, and their recorded spelling differs.
- [ ] `1.2.3-rc.1`, `v1.2.3+build`, `1.2`, `01.2.3` and `version-1.2.3` are
      rejected as they are today.
- [ ] Existing parser tests pass unedited.

## Context

- interface: `internal/releaseplan/version.go`

## Verification

- `grep -q "Prefixed" internal/releaseplan/version.go` — expected: exit 0; the parsed value records its spelling. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestStableVersionAcceptsBothSpellings$" ./internal/releaseplan 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Feature 1; User Stories 1-2; Goals 1-2; Success Metric 1;
Regression locks;
`_techspec.md` → Implementation Design: Both spellings, one version; Interfaces;
API Contract 1; Build Order 1.
