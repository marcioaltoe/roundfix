---
type: fix
status: promoted
created: 2026-10-21
spec: 0301-the-light-tier-yields-to-the-category-profile
reason: null
---

# The Light Tier fails a Task under agent full access

## Symptom

In a repository whose Project Config sets `defaults.agent_full_access: true`,
with the Light Tier on (the default `openrouter.light_models`), every
light-eligible Task (complexity `low`, type not `qa`, no Governed Path) fails
in about 15 seconds, before any Agent work, with:

```text
Task task_01 failed: Agent failed: requested access policy "full-access" failed predicate "runtime_has_access_mode" for runtime "opencode"; recovery: disable full access or select a runtime that supports it
```

Its dependents are skipped, the QA gate is held, and the Run ends Unresolved.
The category's own profile, which offers a full-access mode, is never tried.

## Where

The Daemon's Agent Session owner in `internal/daemon/agent_session_owner.go`:
`applyLightTier` derives the light profile (`opencode` with
`openrouter/<model>`, then the category's Preferred Selection and Fallback
Chain) and skips the tier only for `key_missing`, `spend_unreadable` and
`ceiling_reached`; `activate` falls to the next candidate only for a selection
start failure, so the access-policy refusal of the light session
(`agent.AccessPolicyError`, predicate `runtime_has_access_mode`, from
`internal/agent/agent.go` `ValidateRequestedAccessPolicy`) fails the Task.
The OpenCode runtime names no full-access mode candidate.

## Expected

A light Task never fails because the light runtime cannot serve the requested
access policy: the tier is not chosen for it, with a warning, and the Task
runs on its category's profile. Any other refusal of the light session before
Agent work also moves the Task to its category's profile instead of failing
it. `roundfix doctor` warns about the combination before a Run finds it.

## Evidence

Report from Fluxus (ai-memory message of 2026-10-21, 20:26 BRT, roundfix
0.104.0): Run `run_20261021T223522Z_9aa628a4520889af` of Fluxus Spec 0113 failed
every light-eligible Task with the message above, with
`defaults.agent_full_access: true` in Fluxus's Project Config and the Light
Tier on. The workaround in force on both of the maintainer's machines is
`openrouter.light_models: []` in User Config, which turns the Light Tier off
everywhere.
