---
task: task_05
spec: 0052-claude-adapter-standardization
status: pending
type: backend
complexity: high
---

# Task 05: Migrate stale Claude overrides through Setup

## Overview

Generalize the Setup Command's stale-adapter migration from codex-only to
per-runtime, so a machine whose acpx configuration pins a deprecated or bare
Claude command gets a diagnosed, proved, confirmation-gated migration to the
official pinned Claude adapter — the same choreography the Codex migration
already has. Decline, no-input, and failure paths preserve every byte.

## Requirements

1. MUST detect a stale Claude adapter (legacy lineage or unsupported
   version) during Setup's adapter evaluation and propose migration instead
   of hard-failing, exactly as the Codex path does.
2. MUST propose the official pinned Claude override and prove the proposal's
   adapter identity before offering it; the confirmation prompt names the
   runtime being migrated.
3. MUST seed fresh-machine acpx override proposals with the official pinned
   Claude form instead of the bare `claude-agent-acp` command.
4. MUST support migrating both runtimes in one Setup pass, each with its own
   diagnosis and prompt, without disturbing unrelated configuration bytes.
5. MUST preserve every existing decline, `--no-input`, proof-failure, and
   write-failure guarantee: unauthorized targets keep their bytes.
6. SHOULD keep the Setup report lines (`adapter:`, `acpx agents override:`)
   in their current statuses and shapes, with the Codex proposal now naming
   the raised pin.

## Subtasks

- [ ] Parameterize the stale-adapter gate, migration proposal, and re-proof
      by runtime.
- [ ] Add the official Claude override constructor beside the Codex one and
      use it in fresh-machine seeding.
- [ ] Make the migration state and prompts per-runtime.
- [ ] Extend the Setup tests: Claude migration accept and decline tables,
      fresh-machine offer writing the pinned Claude form, both-runtimes
      migration, unrelated-bytes preservation.

## Acceptance Criteria

- [ ] Setup on a config whose claude command is a legacy lineage offers a
      migration prompt naming claude; accepting writes the official pinned
      command and args, declining writes nothing.
- [ ] Setup on a fresh machine that accepts offers writes the official
      pinned Claude override rather than a bare command.
- [ ] A config holding both a stale codex and a stale claude override can
      migrate both in one pass, each independently confirmable.
- [ ] The Codex migration still writes the official Codex package at the
      raised pin.
- [ ] Unrelated agents and bytes in the acpx configuration survive every
      path byte-identical.

## Context

- interface: `internal/cli/setup.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/agent/acpx_runner.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/cli/ -run 'TestRunSetup|TestSetupCommand'` — expected: pass, including the new Claude migration cases.
- `grep -n '"claude-agent-acp"' internal/cli/setup.go ; test $? -eq 1` — expected: no matches (exit 1); the bare-command seed is gone.

## References

`_prd.md` → User Story 2, Core Features 1–3, 7; `_techspec.md` → Build
Order 5, API Contracts; ADR-0055.
