---
status: superseded # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-08T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null
superseded_by: ADR-0145
---

# A managed region is trusted by its marker, not by an adoption-day digest

Spec 0082 bound managed-region trust to the digest the Setup Manifest recorded
when the repository adopted the Baseline, so a region whose bytes moved was read
as hand-edited and blocked the refresh. That digest answers a different
question: it records what adoption wrote, not what the region should contain
after the catalog legitimately moves. Every sanctioned refresh advances the
bytes, and any advance that did not also republish the manifest — an older
Roundfix, a merge, a hand-maintained dogfooding repository — makes the
divergence permanent and unrecoverable, because the block fires before the plan
that would repair it exists. Measured on this repository, the manifest was
written once and the managed regions moved across seven subsequent pull
requests; bypassing the check produced a clean verified refresh whose entire
diff was the clauses authored that day, with no authored line lost. A managed
region is therefore trusted because its marker delimits it and the plan carries
its preimage, and its bytes are classified against what the catalog can render
now. Keeping the digest as the trust anchor was rejected because it cannot
distinguish a stale record from damaged content, and adding lineage tracking of
historical catalog renderings was rejected as a durable cost paid to answer a
question the preimage and the removed-line report already answer.
