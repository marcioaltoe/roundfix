---
type: fix
status: promoted
created: 2026-10-27
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
reason: null
---

# The authoring templates and the tooling check disagree on no tooling change

## Symptom

With no protected tooling mutation, the `write-prd` and `write-techspec`
templates tell authors to record the Tooling authority row as
`applicable — no protected tooling mutation proposed or authorized`. With
`paths: []` in `_authorization.md`, `roundfix spec check --strict` reports
`SC-TOOLING-UNBOUNDED` ("declares a protected tooling mutation without bounded
files") on the PRD and the TechSpec. The row passes only as `not applicable —
no protected tooling mutation proposed or authorized`. On 0.124.0 the same
`applicable` row passed.

## Where

`.agents/skills/write-prd/references/prd-template.md:30` and
`.agents/skills/write-techspec/references/techspec-template.md:29`;
`internal/speccheck/constraints.go` (`declaresProtectedToolingMutation`, the
`CodeToolingUnbounded` finding).

## Expected

The template and the checker agree: either the checker reads the documented
"no protected tooling mutation" row as no mutation, or the templates prescribe
the wording the checker accepts. Find the change between 0.124.0 and 0.125.0
that moved the verdict.

## Evidence

Quaestor (ai-memory message `01a17e01-45fd-7a81-9052-54c660114746`,
2026-10-27 12:36 BRT), item 1:

```
1. SC-TOOLING-UNBOUNDED quando não há mudança de ferramenta (Specs 0028 e 0029). O modelo do write-prd e do write-techspec manda escrever, sem mudança de ferramenta, "Tooling authority: applicable — no protected tooling mutation proposed or authorized ... bounded files: none". Com `paths: []` no _authorization.md, o `roundfix spec check --strict` acusou "declares a protected tooling mutation without bounded files" no PRD e na TechSpec. Só passou trocando por "not applicable — no protected tooling mutation proposed or authorized". Na v0.124.0, a 0022 passou com "applicable". Esperado: o modelo e o checker concordarem.
```
