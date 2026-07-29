# Infra flow evidence

Build: `75161e9c3a5f7554cd1e0b9290bce6c61820b5c7`.

Disposable checkouts:

- `/private/tmp/roundfix-qa-0054-75161e9.2MCzTT/skill`
- `/private/tmp/roundfix-qa-0054-75161e9.2MCzTT/stale`
- `/private/tmp/roundfix-qa-0054-75161e9.2MCzTT/clause`

## Unchanged and Skill regeneration

Two current-build invocations of `rtk make baseline-digests` each passed all
five update selectors and printed:

```text
baseline-digests: no changes; derived artifacts already match their canonical sources
```

In the Skill checkout, the canonical and embedded Roundfix Skill copies
received the same probe edit. The target passed and named exactly:

```text
internal/baseline/assets/setups/typescript-bun.json
internal/baseline/testdata/catalog.digest
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json
internal/baseline/testdata/parity-corpus/v1/manifest.json
```

`rtk env GOCACHE=/private/tmp/roundfix-qa-0054-skill-gocache go test
-count=1 ./internal/baseline/ ./skills/` passed. A second target invocation
reported no changes. Git status contained only the Skill pair and the five
named deterministic fallout paths.

## Baseline module chain

The clause checkout changed one module clause and its delimited Source
Baseline corpus entry. `rtk make baseline-digests` passed and named:

```text
internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json
internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json
internal/baseline/assets/source-baselines/index.json
internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md
internal/baseline/assets/profiles/standard-typescript-monorepo.json
internal/baseline/testdata/catalog.normalized.json
internal/baseline/testdata/catalog.digest
```

The target therefore repaired the Source Baseline before any catalog consumer
loaded it. A second invocation reported no changes. The focused catalog,
parity, formatter, maintained-source, and corrupted-span validators passed.

## Stale and safety probes

A deliberately stale `catalog.digest` failed default comparison with:

```text
run 'make baseline-digests'
```

The target regenerated only `catalog.digest`; the same default validator then
passed and Git diff returned empty.

Fresh current-build safety selectors passed:

```text
TestAuthorialSkillSync
TestAuthorialSkillSyncUpdateModeRoundTrip
TestCatalogCompatibility
TestBaselineCompatibilityCorpus
TestFormatterComposition
TestReadoptionCompatibilityMaintainedFixture
TestSourceBaselineRegenerationRejectsCorruptedSpan
TestEmbeddedCatalog
```

The clause probe did not change the parity manifest. The Skill probe's parity
diff changed no `frozenDate`, `testCount`, or inventory-digest line, and the
parity validator passed.

## Cache, Skill, and bare-build observables

```text
env -u GOCACHE make -pn       -> GOCACHE = $(CURDIR)/.gocache
GOCACHE=/private/tmp/... make -pn -> exact exported value preserved
git check-ignore -v roundfix .gocache/probe -> /roundfix and /.gocache/
cmp canonical embedded Skill -> exit 0
make skills-sync-check -> 4 tests passed
```

With the repository-local cache supplied to the raw Go command, two
`go build -buildvcs=false ./cmd/roundfix` invocations exited `0`. Git status
before and after contained only QA artifacts, and `git check-ignore` resolved
the generated root binary to `.gitignore`'s `/roundfix` rule.
