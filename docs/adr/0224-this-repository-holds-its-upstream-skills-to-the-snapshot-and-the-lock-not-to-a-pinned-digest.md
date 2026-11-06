---
status: accepted
created_at: 2026-10-03T00:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# This repository holds its upstream skills to the snapshot and the lock, not to a pinned digest

Three skill contract tests compared one hand-pinned digest of every
upstream-managed skill tree with the installed trees. The supported restore,
`roundfix baseline update`, rewrites eleven of those trees and their
`skills-lock.json` entries to the Setup Snapshot (ADR-0221), so the restore
broke `make verify` until someone hand-edited the constant, and on 2026-10-02
the restore was abandoned instead (PR #350).

The tests now read the two records the restore itself writes. Every
`skills-lock.json` entry must match its installed tree through the external
skills digest (`computedHash`), and every required external skill of this
repository's Baseline Profile must match the `treeDigest` its Setup Snapshot
pins, through the same comparison the Doctor Command uses. No constant names a
skill tree's digest.

**Amended by [ADR-0364](0364-a-required-external-skill-is-checked-only-for-presence.md) (2026-11-05).** The Doctor Command no longer runs this comparison for any repository. This repository's contract tests keep both rules for its own skills: the vendored catalog test holds each tree to its Setup Snapshot, and the live Repository Skill Set test holds each required skill to its lock entry.

## Consequences

A restore and the tests that check it can no longer disagree, and a hand edit
of an upstream skill still fails, now naming the skill. This repository is held
stricter than an adopter: here a trailing skill fails `make verify`, while
Doctor only warns an adopter about it. A change that moves the embedded Setup
Snapshot therefore restores this repository's skills in the same change. A
skill outside the Profile's snapshot is held only to its lock entry.
