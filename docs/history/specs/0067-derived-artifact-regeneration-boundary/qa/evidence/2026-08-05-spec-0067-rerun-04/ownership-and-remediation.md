# Ownership, exception, and remediation evidence

Build: `e91bf4088b7547ab1f1c4a15c78d1427e769f032`.

## Focused boundary suite

The documented Go entry point ran against isolated repository fixtures:

```text
rtk env GOCACHE=/private/tmp/spec0067-qa-rerun04-gocache \
  GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1 \
  -run '^(TestDerivedOwnership.*|TestMeasuredSanctionedOwnershipMatchesRecords|TestDeclaredStepRegenerationAndFrozenBoundaries)$' -v
exit 0
ok roundfix/internal/baseline 66.047s
```

The selected suite passed the real exhaustive scan, unclassified and multiple
record rejection, exception precedence, duplicate and outside-directory
exception rejection, schema validation, whole-tree preservation, measured
sanctioned ownership, synthetic dedicated repair, nonexistent and wrong
commands, successful no-op rejection, sanctioned exception restoration, and
frozen-resolved-path rewrite rejection. The real scan and snapshot checks ran
after the negative probes.

## Real parity boundary and corrected documentation

The reproducible tracked-tree count ran from the build:

```text
rtk git ls-tree -r --name-only HEAD -- internal/baseline/testdata/parity-corpus \
  | rtk proxy rg -v '/_ownership\.yml$' | rtk proxy wc -l
17
```

The record declares the directory `frozen` with exactly two sanctioned
exceptions, `v1/manifest.json` and
`v1/fixtures/asset-sync.json`. The known-boundaries check resolves the other
15 artifacts frozen. `_tasks.md` and `task_07.md` now state the same measured
17/2/15 shape and name `git ls-tree -r` as the count source. The focused suite
is the adjacent canary for the documentation-only correction.

## Public frozen-failure remediation

QA changed only `matrix.json` in the scratch clone and ran the documented
strict test:

```text
rtk env GOCACHE=/private/tmp/spec0067-qa-rerun04-gocache \
  GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1 \
  -run '^TestBaselineCompatibilityCorpus$'
exit 1, expected for the negative probe
```

The failure named `matrix.json`, said `nothing regenerates this artifact`,
named `testdata/parity-corpus/_ownership.yml`, and printed the recorded
2026-07-30 tried-and-reverted reason. It did not suggest
`make baseline-digests`. QA restored the byte and reran the same command; it
exited 0, and `git diff --quiet` confirmed that scratch `matrix.json` matched
the build again.
