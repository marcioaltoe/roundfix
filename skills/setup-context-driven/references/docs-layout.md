# Docs layout

<!-- Seeded by the setup-context-driven skill. Edit repo-specific notes freely;
     a re-run regenerates this file carrying confirmed answers forward. -->

How this repository uses each `docs/` folder. Every folder has one job; a file
that fits two folders goes to the one whose job it serves now — move it
when its job changes (inbox → findings → spec is the normal flow).

| Folder | Job | Lifecycle |
| --- | --- | --- |
| `docs/_inbox/` | Raw incoming notes: pasted reports, half-formed ideas, unprocessed field notes. Nothing here is triaged or trustworthy yet. | Triage each item: promote to `findings/`, `references/`, or a spec; then remove it from the inbox. An empty inbox is the healthy state. |
| `docs/adr/` | Accepted decision records — `NNNN-kebab-slug.md`, 1–3 sentences each (context, decision, why). One numbering sequence for the repo's life. | Append-only. Numbers are never reused; superseding decisions name what they supersede. |
| `docs/agents/` | Agent-facing usage guides: the files seeded by `setup-context-driven` plus repo-authored guides. `AGENTS.md`/`CLAUDE.md` hold only short pointers here, never rule bodies. | Seeded files are owned by the skill and regenerated on re-run; repo-authored guides are owned by the repo. |
| `docs/design/` | Design artifacts: mockups, visual and interaction decisions, UI/TUI explorations, design-review notes. | Kept while the design is live; superseded explorations may be pruned or archived into the spec that consumed them. |
| `docs/findings/` | Dated field reports: dogfood incidents, retrospectives, root-cause investigations. The raw material the spec pipeline consumes. Follow the template below. | Immutable history with addenda: append root causes and spec pointers as they land; never rewrite what was observed. Finding status tracks routing, not Spec implementation progress. |
| `docs/handoffs/` | Session handoff documents: the state snapshot one working session leaves for the next. | Superseded by the next handoff; keep the recent few, prune the rest. |
| `docs/references/` | Pointers to external resources, each with a one-line explanation of why it matters here. | Prune links that stop mattering. |
| `docs/specs/` | The spec workflow tree: `NNNN-<slug>/` feature folders, `_archived/` for shipped specs, `_reviews/` for review-run artifacts. | Owned by the pipeline skills; status lives only in task files. |
| `docs/user-guide/` | Human-facing product documentation and runbooks. | Updated with the behavior it documents in the same PR. |

## Findings template

Findings files are named `YYYY-MM-DD-<kebab-slug>.md`. One file records one
session or investigation.

```yaml
---
status: pending # pending | partial | deferred | done
created_at: YYYY-MM-DD
updated_at: YYYY-MM-DD
---
```

- `pending`: the finding is new and does not yet have an implementation Spec.
- `partial`: a Spec covers only the portion selected for implementation because
  the remaining observations were judged unnecessary. Record that reason and
  link the Spec.
- `deferred`: the finding will not be implemented. Record the reason.
- `done`: an implementation Spec exists. Set this status as soon as the Spec is
  created and linked; do not wait for its Tasks, QA, archive, or release.
- `created_at` records when the document was created.
- `updated_at` records the latest status change or evidence addendum.

```markdown
# <Area> — <short title> (YYYY-MM-DD)

<!-- 2–4 sentences: session or Run context, what was being attempted,
     and pointers to adjacent findings files instead of duplicated content. -->

## 1. <Finding title — the symptom, not the guess>

- **Symptom / evidence**: what was observed, verbatim where possible
  (commands, output, run ids, file paths).
- **Root cause**: when established — how it was proven; otherwise say
  "unknown" and what was ruled out.
- **Action / suggestion**: the fix, workaround applied, or the route
  to a Spec, direct fix, or upstream report. Link the Spec or ADR once it
  exists.

## 2. <Next finding…>

## What worked — keep

<!-- Optional: behaviors that held up under stress, worth preserving. -->
```

Conventions:

- Evidence over narrative: quote the command and output; name Run IDs, commits,
  and paths so a later session can re-verify.
- Append dated root-cause addenda instead of rewriting the original observation.
- When a finding gets an implementation Spec, add the pointer and set
  `status: done` in the same change. Use `partial` instead only for an
  intentionally limited Spec and record why the remaining scope is
  unnecessary.
- Use `deferred` only with an explicit reason for not implementing the finding.
- Update `status` and `updated_at` whenever the document's triage state changes.
