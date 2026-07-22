---
task: task_04
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: backend
complexity: high
---

# Task 04: Evaluate local Repository Capability evidence

## Overview

Introduce a generic, read-only capability evaluator so profiles can require or
recommend tools without hard-coded setup branches. Evidence must come only
from bounded local inspection and must explain both blocking and advisory
outcomes.

## Requirements

1. MUST implement the `context_capabilities.py` Repository Capability schema,
   evidence adapters, strengths, statuses, and diagnostics defined by the
   TechSpec.
2. MUST support evidence from declared local files, executable discovery, and
   installed Repository Skill Set membership.
3. MUST evaluate Context7 and Exa as required capabilities for every profile.
4. MUST evaluate Firecrawl, `rtk`, and `rg` as recommended capabilities that
   warn without blocking setup.
5. MUST make missing required capabilities blocking with a concise explanation
   and useful next action.
6. MUST keep external research subordinate to local code search in rendered
   capability guidance.
7. MUST NOT install tools, access the network, execute repository scripts, or
   mutate the inspected repository while evaluating evidence.
8. MUST produce deterministic machine-readable and human-readable outcomes.

## Subtasks

- [ ] Add immutable capability, evidence, strength, and outcome values.
- [ ] Implement bounded file, executable, and installed-skill evidence
      adapters.
- [ ] Define universal required and recommended capability rules.
- [ ] Add stable blocking diagnostics and advisory warnings.
- [ ] Add side-effect guards and evidence-strength mutation tests.
- [ ] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [ ] Missing Context7 or Exa evidence produces a blocking outcome naming the
      missing capability and next action.
- [ ] Missing Firecrawl, `rtk`, or `rg` produces a non-blocking warning whose
      explanation is available in the result.
- [ ] Stronger evidence satisfies a compatible weaker requirement, while
      insufficient evidence strength does not.
- [ ] Equivalent local evidence produces byte-identical ordered output.
- [ ] Capability evaluation performs no writes, installs, network access, or
      repository-script execution.
- [ ] Rendered guidance states that external research cannot replace local
      code search.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/adr/0061-standard-typescript-monorepo-is-opinionated.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/tests/test_assets.py`
- interface: `.agents/skills/setup-context-driven/tests/test_support.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_capabilities.py'` — expected: required evidence blocks when absent, recommended evidence warns, ordering is stable, and side-effect guards remain untouched.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Core Features 10 and 11; User Story 2; User Experience.
- `_techspec.md` → System Architecture; Implementation Design: Interfaces and
  Data Models; Build Order 4.
- ADR-0061 → universal and profile-specific capability policy.
