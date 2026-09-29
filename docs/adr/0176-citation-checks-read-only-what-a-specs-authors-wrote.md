---
status: accepted
created_at: 2026-09-29T00:00:00Z
updated_at: 2026-09-29T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Citation checks read only what a Spec's authors wrote

`SC-ADR-UNLISTED` requires a Spec's `_prd.md` to list every ADR the Spec
cites (ADR-0093). It read every Markdown file under the Spec folder. That
included two kinds of text no Spec author writes:

- the `## Result` an implementing Agent appends to its Task, and the sections
  the Daemon writes: `## Recorded paths` (ADR-0166) and
  `## Carry-forward provenance`;
- the QA reports and evidence the gate writes under `qa/`.

On 2026-09-29, Spec 0181's QA gate refused before measuring any row. Task 03's
Result named ADR-0002 as the number of a test fixture. The Agent cannot edit
the PRD, and the refused QA report quoted the finding, so the refusal also
renewed itself.

The citation walk now reads only the authored projection: `_idea.md`,
`_prd.md`, `_techspec.md`, `_tasks.md`, `_authorization.md`, the Task files
without their Agent- and Daemon-owned sections, and `references/`. It skips
everything under `qa/`. A citation in authored text is reported exactly as
before.

## Consequences

An Agent's evidence and the gate's own reports can name any decision without
creating an obligation that only the Spec's author could discharge, and that
the author did not create. This follows the rule that a check must fail only
where someone can act on it, which ADR-0168 applied to `SC-ADR-RELATED`.
`SC-CITATION-UNSUPPORTED` reads PRD claims only and is unchanged. The
adoption of `references/` is an authorial act, so adopted sources still count.
