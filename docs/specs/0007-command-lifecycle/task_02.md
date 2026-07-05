---
task: task_02
spec: 0007-command-lifecycle
status: pending
type: backend
complexity: high
---

# Task 02: Build the Setup Command

## Overview

One idempotent command that takes a machine from fresh to Run-ready:
verifies Node, installs the pinned acpx on confirmation, probes the runtime
adapters, offers the acpx local-binary agents override when local adapters
exist, and offers User/Project Config creation through the existing Init
Command flows. Verifiable through buffer-captured tests over faked
npm/acpx/filesystem boundaries.

## Requirements

1. MUST add `roundfix setup [--yes] [--no-input]`: ordered checks, one
   deterministic report line per check (`ok | installed | skipped |
   offered: declined | failed`), exit 0 when nothing failed, 1 on any
   failure, 2 on usage errors.
2. MUST check: Node present at the documented minimum; acpx present at the
   pinned version (offering the exact documented install command when
   missing or mismatched, executed on confirmation or `--yes`); the
   configured Agent probes clean.
3. MUST offer the acpx agents override when local adapter binaries are on
   PATH and `~/.acpx/config.json` lacks matching entries: surgical merge of
   only the `agents` entries (file created via the acpx init shape when
   absent), a printed before/after diff of the change, confirmation
   required (or `--yes`).
4. MUST offer User Config and Project Config creation where missing by
   reusing the existing Init Command flows — no duplicated config
   templates.
5. MUST be idempotent: a second run on a healthy environment prints all-ok
   and changes nothing; `--no-input` skips offers (reporting them skipped)
   instead of prompting.
6. MUST route all execution through seams (npm runner, acpx prober,
   filesystem) so tests never touch the real environment.

## Subtasks

- [ ] Command skeleton, flags, dispatch, report format
- [ ] Node/acpx/agent probe checks with install offer
- [ ] Surgical acpx agents override with diff and confirmation
- [ ] Init-flow reuse for missing configs
- [ ] Idempotency and non-interactive tests over faked boundaries

## Acceptance Criteria

- [ ] Fixture matrix: fresh machine (everything offered/installed),
      healthy machine (all ok, zero writes), mismatched acpx (upgrade
      offered), local adapters without override (merge offered, diff
      printed, file byte-checked after), declined offers (reported, no
      writes).
- [ ] The override merge preserves every unrelated key in the acpx config
      file byte-for-byte.
- [ ] `--yes` accepts every offer; `--no-input` skips them as reported.
- [ ] `setup --help` truthful; full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix setup --help` — expected: usage with the
  flags, exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 1; Core Feature 1; Decisions. `_techspec.md` → Setup
Command, Risks (user-owned file), Build Order 2. Round-1 dogfood findings
22 and 27.
