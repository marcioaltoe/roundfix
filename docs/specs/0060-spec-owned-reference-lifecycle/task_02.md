---
task: task_02
spec: 0060-spec-owned-reference-lifecycle
status: completed
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

## Result

### Implementation

- The docs layout now gives `docs/specs/<slug>/references/` one Spec-local
  purpose and names `_index.md` as the adopted-source provenance index.
- The inbox and findings lifecycle entries now say that adoption moves the
  document out of its former tree and leaves Git history at the old path as
  the discovery trail. The finding entry preserves the shipped status/link
  before move order.
- Spec routing now transfers adopted-source ownership when a Spec commits to
  implementation, assigns primary ownership to the first committing Spec, and
  directs later Specs to link the owner's copy.
- The `Spec` glossary definition now includes adopted sources under
  `references/` and their `references/_index.md` index.

### Focused checks

- Pre-change `rtk rg` probes for the exact Spec-local references path and the
  routing ownership-transfer rule each exited 1; the glossary probe showed the
  former artifact set without references.
- `rtk rg -n -F 'docs/specs/<slug>/references/'
  docs/agents/docs-layout.md` — exit 0; matched the inbox, findings, and
  Spec-local references entries.
- Focused `rtk rg -n` inspections for the departure/history wording, routing
  transfer, and `Spec` glossary entry — all exited 0.
- `rtk git diff --check` — exit 0.
- Raw `HEAD` versus worktree setup-block extraction reported byte-identical
  content for both setup-owned blocks in `docs-layout.md` and the one
  setup-owned block in `spec-routing.md`.
- `rtk git status --short` — exit 0; only `CONTEXT.md`, the two authorized
  repository guides, and this Task file are modified. Git emitted its existing
  fsmonitor IPC warning without changing the command outcome.

### Acceptance evidence

1. The layout table names `docs/specs/<slug>/references/`; its job names
   `_index.md` and the provenance fields it records.
2. Both the inbox and findings entries say an adopted document leaves its
   former tree and that Git history at the old path remains the discovery
   trail.
3. Spec routing says ownership transfers when a Spec commits to implementation
   and identifies the first committing Spec as primary owner.
4. The `Spec` glossary artifact set lists `references/` and
   `references/_index.md`.
5. Raw-byte comparison against `HEAD` proves every setup-owned marker block in
   both edited guides is unchanged; all guide edits are in repository-authored
   sections.

### Daemon-owned verification

Neither command under `## Verification` was run in this Agent turn. The Daemon
owns both declared commands, Task status, and terminal settlement.
