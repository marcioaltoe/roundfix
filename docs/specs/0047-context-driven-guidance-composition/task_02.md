---
task: task_02
spec: 0047-context-driven-guidance-composition
status: completed
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

- [x] Add the ADR lifecycle overlay clauses and template.
- [x] Complete the Findings Operational Contract template.
- [x] Add legacy ADR and lifecycle validation fixtures.
- [x] Add upstream-skill immutability guards.

## Acceptance Criteria

- [x] A new generated ADR template contains every confirmed lifecycle field.
- [x] Legacy active and explicitly inactive ADR fixtures classify correctly.
- [x] The Findings template contains every confirmed section and lifecycle
  state.
- [x] Existing ADR fixture bytes remain unchanged.
- [x] Any upstream skill byte change fails the focused guard.

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

## Result

The active documentation guide now renders the repository-owned ADR lifecycle
overlay on top of the upstream body format. It includes all five statuses, UTC
timestamps, nullable deprecation and superseding fields, accepted-only active
semantics, legacy compatibility, and the no-rewrite rule. The same guide now
renders one copyable Findings template with session context, evidence-first
findings, root-cause and routing fields, retained practices, and dated
append-only addenda.

Verification:

- `rtk go test -count=1 ./internal/baseline ./skills -run
  'TestADRLifecycleContract|TestFindingsOperationalContract|TestUpstreamADRFormatUnchanged'`
  passed 7 tests in 2 packages.
- `rtk make verify` passed 2,122 tests in 22 packages, 4 focused skill contract
  tests, the Roundfix skill check, and the final build.

Acceptance evidence:

- `TestADRLifecycleContract` proves the generated ADR template contains every
  lifecycle field, all accepted statuses, RFC 3339 UTC placeholders,
  accepted-only active semantics, the legacy rule, and the upstream body-format
  reference.
- The accepted, deprecated, legacy-active, and legacy-inactive fixtures prove
  lifecycle and legacy classification. The same test proves planning creates
  no postimage for an existing ADR and leaves its fixture bytes unchanged.
- `TestFindingsOperationalContract` proves the generated Findings template
  contains all four states, frontmatter, session context, symptom/evidence,
  proven-or-unknown root cause, action or Spec routing, retained practices,
  lifecycle rules, and the dated addendum.
- `TestUpstreamADRFormatUnchanged` pins the complete upstream-managed skill tree
  and separately pins `.agents/skills/domain-modeling/ADR-FORMAT.md`; the final
  diff contains no upstream-managed skill file.

Follow-ups: none.
