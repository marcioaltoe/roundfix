---
status: accepted
created_at: 2026-10-22T21:30:00Z
updated_at: 2026-10-22T21:30:00Z
deprecated_at: null
superseded_by: null
---

# A Light Tier that cannot serve the Run yields to the category profile

ADR-0238 gave every `low` Task a light profile on `opencode` and said that a
light session that fails to start falls to the next candidate and finally to
the category's Preferred Selection. Only a selection start failure did. On
2026-10-21 a Fluxus Run with `defaults.agent_full_access: true` failed every
light-eligible Task before any Agent work: OpenCode names no full-access mode,
so the light session's access-policy check refused it
(`runtime_has_access_mode`), and that refusal is not a selection start failure.
ADR-0293 keeps a runtime without full-access candidates refusing a full-access
request, and ADR-0011 forbids silently running a sandboxed session instead. On
2026-10-22 the maintainer chose "Pular no planejamento + escalar
(Recomendado)".

## Decision

**A light runtime that cannot serve the requested access is not chosen.**
When the Daemon derives a light Task's profile, it builds each light
selection's runtime with the Run's runtime factory and keeps only those whose
runtime offers a mode for the requested access policy. When none remains, the
Task skips the Light Tier: it prints `roundfix: warning: light tier skipped for
Task <id>: light runtime <runtime> offers no mode for access policy
"<policy>"; it runs on its <category> profile` and publishes the Task event
phase `light_tier_skipped` with the reason code `access_mode_unavailable`.
This check runs before the key, spend log and ceiling checks, whose reason
codes and messages do not change. Under the runtime-default access policy
every light selection is kept.

**Any other refusal of a light session before Agent work yields to the
category profile.** While a light selection is active and `AGENT_WORK_STARTED`
has not been published, a failure that is not a selection failure, a lost
rollout, a provider system error or a cancellation (the runtime factory's
error, a setup error such as an access-policy refusal, a prompt that fails
before any Agent output) no longer fails the Task. Roundfix records the light
attempt as failed, prints `roundfix: warning: light tier yielded for Task
<id> before Agent work: <cause>; it runs on its <category> profile`,
publishes the Task event phase `light_tier_yielded` with the light model, a
reason code and the cause, publishes the Agent Selection fallback
notification, and opens the category's Preferred Selection, skipping the
remaining light models. Selection start failures keep falling to the next
light model as before. After Agent work begins, the only move off the light
model stays ADR-0238's single escalation after a failed Verification.

**`roundfix doctor` warns before a Run finds it.** When the effective light
model list is not empty and `defaults.agent_full_access` is `true`, Doctor
prints a `light-tier:` line after `recommendations:`. It is `warn` with
`DR-LIGHT-ACCESS-MODE` when no light selection's runtime offers a full-access
mode and `ok` otherwise; it never changes the exit code and opens no Agent
Session. `roundfix profiles validate` is unchanged: it proves configured
profiles, and the light selections are not profiles.

## Considered Options

- Disabling full access for the light session was rejected: ADR-0011 forbids
  running a sandboxed session when full access was requested, and the
  Verification that needed full access would fail later.
- Treating every pre-work light failure as a selection start failure that
  walks the remaining light models was rejected for runtime-level refusals:
  every light model runs on the same runtime, so each would fail the same way.
- Refusing the Run at Preflight when the combination is configured was
  rejected: the category profile can run the Task, and the maintainer chose
  to skip and escalate rather than block.
- Proving the light models in `profiles validate` was rejected: it would open
  OpenRouter sessions and spend on every validation, and the light selections
  are derived at dispatch, not configured profiles.
