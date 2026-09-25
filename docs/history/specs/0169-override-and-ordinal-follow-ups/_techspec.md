---
spec: 0169-override-and-ordinal-follow-ups
status: active
created: 2026-09-25
surfaces: [backend, cli, docs]
---

# Override and ordinal follow-ups

## Executive Summary

Stamp the QA Task status an override waives, bring the two skills in line with
the command, and make ordinal claims unique within one Spec with a visible skip.

## Project Constraints

- Identifier strategy: applicable — claims unique within a Spec. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0154, ADR-0155 and ADR-0156 hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/spec/archive.go`, `internal/spec/archive_test.go`,
  `.agents/skills/archive-spec/SKILL.md`, `skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The stamp and the guidance

`Archive` in `internal/spec/archive.go` adds `qa_override_qa_task_status` when
the QA Task is not completed. The Roundfix skill's override paragraph and the
archive-spec skill's Steps state the current refusal rule and the command as the
only path; `docs/user-guide/commands.md` names the new field.

## The ordinal check

`internal/speccheck/ordinal.go` reports two claims of one Spec that share a
number with different paths. `internal/speccheck/citations.go` adds
`SC-ORDINAL-CLAIMED` to the detectors skipped when a Spec has no Task Graph.

## API Contracts

1. An override record carries `qa_override_qa_task_status` when the QA Task is
   not completed.
2. `SC-ORDINAL-CLAIMED` reports a same-Spec duplicate number.

## Coverage Map

- Goal 1 → The stamp and the guidance; API Contract 1.
- Goal 2 → The stamp and the guidance.
- Goal 3 → The ordinal check; API Contract 2.
- Core Features 1-2 → The stamp and the guidance.
- Core Feature 3 → The ordinal check.
- Success Metrics 1-2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 2.
- API Contracts 1-2 → The stamp and the guidance, The ordinal check.

## Integration Points

- **Spec 0159.** Owns the behaviour this Spec completes.

## Testing Approach

1. **Stamp and guidance.** A failed QA Task with a `pass` report records its
   status; a completed QA Task records none; both skills and the guide carry the
   current rule.
2. **Ordinals.** A same-Spec duplicate is reported; distinct numbers pass; a
   Spec without `_tasks.md` lists the skip.
3. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The stamp and the guidance (depends on: none).
2. The ordinal check (depends on: 1).
3. Terminal QA (depends on: 1, 2, 4).
4. One instruction for overrides, one finding per conflict (depends on: 2).

## Risks & Considerations

- **A stale skill.** The guidance change ships in the same Task as the stamp,
  so the skill cannot lag the code.
