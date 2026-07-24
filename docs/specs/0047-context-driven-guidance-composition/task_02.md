---
task: task_02
spec: 0047-context-driven-guidance-composition
status: pending
type: backend
complexity: medium
---

# Task 02: Complete ADR and Findings contracts

## Overview

Render the repository-owned ADR lifecycle overlay and complete Findings
template through the active documentation guide. The result stays
self-contained while preserving every upstream-managed skill byte.

## Requirements

1. MUST render the accepted ADR statuses, RFC 3339 UTC timestamps, nullable
   deprecation timestamp, and nullable superseding ADR identity.
2. MUST treat only `accepted` as active and preserve the documented legacy ADR
   compatibility rule.
3. MUST render the complete Findings frontmatter, evidence-first body, routing,
   retained-practice, and dated addendum structure.
4. MUST not rewrite existing ADRs solely for metadata adoption.
5. MUST prove that `domain-modeling/ADR-FORMAT.md` and all other
   upstream-managed skills remain byte-identical.

## Subtasks

- [ ] Add the ADR lifecycle overlay clauses and template.
- [ ] Complete the Findings Operational Contract template.
- [ ] Add legacy ADR and lifecycle validation fixtures.
- [ ] Add upstream-skill immutability guards.

## Acceptance Criteria

- [ ] A new generated ADR template contains every confirmed lifecycle field.
- [ ] Legacy active and explicitly inactive ADR fixtures classify correctly.
- [ ] The Findings template contains every confirmed section and lifecycle
  state.
- [ ] Existing ADR fixture bytes remain unchanged.
- [ ] Any upstream skill byte change fails the focused guard.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/templates/guides/docs-layout.md`
- interface: `.agents/skills/domain-modeling/ADR-FORMAT.md`

## Verification

- `rtk go test -count=1 ./internal/baseline ./skills -run 'TestADRLifecycleContract|TestFindingsOperationalContract|TestUpstreamADRFormatUnchanged'` — expected: lifecycle, template, legacy, and ownership guards pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 4; User Stories 4–5; Core Features 9–13 and 17.
- `_techspec.md` → Integration Points; Testing Approach; Build Order 1.
- ADR-0074 → operative rules remain in repository guidance.
