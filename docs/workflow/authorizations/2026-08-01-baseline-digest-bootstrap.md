# Tooling authorization — Baseline digest regeneration bootstrap (2026-08-01)

The maintainer authorized Spec 0062 to break the cycle that stops
`make baseline-digests` from refreshing the derived pins it exists to refresh.
The finding
`docs/findings/_archived/2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md`
records three defects; this authorization covers the two that block the
autonomous Spec loop.

## Authorized paths

- `Makefile` — the `baseline-digests` target, `BASELINE_DIGEST_STEPS`, and
  `DERIVED_DIGEST_PATHS`, only as far as breaking the regeneration cycle and
  reporting an unaddable manifest row require.

No other protected tooling mutation is authorized. Go sources under
`internal/` and fixture data under `internal/baseline/testdata/` are product
code and test data, not protected tooling, and are edited under the ordinary
implementation contract. Deterministic digest fallout of an authorized edit is
sanctioned by ADR-0081.

## Explicitly not authorized

The maintainer declined a marker file at the frozen parity corpus path. Defect
3 of the finding — the frozen corpus reading as a derived artifact because it
shares a directory and a `DERIVED_DIGEST_PATHS` entry with genuinely derived
output — stays out of Spec 0062's scope. The comment added on the exclusion set
by the earlier change remains its mitigation.
