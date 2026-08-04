---
task: task_01
spec: 0070-declared-unreachable-acceptance
status: pending
type: backend
complexity: medium
---

# Task 01: Read a Spec's declared unreachable acceptance

## Overview

Everything downstream reads this: the gate matches rows against it, the archive
boundary trusts a count derived from it, and the archive record stamps the
action it names. This slice delivers the declaration format and its reader,
verifiable on its own through fixture Spec folders.

A Spec that declares nothing is the ordinary case and must stay ordinary — an
absent section is not an error.

## Requirements

1. MUST read an `## Unreachable Acceptance` section from a Spec's PRD into
   declarations, each carrying the criterion it covers, the reason no hermetic
   Verification can reach it, the human action that would satisfy it, and the
   line it was declared on.
2. MUST treat an absent section as zero declarations and no error.
3. MUST report a malformed declaration — a missing criterion, reason, or
   satisfying action — as a typed error naming the file and line, never by
   silently dropping the entry. A dropped declaration would later read as a
   Spec that declared nothing.
4. MUST NOT infer, judge, or supply any of the three fields. The Spec's author
   declares; this code reads.
5. MUST leave every existing `internal/spec` export and behavior unchanged.

## Subtasks

- [ ] Add the declaration type and its reader.
- [ ] Parse the three fields with their declaration lines.
- [ ] Return typed errors for each malformed shape.
- [ ] Add fixtures: present, absent, and each malformed shape.

## Acceptance Criteria

- [ ] A fixture Spec declaring two unreachable acceptances returns both, each
      with its criterion, reason, satisfying action, and 1-based line.
- [ ] A fixture Spec with no such section returns zero declarations and no
      error.
- [ ] A declaration missing its reason returns a typed error naming the file
      and the line.
- [ ] A declaration missing its satisfying action returns a typed error naming
      the file and the line.
- [ ] No malformed declaration is silently dropped, asserted by a test that
      counts returned declarations against the fixture's entries.
- [ ] The existing `internal/spec` tests pass unchanged.

## Context

- interface: `internal/spec/spec.go`
- interface: `internal/spec/errors.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/spec -count=1 -run 'Unreachable|Declaration' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the declaration tests ran and passed.
- `go test ./internal/spec -count=1` — expected: exit 0; nothing in the loader
  regressed.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Feature 1; Goals; Decisions.
- `_techspec.md` → Interfaces; Build Order 1.
- ADR-0080.
