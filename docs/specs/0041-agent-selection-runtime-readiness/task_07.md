---
task: task_07
spec: 0041-agent-selection-runtime-readiness
status: pending
type: backend
complexity: medium
---

# Task 07: Prove profile configuration before persistence

## Overview

Make `profiles configure` prove candidate Agent Selection Profiles before
confirmation or persistence. Structural fallback failures must remain
pre-proof, and every failure or declined confirmation must leave the selected
configuration scope byte-identical.

## Requirements

1. MUST validate the complete candidate profile and distinct Fallback Chain
   before creating a disposable Session.
2. MUST prove every distinct candidate Preferred Selection and fallback through
   the shared readiness coordinator before confirmation or write.
3. MUST run exact proof for interactive, `--file`, `--dry-run`, `--json`, and
   `--yes` inputs without allowing any mode to bypass readiness.
4. MUST report that one additional distinct authorized and proven selection is
   required when fallback configuration is missing or duplicates Preferred.
5. MUST leave User Config or Project Config byte-identical on validation or
   proof failure, cleanup failure, output failure, or declined confirmation.
6. MUST keep recommendations advisory and never insert a ranked model or
   reasoning value into the candidate.
7. MUST report `changed: false` in JSON for every non-persisted outcome.

## Subtasks

- [ ] Separate candidate construction from persistence.
- [ ] Run structural fallback validation before Agent Session work.
- [ ] Prove candidate tuples through the shared readiness coordinator.
- [ ] Align interactive, file, dry-run, JSON, and yes modes.
- [ ] Add distinct-fallback remediation without automatic substitution.
- [ ] Snapshot config bytes across every non-persisted outcome.

## Acceptance Criteria

- [ ] A complete proven candidate reaches confirmation and writes only after
      authorization.
- [ ] Missing or duplicate fallback fails before any disposable Session and
      names the required distinct proof or authorization.
- [ ] Unsupported adapter, model, or reasoning capability fails before
      confirmation and reports the shared classification.
- [ ] `--dry-run` performs exact proof but writes nothing.
- [ ] `--yes` bypasses confirmation only; it cannot bypass structural or
      runtime proof.
- [ ] Decline, proof failure, cleanup failure, and JSON output failure leave
      the target bytes unchanged and report `changed: false` when JSON applies.
- [ ] No recommendation is inserted automatically into Preferred Selection or
      the Fallback Chain.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/cli/profiles_configure.go`
- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/config/profile_config.go`
- interface: `internal/config/profiles.go`

## Verification

- `rtk go test ./internal/cli -run 'TestProfilesConfigure.*(Proof|Fallback|DryRun|JSON|Yes|Decline|NoMutation)' -count=1` — expected: all input modes prove candidates and every non-persisted path preserves config bytes.
- `rtk go test ./internal/config -run 'Test(WriteProfilesConfig|NormalizeProfilesFragment).*Fallback' -count=1` — expected: missing and duplicate fallback contracts remain deterministic and pre-proof.
- `rtk go test -race ./internal/cli -run 'TestProfilesConfigure.*(Proof|NoMutation)' -count=1` — expected: candidate proof and persistence boundaries are race-free.

## References

- `_prd.md` → User Stories 2, 3, and 8; Core Features 4, 9, and 10; Success
  Metrics.
- `_techspec.md` → Shared Readiness Coordinator; Setup Transaction Boundary;
  Error Taxonomy and Diagnostics; Build Order 7.
- `../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md`
  → exact advertised selection proof.

