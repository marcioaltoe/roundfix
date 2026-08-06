# Tooling authorization — `rule.autonomous.loop` in the Baseline module (2026-07-30)

The maintainer authorized delivering the autonomous-loop discipline through the
Context-Driven Baseline so every adopting repository receives it from
`roundfix baseline` rather than from a hand-copied block. PR #43 recorded the
same rules in this repository's own `docs/agents/autonomous-work.md`; this
change makes them a Baseline product.

## Authorized paths

- `internal/baseline/assets/modules/autonomous-work.json` — add
  `rule.autonomous.loop` with its six clauses; bump the module and guide
  versions.
- `internal/baseline/assets/profiles/go-cli-tui.json`
- `internal/baseline/assets/profiles/rust-cli.json`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json` — add
  the rule to `requiredRules`.
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/`
  — the corpus entry and its manifest rows for the six clauses.
- `internal/baseline/preservation_test.go`,
  `internal/baseline/plan_test.go` — the assertions this change invalidates.

Every other rewritten pin is deterministic fallout of the edits above and is
owned by `make baseline-digests`.

## Bounded exception — the `goldenDigest` bootstrap

`internal/baseline/assets/profiles/standard-typescript-monorepo.json`
`formatter.goldenDigest` was updated outside the sanctioned command, because
the command cannot reach it. Every step in `BASELINE_DIGEST_STEPS` loads the
embedded catalog, catalog validation rejects a stale `goldenDigest`, and only
those steps refresh it — so a module edit that changes a generated guide makes
`make baseline-digests` fail with a diagnostic instructing the reader to run
`make baseline-digests`. The value was derived by recomputing
`portableFileDigest` over the profile's `formatter.fixturePaths`, not
transcribed from the diagnostic.

This is a product defect in the regeneration contract, not a property of this
change: it blocks every future Baseline module edit, including the autonomous
loop this rule describes. Reported in
`docs/findings/_archived/2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md`.
