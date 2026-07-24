---
task: task_04
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
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

- [x] Add immutable capability, evidence, strength, and outcome values.
- [x] Implement bounded file, executable, and installed-skill evidence
      adapters.
- [x] Define universal required and recommended capability rules.
- [x] Add stable blocking diagnostics and advisory warnings.
- [x] Add side-effect guards and evidence-strength mutation tests.
- [x] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [x] Missing Context7 or Exa evidence produces a blocking outcome naming the
      missing capability and next action.
- [x] Missing Firecrawl, `rtk`, or `rg` produces a non-blocking warning whose
      explanation is available in the result.
- [x] Stronger evidence satisfies a compatible weaker requirement, while
      insufficient evidence strength does not.
- [x] Equivalent local evidence produces byte-identical ordered output.
- [x] Capability evaluation performs no writes, installs, network access, or
      repository-script execution.
- [x] Rendered guidance states that external research cannot replace local
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

## Result

Added a frozen Repository Capability contract and deterministic evaluator with
bounded adapters for declared repository files, executable discovery, and
installed Repository Skill Set membership. The universal rules require
Context7 and Exa, recommend Firecrawl, `rtk`, and `rg`, and return stable JSON
and text outcomes without connecting to a service or executing repository
code. The canonical and distributed setup skill trees contain identical
implementation and test files.

Verification:

- Pre-change import probe failed with `ModuleNotFoundError` for
  `context_capabilities`, establishing the missing behavior.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s
  .agents/skills/setup-context-driven/tests -p 'test_capabilities.py'` — passed
  9 tests after the final implementation edit.
- The same focused command against `skills/setup-context-driven/tests` —
  passed 9 tests.
- `rtk make skills-sync-check`, direct comparisons of both new canonical and
  distributed files, and `rtk git diff --check` — passed.
- `rtk make verify` — passed on the unchanged elevated rerun: both 195-test
  setup suites passed, 1,694 Go tests passed, asset loading passed, the
  Roundfix skill check passed, and the CLI built. The sandboxed attempt was
  blocked only by access to the host Go build cache.

Acceptance evidence:

- `test_missing_required_research_capabilities_block_with_next_actions`
  proves Context7 and Exa gaps use `capability.required.missing`, block
  readiness, name the capability, and provide a Repository Skill Set action.
- `test_missing_recommended_capabilities_warn_without_blocking` proves the
  Firecrawl, `rtk`, and `rg` gaps are explanatory, non-blocking warnings.
- The two evidence-strength tests prove stronger compatible evidence satisfies
  a weaker minimum, while weak or incompatible evidence does not.
- `test_equivalent_evidence_renders_byte_identical_ordered_output` proves
  equivalent reordered inputs produce identical JSON bytes and human guidance.
- The adapter test proves file, executable, and installed-skill evidence; the
  side-effect test disables writes, subprocesses, network calls, and system
  commands while preserving every repository byte.
- The guidance test proves local `rg` search appears before Context7, Exa, and
  Firecrawl and explicitly says external research cannot replace local code
  search.

Follow-ups: Task 05 owns binding these universal rules to the new profile, and
Task 08 owns audit and Change Plan integration.
