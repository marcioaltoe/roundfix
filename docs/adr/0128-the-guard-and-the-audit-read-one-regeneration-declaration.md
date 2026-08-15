---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-14T22:40:00Z
updated_at: 2026-08-14T22:40:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# The suite guard and the changed-path audit read one regeneration declaration

A sanctioned regeneration writes into the repository on purpose: `make
baseline-digests` and the characterization update flags exist to rewrite derived
artifacts, and the authorization records already declare exactly which paths each
command owns, in the `## Sanctioned regeneration` block the changed-path audit
reads. The suite guard therefore reads that same declaration instead of carrying
a list of its own, and exempts only those paths under only those commands. Two
lists would be free to disagree, and a guard with a private exemption list is one
edit away from silencing a real violation — which is the thing ADR-0126 exists to
make impossible. A path written by a command that declares it is sanctioned
fallout; a path written by anything else is a violation, in both readers.
