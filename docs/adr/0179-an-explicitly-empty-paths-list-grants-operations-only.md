---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An explicitly empty paths list grants operations only

An approved authorization record carried two duties from Spec 0119: bound the
Governed Paths a Spec may change, and, since the Delivery Queue, grant the
operations `implement`, `commit`, `push`, `pull_request` and `merge`. The
reader refused a record whose `paths` held no entry, and ADR-0130 requires
every bounded path to be governed. A Spec that changes no Governed Path could
therefore never record a valid delivery grant: an empty list was refused by the
reader, and ordinary files listed instead were refused by the governed-set
contract. Spec 0186 met exactly that on 2026-09-30.

A Spec-contained record whose `paths` is an explicit empty YAML sequence
(`paths: []`) is now an operative grant of its listed operations that bounds no
Governed Path. An absent or null `paths` is still refused, so omission never
reads as a deliberate grant. RFC 6749 §3.3 makes the same distinction between
an omitted scope and a declared one.

## Consequences

- A governed mutation under such a grant still fails every audit that already
  checks the bounded set: the Spec checker's tooling detectors, Implement
  preflight and the QA mechanical authorization audit. The set is empty, so
  nothing is bounded.
- ADR-0130 holds unchanged: every listed path must be governed, and an empty
  list lists nothing.
- Legacy records keep their own reading.
