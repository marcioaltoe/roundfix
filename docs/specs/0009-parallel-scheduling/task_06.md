---
task: task_06
spec: 0009-parallel-scheduling
status: pending
type: docs
complexity: low
---

# Task 06: Sync docs and the Roundfix skill with parallel scheduling

## Overview

Document the shipped concurrency model in the canonical Roundfix skill and
README, cross-checked against the binary: the scheduler and its cap, Task
Worktrees and their integration semantics, the settle resolution order, the
worktree location hierarchy, and the debris reap. Verifiable through the
skills drift check inside the full gate.

## Requirements

1. MUST document in the canonical Roundfix skill: Wave scheduling with
   `worktree.concurrency` (default 2, 1 = sequential), Task Worktrees and
   the conflict → failed rule, per-Task Agent Sessions, the settle
   resolution order, `worktree.location` hierarchy with the fixed
   slug/run-id segments, and the empty-debris reap; regenerate the embedded
   copy through the sync target.
2. MUST update the README's configuration and behavior notes for both new
   keys and the concurrency expectations (parallel `make verify` load,
   the conservative default).
3. MUST cross-check every documented line shape (header `Concurrency: N`,
   reap report lines, conflict messages) against CLI test fixtures and the
   built binary.
4. MUST verify glossary coverage — Wave and Task Worktree exist; call out
   any further gap instead of inventing language.

## Subtasks

- [ ] Skill updates + `make skills-sync`
- [ ] README config and concurrency notes
- [ ] Fixture and binary cross-check
- [ ] Glossary pass

## Acceptance Criteria

- [ ] Skill text matches shipped behavior exactly; drift check passes
      inside the full gate.
- [ ] Documented shapes appear verbatim in CLI test fixtures.
- [ ] README documents both keys with the hierarchy and the fixed
      segments rule.
- [ ] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → User Experience; Core Features 1–6. `_techspec.md` → Build
Order 6. ADR-0025, ADR-0026. Repo hard rule (canonical skill ships with CLI
behavior changes).
