---
task: task_06
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: backend
complexity: high
---

# Task 06: Inventory incompatible Source Baselines

## Overview

Detect pre-0.0.1 setup content without guessing what it means. Audit must
produce a byte-exhaustive, structurally bounded inventory that a later
Decision Plan can resolve entry by entry.

## Requirements

1. MUST inventory every regular-file byte in the bounded carriers declared by
   the TechSpec for root agent instructions and `docs/agents/` content.
2. MUST identify managed markers, unmarked text spans, files, and relevant
   manifest records as stable Source Baseline Entries.
3. MUST preserve exact source bytes, carrier identity, byte ranges, digests,
   and structural provenance without semantic classification.
4. MUST use deterministic ordering and stable entry identifiers across
   repeated reads of unchanged input.
5. MUST reject symlinks, special files, unsafe paths, oversized inputs, and
   traversal beyond the declared carriers.
6. MUST make `audit` report incompatible Source Baseline identity and its
   complete entry inventory without writing to the repository.
7. MUST block planning until every incompatible entry can receive an explicit
   Readoption disposition.

## Subtasks

- [ ] Implement the bounded byte-exhaustive carrier scanner.
- [ ] Model managed blocks, unmarked spans, files, and manifest records as
      stable entries.
- [ ] Add exact-byte digests, ranges, provenance, and deterministic ordering.
- [ ] Expose the inventory through human-readable and machine-readable audit
      output.
- [ ] Add symlink, special-file, size-limit, unsafe-path, and ignored-tree
      probes.
- [ ] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [ ] Concatenating the reported byte spans for each scanned regular file
      accounts for every source byte exactly once.
- [ ] Repeated audit of unchanged input emits identical entry IDs, ordering,
      byte ranges, and digests.
- [ ] Unmarked text is reported as evidence and is never silently inferred as
      owned, obsolete, or disposable.
- [ ] Symlinks, special files, oversized carriers, and unsafe paths fail closed
      with stable diagnostics.
- [ ] Audit reports the incompatible baseline and complete inventory with exit
      behavior defined by the TechSpec and leaves repository bytes unchanged.
- [ ] No entry is automatically assigned a Readoption disposition.

## Context

- instruction: `docs/adr/0064-baseline-readoption-uses-byte-exhaustive-structural-inventory.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_audit.py`
- interface: `.agents/skills/setup-context-driven/tests/test_upgrade_retention.py`
- interface: `.agents/skills/setup-context-driven/tests/test_support.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_source_inventory.py'` — expected: byte coverage is exhaustive, identities are stable, and unsafe carriers fail closed.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit.py'` — expected: incompatible baselines expose a complete read-only inventory and block unresolved planning.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Core Feature 13; User Stories 1 and 5; User Experience.
- `_techspec.md` → Implementation Design: Interfaces and Data Models; API
  Contracts; Build Order 6.
- ADR-0064 → byte-exhaustive structural inventory without semantic inference.
