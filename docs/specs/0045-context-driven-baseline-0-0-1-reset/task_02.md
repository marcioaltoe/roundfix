---
task: task_02
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: docs
complexity: high
---

# Task 02: Publish the project-agnostic governed corpus

## Overview

Publish the complete 0.0.1 corpus that setup can govern without relying on an
example repository. Every retained Normative Clause, recommendation, template,
and Operational Contract must have one explicit home and one indexed identity.

## Requirements

1. MUST account individually for every retained Normative Clause,
   recommendation, template, and Operational Contract identified by the PRD.
2. MUST publish complete project-agnostic content for autonomous work, backend,
   frontend, domain, docs layout, Spec layout and routing, local issue tracking,
   Secondbrain, and TypeScript/Bun guidance.
3. MUST retain the findings lifecycle template, including its required
   frontmatter and status semantics.
4. MUST keep root instructions compact by delegating durable detail to
   governed agent guides without weakening the governed clauses.
5. MUST remove obsolete catalog-generation compatibility once every bundled
   source document is represented by the strict 0.0.1 Source Baseline.
6. MUST reject project names, repository-specific paths, brands, or copied
   generated artifacts in the governed corpus.
7. MUST keep normative requirements distinguishable from recommendations and
   explanatory guidance.

## Subtasks

- [ ] Reconcile every prior clause and recommendation with a retained or
      intentionally excluded 0.0.1 entry.
- [ ] Author the complete governed guide and template corpus.
- [ ] Preserve the findings lifecycle template as a first-class source entry.
- [ ] Reduce root instruction sources to compact delegations.
- [ ] Add project-token and corpus-completeness checks.
- [ ] Remove obsolete bundled-asset compatibility after migration.
- [ ] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [ ] Every retained source requirement has exactly one stable entry identity
      and is reachable from the independent Source Baseline index.
- [ ] Generated guides contain the full Operational Contracts and the findings
      lifecycle template specified by the PRD.
- [ ] Root instruction output delegates to complete guides without duplicating
      their full content.
- [ ] A denied project token inserted into any governed source entry fails the
      corpus gate.
- [ ] New catalog output contains only strict 0.0.1 schemas and versions.
- [ ] Canonical and distributed setup skill trees are byte-identical.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/docs-layout.md`
- instruction: `docs/agents/spec-routing.md`
- instruction: `docs/agents/skill-governance.md`
- interface: `.agents/skills/setup-context-driven/assets/source-baselines`
- interface: `.agents/skills/setup-context-driven/assets/templates`
- interface: `.agents/skills/setup-context-driven/tests/test_assets.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_governed_corpus.py'` — expected: every governed clause and Operational Contract is indexed, project-agnostic, and rendered from a strict 0.0.1 source.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 2 and 5; Core Features 4, 6, and 12; User Story 4;
  Non-Goals / Out of Scope.
- `_techspec.md` → System Architecture; Coverage Map; Build Order 2.
- ADR-0060 → complete project-agnostic Source Baseline content.
- ADR-0062 → removal of obsolete owned compatibility layers.
