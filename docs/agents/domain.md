# Domain Docs

How the engineering skills should consume this repo's domain documentation when exploring the
codebase.

## Before exploring, read these

- `CONTEXT.md` at the repo root
- `docs/adr/` at the repo root, especially ADRs that touch the area being changed

If a file does not exist, proceed silently. Do not suggest creating it upfront. Producer workflows
such as `grill-with-docs` create domain docs when terms or decisions are resolved.

## File structure

This repo uses a single-context layout:

```text
/
|-- CONTEXT.md
|-- docs/adr/
`-- packages/
```

## Use the glossary's vocabulary

When output names a domain concept in an issue title, refactor proposal, hypothesis, or test name,
use the term defined in `CONTEXT.md`. Do not drift to synonyms the glossary explicitly avoids.

If a needed concept is missing from the glossary, treat that as a signal. Either the output is
inventing language the project does not use, or there is a real gap to resolve with
`grill-with-docs`.

## Flag ADR conflicts

If output contradicts an existing ADR, surface it explicitly rather than silently overriding it:

> Contradicts ADR-0007 (Sync Operation Language), but worth reopening because...
