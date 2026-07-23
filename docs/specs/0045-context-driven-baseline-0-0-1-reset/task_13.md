---
task: task_13
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
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

- [x] Update canonical vocabulary and user-facing setup guidance.
- [x] Document TypeScript profile decisions, capabilities, and exact skill
      activations.
- [x] Document Baseline Readoption and Repository-Specific Normative Rules.
- [x] Document the owned-version boundary and protected version surfaces.
- [x] Update the release runbook and Roundfix skill reset-plan guidance.
- [x] Synchronize distributed skill documentation.
- [x] Validate examples, help references, and project-token exclusions.

## Acceptance Criteria

- [x] A reader can complete the 0.0.1 setup and Readoption journey using only
      the documented commands and decision-file contract.
- [x] TypeScript guidance states the exact topology, required stack,
      architecture, HTTP choices, optional modules, capability policy, and
      activation bundles.
- [x] Release guidance clearly stops at an approval-required read-only plan and
      states that destructive cleanup needs a fresh plan plus separate approval.
- [x] Owned and protected version surfaces are enumerated without ambiguity.
- [x] Setup and Roundfix skill guidance matches shipped flags, outputs, exit
      behavior, and ownership.
- [x] Examples contain no external project identity and canonical/distributed
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

## Result

Documented the 0.0.1 operating contract in the canonical glossary, setup and
configuration guides, release runbook, and repo-owned setup and Roundfix
skills. Added a documentation contract test and refreshed the three profile
snapshot digests after the owned skill guidance changed.

### Acceptance evidence

1. `docs/user-guide/context-driven-development.md` now gives the complete
   audit, structured decision-file, preview, confirmation, apply, formatter,
   Verification, audit, and unconfirmed reapply sequence, including every
   Baseline Readoption disposition and Repository-Specific Normative Rule
   destination.
2. The same guide defines the exact Standard TypeScript Monorepo Profile
   topology, stack, architecture, REST and Post-only decisions, ordered
   exceptions, optional modules, capability policy, research tools, and Skill
   Activation bundles.
3. `docs/user-guide/release-runbook.md` and both Roundfix skill copies document
   read-only `release plan --reset-to v0.0.1`, complete inventory,
   `planDigest`, exit 3, a fresh post-QA plan, and separate explicit authority
   for destructive tag or GitHub Release cleanup.
4. The setup guide, configuration guide, and setup skill enumerate the owned
   0.0.1 baseline surfaces and protect operational state, upstream artifact
   versions, third-party dependencies, and repository product versions from
   the reset.
5. The focused CLI help check and documentation contract test confirm that
   setup and Roundfix guidance matches shipped flags, output, exit behavior,
   and skill ownership.
6. `test_documentation_contract.py` rejects external project identity tokens
   and verifies byte-identical canonical/distributed setup and Roundfix skill
   trees; `skills-sync-check` also passed.

### Verification evidence

- `rtk go run -buildvcs=false ./cmd/roundfix release plan --help` — passed;
  help exposes reset mode as read-only, exit 3, and separately authorized
  deletion.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s
  .agents/skills/setup-context-driven/tests -p
  'test_documentation_contract.py'` — passed, 5 tests.
- `rtk make skills-sync-check` — passed.
- `rtk make verify` — passed: both 251-test setup suites, 1,727 Go tests,
  setup asset validation, Roundfix skill checks, and the build completed.

Follow-ups: none.
