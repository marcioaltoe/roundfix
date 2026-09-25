---
type: feat
status: open
created: 2026-09-25
spec: null
reason: null
---

# Approve Baseline history moves separately from managed refreshes

## Opportunity

`baseline update` ties a managed-byte refresh and the relocation of repository-owned ADRs into history under one Plan Digest, so an adopter cannot take one without the other; Fiscus postponed a whole plan because of this.

## Value

Adopters take managed refreshes without relocating decisions they have not reviewed, and path citations stop breaking silently.

## Shape

Separate approval (or digest) for `HistoryMoves`, or at least show their citation impact before approval; needs a design decision against ADR-0073 and ADR-0103.

Evidence: secondbrain `inbox/roundfix/2026-09-17-um-digest-so-prende-refresh-gerenciado-e-mudanca-de-documento-do-repositorio.md`.
