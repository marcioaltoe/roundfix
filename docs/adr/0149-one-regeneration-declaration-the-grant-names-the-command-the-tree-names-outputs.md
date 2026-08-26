---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# One regeneration declaration: the grant names the command, the tree names its outputs

An authorization grant names only the regeneration command it sanctions; the audit resolves that command's outputs from `_ownership.yml`, where the repository already records which derived artifact each command owns — asking every grant to enumerate outputs is the enumeration of consequences ADR-0081 rejected, drifting the first time a command gains an output. The suite guard reads that same declaration instead of carrying a private exemption list, and exempts only those paths under only those commands: two lists would be free to disagree, and a guard with its own list is a checker that approves what the audit refuses.

Consolidates ADR-0128 and ADR-0129 (2026-08-26); both are archived under docs/history/adr/.
