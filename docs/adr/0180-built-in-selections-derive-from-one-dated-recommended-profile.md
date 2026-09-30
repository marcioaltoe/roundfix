---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Built-in selections derive from one dated Recommended Profile

Roundfix named its models in four separate places: the Model Catalog, the
recommendation ranking, the built-in Agent Selection Profiles with the
generated config, and the Baseline semantic analysis of ADR-0069. They drifted
apart twice. On 2026-09-30 the catalog still offered three retired Codex
models and no GPT-6 model, and the built-in fallback, the legacy Codex default
and the Baseline analysis fallback all named `gpt-5.5`, which leaves Codex on
2026-10-14. The ranking carried a DeepSWE result and an average cost for every
row, and no current leading model has a published figure for both.

Roundfix now ships one dated Recommended Profile per Agent Work Category: a
Preferred Selection and a Fallback Chain with a rationale each, under one
snapshot date. Everything else is derived from it or checked against it:

- The built-in profiles of the five required categories, the generated config
  and the legacy Codex runtime default are the Recommended Profile. The five
  optional categories keep inheriting `general` when absent (ADR-0107).
- The recommendation `profiles show` prints is the Recommended Profile itself,
  in order. The numeric ranking is retired.
- A test requires every recommended Codex or Claude model to be in the Model
  Catalog, and another requires the reference document to state the same
  snapshot date as the binary.
- The Baseline semantic analysis of ADR-0069 keeps its read-only, supervised
  contract. Its two selections become `gpt-6.1-sol` and then `gpt-5.6-sol`,
  both at `xhigh`.

## Consequences

- A model retirement is handled in one place: a new snapshot in a release.
- The built-in `review` profile is no longer a copy of `general`. It leads
  with a cheaper model and stays on one runtime, because every selection in it
  must use the pre-PR review provider's runtime.
- The recommendation carries a rationale and a date, not a benchmark figure.
  Figures and their sources live in the reference document, where a missing
  figure can be written as missing.
- A repository that configures its own profiles is unaffected until it chooses
  to adopt the snapshot.
