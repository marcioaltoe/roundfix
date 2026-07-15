---
task: task_03
spec: 0030-context-driven-agent-instructions
status: pending
type: backend
complexity: high
---

# Task 03: Add safe apply and durable decision migrations

## Overview

Add the explicit mutation path that turns a confirmed manifest and module composition into managed repository guidance. The slice proves atomic, surgical, and idempotent updates while retaining durable answers and refusing unsafe marker or migration states.

## Requirements

1. MUST implement `apply` as an explicit operation that updates only Setup Manifest fields and declared managed Markdown boundaries.
2. MUST preserve repository-authored bytes outside managed markers, including custom sections inside shared files but outside owned blocks.
3. MUST build and validate the complete change plan in memory before writing any destination.
4. MUST use temporary sibling files and atomic replacement so a failed multi-file apply leaves the prior repository state intact.
5. MUST store selected profile, modules, decision values, confirmation dates, managed artifact identities, template versions, and generated digests in the Setup Manifest.
6. MUST reuse compatible decisions without prompting and return `decision.required` only for missing, incompatible, or newly introduced decisions.
7. MUST require a one-time stored adoption decision before taking ownership of existing unmarked content.
8. MUST remove obsolete setup-owned blocks and guides only when their ownership is proven by the manifest and markers; it MUST NOT remove repository-authored artifacts.

## Subtasks

- [ ] Implement the Setup Manifest writer and version migration boundary.
- [ ] Implement change planning for managed blocks and fully managed guides.
- [ ] Implement adoption, stale-template replacement, and obsolete-managed-artifact handling.
- [ ] Implement atomic multi-file application and failure rollback behavior.
- [ ] Add macro tests for preservation, failure atomicity, migration, and repeated apply.
- [ ] Add preview information to audit output so agents can explain planned managed changes before applying them.

## Acceptance Criteria

- [ ] Applying a resolved setup creates or updates the manifest, root managed blocks, and supporting guides declared by the selected modules.
- [ ] Content outside managed markers remains byte-for-byte unchanged.
- [ ] A second apply with the same inputs produces no file changes.
- [ ] Missing or incompatible decisions prevent writes and return exit code `3`.
- [ ] Invalid or duplicate markers prevent writes and return a blocking result.
- [ ] A simulated failure before commit leaves every target file at its original content.
- [ ] Obsolete managed artifacts can be removed without deleting unowned files or unowned sections.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `docs/adr/0046-setup-owned-agent-instructions-are-declarative.md`

## Verification

- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_apply*.py'` — expected: adoption, atomicity, preservation, migration, and idempotency flows pass.
- `rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit*.py'` — expected: the read-only contract remains intact after adding apply.
- `rtk git diff --check` — expected: no whitespace errors.

## References

- `_prd.md` → User Stories 2 and 4; Core Features 3–8; Success Metrics 1–3 and 7.
- `_techspec.md` → Data Models; API Contracts: Safe apply; Testing Approach; Build Order 3.
- ADR-0046.
