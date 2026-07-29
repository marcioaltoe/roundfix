---
task: task_01
spec: 0061-repository-derived-skill-requirements
status: pending
type: backend
complexity: low
---

# Task 01: Accept the external requirement instead of deciding it

## Overview

Give the Repository Skill Set readiness check an entry point that receives the
external requirement from its caller, so the decision can move to a caller that
knows the repository. The existing entry point keeps working unchanged for any
caller without that context.

## Requirements

1. MUST add an exported readiness entry point that accepts the external
   required skill names alongside the repository root, and uses them verbatim
   instead of the embedded recommendation list.
2. MUST accept an empty external set as a valid requirement meaning no
   external skill is required, not as a missing argument.
3. MUST keep the existing entry point's signature and behavior intact.
4. MUST keep the owned requirement resolved from the running binary's embedded
   bundle in both entry points.
5. MUST keep every existing readiness classification, error type, and
   ownership label unchanged.

## Subtasks

- [ ] Add the external-accepting entry point delegating to the shared
      implementation.
- [ ] Cover the explicit set, the empty set, and an unchanged existing entry
      point.

## Acceptance Criteria

- [ ] Calling the new entry point with an explicit external set validates
      exactly those skills and ignores the embedded recommendation list.
- [ ] Calling it with an empty set reports zero external requirements and no
      external failure.
- [ ] The existing entry point returns the same result it returns today for
      the same repository.
- [ ] Owned skill validation is unchanged in both paths.

## Context

- interface: `skills/repository.go`
- interface: `skills/repository_test.go`
- interface: `skills/skills.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./skills/` — expected: pass, including the explicit-set and empty-set cases.

## References

`_prd.md` → Core Features 1–2; `_techspec.md` → Build Order 1, Interfaces.
