---
status: accepted
created_at: 2026-07-28T22:02:41Z
updated_at: 2026-07-28T22:02:41Z
deprecated_at: null
superseded_by: null
---

# Adopted sources move to their owning Spec

Inbox-to-finding-to-Spec promotion created links without transferring the
document, so archived Specs left stale links behind and their shaping
evidence lived in a separate lifecycle. A document a Spec adopts as an
implementation source therefore moves — one move with Git history, never a
copy or a stub — into that Spec's references, owned by exactly one primary
Spec, recorded in a Spec-local reference index, and validated by the
authoring and archive gates; secondary Specs link the owner's copy. The
workflow applies to new promotions only — history stays where it happened.
Duplicating content or leaving stub files was rejected because two
authoritative homes drift and stubs have no lifecycle.
