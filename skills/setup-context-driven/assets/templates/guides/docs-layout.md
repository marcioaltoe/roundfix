# Docs layout

Every `docs/` folder has one job. Put raw notes in `docs/_inbox/`, decisions in
`docs/adr/`, agent-facing guides in `docs/agents/`, design artifacts in
`docs/design/`, dated investigations in `docs/findings/`, session handoffs in
`docs/handoffs/`, external pointers in `docs/references/`, and human-facing
docs in `docs/user-guide/`.

Move files when their job changes. Keep agent-facing operational guidance in
`docs/agents/` so generated root instructions can stay compact.

## Findings lifecycle

Use this frontmatter for every dated finding:

```yaml
---
status: pending # pending | partial | deferred | done
created_at: YYYY-MM-DD
updated_at: YYYY-MM-DD
---
```

- `pending`: the finding is new and has no implementation Spec.
- `partial`: a Spec covers only the selected portion because the remaining
  observations were judged unnecessary; record the reason and link the Spec.
- `deferred`: the finding will not be implemented; record the reason.
- `done`: an implementation Spec exists. Set `status: done` as soon as the Spec
  is created and linked; implementation, QA, archive, and release status remain
  owned by the Spec.

Append evidence and routing links instead of rewriting observations. Update
`updated_at` whenever the finding status or evidence changes.
