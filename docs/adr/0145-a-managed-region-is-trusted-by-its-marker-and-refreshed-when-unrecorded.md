---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# A managed region is trusted by its marker, and an unrecorded one is refreshed and named

Managed-region trust comes from the marker delimiting the region and the plan's preimage — not from the digest the Setup Manifest recorded on adoption day, which records what adoption wrote rather than what the region should contain after the catalog legitimately moves. A region whose bytes differ from the recorded digest is therefore refreshed rather than blocked: guidance inside the markers belongs to the Baseline, and proposing a change to it is an Inbox Entry to the Baseline's owner, not a local edit. What was wrong was never the replacement but the silence around it, so the plan names every unrecorded region and lists each on-disk line the refreshed rendering does not reproduce (bounded at 50, truncation counted).

Consolidates ADR-0101 and ADR-0102 (2026-08-26); both are archived under docs/history/adr/.
