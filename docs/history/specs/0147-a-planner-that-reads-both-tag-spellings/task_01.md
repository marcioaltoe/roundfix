---
task: task_01
spec: 0147-a-planner-that-reads-both-tag-spellings
status: completed
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

## Result

### Implementation

- `ParseStableVersion` accepts the stable numeric core with or without one
  leading `v` and records the accepted spelling in `Version.Prefixed`.
- The parser still validates exactly three canonical numeric identifiers. Its
  prerelease and malformed-input paths retain their existing error reasons and
  next actions.
- Parser coverage now exercises both accepted spellings, bare and prefixed
  rejection shapes, exact rejection messages, and the `version-1.2.3` near
  miss. The pre-existing assertion code remains unchanged; the obsolete table
  row that required a prefix was replaced by the new contract cases.

### Focused-check evidence

- Before implementation,
  `GOCACHE=/tmp/roundfix-spec-0147-go-cache go test -count=1 -run 'TestStableVersionAcceptsBothSpellings/bare' ./internal/releaseplan`
  failed to compile because `Version.Prefixed` did not exist.
- After implementation,
  `GOCACHE=/tmp/roundfix-spec-0147-go-cache go test -count=1 -run '^(TestParseStableVersion|TestStableVersionAcceptsBothSpellings|TestStableVersionPreservesRejectionMessages)$' ./internal/releaseplan`
  exited 0 (`ok roundfix/internal/releaseplan`).

### Acceptance evidence

- `TestStableVersionAcceptsBothSpellings` parses `1.2.3` and `v1.2.3` as the
  same numeric version and observes `Prefixed == false` and `true`,
  respectively.
- `TestParseStableVersion` rejects `1.2.3-rc.1`, `v1.2.3+build`, `1.2`,
  `01.2.3`, and `version-1.2.3`; `TestStableVersionPreservesRejectionMessages`
  locks their existing error text.
- The focused run includes the existing `TestParseStableVersion` assertions and
  exits 0. The Daemon owns the declared Verification commands and terminal
  status.
