---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-08
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A failed proof appends a cleanup error the maintainer cannot act on

When an Exact Agent Selection Proof fails, `roundfix profiles validate` and `roundfix doctor` print the real cause correctly and then append a second, unrelated instruction the maintainer has no way to satisfy.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-08-a-failed-proof-appends-a-cleanup-error-the-maintainer-cannot-act-on.md`.
