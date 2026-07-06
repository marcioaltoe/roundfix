---
spec: 0005-tui-cockpit
status: active
created: 2026-07-05
surfaces: [frontend, cli]
---

# TUI Cockpit Redesign

The Live Run View grew organically: the review path got the full cockpit, the
Implement Command received a functional but minimal Task pane, and issue
detail occupies layout instead of appearing on demand. A terminal UI
redesign was already drafted and mocked before the Implement Command existed
(see `design/ui-redesign-plan.md` and the two reference mockups in
`design/`); this Spec adopts that direction and generalizes it: one
two-surface cockpit for both Run kinds, with Work Items on the left, the
session timeline as the dominant surface on the right, a phase row stating
where the Run is, and detail as a centered terminal modal.

## Goals

- One cockpit serves review and spec Runs alike: a Work Queue of the Run's
  Work Items and a dominant session timeline, matching the approved visual
  contract.
- The Run's position is always visible: a phase row with terminal-text status
  markers (`[done]`, `[run]`, `[wait]`, `[locked]`) appropriate to the Run
  kind.
- Any Work Item's detail — Review Issue artifact or task file — opens as a
  centered modal over the dimmed cockpit and closes back into unchanged
  context.
- The UI stays a terminal tool: monospace, borders, text markers; no icons,
  logos, or web-dashboard elements.
- Attach renders the identical cockpit from the Run Event Journal replay.

## User Stories

1. As a developer watching a review Run, I want the Work Queue grouped by
   Batch with per-item status markers and the timeline as the largest
   surface, so that I can scan progress without losing the live stream.
2. As a developer watching a spec Run, I want Tasks in the same cockpit with
   a spec-appropriate phase row, so that implement Runs are first-class in
   the Live Run View, not a reduced variant.
3. As a developer inspecting one Work Item, I want Enter to open its detail
   as a centered modal — Review Issue artifact or task file body — scrollable
   and dismissable back to the exact same context, so that inspection never
   costs me my place.
4. As a developer attaching to a finished or live Run, I want the same
   cockpit replayed from the Journal, so that live and attached views are one
   experience.
5. As a developer scrolling back through the timeline, I want Follow Mode to
   suspend and resume exactly as today, so that redesign does not change
   stream semantics.
6. As a developer on a small terminal, I want the layout to degrade to a
   usable single-surface fallback, so that the cockpit never renders broken.

## Core Features

1. **Two persistent surfaces.** Work Queue left, session timeline right; the
   timeline is the widest surface; no persistent third pane.
2. **Phase row.** Review Runs: fetch → triage → agent → verify → push; spec
   Runs: agent → verify → commit → QA (QA `[locked]` until every Task is
   completed, absent without the QA opt-in). Markers are text only.
3. **Work Queue rows.** Grouped by Batch with separators and elapsed time;
   each row shows a status marker, severity when the Work Item carries one,
   ordinal, title, and location; the queue footer keeps the compact totals.
4. **Timeline grouping.** Events grouped per Batch and kind (plan, tool,
   thought, status, daemon milestones) at render time, preserving raw
   payloads and Follow Mode behavior.
5. **Modal detail.** Enter opens the selected Work Item's detail centered
   over the dimmed cockpit; `D` toggles it; `Esc` closes; `j/k`, arrows, and
   paging scroll; footer hints change per state. Review Runs show the Review
   Issue artifact; spec Runs show the task file body, read-only.
6. **Attach parity and journal discipline.** The cockpit reads only the Run
   Event Journal and re-read artifacts/task files, live or attached
   (ADR-0009); unknown event kinds stay skippable.
7. **Non-TTY untouched.** The plain-text renderer and stderr behavior do not
   change in this Spec.

## User Experience

The reference mockups in `design/` are the visual contract: normal state
(queue + expanded timeline) and detail state (same screen dimmed behind a
centered modal). Keys: `Tab` focus, arrows/`j/k` move and scroll, `PgUp/PgDn`
page, `Enter` detail, `D` toggle detail, `Esc` close, `End` follow, `Ctrl-C`
stop. Small terminals collapse to the timeline with a one-line queue summary.

## Non-Goals / Out of Scope

- No web dashboard, provider logos, pictographic icons, or decorative
  elements.
- No persistent third detail pane.
- No new review workflow, no source-thread mutation from the UI.
- No changes to Run Event storage, journal schema, or agent contracts.
- No non-TTY/stderr changes (shipped in the watch Spec).
- No CLI framework migration.

## Success Metrics

- A review Run and a spec Run render the same cockpit structure with their
  respective phase rows, validated by synchronous model tests.
- Modal open/close returns to identical queue and timeline state
  (byte-compared render before and after).
- Attach over a finished Run of each kind renders the same surfaces as the
  live cockpit did.
- Existing stop, attach, close, scroll, and follow semantics pass their
  current tests unchanged.

## Decisions

- `Enter` opens the modal for the selected Work Item and `D` toggles it —
  both, resolving the plan's open decision; `Esc` closes.
- The phase row replaces the old progress header; the compact totals live in
  the Work Queue footer (as mocked).
- Timeline grouping is render-time only; daemon events are not enriched for
  the UI (ADR-0008 stays untouched).
- Dimming behind the modal is a light shade — readability beats effect.
- The recovered redesign plan and mockups in `design/` are the binding visual
  contract; deviations require updating them in the same change.

## Open Questions

None.
