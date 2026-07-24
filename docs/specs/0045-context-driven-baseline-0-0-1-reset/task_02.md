---
task: task_02
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
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

- [x] Reconcile every prior clause and recommendation with a retained or
      intentionally excluded 0.0.1 entry.
- [x] Author the complete governed guide and template corpus.
- [x] Preserve the findings lifecycle template as a first-class source entry.
- [x] Reduce root instruction sources to compact delegations.
- [x] Add project-token and corpus-completeness checks.
- [x] Remove obsolete bundled-asset compatibility after migration.
- [x] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [x] Every retained source requirement has exactly one stable entry identity
      and is reachable from the independent Source Baseline index.
- [x] Generated guides contain the full Operational Contracts and the findings
      lifecycle template specified by the PRD.
- [x] Root instruction output delegates to complete guides without duplicating
      their full content.
- [x] A denied project token inserted into any governed source entry fails the
      corpus gate.
- [x] New catalog output contains only strict 0.0.1 schemas and versions.
- [x] Canonical and distributed setup skill trees are byte-identical.

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

## Result

Published the independent Standard TypeScript Monorepo Source Baseline as 56
marker-bounded, individually indexed entries. The corpus distinguishes
Normative Clauses, recommendations, and structured Operational Contracts;
renders complete root and guide carriers; preserves the findings template and
lifecycle; and records an individual disposition for all 51 prior governed
clauses. The current Source Baseline catalog emits only strict `0.0.1`
documents, so prior catalog generations remain evidence for Baseline
Readoption rather than current governed output.

Extended Source Baseline validation with explicit enforcement, carrier, and
Operational Contract structure metadata. Corpus policy now fails closed on
denied project tokens, copied generated managed markers, and machine-specific
paths. Root output remains a compact delegation index while durable behavior
lives in generated `docs/agents/` carriers.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_governed_corpus.py'`
  — passed 6 tests covering indexed identity, complete carrier rendering,
  prior-clause accounting, compact root delegation, corpus policy rejection,
  and strict `0.0.1` documents.
- `rtk make skills-sync-check` — passed; canonical and distributed setup skill
  trees are byte-identical.
- `rtk make verify` — passed on the approved unchanged rerun: 1,694 Go tests,
  183 canonical setup tests, 183 distributed setup tests, both setup asset
  loads, owned-skill validation, and the build passed. The initial sandboxed
  attempt could not read the host Go build cache and did not reach a product
  failure.

Acceptance evidence:

- Stable identity and reachability: the governed-corpus gate compared all 56
  unique manifest entries with the independent index order and required a
  non-empty generated carrier for every entry.
- Complete Operational Contracts: nine first-class contracts retain the root
  template, docs directory matrix, findings template and lifecycle, Spec route
  matrix, Task ownership, autonomous Supervisor/ACP Runtime protocol, research
  procedure, and Secondbrain protocol. Carrier rendering proved the findings
  frontmatter, status semantics, and autonomous prohibition survive intact.
- Compact root delegation: the root carrier stayed below 1,500 bytes, linked
  the complete agent, Spec-routing, and docs-layout guides, and contained no
  copy of the findings template.
- Project-agnostic corpus: mutations containing a denied project token, a
  generated managed marker, or a machine-specific path each failed with its
  asserted stable diagnostic.
- Strict generation: every JSON document below `assets/source-baselines/`
  reports a `setup-context-driven/.../0.0.1` schema and string version
  `0.0.1`; the 51-entry accounting maps every retained target to a current
  indexed entry and gives each rejection its own reason.
- Distribution parity: `rtk make skills-sync-check` passed after regenerating
  the distributed setup skill tree from the canonical authorial tree.

Follow-ups: Task 03 owns typed exact Skill Activation bundle rendering. Later
profile and Readoption Tasks consume this corpus; they are outside this Task's
diff.
