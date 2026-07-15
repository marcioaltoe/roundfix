---
task: task_02
spec: 0030-context-driven-agent-instructions
status: pending
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

- [ ] Implement asset, manifest, and repository-state loaders with strict validation.
- [ ] Implement deterministic profile resolution and ownership-marker parsing.
- [ ] Implement document, reference, language, and template-version audits.
- [ ] Implement the finding model, severity aggregation, text renderer, JSON renderer, and exit mapping.
- [ ] Add behavior-focused micro tests and temporary-repository CLI tests.
- [ ] Prove audit leaves every inspected repository byte unchanged.

## Acceptance Criteria

- [ ] A compliant temporary repository returns exit code `0` in text and JSON modes.
- [ ] Blocking findings return exit code `1`, malformed inputs return `2`, and unresolved decisions return `3`.
- [ ] JSON results conform to `setup-context-driven/audit-v1` and contain stable codes, paths, messages, and actions.
- [ ] Marker errors, missing guides, broken references, stale templates, and non-English generated content are independently detectable.
- [ ] Audit does not create, modify, rename, or delete any target-repository file.
- [ ] Running the audit twice against unchanged input produces identical semantic output.

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
