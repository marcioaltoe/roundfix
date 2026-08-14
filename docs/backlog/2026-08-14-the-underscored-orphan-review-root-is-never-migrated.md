---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The underscored orphan review root is never migrated

## Symptom

Spec 0094 moved the live orphan Review Artifact root from `docs/specs/_reviews/`
to `docs/specs/reviews/` for new writes, and nothing migrates the folders already
sitting at the underscored path. Layout discovery evaluates them for retirement,
never for the rename.

Measured in `roundfix` on 2026-08-14: fifty folders held 501 files at the
underscored path after 0094 shipped. Two of the three states the liveness check
can produce leave them there permanently — an undecidable head and a live head
are both retained — so a repository whose reviews never decide keeps writing to
one path while fifty folders sit at another.

They migrated here only because the local object store happened to hold their
recorded heads, which is the accident recorded in
`docs/findings/2026-08-14-a-review-retires-on-whatever-the-object-store-happens-to-hold.md`.
A rename should not have depended on that.

## Where

`DiscoverHistoryLayout` in `internal/baseline`, which enumerates the legacy
shapes a migration recognises. The underscored orphan review root is not among
them.

## Expected

`docs/specs/_reviews/` is a legacy layout the migration recognises and renames to
the live root, independent of any judgement about whether a review finished. A
rename needs no liveness answer; retirement is a separate question the same
discovery can ask afterwards.

## Evidence

`docs/findings/2026-08-14-a-review-retires-on-whatever-the-object-store-happens-to-hold.md`,
finding 3, which measured the gap while running the shipped migration against the
repository that shipped it.
