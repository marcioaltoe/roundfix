---
task: task_01
spec: 0148-a-profile-that-declares-both-tiers
status: pending
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
