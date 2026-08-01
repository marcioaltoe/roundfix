---
status: open
created_at: 2026-08-01
updated_at: 2026-08-01
---

# The characterization corpus is a derived artifact outside its regeneration command (2026-08-01)

Spec 0062 shipped a diagnostic characterization corpus at
`internal/baseline/testdata/catalog.diagnostics.golden.json`. The first owned-skill
edit after it merged broke `make verify`, and `make baseline-digests` could not
fix it. This is the same shape as the defect 0062 closed, one artifact over.

## Symptom / evidence

Editing `.agents/skills/write-tasks/SKILL.md` (an `OWNED_SKILLS` member) and
running the sanctioned sequence:

```text
make skills-sync        # ok
make baseline-digests   # {"ok":true,"changed":true} — one invocation, as 0062 intended
make verify             # FAILS
```

with:

```text
catalog_test.go:616: catalog diagnostic characterization differs;
  run with -update-catalog-diagnostics
changed diagnostic codes: catalog.setup.digest.mismatch
- catalog.setup.digest.mismatch subject="rust-cli" detail="got 0f03942f…"
+ catalog.setup.digest.mismatch subject="rust-cli" detail="got 0feba2a0…"
```

Exactly one recorded line changed. `make baseline-digests` had already
succeeded and reported `changed:true`; it does not touch the corpus.

## Root cause

The corpus lives under `internal/baseline/testdata`, which **is** listed in
`DERIVED_DIGEST_PATHS`, so the natural reading is that `make baseline-digests`
owns it. It does not: the corpus is regenerated only by its own dedicated
`-update-catalog-diagnostics` flag, and
`TestCatalogDiagnosticCharacterization` is not a member of
`BASELINE_DIGEST_STEPS`.

The coupling is not incidental. One recorded fixture is a setup digest
mismatch whose *detail* embeds the observed digest, so the corpus changes
whenever any owned skill's content changes — which is the ordinary case for a
skill edit, exactly like the guide-edit case 0062 was built for.

This is the third instance of the same collision: a directory in
`DERIVED_DIGEST_PATHS` holding artifacts with different regeneration owners,
with nothing at the path saying which is which. Specs 0034 and 0035 handled it
with an unexplained exclusion list; the 2026-07-30 finding's third defect
(the frozen parity corpus) is still open for the same reason and was
explicitly excluded from 0062's authorization.

## Action / suggestion

Two options, and they are not exclusive:

1. Add `./internal/baseline:TestCatalogDiagnosticCharacterization` to
   `BASELINE_DIGEST_STEPS` so the one sanctioned regeneration command covers
   every derived artifact under the paths it scans. This is the smaller change
   and restores the property 0062 established — that a maintainer runs one
   command.
2. Make the corpus's failure message name the sanctioned command rather than
   only the raw test flag, so the printed remediation matches the documented
   workflow instead of sending the reader around it.

Both require the `Makefile`, which is protected tooling and needs express
authorization with bounded files. Until then the manual step is:

```bash
go test ./internal/baseline -run '^TestCatalogDiagnosticCharacterization$' \
  -update-catalog-diagnostics -count=1
```

## What worked — keep

- The corpus caught a real diagnostic change on its first live exercise outside
  its own QA, named the changed code, and printed both digests. It behaved
  exactly as Spec 0062 designed it.
- `make baseline-digests` completed in **one invocation** with `changed:true`,
  confirming 0062's fix on the first real edit that would previously have hit
  the bootstrap cycle.
