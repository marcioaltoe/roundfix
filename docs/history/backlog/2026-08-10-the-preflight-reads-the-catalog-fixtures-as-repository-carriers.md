---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# The preflight reads the catalog fixtures as repository carriers

Baseline preflight inventories every instruction carrier under the repository and warns when one nests managed markers inside unmanaged bytes. It applies that rule to the Baseline's own embedded assets. Measured on 2026-08-10, both `roundfix baseline plan` and `roundfix baseline apply` emitted:

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-10-the-preflight-reads-the-catalog-fixtures-as-repository-carriers.md`.
