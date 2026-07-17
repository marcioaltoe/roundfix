# Domain docs

Every pipeline stage reads the shared domain documentation before writing
Specs, code names, test names, prompts, or TUI copy.

## Before exploring, read these

- `CONTEXT.md` at the repository root.
- Relevant ADRs under `docs/adr/`.

If either location is absent, proceed silently. Domain-producing workflows add
terms and decisions when they are resolved.

## Layout

Roundfix uses one context:

```text
/
├── CONTEXT.md
├── docs/adr/
├── cmd/
└── internal/
```

## Use glossary vocabulary

Use the terms defined in `CONTEXT.md`. Do not replace them with synonyms the
glossary rejects.

If a required concept is missing, reconsider whether the output is inventing
language. Record a genuine vocabulary gap for `domain-modeling`.

## Flag ADR conflicts

Surface any conflict with an accepted ADR instead of silently overriding it.
