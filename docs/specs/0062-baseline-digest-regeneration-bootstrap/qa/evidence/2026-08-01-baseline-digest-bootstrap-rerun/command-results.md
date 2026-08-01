# Command results — Baseline digest bootstrap rerun

Build under test: `bb18c28981b22ae51e21d4b8ae246298edb6c29f`.

## Repository Verification

`rtk make verify` exited 0.

```text
rtk go test ./...
Go test: 2954 passed in 24 packages
rtk go test -count=1 ./skills -run 'TestNoPythonBaselineRuntime|TestThinSetupSkill|TestCheckRejectsExecutableSetupEngineArtifacts|TestRecommendedSkillsMatchLock'
Go test: 4 passed in 1 packages
rtk go run -buildvcs=false ./cmd/roundfix skills check
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec, write-tasks, setup-context-driven, implement-task, implement-spec, brainstorming, council, business-analyst, archive-spec, qa-gate, evidence-gate
rtk go build -buildvcs=false -ldflags "-X 'roundfix/internal/app.BuildCommit=bb18c28-dirty' -X 'roundfix/internal/app.BuildTime=2026-08-01 14:08:04 -0300'" -o bin/roundfix ./cmd/roundfix
```

The build metadata is dirty because the in-progress QA report and this evidence
artifact exist in the worktree; the implementation build is the recorded HEAD.

## Project Constraints and tooling authority

Both `_prd.md` and `_techspec.md` contain all four Project Constraints. They
classify identifier strategy and authentication/HTTP as inapplicable with
operative sources under `docs/agents/`, bind accepted ADR-0081 and ADR-0085,
and record the exact `Makefile` authorization at
`docs/workflow/authorizations/2026-08-01-baseline-digest-bootstrap.md`. ADR-0063
is inapplicable because the Spec has no transport. The two snapshots satisfy
ADR-0077.

The authorization/Spec commit preceded the tooling Task commit:

```text
e03fa07 2026-08-01T12:28:46-03:00 docs: spec the baseline digest regeneration bootstrap
0775ad4 2026-08-01T13:15:35-03:00 feat: re-validate strictly after regeneration
bb18c28 2026-08-01T14:04:14-03:00 feat: load the catalog once on the regeneration path
```

Both `git merge-base --is-ancestor e03fa07 0775ad4` and
`git merge-base --is-ancestor 0775ad4 bb18c28` exited 0.
`git diff-tree --no-commit-id --name-only -r 0775ad4` listed only:

```text
Makefile
docs/specs/0062-baseline-digest-regeneration-bootstrap/task_04.md
```

The six other Task commits contain no protected-tooling path. The authorization
was separate and earlier, no prerequisite or consequent tooling fix was folded
into the Task commit, and the three-line `Makefile` diff adds only the final
strict stage. No derived pin occurs in that commit. Current status contains
only this QA report and its evidence directory.

## Stale-pin precondition

A Git archive of exact build `bb18c28` was extracted under `/private/tmp`. The
probe appended the same sentence to the existing
`rule.backend.boundary-contracts` guidance and its generated formatter fixture,
leaving only the formatter digest pin stale. No Source Baseline identity was
added.

The first direct Go command reached the sandbox-denied host Go cache before
compilation. All scratch commands therefore used
`GOCACHE=/private/tmp/roundfix-qa0062-rerun.vh1jWo/<case>/.gocache`; this is the
only production-parity deviation.

The ordinary strict command then exited 1 before regeneration:

```text
catalog.profile.formatter.goldenDigest.mismatch:
standard-typescript-monorepo:
got 9bc59a439dd0f280aff81057e9ad57b663083424f7a5f6f6e941a1531e9ed3c5,
want 1f8e690949c4f49cdc31991d9727077b351bd89a87e815e46940b9cfe76eec25;
run 'make baseline-digests' to regenerate stale derived artifacts
```

## One-invocation regeneration retest

The first public target invocation against that fresh stale-pin archive exited
0 and rewrote the expected three derived artifacts:

```text
baseline-digests: regenerated
  internal/baseline/assets/profiles/standard-typescript-monorepo.json
  internal/baseline/testdata/catalog.digest
  internal/baseline/testdata/catalog.normalized.json
{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":true}
```

The ordinary strict `TestCatalogCompatibility` command then exited 0. A second
consecutive public invocation independently exited 0 with:

```text
baseline-digests: no changes; derived artifacts already match their canonical sources
{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":false}
```

This retest closes the previous report's F-01 on the current build: the first
invocation, not a partially rewritten retry, reaches the documented goal.

## Clean-tree idempotence and target shape

Two consecutive `make baseline-digests` invocations in a second, unmodified
archive each exited 0 with `"changed":false`. `rtk make -n baseline-digests`
showed the five `BASELINE_DIGEST_STEPS`, both `DERIVED_DIGEST_PATHS` scans, the
comparison, and ordinary `TestCatalogCompatibility` under
`stage: strict-validation` before the existing success output. The failure trap
sets `errorCode: strict_validation_failed`, `stage: strict-validation`, and
`retryable:false`.

After those clean runs, a deliberate formatter inconsistency made that same
ordinary strict command exit 1 with
`catalog.profile.formatter.goldenDigest.mismatch`. This confirms the command
used by the final stage still rejects a catalog that remains inconsistent.

## Focused contract checks

The following focused command exited 0 with 13 passing tests and subtests:

```text
go test ./internal/baseline \
  -run '^(TestCatalogDiagnosticCharacterization|TestCatalogRegenerationMode|TestRegenerationBreaksGoldenDigestCycle|TestRegenerationLoadsCatalogOnce|TestSourceBaselineManifestRowGuidance)$' \
  -count=1
```

`TestRegenerationLoadsCatalogOnce` runs the real plan/apply/re-plan/reapply path
with a filesystem acquisition counter, requires one catalog load across all
stages, requires the second plan to contain no changes, and finally requires a
strict load of the unrepaired fixture to emit the formatter mismatch. The
regeneration identifiers remain unexported and their callers are confined to
`internal/baseline/`. `TestCatalogRegenerationMode` asserts the allowlist has
exactly one member and proves a non-allowlisted unknown-module diagnostic
remains fatal.

A second ordinary characterization run exited 0. The golden remained
byte-stable before and after at SHA-256
`d11653720f83b82a45fcefad66e2ac3cd11fdc63185430f6004b3ab4d54381a4`.
Commit `a069e95` changes only three golden details under the two existing
missing-clause/missing-rule codes.

## Missing manifest row

In a third fresh archive, adding `clause.core.qa-gate-new-clause` to
`modules/core.json` without a Source Baseline manifest row made the public
target exit 2 at `TestReadoptionCompatibilityMaintainedFixture` with structured
`errorCode: regeneration_failed` output. The unfiltered failing step reported:

```text
catalog.sourceBaseline.required-clause.missing:
baseline.standard-typescript-monorepo-0.0.1:
clause.core.qa-gate-new-clause; the regenerator maintains manifest rows but
never creates them, so add this row first
```

A recursive search of `internal/baseline/assets/source-baselines` for that
identity returned no matches, confirming that the target created no manifest
row. The focused clause and rule guidance cases both passed with their original
codes and subjects.

## Documentation and scope

The operating document exists, its ADR-0085 target exists, and a focused scan
found the strict boundary, exact one-code allowlist, strict close, manual
manifest-row step, and `make baseline-digests` recovery command. Both documented
branches were exercised above: the stale-pin branch succeeded on its first
invocation and the missing-row branch stopped with the documented human action.

The committed Spec diff contains no Baseline module, Profile, Source Baseline,
retention, or frozen-parity-marker change. The missing-row flow created no row;
the public strict flow still rejected drift; and the full Verification gate
passed. The only protected-tooling path is the authorized `Makefile` change.

## Pull Request environment

The QA prompt states that no Pull Request is open for
`ma/baseline-digest-regeneration-bootstrap` and that Pull Request journeys are
environment-blocked. The prompt is supervised evidence for the cause. Local
equivalent evidence covers build and ancestry integrity through the passing
repository gate and the complete commit/path audit above. GitHub approval,
checks, unresolved threads, and Merge-Ready state do not exist to observe until
the target-branch Pull Request is opened; the per-Run branch was not queried.
