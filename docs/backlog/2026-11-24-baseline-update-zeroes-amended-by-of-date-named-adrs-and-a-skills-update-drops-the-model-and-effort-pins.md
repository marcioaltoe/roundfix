---
type: fix
status: open
created: 2026-11-24
spec: null
reason: null
---

# `roundfix baseline update` zeroes `amended_by` of date-named ADRs, and a skills update drops the `model` and `effort` pins

## Symptom

1. The Pantheon project names ADRs by date (`docs/adr/YYYY-MM-DD-slug.md`, no number). Its ADR `2026-11-01-enderecos-publicos-por-produto.md` has `amended_by: [2026-11-20-tela-do-fiscus-na-raiz-de-fiscus-gesttione-com]` and the body notice `Amended by [2026-11-20-...](2026-11-20-....md) (2026-11-20): ...`. On roundfix 0.152.0 the documentation part of the Baseline plan (sha256:f009d0b1..., Pantheon PR #340) changes the field to `amended_by: []`. Rewriting the notice without a link did not stop the plan from zeroing it. By `docs-layout.md` the reader recognizes only `Amended by ADR-NNNN (YYYY-MM-DD)`, so applying the plan erases the amendment link in any repository with date-named ADRs. It is the same family as `SC-ORDINAL-CLAIMED`, which also presumes a numbered ADR.
2. `bunx skills update -p -y` removes the `model` and `effort` pins of the authoring skills (`write-idea`, `write-prd`, `write-tasks`, `write-techspec`) on every refresh, and the repository must reapply them by hand.

## Expected

For 1: accept the ADR identifier by file name (a slug with a date), with or without a Markdown link, in both the notice and the field; or at least never zero an `amended_by` that the reader does not know how to fill. For 2: not stated by the sender; what the update should preserve is open.

## Source

Cross-project message `01a1f99f-0904-72e2-9679-dd59138343ce` from the Pantheon project, 2026-11-20, about 12:50 BRT. Nature: measured (plan hash and PR cited). Recorded by triage on 2026-11-24.
