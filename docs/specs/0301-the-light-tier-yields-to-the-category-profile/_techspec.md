---
spec: 0301-the-light-tier-yields-to-the-category-profile
prd: _prd.md
created: 2026-10-22
---

# The Light Tier yields to the category profile — Technical Spec

## Executive Summary

`applyLightTier` in `internal/daemon/agent_session_owner.go` prepends one
`opencode` selection per light model to a light Task's profile without asking
whether that runtime can serve the Run's requested access policy. Under
`defaults.agent_full_access: true` OpenCode cannot: `agent.RuntimeFor` gives
it no full-access candidate, so the session setup's first access check returns
an `*agent.AccessPolicyError` with the predicate `runtime_has_access_mode`.
`activate` falls to the next candidate only for a selection start failure, so
that error, like any other error the runtime factory or the setup returns,
fails the Task before Agent work. This Spec adds one shared helper in
`internal/lighttier` that keeps the light selections whose runtime serves the
requested access; `applyLightTier` skips the Light Tier with the reason code
`access_mode_unavailable` when none remains; the owner yields a light
selection's other pre-work failures to the category's Preferred Selection; and
Doctor uses the same helper to print a `light-tier:` line.

The trade-off accepted: a light Task under full access always runs on its
category profile, at its cost, until OpenCode gains a mode Roundfix can
request; running OpenCode without the requested access stays rejected,
because the Verification that needed full access would fail later.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change adds a reason code and an event phase to the existing
  `daemon.task` Run Event and one Doctor line. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses fake runners, fake or
  static runtime factories and temporary homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 2 and 3 follow ADR-0314,
  named with its title in `_prd.md` → Project Constraints, and keep the
  decisions it refines. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 ("same as cycle AO"), carrying the cycle AJ grant
  "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0301-the-light-tier-yields-to-the-category-profile/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at `42f523c2`. Line numbers move when
other Specs merge; the function names do not.

1. `taskAgentSessionOwner` (`internal/daemon/agent_session_owner.go:141`)
   calls `applyLightTier` (`:1022`) for a Task that `lighttier.TierFor` calls
   `light` while `plan.LightTier.Enabled()`. `applyLightTier` skips the tier
   for `key_missing`, `spend_unreadable` and `ceiling_reached`, printing
   `roundfix: warning: light tier skipped for Task <id>: <reason>; it runs on
   its <category> profile` and publishing a `daemon.task` event with
   `{"task", "phase": "light_tier_skipped", "reason_code"}`. Otherwise it
   builds `{Runtime: "opencode", Model: "openrouter/" + model}` per light
   model, followed by the category's Preferred Selection and Fallback Chain,
   sets the profile source `light-tier` and `lightCandidates`.
2. `activate` (`:397`) returns the runtime factory's error as is (`:411`) and,
   when `prepareSession` fails, falls to the next candidate only when
   `isSelectionStartFailure` (`:844`) holds: a `SelectionFailureError`,
   `SelectionPreflightError` or `AdapterProbeError` (`:434`). Any other error
   is returned and fails the Task.
3. `Run` (`:217`) handles a lost rollout and a provider system error first,
   falls back after a `SelectionFailureError` while `AGENT_WORK_STARTED` has
   not been published, and returns every other error (`:264`).
4. `agent.RuntimeFor` (`internal/agent/agent.go`, around `:380-420`) names
   full-access candidates only for `codex` (`agent-full-access`,
   `full-access`) and `claude` (`bypassPermissions`); with
   `EnableFullAccess` every runtime gets `RequestedAccessPolicy` `full-access`.
   `ValidateRequestedAccessPolicy` (`:90`) then returns `*AccessPolicyError`
   with `Kind` `access_policy_unsupported` and `Predicate`
   `runtime_has_access_mode` for `opencode`. `ACPXRunner.ensureSession`
   (`internal/agent/acpx_runner.go`, around `:1305`) calls `applyFullAccess`
   (`:1363`), whose first step is that check.
5. Characterization, a `go test -overlay` probe in `internal/daemon` that
   wrote nothing to the repository, using `lightPlanForTest`, a runtime
   factory over `agent.RuntimeFor` and a runner whose `PrepareSession`
   returns `req.Runtime.ValidateRequestedAccessPolicy()` after the fake
   runner's own: with full access, the Task cycle reported 0 completed and 1
   failed, prepared one `opencode/openrouter/deepseek/deepseek-v4.1-flash`
   session, ran no prompt, published no skip event, and printed `Task task_01
   failed: Agent failed: requested access policy "full-access" failed
   predicate "runtime_has_access_mode" for runtime "opencode"; recovery:
   disable full access or select a runtime that supports it`, the Fluxus
   message. Without full access it completed on the light model. With a
   factory that refuses `opencode`, it failed with the factory's error and
   prepared nothing.
6. Characterization through the command, an overlay probe in `internal/cli`:
   `roundfix implement` with `defaults.agent_full_access: true` in
   `.roundfixrc.yml`, `openrouter.light_models: [x-ai/grok-4.5]` in User
   Config and a `low` docs Task sent its one work request to `opencode` with
   the requested policy `full-access`, whose access check fails. The command
   fake has no session setup, so it completed; the real `ACPXRunner` refuses
   there.
7. The Verification escalation (`internal/daemon/task_engine.go`, around
   `:2079`) calls `escalate` (`:1059`) only while `isLightActive` and not yet
   escalated, and publishes `light_tier_escalated`.
8. `runDoctorCommand` (`internal/cli/doctor.go:60`) prints
   `recommendations:` (`:102`) then `pre-pr-review:` (`:103`); no check reads
   `openrouter.light_models` or `defaults.agent_full_access`. The only
   conditional line is `linux-lane:`.
   `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`
   (`internal/cli/doctor_linux_lane_test.go`) and the exact transcripts in
   `internal/cli/doctor_test.go` run with the built-in configuration, where
   `defaults.agent_full_access` is `false`.

Sweep of every place that builds light selections or decides the access
policy of a session (`rg -n 'openrouter/' internal -g '!*_test.go'`,
`rg -n 'ValidateRequestedAccessPolicy|AccessModeFor' internal -g '!*_test.go'`,
`rg -n 'lighttier\.' internal -g '!*_test.go'`):

- `applyLightTier` is the only builder of light selections; the helper
  replaces its inline loop.
- `ACPXRunner.applyFullAccess`, `ACPXRunner.ProveExactSelection`
  (`internal/agent/selection_assignment.go:109`) and
  `proveProfileSelections` (`internal/cli/profiles_validate.go:209`) check
  configured profiles and Run sessions; unchanged.
- `implementLightTierPlan` (`internal/cli/implement.go:1432`) builds the plan
  from configuration; unchanged.
- Text: `docs/user-guide/configuration.md` (Light implementation tier: "names
  one of the three reasons"), `.agents/skills/roundfix/references/runtime.md`
  (Light implementation tier: three reason codes),
  `docs/user-guide/commands/doctor.md` and
  `.agents/skills/roundfix/references/setup.md` (Doctor lines and codes), and
  the glossary's **Light Tier** entry.

## System Architecture

- `internal/lighttier/access.go` (new): `Selections` and `ServingSelections`.
- `internal/daemon/agent_session_owner.go`: the access skip in
  `applyLightTier`; the yield in `activate` and `Run`.
- `internal/cli/health.go`: `HealthCheckLightTier`.
- `internal/cli/doctor.go`: `doctorLightTierResult` and its place in
  `runDoctorCommand`.
- `docs/user-guide/configuration.md`, `docs/user-guide/commands/doctor.md`,
  the Roundfix `runtime` and `setup` references and their mirrors, the
  Roundfix skill's versions, `CONTEXT.md`.

`internal/lighttier` gains an import of `internal/agent`, which imports
neither `internal/lighttier` nor anything that does; no other import edge
changes.

## Implementation Design

### Interfaces

```go
// Selections returns one light Agent Selection per model, in order: runtime
// "opencode", model "openrouter/<model>", no reasoning effort.
func Selections(models []string) []roundconfig.AgentSelection

// AccessUnavailable is a light selection whose runtime offers no mode for the
// access policy its runtime requests.
type AccessUnavailable struct {
	Selection roundconfig.AgentSelection
	Err       *agent.AccessPolicyError
}

// ServingSelections builds each light selection's runtime with runtimeFor and
// keeps it unless ValidateRequestedAccessPolicy fails with the predicate
// runtime_has_access_mode. A runtimeFor error keeps the selection, so the
// session start decides it.
func ServingSelections(models []string, runtimeFor func(roundconfig.AgentSelection) (agent.RuntimeSpec, error)) ([]roundconfig.AgentSelection, []AccessUnavailable)

// doctorLightTierResult reports the light-tier line, and whether Doctor prints it.
func doctorLightTierResult(config roundconfig.Config, runtimeFor func(roundconfig.AgentSelection) (agent.RuntimeSpec, error)) (CheckResult, bool)
```

### Data Models

None changes shape. The `daemon.task` Run Event payload keeps its keys for
`light_tier_skipped` and gains the value `access_mode_unavailable` for
`reason_code`. The new phase `light_tier_yielded` carries `task`, `phase`,
`model`, `reason_code` and `reason`.

### API Contracts

1. API Contract: the helper. `Selections` returns exactly the selections
   `applyLightTier` builds today. `ServingSelections` keeps the input order;
   an empty model list returns two empty results; under the runtime-default
   policy it keeps every selection.
2. API Contract: the access skip. `applyLightTier` first calls
   `lighttier.ServingSelections(plan.Models, owner.runtimeFactory)`. When it
   keeps none and refused at least one, the tier is skipped with the reason
   code `access_mode_unavailable` and the reason `light runtime <runtime>
   offers no mode for access policy "<policy>"`, from the first refusal's
   `Runtime` and `Policy`, so the warning reads `roundfix: warning: light
   tier skipped for Task <id>: light runtime opencode offers no mode for
   access policy "full-access"; it runs on its <category> profile`, and the
   event is `{"task", "phase": "light_tier_skipped", "reason_code":
   "access_mode_unavailable"}`. This check runs before the key, spend log and
   ceiling checks, whose messages and codes keep their bytes. Otherwise the
   light candidates are the kept selections, followed by the category's
   Preferred Selection and Fallback Chain, and `lightCandidates` is their
   count.
3. API Contract: the yield. While the active or starting candidate's index is
   below `lightCandidates` and `AGENT_WORK_STARTED` has not been published,
   a failure is yielded when the context is not done and the error does not
   match `context.Canceled`, and it is none of a `SelectionFailureError`,
   `SelectionPreflightError` or `AdapterProbeError` (which keep falling to
   the next candidate), a lost rollout or a provider system error (which
   keep their recovery). Yielded failures are the runtime factory's error and
   a `prepareSession` error in `activate`, and an error the light prompt
   returns in `Run`. The yield closes the light session, records the attempt
   `failed` with the cause, publishes the `daemon.task` event `{"task",
   "phase": "light_tier_yielded", "model": <light model>, "reason_code":
   <AccessPolicyError Kind, else "light_session_failed">, "reason": <cause>}`,
   prints `roundfix: warning: light tier yielded for Task <id> before Agent
   work: <cause>; it runs on its <category> profile`, publishes the Agent
   Selection fallback notification to the category's Preferred Selection,
   and continues at that candidate, never opening the remaining light
   models. A failure after `AGENT_WORK_STARTED`, a cancellation, and every
   failure of a non-light candidate keep today's handling.
4. API Contract: the Doctor line. `doctorLightTierResult` returns `false`
   unless the configuration's `openrouter.light_models` is not empty and
   `defaults.agent_full_access` is `true`. It then calls `ServingSelections`
   with `runtimeFor`; with no kept selection the line is `light-tier: warn
   (DR-LIGHT-ACCESS-MODE: light runtime opencode offers no full-access mode;
   light Tasks run on their category profile; next: set
   openrouter.light_models: [] in User Config, or defaults.agent_full_access:
   false)`; otherwise `light-tier: ok (light runtime <runtime> serves full
   access)`, naming the first kept selection's runtime. `runDoctorCommand`
   passes a factory over `runtimeForProfileSelectionWithOptions` with
   `EnableFullAccess` from the configuration and prints the line right after
   `recommendations:` and before `pre-pr-review:`. The line never fails
   Doctor, opens no Agent Session and reads no key.
5. API Contract: unchanged behavior. `TierFor`, `implementLightTierPlan`, the
   three existing skip reasons, the selection fallback between light models,
   `escalate` and `light_tier_escalated`, the no-fallback-after-work rule,
   `internal/agent` and `roundfix profiles validate` keep their bytes and
   behavior; every existing Doctor line keeps its bytes.
6. API Contract: documentation. The configuration guide's Light
   implementation tier section and the Roundfix `runtime` reference describe
   the access skip with its reason code and the pre-work yield with its
   phase; the Doctor guide and the Roundfix `setup` reference describe the
   `light-tier:` line and `DR-LIGHT-ACCESS-MODE`; the glossary revises
   **Light Tier**.

### Surface Transcripts

None. No command, flag, exit code or stream placement changes. Doctor gains
one conditional line and the Daemon one warning form and one skip reason; the
tests of API Contracts 2-4 assert those exact lines, and every existing line
keeps its bytes.

### Invariants

```text
1. No light session is opened when no light runtime serves the requested access policy.
2. A light session never runs with an access policy weaker than the one the Run requested.
3. No yield happens after AGENT_WORK_STARTED, on cancellation, or for a non-light candidate.
4. The key_missing, spend_unreadable and ceiling_reached skips keep their messages and codes.
5. A selection start failure of a light model still falls to the next light model.
6. Doctor's light-tier line appears only with a non-empty light model list and full access, never fails Doctor, and leaves every other line byte-identical.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metrics 1-3, Invariants 1-2 → the helper and the access skip (API Contracts 1-2)
- Goal 2, Core Feature 2, Success Metric 4, Invariants 3 and 5 → the yield (API Contract 3)
- Goal 3, Success Metric 5, Invariants 4-5 → unchanged behavior (API Contract 5)
- Goal 4, Core Feature 3, Success Metric 6, Invariant 6 → the Doctor line (API Contract 4)
- Goal 5, Core Feature 4 → the guides, the references and the glossary (API Contract 6)
- Success Metric 7 → the PRD's Unreachable Acceptance

## Absence Evidence

1. `ServingSelections` and `Selections` — command: `grep -rn -e ServingSelections -e 'func Selections' internal`; result: no match at `42f523c2`; `applyLightTier` builds the light selections inline.
2. The `access_mode_unavailable` reason and the `light_tier_yielded` phase — command: `grep -rn -e access_mode_unavailable -e light_tier_yielded internal docs .agents skills`; result: no match at `42f523c2`.
3. A Doctor check for the Light Tier — command: `grep -rn -e light-tier -e DR-LIGHT internal/cli docs/user-guide .agents/skills/roundfix`; result: no match at `42f523c2`.

## ADR Alternatives

- ADR-0314: no rejected alternative is taken; the light session never runs without the requested access, pre-work refusals of the light runtime do not walk the remaining light models, Preflight does not refuse the Run, and `profiles validate` does not prove light models.

## Integration Points

The Run's runtime factory, the Agent Session owner's selection attempts in the
Run Database, the Run Event stream and the Doctor report. No test,
Verification command or QA row starts a real ACP runtime, reaches OpenRouter
or the network, or reads the real `~/.roundfix`.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary homes and repositories, and passes no real-time duration.

- `internal/lighttier/access_test.go` (new):
  `TestServingSelectionsKeepsOnlyRuntimesThatServeTheRequestedAccess`, with
  `agent.RuntimeFor` as the factory: under full access the `opencode`
  selections are all refused with `runtime_has_access_mode`, runtime
  `opencode`, policy `full-access` and the exact message of the adopted
  Backlog Entry (subtest `full-access`); under the runtime default all are
  kept (`runtime-default`); a factory error keeps the selection
  (`factory-error`); a fake factory that gives only the second model a
  full-access candidate keeps only that one (`second-model-serves`); an
  empty list returns nothing (`empty`); `Selections` returns the runtime,
  model prefix and empty effort of API Contract 1 (`selections`).
- `internal/daemon/light_tier_access_test.go` (new), over the existing Task
  cycle fixture, `lightPlanForTest` and a runner wrapper whose
  `PrepareSession` also returns `req.Runtime.ValidateRequestedAccessPolicy()`:
  - `TestLightTierSkipsWhenTheLightRuntimeCannotServeTheRequestedAccess`,
    subtests `full-access` (no `opencode` session prepared, one work request
    on `category-model`, one skip event with `access_mode_unavailable`, the
    exact warning line, no yield event), `full-access-key-missing` (the
    reason code is still `access_mode_unavailable`), `runtime-default` (the
    work request is on the first light model, no skip event) and
    `second-model-serves` (a factory that gives only the second light model a
    full-access candidate: the work request is on that model).
  - `TestLightTierYieldsToTheCategoryProfileBeforeAgentWork`, subtests
    `factory-refused`, `setup-refused` (an `*agent.AccessPolicyError` of kind
    `access_policy_rejected` from the first light model's setup) and
    `prompt-failed-before-work` (the first light prompt returns an error and
    publishes nothing): each completes with its last work request on
    `category-model`, prepares no second light model, records the light
    attempt `failed`, publishes one `light_tier_yielded` event with the light
    model, the reason code of API Contract 3 and the cause, one fallback
    notification and one warning line.
  - `TestLightTierDoesNotYieldAfterWorkOrOutsideTheLightTier`, subtests
    `after-work-started` (the light prompt publishes an Agent message and
    then returns an error: the Task fails, no yield event, no category
    session) and `category-refused` (the tier off and the category
    Preferred Selection's setup returning a plain error: the Task fails as
    today, no yield event).
- `internal/cli/implement_light_tier_access_test.go` (new):
  `TestImplementSkipsTheLightTierUnderProjectFullAccess`: `roundfix
  implement` with `defaults.agent_full_access: true` in `.roundfixrc.yml`, a
  light model in User Config, a fixture key and a `low` docs Task exits 0,
  sends its one work request to a runtime whose requested policy is
  `full-access` and whose access check passes, prints the skip warning of API
  Contract 2 on stderr, and prints no key value.
- `internal/cli/doctor_light_tier_test.go` (new):
  - `TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess`:
    `doctorLightTierResult` is not printed with full access off or an empty
    light model list (subtests `full access off`, `no light model`), is the
    exact `warn` line of API Contract 4 with the static factory (`static
    runtime`), and the exact `ok` line with a fake factory that gives
    `opencode` a full-access candidate (`serving runtime`).
  - `TestDoctorPrintsTheLightTierLineOnlyWithFullAccessAndLightModels`:
    through `roundfix doctor` with the fake checker and readiness of the
    `linux-lane` test, the built-in configuration prints no `light-tier:`
    line (subtest `built-in`, run first); with `defaults.agent_full_access` `true` it prints the `warn` line
    right after the `recommendations:` line and before `pre-pr-review:`,
    removing that line leaves the output byte-identical to the first, and the
    exit code is 0; with an empty light model list no line is printed.
- Unchanged and run by name: `TestTierFor`, `TestPlanEnabled`,
  `TestLightTierDispatchOrder`, `TestLightTierSkipsWithAWarning`,
  `TestLightTierRecordsSpend`, `TestLightTierEscalatesOnce`,
  `TestNoFallbackAfterAgentWorkStarted`,
  `TestAgentSelectionFallbackPublishesBeforeNextSession`,
  `TestImplementBuildsTheLightTierPlan`,
  `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`,
  `TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts`,
  `TestRunDoctorMissingRepositoryRoot` and `TestDoctorReportsPrePRReviewPolicy`.
- No existing test, golden or fixture changes.

## Build Order

1. The configuration and Doctor guides, the Roundfix `runtime` and `setup`
   references, the Roundfix skill's versions, the Behavior Surface Record and
   **Light Tier** (no dependency: the text states the contracts of
   `_techspec.md`, and its Verification reads only documents).
2. The helper, the access skip, the yield and their tests (depends on: 1,
   which writes the guide and the references for the surface this step
   changes, so the code and its text land in one chain).
3. The Doctor line and its tests (depends on: 2, which adds the helper the
   line calls).
4. QA gate (depends on: 1, 2, 3).

No Spec is listed in `requires`. The cycle AO Specs 0299 and 0300 deliver
first; Spec 0300 raises the Roundfix skill's version, and this Spec rebases on
it and touches none of its files apart from the skill's version fields.

## Risks & Considerations

- A prompt failure before any Agent output now runs the category profile
  instead of failing the Task, so an infrastructure fault of the light
  runtime costs one category session; the Task still commits only after its
  Verification passes.
- A future OpenCode adapter that names a full-access mode needs only a
  candidate in `agent.RuntimeFor`; the skip then stops applying and the
  advertised-mode check governs the session.
- The skip is decided from the runtime's declared candidates, not from a
  live session; a runtime that declares a candidate its adapter does not
  advertise still reaches the degraded report, as today.
- The Doctor line reads only configuration and static runtime capabilities,
  so it cannot see the `--agent-full-access` flag of a later Run.

## Glossary

- changes: **Light Tier**

## Decisions

- The access check is static, from the runtime factory, and runs first in
  `applyLightTier`; see ADR-0314.
- The yield reuses the `AGENT_WORK_STARTED` boundary and skips the remaining
  light models; selection start failures keep walking them; see ADR-0314.
- Doctor prints `light-tier:` only for the combination and never fails;
  `profiles validate` is unchanged; see ADR-0314.
