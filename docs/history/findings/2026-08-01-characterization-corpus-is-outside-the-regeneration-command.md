---
status: done
created_at: 2026-08-01
updated_at: 2026-09-08
absorbed_by: 0067-derived-artifact-regeneration-boundary
---

# The characterization corpus is a derived artifact outside its regeneration command (2026-08-01)

Spec 0062 shipped a diagnostic characterization corpus at `internal/baseline/testdata/catalog.diagnostics.golden.json`. The first owned-skill edit after it merged broke `make verify`, and `make baseline-digests` could not fix it. This is the same shape as the defect 0062 closed, one artifact over.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-01-characterization-corpus-is-outside-the-regeneration-command.md`.
