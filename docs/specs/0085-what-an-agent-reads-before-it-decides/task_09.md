---
task: task_09
spec: 0085-what-an-agent-reads-before-it-decides
status: pending
type: backend
complexity: high
---

# Task 09: Move the consumers the relocation left behind

## Overview

Task 04 moved archived Specs and findings under `_archived/`, and the resolver
answers the new root. Everything that reads the old layout stayed where it was,
so the Spec's own authoritative gate is red and two contracts still teach the
layout the Spec removed.

The QA gate recorded three findings from one cause:

- **F-01 (Blocks-Completion).** `make verify` exits 2 with nineteen failing
  top-level tests across `internal/cli`, `internal/spec`, `internal/specaudit`,
  `internal/speccheck`, and `internal/worktree`. Their fixtures and assertions
  compose `docs/specs/_archived` directly.
- **F-02 (Trust-Damage).** `_archived/adr/` does not exist, so a retired ADR is
  still in the active read path an Agent loads from `docs/adr/`. The resolver
  answers a directory nothing populated.
- **F-03 (Trust-Damage).** A production sweep found eighteen active files —
  guides, skills, catalog sources, fixtures, source-baseline corpora — naming
  the old layout. An Agent reading them is taught the path the Spec just
  retired.

A relocation whose consumers do not follow is worse than no relocation: it
leaves two layouts, one true and one taught.

## Requirements

1. MUST make `make verify` exit 0 on a genuinely cold cache — the one the
   Makefile exports, not the user-level one.
2. MUST resolve every archive path through the resolver in production code, and
   through a helper in tests, so no fixture or assertion composes the layout
   itself.
3. MUST populate `_archived/adr/` with the retired ADRs and leave the active
   successor in `docs/adr/`, so an Agent loading the active tree reads only
   decisions still in force.
4. MUST update every active guide, skill, catalog source, fixture, and
   source-baseline reference to the single root, and re-record any derived pin
   the change invalidates with the sanctioned command in the same commit.
5. MUST NOT reintroduce the old layout anywhere it was removed, and MUST NOT
   change what any test asserts beyond the path it resolves. A test that passes
   because its expectation was weakened proves nothing about the move.

## Subtasks

- [ ] Move every consumer and fixture onto the resolver.
- [ ] Populate `_archived/adr/` and leave the successor active.
- [ ] Update the guides, skills, catalog sources, and baselines.

## Acceptance Criteria

- [ ] `make verify` exits 0 on a cold cache.
- [ ] No active file composes `docs/specs/_archived` or `docs/findings/_archived`.
- [ ] `_archived/adr/` holds the retired ADRs.
- [ ] Every derived pin the change invalidated is re-recorded.

## Bounded scope

This Task may create or modify only:

- `internal/**`
- `_archived/**`
- `docs/adr/**`
- `docs/agents/**`
- `.agents/skills/**` and their `skills/` mirrors
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_09.md`

## Verification

- `test -d _archived/adr` — expected: exits 0. The directory does not exist before this Task, so a retired ADR still sits in the active read path.
- `test -z "$(grep -rln 'docs/specs/_archived\|docs/findings/_archived' --include='*.go' --include='*.md' --include='*.json' --include='*.yml' internal/ docs/agents/ .agents/ 2>/dev/null | grep -v '^_archived/')"` — expected: exits 0. Eighteen active files name the old layout today.
- `GOCACHE="$PWD/.gocache" go test ./internal/spec ./internal/speccheck ./internal/specaudit ./internal/worktree -count=1 2>&1 | tee /dev/stderr | grep -c '^ok' | grep -q '^4$'` — expected: exits 0. All four packages fail today, so this is a proof rather than a formality.

## References

- `_prd.md` → Goal 1.
- `task_04.md` → the relocation this Task completes.
- `qa/qa-report-2026-08-11.md` → F-01, F-02, F-03.
