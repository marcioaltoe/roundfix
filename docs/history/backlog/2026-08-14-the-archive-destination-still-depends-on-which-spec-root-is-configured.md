---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The archive destination still depends on which Spec Root is configured

`roundfix archive` sends a completed Spec to `docs/history/specs/<slug>` when the Spec Root is the built-in `docs/specs`, and to `<spec-root>/_archived/<slug>` when it is anything else. Spec 0094 moved the built-in destination and left the asymmetry: the default root is still the only configuration whose archive lives outside it, and a non-default root still archives under the pre-0094 name.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-14-the-archive-destination-still-depends-on-which-spec-root-is-configured.md`.
