---
status: accepted
created_at: 2026-10-30T15:00:00Z
updated_at: 2026-10-30T15:00:00Z
deprecated_at: null
superseded_by: null
---

# A derived-record writer names its command to the suite guard

On 2026-10-27 Pull Request #709 conflicted with `main` only in derived
records. The Delivery Queue merged the default branch and ran each matched
regeneration command. The coverage record command exited 1 with `suiteguard:
modified: docs/references/coverage-record.json`, and the item parked
`delivery-error`. The operator read the refusal as "the archived Spec no
longer sanctions the command". The tree said otherwise: the suite guard
already reads the `regeneration` list of every Archive Record in the tree, as
ADR-0247 placed it there, and the item's own record listed the command with
that output. The guard exempts a write only under the command the running
test names, and the coverage record writer in `internal/spec` never named
one. Its write was refused under an active authorization as well. Each
earlier re-record wrote the right bytes and then exited 1.

## Decision

**Every test that writes a derived record names its command.** A test in a
package that installs the suite guard, and that rewrites a repository file
when its record flag is set, calls `suiteguard.DeclareSanctionedRegeneration`
with the exact command that runs it before it writes, and only while that
flag is set. The module version record and the Baseline digest writers
already do this. The coverage record writer joins them. A run without the
flag names no command, so a write it makes is still refused.

**After archive, the item's Archive Record is the sanction.** A regeneration
the Delivery Queue runs after the archive commit is sanctioned by the
`regeneration` list the archive copied from the item's authorization record
into its Archive Record. The suite guard reads that list from the tree where
the command runs, as it reads an active Spec's authorization record. Nothing
new is read.

**The guard is not widened.** It still exempts only the outputs a record lists
for the command the running test names, as ADR-0149 has it. No Project Config
declaration sanctions a write, and the queue's derived merge never turns the
guard off.

## Considered Options

- Restore the item's `_authorization.md` from the Archive Record's
  `source_revision` before regenerating, as the operator did by hand. Rejected:
  the guard would need to read Git objects for one item, and the record
  already carries the same commands and outputs. The restored record would not
  have helped anyway, because the writer named no command.
- Sanction the outputs of the queue's own Derived Path Declarations.
  Rejected: ADR-0192 says such a declaration "grants no authority". Every test
  of a matched package could then write those paths, which is a blanket
  allowance.
- Skip the guard during the queue's derived merge. Rejected for the same
  reason. An unintended write by a test would land in the merge commit.
- Keep the declaration injected by the regeneration contract test's fixture.
  Rejected: that injected declaration is why the contract passed while the
  real writer named nothing.

## Consequences

The coverage record command exits 0 when an authorization record or an
Archive Record in the tree lists it with `docs/references/coverage-record.json`.
It still exits 1 when no record does. The guard keeps reading every record the
tree holds. A command one archived Spec sanctioned therefore stays sanctioned
for later trees, which is the set that only grows under ADR-0149 and ADR-0247.
This decision does not narrow that set. A writer in a package without the
suite guard (`./skills`, `./internal/docscontract`, `go run ./cmd/roundfix`)
is not refused by it and needs no declaration. This refines ADR-0149, whose
exemption needed a command that one writer never named. It relies on ADR-0247's
record and supersedes neither. ADR-0081 still draws the authorization around
the edit that causes the regeneration.
