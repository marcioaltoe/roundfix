# Regeneration evidence

The fresh Go 1.26.7 focused reader run exited 0 for:

- `TestEnumeratedOutputsAreAuthoritative`;
- `TestCommandOnlyDeclarationStillResolvesOwnership`;
- `TestAuditAndSuiteGuardAgreeOnAllowedOutputs`;
- `TestMechanicalAuthPathsAcceptsDeclaredRegenerationOutput`;
- `TestMechanicalAuthPathsStillRefusesAnUndeclaredPath`;
- `TestMechanicalAuthPathsRefusesInvalidRegenerationDeclaration`;
- `TestCleanupRegenerationDiscovery`;
- `TestOutputsForCommand`;
- `TestCleanupRegenerationOwnershipParity`.

The enumerated-list reader allowed the listed output and emitted
`QA-AUTH-PATHS` for the owner-derived sibling excluded from the enumeration.
The command-only reader resolved the repository-owned set. The mechanical audit
and `suiteguardcontract.ReadSanctionedRegenerations` returned the same allowed
set for both declaration shapes.

The outside source is the repository's Baseline ownership tree under
`internal/baseline/assets/**/_ownership.yml` and
`internal/baseline/testdata/**/_ownership.yml`, authored before Spec 0134.
`TestOutputsForCommand/make_baseline-digests_matches_the_2026-08-06_enumeration`
independently compared the current `baseline.OutputsFor` result with commit
`81a6afb48f4a3683d0e5fad52f3919cf1bdfbbf4`, path
`docs/workflow/authorizations/2026-08-06-proof-cost.md`, and passed. That
committed 2026-08-06 list is outside this Spec's artifacts.
