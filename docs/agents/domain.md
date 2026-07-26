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

<!-- setup-context-driven:begin id=guide.domain version=0.0.1 -->

# Domain docs

This repository uses a single `CONTEXT.md`.

- **mandatory**: Read the repository's selected domain context and relevant accepted ADRs before naming domain concepts or changing behavior. Flag conflicts instead of silently overriding repository decisions.

<!-- setup-context-driven:end id=guide.domain -->

<!-- roundfix:repository-rule:begin id=rule.3b1da2dff813927c9add235124eecbb89296b9d0399471c6799722d9ed8eb557 -->
- `CONTEXT.md` — the project glossary (vocabulary contract for code, docs,
  prompts, and TUI copy)

<!-- roundfix:repository-rule:end id=rule.3b1da2dff813927c9add235124eecbb89296b9d0399471c6799722d9ed8eb557 -->

<!-- roundfix:repository-rule:begin id=rule.d8b07ca0330a68aae05aa9ca33369ae6e78b7a66e3a8b50bc2769506f8258857 -->
- `docs/adr/` — accepted architectural decisions and the living contract;
  flag conflicts before overriding them

<!-- roundfix:repository-rule:end id=rule.d8b07ca0330a68aae05aa9ca33369ae6e78b7a66e3a8b50bc2769506f8258857 -->

<!-- roundfix:repository-rule:begin id=rule.a386a34b4e5d005bc9982eef99310ee74d2b8b8caeacc8c7ca8c877240cb0d09 -->
**ALWAYS** use canonical terms from `CONTEXT.md` in command names, help text,
issue titles, test names, and user-facing explanations. If the right term is
missing, call out the gap instead of inventing new language.


<!-- roundfix:repository-rule:end id=rule.a386a34b4e5d005bc9982eef99310ee74d2b8b8caeacc8c7ca8c877240cb0d09 -->

<!-- roundfix:repository-rule:begin id=rule.12c3a84f0277df7a6904c307bc2d7de469129ae284e28a8159d6c3fe2d95e84b -->
### Domain docs

This is a single-context repo: root `CONTEXT.md` plus ADRs in `docs/adr/`. See
`docs/agents/domain.md`.


<!-- roundfix:repository-rule:end id=rule.12c3a84f0277df7a6904c307bc2d7de469129ae284e28a8159d6c3fe2d95e84b -->
