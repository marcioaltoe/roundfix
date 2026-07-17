---
task: task_02
spec: 0030-context-driven-agent-instructions
status: completed
type: backend
complexity: high
---

# Task 02: Build the read-only agent-instruction audit engine

## Overview

Deliver the safe default path: a standard-library Python audit that reads bundled assets and repository state without mutating anything. The slice is observable through stable text/JSON results and exit codes against real temporary repositories.

## Requirements

1. MUST implement the `audit` behavior as the default when no subcommand is supplied and guarantee that it performs no writes.
2. MUST validate the Setup Manifest, selected profile and modules, ownership markers, managed blocks, supporting guides, internal references, generated language policy, and asset versions.
3. MUST emit stable finding codes and classify errors, decisions, warnings, and informational findings according to the TechSpec.
4. MUST support concise text output and schema-versioned JSON output with stdout reserved for results and stderr reserved for diagnostics.
5. MUST implement the documented exit-code contract for compliant state, blocking findings, invalid input, and unresolved decisions.
6. MUST identify unbalanced, nested, duplicate, missing, or stale managed blocks without matching their prose.
7. MUST use Python 3 standard library only and resolve all paths relative to an explicit repository root or the current working directory.

## Subtasks

- [x] Implement asset, manifest, and repository-state loaders with strict validation.
- [x] Implement deterministic profile resolution and ownership-marker parsing.
- [x] Implement document, reference, language, and template-version audits.
- [x] Implement the finding model, severity aggregation, text renderer, JSON renderer, and exit mapping.
- [x] Add behavior-focused micro tests and temporary-repository CLI tests.
- [x] Prove audit leaves every inspected repository byte unchanged.

## Acceptance Criteria

- [x] A compliant temporary repository returns exit code `0` in text and JSON modes.
- [x] Blocking findings return exit code `1`, malformed inputs return `2`, and unresolved decisions return `3`.
- [x] JSON results conform to `setup-context-driven/audit-v1` and contain stable codes, paths, messages, and actions.
- [x] Marker errors, missing guides, broken references, stale templates, and non-English generated content are independently detectable.
- [x] Audit does not create, modify, rename, or delete any target-repository file.
- [x] Running the audit twice against unchanged input produces identical semantic output.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `docs/specs/0030-context-driven-agent-instructions/_techspec.md`

## Verification

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'` — expected: audit contracts, finding codes, output separation, exit codes, and read-only behavior all pass.
- `rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py --help` — expected: concise help documents audit, apply, sync-setups, output formats, and exit behavior.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → User Stories 2–4; Core Features 3–8, 13–14; User Experience.
- `_techspec.md` → Interfaces; Data Models; API Contracts; Testing Approach; Build Order 2.
- ADR-0046.

## Result

Implemented the read-only setup-context-driven audit engine:

- Added `scripts/context_setup.py` with default `audit` behavior, text/JSON renderers, finding aggregation, and the documented exit-code contract.
- Added Setup Manifest validation, profile/module resolution through bundled assets, required-decision checks, managed marker parsing, managed block/guide validation, internal reference checks, generated-language checks, and stale template detection.
- Added behavior-focused temporary-repository tests in `tests/test_audit.py`.
- Mirrored the new script and tests into `skills/setup-context-driven/` so `skills-sync-check` stays aligned.

Verification:

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'` — passed, 6 tests.
- `rtk python3 .agents/skills/setup-context-driven/scripts/context_setup.py --help` — passed and documents audit, apply, sync-setups, output formats, and exit behavior.
- `rtk git diff --check` — passed.
- `rtk make verify` — passed after removing generated Python `__pycache__` directories that caused a skill-sync drift check failure.

Acceptance evidence:

- Compliant repository: `test_compliant_repository_returns_zero_in_text_and_json_modes` asserts exit code `0`, empty stderr, text output, JSON schema `setup-context-driven/audit-v1`, and an empty findings list.
- Exit codes: `test_blocking_malformed_and_decision_exit_codes` asserts `1` for `manifest.missing`, `2` for malformed JSON, and `3` for `decision.required`.
- JSON fields: `test_json_finding_shape_contains_stable_fields` checks `code`, `severity`, `path`, `managedId`, `message`, and `action`.
- Independent findings: `test_marker_document_and_template_findings_are_independent` covers `managed.marker.invalid`, `docs.guide.missing`, `docs.reference.broken`, `managed.template.stale`, and `docs.language.non-english`.
- Read-only behavior: `test_default_subcommand_is_read_only_audit` snapshots every file before audit and verifies the snapshot is unchanged afterward.
- Determinism: `test_audit_twice_produces_identical_semantic_output` compares JSON output from two unchanged audit runs.

Follow-up for later Tasks:

- Implement safe apply and migration behavior in Task 03.
- Implement installed-skill classification and setup snapshot drift checks in Task 04.
