---
task: task_07
spec: 0044-upgrade-retention-and-formatter-compatibility
status: pending
type: backend
complexity: medium
---

# Task 07: Report uncovered instruction delegation

## Overview

Surface repository-authored documents that delegate categories the selected
Context-Driven Baseline does not cover. The signal must state that the
baseline is a floor while remaining deterministic, read-only, and
non-blocking.

## Requirements

1. MUST scan only root and nested `AGENTS.md` and `CLAUDE.md` repository
   instruction documents outside setup-managed marker spans.
2. MUST require both an explicit delegation signal and a declared coverage
   alias before reporting a category.
3. MUST compare detected categories with the active profile's selected clause
   coverage.
4. MUST emit one deterministic informational finding per document and missing
   category, naming the affected document and the baseline-floor action.
5. MUST exclude VCS, dependency, vendor, and setup skill mirror trees; avoid
   symlink traversal; and enforce declared file-count and byte limits.
6. MUST never alter Change Plan authority, exit status, or repository bytes.
7. MUST keep canonical and distributed setup skill trees synchronized.

## Subtasks

- [ ] Declare bounded aliases for delegable coverage categories.
- [ ] Discover eligible repository instruction documents safely.
- [ ] Match delegation paragraphs outside managed spans.
- [ ] Compare matched categories with selected profile coverage.
- [ ] Add covered, uncovered, excluded, symlink, and limit fixtures.
- [ ] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [ ] A nested instruction document delegating an uncovered category emits
      `delegation.baseline-floor` with informational severity and its path.
- [ ] A document with the same category but no delegation signal emits no
      finding.
- [ ] A fully covered delegation emits no floor finding.
- [ ] Text inside setup markers, ignored trees, and symlink targets is not
      treated as repository-authored delegation.
- [ ] Multiple findings have deterministic document/category ordering and
      duplicates collapse to one finding.
- [ ] Scan limits fail safely without writes, process execution, or network
      access.
- [ ] Informational findings do not block preview, apply, or a clean exit.
- [ ] Canonical and distributed setup skill trees are byte-identical.

## Context

- interface: `.agents/skills/setup-context-driven/assets/coverage.json`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_audit.py`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_delegation.py'` — expected: only eligible uncovered delegation emits deterministic non-blocking floor findings.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 5; User Story 6; Core Feature 8; User Experience; Success
  Metrics.
- `_techspec.md` → System Architecture; Data Models; Testing Approach; Risks &
  Considerations; Build Order 4.
