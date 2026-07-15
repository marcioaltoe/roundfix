# Docs layout

<!-- Seeded by the setup-context-driven skill. Edit repo-specific notes freely;
     a re-run regenerates this file carrying confirmed answers forward. -->

How this repository uses each `docs/` folder. Every folder has one job; a file
that fits two folders goes to the one whose job it serves **now** — move it
when its job changes (inbox → findings → spec is the normal flow).

| Folder | Job | Lifecycle |
| --- | --- | --- |
| `docs/_inbox/` | Raw incoming notes: pasted reports, half-formed ideas, unprocessed field notes. Nothing here is triaged or trustworthy yet. | Triage each item: promote to `findings/`, `references/`, or a spec; then remove it from the inbox. An empty inbox is the healthy state. |
| `docs/adr/` | Accepted decision records — `NNNN-kebab-slug.md`, 1–3 sentences each (context, decision, why). One numbering sequence for the repo's life. | Append-only. Numbers are never reused; superseding decisions name what they supersede. |
| `docs/agents/` | Agent-facing usage guides: the files seeded by `setup-context-driven` plus repo-authored guides. `AGENTS.md`/`CLAUDE.md` hold only short pointers here, never rule bodies. | Seeded files are owned by the skill and regenerated on re-run; repo-authored guides are owned by the repo. |
| `docs/design/` | Design artifacts: mockups, visual and interaction decisions, UI/TUI explorations, design-review notes. | Kept while the design is live; superseded explorations may be pruned or archived into the spec that consumed them. |
| `docs/findings/` | Dated field reports: dogfood incidents, retrospectives, root-cause investigations. The raw material the spec pipeline consumes. Follow the template below. | Immutable history with addenda: append root causes and spec pointers as they land; never rewrite what was observed. |
| `docs/handoffs/` | Session handoff documents: the state snapshot one working session leaves for the next (what was done, what is in flight, exact next actions). | Superseded by the next handoff; keep the recent few, prune the rest. |
| `docs/references/` | Pointers to external resources — upstream docs, dashboards, tickets, papers — each with a one-line "why it matters here". | Prune links that stop mattering; a dead link with no why-line is noise. |
| `docs/specs/` | The spec workflow tree: `NNNN-<slug>/` feature folders, `_archived/` for shipped specs, `_reviews/` for review-run artifacts. Conventions live in the issue-tracker guide. | Owned by the pipeline skills (`write-prd` → … → `archive-spec`); status lives only in task files. |
| `docs/user-guide/` | Human-facing product documentation: usage guides, runbooks, method explainers (for example `context-driven-development.md`). Shipped documentation — written for users, not agents. | Updated with the behavior it documents, same-PR (like the skill-sync rule). |

## Findings template

Findings files are named `YYYY-MM-DD-<kebab-slug>.md` (absolute dates, never
"today"). One file per session or investigation; number the findings inside it.

```markdown
# <Area> — <short title> (YYYY-MM-DD)

<!-- 2-4 sentences: session or run context, what was being attempted,
     and pointers to adjacent findings files instead of duplicated content. -->

## 1. <Finding title — the symptom, not the guess>

- **Symptom / evidence**: what was observed, verbatim where possible
  (commands, output, run ids, file paths).
- **Root cause**: when established — how it was proven; otherwise say
  "unknown" and what was ruled out.
- **Action / suggestion**: the fix, workaround applied, or the route
  (spec, direct fix, upstream report). Link the spec or ADR once it exists.

## 2. <Next finding…>

## What worked — keep

<!-- Optional: behaviors that held up under stress, worth preserving. -->
```

Conventions:

- Evidence over narrative: quote the command and its output; name run ids,
  commits, and paths so a later session can re-verify.
- Root-cause addenda are appended to the original finding (marked with their
  date), never rewritten over the observation.
- When a finding becomes a spec, add the pointer in place — the findings file
  is where future sessions look first.
