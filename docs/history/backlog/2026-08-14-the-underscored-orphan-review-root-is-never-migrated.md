---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The underscored orphan review root is never migrated

Spec 0094 moved the live orphan Review Artifact root from `docs/specs/_reviews/` to `docs/specs/reviews/` for new writes, and nothing migrates the folders already sitting at the underscored path. Layout discovery evaluates them for retirement, never for the rename.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-14-the-underscored-orphan-review-root-is-never-migrated.md`.
