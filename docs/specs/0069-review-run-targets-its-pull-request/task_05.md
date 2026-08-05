---
task: task_05
spec: 0069-review-run-targets-its-pull-request
status: pending
type: backend
complexity: low
---

# Task 05: Make the refusal's no-side-effects claim true

## Overview

Corrective Task from the QA gate's F-001 (`Trust-Damage`). The Preflight
refusal exits `2` and prints

```text
Roundfix did not create a Run, fetch Review Source issues, start an Agent,
commit, or push.
```

while creating the Artifact Directory on the way there. The refusal is correct;
the sentence describing it is not, and this Spec's whole subject is a Run that
must have no observable effect before it validates its target.

The claim is the contract. Make it true rather than narrowing it.

## Requirements

1. MUST NOT create the Artifact Directory, or any other filesystem state, on a
   Preflight path that ends in refusal.
2. MUST create it lazily, at the point a Run genuinely needs it, so a refusal
   leaves the filesystem exactly as it found it.
3. MUST keep the existing refusal message and exit `2` unchanged; the defect is
   the behaviour, not the wording.
4. MUST hold for every refusal reason on `fetch`, `resolve`, and `watch`, not
   only the target mismatch this Spec added.
5. MUST assert absence from the filesystem, not from a log line — the log line
   is what was wrong.
6. MUST NOT change what a passing Preflight does, asserted over the existing
   tests unchanged.

## Subtasks

- [ ] Move Artifact Directory creation past every refusal path.
- [ ] Assert the directory stays absent after each refusal reason.

## Acceptance Criteria

- [ ] A target-mismatch refusal leaves a previously absent Artifact Directory
      absent, asserted by inspecting the filesystem.
- [ ] The same holds for every other Preflight refusal reason.
- [ ] The refusal message and exit `2` are unchanged.
- [ ] A Run whose Preflight passes still gets its Artifact Directory.
- [ ] `fetch`, `resolve`, and `watch` all behave identically.

## Context

- interface: `internal/preflight/preflight.go`
- interface: `internal/cli/cli.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli ./internal/preflight -count=1 -run 'Refus|Mismatch|ArtifactDir' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the refusal tests ran and passed.
- `go test ./internal/preflight ./internal/cli -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.
- `git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; the Skill already shipped in task_03 and needs no change.

## References

- `qa/qa-report-2026-08-05.md` → F-001.
- `_prd.md` → Core Feature 2; Goal 1.
- `_techspec.md` → API Contracts.
