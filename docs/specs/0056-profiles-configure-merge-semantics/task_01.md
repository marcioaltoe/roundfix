---
task: task_01
spec: 0056-profiles-configure-merge-semantics
status: pending
type: test
complexity: medium
---

# Task 01: Record how the writer behaves today

## Overview

Every later slice changes how the config file is written, and the Spec promises
that a file the current writer handles can never fail or be corrupted by the
new one. This Task records that behavior first, over the config shapes that
actually break naive YAML rewriting. It changes no product behavior — it is the
regression gate the rest of the Spec is measured against, and it only works if
it lands before anything else moves.

## Requirements

1. MUST assemble a corpus of real config shapes and record, for each, the exact
   bytes today's writer produces for a representative configure operation.
2. MUST cover at minimum: a five-category config with Fallback Chains, comments
   above and beside entries, non-default indentation, key order that is not
   alphabetical, a multiline scalar, a YAML anchor and an alias to it, non-ASCII
   values, an empty file, and a file with no `profiles` section.
3. MUST fail with a readable diff naming the affected config and the changed
   region when a later change alters any recorded output.
4. MUST be regenerable through an explicit flag so an intended change is
   re-recorded deliberately, never silently.
5. MUST NOT change any production behavior, exported API, command surface, or
   output in this Task.

## Subtasks

- [ ] Assemble the config corpus covering every listed shape.
- [ ] Record today's writer output for each as a golden.
- [ ] Add the comparison with a readable failure diff.
- [ ] Add the explicit regeneration flag.

## Acceptance Criteria

- [ ] The corpus contains a case for each shape named in Requirement 2.
- [ ] A test writes each corpus config through the current writer and compares
      against its golden, passing on the unmodified tree.
- [ ] Deliberately altering one recorded golden makes the test fail and name the
      affected config.
- [ ] Running the comparison twice in a row produces the same result and does
      not rewrite its own goldens.
- [ ] The goldens can be re-recorded only through the explicit flag.
- [ ] `git status --porcelain` shows no path outside `internal/config/` and this
      task file.

## Context

- interface: `internal/config/profile_config.go`
- interface: `internal/config/config_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0; the corpus matches the current writer.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0 on a second consecutive run, proving the comparison is
  stable and self-recording is gated.
- `grep -rq "TestProfilesConfigWriterCharacterization" internal/config` —
  expected: exit 0; the test exists at the declared name.
- `go test ./internal/config -count=1` — expected: exit 0; the existing suite is
  unaffected.

## References

- `_prd.md` → User Story 6; Core Features 7; Success Metrics (round-trips with
  only the intended change).
- `_techspec.md` → Testing Approach: characterization corpus; Build Order 1;
  Risks (anchors and aliases).
