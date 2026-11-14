---
status: accepted
created_at: 2026-07-24T21:27:41Z
updated_at: 2026-11-14T02:12:42Z
deprecated_at: null
superseded_by: null
---

# ADR-0070: Baseline audits all carriers but preserves root instructions

Baseline adoption inventories every bounded instruction and agent-document
carrier, but automatic backup, classification, and mutation of pre-existing
instruction carriers are limited to the repository root. Root carriers receive
immutable content-addressed backups named `AGENTS.<digest>.md` or
`CLAUDE.<digest>.md`; a safe root alias preserves its target once. Nested
carriers remain unchanged and detected conflicts are warnings, preserving
complete audit evidence without expanding automatic remediation across the
repository.

**Amended by [ADR-0399](0399-the-baseline-inventories-a-nested-carrier-only-when-git-lists-it.md) (2026-11-13).** A nested `AGENTS.md` or `CLAUDE.md` is inventoried only when `git ls-files --cached --others --exclude-standard` lists it; a carrier in a path Git ignores, or inside a linked worktree, a nested repository or a submodule, is not a carrier of this repository. Root carriers, their backups and the warning-only handling of the nested carriers Git lists still hold.
