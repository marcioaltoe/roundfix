---
status: accepted
created_at: 2026-08-08T00:00:00Z
updated_at: 2026-08-08T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A capability projection bounds what it keeps, not what a runtime advertises

The capability projection refused any advertised select option carrying more
than 64 values and discarded the whole option rather than any part of it, so a
runtime that legitimately advertises a large catalog produced no model state at
all. Measured on 2026-08-08, `acpx opencode sessions show` advertised 417 model
values — 339 under `openrouter`, 60 under `opencode`, 18 under `opencode-go` —
and Roundfix answered `capability evidence invalid: contradictory_response,
missing_model_state, too_many_option_values`, which left no Exact Agent
Selection Proof and therefore no Run on the only runtime that had quota. The cap
was reading size as malformation, and size is exactly what a model aggregator
advertises. A projection therefore accepts a large advertised set and bounds
what it retains: the option's current value, every value the requested Agent
Selection references, and a bounded diagnostic window, with the retained count
and the advertised total both recorded. Fail-closed behavior is unchanged
because retention is by reference rather than by position — a model that is
genuinely not advertised is still absent from the retained set and still
produces `SelectionUnsupportedError` rather than a pass. Raising the cap to fit
today's catalog was rejected because the number stays arbitrary, moves the
failure to the next catalog that grows past it, and carries hundreds of values
into every diagnostic that reports advertised state.
