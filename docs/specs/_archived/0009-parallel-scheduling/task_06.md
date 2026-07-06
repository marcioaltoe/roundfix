---
task: task_06
spec: 0009-parallel-scheduling
status: completed
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

- [x] Skill updates + `make skills-sync`
- [x] README config and concurrency notes
- [x] Fixture and binary cross-check
- [x] Glossary pass

## Acceptance Criteria

- [x] Skill text matches shipped behavior exactly; drift check passes
      inside the full gate.
- [x] Documented shapes appear verbatim in CLI test fixtures.
- [x] README documents both keys with the hierarchy and the fixed
      segments rule.
- [x] No new un-glossaried term.

## Verification

- `rtk go run ./cmd/roundfix skills check` — expected: skill artifacts
  validate.
- `make verify` — expected: full gate passes, including the drift check.

## References

`_prd.md` → User Experience; Core Features 1–6. `_techspec.md` → Build
Order 6. ADR-0025, ADR-0026. Repo hard rule (canonical skill ships with CLI
behavior changes).

## Result

- Updated the canonical Roundfix skill with Wave scheduling, `worktree.concurrency`
  default/sequential behavior, Task Worktrees, per-Task Agent Sessions, Task
  Worktree integration conflict semantics, settle resolution order,
  `worktree.location` hierarchy, fixed slug/run-id/task suffix segments, and
  empty-debris reap reporting. Ran `rtk make skills-sync`; `.agents/skills/roundfix`
  and `skills/roundfix` are byte-identical.
- Updated the README configuration and behavior notes for `worktree.location`
  and `worktree.concurrency`, including the fixed hierarchy and parallel
  Verification load expectation for commands such as `make verify`.
- Cross-checked line shapes against fixtures and the built binary:
  `rtk rg -n "Concurrency: 4|roundfix: reaped terminal Worktree path=|roundfix: settle failed after verification: task worktree integration conflict on|integration conflict: shared\\.txt" internal/cli internal/daemon`
  found the documented shapes in CLI/daemon fixtures, and
  `rtk proxy strings bin/roundfix | rtk rg "Concurrency:|reaped terminal Worktree|task worktree integration conflict|integration conflict:|Re-runs one failed Task"`
  found the shipped binary strings.
- Glossary pass: `CONTEXT.md` already defines `Wave` and `Task Worktree`; no
  additional un-glossaried term was introduced.
- Verification: `rtk go test ./internal/cli/ ./internal/daemon/` passed
  (`304 passed in 2 packages`) after rerunning outside the sandbox-restricted
  Go build cache; `rtk go run ./cmd/roundfix skills check` passed; `rtk make verify`
  passed (`711 passed in 17 packages`, skill check passed, build passed).
