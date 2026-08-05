---
task: task_03
spec: 0075-typed-docs-backlog
status: pending
type: docs
complexity: low
---

# Task 03: Adopt the contract in this repository

## Overview

The clauses are what adopting repositories receive; this repository is one of
them. This slice re-applies the baseline so the checked-in guide carries the new
contract, and seeds `docs/backlog/` with its first real entry using the template
verbatim — which is also the cheapest proof the template is usable.

## Requirements

1. MUST re-apply the baseline so `docs/agents/docs-layout.md` carries the
   backlog contract.
2. MUST create `docs/backlog/` with at least one real entry, using the template
   verbatim rather than an invented shape.
3. MUST make that first entry a genuine item, not a placeholder. The
   verification-performance contract the 2026-08-03 handoff describes as
   belonging in `docs/backlog/` as a `perf` entry is the natural candidate.
4. MUST NOT migrate any existing finding into the backlog. Deprecating findings
   into a backlog type is explicitly deferred.
5. MUST leave every other generated guide unchanged.

## Subtasks

- [ ] Re-apply the baseline and confirm the guide updated.
- [ ] Write the first backlog entry from the template.
- [ ] Confirm no finding was migrated and no other guide moved.

## Acceptance Criteria

- [ ] `docs/agents/docs-layout.md` documents `docs/backlog/` and its contract.
- [ ] `docs/backlog/` holds at least one entry whose frontmatter matches the
      contract exactly.
- [ ] That entry is a real item with real content, not a placeholder.
- [ ] No file under `docs/findings/` moved or changed.
- [ ] No generated guide other than the layout guide changed.

## Verification

- `grep -q "docs/backlog" docs/agents/docs-layout.md` — expected: exit 0.
- `ls docs/backlog/*.md | grep -q .` — expected: exit 0; the directory holds an
  entry.
- `head -8 docs/backlog/*.md | grep -qE "^type: (feat|fix|perf|refactor)( +#.*)?$"`
  — the optional trailing comment is not laxity: the template this Spec adopts
  documents the enum inline, as `type: perf # feat | fix | perf | refactor`, so
  an end-anchored pattern rejects the very shape the contract defines. The
  value is still pinned to the four members at the start of the line. Measured
  on 2026-08-05: the original pattern failed against a conforming file.
  — expected: exit 0; the entry's type is one of the four.
- `git diff --name-only HEAD -- docs/findings | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no finding moved.
- `make verify` — expected: exit 0.

## References

- `_prd.md` → Core Features; Non-Goals (no finding migration).
- `_techspec.md` → Build Order 3.
- `docs/handoffs/2026-08-03-after-the-0.3.1-release.md` → the `perf` entry it
  describes.
