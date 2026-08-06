---
task: task_04
spec: 0059-run-storage-compaction-and-global-sanitation
status: completed
type: backend
complexity: medium
---

# Task 04: Give every durable table a stated lifecycle

## Overview

Tables added after the retention contract — Run summaries, Agent Selection
records — have no defined long-term lifecycle. They grow, and nothing states
who owns them or when they may go.

This slice writes the policy and makes it checkable, so a table added tomorrow
cannot quietly become permanent.

## Requirements

1. MUST state, for every durable table, its owner and its retention rule.
2. MUST preserve the compact Run index by default.
3. MUST keep Active Run locks governed by the Run lifecycle, not by retention.
4. MUST prune Agent Selection records only with their owning Run, or under an
   explicit evidence-retention rule the policy states.
5. MUST uphold the Spec 0014 promise never to delete `runs` rows or Active Run
   locks, unless this policy explicitly bounds a table with measured
   justification recorded alongside it.
6. MUST fail when a durable table exists with no stated policy, so the next
   table added cannot skip the decision.
7. MUST NOT change what retention deletes today.

## Subtasks

- [ ] Write the per-table owner and rule.
- [ ] Add the check that every durable table appears in the policy.
- [ ] Assert retention behaviour is unchanged.

## Acceptance Criteria

- [ ] Every durable table appears in the documented lifecycle policy with an
      owner and a rule.
- [ ] A durable table with no stated policy fails the check, asserted with a
      fixture table.
- [ ] The compact Run index is preserved by default.
- [ ] Agent Selection records prune only with their owning Run.
- [ ] `runs` rows and Active Run locks are not deleted, asserted.
- [ ] Retention behaviour is unchanged, asserted over the existing tests.

## Context

- interface: `internal/store/journal.go`
- interface: `internal/store/store.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/store -count=1 -run 'Lifecycle|Policy|Retention' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the lifecycle tests ran and passed.
- `go test ./internal/store -count=1` — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.

## References

- `_prd.md` → Core Feature 3; User Story 4; Success Metric 4.
- `_techspec.md` → Build Order 4.
- ADR-0033.

## Result

Implemented a checked Run Database lifecycle policy without changing a
retention mutation path. The policy names the lifecycle owner and retention
rule for `runs`, `active_run_locks`, `interactive_defaults`, `run_events`, and
`run_agent_selections`; it also records why SQLite-owned `sqlite_` metadata is
outside the Roundfix table policy. The configuration guide links the policy
from its existing Logs and retention section.

The store integration suite reads the documented policy as the contract and
compares it with every Roundfix-owned table in a freshly migrated SQLite
database. Its negative companion creates `lifecycle_fixture` and requires the
check to reject that unstated durable table. A retention scenario seeds Active
and terminal Runs with Agent Selection records, then proves Journal Retention
deletes only the eligible terminal Run Event while preserving both Run rows,
the Active Run lock, the Active Run journal, and all Agent Selection records.

Pre-change signal:

- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/store -count=1 -run "^TestDurableTableLifecyclePolicyCoversEveryTable$"'`
  — failed as expected because the lifecycle policy did not exist yet.

Focused checks run after the final implementation edits:

- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/store -count=1 -run "^(TestDurableTableLifecyclePolicyCoversEveryTable|TestDurableTableLifecyclePolicyRejectsUnstatedFixtureTable|TestRetentionPreservesRunLifecycleRecords|TestPruneTerminalRunsDeletesOnlyEligibleJournalRows|TestPruneTerminalRunsNoOpsWhenCutoffSelectsNothing)$"'`
  — passed (`5` tests).
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go test ./internal/store -count=1 -run "^$"'`
  — exited `0`; the package and tests compiled without selecting a test.
- `rtk zsh -c 'GOCACHE="$PWD/.gocache" rtk go vet ./internal/store'`
  — exited `0`.
- `rtk git -c core.fsmonitor=false diff --check` — exited `0`.

Acceptance evidence:

- Every durable table has an owner and rule: the positive lifecycle-policy
  test reconciles the five documented policy rows with the migrated schema and
  rejects empty owners, empty rules, duplicates, unknown policy rows, and
  undocumented Roundfix-owned tables.
- An unstated durable table fails: the negative test adds
  `lifecycle_fixture` and asserts the specific missing-policy diagnostic.
- Compact Run index preserved: the policy retains `runs` by default, citing
  the measured 279-row, 118,784-byte index, and the retention test proves the
  Run row count is unchanged after pruning.
- Agent Selection records follow their owning Run: the policy permits deletion
  only with that Run through the existing foreign-key lifecycle unless a
  future measured evidence-retention rule is stated; the retention test proves
  the current Journal Retention path leaves both Active and terminal Run
  selection records unchanged while their owning Run rows remain.
- Run rows and Active Run locks are not deleted: the new retention test records
  both counts before pruning and proves them unchanged afterward. The existing
  eligible-journal test independently reasserts that every Run row survives
  and the Active Run lock count is untouched.
- Retention behavior is unchanged: the existing eligible-only and no-candidate
  prune scenarios both pass unchanged, alongside the new preservation
  scenario. No production retention code was edited.

The commands under `## Verification` were not run; the Daemon owns those
commands and Task settlement.
