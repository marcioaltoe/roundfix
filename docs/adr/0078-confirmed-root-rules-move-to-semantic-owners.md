---
status: accepted
created_at: 2026-07-26T10:45:10Z
updated_at: 2026-07-26T10:45:10Z
deprecated_at: null
superseded_by: null
---

# ADR-0078: Confirmed root rules move to semantic owners

After a Preservation Change Plan backs up and accounts for every unmarked root
instruction byte, apply removes those source bytes from the live `AGENTS.md`
and retains their operative meaning only in the confirmed semantic owner or
Repository-Specific Normative Rules. A later Preservation inventories only new
unmarked root bytes, creates at most one content-addressed backup for each new
root content identity, and otherwise produces no root migration change.

This decision supersedes ADR-0070's live-root preservation outcome after
confirmed redistribution. ADR-0070's bounded root-carrier mutation scope,
immutable backup requirement, safe-alias treatment, and warning-only handling
of arbitrary nested carriers remain active.
