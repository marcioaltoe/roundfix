---
task: task_03
spec: 0233-a-light-tier-on-open-models
status: pending
type: backend
complexity: high
---

# Task 03: Light Tasks dispatch on an open model with the implementation key, record their spend and escalate once to their profile

## Overview

With the tier rule, the User Config keys and the Light Spend Log in place,
dispatch still gives every Task its category's profile. This Task makes the
daemon run a light Task on a derived Agent Selection Profile whose first
candidates are the machine's light models on `opencode` with no effort, makes
the ACPX runner give that session the implementation key by name, records each
light prompt's cost, skips the tier with a warning when the key is missing,
the log is unreadable or the month is at its ceiling, and moves the one
Verification Feedback turn of a failed light Task to the category's profile.
`roundfix implement` builds the plan from User Config and the Run's
environment. It is verifiable through the fake acpx harness, the Task engine
harness and the `implement` command with a fake runner.

## Requirements

1. MUST add to the ACP Runtime description the key variable a light session
   reads, and, only when it is set, MUST give every acpx command of that
   session an `OPENCODE_CONFIG_CONTENT` that merges
   `provider.openrouter.options.apiKey` = `{env:<variable>}` into the inherited
   JSON object (or a new object when none is inherited) and MUST remove
   `OPENROUTER_API_KEY` from that command's environment (TechSpec Invariant
   6); an inherited value that is not a JSON object MUST fail the session
   before any prompt as a selection failure (Invariant 7). A runtime without
   the variable MUST keep today's environment byte for byte.
2. MUST build a light Task's owner from the derived profile of Invariant 4,
   with profile source `light-tier` (API Contract 6), when the Task's tier is
   `light` and the Run's light plan is enabled; every other Task, every QA gate
   and every review Batch MUST keep today's profile and sessions.
3. MUST, before building a light Task's owner, apply the skip order of
   Invariant 5 and, on a skip, use the category's profile unchanged, print API
   Contract 4's warning on the Run's progress stream and publish its Task Run
   Event with the reason code.
4. MUST, after every prompt of a light session, append one Light Spend Log
   line with the Run, repository, Spec, Task, session and model and the
   increase of the session's reported USD cost since the previous line for
   that session, or 0 with `cost_source` `unreported` when no USD cost was
   reported (Invariant 8); a failed append MUST print a warning and never
   fail the Task.
5. MUST, when a light Task's first Verification fails with a command failure,
   close the light session and run the one Verification Feedback turn on the
   first non-light candidate in a new Agent Session, with that candidate's own
   start failures falling through the rest of the chain, and with a prompt
   that carries the Task prompt, the diagnostics and the statement that the
   working tree holds another model's attempt; MUST print and publish API
   Contract 5; MUST NOT escalate a second time, and MUST leave the retry
   budget, the commit rule and every non-light repair unchanged (Invariant 9).
6. MUST make `roundfix implement` build the Run's light plan from the loaded
   User Config, the Run's environment through `OpenRouterImplementKey`, Roundfix
   Home and the repository, and MUST leave it disabled when the Run has a
   one-Run Agent Selection override.
7. MUST never write the key value to a file, a log, a Run Event, a progress
   line or the Light Spend Log (Invariant 10).
8. MUST add tests: through the fake acpx harness, a light runtime's child
   environment holds the merged inline option naming the variable, holds
   neither the key value nor `OPENROUTER_API_KEY`, and an inherited non-object
   fails before any prompt, while the session sends only the work prompt and
   no `Session setup.` prompt; through the Task engine harness with a fake
   runner and a temporary Home, the dispatch order for a light Task and the
   unchanged dispatch for `medium`, `qa` and Governed-Path Tasks, each skip
   reason with its warning and Run Event, a light start failure taken by the
   category's Preferred Selection before Agent work, spend lines whose costs
   add up to the last reported reading, and one escalation in a new session
   followed by a failed Task after a second failure; and through `implement`
   with a fake runner and a disposable Home, Surface Transcript 3 and an
   `--agent` override that keeps every Task on the override.
9. MUST NOT change the judge, the configuration keys, any guide or skill, any
   built-in or Recommended Profile or the Run Database schema; no test may
   reach the network or the real `~/.roundfix`.

## Subtasks

- [ ] Give a light session the key by name in the runner.
- [ ] Derive the light profile, apply the skips and record spend in the
      daemon.
- [ ] Escalate a failed light Task once.
- [ ] Build the light plan in `implement`.
- [ ] Add the runner, daemon and command tests.

## Acceptance Criteria

- [ ] A light Task's first selection is the first light model with no effort,
      and its session never sends a warm-up prompt.
- [ ] A missing key, an unreadable log and a month at its ceiling each run the
      Task on its profile with one warning.
- [ ] Every light prompt leaves one Light Spend Log line and no key value
      anywhere.
- [ ] A failed light Task gets one repair on its profile in a new session.

## Context

- instruction: `docs/adr/0238-a-light-tier-runs-low-complexity-tasks-on-an-open-model.md`
- instruction: `docs/adr/0114-opening-an-agent-session-is-not-agent-work.md`
- instruction: `docs/user-guide/configuration.md`
- instruction: `.agents/skills/roundfix/references/runtime.md`
- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/implement.go`
- creates: `internal/agent/light_session_test.go`
- creates: `internal/daemon/light_tier_test.go`
- creates: `internal/cli/implement_light_tier_test.go`

## Verification

- `grep -q 'light tier skipped for Task' internal/daemon/agent_session_owner.go && grep -q 'light_tier_escalated' internal/daemon/task_engine.go && go build -buildvcs=false ./... && go vet ./internal/agent ./internal/daemon ./internal/cli` — expected: exit 0; before this Task the daemon prints no light tier warning and publishes no escalation, so the command fails.
- `ag="$(go test -count=1 -v -run '^(TestLightSessionReadsTheImplementKeyByName|TestLightSessionSendsNoWarmup)$' ./internal/agent 2>&1)" || { printf '%s\n' "$ag"; exit 1; }; dm="$(go test -count=1 -v -run '^(TestLightTierDispatchOrder|TestLightTierSkipsWithAWarning|TestLightTierRecordsSpend|TestLightTierEscalatesOnce)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$dm"; exit 1; }; cl="$(go test -count=1 -v -run '^TestImplementBuildsTheLightTierPlan$' ./internal/cli 2>&1)" || { printf '%s\n' "$cl"; exit 1; }; for name in TestLightSessionReadsTheImplementKeyByName TestLightSessionSendsNoWarmup TestLightTierDispatchOrder TestLightTierSkipsWithAWarning TestLightTierRecordsSpend TestLightTierEscalatesOnce TestImplementBuildsTheLightTierPlan; do printf '%s\n%s\n%s\n' "$ag" "$dm" "$cl" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the seven tests exists, so the command fails.

## References

- `_prd.md` → Goals; User Stories 1, 4 and 5; Core Features 1, 2, 3, 4, 5 and 6; Success Metrics 1, 2, 3 and 4
- `_techspec.md` → Interfaces; Invariants; Data Models; API Contract 4; API Contract 5; API Contract 6; Surface Transcript 3; Testing Approach; Build Order 3
- ADR-0238; ADR-0114; ADR-0050; ADR-0108
