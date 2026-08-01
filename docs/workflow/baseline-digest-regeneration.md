# Baseline digest regeneration contract

Use regeneration mode only through `make baseline-digests`. The mode is
unreachable outside that regeneration path; this operating boundary implements
[ADR-0085](../adr/0085-a-regeneration-run-is-not-gated-on-the-pins-it-rewrites.md)
without changing the ordinary validation contract.

## Boundary

`make baseline-digests` reaches regeneration mode through the Baseline
package's `-update` path and its unexported regeneration loader. Every other
catalog load stays strict, including production, CLI, CI, and the Verification
gate. Regeneration mode is not a general-purpose way to load an inconsistent
catalog.

## Deferred diagnostic allowlist

Regeneration mode defers exactly one diagnostic:

- `catalog.profile.formatter.goldenDigest.mismatch`

The allowlist is diagnostic-by-diagnostic; it does not defer a severity,
diagnostic category, or every derived-pin failure. After rewriting the derived
artifacts, `make baseline-digests` always performs a strict catalog
re-validation. Any remaining inconsistency fails the run, so a deferred pin is
checked later in the same run rather than ignored.

## Adding a Normative Clause

Add the new clause's Source Baseline manifest row by hand before regenerating.
The regenerator maintains digests and spans for existing manifest rows but
never creates rows. Without the row,
`catalog.sourceBaseline.required-clause.missing` names the clause and stops the
run.

## Recovering from the bootstrap cycle

The ordinary failure is an embedded-catalog load that reports
`catalog.profile.formatter.goldenDigest.mismatch` and tells the maintainer to
run `make baseline-digests`. Run the sanctioned target once:

```bash
rtk make baseline-digests
```

The target loads through regeneration mode, rewrites the stale pin, and closes
with strict re-validation. If the failure instead reports a missing Source
Baseline manifest row for a new clause, add that row by hand first, then rerun
the target.
