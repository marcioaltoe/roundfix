---
type: fix
status: open
created: 2026-11-26
spec: null
reason: null
---

# `roundfix review dispose` refuses with outcome "reviewed" right after a review parked the item as review-findings

## Symptom

Right after a review that parked the item as review-findings, `roundfix review dispose` refuses with `pre-pr review record outcome is "reviewed", not findings`. It happened on three Vortex items (0213, 0218, 0222). A retry with a changed head reviews again, but the disposition stays unrecorded.

## Expected

The record the review writes when it parks the item as review-findings carries the outcome that `review dispose` accepts, or the refusal names how to record the disposition.

## Source

Cross-project message `01a2168c-7b8e-7580-9d8e-8f08971b2f10` from the Vortex project, 2026-11-26. Nature: reported (three items, Roundfix 0.153.3, no response expected). Recorded by triage on 2026-11-26.

## Addendum 2026-11-26: the refusal points at the wrong checkout

The reporter corrected the report the same day: `roundfix review dispose` refused only because it ran from the canonical checkout, where it reads main's review record, already reviewed. Run from the Delivery Queue item's worktree, it recorded the dispositions normally. The remaining defect is the refusal: when the slug passed in `--spec` has a Delivery Queue item, the refusal should name that item's worktree instead of only `outcome is "reviewed", not findings`.

Source: cross-project message `01a21794-19ab-77d2-a371-2ee2430f33e4` from the Vortex project, 2026-11-26 08:20 BRT. Nature: reported.
