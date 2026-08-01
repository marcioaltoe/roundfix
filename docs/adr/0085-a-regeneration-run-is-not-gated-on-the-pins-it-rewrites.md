---
status: accepted
created_at: 2026-08-01T00:00:00Z
updated_at: 2026-08-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A regeneration run is not gated on the pins it rewrites

Every step of `make baseline-digests` loads the embedded catalog, and catalog
validation rejects a formatter `goldenDigest` that disagrees with the goldens on
disk — but only those steps rewrite that pin, so an ordinary Baseline module
edit changes the generated guide, which changes the goldens, which invalidates
the pin, which refuses the load, which prevents the refresh. The command's own
remediation is the command itself, and no ordering of steps can fix it because
the cycle lives inside a single step. Catalog validation therefore gains an
explicit regeneration mode, enabled only by the update path, in which a derived
pin that disagrees with its regenerated source is deferred rather than fatal;
every other load, including CI and the Verification gate, stays strict and
still fails closed on the same mismatch. Relaxing the check globally was
rejected because the pin exists to catch exactly this drift outside
regeneration, and reordering the steps was rejected because it cannot break a
cycle contained in one step.
