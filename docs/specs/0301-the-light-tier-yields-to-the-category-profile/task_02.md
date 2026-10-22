---
task: task_02
spec: 0301-the-light-tier-yields-to-the-category-profile
status: pending
type: backend
complexity: medium
---

# Task 02: A light Task skips or yields the Light Tier when the light session cannot serve it

## Overview

Under `defaults.agent_full_access: true` every light-eligible Task fails before
Agent work, because OpenCode offers no full-access mode and the Daemon treats
the light session's access-policy refusal, like any non-selection failure, as
a Task failure. This Task adds a shared helper that keeps the light selections
whose runtime serves the requested access, makes `applyLightTier` skip the
Light Tier with the reason code `access_mode_unavailable` when none remains,
and makes any other pre-work refusal of a light session yield to the
category's Preferred Selection.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-21,
   [`references/2026-10-21-the-light-tier-fails-a-task-under-agent-full-access.md`](references/2026-10-21-the-light-tier-fails-a-task-under-agent-full-access.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-3 and 5 of
   `_techspec.md`.
2. MUST add `Selections`, `AccessUnavailable` and `ServingSelections` in the
   new file `internal/lighttier/access.go` with the signatures of
   `_techspec.md` → Interfaces and API Contract 1. A runtime factory error
   keeps the selection; only an `*agent.AccessPolicyError` whose `Predicate`
   is `agent.AccessModeAvailablePredicate` refuses it.
3. MUST make `applyLightTier` in `internal/daemon/agent_session_owner.go`
   call `lighttier.ServingSelections(plan.Models, owner.runtimeFactory)`
   before the key, spend log and ceiling checks and, when it keeps none and
   refused at least one, skip the tier exactly as API Contract 2 states: the
   warning `roundfix: warning: light tier skipped for Task <id>: light runtime
   <runtime> offers no mode for access policy "<policy>"; it runs on its
   <category> profile` and the `light_tier_skipped` event with `reason_code`
   `access_mode_unavailable`. Otherwise the light candidates are the kept
   selections, built through `lighttier.Selections`, and `lightCandidates`
   is their count. The three existing skip reasons keep their messages and
   codes.
4. MUST make the owner yield a light candidate's pre-work failure exactly as
   API Contract 3 states, in `activate` (the runtime factory's error and a
   `prepareSession` error that `isSelectionStartFailure` does not cover) and
   in `Run` (an error of the light prompt, after the lost rollout and provider
   system error branches, that is not a `SelectionFailureError`): close the
   light session, record the attempt `failed`, publish the
   `light_tier_yielded` Task event with `model`, `reason_code` (the
   `AccessPolicyError`'s `Kind`, else `light_session_failed`) and `reason`,
   print `roundfix: warning: light tier yielded for Task <id> before Agent
   work: <cause>; it runs on its <category> profile`, publish the Agent
   Selection fallback notification to the category's Preferred Selection and
   continue at that candidate without opening the remaining light models.
5. MUST state the hazard in the code's behavior: the yield sends the Task's
   prompt to the category's runtime, which receives the same requested access
   policy; it never runs a light session with a weaker policy, never yields
   after `AGENT_WORK_STARTED`, never yields when the context is
   done or the error matches `context.Canceled`, and never changes a non-light
   candidate's handling.
6. MUST keep API Contract 5: `internal/lighttier/tier.go`,
   `internal/daemon/task_engine.go`, `internal/agent` and
   `internal/cli/implement.go` keep their bytes; the escalation after a failed
   Verification and the selection fallback between light models are
   unchanged. No existing test, golden or fixture changes.
7. MUST add the tests `_techspec.md` → Testing Approach lists for this Task,
   each top-level test calling `t.Parallel()` first, compiling for Windows,
   using temporary homes and repositories, passing no real-time duration and
   starting no real ACP runtime:
   - `TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess` in
     the new file `internal/lighttier/access_test.go`, whose full-access case
     asserts that the refused `opencode` selection's error message is exactly
     `requested access policy "full-access" failed predicate
     "runtime_has_access_mode" for runtime "opencode"; recovery: disable full
     access or select a runtime that supports it`, the Fluxus text;
   - `TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess`,
     `TestLightTierYieldsToTheCategoryProfileBeforeAgentWork` and
     `TestLightTierDoesNotYieldAfterWorkOrOutsideTheLightTier` in the new file
     `internal/daemon/light_tier_access_test.go`, with the subtests listed
     there; the runner wrapper's `PrepareSession` returns
     `req.Runtime.ValidateRequestedAccessPolicy()` after the fake runner's own
     result, and the full-access subtests build runtimes with
     `agent.RuntimeFor` and `EnableFullAccess: true`;
   - `TestImplementSkipsTheLightTierUnderProjectFullAccess` in the new file
     `internal/cli/implement_light_tier_access_test.go`, through `roundfix
     implement` with `defaults.agent_full_access: true` in the temporary
     repository's `.roundfixrc.yml`, a light model in the temporary User
     Config and a fixture key value built at runtime, asserting that no key
     value reaches stdout or stderr.

## Subtasks

- [ ] Add the light selection helper and its test.
- [ ] Skip the Light Tier with `access_mode_unavailable` when no light runtime serves the requested access.
- [ ] Yield a light candidate's pre-work refusal to the category's Preferred Selection.
- [ ] Add the Daemon and `roundfix implement` tests.

## Acceptance Criteria

- [ ] A `low` Task under full access completes on its category profile with one `access_mode_unavailable` skip event and no light session opened.
- [ ] A light runtime factory error, an access-policy setup refusal and a light prompt failure before Agent output each yield to the category profile with one `light_tier_yielded` event.
- [ ] A failure after Agent output and a non-light setup failure still fail the Task, and every existing Light Tier and selection fallback test named in the Verification passes unchanged.

## Consumers

- `lighttier.Selections` and `lighttier.ServingSelections` — consumers: `./internal/daemon`, `./internal/cli`
- The `access_mode_unavailable` skip reason and the `light_tier_yielded` phase of the `daemon.task` Run Event — consumers: `./internal/daemon`, `./internal/cli`

## Context

- instruction: `docs/adr/0314-a-light-tier-that-cannot-serve-the-run-yields-to-the-category-profile.md`
- instruction: `internal/lighttier/tier.go`
- instruction: `internal/lighttier/tier_test.go`
- instruction: `internal/agent/agent.go`
- instruction: `internal/agent/acpx_runner.go`
- instruction: `internal/daemon/task_engine.go`
- instruction: `internal/daemon/task_engine_test.go`
- instruction: `internal/daemon/light_tier_test.go`
- instruction: `internal/cli/implement.go`
- instruction: `internal/cli/profile_preflight.go`
- instruction: `internal/cli/implement_test.go`
- instruction: `internal/cli/implement_light_tier_test.go`
- creates: `internal/lighttier/access.go`
- creates: `internal/lighttier/access_test.go`
- creates: `internal/daemon/light_tier_access_test.go`
- creates: `internal/cli/implement_light_tier_access_test.go`
- interface: `internal/daemon/agent_session_owner.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess|TestTierFor|TestPlanEnabled|TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess|TestLightTierYieldsToTheCategoryProfileBeforeAgentWork|TestLightTierDoesNotYieldAfterWorkOrOutsideTheLightTier|TestLightTierDispatchOrder|TestLightTierSkipsWithAWarning|TestLightTierRecordsSpend|TestLightTierEscalatesOnce|TestNoFallbackAfterAgentWorkStarted|TestAgentSelectionFallbackPublishesBeforeNextSession|TestImplementSkipsTheLightTierUnderProjectFullAccess|TestImplementBuildsTheLightTierPlan)$' ./internal/lighttier ./internal/daemon ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess TestTierFor TestPlanEnabled TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess TestLightTierYieldsToTheCategoryProfileBeforeAgentWork TestLightTierDoesNotYieldAfterWorkOrOutsideTheLightTier TestLightTierDispatchOrder TestLightTierSkipsWithAWarning TestLightTierRecordsSpend TestLightTierEscalatesOnce TestNoFallbackAfterAgentWorkStarted TestAgentSelectionFallbackPublishesBeforeNextSession TestImplementSkipsTheLightTierUnderProjectFullAccess TestImplementBuildsTheLightTierPlan; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/lighttier ./internal/daemon ./internal/cli` — expected: exit 0. Before this Task the five new tests and `lighttier.ServingSelections` do not exist, so the command fails at its first pass check. Afterwards the helper, skip, yield and command tests pass with the unchanged Light Tier, selection fallback and implement tests, and the three packages' tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-2; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior; API Contracts 1-3 and 5; Invariants 1-5; Testing Approach; Build Order 2
- ADR-0314
