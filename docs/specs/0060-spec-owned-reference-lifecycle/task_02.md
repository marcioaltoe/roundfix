---
task: task_02
spec: 0060-spec-owned-reference-lifecycle
status: pending
type: docs
complexity: low
---

# Task 02: Land the layout, the routing rule, and the glossary

## Overview

The Skills now perform and validate the transition; the guides still describe a
world where a finding stays in the findings tree forever. Update the three
documents that define where things live, so an author reading the layout reaches
the same conclusion the Skills enforce.

## Requirements

1. MUST add `docs/specs/<slug>/references/` to the layout in
   `docs/agents/docs-layout.md`, with its index.
2. MUST state, in the inbox and findings entries of that layout, that an adopted
   document leaves the tree and that its Git history at the old path remains the
   discovery trail.
3. MUST state in `docs/agents/spec-routing.md` that committing to implementation
   transfers ownership of the adopted sources to the Spec.
4. MUST update the `CONTEXT.md` `Spec` glossary entry's artifact set to include
   the references directory and its index.
5. MUST NOT describe behavior that task_01 did not ship.
6. MUST keep the one-purpose-per-folder layout intact — this adds a directory
   inside the Spec folder, it does not reorganize `docs/`.
7. MUST edit only repository-authored sections; these guides carry setup-owned
   blocks that are not authorized here.

## Subtasks

- [ ] Add the references directory to the docs layout and amend the inbox and
      findings entries.
- [ ] Add the ownership-transfer rule to spec routing.
- [ ] Update the `Spec` glossary entry.

## Acceptance Criteria

- [ ] The layout names `docs/specs/<slug>/references/` and its index.
- [ ] The inbox and findings entries state that an adopted document leaves.
- [ ] Spec routing states when ownership transfers.
- [ ] The `Spec` glossary entry lists the references directory.
- [ ] `git diff` touches no setup-owned marker block.

## Context

- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/spec-routing.md`
- interface: `CONTEXT.md`

## Verification

- `make verify` — expected: exit 0.
- `git diff -- docs/agents/` — expected: no change inside a setup-owned marker
  block.

## References

`_prd.md` → Goals 1–2, Decisions (glossary entries update in this Spec's
documentation task); `_techspec.md` → Build Order 4, System Architecture.
