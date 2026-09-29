---
type: fix
status: promoted
created: 2026-09-29
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
reason: null
---

# The Pull Request row blocks every qualifying partial

## Symptom

The authored QA gate always runs before a Pull Request exists (ADR-0088). The `qa-gate` skill and the Daemon's QA prompt tell the Agent to record every Pull Request journey as `blocked (environment: no open Pull Request)`. A qualifying `partial` requires `rows_blocked_environment: 0`. A Spec with a real Unreachable Acceptance declaration therefore closes `partial` and can never qualify: `rows_blocked_environment is 1; expected 0`. Spec 0179 needed a boilerplate fourth declaration covering the Pull Request row. It then also needed a new QA Requirement forcing that row to `blocked (declared: …)`, and two extra QA Runs.

## Where

- `QAReportEligibility` in `internal/spec/qa.go`.
- The count cross-check in `internal/speccheck/mechanical.go`.
- The Pull Request line in `internal/agent/spec_prompt.go`.
- The Pull Request journey rule in `skills/qa-gate/SKILL.md`.

## Expected

The row that no pre-PR gate can observe does not decide whether a `partial` qualifies. A Spec with genuine declared-unreachable rows and a Pull Request row backed by equivalent evidence reaches a qualifying `partial` without a boilerplate declaration or an extra Requirement.

## Evidence

- `docs/history/specs/0179-review-findings-with-evidence-and-no-unselected-providers/qa/qa-report-2026-09-29.md` (`rows_blocked_environment: 1`) against `qa-report-2026-09-29-01.md` (`0`).
- The same Spec's `_prd.md` fourth Unreachable Acceptance declaration and `task_05.md` Requirement 11.
- `docs/history/specs/0180-a-prepared-queue-that-revalidates-before-each-spec/qa/qa-report-2026-09-28.md`, row R11.
