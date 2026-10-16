---
type: fix
status: promoted
created: 2026-10-14
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
reason: null
---

# Deliver status does not say the queue owner is gone

## Symptom

The machine rebooted at 19:57 BRT on 2026-10-13 and the Delivery Owner
(`roundfix deliver resume`) died with it. `roundfix deliver status` kept
printing the 0266 item as `checking`, with no warning, for about 11 hours,
while its pull request was green. `roundfix deliver resume` recovered it at
once.

## Where

`roundfix deliver status` and the Delivery Owner record under the Roundfix
home.

## Expected

`deliver status` says when the recorded Delivery Owner is not running and
names `roundfix deliver resume`. Optionally, a login or launchd hook resumes
an unfinished queue.

## Evidence

Operator log of cycle AG, morning of 2026-10-14: no `roundfix deliver`
process was running, the item read `checking`, and after `deliver resume`
the item merged as #527.
