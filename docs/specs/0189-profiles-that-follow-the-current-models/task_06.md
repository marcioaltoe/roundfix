---
task: task_06
spec: 0189-profiles-that-follow-the-current-models
status: pending
type: backend
complexity: medium
---

# Task 06: An override equal to a fallback swaps places with the preferred

## Overview

The first delivery Run's QA gate failed on three `internal/cli` tests after task_03 made the built-in selections the Recommended Profile (verification log `batch-005-attempt-1.log` of Run `run_20260930T215154Z_eb58c6882d4f3968`).

- `TestRunImplementInteractiveInputPicksSpecThroughCollector` exposes a user-visible regression. An operator who picks `claude` for the `general` category gets the invocation Preferred Selection `claude / opus / high`, which is now also that category's built-in fallback. `ResolveProfile` then refuses the invocation profile as holding a duplicate Agent Selection. Before this Spec, the built-in fallback differed, so the same choice worked.
- `TestProveProfileSelectionsDeduplicatesReferencesAndStartsFreshProofPass` and `TestProveProfileSelectionsRetainsStableFallbackPositions` assert the proof deduplication and fallback positions through `roundconfig.Builtin()`, so they pin the replaced built-in data rather than the behavior they name.

## Requirements

1. MUST change `ResolveProfile` in `internal/config/profiles.go`: when a complete invocation override (`preferredOverride`) equals a Fallback Selection of the resolved profile, that fallback position takes the configured Preferred Selection. The chain keeps its length and order, and it stays free of duplicates. Every other override MUST keep today's behavior, and a configured profile that itself holds a duplicate MUST still be refused.
2. MUST add `internal/config/override_swap_test.go` with `TestAnOverrideEqualToAFallbackSwapsWithThePreferred` (the displaced preferred lands at the matched fallback index, the chain is otherwise unchanged) and `TestAnOverrideDistinctFromTheChainKeepsTheFallbacks`.
3. MUST change the two proof tests in `internal/cli/cli_test.go` to build an explicit fixture `roundconfig.Config` with the shape their assertions describe, instead of `roundconfig.Builtin()`:
   - three unique tuples shared across the required categories, with the listed references;
   - backend and frontend sharing one fallback tuple at the stated positions, with a distinct frontend preferred.
   Their assertions MUST stay as strict as today.
4. MUST add `TestBuiltInProfilesProveOncePerUniqueTuple` in `internal/cli/cli_test.go`. It computes the unique tuples of `roundconfig.Builtin()` over `RequiredWorkCategories()`, then asserts one exact proof request per unique tuple and no error.
5. MUST replace, in `docs/user-guide/usage.md`, the sentence "A complete override replaces only the Preferred Selection for each relevant category and keeps its configured Fallback Chain." with "A complete override replaces only the Preferred Selection for each relevant category and keeps its configured Fallback Chain; when the override equals one of those fallbacks, that fallback's position takes the configured Preferred Selection, so the chain never holds the same selection twice."
6. `TestRunImplementInteractiveInputPicksSpecThroughCollector` MUST pass unchanged.

## Subtasks

- [ ] Add the swap to `ResolveProfile` and its two tests.
- [ ] Move the two proof tests to explicit fixtures and add the built-in proof test.
- [ ] Update the user guide sentence.

## Acceptance Criteria

- [ ] The three tests the first Run failed pass, together with the new tests, and `internal/config` stays green.

## Context

- interface: `internal/config/profiles.go`
- interface: `internal/cli/cli_test.go`
- interface: `docs/user-guide/usage.md`
- creates: `internal/config/override_swap_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestAnOverrideEqualToAFallbackSwapsWithThePreferred|TestAnOverrideDistinctFromTheChainKeepsTheFallbacks)$' ./internal/config 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestAnOverrideEqualToAFallbackSwapsWithThePreferred TestAnOverrideDistinctFromTheChainKeepsTheFallbacks; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && out="$(go test -count=1 -v -run '^(TestRunImplementInteractiveInputPicksSpecThroughCollector|TestProveProfileSelectionsRetainsStableFallbackPositions|TestProveProfileSelectionsDeduplicatesReferencesAndStartsFreshProofPass|TestBuiltInProfilesProveOncePerUniqueTuple)$' ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestRunImplementInteractiveInputPicksSpecThroughCollector TestProveProfileSelectionsRetainsStableFallbackPositions TestProveProfileSelectionsDeduplicatesReferencesAndStartsFreshProofPass TestBuiltInProfilesProveOncePerUniqueTuple; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go test -count=1 ./internal/config >/dev/null && tr -s '[:space:]' ' ' < docs/user-guide/usage.md | grep -qF "that fallback's position takes the configured Preferred Selection"` — expected: exit 0; before this Task the four new tests do not exist and three of the named `internal/cli` tests fail.

## References

- task_03, task_04
- `_techspec.md` → Build Order
- Run `run_20260930T215154Z_eb58c6882d4f3968` (QA precondition refusal)
