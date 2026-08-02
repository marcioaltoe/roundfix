---
status: open
created_at: 2026-07-30
updated_at: 2026-07-30
---

# `make baseline-digests` cannot bootstrap a Baseline module edit (2026-07-30)

Delivering `rule.autonomous.loop` through the Baseline module hit three
distinct failures in the derived-artifact contract. All three surface on any
module edit that changes a generated guide, which is the ordinary case, so they
block the autonomous Spec loop rather than being incidental to one change.

## 1. The regeneration command cannot refresh the pin that blocks it

- **Symptom / evidence**: after adding the rule to the module and the three
  profiles, `make baseline-digests` failed at its first step with

  ```text
  stage: ./internal/baseline:TestReadoptionCompatibilityMaintainedFixture
  catalog.profile.formatter.goldenDigest.mismatch: standard-typescript-monorepo:
    got 49337f03…, want 0d369576…; run 'make baseline-digests'
  ```

  Running the step that owns that pin, `TestFormatterComposition -update`,
  failed with the same diagnostic. The command's remediation is the command
  itself.

- **Root cause**: every entry in `BASELINE_DIGEST_STEPS` loads the embedded
  catalog, `catalog_validate.go:488` rejects a `goldenDigest` that disagrees
  with the goldens on disk, and only those steps rewrite it. A module edit
  changes the generated guide, which changes the goldens, which invalidates the
  pin, which refuses the load, which prevents the refresh. Spec 0054 ordered the
  steps so Source Baseline runs before the catalog; this is the same class of
  ordering defect one level deeper, and reordering cannot fix it because the
  cycle is within a single step.

- **Action / suggestion**: let the update mode defer or skip
  `catalog.profile.formatter.goldenDigest.mismatch` — a run that exists to
  rewrite derived pins must not be gated on those pins already being correct.
  The mismatch remains blocking outside update mode, where it belongs.

## 2. Adding a clause is not something the command can do

- **Symptom / evidence**: the first run reported
  `catalog.sourceBaseline.required-clause.missing` for all six new clauses. The
  regenerator maintains digests and byte ranges for manifest entries that
  already exist; it never adds one. The six manifest rows were written by hand
  with marker-computed spans before the command could proceed.
- **Root cause**: the Source Baseline manifest is treated as maintained input,
  not as derived output, while its digests and spans are derived.
- **Action / suggestion**: either generate the manifest rows from the corpus
  markers, or state in the failure's `nextSteps` that a new clause needs its
  manifest row added first. Today the message names the missing clause without
  saying that the tool cannot supply it.

## 3. The frozen parity corpus reads as a derived artifact

- **Symptom / evidence**: `TestPlanDeterminismMatchesMaintainedManagedEntryFixture`
  failed comparing the planned guide bytes against the parity corpus. The
  corpus lives under `internal/baseline/testdata`, which is listed in
  `DERIVED_DIGEST_PATHS`, so the natural reading is that regeneration owns it.
  Acting on that reading — repointing the blob and twelve fixtures at the new
  content — was wrong: `regenerateBaselineCompatibilityCorpus` rewrites only
  `fixtures/asset-sync.json`, and `baselineCompatibilityFrozenFields` fails if a
  frozen field moves. The edits were reverted.
- **Root cause**: a frozen parity record and genuinely derived artifacts share
  one directory and one `DERIVED_DIGEST_PATHS` entry, with nothing at the path
  saying which is which. Specs 0034 and 0035 handled the same collision by
  adding paths to an unexplained exclusion list in the test.
- **Action / suggestion**: mark the frozen corpus at the path — a `FROZEN.md` or
  a manifest field the test cites — so the next reader does not have to infer
  the boundary from a helper's implementation. This change added the reasoning
  as a comment on the exclusion set, which is a mitigation, not a fix.

## What worked — keep

- The structured failure JSON from Spec 0054 named the exact
  `package:test` stage every time, which is what made the cycle in finding 1
  provable instead of merely confusing.
- `make verify` caught two consequences the change had not anticipated: a
  wholesale reformat of three profile assets, and a mutation test keyed to
  incidental JSON adjacency that silently stopped mutating.

## Routing — 2026-08-01

Defects 1 and 2 shipped in [Spec 0062](../specs/_archived/0062-baseline-digest-regeneration-bootstrap/_prd.md).
Defect 3 — the frozen parity corpus reading as a derived artifact — was
explicitly excluded from that Spec's authorization and is routed to
[Spec 0067](../specs/0067-derived-artifact-regeneration-boundary/_prd.md).
