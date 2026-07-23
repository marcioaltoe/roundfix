<!-- setup-context-driven:begin id=guide.docs-layout version=0.0.1 -->

# Docs layout

- **mandatory**: Give each documentation directory one job: `docs/_inbox/` for raw notes, `docs/adr/` for decisions, `docs/agents/` for agent guidance, `docs/design/` for design artifacts, `docs/findings/` for dated investigations, `docs/handoffs/` for session continuity, `docs/references/` for external pointers, and `docs/user-guide/` for human documentation. Preserve repository-authored extensions outside setup markers.

- **mandatory**: Use this copyable frontmatter contract for each findings file:

```yaml
---
status: pending # pending | partial | deferred | done
created_at: YYYY-MM-DD
updated_at: YYYY-MM-DD
---
```

- **mandatory**: Use `pending` when the finding is new and has no implementation Spec.

- **mandatory**: Use `partial` when a linked Spec covers only the selected implementation scope. Record the reason the remaining observations are unnecessary and link the covering Spec.

- **mandatory**: Use `deferred` only when the finding will not be implemented. Record the reason for deferral.

- **mandatory**: Set `status: done` as soon as the Spec is created and linked; do not wait for its Tasks, QA, archive, or release.

- **mandatory**: Treat findings as immutable history: append evidence and routing links as dated addenda instead of rewriting the original observations.

- **mandatory**: Update `updated_at` whenever status changes or an evidence addendum is appended; keep `created_at` as the document creation date.

<!-- setup-context-driven:end id=guide.docs-layout -->

<!-- setup-context-driven:begin id=guide.spec-docs-layout version=0.0.1 -->

# Spec docs layout

- **mandatory**: Keep `_idea.md`, `_prd.md`, `_techspec.md`, `_tasks.md`, Task files, and `qa/` evidence under the Spec folder. Archive only completed Specs with a passing QA verdict under `docs/specs/_archived/`.

<!-- setup-context-driven:end id=guide.spec-docs-layout -->
