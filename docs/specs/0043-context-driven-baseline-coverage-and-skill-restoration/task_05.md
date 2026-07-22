---
task: task_05
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: completed
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

- [x] Add the portable complete-tree digest and unsafe-entry validation.
- [x] Migrate setup snapshot generation to immutable normalized records.
- [x] Prove synchronized external bytes against their declared commit.
- [x] Compare installed external skill trees against snapshot authority.
- [x] Emit structured missing/drift remediation without fetching content.
- [x] Add nested-tree, provenance, source-precedence, and no-mutation fixtures.

## Acceptance Criteria

- [x] Every external required skill in every bundled snapshot has complete,
      immutable, portable provenance and a valid tree digest.
- [x] Editing, adding, or removing a nested file changes the digest and makes
      audit report the affected skill as drifted.
- [x] A symlink, hard link, device, escaping path, mutable ref, absolute source,
      or content/ref mismatch blocks synchronization or audit at the owning
      boundary.
- [x] Roundfix-owned skill digests cannot be replaced by same-path external
      checkout content.
- [x] Drift JSON contains the exact restoration provider, skill, source, full
      ref, source path, expected digest, and preview argv.
- [x] Audit performs no Git invocation, network access, lock mutation, or skill
      mutation.
- [x] Re-synchronizing unchanged committed sources produces byte-identical
      snapshots.
- [x] Canonical and embedded setup skill trees are synchronized after the
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

## Result

Implemented snapshot v2 generation and offline Repository Skill Set audit.
External records now carry a GitHub repository, full commit, safe source
directory, and framed complete-tree digest. Synchronization proves clean
working bytes against the declared commit before writing. Repo-owned records
retain their local `contentDigest` authority. Audit compares complete external
directories and emits immutable structured remediation without invoking Git or
changing repository state.

Verification:

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_sync_setups.py`
  passed: 7 tests cover complete-tree framing, excluded directories, unsafe
  entries, immutable provenance, content/ref mismatch, repo-owned precedence,
  and byte-idempotent synchronization.
- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_skills.py`
  passed: 11 tests cover nested additions, edits, removals, unsafe installed
  entries, exact remediation JSON, lock-hash separation, and offline
  no-mutation audit behavior.
- `rtk python3 -B .agents/skills/setup-context-driven/scripts/context_setup.py sync-setups --source-dir /Users/marcio/dev/skills/setups --check --format json`
  passed with zero findings against commit
  `236847f6956134bf468abb641bac0493a899bca5`.
- `rtk make verify` passed after allowing access to the existing Go build
  cache; the initial sandboxed run could not open that cache.
- `rtk git diff --check` passed.

Acceptance evidence:

1. All three bundled setup snapshots load as snapshot v2 records at version 2;
   asset mutation tests reject incomplete provenance, mutable refs, unsafe
   paths, invalid digests, machine-local fields, and record/digest drift.
2. Installed-tree tests independently edit, add, and remove nested files and
   observe `skills.required.drift` for the affected external skill.
3. Portable-tree and synchronization tests reject symbolic links, hard links,
   special files, absolute or escaping paths, mutable refs, and source bytes
   that do not match the declared commit. The implementation rejects every
   non-regular filesystem entry, including device nodes.
4. The synchronization fixture changes same-path external bytes for
   `setup-context-driven` and proves the Roundfix-owned digest remains
   unchanged.
5. Missing and drift findings include `provider`, `skill`, `source`, `ref`,
   `sourcePath`, `expectedDigest`, and `previewArgv`.
6. The audit fixture places a failing Git shim first on `PATH`, snapshots the
   repository before and after, and proves no Git call or file mutation occurs;
   lock compatibility hashes do not replace snapshot authority.
7. A second synchronization of unchanged committed fixtures produces no
   findings and byte-identical assets; the canonical checkout check also
   reports no drift.
8. `make skills-sync` refreshed `skills/setup-context-driven`, and both
   `make skills-sync-check` and `make verify` passed.

Follow-up: Task 06 owns execution of the `restore-skills` preview and apply
argv emitted by this slice.
