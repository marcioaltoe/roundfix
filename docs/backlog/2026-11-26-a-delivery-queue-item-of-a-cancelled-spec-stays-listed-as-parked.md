---
type: fix
status: open
created: 2026-11-26
spec: null
reason: null
---

# A Delivery Queue item of a cancelled Spec stays listed as parked

## Symptom

Indexus cancelled its Spec 0081 (pull request #297; the record is under `docs/history/specs/0081-…`). Roundfix already treats it as cancelled ("a cancelled Spec is never delivered") and refuses a retry, but `roundfix deliver status` keeps listing the item as `parked` (`qa-environment-partial`) until the next `roundfix deliver start` records a new queue. No command removes it: `deliver` offers only `plan`, `start`, `status`, `resume`, `retry` and `stop`. Whoever reads the status, the coordinating session's sweep included, sees a parked item that no longer exists and has to silence it by hand.

## Expected

One of the two forms the reporter proposed:

- `roundfix deliver retire <slug>`, accepted only for an item whose Spec is cancelled or archived, touching no worktree or branch beyond what the cancellation already did; or
- `roundfix deliver status` shows that item as `cancelled` instead of `parked`.

The Spec that takes this entry chooses the form.

## Source

Cross-project message `01a2178a-f2de-7992-aa5b-aee7bac7cecf` from the Praetor (the session that coordinates the projects), relaying the Indexus case, 2026-11-26 ~08:15 BRT. Nature: reported, no urgency. Recorded by triage on 2026-11-26.
