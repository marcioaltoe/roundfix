---
type: fix
status: promoted
created: 2026-09-30
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
reason: null
---

# The backend guide renders one clause twice

## Symptom

The `backend` Baseline module lists the same guidance under two entries of one rule: `clause.backend.boundary-contracts` and a second entry whose identifier is the rule's own, `rule.backend.boundary-contracts`. Every adopter with the backend module gets the paragraph "Keep blocking, network, process, database, and daemon boundaries explicit…" twice in `docs/agents/backend.md`.

## Where

`internal/baseline/assets/modules/backend.json`, the last entry of `rule.backend.boundary-contracts`.

## Expected

The rule carries the paragraph once. No two clauses of the shipped Baseline have the same text, and a test refuses a catalog where they do.

## Evidence

The Secondbrain mirror `projects/conexus/mirror/docs/agents/backend.md`, read on 2026-09-30, holds the sentence twice. A scan of all clauses on `9e439dbb` found this as the only duplicated text.
