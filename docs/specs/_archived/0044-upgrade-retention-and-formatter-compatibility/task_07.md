---
task: task_07
spec: 0044-upgrade-retention-and-formatter-compatibility
status: completed
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

- [x] Declare bounded aliases for delegable coverage categories.
- [x] Discover eligible repository instruction documents safely.
- [x] Match delegation paragraphs outside managed spans.
- [x] Compare matched categories with selected profile coverage.
- [x] Add covered, uncovered, excluded, symlink, and limit fixtures.
- [x] Synchronize the distributed setup skill copy.

## Acceptance Criteria

- [x] A nested instruction document delegating an uncovered category emits
      `delegation.baseline-floor` with informational severity and its path.
- [x] A document with the same category but no delegation signal emits no
      finding.
- [x] A fully covered delegation emits no floor finding.
- [x] Text inside setup markers, ignored trees, and symlink targets is not
      treated as repository-authored delegation.
- [x] Multiple findings have deterministic document/category ordering and
      duplicates collapse to one finding.
- [x] Scan limits fail safely without writes, process execution, or network
      access.
- [x] Informational findings do not block preview, apply, or a clean exit.
- [x] Canonical and distributed setup skill trees are byte-identical.

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

## Result

Implemented a bounded, read-only repository instruction scan that requires an
explicit delegation verb, a root or setup-guidance target, and a declared
coverage alias. It compares matched categories with the resolved active
profile modules and emits deduplicated `delegation.baseline-floor` information
without adding Change Plan authority or changing exit behavior. Limit and read
failures stop the scan without partial floor findings.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_delegation.py'` — passed, 6 tests.
- The same focused command against `skills/setup-context-driven/tests` — passed, 6 tests.
- Focused existing `test_audit.py`, `test_upgrade_contracts.py`, and
  `test_macro_profiles.py` suites — passed, 25 tests total.
- `rtk make skills-sync-check` — passed; canonical and distributed trees are
  byte-identical.
- `rtk make verify` — passed after rerunning with access to the host Go build
  cache; 1,694 Go tests, 168 canonical Python tests, 168 distributed Python
  tests, asset validation, skill checks, and the build passed. The first
  sandboxed attempt could not read the host Go build cache and exited before
  repository verification completed.

Acceptance evidence:

- `test_uncovered_nested_delegation_emits_one_informational_finding` proves the
  exact code, informational severity, nested path, category, deduplication, and
  read-only repository bytes.
- `test_alias_without_delegation_signal_and_covered_category_emit_no_floor`
  proves that an alias alone is insufficient and active Rust coverage emits no
  floor finding.
- `test_managed_ignored_and_symlinked_instructions_are_excluded` proves managed
  spans, VCS, dependency, vendor, both setup-skill trees, and symlink targets
  are excluded.
- `test_findings_sort_by_document_and_category_and_collapse_duplicates` proves
  deterministic document/category order and one finding per pair.
- `test_scan_limits_stop_without_partial_floor_findings_or_writes` proves both
  declared limits stop safely with one informational scan finding and no
  partial category findings or writes.
- `test_floor_finding_does_not_change_plan_authority_or_clean_exit` proves the
  plan digest stays identical, preview and apply remain non-blocking, apply
  exits zero, and repository bytes stay unchanged.
- The canonical and distributed focused suites plus `skills-sync-check` prove
  the shipped skill copies remain synchronized.

Follow-ups: none.
