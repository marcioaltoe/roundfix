---
task: task_03
spec: 0052-claude-adapter-standardization
status: completed
type: backend
complexity: low
---

# Task 03: Default the frontend profile to the proven Claude tuple

## Overview

Change the built-in `frontend` Agent Selection Profile's Preferred Selection
from the unadvertised `claude / claude-opus-5 / xhigh` to the proven
`claude / opus / xhigh`, so a fresh installation's frontend routing proves
against the official adapter without configuration edits. The Fallback Chain
and every other built-in profile are untouched.

## Requirements

1. MUST change the built-in `frontend` Preferred Selection model to `opus`,
   keeping runtime `claude` and reasoning effort `xhigh`.
2. MUST keep the `frontend` Fallback Chain and all other built-in profiles
   byte-identical in behavior.
3. MUST keep the rendered default Project Config in agreement with the
   built-in value.
4. MUST NOT add adapter aliases to the Model Catalog; `opus` is proven
   through Exact Agent Selection Proof, not catalog membership.

## Subtasks

- [ ] Update the built-in `frontend` profile value.
- [ ] Update the default Project Config rendering that mirrors it.
- [ ] Update every test that pins the old `claude-opus-5` frontend default.

## Acceptance Criteria

- [ ] Resolving the `frontend` profile with no user or project configuration
      yields preferred `claude / opus / xhigh` with the unchanged fallback.
- [ ] The generated default Project Config names `opus` for the frontend
      preferred model.
- [ ] The Model Catalog's Claude identifiers are unchanged.

## Context

- interface: `internal/config/profiles.go`
- interface: `internal/config/config.go`
- interface: `internal/agent/catalog.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/config/` — expected: pass.
- `grep -n 'claude-opus-5' internal/config/profiles.go ; test $? -eq 1` — expected: no matches (exit 1).
- `grep -c 'claude-opus-5' internal/agent/catalog.go | grep -x 1` — expected: `1`; the catalog identifier is preserved.

## References

`_prd.md` → User Story 5, Core Feature 6; `_techspec.md` → Build Order 3,
Data Models; ADR-0049.

## Result

### Implementation

- Changed only the built-in `frontend` Preferred Selection model from
  `claude-opus-5` to `opus`; runtime `claude`, reasoning effort `xhigh`, and
  the existing Codex fallback remain unchanged.
- Changed the rendered default Project Config to emit the same
  `claude / opus / xhigh` Preferred Selection.
- Updated the config and CLI expectations that derive from the built-in
  frontend default. Catalog and documentation-contract expectations for the
  separate `claude-opus-5` catalog identifier remain unchanged.

### Focused checks

- Red signal before the production edit:
  `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run
  'Test(AgentSelectionProfileBuiltinsResolveRequiredCategories|InitCreatesUserConfig|ProfileGeneratedConfigUsesCompleteProfilesSchema|InitForceOverwritesExistingConfig)$'
  ./internal/config` failed because the resolver and rendered config still
  returned `claude-opus-5`; the matching focused CLI tests failed for the same
  old default.
- `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run
  'Test(AgentSelectionProfileBuiltinsResolveRequiredCategories|InitCreatesUserConfig|ProfileGeneratedConfigUsesCompleteProfilesSchema|InitForceOverwritesExistingConfig)$'
  ./internal/config` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run
  'Test(RunDoctorProfileReadiness|RunInitForceOverwritesExistingConfig|ProfilesValidateDeduplicatesProofsAndReportsEveryReference|InvocationProfileOverrideOmittedUsesTaskQAAndReviewProfiles|RunSetupProfileProofsEveryDistinctTupleOnceBeforePersistence)'
  ./internal/cli` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -run
  '^TestModelCatalogsExposeOrderedPickerData$' ./internal/agent` — passed.
- `rtk git -c core.fsmonitor=false diff --exit-code --
  internal/agent/catalog.go` — passed with no diff.
- The first focused test attempt used Go's default cache and could not start
  because the sandbox denied writes under `~/Library/Caches/go-build`; the
  Task-local `GOCACHE` reruns above reached the tests.
- The Task's declared `## Verification` commands were not run; the Daemon
  owns that gate.

### Acceptance criteria evidence

- The built-in resolver test passes against the complete `frontend` profile:
  preferred `claude / opus / xhigh` and fallback
  `codex / gpt-5.6-sol / high`.
- The generated-config test passes against the exact frontend YAML block,
  including `model: opus` and the unchanged fallback; the CLI init path also
  passes with the new rendered value.
- The Claude Model Catalog behavior test passes, and
  `internal/agent/catalog.go` has no diff.

### Follow-ups

None discovered within this Task's slice.
