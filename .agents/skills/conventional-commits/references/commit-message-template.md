# Commit message template

```text
type(scope): imperative subject

Explain why the change exists when the reason is not obvious.
Keep the body focused on intent, trade-offs, and user impact.

Refs #123
```

Use `BREAKING CHANGE:` as a footer when the change breaks public behavior:

```text
feat(api): require tenant id in report routes

Report routes now require explicit tenant identity to avoid cross-tenant reads.

BREAKING CHANGE: report routes now require the `tenantId` path parameter.
```
