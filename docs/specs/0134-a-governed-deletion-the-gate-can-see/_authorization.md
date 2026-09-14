---
status: approved
granted: 2026-09-14
action: make a governed deletion visible to the gate and an enumerated regeneration list authoritative
consuming: 0134-a-governed-deletion-the-gate-can-see
paths:
  - internal/speccheck/mechanical_test.go
  - internal/suiteguardcontract/regeneration_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0134

On 2026-09-14 the maintainer answered the explicit request naming these two
exact paths with "Aprovar os dois de teste", approving the bounded scope above.

## Why a governed path is unavoidable

ADR-0130 holds every path an authorization has once bounded in the governed
class, and two of those paths are the test files that already carry the
contracts this Spec corrects: `internal/speccheck/mechanical_test.go` holds the
audit's changed-path and regeneration-output contract, and
`internal/suiteguardcontract/regeneration_test.go` holds the suite guard's
reading of the same record.

The defect is precisely that the two readers disagree — the audit unions an
enumerated `outputs` list with every owner-derived output for the same command,
restoring what the record excluded, while the suite guard treats a present list
as final. An assertion that the two agree has to be written where their
contracts are already asserted. Moving the new assertions to an ungoverned file
in the same package would split one package's test contract into an audited half
and an unaudited half, which is the evasion the monotonic set exists to prevent.

The production repair is ordinary source: `internal/speccheck/mechanical.go`,
`internal/daemon/task_engine.go` and `internal/cli/settle.go` are ungoverned and
need no grant.

## Approved bounded mutation

Add assertions to the two governed test files: that an owner-derived output
outside an enumerated list is refused, that a command-only declaration still
resolves ownership unchanged, and that the audit and the suite guard return the
same allowed set for the same record.

## Limits

- No action, operation or path beyond the two above.
- Assertions may be added; no existing assertion may be weakened or deleted
  except the two characterization outcomes this Spec declares it moves.
- No canonical Baseline module edit, so no derivative regeneration is
  sanctioned by this record.
- No paid API use, release, tag, deployment, or branch-policy exception.
- Verification remains Daemon-owned; Task status remains Daemon-written.

## Commit order

This record lands in `main` ancestry before the consuming Task commit. The
audit resolves a consuming commit's authorizing revision as the merge base of
the delivery target and that commit's parent, so a record committed on the
consuming branch is never the revision that authorizes it, and the squash
delivery would flatten the grant into the commit it must precede.
