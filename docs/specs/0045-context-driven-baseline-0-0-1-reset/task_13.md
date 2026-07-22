---
task: task_13
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: docs
complexity: medium
---

# Task 13: Document the 0.0.1 operating contract

## Overview

Publish the maintainer and user guidance needed to operate the new baseline
without relying on implementation archaeology. Documentation must distinguish
safe read-only planning from separately approved destructive release cleanup
and describe Readoption as an explicit user decision workflow.

## Requirements

1. MUST document the Source Baseline, Normative Clause Manifest, Operational
   Contract, Repository Capability, Skill Activation, HTTP Contract Decision,
   Baseline Readoption, and Repository-Specific Normative Rules using canonical
   glossary terms.
2. MUST document how to select the Standard TypeScript Monorepo Profile,
   provide REST or Post-only policy and exceptions, satisfy required
   capabilities, and interpret recommended-capability warnings.
3. MUST document the complete audit, structured decision-file, preview,
   confirm-plan, apply, formatter, Verification, audit, and reapply workflow.
4. MUST document that local code search precedes external research and that
   Context7/Exa are required while Firecrawl/`rtk`/`rg` are recommended.
5. MUST document the 0.0.1 owned-version boundary and the operational,
   upstream, and third-party versions that do not reset.
6. MUST update the release runbook for read-only
   `release plan --reset-to v0.0.1`, exit 3, plan digest, complete inventory,
   and a separate approval before any tag or GitHub Release deletion.
7. MUST update the shipped setup-context-driven and Roundfix skill guidance to
   match the implemented CLI behavior and ownership contract.
8. MUST keep examples project-agnostic and MUST NOT perform live repository,
   tag, Release, package, or publication mutation as part of documentation.

## Subtasks

- [ ] Update canonical vocabulary and user-facing setup guidance.
- [ ] Document TypeScript profile decisions, capabilities, and exact skill
      activations.
- [ ] Document Baseline Readoption and Repository-Specific Normative Rules.
- [ ] Document the owned-version boundary and protected version surfaces.
- [ ] Update the release runbook and Roundfix skill reset-plan guidance.
- [ ] Synchronize distributed skill documentation.
- [ ] Validate examples, help references, and project-token exclusions.

## Acceptance Criteria

- [ ] A reader can complete the 0.0.1 setup and Readoption journey using only
      the documented commands and decision-file contract.
- [ ] TypeScript guidance states the exact topology, required stack,
      architecture, HTTP choices, optional modules, capability policy, and
      activation bundles.
- [ ] Release guidance clearly stops at an approval-required read-only plan and
      states that destructive cleanup needs a fresh plan plus separate approval.
- [ ] Owned and protected version surfaces are enumerated without ambiguity.
- [ ] Setup and Roundfix skill guidance matches shipped flags, outputs, exit
      behavior, and ownership.
- [ ] Examples contain no external project identity and canonical/distributed
      skill trees are byte-identical.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/docs-layout.md`
- instruction: `docs/agents/skill-governance.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/release-runbook.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`

## Verification

- `rtk go run -buildvcs=false ./cmd/roundfix release plan --help` — expected: help and runbook describe the same reset flags, read-only behavior, and exit category.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_documentation_contract.py'` — expected: canonical terms, complete workflows, capability policy, ownership boundaries, and project-token exclusions are present.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–5; User Stories 1–6; User Experience; Non-Goals / Out
  of Scope.
- `_techspec.md` → API Contracts; Integration Points; Risks & Considerations;
  Build Order 8.
- ADR-0060 through ADR-0065 → documented baseline, profile, version,
  Readoption, HTTP, and release-reset contracts.
