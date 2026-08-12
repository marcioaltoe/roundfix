---
status: accepted
created_at: 2026-07-15T16:53:06Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# A missing Review Source check ends the Run Clean Unverified, never Clean

After the Final Push, watch polls for the Review Source's status check on the pushed head through a grace period (the same settle-wait mechanism the fetch step uses); if the check never appears, the Run ends Clean Unverified — a distinct terminal outcome with its own exit code — instead of Clean with a stderr note. This refines ADR-0019: Clean keeps meaning affirmatively Merge-Ready, and a caller (human or script) can distinguish "verified merge-ready" from "pushed, but the Review Source never looked". Observed motivation: a dogfood watch declared Clean immediately after the Final Push while CodeRabbit had not yet begun re-analyzing the new head, and script callers could not tell the two Cleans apart.
