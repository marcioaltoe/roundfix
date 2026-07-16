---
task: task_06
spec: 0031-decision-driven-setup-generation
status: pending
type: backend
complexity: high
---

# Task 06: Migrate existing manifests and synchronize the workflow

## Overview

Complete the existing-repository journey and ship the coherent portable skill.
A repository configured by spec 0030 will reuse its compatible answers, replace
only obsolete setup-owned inventory and blocks, reach a clean audit, and rerun
idempotently while the skill asks only newly activated dependent questions.

## Requirements

1. MUST keep Setup Manifest schema version 1 and preserve every compatible
   stored `{value, confirmedAt}` decision without reconfirmation.
2. MUST migrate old mixed managed IDs, template versions, rendered digests, and
   module inventory to the concrete Decision Plan.
3. MUST remove an obsolete block or guide only when the manifest and ownership
   markers prove setup ownership, preserving all repository-authored bytes.
4. MUST stop with a blocking finding instead of guessing when legacy ownership
   or marker state is ambiguous.
5. MUST update the skill workflow to present selection and definite or
   conditional preview, then ask one unresolved decision at a time, including
   dependent questions introduced by enabled capabilities.
6. MUST synchronize the canonical and embedded skill copies and retain
   portable execution without a canonical setups checkout.
7. MUST cover all three profiles, all nine decision effects, TypeScript/Bun
   priority variations, required/extra skills, and repeat apply through real
   temporary-repository subprocesses.

## Subtasks

- [ ] Implement schema-v1 inventory and managed-artifact migration from spec
      0030 output.
- [ ] Preserve compatible answers and route only missing dependent decisions.
- [ ] Add ambiguous-ownership failure and owner-byte preservation coverage.
- [ ] Update the setup skill's preview, question, confirmation, and apply flow.
- [ ] Add cross-profile and every-decision macro regression cases.
- [ ] Regenerate the embedded skill and keep snapshot checks portable.
- [ ] Ensure the repository gate includes the focused migration coverage.

## Acceptance Criteria

- [ ] A spec 0030 manifest with all compatible answers migrates without a
      decision finding, preserves confirmation dates, and reaches a clean audit.
- [ ] Enabling a capability whose dependent value is absent requests only that
      missing decision and persists it for the next run.
- [ ] Legacy mixed blocks are removed or replaced only when ownership is proven;
      ambiguous markers block migration without partial writes.
- [ ] Repository-authored bytes before, between, and after migrated managed
      blocks remain byte-for-byte unchanged.
- [ ] TypeScript/Bun, Go CLI/TUI, and Rust CLI repositories each apply, audit,
      and reapply successfully with representative decision combinations.
- [ ] Missing required skills remain blocking, extra skills remain opt-in
      information, and no removal command or file deletion is introduced.
- [ ] Canonical and embedded skill copies match, the normal workflow needs no
      developer-specific checkout, and the full repository gate passes.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/no-workarounds/SKILL.md`
- instruction: `.agents/skills/setup-context-driven/SKILL.md`
- instruction: `.agents/skills/systematic-debugging/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`
- interface: `.agents/skills/setup-context-driven/tests/test_workflow.py`
- interface: `docs/specs/0030-context-driven-agent-instructions/qa/manual_variation_qa.py`
- interface: `Makefile`

## Verification

- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_manifest_migration*.py'` — expected: compatible-answer reuse, inventory migration, ambiguous ownership, no-partial-write behavior, owner-byte preservation, and idempotency pass.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_macro_profiles*.py'` — expected: all profiles and representative decision combinations apply, audit, and reapply successfully with correct skill semantics.
- `PYTHONDONTWRITEBYTECODE=1 rtk python3 -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_workflow*.py'` — expected: conditional preview, one-question routing, confirmation, and portable orchestration pass.
- `rtk make skills-sync-check` — expected: canonical and embedded workflow skill copies match.
- `rtk make verify` — expected: formatting, Go tests, all setup-context-driven tests, asset validation, skill checks, sync check, and build pass.

## References

- `_prd.md` → Goals 4–5; Core Feature 7; Success Criteria for spec 0030 migration, idempotency, and owner-byte preservation; Non-goals.
- `_techspec.md` → Data models: Setup Manifest compatibility; Integration points; Testing approach; Build Order 5–6; Risks and considerations.
- ADR-0046 and ADR-0047.
