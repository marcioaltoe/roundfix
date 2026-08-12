# R06 — outside evidence

The source independent of Spec 0080 resolves at
`docs/specs/_archived/0071-verification-cost/baseline/2026-08-03-after.md`.
Its frozen table records:

- unchanged-tree `make verify`: 4.9 seconds;
- complete fresh `go test ./... -count=1`: 88.9 seconds;
- light edit: about 5 seconds;
- `internal/cli` edit: 48.5 seconds;
- `internal/baseline` edit: 54.5 seconds.

Older findings under `docs/findings/_archived/` independently record a
roughly 20-minute gate cycle, 333 tool calls, 95 reasoning blocks, four QA
cycles on one Roundfix Spec, and seven on one Vortex Spec. The active rollup
`docs/findings/2026-08-06-rollup-qa-gates-and-verification-evidence.md`
contains nineteen pre-Spec observations.

The adopted source's exact 92/29/30-minute Spec 0079 numbers were not found in
the archived Spec's squashed Git history; they remain attributable to the
promoted backlog source, not independently reconstructed here. The outside
acceptance rests on the original Spec 0071 measurement and the older findings,
not on Spec 0080's own premise.
