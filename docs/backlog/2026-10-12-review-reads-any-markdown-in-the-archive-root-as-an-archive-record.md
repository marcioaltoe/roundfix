---
type: fix
status: open
created: 2026-10-12
spec: null
---

# `review` reads any Markdown file in the archive root as an Archive Record

## Symptom

On 2026-10-09 a pantheon PR archived Spec 0043 and also changed
`docs/history/specs/README.md`, a prose index with no frontmatter.
`roundfix review` (0.65.0) stopped before it called the provider:

```text
outcome: blocked
reason: review runtime failure: read review archive record: missing YAML frontmatter opening marker
```

The same branch was reviewed normally once the README change was removed.
Confirmed in code on 2026-10-12: `internal/cli/review.go` derives a slug from
every `.md` changed directly under an archive root and calls
`readArchivedSpecAt`. It does this in two places: review scope (around line
415) and Spec context (around line 1466).

## Where

`roundfix review`: `internal/cli/review.go`, review scope and Spec context.

## Expected

Treat a file as an Archive Record only when its name matches a Spec slug
(`NNNN-<slug>.md`) or its frontmatter declares
`schema: roundfix/archive-record/v1`, and ignore other Markdown files in that
root. Sweep every reader of the archive root for the same assumption.

## Evidence

- Secondbrain `inbox/roundfix/_triaged/2026-10-09-review-le-o-readme-da-raiz-de-arquivo-como-archive-record.md`
  (pantheon, 2026-10-09).
