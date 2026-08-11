---
status: accepted
created_at: 2026-08-06T21:23:58Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# A QA row carries forward only on declared, unmoved evidence

A QA Report row may be carried forward from an earlier report of the same Spec
instead of re-observed, and only when every one of these holds: the earlier
report established the row as `pass`; the row declared a non-empty typed list
of evidence inputs and every input is a `repository_path`; the earlier report's
head is an ancestor of the current head; no declared input appears in the
changed-path set between those heads; and every cited evidence path still
resolves with unchanged content. A carried row names the report and the head
that established it, so the record states which evidence is fresh and which is
inherited.

The rule fails closed in every direction that matters. A row that failed, was
blocked by any cause, was skipped, or declared no inputs is never carriable, so
carry-forward is opt-in per row rather than a default. The input kinds are
`repository_path`, `external_repository`, `live_service`, and `elapsed_time`.
A row whose truth depends on state outside the repository declares the matching
non-repository input and always re-observes; a mixed list of repository and
non-repository inputs is likewise never carriable. The alternative rules were
rejected for the same reason: recency alone would carry a row because it passed
lately rather than because nothing changed, and a whole-report carry would
inherit a verdict rather than an observation. Speed may never buy a stale pass;
ADR-0080 keeps verdict semantics, and no carried row may make a verdict more
permissive than a fresh observation would.
