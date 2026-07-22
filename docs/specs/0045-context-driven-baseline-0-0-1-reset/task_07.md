---
task: task_07
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: backend
complexity: high
---

# Task 07: Resolve individual Readoption dispositions

## Overview

Turn the structural inventory into an explicit, reviewable Decision Plan.
Every source entry must receive a valid individual disposition, and any new
repository-specific rule must be proposed as exact bytes before setup may
write it.

## Requirements

1. MUST add `--decision-file <path>` support for the structured decision
   document defined by the TechSpec while retaining scalar `--decision` for
   decisions that remain scalar.
2. MUST require exactly one valid disposition for every incompatible Source
   Baseline Entry and reject missing, duplicate, unknown, or stale entry IDs.
3. MUST support only the typed dispositions and destination constraints
   declared by the TechSpec; no disposition may be inferred from source text.
4. MUST require exact proposed bytes, target path, and digest for content moved
   into a typed documentation destination.
5. MUST use `docs/agents/repository-rules.md` as the default home for
   Repository-Specific Normative Rules.
6. MUST require confirmation before the first Repository-Specific Normative
   Rules file is created, keep the file unmarked, and preserve it on later
   setup runs.
7. MUST expose stable rejection reasons and a deterministic plan preview while
   performing no writes.
8. MUST include the full normalized decision document in the plan digest
   input.

## Subtasks

- [ ] Define and parse the strict structured decision-file schema.
- [ ] Validate one-to-one coverage of the Source Baseline inventory.
- [ ] Implement typed dispositions and destination restrictions.
- [ ] Model exact proposed bytes and digests for documentation moves.
- [ ] Add first-write confirmation and preservation rules for the repository
      rules file.
- [ ] Add stale, incomplete, duplicate, and unsafe decision mutations.
- [ ] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [ ] A complete decision file maps every inventory entry exactly once and
      produces a stable preview without changing repository bytes.
- [ ] Missing, duplicate, unknown, stale, or structurally invalid dispositions
      fail with an entry-specific reason.
- [ ] Documentation moves expose their exact proposed bytes, target, and digest
      before confirmation and reject unsafe or untyped destinations.
- [ ] The first repository-rules write requires explicit confirmation; an
      existing file is treated as repository-owned and preserved.
- [ ] Changing any normalized disposition or proposed target changes the plan
      digest.
- [ ] No source entry receives an automatic semantic classification.

## Context

- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- instruction: `docs/adr/0063-repositories-own-their-http-contract.md`
- instruction: `docs/adr/0064-baseline-readoption-uses-byte-exhaustive-structural-inventory.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_decision_rendering.py`
- interface: `.agents/skills/setup-context-driven/tests/test_preview.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_readoption_decisions.py'` — expected: complete decisions preview deterministically and every incomplete, stale, duplicate, or unsafe mapping fails with the asserted reason.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_preview.py'` — expected: structured decisions affect the plan digest without producing writes.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Core Features 5, 7, 8, and 13; User Stories 1, 3, and 5.
- `_techspec.md` → Implementation Design: Interfaces, Data Models, and API
  Contracts; Build Order 7.
- ADR-0047 → declared decision effects and digest inputs.
- ADR-0063 → durable repository-owned HTTP policy.
- ADR-0064 → individual explicit Readoption dispositions.
