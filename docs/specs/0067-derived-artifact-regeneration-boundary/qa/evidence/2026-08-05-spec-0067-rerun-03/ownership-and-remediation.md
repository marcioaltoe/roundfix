# Ownership, exception, and remediation evidence

Build: `9e668cfe4649c323694c55f8124ed73840260910`.

## Focused boundary suite

The documented Go entry point ran against isolated repository fixtures:

```text
rtk env GOCACHE=/private/tmp/spec0067-qa-rerun03-gocache \
  GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1 \
  -run '^(TestDerivedOwnership.*|TestMeasuredSanctionedOwnershipMatchesRecords|TestDeclaredStepRegenerationAndFrozenBoundaries)$' -v
exit 0
```

The selected suite passed the real exhaustive scan, unclassified and multiple
record rejection, exception precedence, duplicate and outside-directory
exception rejection, schema validation, whole-tree preservation, measured
sanctioned ownership, synthetic dedicated repair, nonexistent and wrong
commands, successful no-op rejection, sanctioned exception restoration, and
frozen-resolved-path rewrite rejection. The real scan and snapshot checks ran
after the negative probes.

## Real parity boundary

`internal/baseline/testdata/parity-corpus/` contains one ownership record and
17 artifact files. The record declares the directory `frozen` with exactly two
sanctioned exceptions:

- `v1/manifest.json`
- `v1/fixtures/asset-sync.json`

The focused known-boundaries check resolves the other 15 artifact files as
frozen. The scratch owned-Skill journey changed only the two exception files
under the parity directory.

The Spec task prose is stale: `_tasks.md` and `task_07.md` say the directory
contains 15 files and 13 frozen files. The real artifact count is 17, of which
15 are frozen. See finding F-003 in the report.

## Public frozen-failure remediation

QA changed `matrix.json` only in the scratch clone and ran the documented strict
test:

```text
rtk env GOCACHE=/private/tmp/spec0067-qa-rerun03-gocache \
  GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1 \
  -run '^TestBaselineCompatibilityCorpus$'
exit 1, as expected for the negative probe
```

The failure named `matrix.json`, said `nothing regenerates this artifact`,
named `testdata/parity-corpus/_ownership.yml`, and printed the recorded
2026-07-30 tried-and-reverted reason. It did not suggest
`make baseline-digests`. QA restored the byte and reran the same command; it
exited 0, and `git diff --quiet` confirmed the scratch matrix matched the build.
