---
task: task_05
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: pending
type: backend
complexity: high
---

# Task 05: Prove portable Repository Skill Set snapshots

## Overview

Make bundled external skill snapshots reproducible and make audit prove each
complete installed directory offline. Drift findings expose immutable portable
provenance and an exact restoration preview without treating external lock
compatibility metadata as the setup content authority.

## Requirements

1. MUST migrate every external required skill to snapshot v2 with a portable
   GitHub identity, full immutable commit, safe source-relative directory, and
   lowercase complete-tree digest.
2. MUST make the top-level snapshot digest cover canonical serialization of
   complete normalized records rather than path names alone.
3. MUST hash every regular file in bytewise POSIX-path order with unambiguous
   length framing, exclude `.git` and `node_modules`, and reject links and
   special files.
4. MUST make `sync-setups` prove external bytes at the declared commit and
   reject empty provenance, mutable refs, unsafe paths, dirty/unmatched source
   bytes, and machine-local persisted sources.
5. MUST preserve the repo-owned source-precedence contract and keep repo-owned
   skills distinct from external restoration authority.
6. MUST make audit compare complete external skill directories locally and
   report nested additions, edits, removals, and unsafe entries as drift.
7. MUST add structured remediation with exact source/ref/path/digest and
   preview argv while keeping audit read-only, network-free, and free of generic
   branch refresh advice.

## Subtasks

- [ ] Add the portable complete-tree digest and unsafe-entry validation.
- [ ] Migrate setup snapshot generation to immutable normalized records.
- [ ] Prove synchronized external bytes against their declared commit.
- [ ] Compare installed external skill trees against snapshot authority.
- [ ] Emit structured missing/drift remediation without fetching content.
- [ ] Add nested-tree, provenance, source-precedence, and no-mutation fixtures.

## Acceptance Criteria

- [ ] Every external required skill in every bundled snapshot has complete,
      immutable, portable provenance and a valid tree digest.
- [ ] Editing, adding, or removing a nested file changes the digest and makes
      audit report the affected skill as drifted.
- [ ] A symlink, hard link, device, escaping path, mutable ref, absolute source,
      or content/ref mismatch blocks synchronization or audit at the owning
      boundary.
- [ ] Roundfix-owned skill digests cannot be replaced by same-path external
      checkout content.
- [ ] Drift JSON contains the exact restoration provider, skill, source, full
      ref, source path, expected digest, and preview argv.
- [ ] Audit performs no Git invocation, network access, lock mutation, or skill
      mutation.
- [ ] Re-synchronizing unchanged committed sources produces byte-identical
      snapshots.
- [ ] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `docs/agents/skill-governance.md`
- instruction: `docs/specs/0036-doctor-skill-readiness/_techspec.md`
- interface: `.agents/skills/setup-context-driven/assets/setups`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_sync_setups.py`
- interface: `.agents/skills/setup-context-driven/tests/test_skills.py`
- interface: `skills-lock.json`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_sync_setups.py`
  — expected: immutable provenance, complete-tree digests, clean-source proof,
  and repo-owned precedence cases pass.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_skills.py`
  — expected: complete installed directories, structured remediation, unsafe
  entries, and offline no-mutation cases pass.
- `rtk make verify` — expected: the full repository gate passes with portable
  snapshot v2 assets in canonical and embedded catalogs.

## References

- `_prd.md` → Goal 5; User Story 5; Core Features 7, 9; User Experience;
  Success Metrics.
- `_techspec.md` → Data Models: snapshot v2 and remediation; Integration
  Points: skills-lock.json and Spec 0036; Testing Approach; Build Order 5.
- Spec 0036 TechSpec → external lock compatibility boundary retained outside
  setup content integrity.
- `docs/agents/skill-governance.md` → repo-owned and external skill authority.
