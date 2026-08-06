---
task: task_04
spec: 0073-skill-versions-decoupled-from-the-binary
status: pending
type: backend
complexity: medium
---

# Task 04: Make every surface use the one comparison

## Overview

Readiness is reported by the Doctor Command and gated on by every command that
requires skills. If those surfaces answer the question separately, they will
eventually answer it differently, and an operator will be told a skill is fine
by one command and blocking by another.

This slice makes them call the same comparison, and pins the boundary that
keeps third-party skills out of it.

## Requirements

1. MUST make the Doctor Command report owned-skill readiness through task_02's
   comparison.
2. MUST make every command that gates on skills use that same comparison, so
   two surfaces cannot disagree.
3. MUST name, when a skill is below the minimum, all four facts: the skill, the
   required minimum, the version found, and how to upgrade.
4. MUST report `unversioned` distinctly from both satisfying and failing, at
   every surface.
5. MUST leave third-party skills with their present treatment and MUST NOT hold
   them to a version Roundfix invented for them. Finding
   `2026-07-29-doctor-requires-roundfix-own-development-skills` records what
   happened the last time Roundfix imposed its own needs on repositories that
   had no reason to hold them.
6. MUST assert the boundary directly: a third-party skill without a version
   passes, one below an owned skill's minimum is not consulted at all.

## Subtasks

- [ ] Route Doctor through the shared comparison.
- [ ] Route every gating command through it.
- [ ] Assert the four-fact diagnostic and the third-party boundary.

## Acceptance Criteria

- [ ] Doctor reports readiness through the shared comparison.
- [ ] Every gating command reports the same state for the same skill, asserted
      by comparing two surfaces on one fixture.
- [ ] A below-minimum skill names skill, minimum, found version, and upgrade
      path.
- [ ] `unversioned` is reported distinctly at every surface.
- [ ] A third-party skill without a version passes.
- [ ] A third-party skill is never compared against an owned minimum.

## Context

- interface: `internal/cli/cli.go`
- interface: `internal/baseline/catalog_validate.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/cli ./internal/baseline -count=1 -run 'Doctor|Skills|Readiness|ThirdParty' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the surface tests ran and passed.
- `go test ./internal/cli ./internal/baseline -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Features 4, 5 and 6; Success Metrics 4 and 5.
- `_techspec.md` → System Architecture; Build Order 4.
