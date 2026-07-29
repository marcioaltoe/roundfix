---
task: task_05
spec: 0052-claude-adapter-standardization
status: completed
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

## Result

Setup now evaluates Adapter Readiness for every distinct ACP Runtime referenced
by the required Agent Selection Profiles. A typed lineage or version failure
for an explicit Codex or Claude acpx override becomes a per-runtime migration:
Setup builds the official pinned override, proves that exact command before
offering it, and records one runtime-named confirmation per migration. Any
decline aborts the acpx write, preserving the existing confirmation
choreography and every target byte.

Fresh-machine Claude seeding now writes the official pinned
`@agentclientprotocol/claude-agent-acp@0.63.0` `npx` form. The migration merge
continues to replace only the selected `agents` members, so unrelated root
properties, unrelated agents, and their formatting survive. The existing
default-runtime validation gate remains ahead of the new multi-runtime
evaluation.

Focused checks used a disposable cache outside the worktree:

- Pre-change signal:
  `rtk env GOCACHE=/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260728T230436Z_3c8ceccb79af5226/.tmp/gocache go test ./internal/cli -run '^(TestRunSetupFreshMachineAcceptsOffers|TestRunSetupClaudeAdapterMigrationAcceptAndDecline|TestRunSetupMigratesBothStaleAdapterOverrides)$' -count=1`
  failed because Setup emitted no Claude migration prompt, proposed only the
  Codex migration, and did not seed the pinned Claude override.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli -run '^(TestRunSetupClaudeAdapterMigrationFailurePathsPreserveAllTargets|TestRunSetupClaudeAdapterMigrationAcceptAndDecline|TestRunSetupMigratesBothStaleAdapterOverrides)$' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli -run '^(TestRunSetup(AdapterMigration.*|ClaudeAdapterMigration.*|MigratesBoth.*|FreshMachine.*|MergesACPX.*|ProfileProofs.*|ProfileProofFailure.*|ProfileCleanup.*|ProfileWriteFailure.*|NoInputProfile.*|ReportsAdapterFailures.*))$' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli -run '^(TestRunSetup(HealthyMachineIsIdempotent|NewerACPXIsReadyWithoutInstall|ProfilePersistenceMatchesSubsequentValidation|ProfileProofUsesProposedProfilesAndWorkDir|AcceptsConfiguredEmptyReasoningEffort|RejectsMissingConfiguredAgentSelection|MismatchedACPXUpgradeOffer|DeclinedOffersReportNoWrites|NoInputSkipsOffers|ExitCodes)|TestSetupCommandCompatibility)$' -count=1`
  — passed.
- `rtk git -c core.fsmonitor=false diff --check` — passed.
- `rtk proxy rg -n '"claude-agent-acp"' internal/cli/setup.go` — exited `1`
  with no matches.

Acceptance evidence:

1. `TestRunSetupClaudeAdapterMigrationAcceptAndDecline` covers a legacy Claude
   override, a prompt naming `claude`, official-proposal re-proof, accepted
   persistence, and byte-identical decline.
2. `TestRunSetupFreshMachineAcceptsOffers` asserts the persisted Claude
   command is `npx` with the official package and `0.63.0` pin.
3. `TestRunSetupMigratesBothStaleAdapterOverrides` starts with stale Claude
   and Codex entries, observes one prompt for each runtime, and verifies both
   official pinned forms in one write.
4. `TestRunSetupAdapterMigrationPersistsSupportedCommand` continues to assert
   the official Codex package at the raised `1.1.5` pin and confirms the
   proposal was re-proved.
5. The Claude accept/decline table, the both-runtime case,
   `TestRunSetupMergesACPXAgentsOverridePreservingUnrelatedBytes`, and
   `TestRunSetupClaudeAdapterMigrationFailurePathsPreserveAllTargets` cover
   unrelated-byte preservation across accept, decline, `--no-input`,
   proposal-proof failure, and write failure.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
