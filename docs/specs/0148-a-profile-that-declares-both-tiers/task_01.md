---
task: task_01
spec: 0148-a-profile-that-declares-both-tiers
status: completed
type: backend
complexity: medium
---

# Task 01: Declare the incremental tier in the Profile

## Overview

Two mandatory clauses name a Profile decision that does not exist. This slice
creates it beside the complete gate, in the shape every other verification
decision already has, and exposes it without deriving it from anything.

## Requirements

1. MUST add a verification decision for the incremental tier to the Profile, with
   the identifier `verification.incremental`.
2. MUST expose it through the Profile reader beside the complete gate, each with
   its own value and source.
3. MUST NOT derive either tier from the other: changing one leaves the other's
   value and source unchanged.
4. MUST report an absent incremental decision as absent, never as an empty value
   and never as the complete gate's value.
5. MUST leave every other verification decision's identifier, kind, tool and
   command unchanged.

## Subtasks

- [ ] Add the decision to the Profile asset.
- [ ] Expose it through the reader with its own source.
- [ ] Cover independence, absence and the unchanged decisions.

## Acceptance Criteria

- [ ] The Profile declares `verification.incremental` beside the complete gate.
- [ ] Changing one tier leaves the other's value and source untouched.
- [ ] A Profile without the decision reports absence.
- [ ] Existing verification decisions are unchanged.

## Context

- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/profile_alignment.go`

## Verification

- `grep -q "verification.incremental" internal/baseline/assets/profiles/standard-typescript-monorepo.json` — expected: exit 0; the Profile declares the tier. Before this Task it does not.
- `out="$(go test -count=1 -run "^TestProfileDeclaresBothVerificationTiers$" ./internal/baseline 2>&1)" || { printf "%s\n" "$out"; exit 1; }; missing="$(printf "%s\n" "$out" | grep "no tests to run")"; test -z "$missing"` — expected: exit 0; before this Task the run reports no tests to run, so the command fails.

## References

`_prd.md` → Core Features 1-2; User Stories 2-3; Goals 1-2; Success Metric 3;
Project Constraints: Tooling authority;
`_techspec.md` → Implementation Design: The decision; API Contract 3;
Build Order 1; `_authorization.md`.

## Result

### Implementation

- Added the Profile-owned `verification.incremental` declaration with kind
  `incremental`, tool `Make`, and command `make verify-incremental`, after the
  five existing verification declarations.
- Kept the Profile reader generic: its existing per-declaration projection
  exposes the incremental tier with Profile provenance, while
  `verification.gate` retains repository-command provenance. No fallback or
  derivation connects the tiers.
- Added `TestProfileDeclaresBothVerificationTiers` at the Profile reader seam.
  It locks the five existing declarations, changes each tier independently,
  and removes the incremental declaration to exercise true absence.

### Focused checks

- Before the Profile edit,
  `rtk env GOCACHE=/private/tmp/roundfix-0148-task01-gocache go test -count=1 ./internal/baseline`
  reached the new regression and failed because the embedded Profile exposed
  only the five existing declarations.
- After the implementation,
  `rtk env GOCACHE=/private/tmp/roundfix-0148-task01-gocache go test -count=1 -run 'Test(ProfileDeclaresBothVerificationTiers|ExecutableVerificationCommandRequiresLocalDeclaration|PortableVerificationRoleMapping)$' ./internal/baseline`
  exited 0.
- `rtk env GOCACHE=/private/tmp/roundfix-0148-task01-gocache make verify-incremental`
  passed formatting and `go vet`, then exited 2 in the test phase. The Profile
  edit invalidates derived catalog and plan-characterization digests assigned
  to Task 03; two `internal/cli` process-table tests also received `operation
  not permitted` from the sandbox. This slice did not regenerate derived pins
  or alter those tests.
- The Daemon-owned commands under `## Verification` were not run.

### Acceptance evidence

- **The Profile declares `verification.incremental` beside the complete
  gate.** The asset carries the independent declaration, and the focused reader
  test observes it beside `verification.gate` with its command, role, tool,
  classification, declaration path, and digest.
- **Changing one tier leaves the other's value and source untouched.** The test
  changes the complete command and compares the full incremental projection,
  then changes the incremental command and compares the full complete
  projection.
- **A Profile without the decision reports absence.** The test removes only the
  incremental declaration, finds no incremental projection, and compares every
  remaining projection with the declared case.
- **Existing verification decisions are unchanged.** The test locks the prior
  five identifiers, kinds, tools, commands, and order; the asset diff adds one
  row without rewriting them.

### Follow-up

- Task 03 owns the sanctioned catalog and plan-characterization digest
  regeneration triggered by the authorized Profile asset change.

## Carry-forward provenance

- Source Run: `run_20260919T131447Z_ccee00b5df999ca1`
- Source commit: `828baf28e8af688bf73b06c2a1483ae26993276d`
