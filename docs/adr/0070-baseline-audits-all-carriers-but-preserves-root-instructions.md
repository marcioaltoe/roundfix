# ADR-0070: Baseline audits all carriers but preserves root instructions

Status: Accepted

Baseline adoption inventories every bounded instruction and agent-document
carrier, but automatic backup, classification, and mutation of pre-existing
instruction carriers are limited to the repository root. Root carriers receive
immutable content-addressed backups named `AGENTS.<digest>.md` or
`CLAUDE.<digest>.md`; a safe root alias preserves its target once. Nested
carriers remain unchanged and detected conflicts are warnings, preserving
complete audit evidence without expanding automatic remediation across the
repository.
