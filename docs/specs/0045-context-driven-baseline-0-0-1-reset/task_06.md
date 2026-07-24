---
task: task_06
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
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

- [x] Implement the bounded byte-exhaustive carrier scanner.
- [x] Model managed blocks, unmarked spans, files, and manifest records as
      stable entries.
- [x] Add exact-byte digests, ranges, provenance, and deterministic ordering.
- [x] Expose the inventory through human-readable and machine-readable audit
      output.
- [x] Add symlink, special-file, size-limit, unsafe-path, and ignored-tree
      probes.
- [x] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [x] Concatenating the reported byte spans for each scanned regular file
      accounts for every source byte exactly once.
- [x] Repeated audit of unchanged input emits identical entry IDs, ordering,
      byte ranges, and digests.
- [x] Unmarked text is reported as evidence and is never silently inferred as
      owned, obsolete, or disposable.
- [x] Symlinks, special files, oversized carriers, and unsafe paths fail closed
      with stable diagnostics.
- [x] Audit reports the incompatible baseline and complete inventory with exit
      behavior defined by the TechSpec and leaves repository bytes unchanged.
- [x] No entry is automatically assigned a Readoption disposition.

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

## Result

Implemented a bounded, read-only Source Baseline inventory for incompatible
repository state. The scanner discovers root and nested instruction carriers
plus `docs/agents/`, excludes declared dependency, VCS, and skill-mirror
trees, reads only regular files without following symlinks, and enforces
carrier-count and byte limits. It partitions each carrier into exact managed
blocks, unmarked spans, whole-file evidence, or Setup Manifest records while
retaining byte ranges, SHA-256 digests, carrier digests, base64 source bytes,
and structural provenance.

Audit now reports the incompatible Source Baseline identity and ordered
`sourceEntries` in JSON and the same identity, range, digest, kind, and
provenance inventory in text. Every entry produces an unresolved Readoption
decision, so audit exits `3` without writing. Unsafe carriers exit `1` without
returning a partial inventory. No entry contains a classification,
destination, or disposition.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_source_inventory.py'` — passed, 4 tests. Proves exact byte reconstruction, stable identities and ordering, structural evidence, ignored-tree bounds, fail-closed filesystem and manifest paths, deterministic audit output, exit `3`, and no writes.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_audit.py'` — passed, 11 tests. Proves existing clean, blocking, decision, output-shape, deterministic, and read-only audit contracts remain intact.
- `rtk make skills-sync-check` — passed. Canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — passed. The gate reported 1,694 Go tests, 209 canonical setup tests, 209 distributed setup tests, valid setup assets, a passing Roundfix skill check, and a successful CLI build. The first sandboxed invocation could not access the external Go build cache; the authorized rerun completed with exit `0`.

Acceptance evidence:

- Byte exhaustion and stable identity: `test_inventory_is_byte_exhaustive_structural_and_stable` reconstructs every carrier from contiguous reported spans and compares two unchanged reads exactly.
- Unmarked evidence and no inference: the same test observes repository-owned and unmarked bytes while asserting that serialized entries contain no disposition.
- Fail-closed bounds: `test_inventory_ignores_declared_trees_and_rejects_unsafe_carriers` covers ignored dependency trees, symlinks, special files, and carrier size; `test_inventory_rejects_unsafe_manifest_record_paths` covers traversal in manifest evidence.
- Read-only audit: `test_audit_reports_complete_unresolved_inventory_without_writes` checks incompatible identity, complete entry count and byte total, identical repeated JSON, human-readable inventory, exit `3`, absent dispositions, and unchanged repository bytes.

Follow-ups:

- Task 07 owns structured decision files, semantic classifications, and
  individual Readoption dispositions.
- Task 08 owns Decision Plan and Change Plan integration and all confirmed
  Baseline Readoption writes.
