---
status: deferred
created_at: 2026-08-26
updated_at: 2026-08-26
kind: finding
---

# The Daemon has a green-tree precondition and no postcondition, so a Task settles Clean while breaking the tree

A Task's scoped Verification passed and the Task settled Clean; in the same
tree state the Supervisor's `make verify` exited 2 with 14 typecheck errors born
from the change (baseline the same day: main exits 0). The Agent classified them
as pre-existing and deferred them. The Daemon checks tree health on entry and
never on exit, so a Task can hand the next Task a broken tree with a Clean
verdict on it.

Measured in fluxus, Spec 0052, task_01, 2026-08-20. Composes with the fiscus
observation that green-on-entry inherits the declared command's cache.

Source: secondbrain `inbox/roundfix/2026-08-25-o-daemon-nao-tem-pos-condicao-e-a-task-settla-quebrando-a-arvore.md`
(origin fluxus). Deferred at the 2026-08-26 triage: needs reproduction in this
repository; strongest candidate to join the active queue after it, since a
false Clean is rework material by definition.
