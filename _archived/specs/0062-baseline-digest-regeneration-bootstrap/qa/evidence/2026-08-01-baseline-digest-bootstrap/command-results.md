# Command results — Baseline digest regeneration bootstrap

Build: `76709c563468623d6dbc33f8d09e361e1645861f`

## Repository gate

`rtk make verify` exited 0.

```text
Go test: 2953 passed in 24 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed: 14 required skills
go build -buildvcs=false ... -o bin/roundfix ./cmd/roundfix
```

## Stale-pin strict precondition

An archive of the exact build was extracted under `/private/tmp`. The probe
changed the existing `rule.backend.boundary-contracts` guidance and the
corresponding generated `docs/agents/backend.md` formatter fixture, leaving the
formatter pin stale. No Source Baseline identity was added.

`go test ./internal/baseline -run TestCatalogCompatibility -count=1` exited 1.

```text
catalog.profile.formatter.goldenDigest.mismatch:
standard-typescript-monorepo:
got 9bc59a439dd0f280aff81057e9ad57b663083424f7a5f6f6e941a1531e9ed3c5,
want 964a7ea17fc6d4be934ea2c2bd5b1fee30d0cd23eac13db4d2edbde171583a97;
run 'make baseline-digests' to regenerate stale derived artifacts
```

## One-invocation regeneration reproduction

First `make baseline-digests` invocation on that state exited 2.

```text
TestFormatterComposition regenerated:
  assets/profiles/standard-typescript-monorepo.json
  testdata/catalog.normalized.json
  testdata/catalog.digest

ApplyPlan() formatter composition error = validate supplied Baseline Plan:
load Baseline catalog for plan validation: load embedded Baseline catalog:
catalog.profile.formatter.goldenDigest.mismatch:
got 9bc59a439dd0f280aff81057e9ad57b663083424f7a5f6f6e941a1531e9ed3c5,
want 964a7ea17fc6d4be934ea2c2bd5b1fee30d0cd23eac13db4d2edbde171583a97

{"schemaVersion":1,"type":"baseline-digests","ok":false,
"changed":false,"errorCode":"regeneration_failed",
"stage":"./internal/baseline:TestFormatterComposition","retryable":false,
"nextSteps":"Read the failing test output above, fix the canonical source it validates, then rerun make baseline-digests."}
```

A second full invocation on the state partially rewritten by the failed first
run exited 0.

```text
baseline-digests: no changes; derived artifacts already match their canonical sources
{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}
```

The same first-invocation failure reproduced from a second fresh archive.

## Clean-tree idempotence and target shape

Two consecutive `make baseline-digests` invocations in an unmodified archive
both exited 0 with:

```json
{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}
```

`rtk make -n baseline-digests` showed the unchanged five-entry
`BASELINE_DIGEST_STEPS`, both `DERIVED_DIGEST_PATHS` scans, the comparison,
then ordinary `TestCatalogCompatibility` under `stage: strict-validation`
before the existing success output.

## Missing manifest row

Adding `clause.core.qa-gate-new-clause` to `modules/core.json` without a Source
Baseline manifest row caused the public `make baseline-digests` target to exit
2 at `TestReadoptionCompatibilityMaintainedFixture`. The unfiltered failing
step reported:

```text
catalog.sourceBaseline.required-clause.missing:
baseline.standard-typescript-monorepo-0.0.1:
clause.core.qa-gate-new-clause; the regenerator maintains manifest rows but
never creates them, so add this row first
```

No manifest row was created.

## Focused contract checks

The following focused command exited 0 with 12 passing tests and subtests:

```text
go test ./internal/baseline \
  -run '^(TestCatalogDiagnosticCharacterization|TestCatalogRegenerationMode|TestRegenerationBreaksGoldenDigestCycle|TestSourceBaselineManifestRowGuidance)$' \
  -count=1
```

A second ordinary characterization run exited 0. The characterization golden
remained byte-stable at SHA-256
`d11653720f83b82a45fcefad66e2ac3cd11fdc63185430f6004b3ab4d54381a4`.

Source/caller inspection found both regeneration entry points only under
`internal/baseline/`; both identifiers are unexported. The allowlist contains
one code. Commit `a069e95` changes only three characterization details under
the existing missing-clause/missing-rule codes.

## Governance and scope

Authorization/Spec commit `e03fa07` is an ancestor of Task 04 commit
`0775ad4`. `git diff-tree --no-commit-id --name-only -r 0775ad4` listed only:

```text
Makefile
docs/specs/0062-baseline-digest-regeneration-bootstrap/task_04.md
```

The remaining five Task commits contain no protected-tooling path. The full
Spec diff contains no change under Baseline module, profile, or Source Baseline
asset directories and adds no frozen-parity marker.
