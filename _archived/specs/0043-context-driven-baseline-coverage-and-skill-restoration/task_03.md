---
task: task_03
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: completed
type: backend
complexity: high
---

# Task 03: Reject unresolved Decision Plan references

## Overview

Make audit prove references against the exact artifacts selected by the
Decision Plan rather than incidental files or Markdown syntax. The slice
blocks managed pointers to excluded targets and repository-owned pointers that
are missing, unsafe, or outside the repository.

## Requirements

1. MUST bind every declared reference token to either one managed target or one
   repository-owned relative path.
2. MUST reject a definite setup-owned source whose target is absent from the
   exact definite artifact set, even when a stale target file exists on disk.
3. MUST validate repository-owned targets inside the repository without
   treating them as setup-owned or generating their content.
4. MUST validate the future tree represented by the Decision Plan before apply
   can create a broken pointer.
5. MUST retain relative Markdown-link scanning as defense in depth while using
   typed declarations as the primary path-pointer contract.
6. MUST cover every finite decision branch that includes or excludes a
   referenced artifact across all bundled profiles.

## Subtasks

- [x] Resolve declared reference tokens into selected target paths.
- [x] Validate managed targets against definite Decision Plan artifacts.
- [x] Validate repository-owned paths against repository boundaries and
      existence.
- [x] Integrate future-tree reference findings into audit and apply planning.
- [x] Add profile and decision-transition regressions for missing, excluded,
      stale, and external targets.

## Acceptance Criteria

- [x] Audit blocks a managed reference when its target is excluded or absent
      from the selected Decision Plan.
- [x] A stale on-disk guide cannot satisfy a managed reference omitted from the
      selected artifact set.
- [x] The single-context monorepo case audits clean because its referenced
      monorepo guide is selected.
- [x] A frontend profile with no repository-owned `DESIGN.md` reports one
      blocking path-specific finding and does not create that file.
- [x] Absolute, escaping, and repository-external declared paths fail without
      traversal.
- [x] Every artifact-changing boolean or enum branch either resolves all
      references or produces the expected stable blocking finding.
- [x] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `docs/adr/0046-setup-owned-agent-instructions-are-declarative.md`
- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_assets.py`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_decision_plan_contracts.py`
- interface: `.agents/skills/setup-context-driven/tests/test_audit.py`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_decision_plan_contracts.py`
  — expected: every selected reference resolves across all finite artifact
  branches and excluded targets fail deterministically.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_audit.py`
  — expected: managed and repository-owned broken pointers block audit without
  writes or path traversal.
- `rtk make verify` — expected: the full repository gate passes with exact
  Decision Plan reference validation.

## References

- `_prd.md` → Goal 3; User Story 3; Core Features 4, 9; Success Metrics.
- `_techspec.md` → Data Models: typed references; API Contracts: audit;
  Testing Approach; Build Order 3.
- ADR-0046 → managed identity and repository-owned boundary.
- ADR-0047 → exact selected artifact composition.

## Result

Implemented exact Decision Plan reference validation before audit or apply can
accept a future tree. Each typed template token now has one unambiguous target
binding. Definite setup-owned sources resolve only against definite selected
managed artifacts, while repository-owned references must resolve to an
existing path inside the repository and remain outside setup's generated
content.

Relative Markdown-link scanning remains a second check and now evaluates the
future selected and removed paths, so a stale file scheduled out of the plan
cannot satisfy either a typed managed pointer or an incidental Markdown link.
Unresolved previews retain conditional findings until the Decision Plan is
fully resolved.

Acceptance evidence:

- Excluded `guide.domain` fixtures produce the stable blocking
  `reference.managed.missing` finding for `root.context-workflow`; an existing
  stale `docs/agents/domain.md` does not satisfy it.
- All finite boolean and enum combinations for all three bundled profiles
  resolve without reference findings. The TypeScript/Bun single-context case
  proves `guide.monorepo` remains in the definite artifact set.
- A TypeScript/Bun repository without `DESIGN.md` produces exactly one
  `reference.repository.missing` finding at `DESIGN.md`; audit and apply leave
  the repository unchanged. A symlink outside the repository produces
  `reference.repository.outside`.
- Absolute and escaping repository paths fail asset loading with
  `reference.repository.path.invalid`, and duplicate token bindings fail with
  `reference.token.duplicate`.
- `make skills-sync` regenerated the embedded setup skill from the canonical
  tree; the full gate's `skills-sync-check` passed.

Verification:

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_decision_plan_contracts.py`
  — passed, 10 tests.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_audit.py` —
  passed, 11 tests.
- `rtk python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_*.py'`
  — passed, 99 tests.
- `rtk make verify` — passed after granting access to the host Go build cache;
  1,687 Go tests and 99 setup tests passed, followed by asset validation,
  canonical/embedded synchronization checks, Repository Skill Set checks, and
  the Roundfix build.

### Verification Feedback repair

Daemon Verification attempt 1 reached `skills-sync-check` with a generated
canonical `context_assets` bytecode cache that was absent from the embedded
skill tree. The setup verification commands now use `rtk env` to propagate
`PYTHONDONTWRITEBYTECODE=1` through child processes and also pass `-B` directly
to Python. The audit test helper also passes the same environment invariant to
every CLI subprocess it owns. This prevents the interpreter and its children
from creating runtime cache files without weakening the exact
canonical/embedded tree comparison.

Focused repair evidence:

- `rtk make setup-context-check` — passed, 99 tests and both canonical and
  embedded asset catalogs loaded with interpreter-level bytecode suppression.
- `rtk diff -qr .agents/skills/setup-context-driven skills/setup-context-driven`
  — produced no differences after each focused test and after the setup check.
- `rtk make skills-sync-check` — passed without regenerating or excluding any
  skill-tree content.

The Daemon owns the next full configured Verification attempt.
