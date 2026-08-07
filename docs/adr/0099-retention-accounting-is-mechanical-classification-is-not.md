---
status: accepted
created_at: 2026-08-07T00:00:00Z
updated_at: 2026-08-07T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Retention accounting is mechanical; instruction classification is not

ADR-0058 requires every Baseline transition to account for each prior managed
Normative Clause, and the only implementation that ever satisfied it re-ran
supervised semantic analysis over the whole root instruction corpus — so an
already-adopted repository paid an ACP turn and a rule-by-rule review to
refresh generated blocks whose repository-authored neighbours had not moved.
The two obligations are separable: retention accounting compares prior managed
clauses against the current catalog's clauses, and both sides are generated
artifacts with recorded identities, so it is a mechanical comparison that stays
fail-closed without a model; instruction classification decides where
repository-authored prose belongs, which is a judgment and keeps requiring a
supervised analyzer. An update therefore performs retention accounting and
skips classification, because repository-authored bytes already sit in their
canonical carriers and the update leaves them byte-identical. Re-running
classification on every update was rejected because it spends a model turn to
re-derive a settled answer and invites drift in artifacts nobody asked to
change; dropping retention accounting alongside it was rejected because that is
exactly the unaccounted-removal failure ADR-0058 exists to block.
