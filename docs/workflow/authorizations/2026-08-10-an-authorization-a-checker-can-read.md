---
granted: 2026-08-10
action: add-spec-check-detector
paths:
  - internal/speccheck/constraints.go
  - internal/speccheck/coherence.go
  - .agents/skills/qa-gate/SKILL.md
consuming: direct
---

# Tooling authorization — an authorization a checker can read (2026-08-10)

On 2026-08-10, after the Secondbrain evaluation of the greenfield Baseline, the
maintainer directed:

> Faça os ajustes indicados pelo brain.

The adjustment this record covers comes from
`wiki/sources/arxiv-csai-agent-software-reliability-2026-08-08.md:96`:

> uma autorização de ferramenta deve ser um artefato tipado e enumerável —
> ação, recurso, escopo, estado necessário, efeitos, orçamento, replay key,
> preconditions e postconditions — e não uma frase no prompt.

## What this covers

Every authorization record in this repository states its grant in prose. The
`SC-TOOLING-UNAUTHORIZED` detector can only confirm that the record mentions the
Spec, and `SC-TOOLING-UNBOUNDED` can only confirm that the *citing artifact*
claims bounded files. Neither can read what the grant actually permits, so the
bounded-path audit stays a human reading task — which is exactly the QA gate row
that has cost the most rereading.

`SC-TOOLING-UNTYPED` requires an authorization record to open with frontmatter
carrying `granted`, `action`, `paths`, and `consuming`. A record whose filename
dates it before 2026-08-10 is historical evidence and stays byte-identical; the
detector exempts it by date rather than by content, so no past record is
rewritten to satisfy a rule it predates.

This record is the first typed one, and its own frontmatter is what the new
detector requires.

## Authorized paths

- `internal/speccheck/constraints.go`, limited to the `SC-TOOLING-UNTYPED` code,
  its detector, and the frontmatter reader it needs.
- `internal/speccheck/coherence.go`, limited to registering that code against
  the `prd` stage.
- `.agents/skills/qa-gate/SKILL.md`, limited to adding the new code to the row
  that maps authoring rules onto `roundfix spec check`.

The generated copy under `skills/qa-gate/SKILL.md`, rewritten by
`make skills-sync`, is sanctioned fallout under ADR-0081 rather than a separate
target. Test files and fixtures under `internal/speccheck/testdata/` are test
scaffolding for the same change, not repository tooling.

## Bounded by purpose

This grant covers making a tooling authorization machine-readable going forward.
It does not authorize rewriting any existing authorization record, changing what
the tooling-authority rule requires, or widening any other detector.

## Consuming Spec

Applied directly rather than through a Spec: it adds one detector beside two
that already read the same record, and it leaves every historical record
untouched.

## Commit choreography

This record lands as its own commit, before the commit that adds the detector.
