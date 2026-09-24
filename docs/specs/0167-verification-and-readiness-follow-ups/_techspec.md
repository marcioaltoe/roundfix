---
spec: 0167-verification-and-readiness-follow-ups
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# Verification and readiness follow-ups

## Executive Summary

Let the exclusive retry's verdicts replace the first run's, skip completed Tasks
in the verbatim planning check, and print a degraded access policy in text
readiness output.

## Project Constraints

- Identifier strategy: not applicable — no identifier changes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes and ACP adapters
  only; no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0038, ADR-0056, ADR-0080, ADR-0091, ADR-0093,
  ADR-0096, ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155, ADR-0156,
  ADR-0159 and ADR-0160 hold. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — empty intersection with `GovernedPath`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Retry verdicts

`retainCollectedVerificationFailures` in `internal/daemon/task_engine.go` keys
failures by command; the retry's outcome for every command it reached replaces
the first run's, and only first-run failures of commands it did not reach are
carried.

## Completed repairs

`ValidatePreconditionRepairs` in `internal/daemon/task_engine.go` skips named
Tasks whose status is completed.

## Degraded access

`printProfilesValidateSuccess` in `internal/cli/profiles_validate.go` and
Doctor's profile readiness line print the effective access policy when it is
degraded.

## API Contracts

1. The repair request lists each current failure once.
2. `implement` accepts a completed repair Task without the verbatim command.
3. Text readiness output names a degraded access policy.

## Coverage Map

- Goal 1 → Retry verdicts; API Contract 1.
- Goal 2 → Completed repairs; API Contract 2.
- Goal 3 → Degraded access; API Contract 3.
- Core Feature 1 → Retry verdicts.
- Core Feature 2 → Completed repairs.
- Core Feature 3 → Degraded access.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- API Contracts 1-3 → Retry verdicts, Completed repairs, Degraded access.

## Integration Points

- **Spec 0158.** Owns the behaviour this Spec corrects.

## Testing Approach

1. **Daemon.** A command failing in both runs appears once; a command passing on
   retry is absent; a completed repair Task without the verbatim command does not
   refuse planning, while a pending one still does.
2. **Readiness.** The text output of `profiles validate` and Doctor name a
   degraded policy, and print nothing extra when it is not degraded.
3. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Retry verdicts and completed repairs (depends on: none).
2. Degraded access (depends on: 1).
3. Terminal QA (depends on: 1, 2).

## Risks & Considerations

- **Losing a real failure.** Only commands the retry actually ran have their
  first-run verdict replaced.
