---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-14
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A Spec archives before continuous integration answers

The mandatory order per Spec is: implement the graph with its authored gate, archive, open the Pull Request, watch until Clean, and merge. Archiving therefore happens **before** any answer from continuous integration, so a CI failure arrives when the Spec is already history.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-14-a-spec-archives-before-continuous-integration-answers.md`.
